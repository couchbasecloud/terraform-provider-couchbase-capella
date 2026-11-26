package acceptance_tests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const outboundPrivateEndpointServiceVerificationType = "data.couchbase-capella_outbound_private_endpoint_service_verification."

// TestAccOutboundPrivateEndpointServiceVerificationDataSourceInvalidFormat checks a service name
// that isn't an AWS endpoint service name against the shared test cluster, which is on AWS. Format
// fails, so no later check runs and nothing reaches AWS: the answer is the same for any cluster
// region. A failed check is a 200, so the read still succeeds, with ready false.
func TestAccOutboundPrivateEndpointServiceVerificationDataSourceInvalidFormat(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_service_verification_")
	resourceReference := outboundPrivateEndpointServiceVerificationType + resourceName
	serviceName := "not-an-endpoint-service"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOutboundPrivateEndpointServiceVerificationConfig(resourceName, globalProjectId, globalClusterId, serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "service_name", serviceName),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttrSet(resourceReference, "cloud_provider.region"),
					resource.TestCheckResourceAttr(resourceReference, "ready", "false"),
					testAccOutboundPrivateEndpointServiceVerificationChecks(resourceReference, []outboundPrivateEndpointVerificationCheck{
						{"format", "failed", "Enter an endpoint service name, com.amazonaws.vpce.<region>.vpce-svc-<id>."},
						{"region", "notApplicable", ""},
						{"accessGranted", "notApplicable", ""},
						{"serviceZones", "notApplicable", ""},
						{"notAlreadyConnected", "notApplicable", ""},
						{"privateDnsName", "notApplicable", ""},
					}),
					resource.TestCheckNoResourceAttr(resourceReference, "provider_config"),
				),
			},
		},
	})
}

// TestAccOutboundPrivateEndpointServiceVerificationDataSource checks a real endpoint service, and
// one in the same region that doesn't exist. It needs ACC_OUTBOUND_PE_PROJECT_ID,
// ACC_OUTBOUND_PE_CLUSTER_ID (an AWS cluster on Couchbase Server 8.5 or later with the
// SelfManagedPrivateEndpoints flag on) and ACC_OUTBOUND_PE_SERVICE_NAME (an AWS endpoint service
// in the cluster's region that allows Capella's AWS account, accepts connections itself and has
// no endpoint from this cluster).
func TestAccOutboundPrivateEndpointServiceVerificationDataSource(t *testing.T) {
	projectId := os.Getenv("ACC_OUTBOUND_PE_PROJECT_ID")
	clusterId := os.Getenv("ACC_OUTBOUND_PE_CLUSTER_ID")
	serviceName := os.Getenv("ACC_OUTBOUND_PE_SERVICE_NAME")
	if projectId == "" || clusterId == "" || serviceName == "" {
		t.Skip("skipping: set ACC_OUTBOUND_PE_PROJECT_ID, ACC_OUTBOUND_PE_CLUSTER_ID and ACC_OUTBOUND_PE_SERVICE_NAME to run this test")
	}
	// com.amazonaws.vpce.<region>.vpce-svc-<id>
	parts := strings.Split(serviceName, ".")
	if len(parts) != 5 {
		t.Fatalf("ACC_OUTBOUND_PE_SERVICE_NAME %q is not com.amazonaws.vpce.<region>.vpce-svc-<id>", serviceName)
	}
	region := parts[3]
	missingServiceName := fmt.Sprintf("com.amazonaws.vpce.%s.vpce-svc-00000000000000000", region)

	found := randomStringWithPrefix("tf_acc_outbound_pe_service_verification_")
	foundReference := outboundPrivateEndpointServiceVerificationType + found
	missing := randomStringWithPrefix("tf_acc_outbound_pe_service_verification_")
	missingReference := outboundPrivateEndpointServiceVerificationType + missing

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOutboundPrivateEndpointServiceVerificationConfig(found, projectId, clusterId, serviceName) +
					testAccOutboundPrivateEndpointServiceVerificationBlock(missing, projectId, clusterId, missingServiceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(foundReference, "service_name", serviceName),
					resource.TestCheckResourceAttr(foundReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(foundReference, "cloud_provider.region", region),
					resource.TestCheckResourceAttr(foundReference, "ready", "true"),
					testAccOutboundPrivateEndpointServiceVerificationChecks(foundReference, []outboundPrivateEndpointVerificationCheck{
						{"format", "passed", ""},
						{"region", "passed", ""},
						{"accessGranted", "passed", ""},
						{"serviceZones", "passed", ""},
						{"notAlreadyConnected", "passed", ""},
					}),
					// The last check depends on the service's private DNS name, so only its name is fixed.
					resource.TestCheckResourceAttr(foundReference, "checks.5.name", "privateDnsName"),
					resource.TestMatchResourceAttr(foundReference, "checks.5.result", regexp.MustCompile(`^(passed|warning)$`)),
					resource.TestMatchResourceAttr(foundReference, "provider_config.aws.availability_zones.0", regexp.MustCompile(`^`+regexp.QuoteMeta(region)+`[a-z]$`)),
					resource.TestCheckResourceAttr(foundReference, "provider_config.aws.acceptance_required", "false"),
					resource.TestMatchResourceAttr(foundReference, "provider_config.aws.service_owner", regexp.MustCompile(`^\d{12}$`)),
					resource.TestMatchResourceAttr(foundReference, "provider_config.aws.private_dns_name_verification", regexp.MustCompile(`^(none|pendingVerification|verified|failed)$`)),

					resource.TestCheckResourceAttr(missingReference, "service_name", missingServiceName),
					resource.TestCheckResourceAttr(missingReference, "ready", "false"),
					testAccOutboundPrivateEndpointServiceVerificationChecks(missingReference, []outboundPrivateEndpointVerificationCheck{
						{"format", "passed", ""},
						{"region", "passed", ""},
						{"accessGranted", "failed", fmt.Sprintf("Service not found in %s. Check the name, and add the Capella AWS account to its allowed principals.", region)},
						{"serviceZones", "notApplicable", ""},
						{"notAlreadyConnected", "passed", ""},
						{"privateDnsName", "notApplicable", ""},
					}),
					// A service AWS can't see has nothing to report.
					resource.TestCheckNoResourceAttr(missingReference, "provider_config"),
				),
			},
		},
	})
}

type outboundPrivateEndpointVerificationCheck struct {
	name, result, message string
}

// testAccOutboundPrivateEndpointServiceVerificationChecks asserts each check in order. An empty
// message must be absent, not just empty.
func testAccOutboundPrivateEndpointServiceVerificationChecks(resourceReference string, want []outboundPrivateEndpointVerificationCheck) resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for i, c := range want {
		prefix := fmt.Sprintf("checks.%d.", i)
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceReference, prefix+"name", c.name),
			resource.TestCheckResourceAttr(resourceReference, prefix+"result", c.result),
		)
		if c.message == "" {
			checks = append(checks, resource.TestCheckNoResourceAttr(resourceReference, prefix+"message"))
		} else {
			checks = append(checks, resource.TestCheckResourceAttr(resourceReference, prefix+"message", c.message))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(append(checks, resource.TestCheckResourceAttr(resourceReference, "checks.#", "6"))...)
}

func testAccOutboundPrivateEndpointServiceVerificationConfig(resourceName, projectId, clusterId, serviceName string) string {
	return globalProviderBlock + testAccOutboundPrivateEndpointServiceVerificationBlock(resourceName, projectId, clusterId, serviceName)
}

func testAccOutboundPrivateEndpointServiceVerificationBlock(resourceName, projectId, clusterId, serviceName string) string {
	return fmt.Sprintf(`
data "couchbase-capella_outbound_private_endpoint_service_verification" "%[1]s" {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  service_name    = %[5]q
}
`, resourceName, globalOrgId, projectId, clusterId, serviceName)
}
