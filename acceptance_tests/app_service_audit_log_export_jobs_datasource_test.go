package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccDatasourceAppServiceAuditLogExportJobsRegistered reads the export jobs list of the shared App Service. It
// guards that the data source stays registered in the provider: it was missing, failing with "Invalid data source"
// (AV-153097).
func TestAccDatasourceAppServiceAuditLogExportJobsRegistered(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
%s
data "couchbase-capella_app_service_audit_log_export_jobs" "x" {
  organization_id = "%s"
  project_id      = "%s"
  cluster_id      = "%s"
  app_service_id  = "%s"
}`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId),
			Check: resource.TestCheckResourceAttrSet("data.couchbase-capella_app_service_audit_log_export_jobs.x", "data.#"),
		}},
	})
}

// TestAccDatasourceAppServiceAuditLogExportJobs creates an export job, then lists the App Service's export jobs
// filtered by its ID and by an ID that does not exist. Right after create the API can return the job without a
// status, so status is checked in a second step, after the job has run.
func TestAccDatasourceAppServiceAuditLogExportJobs(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_as_audit_export_")
	resourceReference := "couchbase-capella_app_service_audit_log_export_job." + name
	matching := "data.couchbase-capella_app_service_audit_log_export_jobs.matching_" + name
	unknown := "data.couchbase-capella_app_service_audit_log_export_jobs.unknown_" + name
	start, end := appServiceAuditLogExportWindow()
	config := testAccAppServiceAuditLogExportJobConfig(name, start, end) + fmt.Sprintf(`
data "couchbase-capella_app_service_audit_log_export_jobs" "matching_%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  export_id       = couchbase-capella_app_service_audit_log_export_job.%[1]s.id
}

data "couchbase-capella_app_service_audit_log_export_jobs" "unknown_%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  export_id       = "00000000-0000-0000-0000-000000000000"
  depends_on      = [couchbase-capella_app_service_audit_log_export_job.%[1]s]
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
					resource.TestCheckResourceAttr(matching, "data.#", "1"),
					resource.TestCheckResourceAttrPair(matching, "data.0.id", resourceReference, "id"),
					resource.TestCheckResourceAttr(matching, "data.0.cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(matching, "data.0.app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(unknown, "data.#", "0"),
				),
			},
			// Once the job has run, the list reports the same status as the resource.
			{
				PreConfig: func() { waitForAppServiceAuditLogExportJobToRun(t, jobId) },
				Config:    config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(matching, "data.0.status"),
					resource.TestCheckResourceAttrPair(matching, "data.0.status", resourceReference, "status"),
				),
			},
		},
	})
}
