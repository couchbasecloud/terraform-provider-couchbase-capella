package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccPrivateEndpointServiceEnableDisable tests enabling, reading, and disabling private endpoint service for a cluster.
func TestAccPrivateEndpointServiceEnableDisable(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_private_endpoint_service_")
	dataSourceName := randomStringWithPrefix("tf_acc_private_endpoints_ds_")
	resourceReference := "couchbase-capella_private_endpoint_service." + resourceName
	dataSourceReference := "data.couchbase-capella_private_endpoints." + dataSourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateEndpointsDataSourceNoEndpointConfig(resourceName, dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "private_endpoint_dns"),
					resource.TestCheckResourceAttr(dataSourceReference, "data.#", "0"),
				),
			},
			{
				Config: testAccPrivateEndpointServiceEnableConfig(resourceName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "enabled", "false"),
				),
			},
		},
	})
}

// testAccPrivateEndpointServiceEnableConfig returns terraform config for enabling/disabling private endpoint service
func testAccPrivateEndpointServiceEnableConfig(resourceName string, enabled bool) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_private_endpoint_service" "%[2]s" {
			organization_id = "%[3]s"
			project_id      = "%[4]s"
			cluster_id      = "%[5]s"
			enabled         = %[6]t
		}
		`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, enabled)
}

func testAccPrivateEndpointsDataSourceNoEndpointConfig(serviceResourceName, dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_private_endpoint_service" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  enabled         = true
}

data "couchbase-capella_private_endpoints" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"

  depends_on = [couchbase-capella_private_endpoint_service.%[2]s]
}
`, globalProviderBlock, serviceResourceName, dataSourceName, globalOrgId, globalProjectId, globalClusterId)
}

// TestAccPrivateEndpointsInvalidEndpointID verifies that the private endpoints resource rejects
// an empty endpoint_id at plan time. The LengthAtLeast(1) validator fires before any API call,
// so dummy org/project/cluster IDs are sufficient.
func TestAccPrivateEndpointsInvalidEndpointID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_private_endpoints_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccPrivateEndpointsInvalidEndpointIDConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)endpoint_id.*string length must be at least 1`),
			},
		},
	})
}

func testAccPrivateEndpointsInvalidEndpointIDConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_private_endpoints" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  endpoint_id     = ""
}
`, globalProviderBlock, resourceName)
}
