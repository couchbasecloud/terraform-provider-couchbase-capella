package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testOrgID     = "org-1"
	testProjectID = "proj-1"
	testClusterID = "cluster-1"
)

// fakeBackend is an in-memory control-plane stand-in for the private endpoint
// service API. It records every request it receives and answers GET status
// polls from a scripted queue of responses (the final entry repeats once the
// queue is drained) so tests can drive the poll loop deterministically.
type fakeBackend struct {
	mu sync.Mutex

	// statuses is the queue of GET /privateEndpointService responses. The last
	// element is returned for any poll beyond the queue length.
	statuses []api.GetPrivateEndpointServiceStatusResponse

	// deleteStatus is the HTTP status returned for DELETE requests.
	deleteStatus int
	// deleteBody is the body returned for DELETE requests (used to simulate an
	// API error on cleanup).
	deleteBody string

	// recorded request methods, in order.
	methods []string
	getIdx  int
}

func (b *fakeBackend) counts() (gets, deletes int) {
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

func (b *fakeBackend) handler(w http.ResponseWriter, r *http.Request) {
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

// newTestResource wires a PrivateEndpointService to an httptest server backed
// by the supplied fakeBackend.
func newTestResource(t *testing.T, b *fakeBackend) *PrivateEndpointService {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(b.handler))
	t.Cleanup(srv.Close)

	return &PrivateEndpointService{
		Data: &providerschema.Data{
			ClientV1: &api.Client{Client: srv.Client()},
			HostURL:  srv.URL,
		},
	}
}

func strPtr(s string) *string { return &s }

// fastPolling shrinks the package-level poll timing for the duration of a test
// so the poll loops resolve in milliseconds rather than seconds.
func fastPolling(t *testing.T) {
	t.Helper()
	origPoll, origCleanup, origStatus := pollInterval, cleanupTimeout, statusChangeTimeout
	pollInterval = time.Millisecond
	cleanupTimeout = 2 * time.Second
	statusChangeTimeout = 2 * time.Second
	t.Cleanup(func() {
		pollInterval = origPoll
		cleanupTimeout = origCleanup
		statusChangeTimeout = origStatus
	})
}

func TestWaitUntilStatusChanges(t *testing.T) {
	fastPolling(t)

	tests := []struct {
		name       string
		finalState bool
		statuses   []api.GetPrivateEndpointServiceStatusResponse
		wantErr    error
	}{
		{
			// A terminal enableFailed that is present from the first poll and
			// never transitions cannot be told apart from a stale one while it
			// is happening, so the loop keeps polling. But on the overall
			// timeout we trust the persistent terminal status and surface the
			// typed failure (not a generic timeout) so the caller still routes
			// to cleanup / state removal instead of leaving orphaned infra.
			name:       "enableFailed stuck without any progress surfaces enableFailed on timeout",
			finalState: true,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: false, Status: strPtr(statusEnableFailed)},
			},
			wantErr: errors.ErrPrivateEndpointServiceEnableFailed,
		},
		{
			// Symmetric case on the disable path.
			name:       "disableFailed stuck without any progress surfaces disableFailed on timeout",
			finalState: false,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: true, Status: strPtr(statusDisableFailed)},
			},
			wantErr: errors.ErrPrivateEndpointServiceDisableFailed,
		},
		{
			// Once the backend shows progress (enabling) the earlier terminal
			// status was genuinely stale, so a later stall is a real transition
			// timeout — NOT a failure. This proves the deferred-failure latch is
			// cleared when sawInFlight flips.
			name:       "stale enableFailed then stuck enabling times out generically",
			finalState: true,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: false, Status: strPtr(statusEnableFailed)},
				{Enabled: false, Status: strPtr(statusEnabling)},
			},
			wantErr: errors.ErrPrivateEndpointServiceTimeout,
		},
		{
			name:       "transient enabling resolves to enabled",
			finalState: true,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: false, Status: strPtr(statusEnabling)},
				{Enabled: false, Status: strPtr(statusUnknown)},
				{Enabled: true, Status: strPtr(statusEnabled)},
			},
			wantErr: nil,
		},
		{
			name:       "transient disabling resolves to disabled",
			finalState: false,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: true, Status: strPtr(statusDisabling)},
				{Enabled: false, Status: strPtr(statusDisabled)},
			},
			wantErr: nil,
		},
		{
			name:       "older control plane falls back to enabled boolean",
			finalState: true,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: false, Status: nil},
				{Enabled: true, Status: nil},
			},
			wantErr: nil,
		},
		{
			name:       "resolved status but wrong boolean keeps polling until match",
			finalState: true,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				// Status says enabled but the boolean has not flipped yet, so the
				// loop must keep polling rather than return early.
				{Enabled: false, Status: strPtr(statusEnabled)},
				{Enabled: true, Status: strPtr(statusEnabled)},
			},
			wantErr: nil,
		},
		{
			// Reproduces the stale-state race: a prior failed attempt left the
			// backend in enableFailed; our POST has been accepted but the GET we
			// fire immediately afterward may still see the residual terminal
			// state before the new enable job is observed. The poll must wait
			// for evidence the current operation is in flight (a transient
			// state) before short-circuiting on a terminal status.
			name:       "stale enableFailed before current operation in flight resolves to enabled",
			finalState: true,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: false, Status: strPtr(statusEnableFailed)},
				{Enabled: false, Status: strPtr(statusEnabling)},
				{Enabled: true, Status: strPtr(statusEnabled)},
			},
			wantErr: nil,
		},
		{
			// Symmetric case on the disable path.
			name:       "stale disableFailed before current operation in flight resolves to disabled",
			finalState: false,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: true, Status: strPtr(statusDisableFailed)},
				{Enabled: true, Status: strPtr(statusDisabling)},
				{Enabled: false, Status: strPtr(statusDisabled)},
			},
			wantErr: nil,
		},
		{
			// After we have evidence the current operation is in flight, a
			// terminal status IS authoritative — the existing fail-fast
			// behavior must be preserved.
			name:       "enableFailed after transient enabling is terminal",
			finalState: true,
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: false, Status: strPtr(statusEnabling)},
				{Enabled: false, Status: strPtr(statusEnableFailed)},
			},
			wantErr: errors.ErrPrivateEndpointServiceEnableFailed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBackend{statuses: tc.statuses}
			p := newTestResource(t, b)

			err := p.waitUntilStatusChanges(context.Background(), tc.finalState, testOrgID, testProjectID, testClusterID)

			if tc.wantErr != nil {
				assert.Assert(t, stderrors.Is(err, tc.wantErr), "got %v, want %v", err, tc.wantErr)
				return
			}
			assert.NilError(t, err)
		})
	}
}

func TestWaitUntilCleanedUp(t *testing.T) {
	fastPolling(t)

	tests := []struct {
		name     string
		statuses []api.GetPrivateEndpointServiceStatusResponse
		wantErr  error
	}{
		{
			name: "reaches disabled after teardown",
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: true, Status: strPtr(statusEnableFailed)},
				{Enabled: false, Status: strPtr(statusDisabling)},
				{Enabled: false, Status: strPtr(statusDisabled)},
			},
			wantErr: nil,
		},
		{
			name: "teardown reports disableFailed",
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: true, Status: strPtr(statusDisableFailed)},
			},
			wantErr: errors.ErrPrivateEndpointServiceDisableFailed,
		},
		{
			name: "older control plane resolves via enabled boolean",
			statuses: []api.GetPrivateEndpointServiceStatusResponse{
				{Enabled: true, Status: nil},
				{Enabled: false, Status: nil},
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBackend{statuses: tc.statuses}
			p := newTestResource(t, b)

			err := p.waitUntilCleanedUp(context.Background(), testOrgID, testProjectID, testClusterID)

			if tc.wantErr != nil {
				assert.Assert(t, stderrors.Is(err, tc.wantErr), "got %v, want %v", err, tc.wantErr)
				return
			}
			assert.NilError(t, err)
		})
	}
}

func TestCleanupFailedEnableIssuesDelete(t *testing.T) {
	fastPolling(t)

	b := &fakeBackend{
		statuses: []api.GetPrivateEndpointServiceStatusResponse{
			{Enabled: false, Status: strPtr(statusDisabled)},
		},
	}
	p := newTestResource(t, b)

	err := p.cleanupFailedEnable(context.Background(), testOrgID, testProjectID, testClusterID)
	assert.NilError(t, err)

	_, deletes := b.counts()
	assert.Equal(t, deletes, 1, "cleanup must issue exactly one DELETE on the user's behalf")
}

func TestCleanupFailedEnablePropagatesDeleteError(t *testing.T) {
	fastPolling(t)

	b := &fakeBackend{
		deleteStatus: http.StatusInternalServerError,
		deleteBody:   `{"code":4000,"message":"boom","httpStatusCode":500}`,
	}
	p := newTestResource(t, b)

	err := p.cleanupFailedEnable(context.Background(), testOrgID, testProjectID, testClusterID)
	assert.Assert(t, err != nil, "a failed DELETE must surface an error")

	gets, _ := b.counts()
	assert.Equal(t, gets, 0, "must not poll for cleanup completion when the DELETE itself failed")
}

func TestHandleFailedEnableCleansUpAndRemovesState(t *testing.T) {
	fastPolling(t)

	b := &fakeBackend{
		statuses: []api.GetPrivateEndpointServiceStatusResponse{
			{Enabled: false, Status: strPtr(statusDisabled)},
		},
	}
	p := newTestResource(t, b)

	state := &tfsdk.State{Schema: PrivateEndpointServiceSchema()}
	var diags diag.Diagnostics

	p.handleFailedEnable(context.Background(), state, &diags,
		testOrgID, testProjectID, testClusterID, errors.ErrPrivateEndpointServiceEnableFailed)

	_, deletes := b.counts()
	assert.Equal(t, deletes, 1, "must issue a DELETE to clean up the failed enable")
	assert.Assert(t, state.Raw.IsNull(), "resource must be removed from state for a clean re-create")
	assert.Equal(t, diags.HasError(), true, "must surface an actionable error")
	assert.Equal(t, diags.Errors()[0].Summary(), "Private endpoint service enablement failed")
}

func TestHandleFailedEnableRemovesStateEvenWhenCleanupFails(t *testing.T) {
	fastPolling(t)

	// DELETE fails, so automatic cleanup cannot complete.
	b := &fakeBackend{
		deleteStatus: http.StatusInternalServerError,
		deleteBody:   `{"code":4000,"message":"boom","httpStatusCode":500}`,
	}
	p := newTestResource(t, b)

	state := &tfsdk.State{Schema: PrivateEndpointServiceSchema()}
	var diags diag.Diagnostics

	p.handleFailedEnable(context.Background(), state, &diags,
		testOrgID, testProjectID, testClusterID, errors.ErrPrivateEndpointServiceEnableFailed)

	assert.Assert(t, state.Raw.IsNull(), "state must be removed even when cleanup fails to avoid a wedged resource")
	assert.Equal(t, diags.HasError(), true)
	assert.Equal(t, diags.Errors()[0].Summary(),
		"Private endpoint service enablement failed and automatic cleanup did not complete")
}

func TestGetServiceState(t *testing.T) {
	tests := []struct {
		name        string
		status      api.GetPrivateEndpointServiceStatusResponse
		wantEnabled bool
		wantStatus  string
		wantNull    bool
	}{
		{
			name:        "maps status when present",
			status:      api.GetPrivateEndpointServiceStatusResponse{Enabled: true, Status: strPtr(statusEnabled)},
			wantEnabled: true,
			wantStatus:  statusEnabled,
		},
		{
			name:        "leaves status null for older control plane",
			status:      api.GetPrivateEndpointServiceStatusResponse{Enabled: true, Status: nil},
			wantEnabled: true,
			wantNull:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBackend{statuses: []api.GetPrivateEndpointServiceStatusResponse{tc.status}}
			p := newTestResource(t, b)

			got, err := p.getServiceState(context.Background(), testOrgID, testProjectID, testClusterID)
			assert.NilError(t, err)
			assert.Equal(t, got.Enabled.ValueBool(), tc.wantEnabled)
			assert.Equal(t, got.Status.IsNull(), tc.wantNull)
			if !tc.wantNull {
				assert.Equal(t, got.Status.ValueString(), tc.wantStatus)
			}
		})
	}
}
