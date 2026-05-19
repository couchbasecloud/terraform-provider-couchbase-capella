package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceAuditLogExport drives both the audit_log_export resource
// and its datasource in the same plan: the resource creates an export job,
// the datasource lists exports and must contain at least one entry. The
// datasource's `data` attribute is a SetNestedAttribute, so elements aren't
// addressable by stable numeric indices; assertions use the `data.*`
// set-aware checks and confirm membership of our just-created export.
func TestAccDatasourceAuditLogExport(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_export_for_ds_")
	dsName := randomStringWithPrefix("tf_acc_audit_log_export_ds_")
	dsReference := "data.couchbase-capella_audit_log_export." + dsName

	start, end := auditLogExportWindow()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditLogExportResourceAndDatasourceConfig(resourceName, dsName, start, end),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					// At least one export exists because the resource block
					// in this step just created one. Use set-aware membership
					// assertion to find an element whose ids match ours.
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
						"start":           start,
						"end":             end,
					}),
				),
			},
		},
	})
}

func TestAccDatasourceAuditLogExportInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_audit_log_export_ds_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_audit_log_export" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Audit Log Exports|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func testAccAuditLogExportResourceAndDatasourceConfig(resourceName, dsName, start, end string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_export" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  start           = "%[7]s"
  end             = "%[8]s"
}

data "couchbase-capella_audit_log_export" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"

  depends_on = [couchbase-capella_audit_log_export.%[2]s]
}
`, globalProviderBlock, resourceName, dsName, globalOrgId, globalProjectId, globalClusterId, start, end)
}
