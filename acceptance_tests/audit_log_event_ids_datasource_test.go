package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					// `data` is a SetNestedAttribute — elements aren't
					// addressable by stable numeric indices, so we scan all
					// elements and require at least one with id/name/module
					// populated rather than asserting on data.0.*.
					testAccCheckAuditLogEventIDsNonEmpty(dsReference),
				),
			},
		},
	})
}

func testAccCheckAuditLogEventIDsNonEmpty(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		count, _ := strconv.Atoi(ds.Primary.Attributes["data.#"])
		if count == 0 {
			return fmt.Errorf("datasource %s returned no audit log event ids", dsReference)
		}
		for i := 0; i < count; i++ {
			id := ds.Primary.Attributes[fmt.Sprintf("data.%d.id", i)]
			name := ds.Primary.Attributes[fmt.Sprintf("data.%d.name", i)]
			module := ds.Primary.Attributes[fmt.Sprintf("data.%d.module", i)]
			if id != "" && name != "" && module != "" {
				return nil
			}
		}
		return fmt.Errorf("datasource %s has %d elements but none had id/name/module all set", dsReference, count)
	}
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
