package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccPrivateEndpointServiceEnable tests enabling private endpoint service for a cluster
func TestAccPrivateEndpointServiceEnable(t *testing.T) {
	resourceName := "new_private_endpoint_service"
	resourceReference := "couchbase-capella_private_endpoint_service." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateEndpointServiceEnableConfig(resourceName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "enabled", "true"),
				),
			},
		},
	})
}

// TestAccPrivateEndpointServiceDisable tests disabling private endpoint service for a cluster
func TestAccPrivateEndpointServiceDisable(t *testing.T) {
	resourceName := "disable_private_endpoint_service"
	resourceReference := "couchbase-capella_private_endpoint_service." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// First enable the service
			{
				Config: testAccPrivateEndpointServiceEnableConfig(resourceName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "enabled", "true"),
				),
			},
			// Then disable it
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

// TestAccPrivateEndpointServiceUpdate tests updating private endpoint service status
func TestAccPrivateEndpointServiceUpdate(t *testing.T) {
	resourceName := "update_private_endpoint_service"
	resourceReference := "couchbase-capella_private_endpoint_service." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create with enabled=true
			{
				Config: testAccPrivateEndpointServiceEnableConfig(resourceName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "enabled", "true"),
				),
			},
			// Update to enabled=false
			{
				Config: testAccPrivateEndpointServiceEnableConfig(resourceName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "enabled", "false"),
				),
			},
			// Update back to enabled=true
			{
				Config: testAccPrivateEndpointServiceEnableConfig(resourceName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "enabled", "true"),
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
