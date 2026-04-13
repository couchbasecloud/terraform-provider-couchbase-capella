package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestDataAPIResource(t *testing.T) {
	t.Parallel()
	resourceName := randomStringWithPrefix("tf_acc_data_api_")
	resourceReference := "couchbase-capella_data_api." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDataAPIResourceConfig(resourceName, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttrSet(resourceReference, "state"),
					resource.TestCheckResourceAttrSet(resourceReference, "state_for_network_peering"),
					resource.TestCheckResourceAttrSet(resourceReference, "connection_string"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDataAPIImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_data_api",
					"enable_network_peering",
				},
			},
			// Update testing - enable network peering
			{
				Config: testAccDataAPIResourceConfig(resourceName, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "true"),
					resource.TestCheckResourceAttrSet(resourceReference, "state"),
					resource.TestCheckResourceAttrSet(resourceReference, "state_for_network_peering"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDataAPIResourceConfig(resourceName string, enableDataApi, enableNetworkPeering bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id        = "%[3]s"
  project_id             = "%[4]s"
  cluster_id             = "%[5]s"
  enable_data_api        = %[6]t
  enable_network_peering = %[7]t
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, enableDataApi, enableNetworkPeering)
}

func generateDataAPIImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}
