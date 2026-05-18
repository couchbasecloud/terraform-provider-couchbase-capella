package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceAuditLogEventIDs(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_audit_log_event_ids_ds_")
	dsReference := "data.couchbase-capella_audit_log_event_ids." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditLogEventIDsDataSourceConfig(dsName, globalClusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					// The catalog of audit event ids is fixed and non-empty
					// for any supported cluster, so the first element must
					// expose the documented attributes.
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.name"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.module"),
				),
			},
		},
	})
}

func TestAccDatasourceAuditLogEventIDsInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_audit_log_event_ids_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogEventIDsDataSourceConfig(dsName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading audit log event ids|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceAuditLogEventIDsMissingCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_audit_log_event_ids_missing_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_audit_log_event_ids" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)cluster_id|argument.*required`),
			},
		},
	})
}

func testAccAuditLogEventIDsDataSourceConfig(dsName, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_audit_log_event_ids" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, clusterID)
}
