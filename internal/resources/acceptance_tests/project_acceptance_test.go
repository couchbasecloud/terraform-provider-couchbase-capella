package acceptance_tests

import (
	"fmt"
	"testing"

	"terraform-provider-capella/internal/provider"
	acctest "terraform-provider-capella/internal/testing"

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

func TestAccProjectResource(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig(acctest.Cfg, rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: generateProjectImportIdForResource(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfigUpdate(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project_name_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "description_update"),
				),
			},
			{
				Config: testAccProjectResourceConfigUpdateWithIfMatch(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project_name_update_with_if_match"),
					resource.TestCheckResourceAttr(resourceName, "description", "description_update_with_match"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 3"),
					resource.TestCheckResourceAttr(resourceName, "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfig(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "%[2]s"
	description     = "description"
}
`, cfg, rnd)
}

func testAccProjectResourceConfigUpdate(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
   organization_id = var.organization_id
	name            = "acc_test_project_name_update"
	description     = "description_update"
}
`, cfg, rnd)
}

func testAccProjectResourceConfigUpdateWithIfMatch(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name_update_with_if_match"
	description     = "description_update_with_match"
	if_match        =  2
}
`, cfg, rnd)
}

// generateProjectImportIdForResource generates a project import ID based on the provided resource name
// and the attributes in the Terraform state.
//
// This function takes a resource name as input and returns a function of type `resource.ImportStateIdFunc`.
// The generated import ID is in the format "id=<value>,organization_id=<value>".
func generateProjectImportIdForResource(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
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
}
