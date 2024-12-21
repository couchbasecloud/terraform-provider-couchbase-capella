package acceptance_tests

import (
	"fmt"
	"testing"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseCredentialTestCases(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//database_credential with required fields
			{
				Config: testAccAddDatabaseCredWithReqFieldsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_req", "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_req", "access.0.privileges.0", "data_writer"),
				),
			},
			//database_credential with optional fields
			{
				Config: testAccAddDatabaseCredWithOptionalFieldsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_opt", "name", "acc_test_database_credential_name2"),
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_opt", "password", "Secret12$#"),
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_opt", "access.0.privileges.0", "data_writer"),
				),
			},
		},
	})
}

// Delete the database_credential when the cluster is destroyed through api
func testAccAddDatabaseCredWithReqFieldsConfig() string {
	return fmt.Sprintf(
		`
		%[1]s
	
		output "add_database_credential_req"{
			value = couchbase-capella_database_credential.add_database_credential_req
			sensitive = true
		}
		
		resource "couchbase-capella_database_credential" "add_database_credential_req" {
			name            = "acc_test_database_credential_name"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, ProviderBlock, OrgId, ProjectId, ClusterId)
}

func testAccAddDatabaseCredWithOptionalFieldsConfig() string {
	return fmt.Sprintf(
		`
		%[1]s
	
		output "add_database_credential_opt"{
			value = couchbase-capella_database_credential.add_database_credential_opt
			sensitive = true
		}
		
		resource "couchbase-capella_database_credential" "add_database_credential_opt" {
			name            = "acc_test_database_credential_name2"
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
		`, ProviderBlock, OrgId, ProjectId, ClusterId)
}
