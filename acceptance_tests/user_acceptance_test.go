package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	resourceReference := "couchbase-capella_user." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserResourceConfig(resourceName, "terraform_acceptance_test1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "terraform_acceptance_test1"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraform_acceptance_test1@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationOwner"),
				),
			},
			// Import state
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateUserImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccUserResourceConfigUpdate(resourceName, "terraform_acceptance_test1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "terraform_acceptance_test1"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraform_acceptance_test1@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectViewer"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

func testAccUserResourceConfig(resourceName, username string) string {
	return fmt.Sprintf(`
	%[1]s
	
	resource "couchbase-capella_user" "%[2]s" {
		organization_id = "%[3]s"
	  
		name  = "%[4]s"
		email = "%[5]s"
	  
		organization_roles = [
			"organizationOwner"
		]
	  }
	`, globalProviderBlock, resourceName, globalOrgId, username, username+"@couchbase.com")
}

func testAccUserResourceConfigUpdate(resourceName, username string) string {
	return fmt.Sprintf(`
	%[1]s
	resource "couchbase-capella_user" "%[2]s" {
		organization_id = "%[3]s"
	  
		name  = "%[5]s"
		email = "%[6]s"
	  
		organization_roles = [
			"organizationMember"
		]
	  
		resources = [
		  {
			type = "project"
			id   = "%[4]s"
			roles = [
			  "projectViewer",
			]
		  }
		]
	  }
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, username, username+"@couchbase.com")
}

func generateUserImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], globalOrgId), nil
	}
}
