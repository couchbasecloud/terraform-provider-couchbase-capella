package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceDatabaseCredentials(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credentials_")
	dsName := randomStringWithPrefix("tf_acc_database_credentials_ds_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	dsReference := "data.couchbase-capella_database_credentials." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialsDataSourceConfig(resourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            resourceName,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
				),
			},
		},
	})
}

// TestAccDatasourceDatabaseCredentialsAdvanced verifies that the datasource returns the
// user roles assigned to an advanced database credential.
func TestAccDatasourceDatabaseCredentialsAdvanced(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credentials_adv_")
	dsName := randomStringWithPrefix("tf_acc_database_credentials_ds_")
	roleName := randomStringWithPrefix("tf_acc_db_role_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	dsReference := "data.couchbase-capella_database_credentials." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialsAdvancedDataSourceConfig(resourceName, dsName, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "credential_type", "advanced"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":         resourceName,
						"user_roles.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(dsReference, "data.*.user_roles.*", roleName),
				),
			},
		},
	})
}

func TestAccDatasourceDatabaseCredentialsInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_database_credentials_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_credentials" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				// a bogus cluster UUID returns 403/404 from the database credentials endpoint.
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Database Credentials.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

func testAccDatabaseCredentialsDataSourceConfig(resourceName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[5]s" {
  name            = "%[5]s"
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  access = [
    {
      privileges = ["data_writer"]
    },
  ]
}

data "couchbase-capella_database_credentials" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [couchbase-capella_database_credential.%[5]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, dsName)
}

func testAccDatabaseCredentialsAdvancedDataSourceConfig(resourceName, dsName, roleName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" "role1" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[7]s"
  %[8]s
}

resource "couchbase-capella_database_credential" "%[5]s" {
  name            = "%[5]s"
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  credential_type = "advanced"
  user_roles      = [couchbase-capella_database_role.role1.name]
}

data "couchbase-capella_database_credentials" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [couchbase-capella_database_credential.%[5]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, dsName,
		roleName, databaseRoleAccessBlock(`"dataRead"`))
}
