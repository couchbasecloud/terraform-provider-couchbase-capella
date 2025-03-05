package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseCredentialWithReqFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddDatabaseCredWithReqFieldsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_writer"),
				),
			},
		},
	})
}

func TestAccDatabaseCredentialWithOptionalFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddDatabaseCredWithOptionalFieldsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "password", "Secret12$#"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_writer"),
				),
			},
		},
	})
}

func testAccAddDatabaseCredWithReqFieldsConfig(resourceName string) string {
	return fmt.Sprintf(
		`
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
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccAddDatabaseCredWithOptionalFieldsConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s
		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			password        = "Secret12$#"
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}
