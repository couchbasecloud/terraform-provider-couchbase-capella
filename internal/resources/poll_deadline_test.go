package resources

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// stalledHandler never answers. It returns when the client cancels the request,
// or after a guard so a regression fails the test instead of hanging it.
func stalledHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
	case <-time.After(10 * time.Second):
	}
}

// These pollers return request errors. When their deadline passes while a GET is in
// flight, the cancelled request must still report the poller's own timeout error.

func TestCheckForFreeTierClusterStatus_DeadlineDuringRequest(t *testing.T) {
	fastFreeTierClusterPolling(t)
	f := newTestFreeTierCluster(t, stalledHandler)

	start := time.Now()
	_, err := f.checkForFreeTierClusterStatus(t.Context(), "org-1", "proj-1", "cluster-1")

	require.Error(t, err)
	assert.ErrorContains(t, err, "free tier cluster creation status transition timed out")
	assert.Less(t, time.Since(start), 5*time.Second)
}

func TestWaitUntilCleanedUp_DeadlineDuringRequest(t *testing.T) {
	fastPolling(t)
	p := &PrivateEndpointService{Data: newTestProviderData(t, stalledHandler)}

	start := time.Now()
	err := p.waitUntilCleanedUp(t.Context(), "org-1", "proj-1", "cluster-1")

	assert.ErrorIs(t, err, errors.ErrPrivateEndpointServiceTimeout)
	assert.Less(t, time.Since(start), 5*time.Second)
}

func TestWaitUntilLinked_DeadlineDuringRequest(t *testing.T) {
	fastPolling(t)
	a := &AppServicePrivateEndpoint{Data: newTestProviderData(t, stalledHandler)}

	start := time.Now()
	_, err := a.waitUntilLinked(t.Context(), "org-1", "proj-1", "cluster-1", "app-1", "endpoint-1")

	require.Error(t, err)
	assert.ErrorContains(t, err, "waiting for status linked")
	assert.Less(t, time.Since(start), 5*time.Second)
}
