package acceptance_tests

import (
	"fmt"
	"testing"

	cfg "terraform-provider-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAppServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppServiceResourceConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.ram", "4"),

					//resource.TestCheckResourceAttr("capella_app_service.new_app_service", "description", "description"),

					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "etag", "Version: 8"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      "capella_app_service.new_app_service",
				ImportStateIdFunc: generateAppServiceImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAppServiceResourceConfigUpdate(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "name", "test-terraform-app-service-update"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "description", "description-update"),
				),
			},
			{
				Config: testAccAppServiceResourceConfigUpdateWithIfMatch(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "name", "acc_test_project_name_update_with_if_match"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "description", "description_update_with_match"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "etag", "Version: 3"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppServiceResourceConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = "test-terraform-app-service"
  description     = "description"
  compute = {
    cpu = 2
    ram = 4
}
}
`, cfg)
}

func testAccAppServiceResourceConfigUpdate(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = "test-terraform-app-service-update"
  description     = "description-update"
}
`, cfg)
}

func testAccAppServiceResourceConfigUpdateWithIfMatch(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = "test-terraform-app-service-update-with-if-match"
  description     = "description-update-with-if-match"
  if_match        =  2
}
`, cfg)
}

func generateAppServiceImportId(state *terraform.State) (string, error) {
	resourceName := "capella_app_service.new_app_service"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	fmt.Printf("raw state %s", rawState)
	return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
}
