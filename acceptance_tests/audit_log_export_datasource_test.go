package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceAuditLogExport drives both the audit_log_export resource
// and its datasource in the same plan: the resource creates an export job,
// the datasource lists exports and must contain at least one entry. We
// don't try to match by id because the datasource returns the full list;
// the count assertion `data.#` >= 1 plus the cluster/project/org IDs are
// enough to guard the datasource behaviour.
func TestAccDatasourceAuditLogExport(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_export_for_ds_")
	dsName := randomStringWithPrefix("tf_acc_audit_log_export_ds_")
	dsReference := "data.couchbase-capella_audit_log_export." + dsName

	end := time.Now().UTC().Add(-1 * time.Hour).Truncate(time.Second)
	start := end.Add(-4 * time.Hour)
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditLogExportResourceAndDatasourceConfig(resourceName, dsName, startStr, endStr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					// At least one export exists because the resource block
					// in this step just created one.
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.start"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.end"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.created_at"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.status"),
					resource.TestCheckResourceAttr(dsReference, "data.0.organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "data.0.project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "data.0.cluster_id", globalClusterId),
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
