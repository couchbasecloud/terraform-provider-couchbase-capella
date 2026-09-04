package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppServicePrivateEndpointServiceEnableDisable tests enabling, reading, and disabling
// the private endpoint service for an App Service (Sync Gateway). It reuses the shared App
// Service created by TestMain (globalAppServiceId) rather than creating its own.
func TestAccAppServicePrivateEndpointServiceEnableDisable(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_private_endpoint_service_")
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_private_endpoints_ds_")
	resourceReference := "couchbase-capella_app_service_private_endpoint_service." + resourceName
	dataSourceReference := "data.couchbase-capella_app_service_private_endpoints." + dataSourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicePrivateEndpointsDataSourceNoEndpointConfig(resourceName, dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "state", "enabled"),
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(dataSourceReference, "data.#", "0"),
				),
			},
			{
				Config: testAccAppServicePrivateEndpointServiceEnableConfig(resourceName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "state", "disabled"),
				),
			},
		},
	})
}

// TestAccAppServiceAWSPrivateEndpointCommandInvalidVPCID verifies that the App Service AWS
// private endpoint command data source rejects a vpc_id shorter than the OpenAPI minimum
// length at plan time, mirroring the operational-cluster equivalent.
func TestAccAppServiceAWSPrivateEndpointCommandInvalidVPCID(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_aws_pe_command_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceAWSPrivateEndpointCommandInvalidVPCIDConfig(dataSourceName),
				ExpectError: regexp.MustCompile(`(?s)vpc_id.*string length must be between 12 and 21`),
			},
		},
	})
}

func testAccAppServiceAWSPrivateEndpointCommandInvalidVPCIDConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_aws_private_endpoint_command" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  app_service_id  = "33333333-3333-3333-3333-333333333333"
  vpc_id          = "vpc-short"
  subnet_ids      = ["subnet-1234567890abcdef0"]
}
`, globalProviderBlock, dataSourceName)
}

// TestAccAppServiceAzurePrivateEndpointCommandInvalidVirtualNetwork verifies that the App
// Service Azure private endpoint command data source rejects a virtual_network shorter than
// the OpenAPI minimum length at plan time.
func TestAccAppServiceAzurePrivateEndpointCommandInvalidVirtualNetwork(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_azure_pe_command_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceAzurePrivateEndpointCommandInvalidVirtualNetworkConfig(dataSourceName),
				ExpectError: regexp.MustCompile(`(?s)virtual_network.*string length must be between 2 and 64`),
			},
		},
	})
}

func testAccAppServiceAzurePrivateEndpointCommandInvalidVirtualNetworkConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_azure_private_endpoint_command" "%[2]s" {
  organization_id     = "00000000-0000-0000-0000-000000000000"
  project_id          = "11111111-1111-1111-1111-111111111111"
  cluster_id          = "22222222-2222-2222-2222-222222222222"
  app_service_id      = "33333333-3333-3333-3333-333333333333"
  resource_group_name = "my-resource-group"
  virtual_network     = "a"
}
`, globalProviderBlock, dataSourceName)
}

// TestAccAppServiceGCPPrivateEndpointCommandInvalidVPCNetworkID verifies that the App Service
// GCP private endpoint command data source rejects a vpc_network_id shorter than the OpenAPI
// minimum length at plan time.
func TestAccAppServiceGCPPrivateEndpointCommandInvalidVPCNetworkID(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_gcp_pe_command_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceGCPPrivateEndpointCommandInvalidVPCNetworkIDConfig(dataSourceName),
				ExpectError: regexp.MustCompile(`(?s)vpc_network_id.*string length must be between 12 and 21`),
			},
		},
	})
}

func testAccAppServiceGCPPrivateEndpointCommandInvalidVPCNetworkIDConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_gcp_private_endpoint_command" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  app_service_id  = "33333333-3333-3333-3333-333333333333"
  vpc_network_id  = "vpc-short"
  subnet_ids      = ["subnet-1234567890abcdef0"]
}
`, globalProviderBlock, dataSourceName)
}

// testAccAppServicePrivateEndpointServiceEnableConfig returns terraform config for
// enabling/disabling the App Service private endpoint service.
func testAccAppServicePrivateEndpointServiceEnableConfig(resourceName string, enabled bool) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_app_service_private_endpoint_service" "%[2]s" {
			organization_id = "%[3]s"
			project_id      = "%[4]s"
			cluster_id      = "%[5]s"
			app_service_id  = "%[6]s"
			enabled         = %[7]t
		}
		`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, enabled)
}

func testAccAppServicePrivateEndpointsDataSourceNoEndpointConfig(serviceResourceName, dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_private_endpoint_service" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  app_service_id  = "%[7]s"
  enabled         = true
}

data "couchbase-capella_app_service_private_endpoints" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  app_service_id  = "%[7]s"

  depends_on = [couchbase-capella_app_service_private_endpoint_service.%[2]s]
}
`, globalProviderBlock, serviceResourceName, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

// TestAccAppServicePrivateEndpointsInvalidEndpointID verifies that the App Service private
// endpoints resource rejects an empty endpoint_id at plan time.
func TestAccAppServicePrivateEndpointsInvalidEndpointID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_private_endpoints_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicePrivateEndpointsInvalidEndpointIDConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)endpoint_id.*string length must be at least 1`),
			},
		},
	})
}

func testAccAppServicePrivateEndpointsInvalidEndpointIDConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_private_endpoints" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  app_service_id  = "33333333-3333-3333-3333-333333333333"
  endpoint_id     = ""
}
`, globalProviderBlock, resourceName)
}

// NOTE: unlike the operational cluster's private_endpoint_service (which exposes
// service_name), the App Service private endpoint service status API only returns
// state/targetState (see internal/api/app_service_private_endpoint_service.go). There is
// therefore no Terraform-native way to discover the CSP-side VPC endpoint service name to
// wire a real aws_vpc_endpoint resource against, so an end-to-end
// "TestAccAppServicePrivateEndpointsDataSourceWithEndpoint" test analogous to the cluster-level
// TestAccPrivateEndpointsDataSourceWithEndpoint is intentionally omitted here. A human reviewer
// at Couchbase should confirm whether the live API is expected to stay this way, or whether a
// service_name-equivalent field should be added so an end-to-end test can be written.
