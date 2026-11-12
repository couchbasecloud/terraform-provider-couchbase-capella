package resources

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testMaintenanceID = "req1"
	// testMaintenanceBody is a pending, Couchbase-initiated record queued at 2026-09-05T02:30:00Z.
	testMaintenanceBody = `{"id":"req1","type":"upgrade","initiatedBy":"couchbase","currentState":"pending",` +
		`"queueAt":"2026-09-05T02:30:00.000+00:00"}`
)

// fakeMaintenanceServer answers GET with getStatus and getBody and PUT with 204, and records each
// request's method and body.
type fakeMaintenanceServer struct {
	mu        sync.Mutex
	methods   []string
	putBodies []string
}

func (f *fakeMaintenanceServer) calls() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.methods...)
}

func newTestMaintenanceSchedule(t *testing.T, getStatus int, getBody string) (*AppServiceMaintenanceSchedule, *fakeMaintenanceServer) {
	t.Helper()
	fake := &fakeMaintenanceServer{}
	data := newTestProviderData(t, func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/appservices/"+testAppServiceID+"/maintenances/"+testMaintenanceID), r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		fake.mu.Lock()
		fake.methods = append(fake.methods, r.Method)
		if r.Method == http.MethodPut {
			fake.putBodies = append(fake.putBodies, string(body))
		}
		fake.mu.Unlock()

		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(getStatus)
		_, _ = w.Write([]byte(getBody))
	})
	return &AppServiceMaintenanceSchedule{Data: data}, fake
}

func testMaintenanceScheduleModel(queueAt string) providerschema.AppServiceMaintenanceSchedule {
	return providerschema.AppServiceMaintenanceSchedule{
		OrganizationId: types.StringValue(testOrgID),
		ProjectId:      types.StringValue(testProjectID),
		ClusterId:      types.StringValue(testClusterID),
		AppServiceId:   types.StringValue(testAppServiceID),
		MaintenanceId:  types.StringValue(testMaintenanceID),
		QueueAt:        timetypes.NewRFC3339ValueMust(queueAt),
		Kind:           types.StringUnknown(),
		InitiatedBy:    types.StringUnknown(),
		VersionId:      types.StringUnknown(),
		CurrentState:   types.StringUnknown(),
		CompletedAt:    timetypes.NewRFC3339Unknown(),
		Audit:          types.ObjectUnknown(providerschema.CouchbaseAuditData{}.AttributeTypes()),
	}
}

// testMaintenanceScheduleState is a stored state, as Read and Update receive it.
func testMaintenanceScheduleState(initiatedBy string) providerschema.AppServiceMaintenanceSchedule {
	model := testMaintenanceScheduleModel("2026-09-05T02:30:00Z")
	model.Kind, model.VersionId, model.CurrentState = types.StringNull(), types.StringNull(), types.StringNull()
	model.InitiatedBy = types.StringValue(initiatedBy)
	model.CompletedAt = timetypes.NewRFC3339Null()
	model.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	return model
}

func runMaintenanceScheduleCreate(t *testing.T, r *AppServiceMaintenanceSchedule, queueAt string) *resource.CreateResponse {
	t.Helper()
	ctx := context.Background()
	s := AppServiceMaintenanceScheduleSchema()
	plan := tfsdk.Plan{Schema: s}
	require.False(t, plan.Set(ctx, testMaintenanceScheduleModel(queueAt)).HasError())
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: s}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)
	return resp
}

func TestAppServiceMaintenanceScheduleComputedAttributesHaveNoValidators(t *testing.T) {
	attrs := AppServiceMaintenanceScheduleSchema().Attributes
	for _, name := range []string{"kind", "initiated_by", "version_id", "current_state", "completed_at"} {
		attr, ok := attrs[name].(*schema.StringAttribute)
		require.True(t, ok, name)
		assert.True(t, attr.Computed, name)
		assert.Empty(t, attr.Validators, "a validator on %s would break Read on a new server-side value", name)
	}
}

func TestAppServiceMaintenanceScheduleCreate(t *testing.T) {
	tests := []struct {
		name            string
		getStatus       int
		getBody         string
		queueAt         string
		expectedCalls   []string
		expectedError   string
		expectedPutBody string
	}{
		{
			name:          "semantically equal queue_at skips the PUT",
			getStatus:     http.StatusOK,
			getBody:       testMaintenanceBody,
			queueAt:       "2026-09-05T02:30:00Z",
			expectedCalls: []string{http.MethodGet},
		},
		{
			name:            "different queue_at sends the reschedule",
			getStatus:       http.StatusOK,
			getBody:         testMaintenanceBody,
			queueAt:         "2026-09-06T04:00:00Z",
			expectedCalls:   []string{http.MethodGet, http.MethodPut, http.MethodGet},
			expectedPutBody: `{"queueAt":"2026-09-06T04:00:00Z"}`,
		},
		{
			name:      "self-service record is refused without a PUT",
			getStatus: http.StatusOK,
			getBody: `{"id":"req1","type":"upgrade","initiatedBy":"selfService","currentState":"pending",` +
				`"queueAt":"2026-09-05T02:30:00Z"}`,
			queueAt:       "2026-09-06T04:00:00Z",
			expectedCalls: []string{http.MethodGet},
			expectedError: "couchbase-capella_app_service_self_service_upgrade",
		},
		{
			name:          "missing record is a hard error",
			getStatus:     http.StatusNotFound,
			getBody:       `{"code":404,"message":"not found","httpStatusCode":404}`,
			queueAt:       "2026-09-06T04:00:00Z",
			expectedCalls: []string{http.MethodGet},
			expectedError: "does not exist on app service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, fake := newTestMaintenanceSchedule(t, tt.getStatus, tt.getBody)

			resp := runMaintenanceScheduleCreate(t, r, tt.queueAt)

			assert.Equal(t, tt.expectedCalls, fake.calls())
			if tt.expectedError != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.expectedError)
				assert.True(t, resp.State.Raw.IsNull(), "nothing may be written to state")
				return
			}
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			if tt.expectedPutBody != "" {
				require.Len(t, fake.putBodies, 1)
				assert.JSONEq(t, tt.expectedPutBody, fake.putBodies[0])
			}
		})
	}
}

func TestAppServiceMaintenanceScheduleCreateStateIsSemanticallyEqualToPlan(t *testing.T) {
	ctx := context.Background()
	r, _ := newTestMaintenanceSchedule(t, http.StatusOK, testMaintenanceBody)

	resp := runMaintenanceScheduleCreate(t, r, "2026-09-05T02:30:00Z")

	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
	var got providerschema.AppServiceMaintenanceSchedule
	require.False(t, resp.State.Get(ctx, &got).HasError())
	equal, diags := timetypes.NewRFC3339ValueMust("2026-09-05T02:30:00Z").StringSemanticEquals(ctx, got.QueueAt)
	require.False(t, diags.HasError())
	assert.True(t, equal, "queue_at %s must not produce a diff", got.QueueAt)
}

func TestAppServiceMaintenanceScheduleRead(t *testing.T) {
	tests := []struct {
		name          string
		getStatus     int
		getBody       string
		expectRemoved bool
		expectError   string
		expected      func(*testing.T, providerschema.AppServiceMaintenanceSchedule)
	}{
		{
			name:      "maps the response to state",
			getStatus: http.StatusOK,
			getBody:   testMaintenanceBody,
			expected: func(t *testing.T, got providerschema.AppServiceMaintenanceSchedule) {
				assert.Equal(t, "upgrade", got.Kind.ValueString())
				assert.Equal(t, "couchbase", got.InitiatedBy.ValueString())
				assert.True(t, got.VersionId.IsNull())
				assert.True(t, got.CompletedAt.IsNull())
				assert.Equal(t, testMaintenanceID, got.MaintenanceId.ValueString())
			},
		},
		{
			name:      "unknown current_state is kept as-is",
			getStatus: http.StatusOK,
			getBody:   `{"id":"req1","type":"upgrade","initiatedBy":"couchbase","currentState":"paused","queueAt":"2026-09-05T02:30:00Z"}`,
			expected: func(t *testing.T, got providerschema.AppServiceMaintenanceSchedule) {
				assert.Equal(t, "paused", got.CurrentState.ValueString())
			},
		},
		{
			name:        "self-service record is refused, which also blocks its import",
			getStatus:   http.StatusOK,
			getBody:     `{"id":"req1","type":"upgrade","initiatedBy":"selfService","currentState":"pending","queueAt":"2026-09-05T02:30:00Z"}`,
			expectError: "couchbase-capella_app_service_self_service_upgrade",
		},
		{
			name:          "404 removes the resource from state",
			getStatus:     http.StatusNotFound,
			getBody:       `{"code":404,"message":"not found","httpStatusCode":404}`,
			expectRemoved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r, _ := newTestMaintenanceSchedule(t, tt.getStatus, tt.getBody)
			s := AppServiceMaintenanceScheduleSchema()
			state := tfsdk.State{Schema: s}
			require.False(t, state.Set(ctx, testMaintenanceScheduleState("couchbase")).HasError())
			resp := &resource.ReadResponse{State: state}

			r.Read(ctx, resource.ReadRequest{State: state}, resp)

			if tt.expectError != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.expectError)
				return
			}
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			if tt.expectRemoved {
				assert.True(t, resp.State.Raw.IsNull())
				return
			}
			var got providerschema.AppServiceMaintenanceSchedule
			require.False(t, resp.State.Get(ctx, &got).HasError())
			tt.expected(t, got)
		})
	}
}

func TestAppServiceMaintenanceScheduleDeleteMakesNoCall(t *testing.T) {
	r, fake := newTestMaintenanceSchedule(t, http.StatusOK, testMaintenanceBody)
	resp := &resource.DeleteResponse{}

	r.Delete(context.Background(), resource.DeleteRequest{}, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Empty(t, fake.calls())
}

func TestAppServiceMaintenanceScheduleUpdate(t *testing.T) {
	tests := []struct {
		name          string
		initiatedBy   string
		expectedCalls []string
		expectedError string
	}{
		{
			name:          "reschedules then reads back",
			initiatedBy:   "couchbase",
			expectedCalls: []string{http.MethodPut, http.MethodGet},
		},
		{
			name:          "self-service record is refused without a PUT",
			initiatedBy:   "selfService",
			expectedError: "couchbase-capella_app_service_self_service_upgrade",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r, fake := newTestMaintenanceSchedule(t, http.StatusOK, testMaintenanceBody)
			s := AppServiceMaintenanceScheduleSchema()
			plan := tfsdk.Plan{Schema: s}
			require.False(t, plan.Set(ctx, testMaintenanceScheduleModel("2026-09-06T04:00:00Z")).HasError())
			state := tfsdk.State{Schema: s}
			require.False(t, state.Set(ctx, testMaintenanceScheduleState(tt.initiatedBy)).HasError())
			resp := &resource.UpdateResponse{State: state}

			r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, resp)

			assert.Equal(t, tt.expectedCalls, fake.calls())
			if tt.expectedError != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.expectedError)
				return
			}
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			require.Len(t, fake.putBodies, 1)
			assert.JSONEq(t, `{"queueAt":"2026-09-06T04:00:00Z"}`, fake.putBodies[0])
		})
	}
}
