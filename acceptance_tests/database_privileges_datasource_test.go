package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatasourceDatabasePrivileges(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_ds_")
	dsReference := "data.couchbase-capella_database_privileges." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabasePrivilegesDataSourceConfig(dsName, globalClusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					testAccCheckDatabasePrivilegesNonEmpty(dsReference),
				),
			},
		},
	})
}

func TestAccDatasourceDatabasePrivilegesInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccDatabasePrivilegesDataSourceConfig(dsName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Database Privileges|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceDatabasePrivilegesMissingCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_missing_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_privileges" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)cluster_id|argument.*required`),
			},
		},
	})
}

func TestAccDatasourceDatabasePrivilegesMissingProject(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_missing_project_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_privileges" "%[2]s" {
  organization_id = "%[3]s"
  cluster_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)project_id|argument.*required`),
			},
		},
	})
}

// testAccCheckDatabasePrivilegesNonEmpty verifies the datasource returned at
// least one privilege item with name and group populated.
func testAccCheckDatabasePrivilegesNonEmpty(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		count, _ := strconv.Atoi(ds.Primary.Attributes["data.#"])
		if count == 0 {
			return fmt.Errorf("datasource %s returned no database privileges", dsReference)
		}
		for i := 0; i < count; i++ {
			name := ds.Primary.Attributes[fmt.Sprintf("data.%d.name", i)]
			group := ds.Primary.Attributes[fmt.Sprintf("data.%d.group", i)]
			if name != "" && group != "" {
				return nil
			}
		}
		return fmt.Errorf("datasource %s has %d elements but none had name and group both set", dsReference, count)
	}
}

func testAccDatabasePrivilegesDataSourceConfig(dsName, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_privileges" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, clusterID)
}
