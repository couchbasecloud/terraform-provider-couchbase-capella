package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppEndpointResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_ep_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppEndpointResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "bucket", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "scope", globalScopeName),
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppEndpointImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: testAccAppEndpointResourceConfigUpdate(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppEndpointResourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[6]s" {
  organization_id    = "%[2]s"
  project_id         = "%[3]s"
  cluster_id         = "%[4]s"
  app_service_id     = "%[5]s"
  bucket             = "%[7]s"
  name               = "%[6]s"
  delta_sync_enabled = false
  scope              = "%[8]s"
  collections = {
    "%[9]s" = {}
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, resourceName, globalBucketName, globalScopeName, globalCollectionName)
}

func testAccAppEndpointResourceConfigUpdate(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[6]s" {
  organization_id    = "%[2]s"
  project_id         = "%[3]s"
  cluster_id         = "%[4]s"
  app_service_id     = "%[5]s"
  bucket             = "%[7]s"
  name               = "%[6]s"
  delta_sync_enabled = true
  scope              = "%[8]s"
  collections = {
    "%[9]s" = {}
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, resourceName, globalBucketName, globalScopeName, globalCollectionName)
}

func generateAppEndpointImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("name=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["name"], rawState["app_service_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}
