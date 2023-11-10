package acceptance_tests

import (
	"fmt"
	cfg "terraform-provider-capella/internal/testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestDatabaseCredentialResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccDatabaseCredentialConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_database_credential.acc_test", "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr("capella_database_credential.acc_test", "password", "password"),
					resource.TestCheckResourceAttr("capella_database_credential.acc_test", "access", "access"),
				),
			},
			// Import state
			{
				ResourceName:      "capella_database_credential.acc_test",
				ImportStateIdFunc: generateDatabaseCredentialImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccDatabaseCredentialConfigUpdate(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_project.acc_test", "name", "acc_test_project_name_update_with_if_match"),
					resource.TestCheckResourceAttr("capella_database_credential.acc_test", "password", "updated_password"),
					resource.TestCheckResourceAttr("capella_database_credential.acc_test", "access", "access"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically

		},
	})
}

func testAccDatabaseCredentialConfig(cfg string) string {
	return fmt.Sprintf(`
	%[1]s
	
	resource "capella_database_credential" "new_database_credential" {
		name            = var.database_credential_name
		organization_id = var.organization_id
		project_id      = var.project_id
		cluster_id      = var.cluster_id
		password        = "password"
		access          = "access"
	  }
	`, cfg)
}

func testAccDatabaseCredentialConfigUpdate(cfg string) string {
	return fmt.Sprintf(`
	%[1]s
	
	resource "capella_database_credential" "new_database_credential" {
		name            = var.database_credential_name
		organization_id = var.organization_id
		project_id      = var.project_id
		cluster_id      = var.cluster_id
		password        = "updated_password"
		access          = "access"
	  }
	`, cfg)
}

func generateDatabaseCredentialImportId(state *terraform.State) (string, error) {
	resourceName := "capella_project.acc_test"
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
			"id=%s,organization_id=%s,project_id=%s,cluster_id=%s",
			rawState["id"],
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"]),
		nil
}
