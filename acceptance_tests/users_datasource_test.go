package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The happy-path users datasource test is not included here because the datasource
// has no pagination or filtering, so it scans the entire org and exceeds the test budget
// on tenants with many users.

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
