package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatabaseCredentialPlanOnly guards the database credential schema against a
// declared attribute type drifting away from what the provider code reads back.
//
// AV-139978: `access` was declared as a ListNestedAttribute while ValidateConfig read it
// into a types.Set. That combination builds and vets cleanly but fails every plan with a
// framework "Value Conversion Error", including plain basic credentials. These steps stop
// at plan, so they exercise ValidateConfig and the schema validators without provisioning
// anything or calling the V4 API.
func TestAccDatabaseCredentialPlanOnly(t *testing.T) {
	basicName := randomStringWithPrefix("tf_acc_db_cred_plan_basic_")
	advancedName := randomStringWithPrefix("tf_acc_db_cred_plan_adv_")
	emptyAccessName := randomStringWithPrefix("tf_acc_db_cred_plan_empty_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// A basic credential with no credential_type is the case AV-139978 broke.
			{
				Config:             testAccDatabaseCredentialPlanOnlyBasicConfig(basicName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// An advanced credential reads user_roles down the same ValidateConfig path.
			{
				Config:             testAccDatabaseCredentialPlanOnlyAdvancedConfig(advancedName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// `access = []` satisfies ExactlyOneOf, so without SizeAtLeast it reaches the
			// API as a body carrying neither access nor userRoles.
			{
				Config:      testAccDatabaseCredentialPlanOnlyEmptyAccessConfig(emptyAccessName),
				ExpectError: regexp.MustCompile(`(?s)access.*must contain at least 1`),
			},
		},
	})
}

func testAccDatabaseCredentialPlanOnlyBasicConfig(resourceName string) string {
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
					resources = {
						buckets = [
							{
								name = "test_bucket"
								scopes = [
									{
										name        = "test_scope"
										collections = ["test_collection"]
									},
								]
							},
						]
					}
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccDatabaseCredentialPlanOnlyAdvancedConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			password        = "Secret12$#"
			credential_type = "advanced"
			user_roles      = ["tf_acc_plan_only_role"]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccDatabaseCredentialPlanOnlyEmptyAccessConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			access          = []
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}
