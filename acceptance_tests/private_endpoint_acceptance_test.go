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

// TestAccPrivateEndpointAccept tests accepting a private endpoint
func TestAccPrivateEndpointAccept(t *testing.T) {
	resourceName := "accept_private_endpoint"
	resourceReference := "couchbase-capella_private_endpoints." + resourceName
	// Note: In a real test, you would need to create an actual endpoint ID from your cloud provider
	// This is a placeholder and would need to be replaced with a valid endpoint ID
	endpointId := "vpce-062b45ce286b28577"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateEndpointAcceptConfig(resourceName, endpointId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "endpoint_id", endpointId),
					resource.TestCheckResourceAttrSet(resourceReference, "status"),
				),
			},
		},
	})
}

// // TestAccPrivateEndpointImport tests importing a private endpoint
// func TestAccPrivateEndpointImport(t *testing.T) {
// 	resourceName := "import_private_endpoint"
// 	resourceReference := "couchbase-capella_private_endpoints." + resourceName
// 	endpointId := "endpoint-12345678-1234-1234-1234-123456789012"

// 	resource.ParallelTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccPrivateEndpointAcceptConfig(resourceName, endpointId),
// 			},
// 			{
// 				ResourceName:      resourceReference,
// 				ImportStateId:     fmt.Sprintf("organization_id=%s,project_id=%s,cluster_id=%s,endpoint_id=%s", globalOrgId, globalProjectId, globalClusterId, endpointId),
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// // TestAccPrivateEndpointServiceImport tests importing a private endpoint service
// func TestAccPrivateEndpointServiceImport(t *testing.T) {
// 	resourceName := "import_private_endpoint_service"
// 	resourceReference := "couchbase-capella_private_endpoint_service." + resourceName

// 	resource.ParallelTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccPrivateEndpointServiceEnableConfig(resourceName, true),
// 			},
// 			{
// 				ResourceName:      resourceReference,
// 				ImportStateId:     fmt.Sprintf("organization_id=%s,project_id=%s,cluster_id=%s", globalOrgId, globalProjectId, globalClusterId),
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

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

// testAccPrivateEndpointAcceptConfig returns terraform config for accepting a private endpoint
func testAccPrivateEndpointAcceptConfig(resourceName string, endpointId string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_private_endpoints" "%[2]s" {
			organization_id = "%[3]s"
			project_id      = "%[4]s"
			cluster_id      = "%[5]s"
			endpoint_id     = "%[6]s"
		}
		`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, endpointId)
}
