package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
)

// This test is not parallel so it's done before all the other tests.
// The reason is buckets cannot be deleted while app service is being deleted.
func TestAppServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppServiceResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.ram", "4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "couchbase-capella_app_service.new_app_service",
				ImportStateIdFunc: generateAppServiceImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppServiceResourceConfig() string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "test-terraform-app-service"
  description     = "description"
  compute = {
    cpu = 2
    ram = 4
  }
}
`, ProviderBlock, OrgId, ProjectId, ClusterId)
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
	return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
}
