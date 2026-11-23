package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// TestAccAppServiceAuditLogExportJobNoAuditData exports a window with no audit data. Audit logging is never enabled
// on the shared App Service, so every window is empty. Unlike the cluster export API, which answers GET with 404 for
// such a job (AV-128951), this API returns the job with status "No Log Files In Requested Time Range". The second
// step guards that the job stays in state once it has run: a 404 there would make Read drop it and every apply
// create a new job. Each run leaves an export job behind, because export jobs have no delete API.
func TestAccAppServiceAuditLogExportJobNoAuditData(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_as_audit_export_")
	resourceReference := "couchbase-capella_app_service_audit_log_export_job." + name

	end := time.Now().UTC().Add(-1 * time.Hour).Truncate(time.Second)
	start := end.Add(-4 * time.Hour)
	config := fmt.Sprintf(`
%s
resource "couchbase-capella_app_service_audit_log_export_job" "%s" {
  organization_id = "%s"
  project_id      = "%s"
  cluster_id      = "%s"
  app_service_id  = "%s"
  start           = "%s"
  end             = "%s"
}`, globalProviderBlock, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId,
		start.Format(time.RFC3339), end.Format(time.RFC3339))

	var jobId string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "status"),
					func(s *terraform.State) error {
						jobId = s.RootModule().Resources[resourceReference].Primary.ID
						return nil
					},
				),
			},
			// Let the job run, then refresh and plan. A 404 after the job has run drops it from state.
			{
				PreConfig: func() { waitForAppServiceAuditLogExportJobToRun(t, jobId) },
				Config:    config,
				PlanOnly:  true,
			},
		},
	})
}

// waitForAppServiceAuditLogExportJobToRun polls the export job until it reaches a final status or GET returns 404,
// logging every status it sees. It gives up after 20 minutes and lets the plan step run regardless.
func waitForAppServiceAuditLogExportJobToRun(t *testing.T, jobId string) {
	t.Helper()
	data := newTestClient(t)
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogExports/%s",
		globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, jobId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	lastStatus := ""
	deadline := time.Now().Add(20 * time.Minute)
	for time.Now().Before(deadline) {
		response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
		if err != nil {
			if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
				t.Logf("export job %s: GET returned 404 after status %q: %v", jobId, lastStatus, err)
				return
			}
			t.Fatalf("export job %s: GET failed: %v", jobId, err)
		}

		var job apigen.GetAuditExportDocResponse
		if err := json.Unmarshal(response.Body, &job); err != nil {
			t.Fatalf("export job %s: unmarshalling GET response: %v", jobId, err)
		}
		status := ""
		if job.Status != nil {
			status = *job.Status
		}
		if status != lastStatus {
			t.Logf("export job %s: status %q", jobId, status)
			lastStatus = status
		}
		// Final statuses are free text, such as "No Log Files In Requested Time Range", so wait only while pending.
		if !isPendingAppServiceAuditLogExportStatus(status) {
			return
		}
		time.Sleep(15 * time.Second)
	}
	t.Logf("export job %s: still %q after 20 minutes; planning anyway", jobId, lastStatus)
}

func isPendingAppServiceAuditLogExportStatus(status string) bool {
	s := strings.ToLower(status)
	return s == "" || strings.Contains(s, "queued") || strings.Contains(s, "progress") ||
		strings.Contains(s, "pending") || strings.Contains(s, "processing")
}
