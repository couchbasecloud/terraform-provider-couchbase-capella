package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
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
// create a new job. Right after create the API can return the job without a status, so status is only checked
// after the job has run. Each run leaves an export job behind, because export jobs have no delete API.
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
			{
				RefreshState: true,
				Check:        resource.TestCheckResourceAttr(resourceReference, "status", "No Log Files In Requested Time Range"),
			},
		},
	})
}

// TestAccAppServiceAuditLogExportJob creates an export job and reads it back through the single export job data
// source. Right after create the API can return the job without a status, so status is checked in a second step,
// after the job has run. The resource has no import, because the API does not return start and end.
func TestAccAppServiceAuditLogExportJob(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_as_audit_export_")
	resourceReference := "couchbase-capella_app_service_audit_log_export_job." + name
	dataSourceReference := "data.couchbase-capella_app_service_audit_log_export_job." + name
	start, end := appServiceAuditLogExportWindow()
	config := testAccAppServiceAuditLogExportJobConfig(name, start, end) + fmt.Sprintf(`
data "couchbase-capella_app_service_audit_log_export_job" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  id              = couchbase-capella_app_service_audit_log_export_job.%[1]s.id
}
`, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)

	var jobId string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						jobId = s.RootModule().Resources[resourceReference].Primary.ID
						return nil
					},
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "start", start),
					resource.TestCheckResourceAttr(resourceReference, "end", end),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttrPair(dataSourceReference, "id", resourceReference, "id"),
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
				),
			},
			// Once the job has run, the resource and the data source both report its status.
			{
				PreConfig: func() { waitForAppServiceAuditLogExportJobToRun(t, jobId) },
				Config:    config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "status"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "status", resourceReference, "status"),
				),
			},
		},
	})
}

func TestAccAppServiceAuditLogExportJobEndBeforeStart(t *testing.T) {
	start, end := appServiceAuditLogExportWindow()
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceAuditLogExportJobConfig(randomStringWithPrefix("tf_acc_as_audit_export_"), end, start),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`end must not be earlier than start`),
			},
		},
	})
}

func TestAccAppServiceAuditLogExportJobInvalidTimestamp(t *testing.T) {
	_, end := appServiceAuditLogExportWindow()
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceAuditLogExportJobConfig(randomStringWithPrefix("tf_acc_as_audit_export_"), "2026-11-23 10:00:00", end),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`start must be a valid RFC3339 timestamp`),
			},
		},
	})
}

// appServiceAuditLogExportWindow returns a 4-hour window ending an hour ago, inside the 30-day retention period.
func appServiceAuditLogExportWindow() (string, string) {
	end := time.Now().UTC().Add(-1 * time.Hour).Truncate(time.Second)
	return end.Add(-4 * time.Hour).Format(time.RFC3339), end.Format(time.RFC3339)
}

func testAccAppServiceAuditLogExportJobConfig(name, start, end string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service_audit_log_export_job" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  start           = "%[7]s"
  end             = "%[8]s"
}
`, globalProviderBlock, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, start, end)
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
