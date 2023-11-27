package acceptance_tests

import (
	"fmt"
	"testing"

	acctest "terraform-provider-capella/internal/testing"
	cfg "terraform-provider-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestUserResource(t *testing.T) {
	resourceName := "acc_user" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_user." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserResourceConfig(cfg.Cfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr(resourceReference, "email", "acc_test_email"),
					resource.TestCheckResourceAttr(
						resourceReference, "organization_roles", "acc_test_organization_roles",
					),
					resource.TestCheckResourceAttr(resourceReference, "resources", "acc_test_resources"),
				),
			},
			// Import state
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateUserImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccUserResourceConfig(cfg.Cfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr(resourceReference, "email", "acc_test_email"),
					resource.TestCheckResourceAttr(
						resourceReference, "organization_roles", "acc_test_organization_roles",
					),
					resource.TestCheckResourceAttr(resourceReference, "resources", "acc_test_resources"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

func testAccUserResourceConfig(cfg, resourceReference, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
	%[1]s
	  
	resource "capella_project" "%[3]s" {
		organization_id = var.organization_id
		name            = "acc_test_project_name"
		description     = "description"
	}
	
	resource "capella_user" "%[2]s" {
		organization_id = var.organization_id
	  
		name  = "Terraform Acceptance Test User"
		email = "terraformacceptancetest@couchbase.com"
	  
		organization_roles = [
			"organizationOwner"
		]
	  
		resources = [
		  {
			type = "project"
			id   = %[4]s.id
			roles = [
			  "projectViewer",
			  "projectDataReaderWriter"
			]
		  }
		]
	  }
	`, cfg, resourceReference, projectResourceName, projectResourceReference)
}

func testAccUserResourceConfigUpdate(cfg, resourceReference, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
	%[1]s
	  
	resource "capella_project" "%[3]s" {
		organization_id = var.organization_id
		name            = "acc_test_project_name"
		description     = "description"
	}
	
	resource "capella_user" "%[2]s" {
		organization_id = var.organization_id
	  
		name  = "Terraform Acceptance Test User"
		email = "terraformacceptancetest@couchbase.com"
	  
		organization_roles = [
			"organizationOwner"
		]
	  
		resources = [
		  {
			type = "project"
			id   = %[4]s.id
			roles = [
			  "projectViewer",
			]
		  }
		]
	  }
	`, cfg, resourceReference, projectResourceName, projectResourceReference)
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
