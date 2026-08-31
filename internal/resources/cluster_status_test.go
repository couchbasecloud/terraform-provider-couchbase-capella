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

// clusterStatusBackend is an in-memory control-plane stand-in that answers GET
// cluster status polls from a scripted queue of states (the last entry repeats
// once the queue is drained) so tests can drive checkClusterStatus's poll loop
// deterministically.
type clusterStatusBackend struct {
	mu     sync.Mutex
	states []clusterapi.State
	idx    int
}

func (b *clusterStatusBackend) handler(w http.ResponseWriter, r *http.Request) {
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

// fastClusterPolling shrinks the package-level poll timing for the duration of
// a test so checkClusterStatus's poll loop resolves in milliseconds.
func fastClusterPolling(t *testing.T) {
	t.Helper()
	origInterval, origTimeout := clusterStatusPollInterval, clusterStatusTimeout
	clusterStatusPollInterval = time.Millisecond
	clusterStatusTimeout = 2 * time.Second
	t.Cleanup(func() {
		clusterStatusPollInterval = origInterval
		clusterStatusTimeout = origTimeout
	})
}

func newTestCluster(t *testing.T, handler http.HandlerFunc) *Cluster {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return &Cluster{
		Data: &providerschema.Data{
			ClientV1: &api.Client{Client: srv.Client()},
			HostURL:  srv.URL,
		},
	}
}

func TestCheckClusterStatus(t *testing.T) {
	fastClusterPolling(t)

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
			name:        "scaleFailed surfaces an error",
			states:      []clusterapi.State{clusterapi.Scaling, clusterapi.ScaleFailed},
			wantErr:     true,
			errContains: "scaleFailed",
		},
		{
			name:        "rebalanceFailed surfaces an error",
			states:      []clusterapi.State{clusterapi.RebalanceFailed},
			wantErr:     true,
			errContains: "rebalanceFailed",
		},
		{
			name:        "upgradeFailed surfaces an error",
			states:      []clusterapi.State{clusterapi.UpgradeFailed},
			wantErr:     true,
			errContains: "upgradeFailed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := &clusterStatusBackend{states: tt.states}
			c := newTestCluster(t, backend.handler)

			err := c.checkClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.errContains)
				return
			}
			assert.NilError(t, err)
		})
	}
}

func TestCheckClusterStatus_TimesOutWithoutReachingTerminalState(t *testing.T) {
	fastClusterPolling(t)

	backend := &clusterStatusBackend{states: []clusterapi.State{clusterapi.Deploying}}
	c := newTestCluster(t, backend.handler)

	err := c.checkClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

	assert.ErrorContains(t, err, "timed out")
}

func TestCheckClusterStatus_RetriesOnTransientGetErrors(t *testing.T) {
	fastClusterPolling(t)

	var mu sync.Mutex
	calls := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()

		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(clusterapi.GetClusterResponse{CurrentState: clusterapi.Healthy})
	}

	c := newTestCluster(t, handler)

	err := c.checkClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

	assert.NilError(t, err)
	mu.Lock()
	defer mu.Unlock()
	assert.Assert(t, calls >= 3, "expected checkClusterStatus to retry past the transient errors, got %d call(s)", calls)
}

func TestCheckClusterStatus_UsesGetClusterEndpoint(t *testing.T) {
	fastClusterPolling(t)

	var gotPath string
	handler := func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(clusterapi.GetClusterResponse{CurrentState: clusterapi.Healthy})
	}

	c := newTestCluster(t, handler)

	err := c.checkClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(gotPath, "/organizations/org-1/projects/proj-1/clusters/cluster-1"), "unexpected path: %s", gotPath)
}
