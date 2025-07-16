package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This test is not parallel so it's done before all the other tests.
// The reason is buckets cannot be deleted while app service is being deleted.
func TestAppEndpointResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppEndpointResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "bucket", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppEndpointImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppEndpointResourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[7]s
	
resource "couchbase-capella_app_endpoint" "%[6]s" {
  organization_id = "%[1]s"
  project_id      = "%[2]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[3]s"
  bucket          = "%[5]s"
  name            = "%[6]s"
  
}
`, globalOrgId, globalProjectId, globalAppServiceId, globalClusterId, globalBucketName, resourceName, globalProviderBlock)
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
		return fmt.Sprintf("id=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["id"], rawState["app_service_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}
