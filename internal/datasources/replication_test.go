package datasources

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testReplicationID   = "cmVwbC0x"
	testReplicationPath = "/v4/organizations/org-1/projects/proj-1/clusters/cluster-1/replications/" + testReplicationID
)

// testReplicationGetBody is a GET response as the API sends it, with all optional parts set.
const testReplicationGetBody = `{
	"id": "cmVwbC0x",
	"status": "running",
	"changesLeft": 5,
	"source": {
		"project": {"id": "proj-1", "name": "project"},
		"cluster": {"id": "cluster-1", "name": "source"},
		"bucket": {"id": "c291cmNl", "name": "source", "conflictResolutionType": "seqno"}
	},
	"target": {
		"project": {"id": "proj-1", "name": "project"},
		"cluster": {"id": "cluster-2", "name": "target"},
		"bucket": {"id": "dGFyZ2V0", "name": "target", "conflictResolutionType": "seqno"},
		"type": "capella"
	},
	"mappings": [
		{"sourceScope": "sales", "targetScope": "sales", "collections": [{"sourceCollection": "orders", "targetCollection": "orders"}]}
	],
	"direction": "oneWay",
	"priority": "medium",
	"networkUsageLimit": 50,
	"filter": {
		"expressions": {"regEx": "^a"},
		"documentExcludeOptions": {"deletion": false, "expiration": false, "ttl": true, "binary": false}
	},
	"mobileSyncMode": "active",
	"conflictLogging": {
		"disabled": false,
		"bucket": "bG9ncw==",
		"collection": "logs.default",
		"loggingRules": {"sales": {}, "sales.orders": null}
	},
	"audit": {"createdBy": "user-1", "createdAt": "2026-10-01T00:00:00Z"}
}`

func newReplicationTestData(t *testing.T, status int, body string) *providerschema.Data {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != testReplicationPath {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	return &providerschema.Data{
		ClientV1: &api.Client{Client: srv.Client()},
		HostURL:  srv.URL,
	}
}

// replicationConfig is a data source config with only the four required IDs set.
func replicationConfig(t *testing.T, ctx context.Context) tfsdk.Config {
	t.Helper()
	s := ReplicationSchema()
	state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
	diags := state.Set(ctx, providerschema.ReplicationData{
		Id:             types.StringValue(testReplicationID),
		OrganizationId: types.StringValue("org-1"),
		ProjectId:      types.StringValue("proj-1"),
		ClusterId:      types.StringValue("cluster-1"),
		Mappings:       types.SetNull(types.ObjectType{AttrTypes: providerschema.ReplicationScopeMapping{}.AttributeTypes()}),
		SourceDetails:  types.ObjectNull(providerschema.ReplicationSourceDetailsAttributeTypes()),
		TargetDetails:  types.ObjectNull(providerschema.ReplicationTargetDetailsAttributeTypes()),
		Audit:          types.ObjectNull(providerschema.ReplicationAudit{}.AttributeTypes()),
	})
	require.False(t, diags.HasError(), "%v", diags)
	return tfsdk.Config{Schema: s, Raw: state.Raw}
}

func readReplication(t *testing.T, data *providerschema.Data) *datasource.ReadResponse {
	t.Helper()
	ctx := context.Background()
	s := ReplicationSchema()
	resp := &datasource.ReadResponse{State: tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}}
	(&Replication{Data: data}).Read(ctx, datasource.ReadRequest{Config: replicationConfig(t, ctx)}, resp)
	return resp
}

func apiErrorBody(t *testing.T, status int, message string) string {
	t.Helper()
	body, err := json.Marshal(map[string]any{"code": 2, "httpStatusCode": status, "message": message})
	require.NoError(t, err)
	return string(body)
}

func TestReplicationRead(t *testing.T) {
	t.Run("writes all the attributes that the API returns", func(t *testing.T) {
		resp := readReplication(t, newReplicationTestData(t, http.StatusOK, testReplicationGetBody))
		require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())

		var got providerschema.ReplicationData
		require.False(t, resp.State.Get(context.Background(), &got).HasError())

		assert.Equal(t, testReplicationID, got.Id.ValueString())
		assert.Equal(t, "c291cmNl", got.SourceBucket.ValueString())
		assert.Equal(t, &providerschema.ReplicationTarget{
			Cluster: types.StringValue("cluster-2"),
			Bucket:  types.StringValue("dGFyZ2V0"),
			Type:    types.StringValue("capella"),
		}, got.Target)
		assert.Equal(t, "oneWay", got.Direction.ValueString())
		assert.Equal(t, "medium", got.Priority.ValueString())
		assert.Equal(t, int64(50), got.NetworkUsageLimit.ValueInt64())
		assert.Len(t, got.Mappings.Elements(), 1)
		assert.Equal(t, &providerschema.ReplicationDataFilter{
			DocumentExcludeOptions: &providerschema.ReplicationDocumentExcludeOptions{
				Deletion:   types.BoolValue(false),
				Expiration: types.BoolValue(false),
				TTL:        types.BoolValue(true),
				Binary:     types.BoolValue(false),
			},
			Expressions: &providerschema.ReplicationDataFilterExpressions{RegEx: types.StringValue("^a")},
		}, got.Filter)
		assert.Equal(t, "active", got.MobileSyncMode.ValueString())
		require.NotNil(t, got.ConflictLogging)
		assert.Equal(t, "bG9ncw==", got.ConflictLogging.Bucket.ValueString())
		assert.True(t, got.ConflictLogging.LoggingRules.Elements()["sales.orders"].IsNull(), "a null rule stays null")
		assert.Equal(t, "running", got.Status.ValueString())
		assert.Equal(t, int64(5), got.ChangesLeft.ValueInt64())
		assert.True(t, got.Error.IsNull())
		assert.False(t, got.SourceDetails.IsNull())
		assert.False(t, got.TargetDetails.IsNull())
		assert.False(t, got.Audit.IsNull())
	})

	t.Run("a missing replication is an error that explains which cluster to use", func(t *testing.T) {
		resp := readReplication(t, newReplicationTestData(t, http.StatusNotFound,
			apiErrorBody(t, http.StatusNotFound, "The replication was not found.")))

		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "reverse replication ID and the target cluster ID")
	})

	t.Run("an unhealthy source cluster is an error with the API message", func(t *testing.T) {
		resp := readReplication(t, newReplicationTestData(t, http.StatusUnprocessableEntity,
			apiErrorBody(t, http.StatusUnprocessableEntity, "Cluster cluster-1 must be in a healthy state to create or manage replications.")))

		require.True(t, resp.Diagnostics.HasError())
		detail := resp.Diagnostics.Errors()[0].Detail()
		assert.Contains(t, detail, "must be in a healthy state")
		assert.NotContains(t, detail, "reverse replication ID")
	})
}
