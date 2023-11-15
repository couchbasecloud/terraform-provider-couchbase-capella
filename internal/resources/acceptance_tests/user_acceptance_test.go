package acceptance_tests

import (
	"fmt"
	"testing"

	cfg "terraform-provider-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserResourceConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_user.acc_test", "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr("capella_user.acc_test", "email", "acc_test_email"),
					resource.TestCheckResourceAttr(
						"capella_user.acc_test", "organization_roles", "acc_test_organization_roles",
					),
					resource.TestCheckResourceAttr("capella_user.acc_test", "resources", "acc_test_resources"),
				),
			},
			// Import state
			{
				ResourceName:      "capella_user.acc_test",
				ImportStateIdFunc: generateUserImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccUserResourceConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_user.acc_test", "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr("capella_user.acc_test", "email", "acc_test_email"),
					resource.TestCheckResourceAttr(
						"capella_user.acc_test", "organization_roles", "acc_test_organization_roles",
					),
					resource.TestCheckResourceAttr("capella_user.acc_test", "resources", "acc_test_resources"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

func testAccUserResourceConfig(cfg string) string {
	return fmt.Sprintf(`
	%[1]s
	
	resource "capella_user" "new_user" {
		organization_id = var.organization_id
	  
		name  = var.user_name
		email = var.email
	  
		organization_roles = var.organization_roles
	  
		resources = [
		  {
			type = "project"
			id   = var.project_id
			roles = [
			  "projectViewer",
			  "projectDataReaderWriter"
			]
		  }
		]
	  }
	`, cfg)
}

func testAccUserResourceConfigUpdate(cfg string) string {
	return fmt.Sprintf(`
	%[1]s
	
	resource "capella_user" "new_user" {
		organization_id = var.organization_id
	  
		name  = var.user_name
		email = var.email
	  
		organization_roles = var.organization_roles
	  
		resources = [
		  {
			type = "project"
			id   = var.project_id
			roles = [
			  "projectViewer",
			]
		  }
		]
	  }
	`, cfg)
}

func generateUserImportId(state *terraform.State) (string, error) {
	resourceName := "capella_user.acc_test"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	fmt.Printf("raw state %s", rawState)

	return fmt.Sprintf(
			"id=%s,organization_id=%s",
			rawState["id"],
			rawState["organization_id"]),
		nil
}
