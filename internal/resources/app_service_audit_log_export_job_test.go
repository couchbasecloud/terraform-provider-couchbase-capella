package resources_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/provider"
)

// exportJobStub is a fake export jobs API. POST creates job-1; GET job-1 returns the next status in statuses and
// repeats the last one once they run out.
type exportJobStub struct {
	mu       sync.Mutex
	statuses []string
	gets     int
}

func (s *exportJobStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/auditLogExports"):
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"exportId":"job-1"}`))
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/auditLogExports/job-1"):
		status := s.statuses[min(s.gets, len(s.statuses)-1)]
		s.gets++
		_, _ = fmt.Fprintf(w, `{"exportId":"job-1","status":%q}`, status)
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":404,"httpStatusCode":404,"message":"not found"}`))
	}
}

// TestAppServiceAuditLogExportJobCreateWaitsForFailedJob runs Terraform against a fake API whose export job is queued
// when first read and failed on every read after that. Apply should wait for a final status and report the failure.
// Today Create reads the job once, stores status "queued" and reports success, and because status is computed the
// later failure never shows in a plan either (AV-153195).
func TestAppServiceAuditLogExportJobCreateWaitsForFailedJob(t *testing.T) {
	t.Skip("AV-153195: apply reports success for an export job that later fails; unskip once Create waits for the job")

	stub := &exportJobStub{statuses: []string{"queued", "failed"}}
	server := httptest.NewServer(stub)
	defer server.Close()

	const id = "ffffffff-aaaa-1414-eeee-000000000000"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"couchbase-capella": providerserver.NewProtocol6WithError(provider.New()()),
		},
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
provider "couchbase-capella" {
  host                 = %q
  authentication_token = "test-token"
}

resource "couchbase-capella_app_service_audit_log_export_job" "job" {
  organization_id = %[2]q
  project_id      = %[2]q
  cluster_id      = %[2]q
  app_service_id  = %[2]q
  start           = "2026-11-23T00:00:00Z"
  end             = "2026-11-23T04:00:00Z"
}
`, server.URL, id),
			ExpectError: regexp.MustCompile(`(?i)failed`),
		}},
	})
}
