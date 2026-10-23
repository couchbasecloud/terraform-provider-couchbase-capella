package resources

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestValidateClusterUpdateExpressScaling checks that express_scaling is create-only on the
// cluster resource: leaving it out of configuration never blocks an update, and only an
// explicit change is rejected, pointing at the dedicated resource.
func TestValidateClusterUpdateExpressScaling(t *testing.T) {
	tests := []struct {
		name    string
		plan    types.Bool
		state   types.Bool
		wantErr error
	}{
		{name: "omitted from configuration", plan: types.BoolUnknown(), state: types.BoolValue(true)},
		{name: "null in plan", plan: types.BoolNull(), state: types.BoolValue(true)},
		{name: "unchanged true", plan: types.BoolValue(true), state: types.BoolValue(true)},
		{name: "unchanged false", plan: types.BoolValue(false), state: types.BoolValue(false)},
		{name: "enable after create", plan: types.BoolValue(true), state: types.BoolValue(false), wantErr: errors.ErrUnableToUpdateExpressScaling},
		{name: "disable after create", plan: types.BoolValue(false), state: types.BoolValue(true), wantErr: errors.ErrUnableToUpdateExpressScaling},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := awsClusterPlan(types.Int64Value(50), types.Int64Value(3000))
			state := awsClusterPlan(types.Int64Value(50), types.Int64Value(3000))
			plan.ExpressScaling = tt.plan
			state.ExpressScaling = tt.state

			err := (&Cluster{}).validateClusterUpdate(plan, state)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

// expressScalingBackend answers the create POST and every cluster GET, recording the
// expressScaling field of the create body so a test can tell an omitted field from false.
type expressScalingBackend struct {
	mu        sync.Mutex
	state     clusterapi.ExpressScalingState
	sent      json.RawMessage
	sentField bool
}

func (b *expressScalingBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	id := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	version := "7.6.4"

	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		var fields map[string]json.RawMessage
		_ = json.Unmarshal(body, &fields)
		b.sent, b.sentField = fields["expressScaling"]

		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(clusterapi.GetClusterResponse{Id: id})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(clusterapi.GetClusterResponse{
		Id:              id,
		CurrentState:    clusterapi.Healthy,
		ExpressScaling:  b.state,
		CloudProvider:   clusterapi.CloudProvider{Type: clusterapi.Aws},
		CouchbaseServer: clusterapi.CouchbaseServer{Version: &version},
	})
}

// TestClusterCreateExpressScaling drives Create against a local backend. It checks that
// express_scaling reaches the POST only when configured, and that state keeps the requested
// bool while express_scaling_state reports what the API returned. The first GET after create
// can still say disabled, because the enable rebalance only starts once the cluster is healthy.
func TestClusterCreateExpressScaling(t *testing.T) {
	fastClusterPolling(t)

	tests := []struct {
		name          string
		plan          types.Bool
		apiState      clusterapi.ExpressScalingState
		wantSent      string
		wantSentField bool
		wantEnabled   bool
	}{
		{name: "true while the first read still says disabled", plan: types.BoolValue(true), apiState: clusterapi.ExpressScalingDisabled, wantSent: "true", wantSentField: true, wantEnabled: true},
		{name: "true and enabling", plan: types.BoolValue(true), apiState: clusterapi.ExpressScalingEnabling, wantSent: "true", wantSentField: true, wantEnabled: true},
		{name: "explicit false", plan: types.BoolValue(false), apiState: clusterapi.ExpressScalingDisabled, wantSent: "false", wantSentField: true, wantEnabled: false},
		{name: "omitted from configuration", plan: types.BoolUnknown(), apiState: clusterapi.ExpressScalingDisabled, wantSentField: false, wantEnabled: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			backend := &expressScalingBackend{state: tt.apiState}
			c := newTestCluster(t, backend.handler)

			plan := awsClusterPlan(types.Int64Value(50), types.Int64Value(3000))
			plan.Availability = &providerschema.Availability{Type: types.StringValue("multi")}
			plan.Support = &providerschema.Support{Plan: types.StringValue("enterprise"), Timezone: types.StringValue("PT")}
			plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
			plan.CouchbaseServer = types.ObjectNull(providerschema.CouchbaseServer{}.AttributeTypes())
			plan.ExpressScaling = tt.plan

			s := ClusterSchema()
			req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: s}}
			require.False(t, req.Plan.Set(ctx, plan).HasError())
			resp := &resource.CreateResponse{State: tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}}

			c.Create(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())

			backend.mu.Lock()
			assert.Equal(t, tt.wantSentField, backend.sentField)
			if tt.wantSentField {
				assert.JSONEq(t, tt.wantSent, string(backend.sent))
			}
			backend.mu.Unlock()

			var got providerschema.Cluster
			require.False(t, resp.State.Get(ctx, &got).HasError())
			assert.Equal(t, types.BoolValue(tt.wantEnabled), got.ExpressScaling)
			assert.Equal(t, types.StringValue(string(tt.apiState)), got.ExpressScalingState)
		})
	}
}
