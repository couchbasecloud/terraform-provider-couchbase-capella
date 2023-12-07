package acceptance_tests

import (
	"fmt"
	"testing"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	cfg "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAppServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppServiceResourceConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.ram", "4"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "nodes", "2"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      "couchbase-capella_app_service.new_app_service",
				ImportStateIdFunc: generateAppServiceImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAppServiceResourceConfigUpdate(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.ram", "4"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "nodes", "3"),
				),
			},
			{
				Config: testAccAppServiceResourceConfigUpdateWithIfMatch(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.cpu", "4"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.ram", "8"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "nodes", "2"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppServiceResourceConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "new_app_service" {
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

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = "test-terraform-app-service"
  description     = "description"
  compute = {
    cpu = 2
    ram = 4
  }
  nodes = 3
}
`, cfg)
}

func testAccAppServiceResourceConfigUpdateWithIfMatch(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = "test-terraform-app-service"
  description     = "description"
  if_match        =  2
  compute = {
    cpu = 4
    ram = 8
  }
  nodes = 2
}
`, cfg)
}

func generateAppServiceImportId(state *terraform.State) (string, error) {
	resourceName := "couchbase-capella_app_service.new_app_service"
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
