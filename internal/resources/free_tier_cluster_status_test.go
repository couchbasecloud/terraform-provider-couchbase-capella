package resources

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// freeTierClusterStatusBackend is an in-memory control-plane stand-in that
// answers GET free tier cluster status polls from a scripted queue of states
// (the last entry repeats once the queue is drained) so tests can drive
// checkForFreeTierClusterStatus's poll loop deterministically.
type freeTierClusterStatusBackend struct {
	mu     sync.Mutex
	states []clusterapi.State
	idx    int
}

func (b *freeTierClusterStatusBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	idx := b.idx
	if idx >= len(b.states) {
		idx = len(b.states) - 1
	}
	b.idx++

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(clusterapi.GetClusterResponse{CurrentState: b.states[idx]})
}

// fastFreeTierClusterPolling shrinks the package-level poll timing for the
// duration of a test so checkForFreeTierClusterStatus's poll loop resolves in
// milliseconds.
func fastFreeTierClusterPolling(t *testing.T) {
	t.Helper()
	origInterval, origTimeout := freeTierClusterStatusPollInterval, freeTierClusterStatusTimeout
	freeTierClusterStatusPollInterval = time.Millisecond
	freeTierClusterStatusTimeout = 2 * time.Second
	t.Cleanup(func() {
		freeTierClusterStatusPollInterval = origInterval
		freeTierClusterStatusTimeout = origTimeout
	})
}

func newTestFreeTierCluster(t *testing.T, handler http.HandlerFunc) *FreeTierCluster {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return &FreeTierCluster{
		Data: &providerschema.Data{
			ClientV1: &api.Client{Client: srv.Client()},
			HostURL:  srv.URL,
		},
	}
}

func TestCheckForFreeTierClusterStatus(t *testing.T) {
	fastFreeTierClusterPolling(t)

	tests := []struct {
		name        string
		states      []clusterapi.State
		wantErr     bool
		errContains string
	}{
		{
			name:   "reaches healthy after transient deploying states",
			states: []clusterapi.State{clusterapi.Deploying, clusterapi.Deploying, clusterapi.Healthy},
		},
		{
			name:        "deploymentFailed surfaces an error instead of success",
			states:      []clusterapi.State{clusterapi.Deploying, clusterapi.DeploymentFailed},
			wantErr:     true,
			errContains: "deploymentFailed",
		},
		{
			name:        "destroyFailed surfaces an error",
			states:      []clusterapi.State{clusterapi.Destroying, clusterapi.DestroyFailed},
			wantErr:     true,
			errContains: "destroyFailed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := &freeTierClusterStatusBackend{states: tt.states}
			f := newTestFreeTierCluster(t, backend.handler)

			resp, err := f.checkForFreeTierClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.errContains)
				assert.Assert(t, resp == nil, "expected a nil response on failure")
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, resp.CurrentState, clusterapi.Healthy)
		})
	}
}

func TestCheckForFreeTierClusterStatus_TimesOutWithoutReachingTerminalState(t *testing.T) {
	fastFreeTierClusterPolling(t)

	backend := &freeTierClusterStatusBackend{states: []clusterapi.State{clusterapi.Deploying}}
	f := newTestFreeTierCluster(t, backend.handler)

	_, err := f.checkForFreeTierClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

	assert.ErrorContains(t, err, "timed out")
}

func TestCheckForFreeTierClusterStatus_ReturnsImmediatelyOnGetError(t *testing.T) {
	fastFreeTierClusterPolling(t)

	var mu sync.Mutex
	calls := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		mu.Unlock()
		w.WriteHeader(http.StatusInternalServerError)
	}

	f := newTestFreeTierCluster(t, handler)

	_, err := f.checkForFreeTierClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

	assert.Assert(t, err != nil, "expected an error from a failing GET")
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, calls, 1, "expected checkForFreeTierClusterStatus to return on the first GET error, not retry")
}

func TestCheckForFreeTierClusterStatus_UsesFreeTierGetClusterEndpoint(t *testing.T) {
	fastFreeTierClusterPolling(t)

	var gotPath string
	handler := func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(clusterapi.GetClusterResponse{CurrentState: clusterapi.Healthy})
	}

	f := newTestFreeTierCluster(t, handler)

	_, err := f.checkForFreeTierClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(gotPath, "/organizations/org-1/projects/proj-1/clusters/freeTier/cluster-1"), "unexpected path: %s", gotPath)
}
