package acceptance_tests

import (
	"fmt"
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
