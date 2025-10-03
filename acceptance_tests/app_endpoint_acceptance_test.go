package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppEndpoint(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, "syncFnXattr", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "bucket", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "user_xattr_key", "syncFnXattr"),
				),
			},
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, "syncFnXattrUpdated", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "user_xattr_key", "syncFnXattrUpdated"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppEndpointImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppEndpointResourceConfig(resourceName, endpointName, userXattr string, deltaSync bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  bucket          = "%[7]s"
  name            = "%[8]s"
  user_xattr_key  = "%[9]s"
  delta_sync_enabled = %[10]t

  scopes = {
    "%[11]s" = {
      collections = {
        "%[12]s" = {}
      }
    }
  }
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		globalBucketName,
		endpointName,
		userXattr,
		deltaSync,
		globalScopeName,
		globalCollectionName,
	)
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
		// Import uses the endpoint name
		return rawState["name"], nil
	}
}
