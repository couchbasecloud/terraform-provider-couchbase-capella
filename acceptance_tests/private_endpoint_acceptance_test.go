package acceptance_tests

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// TestAccPrivateEndpointsDataSourceWithEndpoint is an end-to-end test for the
// couchbase-capella_private_endpoints data source against a NON-EMPTY endpoint
// list. It enables the private endpoint service, creates a real AWS interface
// VPC endpoint against the service, waits for Capella to register the
// connection, then reads the data source and accepts the endpoint.
//
// This is the path that regressed: the data source's nested `data` schema did
// not match the PrivateEndpointData struct, so state conversion failed with a
// "Value Conversion Error" the moment the list contained an endpoint. The
// existing TestAccPrivateEndpointServiceEnableDisable only covers the empty
// list (data.# = 0), so it never exercised the mismatch.
//
// It runs only when the AWS inputs are supplied (and AWS credentials are set in
// the environment for the aws provider). The VPC must be in the same region as
// the Capella cluster:
//
//	ACC_AWS_REGION      e.g. us-east-1 (must match the cluster region)
//	ACC_AWS_VPC_ID      e.g. vpc-0123456789abcdef0
//	ACC_AWS_SUBNET_IDS  comma-separated, e.g. subnet-aaa,subnet-bbb
//	ACC_AWS_VPC_CIDR    e.g. 10.0.0.0/16
func TestAccPrivateEndpointsDataSourceWithEndpoint(t *testing.T) {
	region := os.Getenv("ACC_AWS_REGION")
	vpcID := os.Getenv("ACC_AWS_VPC_ID")
	subnetCSV := os.Getenv("ACC_AWS_SUBNET_IDS")
	vpcCIDR := os.Getenv("ACC_AWS_VPC_CIDR")
	if region == "" || vpcID == "" || subnetCSV == "" || vpcCIDR == "" {
		t.Skip("skipping: set ACC_AWS_REGION, ACC_AWS_VPC_ID, ACC_AWS_SUBNET_IDS (comma-separated) and ACC_AWS_VPC_CIDR (plus AWS credentials) to run this end-to-end test")
	}

	serviceName := randomStringWithPrefix("tf_acc_pe_svc_")
	dsName := randomStringWithPrefix("tf_acc_pe_ds_")
	acceptName := randomStringWithPrefix("tf_acc_pe_accept_")
	serviceRef := "couchbase-capella_private_endpoint_service." + serviceName
	dsRef := "data.couchbase-capella_private_endpoints." + dsName
	acceptRef := "couchbase-capella_private_endpoints." + acceptName
	vpceRef := "aws_vpc_endpoint.capella"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		ExternalProviders: map[string]resource.ExternalProvider{
			"aws":  {Source: "hashicorp/aws"},
			"time": {Source: "hashicorp/time"},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateEndpointsWithEndpointConfig(
					serviceName, dsName, acceptName, region, vpcID, subnetListHCL(subnetCSV), vpcCIDR,
					globalOrgId, globalProjectId, globalClusterId,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(serviceRef, "enabled", "true"),
					resource.TestCheckResourceAttrSet(serviceRef, "service_name"),
					resource.TestCheckResourceAttrSet(vpceRef, "id"),
					resource.TestCheckResourceAttrSet(dsRef, "private_endpoint_dns"),
					// Regression guard: converting a NON-EMPTY endpoint list to
					// state must succeed. Pre-fix this step errored with a "Value
					// Conversion Error" before any check below could run.
					resource.TestCheckResourceAttrSet(dsRef, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsRef, "data.0.status"),
					resource.TestCheckResourceAttrSet(dsRef, "data.0.service_name"),
					testAccCheckPrivateEndpointInList(dsRef, vpceRef),
					// The accept (associate) resource captured a status.
					resource.TestCheckResourceAttrSet(acceptRef, "status"),
				),
			},
		},
	})
}

// subnetListHCL turns a comma-separated list of subnet IDs into the body of an
// HCL list, e.g. "subnet-a, subnet-b" -> `"subnet-a", "subnet-b"`.
func subnetListHCL(csv string) string {
	parts := strings.Split(csv, ",")
	quoted := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			quoted = append(quoted, strconv.Quote(s))
		}
	}
	return strings.Join(quoted, ", ")
}

// testAccCheckPrivateEndpointInList asserts the private endpoints data source
// returned a non-empty list that includes the VPC endpoint we just created.
func testAccCheckPrivateEndpointInList(dsRef, vpceRef string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vpce, ok := s.RootModule().Resources[vpceRef]
		if !ok {
			return fmt.Errorf("resource %s not found in state", vpceRef)
		}
		vpceID := vpce.Primary.Attributes["id"]

		ds, ok := s.RootModule().Resources[dsRef]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dsRef)
		}

		count, err := strconv.Atoi(ds.Primary.Attributes["data.#"])
		if err != nil {
			return fmt.Errorf("could not read %s data.#: %w", dsRef, err)
		}
		if count < 1 {
			return fmt.Errorf("expected at least one endpoint in %s, got %d (connection may not have registered yet; increase the time_sleep duration)", dsRef, count)
		}

		for i := 0; i < count; i++ {
			if ds.Primary.Attributes[fmt.Sprintf("data.%d.id", i)] == vpceID {
				return nil
			}
		}
		return fmt.Errorf("endpoint %s not present in %s list of %d", vpceID, dsRef, count)
	}
}

// testAccPrivateEndpointsWithEndpointConfig builds the end-to-end config:
// enable the service, create a real interface VPC endpoint against it, wait for
// the connection to register, then list (data source) and accept (resource).
func testAccPrivateEndpointsWithEndpointConfig(
	serviceName, dsName, acceptName, region, vpcID, subnetListHCL, vpcCIDR, orgID, projectID, clusterID string,
) string {
	return fmt.Sprintf(`
%[1]s

provider "aws" {
  region = %[5]q
}

resource "couchbase-capella_private_endpoint_service" %[2]q {
  organization_id = %[9]q
  project_id      = %[10]q
  cluster_id      = %[11]q
  enabled         = true
}

resource "aws_security_group" "capella_pe" {
  name_prefix = "tf-acc-capella-pe-"
  description = "TLS access from the VPC to the Capella private endpoint (acc test)"
  vpc_id      = %[6]q

  ingress {
    description = "Allow from the VPC"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [%[8]q]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_vpc_endpoint" "capella" {
  vpc_id              = %[6]q
  service_name        = couchbase-capella_private_endpoint_service.%[2]s.service_name
  vpc_endpoint_type   = "Interface"
  subnet_ids          = [%[7]s]
  security_group_ids  = [aws_security_group.capella_pe.id]
  private_dns_enabled = false
}

# Give Capella time to register the connection request before we list/accept it.
resource "time_sleep" "wait_for_registration" {
  depends_on      = [aws_vpc_endpoint.capella]
  create_duration = "90s"
}

data "couchbase-capella_private_endpoints" %[3]q {
  organization_id = %[9]q
  project_id      = %[10]q
  cluster_id      = %[11]q

  depends_on = [time_sleep.wait_for_registration]
}

resource "couchbase-capella_private_endpoints" %[4]q {
  organization_id = %[9]q
  project_id      = %[10]q
  cluster_id      = %[11]q
  endpoint_id     = aws_vpc_endpoint.capella.id

  depends_on = [time_sleep.wait_for_registration]
}
`, globalProviderBlock, serviceName, dsName, acceptName, region, vpcID, subnetListHCL, vpcCIDR, orgID, projectID, clusterID)
}
