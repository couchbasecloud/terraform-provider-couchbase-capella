package resources

import (
	stderrors "errors"
	"net/http"
	"testing"

	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

const testAppServiceID = "appservice-1"

// listEndpointsHandler answers every request with the next body in the queue
// and repeats the last body once the queue is drained.
func listEndpointsHandler(bodies ...string) http.HandlerFunc {
	var i int
	return func(w http.ResponseWriter, _ *http.Request) {
		if i >= len(bodies) {
			i = len(bodies) - 1
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(bodies[i]))
		i++
	}
}

func endpointListBody(status string) string {
	return `{"endpoints":[{"id":"vpce-1","serviceName":"svc","status":"` + status + `"}]}`
}

func TestRefreshAppServicePrivateEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus string
		wantErr    error
	}{
		{name: "found", body: endpointListBody("linked"), wantStatus: "linked"},
		{name: "not found", body: `{"endpoints":[]}`, wantErr: errors.ErrNotFound},
		{name: "unmarshal error", body: `not json`, wantErr: errors.ErrUnmarshallingResponse},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AppServicePrivateEndpoint{Data: newTestProviderData(t, listEndpointsHandler(tt.body))}
			got, err := r.refreshAppServicePrivateEndpoint(t.Context(), testOrgID, testProjectID, testClusterID, testAppServiceID, "vpce-1")
			if tt.wantErr != nil {
				assert.Assert(t, stderrors.Is(err, tt.wantErr), "got %v, want %v", err, tt.wantErr)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, got.Status.ValueString(), tt.wantStatus)
			assert.Equal(t, got.ServiceName.ValueString(), "svc")
		})
	}
}

func TestWaitUntilLinked(t *testing.T) {
	fastPolling(t)

	tests := []struct {
		name     string
		statuses []string
		wantErr  bool
	}{
		{name: "pending then linked", statuses: []string{"pendingAcceptance", "pending", "linked"}},
		{name: "rejected fails fast", statuses: []string{"rejected"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodies := make([]string, len(tt.statuses))
			for i, s := range tt.statuses {
				bodies[i] = endpointListBody(s)
			}
			r := &AppServicePrivateEndpoint{Data: newTestProviderData(t, listEndpointsHandler(bodies...))}
			got, err := r.waitUntilLinked(t.Context(), testOrgID, testProjectID, testClusterID, testAppServiceID, "vpce-1")
			if tt.wantErr {
				assert.Assert(t, err != nil)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, got.Status.ValueString(), "linked")
		})
	}
}

func TestIsPrivateEndpointNotEnabled(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":400,"hint":"client error","httpStatusCode":400,"message":"Private endpoints aren't enabled for this cluster. Please select another cluster."}`))
	}
	r := &AppServicePrivateEndpoint{Data: newTestProviderData(t, handler)}
	_, err := r.refreshAppServicePrivateEndpoint(t.Context(), testOrgID, testProjectID, testClusterID, testAppServiceID, "vpce-1")
	assert.Assert(t, isPrivateEndpointNotEnabled(err), "got %v", err)
}
