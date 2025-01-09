package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReadOrganization(t *testing.T) {
	resourceName := RandomStringWithPrefix("tf_acc_org_")
	resourceReference := "data.couchbase-capella_organization." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "name"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", OrgId),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
		},
	})
}

func testAccOrganizationResourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_organization" "%[3]s" {
  organization_id = "%[2]s"
}

`, ProviderBlock, OrgId, resourceName)
}
