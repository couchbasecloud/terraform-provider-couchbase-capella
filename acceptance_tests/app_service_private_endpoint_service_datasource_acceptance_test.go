package acceptance_tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceAppServicePrivateEndpointService reads the private endpoint service state of the
// global App Service. Without the aws block, no command is requested, so command stays null.
func TestAccDatasourceAppServicePrivateEndpointService(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_pe_service_")
	dataSourceReference := "data.couchbase-capella_app_service_private_endpoint_service." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicePrivateEndpointServiceDatasourceConfig(dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "state"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "target_state"),
					resource.TestCheckNoResourceAttr(dataSourceReference, "command"),
				),
			},
		},
	})
}

// appServicePrivateEndpointServiceAWSCommandSteps reads the state and requests the AWS VPC endpoint
// command. Gated on ACC_AWS_VPC_ID / ACC_AWS_SUBNET_IDS; it skips otherwise.
//
// The privateEndpointCommand endpoint is only reachable once the private endpoint service is enabled
// on the App Service - it returns 400 "Private endpoints aren't enabled for this cluster" otherwise -
// so the config enables the service and the data source depends on it. The depends_on also defers the
// read to apply time, so the pre-apply plan does not call the API before the service exists.
//
// These steps run as a subtest of TestAccAppServicePrivateEndpoint, which serialises them against the
// other tests that enable and disable the service on the same shared App Service.
func appServicePrivateEndpointServiceAWSCommandSteps(t *testing.T) []resource.TestStep {
	vpcID := os.Getenv("ACC_AWS_VPC_ID")
	subnetIDs := parseSubnetIDs(os.Getenv("ACC_AWS_SUBNET_IDS"))
	if vpcID == "" || len(subnetIDs) == 0 {
		t.Skip("skipping: set ACC_AWS_VPC_ID and ACC_AWS_SUBNET_IDS (comma-separated, at least one subnet id) to run this test")
	}

	serviceName := randomStringWithPrefix("tf_acc_app_service_pe_command_service_")
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_pe_command_")
	dataSourceReference := "data.couchbase-capella_app_service_private_endpoint_service." + dataSourceName

	return []resource.TestStep{
		{
			Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_private_endpoint_service" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  app_service_id  = "%[7]s"
}

data "couchbase-capella_app_service_private_endpoint_service" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  app_service_id  = "%[7]s"

  aws = {
    vpc_id     = "%[8]s"
    subnet_ids = [%[9]s]
  }

  depends_on = [couchbase-capella_app_service_private_endpoint_service.%[2]s]
}
`, globalProviderBlock, serviceName, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, vpcID, subnetsToHCL(subnetIDs)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(dataSourceReference, "aws.vpc_id", vpcID),
				resource.TestCheckResourceAttr(dataSourceReference, "aws.subnet_ids.#", fmt.Sprint(len(subnetIDs))),
				// The config enables the service and the data source depends on it, so the
				// read always sees it enabled.
				resource.TestCheckResourceAttr(dataSourceReference, "state", "enabled"),
				resource.TestCheckResourceAttrSet(dataSourceReference, "target_state"),
				resource.TestCheckResourceAttrSet(dataSourceReference, "command"),
			),
		},
	}
}

// TestAccDatasourceAppServicePrivateEndpointServiceInvalidVPCID verifies that a vpc_id shorter than the
// OpenAPI minimum length fails at plan time, before any API call, so dummy IDs are sufficient.
func TestAccDatasourceAppServicePrivateEndpointServiceInvalidVPCID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_private_endpoint_service" "test" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  app_service_id  = "33333333-3333-3333-3333-333333333333"

  aws = {
    vpc_id     = "vpc-short"
    subnet_ids = ["subnet-1234567890abcdef0"]
  }
}
`, globalProviderBlock),
				ExpectError: regexp.MustCompile(`(?s)vpc_id.*string length must be between 12 and 21`),
			},
		},
	})
}

func testAccAppServicePrivateEndpointServiceDatasourceConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_private_endpoint_service" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}
`, globalProviderBlock, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}
