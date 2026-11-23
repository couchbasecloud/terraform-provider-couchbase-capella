package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceAppServiceAuditLogExportJobsRegistered reads the export jobs list of the shared App Service.
// It fails with "Invalid data source" until the data source is registered in the provider (AV-153097).
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
