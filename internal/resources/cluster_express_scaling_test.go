package resources

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testExpressScalingOrgId     = "11111111-1111-1111-1111-111111111111"
	testExpressScalingProjectId = "22222222-2222-2222-2222-222222222222"
	testExpressScalingClusterId = "33333333-3333-3333-3333-333333333333"
)

// expressScalingEndpoint answers the cluster GET with a scripted sequence of states, holding
// the last one, and the expressScaling PUT with a fixed status and error body. GETs after the
// PUT can fail with getStatusAfterPut, or stall for getDelayAfterPut.
type expressScalingEndpoint struct {
	mu                sync.Mutex
	states            []clusterapi.ExpressScalingState
	getStatus         int
	getStatusAfterPut int
	getDelayAfterPut  time.Duration
	release           chan struct{}
	putStatus         int
	putError          *api.Error
	puts              []json.RawMessage
}

func (e *expressScalingEndpoint) handler(w http.ResponseWriter, r *http.Request) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if r.Method == http.MethodPut {
		body, _ := io.ReadAll(r.Body)
		var fields map[string]json.RawMessage
		_ = json.Unmarshal(body, &fields)
		e.puts = append(e.puts, fields["expressScaling"])

		if e.putError != nil {
			w.WriteHeader(e.putError.HttpStatusCode)
			_ = json.NewEncoder(w).Encode(e.putError)
			return
		}
		w.WriteHeader(e.putStatus)
		return
	}

	status := e.getStatus
	if len(e.puts) > 0 && e.getStatusAfterPut != 0 {
		status = e.getStatusAfterPut
	}
	if status != 0 && status != http.StatusOK {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(api.Error{HttpStatusCode: status, Code: status, Message: http.StatusText(status)})
		return
	}

	if len(e.puts) > 0 && e.getDelayAfterPut > 0 {
		e.mu.Unlock()
		select {
		case <-r.Context().Done():
		case <-e.release:
		case <-time.After(e.getDelayAfterPut):
		}
		e.mu.Lock()
	}

	state := e.states[0]
	if len(e.states) > 1 {
		e.states = e.states[1:]
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(clusterapi.GetClusterResponse{ExpressScaling: state})
}

func (e *expressScalingEndpoint) putCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.puts)
}

func newTestClusterExpressScaling(t *testing.T, e *expressScalingEndpoint) *ClusterExpressScaling {
	t.Helper()

	orig := expressScalingPollInterval
	expressScalingPollInterval = time.Millisecond
	t.Cleanup(func() { expressScalingPollInterval = orig })

	srv := httptest.NewServer(http.HandlerFunc(e.handler))
	t.Cleanup(srv.Close)

	return &ClusterExpressScaling{
		Data: &providerschema.Data{
			ClientV1: &api.Client{Client: srv.Client()},
			HostURL:  srv.URL,
		},
	}
}

func expressScalingPlan(t *testing.T, enabled bool, timeout string) providerschema.ClusterExpressScaling {
	t.Helper()

	timeoutAttrs := map[string]attr.Value{
		"create": types.StringValue(timeout),
		"update": types.StringValue(timeout),
	}
	obj, diags := types.ObjectValue(map[string]attr.Type{"create": types.StringType, "update": types.StringType}, timeoutAttrs)
	require.False(t, diags.HasError())

	return providerschema.ClusterExpressScaling{
		OrganizationId:      types.StringValue(testExpressScalingOrgId),
		ProjectId:           types.StringValue(testExpressScalingProjectId),
		ClusterId:           types.StringValue(testExpressScalingClusterId),
		ExpressScaling:      types.BoolValue(enabled),
		ExpressScalingState: types.StringUnknown(),
		Timeouts:            timeouts.Value{Object: obj},
	}
}

func createExpressScaling(t *testing.T, r *ClusterExpressScaling, plan providerschema.ClusterExpressScaling) (*resource.CreateResponse, providerschema.ClusterExpressScaling) {
	t.Helper()
	ctx := t.Context()

	s := ClusterExpressScalingSchema()
	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: s}}
	require.False(t, req.Plan.Set(ctx, plan).HasError())
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}}

	r.Create(ctx, req, resp)

	var got providerschema.ClusterExpressScaling
	if !resp.State.Raw.IsNull() {
		require.False(t, resp.State.Get(ctx, &got).HasError())
	}
	return resp, got
}

func TestClusterExpressScalingCreate(t *testing.T) {
	tests := []struct {
		name        string
		enabled     bool
		states      []clusterapi.ExpressScalingState
		putError    *api.Error
		wantPuts    int
		wantState   clusterapi.ExpressScalingState
		wantWarning bool
		wantError   string
	}{
		{
			name:      "enable waits for the rebalance to finish",
			enabled:   true,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabled, clusterapi.ExpressScalingEnabling, clusterapi.ExpressScalingEnabling, clusterapi.ExpressScalingEnabled},
			wantPuts:  1,
			wantState: clusterapi.ExpressScalingEnabled,
		},
		{
			name:      "disable on a cluster that reports no state is a no-op",
			enabled:   false,
			states:    []clusterapi.ExpressScalingState{""},
			wantPuts:  0,
			wantState: "",
		},
		{
			name:      "already enabled sends no request",
			enabled:   true,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingEnabled},
			wantPuts:  0,
			wantState: clusterapi.ExpressScalingEnabled,
		},
		{
			name:      "rejected while already enabling waits instead of failing",
			enabled:   true,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingEnabling, clusterapi.ExpressScalingEnabled},
			putError:  &api.Error{HttpStatusCode: http.StatusUnprocessableEntity, Code: expressScalingInvalidStateCode, Message: "not healthy"},
			wantPuts:  1,
			wantState: clusterapi.ExpressScalingEnabled,
		},
		{
			name:      "rejected while disabling fails when enabling is planned",
			enabled:   true,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabling},
			putError:  &api.Error{HttpStatusCode: http.StatusUnprocessableEntity, Code: expressScalingInvalidStateCode, Message: "not healthy"},
			wantPuts:  1,
			wantError: "The cluster must be healthy to change express scaling.",
		},
		{
			name:      "enterprise plan required surfaces the API message",
			enabled:   true,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabled},
			putError:  &api.Error{HttpStatusCode: http.StatusUnprocessableEntity, Code: http.StatusUnprocessableEntity, Message: "Enterprise plan is required to enable express scaling."},
			wantPuts:  1,
			wantError: "Enterprise plan is required to enable express scaling.",
		},
		{
			name:      "feature not available returns a not found explanation",
			enabled:   true,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabled},
			putError:  &api.Error{HttpStatusCode: http.StatusNotFound, Code: http.StatusNotFound, Message: "not found"},
			wantPuts:  1,
			wantError: "express scaling is not available for this organization",
		},
		{
			name:      "invalid request asks for a provider bug report",
			enabled:   false,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingEnabled},
			putError:  &api.Error{HttpStatusCode: http.StatusBadRequest, Code: http.StatusBadRequest, Message: "Ensure expressScaling field is provided."},
			wantPuts:  1,
			wantError: "Please report this issue to the provider developers.",
		},
		{
			name:      "rate limit surfaces the API message",
			enabled:   true,
			states:    []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabled},
			putError:  &api.Error{HttpStatusCode: http.StatusUnprocessableEntity, Code: http.StatusUnprocessableEntity, Message: "The maximum number of 3 express scales allowed in a 24 hour period has been reached."},
			wantPuts:  1,
			wantError: "The maximum number of 3 express scales allowed in a 24 hour period has been reached.",
		},
		{
			name:        "timeout saves the planned value with a warning",
			enabled:     true,
			states:      []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabled, clusterapi.ExpressScalingEnabling},
			wantPuts:    1,
			wantState:   clusterapi.ExpressScalingEnabling,
			wantWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &expressScalingEndpoint{states: tt.states, putStatus: http.StatusNoContent, putError: tt.putError}
			r := newTestClusterExpressScaling(t, e)

			timeout := "2s"
			if tt.wantWarning {
				timeout = "50ms"
			}
			resp, got := createExpressScaling(t, r, expressScalingPlan(t, tt.enabled, timeout))

			assert.Equal(t, tt.wantPuts, e.putCount())

			if tt.wantError != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, diagDetails(resp.Diagnostics.Errors()), tt.wantError)
				assert.True(t, resp.State.Raw.IsNull(), "state must not be set when the change was refused")
				return
			}

			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())
			assert.Equal(t, tt.wantWarning, resp.Diagnostics.WarningsCount() > 0)
			assert.Equal(t, types.BoolValue(tt.enabled), got.ExpressScaling)
			assert.Equal(t, types.StringValue(string(tt.wantState)), got.ExpressScalingState)
		})
	}
}

// TestClusterExpressScalingCreatePollingFailures checks that losing access while waiting is
// reported as an error rather than a timeout warning, and that a GET which stalls past the
// timeout is cut off at the deadline rather than holding the apply open.
func TestClusterExpressScalingCreatePollingFailures(t *testing.T) {
	t.Run("forbidden while polling is an error", func(t *testing.T) {
		e := &expressScalingEndpoint{
			states:            []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabled},
			putStatus:         http.StatusNoContent,
			getStatusAfterPut: http.StatusForbidden,
		}
		r := newTestClusterExpressScaling(t, e)

		start := time.Now()
		resp, _ := createExpressScaling(t, r, expressScalingPlan(t, true, "10s"))

		assert.Less(t, time.Since(start), 5*time.Second, "polling must stop on 403 instead of retrying until the timeout")
		require.True(t, resp.Diagnostics.HasError())
		assert.Zero(t, resp.Diagnostics.WarningsCount())
		assert.Contains(t, diagDetails(resp.Diagnostics.Errors()), "could not be read afterwards")
		assert.Contains(t, diagDetails(resp.Diagnostics.Errors()), "403")
	})

	t.Run("stalled GET is cut off at the timeout", func(t *testing.T) {
		e := &expressScalingEndpoint{
			states:           []clusterapi.ExpressScalingState{clusterapi.ExpressScalingDisabled, clusterapi.ExpressScalingEnabled},
			putStatus:        http.StatusNoContent,
			getDelayAfterPut: 30 * time.Second,
			release:          make(chan struct{}),
		}
		r := newTestClusterExpressScaling(t, e)
		// Runs before the server closes, which waits for the stalled handler.
		t.Cleanup(func() { close(e.release) })

		start := time.Now()
		resp, got := createExpressScaling(t, r, expressScalingPlan(t, true, "200ms"))

		assert.Less(t, time.Since(start), 5*time.Second, "a stalled GET must not outlive the configured timeout")
		require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())
		assert.Equal(t, 1, resp.Diagnostics.WarningsCount())
		assert.Equal(t, types.BoolValue(true), got.ExpressScaling)
	})
}

func TestClusterExpressScalingCreateSendsConcreteBool(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		current := clusterapi.ExpressScalingEnabled
		want := clusterapi.ExpressScalingDisabled
		if enabled {
			current, want = want, current
		}

		e := &expressScalingEndpoint{states: []clusterapi.ExpressScalingState{current, want}, putStatus: http.StatusNoContent}
		r := newTestClusterExpressScaling(t, e)

		resp, _ := createExpressScaling(t, r, expressScalingPlan(t, enabled, "2s"))
		require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())

		require.Len(t, e.puts, 1)
		if enabled {
			assert.JSONEq(t, "true", string(e.puts[0]))
		} else {
			assert.JSONEq(t, "false", string(e.puts[0]))
		}
	}
}

func TestClusterExpressScalingRead(t *testing.T) {
	tests := []struct {
		name        string
		state       clusterapi.ExpressScalingState
		getStatus   int
		wantRemoved bool
		wantEnabled bool
	}{
		{name: "enabled", state: clusterapi.ExpressScalingEnabled, wantEnabled: true},
		{name: "enabling reads as true", state: clusterapi.ExpressScalingEnabling, wantEnabled: true},
		{name: "disabling reads as false", state: clusterapi.ExpressScalingDisabling, wantEnabled: false},
		{name: "no state reads as false", state: "", wantEnabled: false},
		{name: "cluster gone removes the resource", getStatus: http.StatusNotFound, wantRemoved: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			e := &expressScalingEndpoint{states: []clusterapi.ExpressScalingState{tt.state}, getStatus: tt.getStatus}
			r := newTestClusterExpressScaling(t, e)

			s := ClusterExpressScalingSchema()
			prior := tfsdk.State{Schema: s}
			require.False(t, prior.Set(ctx, expressScalingPlan(t, !tt.wantEnabled, "60m")).HasError())

			resp := &resource.ReadResponse{State: prior}
			r.Read(ctx, resource.ReadRequest{State: prior}, resp)
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())

			if tt.wantRemoved {
				assert.True(t, resp.State.Raw.IsNull())
				return
			}

			var got providerschema.ClusterExpressScaling
			require.False(t, resp.State.Get(ctx, &got).HasError())
			assert.Equal(t, types.BoolValue(tt.wantEnabled), got.ExpressScaling)
			assert.Equal(t, types.StringValue(string(tt.state)), got.ExpressScalingState)
			assert.Equal(t, types.StringValue(testExpressScalingClusterId), got.ClusterId)
		})
	}
}

func TestClusterExpressScalingReadAfterImport(t *testing.T) {
	ctx := t.Context()
	e := &expressScalingEndpoint{states: []clusterapi.ExpressScalingState{clusterapi.ExpressScalingEnabled}}
	r := newTestClusterExpressScaling(t, e)

	s := ClusterExpressScalingSchema()
	imported := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
	importResp := &resource.ImportStateResponse{State: imported}
	r.ImportState(ctx, resource.ImportStateRequest{
		ID: "id=" + testExpressScalingClusterId + ",project_id=" + testExpressScalingProjectId + ",organization_id=" + testExpressScalingOrgId,
	}, importResp)
	require.False(t, importResp.Diagnostics.HasError(), "%v", importResp.Diagnostics.Errors())

	resp := &resource.ReadResponse{State: importResp.State}
	r.Read(ctx, resource.ReadRequest{State: importResp.State}, resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())

	var got providerschema.ClusterExpressScaling
	require.False(t, resp.State.Get(ctx, &got).HasError())
	assert.Equal(t, types.StringValue(testExpressScalingOrgId), got.OrganizationId)
	assert.Equal(t, types.StringValue(testExpressScalingProjectId), got.ProjectId)
	assert.Equal(t, types.StringValue(testExpressScalingClusterId), got.ClusterId)
	assert.Equal(t, types.BoolValue(true), got.ExpressScaling)
}

func diagDetails(diags diag.Diagnostics) string {
	var b strings.Builder
	for _, d := range diags {
		b.WriteString(d.Detail())
		b.WriteString("\n")
	}
	return b.String()
}
