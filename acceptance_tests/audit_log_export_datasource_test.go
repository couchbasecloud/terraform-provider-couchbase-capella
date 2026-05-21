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
	// Skipped for the same reason as TestAccAuditLogExportResource: this
	// test creates a couchbase-capella_audit_log_export resource inline
	// (so the datasource has something to list). The CI cluster does not
	// have audit logging configured, so GET /auditLogExports/{id}
	// returns 404 "No audit log files exist within the requested time
	// frame." right after create. refreshAuditLogExport classifies that
	// 404 as ResourceNotFound, the framework removes the resource from
	// state, and the post-apply refresh plan shows the resource as a
	// new create — failing the step with "the refresh plan was not
	// empty". Re-enable once AV-128951's infra work seeds audit log data
	// on the CI cluster.
	t.Skip("requires audit log data on globalClusterId; see comment above")

	resourceName := randomStringWithPrefix("tf_acc_audit_log_export_for_ds_")
	resourceReference := "couchbase-capella_audit_log_export." + resourceName
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
					// Match the just-created export by its server-assigned id.
					// The datasource's start/end attrs are written via
					// time.Time.String() ("2006-01-02 15:04:05 -0700 UTC")
					// while the resource normalises them to the
					// "2006-01-02T15:04:05-07:00" RFC3339-with-offset layout,
					// so matching by time string here would fail on format
					// alone. Matching on id is both stable and what the test
					// actually cares about (membership of THIS export).
					resource.TestCheckTypeSetElemAttrPair(dsReference, "data.*.id", resourceReference, "id"),
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
