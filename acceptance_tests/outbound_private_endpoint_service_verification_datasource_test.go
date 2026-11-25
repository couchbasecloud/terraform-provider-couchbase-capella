package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccOutboundPrivateEndpointServiceVerificationDataSource checks an AWS endpoint service that
// doesn't exist against the shared test cluster, which is on AWS. The name is well formed, so
// format passes, but the service is never found: the read must still succeed with ready false,
// since a failed check is a 200 and not an error.
func TestAccOutboundPrivateEndpointServiceVerificationDataSource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_service_verification_")
	resourceReference := "data.couchbase-capella_outbound_private_endpoint_service_verification." + resourceName
	serviceName := "com.amazonaws.vpce.us-east-1.vpce-svc-00000000000000000"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOutboundPrivateEndpointServiceVerificationConfig(resourceName, serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "service_name", serviceName),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttrSet(resourceReference, "cloud_provider.region"),
					resource.TestCheckResourceAttr(resourceReference, "ready", "false"),
					resource.TestCheckResourceAttr(resourceReference, "checks.#", "6"),
					resource.TestCheckResourceAttr(resourceReference, "checks.0.name", "format"),
					resource.TestCheckResourceAttr(resourceReference, "checks.0.result", "passed"),
					resource.TestCheckResourceAttr(resourceReference, "checks.1.name", "region"),
					resource.TestCheckResourceAttr(resourceReference, "checks.2.name", "accessGranted"),
					resource.TestCheckResourceAttr(resourceReference, "checks.3.name", "serviceZones"),
					resource.TestCheckResourceAttr(resourceReference, "checks.4.name", "notAlreadyConnected"),
					resource.TestCheckResourceAttr(resourceReference, "checks.5.name", "privateDnsName"),
					// A service that is never found has nothing to report.
					resource.TestCheckNoResourceAttr(resourceReference, "provider_config"),
				),
			},
		},
	})
}

func testAccOutboundPrivateEndpointServiceVerificationConfig(resourceName, serviceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_outbound_private_endpoint_service_verification" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  service_name    = "%[6]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, serviceName)
}
