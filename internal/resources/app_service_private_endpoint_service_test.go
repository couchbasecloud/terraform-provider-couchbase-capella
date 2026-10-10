package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const testAppServiceID = "app-service-1"

// appServicePEFakeBackend is an in-memory control-plane stand-in for the App Service
// private endpoint service API. It records every request it receives and answers GET
// status polls from a scripted queue of responses (the final entry repeats once the
// queue is drained) so tests can drive the poll loop deterministically.
type appServicePEFakeBackend struct {
	mu sync.Mutex

	// statuses is the queue of GET .../privateEndpointService responses. The last
	// element is returned for any poll beyond the queue length.
	statuses []api.GetAppServicePrivateEndpointServiceStatusResponse

	// deleteStatus is the HTTP status returned for DELETE requests.
	deleteStatus int
	// deleteBody is the body returned for DELETE requests (used to simulate an
	// API error on cleanup).
	deleteBody string

	// recorded request methods, in order.
	methods []string
	getIdx  int
}

func (b *appServicePEFakeBackend) counts() (gets, deletes int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, m := range b.methods {
		switch m {
		case http.MethodGet:
			gets++
		case http.MethodDelete:
			deletes++
		}
	}
	return gets, deletes
}

func (b *appServicePEFakeBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.methods = append(b.methods, r.Method)

	switch r.Method {
	case http.MethodDelete:
		status := b.deleteStatus
		if status == 0 {
			status = http.StatusAccepted
		}
		w.WriteHeader(status)
		if b.deleteBody != "" {
			_, _ = w.Write([]byte(b.deleteBody))
		}
	case http.MethodGet:
		idx := b.getIdx
		if idx >= len(b.statuses) {
			idx = len(b.statuses) - 1
		}
		b.getIdx++
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(b.statuses[idx])
	default:
		// Treat POST (enable) as accepted.
		w.WriteHeader(http.StatusAccepted)
	}
}

// newAppServicePETestResource wires an AppServicePrivateEndpointService to an httptest
// server backed by the supplied appServicePEFakeBackend.
func newAppServicePETestResource(t *testing.T, b *appServicePEFakeBackend) *AppServicePrivateEndpointService {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(b.handler))
	t.Cleanup(srv.Close)

	return &AppServicePrivateEndpointService{
		Data: &providerschema.Data{
			ClientV1: &api.Client{Client: srv.Client()},
			HostURL:  srv.URL,
		},
	}
}

func TestAppServicePrivateEndpointServiceWaitUntilStatusChanges(t *testing.T) {
	fastPolling(t)

	tests := []struct {
		name       string
		finalState bool
		statuses   []api.GetAppServicePrivateEndpointServiceStatusResponse
		wantErr    error
	}{
		{
			// A terminal failed state present from the first poll and never
			// transitioning cannot be told apart from a stale one while it is
			// happening, so the loop keeps polling. But on the overall timeout we
			// trust the persistent terminal status and surface the typed failure
			// (not a generic timeout) so the caller still routes to cleanup / state
			// removal instead of leaving orphaned infra.
			name:       "failed with targetState enabled stuck without progress surfaces enable failed on timeout",
			finalState: true,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusFailed), TargetState: strPtr(statusEnabled)},
			},
			wantErr: errors.ErrAppServicePrivateEndpointServiceEnableFailed,
		},
		{
			// Symmetric case on the disable path.
			name:       "failed with targetState disabled stuck without progress surfaces disable failed on timeout",
			finalState: false,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusFailed), TargetState: strPtr(statusDisabled)},
			},
			wantErr: errors.ErrAppServicePrivateEndpointServiceDisableFailed,
		},
		{
			// Once the backend shows progress (enabling) the earlier terminal
			// status was genuinely stale, so a later stall is a real transition
			// timeout — NOT a failure.
			name:       "stale failed then stuck enabling times out generically",
			finalState: true,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusFailed), TargetState: strPtr(statusEnabled)},
				{State: strPtr(statusEnabling)},
			},
			wantErr: errors.ErrAppServicePrivateEndpointServiceTimeout,
		},
		{
			name:       "transient enabling resolves to enabled",
			finalState: true,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusEnabling)},
				{State: strPtr(statusEnabled)},
			},
			wantErr: nil,
		},
		{
			name:       "transient disabling resolves to disabled",
			finalState: false,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusDisabling)},
				{State: strPtr(statusDisabled)},
			},
			wantErr: nil,
		},
		{
			name:       "resolved state but wrong direction keeps polling until match",
			finalState: true,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusDisabled)},
				{State: strPtr(statusEnabled)},
			},
			wantErr: nil,
		},
		{
			// Reproduces the stale-state race: a prior failed attempt left the
			// backend in a failed state; our POST has been accepted but the GET we
			// fire immediately afterward may still see the residual terminal state
			// before the new enable job is observed.
			name:       "stale failed before current operation in flight resolves to enabled",
			finalState: true,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusFailed), TargetState: strPtr(statusEnabled)},
				{State: strPtr(statusEnabling)},
				{State: strPtr(statusEnabled)},
			},
			wantErr: nil,
		},
		{
			// After we have evidence the current operation is in flight, a
			// terminal status IS authoritative.
			name:       "failed after transient enabling is terminal",
			finalState: true,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusEnabling)},
				{State: strPtr(statusFailed), TargetState: strPtr(statusEnabled)},
			},
			wantErr: errors.ErrAppServicePrivateEndpointServiceEnableFailed,
		},
		{
			name:       "absent state keeps polling until timeout",
			finalState: true,
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{},
			},
			wantErr: errors.ErrAppServicePrivateEndpointServiceTimeout,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &appServicePEFakeBackend{statuses: tc.statuses}
			p := newAppServicePETestResource(t, b)

			err := p.waitUntilStatusChanges(context.Background(), tc.finalState, testOrgID, testProjectID, testClusterID, testAppServiceID)

			if tc.wantErr != nil {
				assert.Assert(t, stderrors.Is(err, tc.wantErr), "got %v, want %v", err, tc.wantErr)
				return
			}
			assert.NilError(t, err)
		})
	}
}

func TestAppServicePrivateEndpointServiceWaitUntilCleanedUp(t *testing.T) {
	fastPolling(t)

	tests := []struct {
		name     string
		statuses []api.GetAppServicePrivateEndpointServiceStatusResponse
		wantErr  error
	}{
		{
			name: "reaches disabled after teardown",
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusFailed), TargetState: strPtr(statusEnabled)},
				{State: strPtr(statusDisabling)},
				{State: strPtr(statusDisabled)},
			},
			wantErr: nil,
		},
		{
			name: "teardown reports failed",
			statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
				{State: strPtr(statusFailed), TargetState: strPtr(statusDisabled)},
			},
			wantErr: errors.ErrAppServicePrivateEndpointServiceDisableFailed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &appServicePEFakeBackend{statuses: tc.statuses}
			p := newAppServicePETestResource(t, b)

			err := p.waitUntilCleanedUp(context.Background(), testOrgID, testProjectID, testClusterID, testAppServiceID)

			if tc.wantErr != nil {
				assert.Assert(t, stderrors.Is(err, tc.wantErr), "got %v, want %v", err, tc.wantErr)
				return
			}
			assert.NilError(t, err)
		})
	}
}

func TestAppServicePrivateEndpointServiceCleanupFailedEnableIssuesDelete(t *testing.T) {
	fastPolling(t)

	b := &appServicePEFakeBackend{
		statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
			{State: strPtr(statusDisabled)},
		},
	}
	p := newAppServicePETestResource(t, b)

	err := p.cleanupFailedEnable(context.Background(), testOrgID, testProjectID, testClusterID, testAppServiceID)
	assert.NilError(t, err)

	_, deletes := b.counts()
	assert.Equal(t, deletes, 1, "cleanup must issue exactly one DELETE on the user's behalf")
}

func TestAppServicePrivateEndpointServiceCleanupFailedEnablePropagatesDeleteError(t *testing.T) {
	fastPolling(t)

	b := &appServicePEFakeBackend{
		deleteStatus: http.StatusInternalServerError,
		deleteBody:   `{"code":4000,"message":"boom","httpStatusCode":500}`,
	}
	p := newAppServicePETestResource(t, b)

	err := p.cleanupFailedEnable(context.Background(), testOrgID, testProjectID, testClusterID, testAppServiceID)
	assert.Assert(t, err != nil, "a failed DELETE must surface an error")

	gets, _ := b.counts()
	assert.Equal(t, gets, 0, "must not poll for cleanup completion when the DELETE itself failed")
}

func TestAppServicePrivateEndpointServiceHandleFailedEnableCleansUpAndRemovesState(t *testing.T) {
	fastPolling(t)

	b := &appServicePEFakeBackend{
		statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{
			{State: strPtr(statusDisabled)},
		},
	}
	p := newAppServicePETestResource(t, b)

	state := &tfsdk.State{Schema: AppServicePrivateEndpointServiceSchema()}
	var diags diag.Diagnostics

	p.handleFailedEnable(context.Background(), state, &diags,
		testOrgID, testProjectID, testClusterID, testAppServiceID, errors.ErrAppServicePrivateEndpointServiceEnableFailed)

	_, deletes := b.counts()
	assert.Equal(t, deletes, 1, "must issue a DELETE to clean up the failed enable")
	assert.Assert(t, state.Raw.IsNull(), "resource must be removed from state for a clean re-create")
	assert.Equal(t, diags.HasError(), true, "must surface an actionable error")
	assert.Equal(t, diags.Errors()[0].Summary(), "App Service private endpoint service enablement failed")
}

func TestAppServicePrivateEndpointServiceHandleFailedEnableRemovesStateEvenWhenCleanupFails(t *testing.T) {
	fastPolling(t)

	// DELETE fails, so automatic cleanup cannot complete.
	b := &appServicePEFakeBackend{
		deleteStatus: http.StatusInternalServerError,
		deleteBody:   `{"code":4000,"message":"boom","httpStatusCode":500}`,
	}
	p := newAppServicePETestResource(t, b)

	state := &tfsdk.State{Schema: AppServicePrivateEndpointServiceSchema()}
	var diags diag.Diagnostics

	p.handleFailedEnable(context.Background(), state, &diags,
		testOrgID, testProjectID, testClusterID, testAppServiceID, errors.ErrAppServicePrivateEndpointServiceEnableFailed)

	assert.Assert(t, state.Raw.IsNull(), "state must be removed even when cleanup fails to avoid a wedged resource")
	assert.Equal(t, diags.HasError(), true)
	assert.Equal(t, diags.Errors()[0].Summary(),
		"App Service private endpoint service enablement failed and automatic cleanup did not complete")
}

func TestAppServicePrivateEndpointServiceGetServiceState(t *testing.T) {
	tests := []struct {
		name            string
		status          api.GetAppServicePrivateEndpointServiceStatusResponse
		wantEnabled     bool
		wantState       string
		wantTargetState string
		wantStateNull   bool
	}{
		{
			name:        "maps state and target state when present",
			status:      api.GetAppServicePrivateEndpointServiceStatusResponse{State: strPtr(statusEnabled), TargetState: strPtr(statusEnabled)},
			wantEnabled: true,
			wantState:   statusEnabled,
		},
		{
			name:          "leaves state null when absent",
			status:        api.GetAppServicePrivateEndpointServiceStatusResponse{},
			wantEnabled:   false,
			wantStateNull: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &appServicePEFakeBackend{statuses: []api.GetAppServicePrivateEndpointServiceStatusResponse{tc.status}}
			p := newAppServicePETestResource(t, b)

			got, err := p.getServiceState(context.Background(), testOrgID, testProjectID, testClusterID, testAppServiceID)
			assert.NilError(t, err)
			assert.Equal(t, got.Enabled.ValueBool(), tc.wantEnabled)
			assert.Equal(t, got.State.IsNull(), tc.wantStateNull)
			if !tc.wantStateNull {
				assert.Equal(t, got.State.ValueString(), tc.wantState)
			}
		})
	}
}
