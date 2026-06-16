package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseRoleWithReqFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf-acc-db-role-")
	resourceReference := "couchbase-capella_database_role." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleWithReqFieldsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDatabaseRoleWithOptionalFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf-acc-db-role-")
	resourceReference := "couchbase-capella_database_role." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleWithOptionalFieldsConfig(resourceName, "initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "initial description"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
				),
			},
			{
				Config: testAccDatabaseRoleWithOptionalFieldsConfig(resourceName, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "updated description"),
				),
			},
		},
	})
}

func testAccDatabaseRoleWithReqFieldsConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_role" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			access = [
				{
					privileges = ["data_reader"]
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccDatabaseRoleWithOptionalFieldsConfig(resourceName, description string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_role" "%[6]s" {
			name            = "%[6]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			description     = "%[5]s"
			access = [
				{
					privileges = ["data_reader", "data_writer"]
					resources = {
						buckets = [
							{
								name = "*"
							},
						]
					}
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, description, resourceName)
}
