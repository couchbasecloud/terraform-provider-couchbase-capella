package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{},
			// Import state
			{},
			// Update and Read
			{},
			// NOTE: No delete case is provided - this occurs automatically

		},
	})
}

func testAccDatabaseCredentialConfig(cfg string) string {
	return ""
}

func testAccDatabaseCredentialConfigUpdate(cfg string) string {
	return ""
}

func generateDatabaseCredentialImportId(state *terraform.State) (string, error) {
	return "", nil
}
