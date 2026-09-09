package resources

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// newTestProviderData starts an httptest server backed by handler and returns the
// provider data a resource needs to talk to it. The server is closed on test cleanup.
func newTestProviderData(t *testing.T, handler http.HandlerFunc) *providerschema.Data {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return &providerschema.Data{
		ClientV1: &api.Client{Client: srv.Client()},
		HostURL:  srv.URL,
	}
}
