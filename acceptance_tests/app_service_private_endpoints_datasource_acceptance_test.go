package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
				testAccCheckAppServicePrivateEndpointsWellFormed(dataSourceReference),
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

// testAccCheckAppServicePrivateEndpointsWellFormed asserts that every endpoint the data source
// returned carries an id and a status, which guards the conversion of the API response into state.
//
// It deliberately does not assert the length. The shared App Service can be supplied through
// TF_VAR_app_service_id and may already carry endpoints, and the ACC_APP_SERVICE_ENDPOINT_ID
// scenario requires one on that same App Service, so an exact count is environment-dependent. An
// empty list passes trivially; the emptiness assertion belongs on the filtered step, where the
// filter is what makes the result predictable.
//
// service_name is not checked: the API models it as optional, so a null value is legitimate.
func testAccCheckAppServicePrivateEndpointsWellFormed(dataSourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceReference]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataSourceReference)
		}

		count, err := strconv.Atoi(ds.Primary.Attributes["data.#"])
		if err != nil {
			return fmt.Errorf("could not read %s data.#: %w", dataSourceReference, err)
		}

		for i := 0; i < count; i++ {
			for _, attr := range []string{"id", "status"} {
				key := fmt.Sprintf("data.%d.%s", i, attr)
				if ds.Primary.Attributes[key] == "" {
					return fmt.Errorf("%s %s is empty", dataSourceReference, key)
				}
			}
		}

		return nil
	}
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
