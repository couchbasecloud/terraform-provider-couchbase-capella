package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// TestAccOutboundPrivateEndpointCapacityDataSource reads the shared test cluster's outbound private
// endpoint capacity. Other tests may create endpoints on the same cluster in parallel, so used is
// compared with a direct API read rather than pinned.
func TestAccOutboundPrivateEndpointCapacityDataSource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_ope_capacity_")
	resourceReference := "data.couchbase-capella_outbound_private_endpoint_capacity." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOutboundPrivateEndpointCapacityDataSourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestMatchResourceAttr(resourceReference, "cluster_endpoint_limit.used", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceReference, "cluster_endpoint_limit.limit", regexp.MustCompile(`^[1-9]\d*$`)),
					testAccCheckOutboundPrivateEndpointCapacityMatchAPI(resourceReference),
				),
			},
		},
	})
}

// testAccCheckOutboundPrivateEndpointCapacityMatchAPI reads the route the data source reads and
// compares the limit with it. The response is decoded by property name rather than through the
// provider's own struct, so a property name the provider gets wrong fails here. used is left out,
// since a parallel test can change it between the two reads.
func testAccCheckOutboundPrivateEndpointCapacityMatchAPI(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := api.EndpointCfg{
			Url: fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints/capacity",
				globalHost, globalOrgId, globalProjectId, globalClusterId),
			Method:        http.MethodGet,
			SuccessStatus: http.StatusOK,
		}
		response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
		if err != nil {
			return fmt.Errorf("reading the outbound private endpoint capacity: %w", err)
		}

		var want struct {
			ClusterEndpointLimit map[string]int64 `json:"clusterEndpointLimit"`
		}
		if err := json.Unmarshal(response.Body, &want); err != nil {
			return fmt.Errorf("decoding the outbound private endpoint capacity: %w", err)
		}
		limit, ok := want.ClusterEndpointLimit["limit"]
		if !ok {
			return fmt.Errorf("the capacity response has no clusterEndpointLimit.limit: %s", response.Body)
		}

		return resource.TestCheckResourceAttr(resourceReference, "cluster_endpoint_limit.limit", strconv.FormatInt(limit, 10))(s)
	}
}

func testAccOutboundPrivateEndpointCapacityDataSourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_outbound_private_endpoint_capacity" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}
