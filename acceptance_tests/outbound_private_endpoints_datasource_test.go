package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
