package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// appServicePrivateEndpointsDatasourceSteps lists the private endpoints of the global App Service,
// first without a filter and then filtered by status.
//
// The list endpoint is only reachable once the private endpoint service is enabled on the App
// Service - it returns 400 "Private endpoints aren't enabled for this cluster" otherwise - so every
// config enables the service and the data source depends on it. The depends_on also defers the read
// to apply time, so the pre-apply plan of the first step does not call the API before the service
// exists.
//
// These steps run as a subtest of TestAccAppServicePrivateEndpoint, which serialises them against
// the private endpoint service resource test on the same shared App Service.
func appServicePrivateEndpointsDatasourceSteps() []resource.TestStep {
	serviceName := randomStringWithPrefix("tf_acc_app_service_pes_service_")
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_pes_")
	dataSourceReference := "data.couchbase-capella_app_service_private_endpoints." + dataSourceName

	return []resource.TestStep{
		{
			Config: testAccAppServicePrivateEndpointsDatasourceConfig(serviceName, dataSourceName, ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
				// The test App Service has no CSP private endpoint attached to it, so the
				// unfiltered list is empty.
				resource.TestCheckResourceAttr(dataSourceReference, "data.#", "0"),
			),
		},
		{
			// The test App Service has no endpoint in the unrecognized status, so the
			// filtered list is empty.
			Config: testAccAppServicePrivateEndpointsDatasourceConfig(serviceName, dataSourceName, `status = ["unrecognized"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(dataSourceReference, "status.#", "1"),
				resource.TestCheckResourceAttr(dataSourceReference, "data.#", "0"),
			),
		},
	}
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

func testAccAppServicePrivateEndpointsDatasourceConfig(serviceName, dataSourceName, statusFilter string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_private_endpoint_service" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  app_service_id  = "%[7]s"
}

data "couchbase-capella_app_service_private_endpoints" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  app_service_id  = "%[7]s"
  %[8]s

  depends_on = [couchbase-capella_app_service_private_endpoint_service.%[2]s]
}
`, globalProviderBlock, serviceName, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, statusFilter)
}
