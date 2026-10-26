package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testReplicationID        = "cmVwbC0x"
	testReverseReplicationID = "cmV2LTE="
	testReplicationsPath     = "/v4/organizations/org-1/projects/proj-1/clusters/cluster-1/replications"
	testReplicationPath      = testReplicationsPath + "/" + testReplicationID

	// testClusterNotHealthyMessage is the API message for a 422 while a cluster is not healthy.
	testClusterNotHealthyMessage = "Cluster cluster-1 must be in a healthy state to create or manage replications. " +
		"Currently it is in peering state. Please wait for the cluster to become healthy before attempting replication operations."
)

// testReplicationGetBody is a GET response as the API sends it. The mappings come back in a
// different order from the plan, the all-false documentExcludeOptions are dropped and
// skipRestream is missing, as the real API does.
const testReplicationGetBody = `{
	"id": "cmVwbC0x",
	"status": "running",
	"changesLeft": 0,
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
		{"sourceScope": "inventory", "targetScope": "stock"},
		{"sourceScope": "sales", "targetScope": "sales", "collections": [{"sourceCollection": "orders", "targetCollection": "orders"}]}
	],
	"direction": "twoWay",
	"priority": "medium",
	"networkUsageLimit": 0,
	"filter": {"expressions": {"regEx": "^a"}},
	"mobileSyncMode": "off",
	"audit": {"createdBy": "user-1", "createdAt": "2026-10-01T00:00:00Z"}
}`

// replicationBackend is an in-memory stand-in for the v4 replication endpoints. Create and GET
// replies come from scripted lists of status codes; the last one repeats.
type replicationBackend struct {
	t                  *testing.T
	mu                 sync.Mutex
	createStatuses     []int
	createErrorMessage string
	jobState           string
	jobLastError       string
	getStatuses        []int
	updateStatus       int
	deleteStatus       int
	createBody         string
	updateBody         string
	calls              map[string]int
}

func newReplicationBackend(t *testing.T) *replicationBackend {
	return &replicationBackend{
		t:              t,
		createStatuses: []int{http.StatusAccepted},
		jobState:       "complete",
		getStatuses:    []int{http.StatusOK},
		updateStatus:   http.StatusNoContent,
		deleteStatus:   http.StatusNoContent,
		calls:          map[string]int{},
	}
}

func (b *replicationBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	key := r.Method + " " + r.URL.Path
	b.calls[key]++
	body, err := io.ReadAll(r.Body)
	require.NoError(b.t, err)

	switch key {
	case http.MethodPost + " " + testReplicationsPath:
		b.createBody = string(body)
		status := b.createStatuses[min(b.calls[key]-1, len(b.createStatuses)-1)]
		if status != http.StatusAccepted {
			writeReplicationError(w, status, b.createErrorMessage)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"jobId": "job-1"}`))

	case http.MethodGet + " " + testReplicationsPath + "/jobs/job-1":
		_, _ = fmt.Fprintf(w, `{"jobId": "job-1", "state": %q, "lastError": %q, "replicationId": %q, "reverseReplicationId": %q}`,
			b.jobState, b.jobLastError, testReplicationID, testReverseReplicationID)

	case http.MethodGet + " " + testReplicationPath:
		status := b.getStatuses[min(b.calls[key]-1, len(b.getStatuses)-1)]
		if status != http.StatusOK {
			writeReplicationError(w, status, "")
			return
		}
		_, _ = w.Write([]byte(testReplicationGetBody))

	case http.MethodPut + " " + testReplicationPath:
		b.updateBody = string(body)
		if b.updateStatus != http.StatusNoContent {
			writeReplicationError(w, b.updateStatus, "")
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodDelete + " " + testReplicationPath:
		if b.deleteStatus != http.StatusNoContent {
			writeReplicationError(w, b.deleteStatus, "")
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		b.t.Errorf("unexpected request: %s", key)
		w.WriteHeader(http.StatusTeapot)
	}
}

func (b *replicationBackend) callCount(method, urlPath string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.calls[method+" "+urlPath]
}

// writeReplicationError writes an API error body. An empty message uses the status text.
func writeReplicationError(w http.ResponseWriter, status int, message string) {
	if message == "" {
		message = http.StatusText(status)
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"code": 2, "httpStatusCode": status, "message": message})
}

// fastReplicationCreateRetry shrinks the create retry timing for the duration of a test.
func fastReplicationCreateRetry(t *testing.T, timeout time.Duration) {
	t.Helper()
	origInterval, origTimeout := replicationCreateRetryInterval, replicationCreateRetryTimeout
	replicationCreateRetryInterval = time.Millisecond
	replicationCreateRetryTimeout = timeout
	t.Cleanup(func() {
		replicationCreateRetryInterval = origInterval
		replicationCreateRetryTimeout = origTimeout
	})
}

// fastReplicationDeletePolling shrinks the package-level poll timing for the duration of a test
// so waitForReplicationDeletion's poll loop resolves in milliseconds.
func fastReplicationDeletePolling(t *testing.T, timeout time.Duration) {
	t.Helper()
	origInterval, origTimeout := replicationDeletePollInterval, replicationDeleteTimeout
	replicationDeletePollInterval = time.Millisecond
	replicationDeleteTimeout = timeout
	t.Cleanup(func() {
		replicationDeletePollInterval = origInterval
		replicationDeleteTimeout = origTimeout
	})
}

func testMappingsSet(t *testing.T) types.Set {
	t.Helper()
	collectionType := types.ObjectType{AttrTypes: providerschema.ReplicationCollectionMapping{}.AttributeTypes()}
	scopeType := types.ObjectType{AttrTypes: providerschema.ReplicationScopeMapping{}.AttributeTypes()}

	orders := types.SetValueMust(collectionType, []attr.Value{types.ObjectValueMust(collectionType.AttrTypes, map[string]attr.Value{
		"source_collection": types.StringValue("orders"),
		"target_collection": types.StringValue("orders"),
	})})

	return types.SetValueMust(scopeType, []attr.Value{
		types.ObjectValueMust(scopeType.AttrTypes, map[string]attr.Value{
			"source_scope": types.StringValue("sales"),
			"target_scope": types.StringValue("sales"),
			"collections":  orders,
		}),
		types.ObjectValueMust(scopeType.AttrTypes, map[string]attr.Value{
			"source_scope": types.StringValue("inventory"),
			"target_scope": types.StringValue("stock"),
			"collections":  types.SetNull(collectionType),
		}),
	})
}

// testReplicationPlanModel is the plan Terraform sends to Create: config values with defaults
// applied, and unknown computed values.
func testReplicationPlanModel(t *testing.T) providerschema.Replication {
	t.Helper()
	return providerschema.Replication{
		Id:             types.StringUnknown(),
		OrganizationId: types.StringValue(testOrgID),
		ProjectId:      types.StringValue(testProjectID),
		ClusterId:      types.StringValue(testClusterID),
		SourceBucket:   types.StringValue("c291cmNl"),
		Target: &providerschema.ReplicationTarget{
			Cluster: types.StringValue("cluster-2"),
			Bucket:  types.StringValue("dGFyZ2V0"),
			Type:    types.StringValue(providerschema.ReplicationTargetTypeCapella),
		},
		Direction:         types.StringValue("twoWay"),
		Priority:          types.StringValue("medium"),
		NetworkUsageLimit: types.Int64Value(0),
		MobileSyncMode:    types.StringValue(providerschema.ReplicationMobileSyncModeOff),
		Mappings:          testMappingsSet(t),
		Filter: &providerschema.ReplicationFilter{
			DocumentExcludeOptions: &providerschema.ReplicationDocumentExcludeOptions{
				Deletion: types.BoolValue(false), Expiration: types.BoolValue(false),
				TTL: types.BoolValue(false), Binary: types.BoolValue(false),
			},
			Expressions: &providerschema.ReplicationFilterExpressions{
				RegEx:        types.StringValue("^a"),
				SkipRestream: types.BoolValue(true),
			},
		},
		ReverseReplicationId: types.StringUnknown(),
		Status:               types.StringUnknown(),
		ChangesLeft:          types.Int64Unknown(),
		Error:                types.StringUnknown(),
		SourceDetails:        types.ObjectUnknown(providerschema.ReplicationSourceDetailsAttributeTypes()),
		TargetDetails:        types.ObjectUnknown(providerschema.ReplicationTargetDetailsAttributeTypes()),
		Audit:                types.ObjectUnknown(providerschema.ReplicationAudit{}.AttributeTypes()),
	}
}

// testReplicationStateModel is the state that Create writes for testReplicationPlanModel.
func testReplicationStateModel(t *testing.T) providerschema.Replication {
	t.Helper()
	state := initialReplicationState(testReplicationPlanModel(t))
	state.Id = types.StringValue(testReplicationID)
	state.ReverseReplicationId = types.StringValue(testReverseReplicationID)
	return state
}

func emptyReplicationState(ctx context.Context) tfsdk.State {
	s := ReplicationSchema()
	return tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
}

func replicationStateFrom(t *testing.T, ctx context.Context, model providerschema.Replication) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: ReplicationSchema()}
	assertNoDiags(t, state.Set(ctx, model))
	return state
}

func TestReplicationCreate(t *testing.T) {
	ctx := context.Background()
	fastReplicationJobPolling(t, 5*time.Second)

	t.Run("creates in async mode and keeps the config shape in state", func(t *testing.T) {
		backend := newReplicationBackend(t)
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		planModel := testReplicationPlanModel(t)
		plan := tfsdk.Plan{Schema: ReplicationSchema()}
		assertNoDiags(t, plan.Set(ctx, planModel))

		resp := resource.CreateResponse{State: emptyReplicationState(ctx)}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
		assertNoDiags(t, resp.Diagnostics)

		var got providerschema.Replication
		assertNoDiags(t, resp.State.Get(ctx, &got))

		assert.Equal(t, testReplicationID, got.Id.ValueString())
		assert.Equal(t, testReverseReplicationID, got.ReverseReplicationId.ValueString(),
			"the reverse replication ID comes only from the job")
		assert.Equal(t, "running", got.Status.ValueString())
		assert.True(t, got.Mappings.Equal(planModel.Mappings), "mappings in a different order must give the same set")
		assert.Equal(t, planModel.Filter, got.Filter,
			"skip_restream and the all-false exclude options must be kept from the plan")
		assert.Equal(t, planModel.Target, got.Target)

		var createBody map[string]any
		require.NoError(t, json.Unmarshal([]byte(backend.createBody), &createBody))
		assert.Equal(t, "async", createBody["mode"])
	})

	t.Run("a failed refresh after create is a warning and keeps both IDs", func(t *testing.T) {
		backend := newReplicationBackend(t)
		backend.getStatuses = []int{http.StatusUnprocessableEntity}
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		plan := tfsdk.Plan{Schema: ReplicationSchema()}
		assertNoDiags(t, plan.Set(ctx, testReplicationPlanModel(t)))

		resp := resource.CreateResponse{State: emptyReplicationState(ctx)}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

		require.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics.Errors())
		require.Equal(t, 1, resp.Diagnostics.WarningsCount())

		var got providerschema.Replication
		assertNoDiags(t, resp.State.Get(ctx, &got))
		assert.Equal(t, testReplicationID, got.Id.ValueString())
		assert.Equal(t, testReverseReplicationID, got.ReverseReplicationId.ValueString())
		assert.True(t, got.Status.IsNull())
	})

	t.Run("a failed job is an error and saves no state", func(t *testing.T) {
		backend := newReplicationBackend(t)
		backend.jobState = "failed"
		backend.jobLastError = "job-fatal:source bucket not found"
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		plan := tfsdk.Plan{Schema: ReplicationSchema()}
		assertNoDiags(t, plan.Set(ctx, testReplicationPlanModel(t)))

		resp := resource.CreateResponse{State: emptyReplicationState(ctx)}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "terraform import")
		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("a rejected create call is an error and does not poll the job", func(t *testing.T) {
		backend := newReplicationBackend(t)
		backend.createStatuses = []int{http.StatusConflict}
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		plan := tfsdk.Plan{Schema: ReplicationSchema()}
		assertNoDiags(t, plan.Set(ctx, testReplicationPlanModel(t)))

		resp := resource.CreateResponse{State: emptyReplicationState(ctx)}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "terraform import",
			"a 409 must tell the user to import the existing replication")
		assert.Equal(t, 1, backend.callCount(http.MethodPost, testReplicationsPath), "a 409 must not be retried")
		assert.Equal(t, 0, backend.callCount(http.MethodGet, testReplicationsPath+"/jobs/job-1"))
	})

	t.Run("a create rejected because a cluster is not healthy is retried until it is accepted", func(t *testing.T) {
		fastReplicationCreateRetry(t, 5*time.Second)
		backend := newReplicationBackend(t)
		backend.createStatuses = []int{http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, http.StatusAccepted}
		backend.createErrorMessage = testClusterNotHealthyMessage
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		plan := tfsdk.Plan{Schema: ReplicationSchema()}
		assertNoDiags(t, plan.Set(ctx, testReplicationPlanModel(t)))

		resp := resource.CreateResponse{State: emptyReplicationState(ctx)}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
		assertNoDiags(t, resp.Diagnostics)

		var got providerschema.Replication
		assertNoDiags(t, resp.State.Get(ctx, &got))
		assert.Equal(t, testReplicationID, got.Id.ValueString())
		assert.Equal(t, 3, backend.callCount(http.MethodPost, testReplicationsPath))
	})

	t.Run("a create rejected because a cluster is not healthy fails after the retry timeout", func(t *testing.T) {
		fastReplicationCreateRetry(t, 20*time.Millisecond)
		backend := newReplicationBackend(t)
		backend.createStatuses = []int{http.StatusUnprocessableEntity}
		backend.createErrorMessage = testClusterNotHealthyMessage
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		plan := tfsdk.Plan{Schema: ReplicationSchema()}
		assertNoDiags(t, plan.Set(ctx, testReplicationPlanModel(t)))

		resp := resource.CreateResponse{State: emptyReplicationState(ctx)}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "must be in a healthy state")
		assert.Greater(t, backend.callCount(http.MethodPost, testReplicationsPath), 1)
		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("any other 422 on create is an error and is not retried", func(t *testing.T) {
		fastReplicationCreateRetry(t, 5*time.Second)
		backend := newReplicationBackend(t)
		backend.createStatuses = []int{http.StatusUnprocessableEntity}
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		plan := tfsdk.Plan{Schema: ReplicationSchema()}
		assertNoDiags(t, plan.Set(ctx, testReplicationPlanModel(t)))

		resp := resource.CreateResponse{State: emptyReplicationState(ctx)}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

		require.True(t, resp.Diagnostics.HasError())
		assert.NotContains(t, resp.Diagnostics.Errors()[0].Detail(), "terraform import")
		assert.Equal(t, 1, backend.callCount(http.MethodPost, testReplicationsPath))
	})
}

// TestReplicationSchemaKeepsUnchangedComputedObjects checks that an update plan keeps the values
// that an update cannot change, and leaves the others unknown.
func TestReplicationSchemaKeepsUnchangedComputedObjects(t *testing.T) {
	replicationSchema := ReplicationSchema()

	for _, name := range []string{"source_details", "target_details", "audit"} {
		attribute, ok := replicationSchema.Attributes[name].(*schema.SingleNestedAttribute)
		require.True(t, ok, "%s is %T, want *schema.SingleNestedAttribute", name, replicationSchema.Attributes[name])
		require.Len(t, attribute.PlanModifiers, 1, "%s must keep its state value in an update plan", name)
		assert.Equal(t, objectplanmodifier.UseStateForUnknown().Description(t.Context()),
			attribute.PlanModifiers[0].Description(t.Context()))
	}

	for _, name := range []string{"status", "changes_left", "error"} {
		attribute := replicationSchema.Attributes[name]
		switch a := attribute.(type) {
		case *schema.StringAttribute:
			assert.Empty(t, a.PlanModifiers, "%s can change in an update, so it must stay unknown", name)
		case *schema.Int64Attribute:
			assert.Empty(t, a.PlanModifiers, "%s can change in an update, so it must stay unknown", name)
		default:
			t.Fatalf("%s is %T", name, attribute)
		}
	}
}

func TestReplicationRead(t *testing.T) {
	ctx := context.Background()

	t.Run("refresh keeps the reverse replication ID and the config shape", func(t *testing.T) {
		backend := newReplicationBackend(t)
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		priorModel := testReplicationStateModel(t)
		prior := replicationStateFrom(t, ctx, priorModel)

		resp := resource.ReadResponse{State: prior}
		r.Read(ctx, resource.ReadRequest{State: prior}, &resp)
		assertNoDiags(t, resp.Diagnostics)

		var got providerschema.Replication
		assertNoDiags(t, resp.State.Get(ctx, &got))
		assert.Equal(t, testReverseReplicationID, got.ReverseReplicationId.ValueString())
		assert.Equal(t, priorModel.Filter, got.Filter)
		assert.True(t, got.Mappings.Equal(priorModel.Mappings))
	})

	t.Run("a deleted replication is removed from state", func(t *testing.T) {
		backend := newReplicationBackend(t)
		backend.getStatuses = []int{http.StatusNotFound}
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		prior := replicationStateFrom(t, ctx, testReplicationStateModel(t))
		resp := resource.ReadResponse{State: prior}
		r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

		assertNoDiags(t, resp.Diagnostics)
		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("an unhealthy source cluster (422) is an error and keeps the state", func(t *testing.T) {
		backend := newReplicationBackend(t)
		backend.getStatuses = []int{http.StatusUnprocessableEntity}
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		prior := replicationStateFrom(t, ctx, testReplicationStateModel(t))
		resp := resource.ReadResponse{State: prior}
		r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.False(t, resp.State.Raw.IsNull())
	})

	t.Run("import resolves the IDs and has no reverse replication ID", func(t *testing.T) {
		backend := newReplicationBackend(t)
		r := &Replication{Data: newTestProviderData(t, backend.handler)}

		importModel := emptyReplicationModel()
		importModel.Id = types.StringValue(fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			testReplicationID, testClusterID, testProjectID, testOrgID))
		prior := replicationStateFrom(t, ctx, importModel)

		resp := resource.ReadResponse{State: prior}
		r.Read(ctx, resource.ReadRequest{State: prior}, &resp)
		assertNoDiags(t, resp.Diagnostics)

		var got providerschema.Replication
		assertNoDiags(t, resp.State.Get(ctx, &got))
		assert.Equal(t, testReplicationID, got.Id.ValueString())
		assert.Equal(t, testOrgID, got.OrganizationId.ValueString())
		assert.True(t, got.ReverseReplicationId.IsNull())
		assert.Equal(t, "c291cmNl", got.SourceBucket.ValueString())
		assert.Equal(t, "dGFyZ2V0", got.Target.Bucket.ValueString())
		require.NotNil(t, got.Filter)
		assert.False(t, got.Filter.Expressions.SkipRestream.ValueBool())
	})
}

func TestReplicationUpdate(t *testing.T) {
	ctx := context.Background()

	changePriority := func(plan *providerschema.Replication) { plan.Priority = types.StringValue("low") }

	tests := []struct {
		name         string
		change       func(plan *providerschema.Replication)
		updateStatus int
		getStatus    int
		wantBody     string
		wantError    string
		wantRemoved  bool
	}{
		{
			name:     "only a changed priority is sent with the network usage limit",
			change:   changePriority,
			wantBody: `{"priority": "low", "networkUsageLimit": 0}`,
		},
		{
			name: "removed mappings is sent as allScopes",
			change: func(plan *providerschema.Replication) {
				plan.Mappings = types.SetNull(plan.Mappings.ElementType(ctx))
			},
			wantBody: `{"priority": "medium", "networkUsageLimit": 0, "allScopes": true}`,
		},
		{
			name:         "a replication deleted before the update is removed from state",
			change:       changePriority,
			updateStatus: http.StatusNotFound,
			wantRemoved:  true,
		},
		{
			name:         "a rejected update is an error",
			change:       changePriority,
			updateStatus: http.StatusBadRequest,
			wantError:    "Error updating replication",
		},
		{
			name:      "a failed read after the update is an error",
			change:    changePriority,
			getStatus: http.StatusInternalServerError,
			wantError: "Error reading replication after update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := newReplicationBackend(t)
			if tt.updateStatus != 0 {
				backend.updateStatus = tt.updateStatus
			}
			if tt.getStatus != 0 {
				backend.getStatuses = []int{tt.getStatus}
			}
			r := &Replication{Data: newTestProviderData(t, backend.handler)}

			stateModel := testReplicationStateModel(t)
			state := replicationStateFrom(t, ctx, stateModel)

			planModel := stateModel
			tt.change(&planModel)
			plan := tfsdk.Plan{Schema: ReplicationSchema()}
			assertNoDiags(t, plan.Set(ctx, planModel))

			resp := resource.UpdateResponse{State: state}
			r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)

			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, backend.updateBody)
			}

			if tt.wantError != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Equal(t, tt.wantError, resp.Diagnostics.Errors()[0].Summary())
				assert.False(t, resp.State.Raw.IsNull())
				return
			}
			assertNoDiags(t, resp.Diagnostics)

			if tt.wantRemoved {
				assert.True(t, resp.State.Raw.IsNull())
				return
			}

			var got providerschema.Replication
			assertNoDiags(t, resp.State.Get(ctx, &got))
			assert.Equal(t, testReverseReplicationID, got.ReverseReplicationId.ValueString())
		})
	}
}

func TestReplicationDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		deleteStatus  int
		getStatuses   []int
		timeout       time.Duration
		wantErr       string
		wantGetCalls  int
		checkGetCalls bool
	}{
		{
			name:          "204 deletes with one call and no polling",
			deleteStatus:  http.StatusNoContent,
			checkGetCalls: true,
		},
		{
			name:          "404 means already deleted",
			deleteStatus:  http.StatusNotFound,
			checkGetCalls: true,
		},
		{
			name:          "500 waits until the replication is gone",
			deleteStatus:  http.StatusInternalServerError,
			getStatuses:   []int{http.StatusOK, http.StatusOK, http.StatusNotFound},
			wantGetCalls:  3,
			checkGetCalls: true,
		},
		{
			name:         "500 with the replication still present after the timeout is an error",
			deleteStatus: http.StatusInternalServerError,
			getStatuses:  []int{http.StatusOK},
			timeout:      50 * time.Millisecond,
			wantErr:      "not deleted within the polling timeout",
		},
		{
			name:          "403 is an error and does not poll",
			deleteStatus:  http.StatusForbidden,
			wantErr:       "Forbidden",
			checkGetCalls: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := tt.timeout
			if timeout == 0 {
				timeout = 5 * time.Second
			}
			fastReplicationDeletePolling(t, timeout)

			backend := newReplicationBackend(t)
			backend.deleteStatus = tt.deleteStatus
			if tt.getStatuses != nil {
				backend.getStatuses = tt.getStatuses
			}
			r := &Replication{Data: newTestProviderData(t, backend.handler)}

			state := replicationStateFrom(t, ctx, testReplicationStateModel(t))
			resp := resource.DeleteResponse{State: state}
			r.Delete(ctx, resource.DeleteRequest{State: state}, &resp)

			if tt.wantErr != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.wantErr)
			} else {
				assertNoDiags(t, resp.Diagnostics)
			}

			assert.Equal(t, 1, backend.callCount(http.MethodDelete, testReplicationPath), "only one DELETE is sent")
			if tt.checkGetCalls {
				assert.Equal(t, tt.wantGetCalls, backend.callCount(http.MethodGet, testReplicationPath))
			}
		})
	}
}

// emptyReplicationModel is a replication with every value null.
func emptyReplicationModel() providerschema.Replication {
	return providerschema.Replication{
		Mappings:      types.SetNull(types.ObjectType{AttrTypes: providerschema.ReplicationScopeMapping{}.AttributeTypes()}),
		SourceDetails: types.ObjectNull(providerschema.ReplicationSourceDetailsAttributeTypes()),
		TargetDetails: types.ObjectNull(providerschema.ReplicationTargetDetailsAttributeTypes()),
		Audit:         types.ObjectNull(providerschema.ReplicationAudit{}.AttributeTypes()),
	}
}

func TestReplicationValidateConfig(t *testing.T) {
	ctx := context.Background()

	config := func(sourceBucket, targetBucket, targetType, direction types.String) providerschema.Replication {
		model := emptyReplicationModel()
		model.OrganizationId = types.StringValue(testOrgID)
		model.ProjectId = types.StringValue(testProjectID)
		model.ClusterId = types.StringValue(testClusterID)
		model.SourceBucket = sourceBucket
		model.Target = &providerschema.ReplicationTarget{Cluster: types.StringValue("cluster-2"), Bucket: targetBucket, Type: targetType}
		model.Direction = direction
		return model
	}

	tests := []struct {
		name    string
		model   providerschema.Replication
		wantErr string
	}{
		{
			name:  "capella target with bucket IDs and defaults is valid",
			model: config(types.StringValue("c291cmNl"), types.StringValue("dGVzdA=="), types.StringNull(), types.StringNull()),
		},
		{
			name:  "external target takes a bucket name",
			model: config(types.StringValue("c291cmNl"), types.StringValue("travel-sample"), types.StringValue("external"), types.StringNull()),
		},
		{
			name:    "twoWay with an external target is rejected",
			model:   config(types.StringValue("c291cmNl"), types.StringValue("travel-sample"), types.StringValue("external"), types.StringValue("twoWay")),
			wantErr: "must be oneWay",
		},
		{
			name:    "a bucket name as source_bucket is rejected",
			model:   config(types.StringValue("orders"), types.StringValue("dGVzdA=="), types.StringNull(), types.StringNull()),
			wantErr: "source_bucket must be a bucket ID",
		},
		{
			name:    "a bucket name as a capella target bucket is rejected",
			model:   config(types.StringValue("c291cmNl"), types.StringValue("orders"), types.StringNull(), types.StringNull()),
			wantErr: "target.bucket must be a bucket ID",
		},
		{
			name:  "unknown values are checked later",
			model: config(types.StringUnknown(), types.StringUnknown(), types.StringUnknown(), types.StringUnknown()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := replicationStateFrom(t, ctx, tt.model).Raw
			req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: ReplicationSchema(), Raw: raw}}
			resp := resource.ValidateConfigResponse{}

			(&Replication{}).ValidateConfig(ctx, req, &resp)

			if tt.wantErr == "" {
				assertNoDiags(t, resp.Diagnostics)
				return
			}
			require.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.wantErr)
		})
	}
}

func TestReplicationValidateConfigConflictLogging(t *testing.T) {
	ctx := context.Background()
	ruleType := types.ObjectType{AttrTypes: providerschema.ReplicationLoggingRule{}.AttributeTypes()}

	conflictLogging := func(disabled bool, bucket, collection string, rules map[string]attr.Value) *providerschema.ReplicationConflictLogging {
		loggingRules := types.MapNull(ruleType)
		if rules != nil {
			loggingRules = types.MapValueMust(ruleType, rules)
		}
		return &providerschema.ReplicationConflictLogging{
			Disabled:     types.BoolValue(disabled),
			Bucket:       types.StringValue(bucket),
			Collection:   types.StringValue(collection),
			LoggingRules: loggingRules,
		}
	}
	rule := func(bucket, collection types.String) attr.Value {
		return types.ObjectValueMust(ruleType.AttrTypes, map[string]attr.Value{"bucket": bucket, "collection": collection})
	}

	tests := []struct {
		name            string
		conflictLogging *providerschema.ReplicationConflictLogging
		wantErr         string
	}{
		{
			name:            "all four rule states are valid",
			conflictLogging: testConflictLogging(),
		},
		{
			name:            "disabled with an empty bucket and collection is valid",
			conflictLogging: conflictLogging(true, "", "", nil),
		},
		{
			name:            "enabled with an empty bucket is rejected",
			conflictLogging: conflictLogging(false, "", "logs.default", nil),
			wantErr:         "conflict_logging.bucket must not be empty",
		},
		{
			name:            "enabled with an empty collection is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "", nil),
			wantErr:         "conflict_logging.collection must not be empty",
		},
		{
			name:            "a bucket name as the log bucket is rejected",
			conflictLogging: conflictLogging(false, "conflict_logs", "logs.default", nil),
			wantErr:         "conflict_logging.bucket must be a bucket ID",
		},
		{
			name: "a rule with only a bucket is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "logs.default", map[string]attr.Value{
				"inventory": rule(types.StringValue("bG9ncw=="), types.StringNull()),
			}),
			wantErr: `The logging rule for "inventory" must set both bucket and collection, or neither`,
		},
		{
			name: "a rule with a bucket name is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "logs.default", map[string]attr.Value{
				"inventory": rule(types.StringValue("conflict_logs"), types.StringValue("logs.inv")),
			}),
			wantErr: `conflict_logging.logging_rules["inventory"].bucket must be a bucket ID`,
		},
		{
			name:            "a log collection without a collection name is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "logs.", nil),
			wantErr:         `conflict_logging.collection "logs." must be in the format "scope.collection"`,
		},
		{
			name:            "a log collection with three names is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "a.b.c", nil),
			wantErr:         `conflict_logging.collection "a.b.c" must be in the format "scope.collection"`,
		},
		{
			name: "a rule collection without a scope is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "logs.default", map[string]attr.Value{
				"inventory": rule(types.StringValue("bG9ncw=="), types.StringValue("nodot")),
			}),
			wantErr: `conflict_logging.logging_rules["inventory"].collection "nodot" must be in the format "scope.collection"`,
		},
		{
			name: "an empty rule key is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "logs.default", map[string]attr.Value{
				"": rule(types.StringNull(), types.StringNull()),
			}),
			wantErr: `The logging rule key "" must be in the format "scope" or "scope.collection"`,
		},
		{
			name: "a rule key with three names is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "logs.default", map[string]attr.Value{
				"a.b.c": rule(types.StringNull(), types.StringNull()),
			}),
			wantErr: `The logging rule key "a.b.c" must be in the format "scope" or "scope.collection"`,
		},
		{
			name: "a rule key without a collection name is rejected",
			conflictLogging: conflictLogging(false, "bG9ncw==", "logs.default", map[string]attr.Value{
				"inventory.": types.ObjectNull(ruleType.AttrTypes),
			}),
			wantErr: `The logging rule key "inventory." must be in the format "scope" or "scope.collection"`,
		},
		{
			name: "unknown values are checked later",
			conflictLogging: &providerschema.ReplicationConflictLogging{
				Disabled:   types.BoolValue(false),
				Bucket:     types.StringUnknown(),
				Collection: types.StringUnknown(),
				LoggingRules: types.MapValueMust(ruleType, map[string]attr.Value{
					"inventory": rule(types.StringUnknown(), types.StringValue("logs.inv")),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := emptyReplicationModel()
			model.OrganizationId = types.StringValue(testOrgID)
			model.ProjectId = types.StringValue(testProjectID)
			model.ClusterId = types.StringValue(testClusterID)
			model.SourceBucket = types.StringValue("c291cmNl")
			model.Target = &providerschema.ReplicationTarget{
				Cluster: types.StringValue("cluster-2"), Bucket: types.StringValue("dGVzdA=="), Type: types.StringNull(),
			}
			model.ConflictLogging = tt.conflictLogging

			raw := replicationStateFrom(t, ctx, model).Raw
			req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: ReplicationSchema(), Raw: raw}}
			resp := resource.ValidateConfigResponse{}

			(&Replication{}).ValidateConfig(ctx, req, &resp)

			if tt.wantErr == "" {
				assertNoDiags(t, resp.Diagnostics)
				return
			}
			require.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.wantErr)
		})
	}
}

func TestReplicationValidateConfigMappings(t *testing.T) {
	ctx := context.Background()
	collectionType := types.ObjectType{AttrTypes: providerschema.ReplicationCollectionMapping{}.AttributeTypes()}
	scopeType := types.ObjectType{AttrTypes: providerschema.ReplicationScopeMapping{}.AttributeTypes()}

	collection := func(source, target string) attr.Value {
		return types.ObjectValueMust(collectionType.AttrTypes, map[string]attr.Value{
			"source_collection": types.StringValue(source),
			"target_collection": types.StringValue(target),
		})
	}
	mapping := func(source types.String, target string, collections ...attr.Value) attr.Value {
		collectionSet := types.SetNull(collectionType)
		if len(collections) > 0 {
			collectionSet = types.SetValueMust(collectionType, collections)
		}
		return types.ObjectValueMust(scopeType.AttrTypes, map[string]attr.Value{
			"source_scope": source,
			"target_scope": types.StringValue(target),
			"collections":  collectionSet,
		})
	}

	tests := []struct {
		name     string
		mappings []attr.Value
		wantErr  string
	}{
		{
			name: "each scope mapped once is valid",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "_default"),
				mapping(types.StringValue("tenant_agent_00"), "tenant_agent_00"),
			},
		},
		{
			name: "same scope mapped twice without collections is rejected",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "_default"),
				mapping(types.StringValue("inventory"), "tenant_agent_00"),
			},
			wantErr: `source_scope "inventory" is in more than one mapping without collections`,
		},
		{
			name: "whole scope and one of its collections in another mapping is valid",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "backup"),
				mapping(types.StringValue("inventory"), "archive", collection("airline", "airline")),
			},
		},
		{
			name: "one scope split across two target scopes is valid",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "inventory", collection("airline", "airline")),
				mapping(types.StringValue("inventory"), "archive", collection("route", "route")),
			},
		},
		{
			name: "same source and target scope in two mappings with collections is rejected",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "inventory", collection("airline", "airline")),
				mapping(types.StringValue("inventory"), "inventory", collection("route", "route")),
			},
			wantErr: `source_scope "inventory" with target_scope "inventory" is in more than one mapping with collections`,
		},
		{
			name: "same collection mapped twice in one scope is rejected",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "inventory", collection("airline", "airline"), collection("airline", "route")),
			},
			wantErr: `source_collection "airline" of source_scope "inventory" is mapped more than once`,
		},
		{
			name: "same collection mapped to two target scopes is rejected",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "inventory", collection("airline", "airline")),
				mapping(types.StringValue("inventory"), "archive", collection("airline", "airline")),
			},
			wantErr: `source_collection "airline" of source_scope "inventory" is mapped more than once`,
		},
		{
			name: "two collections mapped to one target collection is rejected",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "inventory", collection("airline", "airline"), collection("route", "airline")),
			},
			wantErr: `target_collection "airline" of target_scope "inventory" is the target of more than one collection mapping`,
		},
		{
			name: "same target collection from two source scopes is rejected",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "backup", collection("airline", "airline")),
				mapping(types.StringValue("tenant"), "backup", collection("airline", "airline")),
			},
			wantErr: `target_collection "airline" of target_scope "backup" is the target of more than one collection mapping`,
		},
		{
			name: "same collection name in two different scopes is valid",
			mappings: []attr.Value{
				mapping(types.StringValue("inventory"), "inventory", collection("airline", "airline")),
				mapping(types.StringValue("archive"), "archive", collection("airline", "airline")),
			},
		},
		{
			name: "scope names that are not known yet are not checked",
			mappings: []attr.Value{
				mapping(types.StringUnknown(), "_default"),
				mapping(types.StringUnknown(), "tenant_agent_00"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := emptyReplicationModel()
			model.OrganizationId = types.StringValue(testOrgID)
			model.ProjectId = types.StringValue(testProjectID)
			model.ClusterId = types.StringValue(testClusterID)
			model.SourceBucket = types.StringValue("c291cmNl")
			model.Target = &providerschema.ReplicationTarget{
				Cluster: types.StringValue("cluster-2"), Bucket: types.StringValue("dGVzdA=="), Type: types.StringNull(),
			}
			model.Mappings = types.SetValueMust(scopeType, tt.mappings)

			raw := replicationStateFrom(t, ctx, model).Raw
			req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: ReplicationSchema(), Raw: raw}}
			resp := resource.ValidateConfigResponse{}

			(&Replication{}).ValidateConfig(ctx, req, &resp)

			if tt.wantErr == "" {
				assertNoDiags(t, resp.Diagnostics)
				return
			}
			require.Len(t, resp.Diagnostics.Errors(), 1)
			assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.wantErr)
		})
	}
}

// TestReplicationValidateConfigReportsAllErrors checks that an error in another attribute does not
// hide the conflict logging errors.
func TestReplicationValidateConfigReportsAllErrors(t *testing.T) {
	ctx := context.Background()

	model := emptyReplicationModel()
	model.OrganizationId = types.StringValue(testOrgID)
	model.ProjectId = types.StringValue(testProjectID)
	model.ClusterId = types.StringValue(testClusterID)
	model.SourceBucket = types.StringValue("orders")
	model.Target = &providerschema.ReplicationTarget{
		Cluster: types.StringValue("cluster-2"), Bucket: types.StringValue("dGVzdA=="), Type: types.StringNull(),
	}
	model.ConflictLogging = testConflictLogging()
	model.ConflictLogging.Bucket = types.StringValue("conflict_logs")

	raw := replicationStateFrom(t, ctx, model).Raw
	req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: ReplicationSchema(), Raw: raw}}
	resp := resource.ValidateConfigResponse{}

	(&Replication{}).ValidateConfig(ctx, req, &resp)

	require.Len(t, resp.Diagnostics.Errors(), 2)
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "source_bucket must be a bucket ID")
	assert.Contains(t, resp.Diagnostics.Errors()[1].Detail(), "conflict_logging.bucket must be a bucket ID")
}

// TestReplicationValidateConfigTargetNotKnown checks that a null or unknown target object gives no
// error of its own. Terraform reports a missing target, and checks an unknown one again later.
func TestReplicationValidateConfigTargetNotKnown(t *testing.T) {
	ctx := context.Background()
	replicationSchema := ReplicationSchema()

	model := emptyReplicationModel()
	model.OrganizationId = types.StringValue(testOrgID)
	model.ProjectId = types.StringValue(testProjectID)
	model.ClusterId = types.StringValue(testClusterID)
	model.SourceBucket = types.StringValue("c291cmNl")
	nullTarget := replicationStateFrom(t, ctx, model).Raw

	var values map[string]tftypes.Value
	require.NoError(t, nullTarget.As(&values))
	objectType, ok := replicationSchema.Type().TerraformType(ctx).(tftypes.Object)
	require.True(t, ok)
	values["target"] = tftypes.NewValue(objectType.AttributeTypes["target"], tftypes.UnknownValue)
	unknownTarget := tftypes.NewValue(objectType, values)

	for name, raw := range map[string]tftypes.Value{"null target": nullTarget, "unknown target": unknownTarget} {
		t.Run(name, func(t *testing.T) {
			req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: replicationSchema, Raw: raw}}
			resp := resource.ValidateConfigResponse{}

			(&Replication{}).ValidateConfig(ctx, req, &resp)

			assertNoDiags(t, resp.Diagnostics)
		})
	}
}

// testScopeMappingsSet builds a mappings set with one scope, and optionally one collection.
func testScopeMappingsSet(t *testing.T, scope, collection string) types.Set {
	t.Helper()
	collectionType := types.ObjectType{AttrTypes: providerschema.ReplicationCollectionMapping{}.AttributeTypes()}
	scopeType := types.ObjectType{AttrTypes: providerschema.ReplicationScopeMapping{}.AttributeTypes()}

	collections := types.SetNull(collectionType)
	if collection != "" {
		collections = types.SetValueMust(collectionType, []attr.Value{types.ObjectValueMust(collectionType.AttrTypes, map[string]attr.Value{
			"source_collection": types.StringValue(collection),
			"target_collection": types.StringValue(collection),
		})})
	}

	return types.SetValueMust(scopeType, []attr.Value{types.ObjectValueMust(scopeType.AttrTypes, map[string]attr.Value{
		"source_scope": types.StringValue(scope),
		"target_scope": types.StringValue(scope),
		"collections":  collections,
	})})
}

func allFalseDocumentExcludeOptions() *providerschema.ReplicationDocumentExcludeOptions {
	return &providerschema.ReplicationDocumentExcludeOptions{
		Deletion: types.BoolValue(false), Expiration: types.BoolValue(false),
		TTL: types.BoolValue(false), Binary: types.BoolValue(false),
	}
}

// testConflictLoggingJSON is the API shape of testConflictLogging.
const testConflictLoggingJSON = `{
	"disabled": false,
	"bucket": "bG9ncw==",
	"collection": "logs.default",
	"loggingRules": {
		"inventory": {"bucket": "bG9ncw==", "collection": "logs.inv"},
		"sales": {},
		"agents.partners": null
	}
}`

// testConflictLogging has a rule of each kind: a custom target, {} (log to the default) and null (do
// not log). A namespace with no rule, the fourth state, is not in the map.
func testConflictLogging() *providerschema.ReplicationConflictLogging {
	ruleType := types.ObjectType{AttrTypes: providerschema.ReplicationLoggingRule{}.AttributeTypes()}
	return &providerschema.ReplicationConflictLogging{
		Disabled:   types.BoolValue(false),
		Bucket:     types.StringValue("bG9ncw=="),
		Collection: types.StringValue("logs.default"),
		LoggingRules: types.MapValueMust(ruleType, map[string]attr.Value{
			"inventory": types.ObjectValueMust(ruleType.AttrTypes, map[string]attr.Value{
				"bucket":     types.StringValue("bG9ncw=="),
				"collection": types.StringValue("logs.inv"),
			}),
			"sales": types.ObjectValueMust(ruleType.AttrTypes, map[string]attr.Value{
				"bucket":     types.StringNull(),
				"collection": types.StringNull(),
			}),
			"agents.partners": types.ObjectNull(ruleType.AttrTypes),
		}),
	}
}

// testReplicationRequestModel is a plan with one mapping and a filter, for the request builders.
func testReplicationRequestModel(t *testing.T) providerschema.Replication {
	t.Helper()
	return providerschema.Replication{
		SourceBucket: types.StringValue("c291cmNl"),
		Target: &providerschema.ReplicationTarget{
			Cluster: types.StringValue("cluster-2"),
			Bucket:  types.StringValue("dGFyZ2V0"),
			Type:    types.StringValue(providerschema.ReplicationTargetTypeCapella),
		},
		Direction:         types.StringValue("twoWay"),
		Priority:          types.StringValue("medium"),
		NetworkUsageLimit: types.Int64Value(0),
		MobileSyncMode:    types.StringValue(providerschema.ReplicationMobileSyncModeOff),
		Mappings:          testScopeMappingsSet(t, "sales", "orders"),
		Filter: &providerschema.ReplicationFilter{
			DocumentExcludeOptions: allFalseDocumentExcludeOptions(),
			Expressions: &providerschema.ReplicationFilterExpressions{
				RegEx:        types.StringValue("^a"),
				SkipRestream: types.BoolValue(false),
			},
		},
	}
}

func TestCreateReplicationRequest(t *testing.T) {
	got, err := createReplicationRequest(context.Background(), testReplicationRequestModel(t))
	require.NoError(t, err)

	body, err := json.Marshal(got)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"sourceBucket": "c291cmNl",
		"target": {"cluster": "cluster-2", "bucket": "dGFyZ2V0", "type": "capella"},
		"mode": "async",
		"direction": "twoWay",
		"priority": "medium",
		"networkUsageLimit": 0,
		"mappings": [{"sourceScope": "sales", "targetScope": "sales", "collections": [{"sourceCollection": "orders", "targetCollection": "orders"}]}],
		"filter": {
			"documentExcludeOptions": {"deletion": false, "expiration": false, "ttl": false, "binary": false},
			"expressions": {"regEx": "^a", "skipRestream": false}
		}
	}`, string(body))
}

func TestCreateReplicationRequestMobileSyncModeAndConflictLogging(t *testing.T) {
	tests := []struct {
		name   string
		change func(plan *providerschema.Replication)
		want   map[string]string
	}{
		{
			name:   "the default off is not sent, and no block sends no conflict logging",
			change: func(*providerschema.Replication) {},
			want:   map[string]string{},
		},
		{
			name: "active is sent",
			change: func(plan *providerschema.Replication) {
				plan.MobileSyncMode = types.StringValue(providerschema.ReplicationMobileSyncModeActive)
			},
			want: map[string]string{"mobileSyncMode": `"active"`},
		},
		{
			name:   "conflict logging is sent with every rule state",
			change: func(plan *providerschema.Replication) { plan.ConflictLogging = testConflictLogging() },
			want:   map[string]string{"conflictLogging": testConflictLoggingJSON},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := testReplicationRequestModel(t)
			tt.change(&plan)

			got, err := createReplicationRequest(context.Background(), plan)
			require.NoError(t, err)

			body, err := json.Marshal(got)
			require.NoError(t, err)
			var fields map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(body, &fields))

			for _, key := range []string{"mobileSyncMode", "conflictLogging"} {
				want, ok := tt.want[key]
				if !ok {
					assert.NotContains(t, fields, key)
					continue
				}
				require.Contains(t, fields, key)
				assert.JSONEq(t, want, string(fields[key]))
			}
		})
	}
}

func TestUpdateReplicationRequest(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		changeState func(state *providerschema.Replication)
		change      func(plan *providerschema.Replication)
		want        string
	}{
		{
			name:   "no change sends only priority and network usage limit",
			change: func(*providerschema.Replication) {},
			want:   `{"priority": "medium", "networkUsageLimit": 0}`,
		},
		{
			name:   "changed priority does not send the unchanged mappings or filter",
			change: func(plan *providerschema.Replication) { plan.Priority = types.StringValue("low") },
			want:   `{"priority": "low", "networkUsageLimit": 0}`,
		},
		{
			name: "removed mappings sends allScopes",
			change: func(plan *providerschema.Replication) {
				plan.Mappings = types.SetNull(plan.Mappings.ElementType(ctx))
			},
			want: `{"priority": "medium", "networkUsageLimit": 0, "allScopes": true}`,
		},
		{
			name: "changed mappings sends the new mappings",
			change: func(plan *providerschema.Replication) {
				plan.Mappings = testScopeMappingsSet(t, "inventory", "")
			},
			want: `{"priority": "medium", "networkUsageLimit": 0, "mappings": [{"sourceScope": "inventory", "targetScope": "inventory"}]}`,
		},
		{
			name:   "removed filter sends the full default filter",
			change: func(plan *providerschema.Replication) { plan.Filter = nil },
			want: `{"priority": "medium", "networkUsageLimit": 0, "filter": {
				"documentExcludeOptions": {"deletion": false, "expiration": false, "ttl": false, "binary": false},
				"expressions": {"regEx": "", "skipRestream": false}}}`,
		},
		{
			name: "removed exclude options with the expression kept sends all four flags as false",
			changeState: func(state *providerschema.Replication) {
				deletion := allFalseDocumentExcludeOptions()
				deletion.Deletion = types.BoolValue(true)
				state.Filter = &providerschema.ReplicationFilter{
					DocumentExcludeOptions: deletion,
					Expressions:            state.Filter.Expressions,
				}
			},
			change: func(plan *providerschema.Replication) {
				plan.Filter = &providerschema.ReplicationFilter{Expressions: plan.Filter.Expressions}
			},
			want: `{"priority": "medium", "networkUsageLimit": 0, "filter": {
				"documentExcludeOptions": {"deletion": false, "expiration": false, "ttl": false, "binary": false},
				"expressions": {"regEx": "^a", "skipRestream": false}}}`,
		},
		{
			name: "changed expression sends the full filter",
			change: func(plan *providerschema.Replication) {
				plan.Filter = &providerschema.ReplicationFilter{
					DocumentExcludeOptions: allFalseDocumentExcludeOptions(),
					Expressions: &providerschema.ReplicationFilterExpressions{
						RegEx:        types.StringValue("^b"),
						SkipRestream: types.BoolValue(true),
					},
				}
			},
			want: `{"priority": "medium", "networkUsageLimit": 0, "filter": {
				"documentExcludeOptions": {"deletion": false, "expiration": false, "ttl": false, "binary": false},
				"expressions": {"regEx": "^b", "skipRestream": true}}}`,
		},
		{
			name: "changed mobile sync mode is sent",
			change: func(plan *providerschema.Replication) {
				plan.MobileSyncMode = types.StringValue(providerschema.ReplicationMobileSyncModeActive)
			},
			want: `{"priority": "medium", "networkUsageLimit": 0, "mobileSyncMode": "active"}`,
		},
		{
			name: "mobile sync mode set back to off is sent",
			changeState: func(state *providerschema.Replication) {
				state.MobileSyncMode = types.StringValue(providerschema.ReplicationMobileSyncModeActive)
			},
			change: func(*providerschema.Replication) {},
			want:   `{"priority": "medium", "networkUsageLimit": 0, "mobileSyncMode": "off"}`,
		},
		{
			name:   "added conflict logging sends the full configuration",
			change: func(plan *providerschema.Replication) { plan.ConflictLogging = testConflictLogging() },
			want:   `{"priority": "medium", "networkUsageLimit": 0, "conflictLogging": ` + testConflictLoggingJSON + `}`,
		},
		{
			name:        "removed conflict logging sends an empty object, which clears it",
			changeState: func(state *providerschema.Replication) { state.ConflictLogging = testConflictLogging() },
			change:      func(*providerschema.Replication) {},
			want:        `{"priority": "medium", "networkUsageLimit": 0, "conflictLogging": {}}`,
		},
		{
			name:        "a changed rule sends the full configuration, because it replaces the old one",
			changeState: func(state *providerschema.Replication) { state.ConflictLogging = testConflictLogging() },
			change: func(plan *providerschema.Replication) {
				plan.ConflictLogging = testConflictLogging()
				plan.ConflictLogging.Disabled = types.BoolValue(true)
			},
			want: `{"priority": "medium", "networkUsageLimit": 0, "conflictLogging": {
				"disabled": true,
				"bucket": "bG9ncw==",
				"collection": "logs.default",
				"loggingRules": {
					"inventory": {"bucket": "bG9ncw==", "collection": "logs.inv"},
					"sales": {},
					"agents.partners": null
				}}}`,
		},
		{
			name: "changed priority does not send the unchanged mobile sync mode or conflict logging",
			changeState: func(state *providerschema.Replication) {
				state.MobileSyncMode = types.StringValue(providerschema.ReplicationMobileSyncModeActive)
				state.ConflictLogging = testConflictLogging()
			},
			change: func(plan *providerschema.Replication) {
				plan.Priority = types.StringValue("low")
				plan.MobileSyncMode = types.StringValue(providerschema.ReplicationMobileSyncModeActive)
				plan.ConflictLogging = testConflictLogging()
			},
			want: `{"priority": "low", "networkUsageLimit": 0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := testReplicationRequestModel(t)
			if tt.changeState != nil {
				tt.changeState(&state)
			}
			plan := testReplicationRequestModel(t)
			tt.change(&plan)

			got, err := updateReplicationRequest(ctx, plan, state)
			require.NoError(t, err)

			body, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(body))
		})
	}
}
