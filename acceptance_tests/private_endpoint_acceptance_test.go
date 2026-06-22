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

// TestAccAWSPrivateEndpointCommandInvalidVPCID verifies that the AWS private endpoint
// command data source rejects a vpc_id shorter than the OpenAPI minimum length at plan
// time. The generated LengthBetween(12, 21) validator fires before any API call, so dummy
// org/project/cluster IDs are sufficient.
func TestAccAWSPrivateEndpointCommandInvalidVPCID(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_aws_pe_command_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAWSPrivateEndpointCommandInvalidVPCIDConfig(dataSourceName),
				ExpectError: regexp.MustCompile(`(?s)vpc_id.*string length must be between 12 and 21`),
			},
		},
	})
}

func testAccAWSPrivateEndpointCommandInvalidVPCIDConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_aws_private_endpoint_command" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  vpc_id          = "vpc-short"
  subnet_ids      = ["subnet-1234567890abcdef0"]
}
`, globalProviderBlock, dataSourceName)
}

// TestAccAzurePrivateEndpointCommandInvalidVirtualNetwork verifies that the Azure private
// endpoint command data source rejects a virtual_network shorter than the OpenAPI minimum
// length at plan time. The generated LengthBetween(2, 64) validator fires before any API call.
func TestAccAzurePrivateEndpointCommandInvalidVirtualNetwork(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_azure_pe_command_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzurePrivateEndpointCommandInvalidVirtualNetworkConfig(dataSourceName),
				ExpectError: regexp.MustCompile(`(?s)virtual_network.*string length must be between 2 and 64`),
			},
		},
	})
}

func testAccAzurePrivateEndpointCommandInvalidVirtualNetworkConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_azure_private_endpoint_command" "%[2]s" {
  organization_id     = "00000000-0000-0000-0000-000000000000"
  project_id          = "11111111-1111-1111-1111-111111111111"
  cluster_id          = "22222222-2222-2222-2222-222222222222"
  resource_group_name = "my-resource-group"
  virtual_network     = "a"
}
`, globalProviderBlock, dataSourceName)
}

// TestAccGCPPrivateEndpointCommandInvalidVPCNetworkID verifies that the GCP private endpoint
// command data source rejects a vpc_network_id shorter than the OpenAPI minimum length at
// plan time. The generated LengthBetween(12, 21) validator fires before any API call.
func TestAccGCPPrivateEndpointCommandInvalidVPCNetworkID(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_gcp_pe_command_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccGCPPrivateEndpointCommandInvalidVPCNetworkIDConfig(dataSourceName),
				ExpectError: regexp.MustCompile(`(?s)vpc_network_id.*string length must be between 12 and 21`),
			},
		},
	})
}

func testAccGCPPrivateEndpointCommandInvalidVPCNetworkIDConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_gcp_private_endpoint_command" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  vpc_network_id  = "vpc-short"
  subnet_ids      = ["subnet-1234567890abcdef0"]
}
`, globalProviderBlock, dataSourceName)
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
