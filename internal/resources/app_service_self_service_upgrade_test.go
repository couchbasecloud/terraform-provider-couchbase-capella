package resources

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_self_service_upgrade"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testUpgradeID  = "upgrade-1"
	testNotFound   = `{"code":404,"message":"not found","httpStatusCode":404}`
	testUpgradeFmt = `{"id":"upgrade-1","type":"upgrade","initiatedBy":"selfService","currentState":"%s","versionId":"3.1.8",` +
		`"queueAt":"2026-08-24T02:30:00.000+00:00"}`
)

func testUpgradeBody(state string) string {
	return strings.Replace(testUpgradeFmt, "%s", state, 1)
}

// scriptedResponse is one canned answer of fakeUpgradeServer.
type scriptedResponse struct {
	status int
	body   string
}

// fakeUpgradeServer answers each method from its own queue, repeating the last entry once the queue
// is drained, and records every request's method and body.
type fakeUpgradeServer struct {
	mu        sync.Mutex
	responses map[string][]scriptedResponse
	methods   []string
	bodies    map[string][]string
}

func (f *fakeUpgradeServer) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		defer f.mu.Unlock()

		if r.Method == http.MethodPost {
			assert.True(t, strings.HasSuffix(r.URL.Path, "/appservices/"+testAppServiceID+"/selfServiceUpgrades"), r.URL.Path)
		} else {
			assert.True(t, strings.HasSuffix(r.URL.Path, "/appservices/"+testAppServiceID+"/selfServiceUpgrades/"+testUpgradeID), r.URL.Path)
		}
		f.methods = append(f.methods, r.Method)
		f.bodies[r.Method] = append(f.bodies[r.Method], string(body))

		queue := f.responses[r.Method]
		require.NotEmpty(t, queue, "unexpected %s", r.Method)
		next := queue[0]
		if len(queue) > 1 {
			f.responses[r.Method] = queue[1:]
		}
		w.WriteHeader(next.status)
		_, _ = w.Write([]byte(next.body))
	}
}

func (f *fakeUpgradeServer) calls() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.methods...)
}

func newTestSelfServiceUpgrade(t *testing.T, responses map[string][]scriptedResponse) (*AppServiceSelfServiceUpgrade, *fakeUpgradeServer) {
	t.Helper()
	fake := &fakeUpgradeServer{responses: responses, bodies: map[string][]string{}}
	return &AppServiceSelfServiceUpgrade{Data: newTestProviderData(t, fake.handler(t))}, fake
}

// fastSelfServiceUpgradePolling shrinks the create poll timing for the duration of a test.
func fastSelfServiceUpgradePolling(t *testing.T) {
	t.Helper()
	origInterval, origTimeout := selfServiceUpgradePollInterval, selfServiceUpgradeCreateTimeout
	selfServiceUpgradePollInterval = time.Millisecond
	selfServiceUpgradeCreateTimeout = 200 * time.Millisecond
	t.Cleanup(func() {
		selfServiceUpgradePollInterval = origInterval
		selfServiceUpgradeCreateTimeout = origTimeout
	})
}

// testSelfServiceUpgradePlan is a create plan; queueAt "" leaves queue_at unknown, as when it is not configured.
func testSelfServiceUpgradePlan(queueAt string) providerschema.AppServiceSelfServiceUpgrade {
	model := providerschema.AppServiceSelfServiceUpgrade{
		OrganizationId: types.StringValue(testOrgID),
		ProjectId:      types.StringValue(testProjectID),
		ClusterId:      types.StringValue(testClusterID),
		AppServiceId:   types.StringValue(testAppServiceID),
		Id:             types.StringUnknown(),
		VersionId:      types.StringValue("3.1.8"),
		QueueAt:        timetypes.NewRFC3339Unknown(),
		CurrentState:   types.StringUnknown(),
		CompletedAt:    timetypes.NewRFC3339Unknown(),
		Audit:          types.ObjectUnknown(providerschema.CouchbaseAuditData{}.AttributeTypes()),
	}
	if queueAt != "" {
		model.QueueAt = timetypes.NewRFC3339ValueMust(queueAt)
	}
	return model
}

// testSelfServiceUpgradeState is a stored state, as Read, Update and Delete receive it.
func testSelfServiceUpgradeState() providerschema.AppServiceSelfServiceUpgrade {
	model := testSelfServiceUpgradePlan("2026-08-24T02:30:00Z")
	model.Id = types.StringValue(testUpgradeID)
	model.CurrentState = types.StringValue("pending")
	model.CompletedAt = timetypes.NewRFC3339Null()
	model.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	return model
}

func runSelfServiceUpgradeCreate(t *testing.T, r *AppServiceSelfServiceUpgrade, plan providerschema.AppServiceSelfServiceUpgrade) *resource.CreateResponse {
	t.Helper()
	s := AppServiceSelfServiceUpgradeSchema()
	p := tfsdk.Plan{Schema: s}
	require.False(t, p.Set(context.Background(), plan).HasError())
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: s}}
	r.Create(context.Background(), resource.CreateRequest{Plan: p}, resp)
	return resp
}

func storedSelfServiceUpgradeState(t *testing.T) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: AppServiceSelfServiceUpgradeSchema()}
	require.False(t, state.Set(context.Background(), testSelfServiceUpgradeState()).HasError())
	return state
}

func TestAppServiceSelfServiceUpgradeSchema(t *testing.T) {
	attrs := AppServiceSelfServiceUpgradeSchema().Attributes

	queueAt, ok := attrs["queue_at"].(*schema.StringAttribute)
	require.True(t, ok)
	assert.True(t, queueAt.Optional)
	assert.True(t, queueAt.Computed)

	versionId, ok := attrs["version_id"].(*schema.StringAttribute)
	require.True(t, ok)
	assert.True(t, versionId.Required)
	for _, m := range versionId.PlanModifiers {
		assert.NotContains(t, m.Description(context.Background()), "replace", "version_id must change in place")
	}

	audit, ok := attrs["audit"].(*schema.SingleNestedAttribute)
	require.True(t, ok)
	assert.True(t, audit.Computed)
}

func TestAppServiceSelfServiceUpgradeCreateRequestBody(t *testing.T) {
	immediate, err := json.Marshal(appservice_self_service_upgrade.CreateRequest{VersionId: "3.1.8"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"versionId":"3.1.8"}`, string(immediate), "an immediate upgrade has no queueAt key")

	queueAt := time.Date(2026, 8, 24, 2, 30, 0, 0, time.UTC)
	scheduled, err := json.Marshal(appservice_self_service_upgrade.CreateRequest{VersionId: "3.1.8", QueueAt: &queueAt})
	require.NoError(t, err)
	assert.JSONEq(t, `{"versionId":"3.1.8","queueAt":"2026-08-24T02:30:00Z"}`, string(scheduled))
}

func TestNormalizeQueueAt(t *testing.T) {
	sameInstant := time.Date(2026, 8, 24, 2, 30, 0, 0, time.UTC)
	moved := time.Date(2026, 8, 24, 3, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		prior    timetypes.RFC3339
		actual   *time.Time
		expected string
	}{
		{name: "same instant keeps the prior string", prior: timetypes.NewRFC3339ValueMust("2026-08-24T08:00:00+05:30"), actual: &sameInstant, expected: "2026-08-24T08:00:00+05:30"},
		{name: "moved time takes the API value", prior: timetypes.NewRFC3339ValueMust("2026-08-24T02:30:00Z"), actual: &moved, expected: "2026-08-24T03:00:00Z"},
		{name: "unknown prior takes the API value", prior: timetypes.NewRFC3339Unknown(), actual: &sameInstant, expected: "2026-08-24T02:30:00Z"},
		{name: "null prior takes the API value", prior: timetypes.NewRFC3339Null(), actual: &sameInstant, expected: "2026-08-24T02:30:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeQueueAt(tt.prior, tt.actual).ValueString())
		})
	}

	assert.True(t, normalizeQueueAt(timetypes.NewRFC3339Unknown(), nil).IsNull())
}

func TestIsTerminalSelfServiceUpgradeState(t *testing.T) {
	for state, expected := range map[string]bool{
		"completed": true,
		"failed":    true,
		"cancelled": true,
		"queued":    false,
		"pending":   false,
	} {
		assert.Equal(t, expected, isTerminalSelfServiceUpgradeState(state), state)
	}
}

func TestAppServiceSelfServiceUpgradeValidate(t *testing.T) {
	importId := "id=" + testUpgradeID + ",app_service_id=a,cluster_id=c,project_id=p,organization_id=o"
	imported := providerschema.AppServiceSelfServiceUpgrade{Id: types.StringValue(importId)}

	IDs, err := imported.Validate()

	require.NoError(t, err)
	assert.Equal(t, map[providerschema.Attr]string{
		providerschema.Id:             testUpgradeID,
		providerschema.AppServiceId:   "a",
		providerschema.ClusterId:      "c",
		providerschema.ProjectId:      "p",
		providerschema.OrganizationId: "o",
	}, IDs)

	missing := providerschema.AppServiceSelfServiceUpgrade{Id: types.StringValue("id=" + testUpgradeID + ",cluster_id=c,project_id=p,organization_id=o")}
	_, err = missing.Validate()
	assert.ErrorIs(t, err, errors.ErrValidatingResource)
}

func TestAppServiceSelfServiceUpgradeCreate(t *testing.T) {
	accepted := scriptedResponse{status: http.StatusAccepted, body: `{"id":"` + testUpgradeID + `"}`}

	tests := []struct {
		name          string
		queueAt       string
		gets          []scriptedResponse
		expectedPost  string
		expectedGets  int
		expectedError string
		expectedQueue string
	}{
		{
			name:          "scheduled upgrade keeps the planned queue_at string",
			queueAt:       "2026-08-24T02:30:00Z",
			gets:          []scriptedResponse{{http.StatusOK, testUpgradeBody("pending")}},
			expectedPost:  `{"versionId":"3.1.8","queueAt":"2026-08-24T02:30:00Z"}`,
			expectedGets:  1,
			expectedQueue: "2026-08-24T02:30:00Z",
		},
		{
			name:          "immediate upgrade sends no queueAt and takes it from the API",
			gets:          []scriptedResponse{{http.StatusOK, testUpgradeBody("queued")}},
			expectedPost:  `{"versionId":"3.1.8"}`,
			expectedGets:  1,
			expectedQueue: "2026-08-24T02:30:00Z",
		},
		{
			name:          "404 twice then the record",
			queueAt:       "2026-08-24T02:30:00Z",
			gets:          []scriptedResponse{{http.StatusNotFound, testNotFound}, {http.StatusNotFound, testNotFound}, {http.StatusOK, testUpgradeBody("pending")}},
			expectedPost:  `{"versionId":"3.1.8","queueAt":"2026-08-24T02:30:00Z"}`,
			expectedGets:  3,
			expectedQueue: "2026-08-24T02:30:00Z",
		},
		{
			name:          "404 until the timeout names the id",
			queueAt:       "2026-08-24T02:30:00Z",
			gets:          []scriptedResponse{{http.StatusNotFound, testNotFound}},
			expectedPost:  `{"versionId":"3.1.8","queueAt":"2026-08-24T02:30:00Z"}`,
			expectedError: testUpgradeID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fastSelfServiceUpgradePolling(t)
			r, fake := newTestSelfServiceUpgrade(t, map[string][]scriptedResponse{
				http.MethodPost: {accepted},
				http.MethodGet:  tt.gets,
			})

			resp := runSelfServiceUpgradeCreate(t, r, testSelfServiceUpgradePlan(tt.queueAt))

			require.Len(t, fake.bodies[http.MethodPost], 1)
			assert.JSONEq(t, tt.expectedPost, fake.bodies[http.MethodPost][0])

			var got providerschema.AppServiceSelfServiceUpgrade
			require.False(t, resp.State.Get(context.Background(), &got).HasError())
			assert.Equal(t, testUpgradeID, got.Id.ValueString(), "the accepted upgrade must stay tracked")

			if tt.expectedError != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.expectedError)
				return
			}
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			assert.Len(t, fake.calls(), 1+tt.expectedGets)
			assert.Equal(t, tt.expectedQueue, got.QueueAt.ValueString())
			assert.Equal(t, "3.1.8", got.VersionId.ValueString())
			assert.False(t, got.CurrentState.IsNull())
		})
	}
}

func TestAppServiceSelfServiceUpgradeRead(t *testing.T) {
	tests := []struct {
		name          string
		get           scriptedResponse
		expectRemoved bool
		expectError   bool
		expected      func(*testing.T, providerschema.AppServiceSelfServiceUpgrade)
	}{
		{
			name: "maps the response to state",
			get:  scriptedResponse{http.StatusOK, testUpgradeBody("pending")},
			expected: func(t *testing.T, got providerschema.AppServiceSelfServiceUpgrade) {
				assert.Equal(t, "pending", got.CurrentState.ValueString())
				assert.Equal(t, "2026-08-24T02:30:00Z", got.QueueAt.ValueString())
				assert.True(t, got.CompletedAt.IsNull())
			},
		},
		{
			name: "missing versionId keeps the prior version_id",
			get: scriptedResponse{http.StatusOK, `{"id":"upgrade-1","type":"upgrade","initiatedBy":"selfService",` +
				`"currentState":"pending","queueAt":"2026-08-24T02:30:00Z"}`},
			expected: func(t *testing.T, got providerschema.AppServiceSelfServiceUpgrade) {
				assert.Equal(t, "3.1.8", got.VersionId.ValueString())
			},
		},
		{
			name:          "404 removes the resource from state",
			get:           scriptedResponse{http.StatusNotFound, testNotFound},
			expectRemoved: true,
		},
		{
			name: "403 keeps the state and returns the error",
			get: scriptedResponse{http.StatusForbidden, `{"code":2,"message":"This feature is disabled for your organization.",` +
				`"httpStatusCode":403}`},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r, _ := newTestSelfServiceUpgrade(t, map[string][]scriptedResponse{http.MethodGet: {tt.get}})
			state := storedSelfServiceUpgradeState(t)
			resp := &resource.ReadResponse{State: state}

			r.Read(ctx, resource.ReadRequest{State: state}, resp)

			if tt.expectError {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "403")
				assert.False(t, resp.State.Raw.IsNull(), "state must be kept")
				return
			}
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			if tt.expectRemoved {
				assert.True(t, resp.State.Raw.IsNull())
				return
			}
			var got providerschema.AppServiceSelfServiceUpgrade
			require.False(t, resp.State.Get(ctx, &got).HasError())
			tt.expected(t, got)
		})
	}
}

func TestAppServiceSelfServiceUpgradeUpdate(t *testing.T) {
	updated := scriptedResponse{http.StatusOK, strings.Replace(testUpgradeBody("pending"), "3.1.8", "3.2.0", 1)}

	tests := []struct {
		name            string
		put             scriptedResponse
		get             scriptedResponse
		expectedCalls   []string
		expectError     bool
		expectedVersion string
	}{
		{
			name:            "replaces version and queue time then reads back",
			put:             scriptedResponse{http.StatusNoContent, ""},
			get:             updated,
			expectedCalls:   []string{http.MethodPut, http.MethodGet},
			expectedVersion: "3.2.0",
		},
		{
			name: "rejection of a started upgrade is returned without retry",
			put: scriptedResponse{http.StatusBadRequest, `{"code":11077,"message":"Only pending self-service upgrades can be updated or cancelled.",` +
				`"httpStatusCode":400}`},
			get:             updated,
			expectedCalls:   []string{http.MethodPut},
			expectError:     true,
			expectedVersion: "3.1.8",
		},
		{
			name:            "failed read after the update keeps the submitted version",
			put:             scriptedResponse{http.StatusNoContent, ""},
			get:             scriptedResponse{http.StatusInternalServerError, `{"code":500,"message":"boom","httpStatusCode":500}`},
			expectedCalls:   []string{http.MethodPut, http.MethodGet},
			expectError:     true,
			expectedVersion: "3.2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r, fake := newTestSelfServiceUpgrade(t, map[string][]scriptedResponse{
				http.MethodPut: {tt.put},
				http.MethodGet: {tt.get},
			})
			s := AppServiceSelfServiceUpgradeSchema()
			planModel := testSelfServiceUpgradeState()
			planModel.VersionId = types.StringValue("3.2.0")
			planModel.CurrentState = types.StringUnknown()
			plan := tfsdk.Plan{Schema: s}
			require.False(t, plan.Set(ctx, planModel).HasError())
			state := storedSelfServiceUpgradeState(t)
			resp := &resource.UpdateResponse{State: state}

			r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, resp)

			assert.Equal(t, tt.expectedCalls, fake.calls())
			require.Len(t, fake.bodies[http.MethodPut], 1)
			assert.JSONEq(t, `{"versionId":"3.2.0","queueAt":"2026-08-24T02:30:00Z"}`, fake.bodies[http.MethodPut][0])
			require.Equal(t, tt.expectError, resp.Diagnostics.HasError(), resp.Diagnostics)
			var got providerschema.AppServiceSelfServiceUpgrade
			require.False(t, resp.State.Get(ctx, &got).HasError())
			assert.Equal(t, tt.expectedVersion, got.VersionId.ValueString())
			assert.Equal(t, "2026-08-24T02:30:00Z", got.QueueAt.ValueString())
		})
	}
}

func TestAppServiceSelfServiceUpgradeDelete(t *testing.T) {
	pending := scriptedResponse{http.StatusOK, testUpgradeBody("pending")}
	notFound := scriptedResponse{http.StatusNotFound, testNotFound}

	tests := []struct {
		name          string
		gets          []scriptedResponse
		deletes       []scriptedResponse
		expectedCalls []string
		expectedError string
	}{
		{
			name:          "pending upgrade is cancelled",
			gets:          []scriptedResponse{pending},
			deletes:       []scriptedResponse{{http.StatusNoContent, ""}},
			expectedCalls: []string{http.MethodGet, http.MethodDelete},
		},
		{
			name:          "missing upgrade makes no cancel",
			gets:          []scriptedResponse{notFound},
			expectedCalls: []string{http.MethodGet},
		},
		{
			name:          "finished upgrade makes no cancel",
			gets:          []scriptedResponse{{http.StatusOK, testUpgradeBody("completed")}},
			expectedCalls: []string{http.MethodGet},
		},
		{
			name:          "cancel failure keeps the state",
			gets:          []scriptedResponse{pending},
			deletes:       []scriptedResponse{{http.StatusInternalServerError, `{"code":500,"message":"boom","httpStatusCode":500}`}},
			expectedCalls: []string{http.MethodGet, http.MethodDelete},
			expectedError: "Could not cancel self-service upgrade " + testUpgradeID,
		},
		{
			name:          "cancel 404 confirmed by a 404 read",
			gets:          []scriptedResponse{pending, notFound},
			deletes:       []scriptedResponse{notFound},
			expectedCalls: []string{http.MethodGet, http.MethodDelete, http.MethodGet},
		},
		{
			name:          "cancel 404 while the upgrade is still queued keeps the state",
			gets:          []scriptedResponse{pending, {http.StatusOK, testUpgradeBody("queued")}},
			deletes:       []scriptedResponse{notFound},
			expectedCalls: []string{http.MethodGet, http.MethodDelete, http.MethodGet},
			expectedError: testUpgradeID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, fake := newTestSelfServiceUpgrade(t, map[string][]scriptedResponse{
				http.MethodGet:    tt.gets,
				http.MethodDelete: tt.deletes,
			})
			state := storedSelfServiceUpgradeState(t)
			resp := &resource.DeleteResponse{State: state}

			r.Delete(context.Background(), resource.DeleteRequest{State: state}, resp)

			assert.Equal(t, tt.expectedCalls, fake.calls())
			if tt.expectedError != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.expectedError)
				return
			}
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
		})
	}
}
