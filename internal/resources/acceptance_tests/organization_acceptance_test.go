package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "name"),
					resource.TestCheckResourceAttr("data.couchbase-capella_organization.get_organization", "organization_id", OrgId),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.created_at"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.modified_by"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.modified_at"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.version"),
				),
			},
		},
	})
}

func testAccOrganizationResourceConfig() string {
	return fmt.Sprintf(`
%[1]s

output "organizations_get" {
  value = data.couchbase-capella_organization.get_organization
}

data "couchbase-capella_organization" "get_organization" {
  organization_id = "%[2]s"
}

`, ProviderBlock, OrgId)
}
