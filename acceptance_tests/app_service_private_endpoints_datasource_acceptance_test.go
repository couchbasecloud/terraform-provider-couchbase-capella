package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceAppServicePrivateEndpoints lists the private endpoints of the global App Service,
// first without a filter and then filtered by status.
func TestAccDatasourceAppServicePrivateEndpoints(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_pes_")
	dataSourceReference := "data.couchbase-capella_app_service_private_endpoints." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicePrivateEndpointsDatasourceConfig(dataSourceName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "data.#"),
				),
			},
			{
				// The test App Service has no endpoint in the unrecognized status, so the
				// filtered list is empty.
				Config: testAccAppServicePrivateEndpointsDatasourceConfig(dataSourceName, `status = ["unrecognized"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "status.#", "1"),
					resource.TestCheckResourceAttr(dataSourceReference, "data.#", "0"),
				),
			},
		},
	})
}

// TestAccDatasourceAppServicePrivateEndpointsInvalidStatus verifies that an unknown status fails at
// plan time, before any API call, so dummy IDs are sufficient.
func TestAccDatasourceAppServicePrivateEndpointsInvalidStatus(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_private_endpoints" "test" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  app_service_id  = "33333333-3333-3333-3333-333333333333"
  status          = ["active"]
}
`, globalProviderBlock),
				ExpectError: regexp.MustCompile(`(?s)status.*value must be one of`),
			},
		},
	})
}

func testAccAppServicePrivateEndpointsDatasourceConfig(dataSourceName, statusFilter string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_private_endpoints" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  %[7]s
}
`, globalProviderBlock, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, statusFilter)
}
