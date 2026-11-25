package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

const outboundPrivateEndpointType = "couchbase-capella_outbound_private_endpoint"

// outboundPrivateEndpointEnv is what the end-to-end tests need beyond the suite's globals. The
// shared test cluster runs a Server version older than 8.5, which outbound private endpoints
// need, so the tests run against their own cluster.
type outboundPrivateEndpointEnv struct {
	projectId   string
	clusterId   string
	serviceName string
}

// outboundPrivateEndpointTestEnv skips the test unless ACC_OUTBOUND_PE_PROJECT_ID,
// ACC_OUTBOUND_PE_CLUSTER_ID and ACC_OUTBOUND_PE_SERVICE_NAME are set. The cluster must be an AWS
// cluster on Couchbase Server 8.5 or later with the SelfManagedPrivateEndpoints flag on, and the
// service an AWS endpoint service in the cluster's region whose allowed principals include
// Capella's AWS account (see the couchbase-capella_cloud_accounts data source).
func outboundPrivateEndpointTestEnv(t *testing.T) outboundPrivateEndpointEnv {
	t.Helper()
	env := outboundPrivateEndpointEnv{
		projectId:   os.Getenv("ACC_OUTBOUND_PE_PROJECT_ID"),
		clusterId:   os.Getenv("ACC_OUTBOUND_PE_CLUSTER_ID"),
		serviceName: os.Getenv("ACC_OUTBOUND_PE_SERVICE_NAME"),
	}
	if env.projectId == "" || env.clusterId == "" || env.serviceName == "" {
		t.Skip("skipping: set ACC_OUTBOUND_PE_PROJECT_ID, ACC_OUTBOUND_PE_CLUSTER_ID (an AWS cluster on Couchbase Server 8.5 or later " +
			"with the SelfManagedPrivateEndpoints flag on) and ACC_OUTBOUND_PE_SERVICE_NAME (an AWS endpoint service in the " +
			"cluster's region that allows Capella's AWS account) to run this test")
	}
	return env
}

// TestAccOutboundPrivateEndpointResource creates an endpoint with an assigned hostname, imports
// it and renames it in place.
func TestAccOutboundPrivateEndpointResource(t *testing.T) {
	env := outboundPrivateEndpointTestEnv(t)

	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_")
	resourceReference := outboundPrivateEndpointType + "." + resourceName
	name := randomStringWithPrefix("tf-acc-outbound-pe-")
	renamed := name + "-renamed"
	var createdId string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckOutboundPrivateEndpointDestroy(t),
		Steps: []resource.TestStep{
			// Create and Read. endpoint_type is left out to check its default.
			{
				Config: testAccOutboundPrivateEndpointResourceConfig(resourceName, env, name, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccOutboundPrivateEndpointAssignedHostnameChecks(resourceReference, env, name),
					testAccOutboundPrivateEndpointRecordId(resourceReference, &createdId),
				),
			},
			// ImportState. The service can accept the connection between create and import, which
			// moves the status, the addresses and private DNS, and may write the record's modified
			// audit fields, so those aren't compared. The steps' checks cover them.
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateOutboundPrivateEndpointImportIdForResource(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore:              outboundPrivateEndpointImportVerifyIgnore,
			},
			// Update. A rename keeps the endpoint, so every other field is unchanged.
			{
				Config: testAccOutboundPrivateEndpointResourceConfig(resourceName, env, renamed, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccOutboundPrivateEndpointAssignedHostnameChecks(resourceReference, env, renamed),
					testAccOutboundPrivateEndpointKeepsId(resourceReference, &createdId),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccOutboundPrivateEndpointResourceCustomHostname creates an endpoint with the customer's own
// hostname, which turns private DNS on.
//
// Also gated on ACC_OUTBOUND_PE_HOSTNAME: the endpoint service's verified private DNS name, or a
// name under it when that name is a wildcard.
func TestAccOutboundPrivateEndpointResourceCustomHostname(t *testing.T) {
	env := outboundPrivateEndpointTestEnv(t)
	hostname := os.Getenv("ACC_OUTBOUND_PE_HOSTNAME")
	if hostname == "" {
		t.Skip("skipping: set ACC_OUTBOUND_PE_HOSTNAME to the endpoint service's verified private DNS name to run this test")
	}

	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_")
	resourceReference := outboundPrivateEndpointType + "." + resourceName
	name := randomStringWithPrefix("tf-acc-outbound-pe-")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckOutboundPrivateEndpointDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOutboundPrivateEndpointResourceConfig(resourceName, env, name, hostname),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccOutboundPrivateEndpointCommonChecks(resourceReference, env, name),
					resource.TestCheckResourceAttr(resourceReference, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceReference, "is_custom_hostname", "true"),
					// A service that accepts each connection itself leaves create at pendingAcceptance,
					// and AWS only turns private DNS on once the endpoint is available.
					testAccOutboundPrivateEndpointByStatus(resourceReference,
						resource.TestCheckResourceAttr(resourceReference, "private_dns_enabled", "false"),
						resource.TestCheckResourceAttr(resourceReference, "private_dns_enabled", "true"),
					),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateOutboundPrivateEndpointImportIdForResource(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore:              outboundPrivateEndpointImportVerifyIgnore,
			},
		},
	})
}

// TestAccOutboundPrivateEndpointInvalidServiceName checks the API's answer to a service name that
// isn't an AWS endpoint service name.
func TestAccOutboundPrivateEndpointInvalidServiceName(t *testing.T) {
	env := outboundPrivateEndpointTestEnv(t)
	env.serviceName = "vpce-svc-0123456789abcdef0"

	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOutboundPrivateEndpointResourceConfig(resourceName, env, "tf-acc-outbound-pe", ""),
				// Match the code: Terraform wraps the message mid-phrase.
				ExpectError: regexp.MustCompile(`(?s)Could not create the outbound private endpoint.*"code":16006`),
			},
		},
	})
}

// The tests below are refused at plan time, before any API call, so they run without the
// end-to-end environment.

func TestAccOutboundPrivateEndpointInvalidName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_")
	env := outboundPrivateEndpointEnv{projectId: globalProjectId, clusterId: globalClusterId, serviceName: "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"}
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccOutboundPrivateEndpointResourceConfig(resourceName, env, "a", ""),
				ExpectError: regexp.MustCompile(`(?s)name.*string length must be between 2 and 128`),
			},
		},
	})
}

func TestAccOutboundPrivateEndpointInvalidEmptyServiceName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_")
	env := outboundPrivateEndpointEnv{projectId: globalProjectId, clusterId: globalClusterId, serviceName: ""}
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccOutboundPrivateEndpointResourceConfig(resourceName, env, "tf-acc-outbound-pe", ""),
				ExpectError: regexp.MustCompile(`(?s)service_name.*string length must be between 1 and 512`),
			},
		},
	})
}

func TestAccOutboundPrivateEndpointInvalidHostname(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_")
	env := outboundPrivateEndpointEnv{projectId: globalProjectId, clusterId: globalClusterId, serviceName: "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"}
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccOutboundPrivateEndpointResourceConfig(resourceName, env, "tf-acc-outbound-pe", strings.Repeat("a", 254)),
				ExpectError: regexp.MustCompile(`(?s)hostname.*string length must be between 1 and 253`),
			},
		},
	})
}

// testAccOutboundPrivateEndpointCommonChecks asserts the fields that don't depend on the hostname.
func testAccOutboundPrivateEndpointCommonChecks(resourceReference string, env outboundPrivateEndpointEnv, name string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceReference, "id"),
		resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
		resource.TestCheckResourceAttr(resourceReference, "project_id", env.projectId),
		resource.TestCheckResourceAttr(resourceReference, "cluster_id", env.clusterId),
		resource.TestCheckResourceAttr(resourceReference, "name", name),
		resource.TestCheckResourceAttr(resourceReference, "service_name", env.serviceName),
		resource.TestCheckResourceAttr(resourceReference, "endpoint_type", "interface"),
		resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
		resource.TestCheckResourceAttrSet(resourceReference, "cloud_provider.region"),
		// Create returns once the connection request reaches the service. A service that accepts
		// connections on its own is already linked.
		resource.TestMatchResourceAttr(resourceReference, "status", regexp.MustCompile(`^(pendingAcceptance|linked)$`)),
		resource.TestMatchResourceAttr(resourceReference, "provider_endpoint_id", regexp.MustCompile(`^vpce-[0-9a-f]+$`)),
		resource.TestCheckNoResourceAttr(resourceReference, "failure_reason"),
		resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
		resource.TestCheckResourceAttrSet(resourceReference, "audit.created_by"),
		resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
		resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_by"),
		resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
		testAccOutboundPrivateEndpointAWSChecks(resourceReference),
	)
}

// outboundPrivateEndpointImportVerifyIgnore are the fields that can change between create and
// import without Terraform doing anything: the service accepting the connection moves them.
var outboundPrivateEndpointImportVerifyIgnore = []string{
	"status", "provider_status", "ip_addresses", "provider_config", "private_dns_enabled",
	"audit.modified_at", "audit.modified_by", "audit.version",
}

var (
	outboundPrivateEndpointIPv4 = regexp.MustCompile(`^\d{1,3}(\.\d{1,3}){3}$`)
	awsAvailabilityZone         = regexp.MustCompile(`^[a-z]{2}(-[a-z]+)+-\d[a-z]$`)
)

// testAccOutboundPrivateEndpointAWSChecks asserts the AWS details. The zones are known from
// create, while the endpoint's addresses, its provider status and its regional DNS name come once
// the service has accepted the connection, so they're only required when the endpoint is linked.
func testAccOutboundPrivateEndpointAWSChecks(resourceReference string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestMatchResourceAttr(resourceReference, "provider_config.aws.availability_zones.0", awsAvailabilityZone),
		testAccOutboundPrivateEndpointByStatus(resourceReference,
			resource.ComposeAggregateTestCheckFunc(),
			resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "provider_status", "available"),
				resource.TestMatchResourceAttr(resourceReference, "provider_config.aws.regional_dns_name", regexp.MustCompile(`\.vpce\.amazonaws\.com$`)),
				resource.TestMatchResourceAttr(resourceReference, "ip_addresses.0", outboundPrivateEndpointIPv4),
				resource.TestMatchResourceAttr(resourceReference, "provider_config.aws.ip_addresses.0.ip_address", outboundPrivateEndpointIPv4),
				resource.TestMatchResourceAttr(resourceReference, "provider_config.aws.ip_addresses.0.availability_zone", awsAvailabilityZone),
			),
		),
	)
}

// testAccOutboundPrivateEndpointByStatus runs linked once the service has accepted the connection,
// and pending before then.
func testAccOutboundPrivateEndpointByStatus(resourceReference string, pending, linked resource.TestCheckFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		if rs.Primary.Attributes["status"] == "linked" {
			return linked(s)
		}
		return pending(s)
	}
}

// testAccOutboundPrivateEndpointAssignedHostnameChecks asserts an endpoint created without a
// hostname: Capella assigns one and private DNS stays off.
func testAccOutboundPrivateEndpointAssignedHostnameChecks(resourceReference string, env outboundPrivateEndpointEnv, name string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testAccOutboundPrivateEndpointCommonChecks(resourceReference, env, name),
		resource.TestCheckResourceAttrSet(resourceReference, "hostname"),
		resource.TestCheckResourceAttr(resourceReference, "is_custom_hostname", "false"),
		resource.TestCheckResourceAttr(resourceReference, "private_dns_enabled", "false"),
	)
}

// testAccOutboundPrivateEndpointRecordId keeps the created endpoint's id for
// testAccOutboundPrivateEndpointKeepsId.
func testAccOutboundPrivateEndpointRecordId(resourceReference string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		*id = rs.Primary.ID
		return nil
	}
}

// testAccOutboundPrivateEndpointKeepsId fails if a rename replaced the created endpoint.
func testAccOutboundPrivateEndpointKeepsId(resourceReference string, createdId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		if rs.Primary.ID != *createdId {
			return fmt.Errorf("endpoint id changed from %s to %s on a rename", *createdId, rs.Primary.ID)
		}
		return nil
	}
}

func generateOutboundPrivateEndpointImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		a := rs.Primary.Attributes
		return fmt.Sprintf(
			"id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			a["id"], a["cluster_id"], a["project_id"], a["organization_id"],
		), nil
	}
}

// testAccCheckOutboundPrivateEndpointDestroy checks that every endpoint the test created is gone.
// Destroy waits for the teardown to finish, so the API must answer not found.
func testAccCheckOutboundPrivateEndpointDestroy(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		data := newTestClient(t)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != outboundPrivateEndpointType {
				continue
			}
			a := rs.Primary.Attributes
			cfg := api.EndpointCfg{
				Url: fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints/%s",
					data.HostURL, a["organization_id"], a["project_id"], a["cluster_id"], a["id"]),
				Method:        http.MethodGet,
				SuccessStatus: http.StatusOK,
			}
			_, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
			if err == nil {
				return fmt.Errorf("outbound private endpoint %s still exists after destroy", a["id"])
			}
			if notFound, msg := api.CheckResourceNotFoundError(err); !notFound {
				return fmt.Errorf("failed to read outbound private endpoint %s after destroy: %s", a["id"], msg)
			}
		}
		return nil
	}
}

func testAccOutboundPrivateEndpointResourceConfig(resourceName string, env outboundPrivateEndpointEnv, name, hostname string) string {
	hostnameLine := ""
	if hostname != "" {
		hostnameLine = fmt.Sprintf("  hostname        = %q\n", hostname)
	}
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_outbound_private_endpoint" "%[2]s" {
  organization_id = %[3]q
  project_id      = %[4]q
  cluster_id      = %[5]q
  name            = %[6]q
  service_name    = %[7]q
%[8]s}
`, globalProviderBlock, resourceName, globalOrgId, env.projectId, env.clusterId, name, env.serviceName, hostnameLine)
}
