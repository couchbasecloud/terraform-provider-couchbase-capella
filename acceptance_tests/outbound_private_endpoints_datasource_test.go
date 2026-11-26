package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// TestAccOutboundPrivateEndpointsDataSource lists the outbound private endpoints of a cluster and
// compares them with the API's own response. The shared test cluster runs a Server version older
// than 8.5, which outbound private endpoints need, so the test reads its own cluster: set
// ACC_OUTBOUND_PE_PROJECT_ID and ACC_OUTBOUND_PE_CLUSTER_ID to an AWS cluster on Couchbase Server
// 8.5 or later with the SelfManagedPrivateEndpoints flag on. The cluster may have no endpoints.
func TestAccOutboundPrivateEndpointsDataSource(t *testing.T) {
	projectId := os.Getenv("ACC_OUTBOUND_PE_PROJECT_ID")
	clusterId := os.Getenv("ACC_OUTBOUND_PE_CLUSTER_ID")
	if projectId == "" || clusterId == "" {
		t.Skip("skipping: set ACC_OUTBOUND_PE_PROJECT_ID and ACC_OUTBOUND_PE_CLUSTER_ID (an AWS cluster on Couchbase Server " +
			"8.5 or later with the SelfManagedPrivateEndpoints flag on) to run this test")
	}

	resourceName := randomStringWithPrefix("tf_acc_outbound_pes_")
	resourceReference := "data.couchbase-capella_outbound_private_endpoints." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOutboundPrivateEndpointsDataSourceConfig(resourceName, projectId, clusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", projectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", clusterId),
					testAccCheckOutboundPrivateEndpointsMatchAPI(resourceReference, projectId, clusterId),
				),
			},
		},
	})
}

// testAccCheckOutboundPrivateEndpointsMatchAPI reads the route the data source reads and compares
// each endpoint with it. The response is decoded by property name rather than through the
// provider's own struct, so a property name the provider gets wrong fails here. Only the fields
// Terraform can't see move are compared: a service accepting a connection between the two reads
// changes the status, the addresses and the audit.
func testAccCheckOutboundPrivateEndpointsMatchAPI(resourceReference, projectId, clusterId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := api.EndpointCfg{
			Url: fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints",
				globalHost, globalOrgId, projectId, clusterId),
			Method:        http.MethodGet,
			SuccessStatus: http.StatusOK,
		}
		response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
		if err != nil {
			return fmt.Errorf("reading the outbound private endpoints: %w", err)
		}

		var want struct {
			Data []map[string]any `json:"data"`
		}
		if err := json.Unmarshal(response.Body, &want); err != nil {
			return fmt.Errorf("decoding the outbound private endpoints: %w", err)
		}

		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceReference, "data.#", strconv.Itoa(len(want.Data))),
		}
		for i, endpoint := range want.Data {
			prefix := fmt.Sprintf("data.%d.", i)
			for attr, property := range map[string]string{
				"id":            "id",
				"name":          "name",
				"endpoint_type": "endpointType",
				"service_name":  "serviceName",
			} {
				value, _ := endpoint[property].(string)
				checks = append(checks, resource.TestCheckResourceAttr(resourceReference, prefix+attr, value))
			}

			cloudProvider, _ := endpoint["cloudProvider"].(map[string]any)
			cloudType, _ := cloudProvider["type"].(string)
			region, _ := cloudProvider["region"].(string)
			isCustom, _ := endpoint["isCustomHostname"].(bool)
			checks = append(checks,
				resource.TestCheckResourceAttr(resourceReference, prefix+"organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, prefix+"project_id", projectId),
				resource.TestCheckResourceAttr(resourceReference, prefix+"cluster_id", clusterId),
				resource.TestCheckResourceAttr(resourceReference, prefix+"cloud_provider.type", cloudType),
				resource.TestCheckResourceAttr(resourceReference, prefix+"cloud_provider.region", region),
				resource.TestCheckResourceAttr(resourceReference, prefix+"is_custom_hostname", strconv.FormatBool(isCustom)),
				resource.TestCheckResourceAttrSet(resourceReference, prefix+"status"),
				resource.TestCheckResourceAttrSet(resourceReference, prefix+"audit.created_at"),
			)
		}

		return resource.ComposeAggregateTestCheckFunc(checks...)(s)
	}
}

// TestAccOutboundPrivateEndpointsDataSourceWithEndpoint creates an endpoint and checks the list
// returns it with every field the resource has. It lists the cluster before, with and after the
// endpoint, so the count goes back to where it started once the endpoint is gone: on a cluster
// with no other endpoints, the first and last reads are the empty list. Needs
// ACC_OUTBOUND_PE_SERVICE_NAME as well: an AWS endpoint service in the cluster's region that
// allows Capella's AWS account.
func TestAccOutboundPrivateEndpointsDataSourceWithEndpoint(t *testing.T) {
	projectId := os.Getenv("ACC_OUTBOUND_PE_PROJECT_ID")
	clusterId := os.Getenv("ACC_OUTBOUND_PE_CLUSTER_ID")
	serviceName := os.Getenv("ACC_OUTBOUND_PE_SERVICE_NAME")
	if projectId == "" || clusterId == "" || serviceName == "" {
		t.Skip("skipping: set ACC_OUTBOUND_PE_PROJECT_ID, ACC_OUTBOUND_PE_CLUSTER_ID and ACC_OUTBOUND_PE_SERVICE_NAME to run this test")
	}

	dataName := randomStringWithPrefix("tf_acc_outbound_pes_")
	dataReference := "data.couchbase-capella_outbound_private_endpoints." + dataName
	resourceName := randomStringWithPrefix("tf_acc_outbound_pe_")
	resourceReference := "couchbase-capella_outbound_private_endpoint." + resourceName
	endpointName := randomStringWithPrefix("tf-acc-outbound-pe-")

	var before int
	var createdId string
	listOnly := testAccOutboundPrivateEndpointsDataSourceConfig(dataName, projectId, clusterId)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: listOnly,
				Check:  testAccOutboundPrivateEndpointsCount(dataReference, &before),
			},
			{
				Config: globalProviderBlock + fmt.Sprintf(`
resource "couchbase-capella_outbound_private_endpoint" "%[1]s" {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[5]q
  service_name    = %[6]q
}

data "couchbase-capella_outbound_private_endpoints" "%[7]s" {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  depends_on      = [couchbase-capella_outbound_private_endpoint.%[1]s]
}
`, resourceName, globalOrgId, projectId, clusterId, endpointName, serviceName, dataName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccOutboundPrivateEndpointsCountFrom(dataReference, &before, 1),
					testAccOutboundPrivateEndpointsListsResource(dataReference, resourceReference, &createdId),
				),
			},
			// Removing the resource destroys the endpoint, but the list is read while planning,
			// before the destroy, so it is checked on the next step.
			{
				Config: listOnly,
			},
			{
				Config: listOnly,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccOutboundPrivateEndpointsCountFrom(dataReference, &before, 0),
					testAccOutboundPrivateEndpointsOmits(dataReference, &createdId),
				),
			},
		},
	})
}

// testAccOutboundPrivateEndpointsCount records how many endpoints the list has.
func testAccOutboundPrivateEndpointsCount(dataReference string, count *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataReference]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataReference)
		}
		n, err := strconv.Atoi(rs.Primary.Attributes["data.#"])
		if err != nil {
			return fmt.Errorf("reading data.# of %s: %w", dataReference, err)
		}
		*count = n
		return nil
	}
}

// testAccOutboundPrivateEndpointsCountFrom checks the list has added more endpoints than before.
func testAccOutboundPrivateEndpointsCountFrom(dataReference string, before *int, added int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(dataReference, "data.#", strconv.Itoa(*before+added))(s)
	}
}

// outboundPrivateEndpointMovingAttributes change on their own when the service accepts the
// connection, so they're only compared when the resource and the list saw the same status.
var outboundPrivateEndpointMovingAttributes = []string{
	"status", "provider_status", "ip_addresses", "provider_config", "private_dns_enabled", "audit.modified_",
	"audit.version",
}

// testAccOutboundPrivateEndpointsListsResource finds the resource's endpoint in the list and
// compares every attribute the resource has with the list entry, including nested ones and
// attributes that are null on both.
func testAccOutboundPrivateEndpointsListsResource(dataReference, resourceReference string, createdId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataReference]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataReference)
		}
		rs, ok := s.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		*createdId = rs.Primary.ID

		n, _ := strconv.Atoi(ds.Primary.Attributes["data.#"])
		prefix := ""
		for i := range n {
			if ds.Primary.Attributes[fmt.Sprintf("data.%d.id", i)] == rs.Primary.ID {
				prefix = fmt.Sprintf("data.%d.", i)
				break
			}
		}
		if prefix == "" {
			return fmt.Errorf("endpoint %s is not in the list", rs.Primary.ID)
		}

		sameStatus := ds.Primary.Attributes[prefix+"status"] == rs.Primary.Attributes["status"]
		listed := map[string]bool{}
		for attr := range ds.Primary.Attributes {
			if strings.HasPrefix(attr, prefix) {
				listed[strings.TrimPrefix(attr, prefix)] = true
			}
		}
		var errs []string
		for attr, want := range rs.Primary.Attributes {
			if attr == "%" || (!sameStatus && outboundPrivateEndpointMoves(attr)) {
				continue
			}
			delete(listed, attr)
			if got := ds.Primary.Attributes[prefix+attr]; got != want {
				errs = append(errs, fmt.Sprintf("%s: resource %q, list %q", attr, want, got))
			}
		}
		for attr := range listed {
			if attr != "%" && (sameStatus || !outboundPrivateEndpointMoves(attr)) {
				errs = append(errs, fmt.Sprintf("%s: only in the list", attr))
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("list entry for %s differs from the resource:\n%s", rs.Primary.ID, strings.Join(errs, "\n"))
		}
		return nil
	}
}

func outboundPrivateEndpointMoves(attr string) bool {
	for _, p := range outboundPrivateEndpointMovingAttributes {
		if strings.HasPrefix(attr, p) {
			return true
		}
	}
	return false
}

// testAccOutboundPrivateEndpointsOmits checks a destroyed endpoint is no longer listed.
func testAccOutboundPrivateEndpointsOmits(dataReference string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataReference]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataReference)
		}
		for attr, v := range ds.Primary.Attributes {
			if strings.HasSuffix(attr, ".id") && v == *id {
				return fmt.Errorf("destroyed endpoint %s is still listed", *id)
			}
		}
		return nil
	}
}

func testAccOutboundPrivateEndpointsDataSourceConfig(resourceName, projectId, clusterId string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_outbound_private_endpoints" "%[2]s" {
  organization_id = %[3]q
  project_id      = %[4]q
  cluster_id      = %[5]q
}
`, globalProviderBlock, resourceName, globalOrgId, projectId, clusterId)
}
