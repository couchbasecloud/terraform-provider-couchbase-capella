package resources_test

import (
	"fmt"
	"os"
	"testing"

	"terraform-provider-capella/internal/provider"
	cfg "terraform-provider-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"capella": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for
	// example assertions about the appropriate environment variables being set
	// are common to see in a pre-check function.
}

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_project.acc_test", "name", "acc_test_project_name"),
					resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description"),
					resource.TestCheckResourceAttr("capella_project.acc_test", "etag", "Version: 1"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      "capella_project.acc_test",
				ImportStateIdFunc: generateProjectImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfigUpdate(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_project.acc_test", "name", "acc_test_project_name_update"),
					resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description_update"),
				),
			},
			{
				Config: testAccProjectResourceConfigUpdateWithIfMatch(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_project.acc_test", "name", "acc_test_project_name_update_with_if_match"),
					resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description_update_with_match"),
					resource.TestCheckResourceAttr("capella_project.acc_test", "etag", "Version: 3"),
					resource.TestCheckResourceAttr("capella_project.acc_test", "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	SetResult(t.Failed())

}

func testAccProjectResourceConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "acc_test" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}
`, cfg)
}

func testAccProjectResourceConfigUpdate(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "acc_test" {
   organization_id = var.organization_id
	name            = "acc_test_project_name_update"
	description     = "description_update"
}
`, cfg)
}

func testAccProjectResourceConfigUpdateWithIfMatch(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "acc_test" {
    organization_id = var.organization_id
	name            = "acc_test_project_name_update_with_if_match"
	description     = "description_update_with_match"
	if_match        =  2
}
`, cfg)
}

func generateProjectImportId(state *terraform.State) (string, error) {
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
	return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], rawState["organization_id"]), nil
}

func SetResult(isFailed bool) {
	if isFailed {
		setEnvVariable("TERRAFORM_ACCEPTANCE_TEST_RESULT", "failed")
	} else {
		setEnvVariable("TERRAFORM_ACCEPTANCE_TEST_RESULT", "passed")
	}
}

func setEnvVariable(varname string, val string) {
	err := os.Setenv(varname, val)
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}
}
