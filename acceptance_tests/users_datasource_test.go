package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceUsers(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_users_")
	dsName := randomStringWithPrefix("tf_acc_users_ds_")
	resourceReference := "couchbase-capella_user." + resourceName
	dsReference := "data.couchbase-capella_users." + dsName

	// Match the pattern in user_acceptance_test.go — the username/email pair is
	// a fixture in the test tenant and the invite flow needs a deterministic
	// value. The resource handle uses a randomised name so parallel tests do
	// not collide on terraform state.
	username := "terraform_acceptance_test_ds"
	email := username + "@couchbase.com"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig(resourceName, dsName, username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", username),
					resource.TestCheckResourceAttr(resourceReference, "email", email),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            username,
						"email":           email,
						"organization_id": globalOrgId,
					}),
				),
			},
		},
	})
}

func TestAccDatasourceUsersInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_users_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_users" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName),
				// Read() wraps API errors with "Error Reading Capella Users"; a
				// bogus org id is rejected by /v4/organizations/.../users with a
				// 403/404. Require both so this test only passes for that exact
				// failure mode, not for any unrelated diagnostic.
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Users.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

func testAccUsersDataSourceConfig(resourceName, dsName, username, email string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[3]s" {
  organization_id = "%[2]s"

  name  = "%[5]s"
  email = "%[6]s"

  organization_roles = [
    "organizationMember"
  ]
}

data "couchbase-capella_users" "%[4]s" {
  organization_id = "%[2]s"

  depends_on = [couchbase-capella_user.%[3]s]
}
`, globalProviderBlock, globalOrgId, resourceName, dsName, username, email)
}
