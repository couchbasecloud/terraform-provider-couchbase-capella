package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/data_api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAccDataApiResource exercises the resource and its paired data source together through the
// enable, enable-peering, import and disable transitions on the shared cluster.
func TestAccDataApiResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_")
	resourceReference := "couchbase-capella_data_api." + resourceName
	dsName := randomStringWithPrefix("tf_acc_data_api_ds_")
	dsReference := "data.couchbase-capella_data_api." + dsName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApiResourceAndDatasourceConfig(resourceName, dsName, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttrSet(resourceReference, "connection_string"),
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(dsReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(dsReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(dsReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttrSet(dsReference, "connection_string"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    generateDataApiImportId(resourceReference),
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
			{
				Config: testAccDataApiResourceAndDatasourceConfig(resourceName, dsName, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "true"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "enabled"),
					resource.TestCheckResourceAttrSet(resourceReference, "connection_string"),
					resource.TestCheckResourceAttr(dsReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(dsReference, "enable_network_peering", "true"),
					resource.TestCheckResourceAttr(dsReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(dsReference, "state_for_network_peering", "enabled"),
					resource.TestCheckResourceAttrSet(dsReference, "connection_string"),
				),
			},
			{
				Config: testAccDataApiResourceAndDatasourceConfig(resourceName, dsName, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "false"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "disabled"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttr(dsReference, "enable_data_api", "false"),
					resource.TestCheckResourceAttr(dsReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(dsReference, "state_for_data_api", "disabled"),
					resource.TestCheckResourceAttr(dsReference, "state_for_network_peering", "disabled"),
				),
			},
		},
	})
}

// TestAccDataApiResourceNetworkPeeringWithoutDataApi verifies that enabling network peering while the Data API is
// disabled is rejected by local config validation instead of being sent to the API.
func TestAccDataApiResourceNetworkPeeringWithoutDataApi(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_peering_only_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataApiResourceConfig(resourceName, false, true),
				ExpectError: regexp.MustCompile(`Network peering cannot be enabled`),
			},
		},
	})
}

// TestAccDataApiResourceInvalidUUIDs verifies that each ID attribute rejects
// non-UUID values via local schema validation, one attribute per subtest.
func TestAccDataApiResourceInvalidUUIDs(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		projectID      string
		clusterID      string
	}{
		{
			name:           "organization_id",
			organizationID: "not-a-uuid",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "project_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "not-a-uuid",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "cluster_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "not-a-uuid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_data_api_non_uuid_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccDataApiResourceIDsConfig(
							resourceName, test.organizationID, test.projectID, test.clusterID),
						ExpectError: regexp.MustCompile(
							`(?s)Invalid Attribute Value Match.*` + test.name + `.*must be a valid UUID`),
					},
				},
			})
		})
	}
}

func testAccDataApiResourceConfig(resourceName string, enableDataApi, enableNetworkPeering bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id        = "%[3]s"
  project_id             = "%[4]s"
  cluster_id             = "%[5]s"
  enable_data_api        = %[6]t
  enable_network_peering = %[7]t
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, enableDataApi, enableNetworkPeering)
}

func testAccDataApiResourceIDsConfig(resourceName, organizationID, projectID, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  enable_data_api = true
}
`, globalProviderBlock, resourceName, organizationID, projectID, clusterID)
}

func testAccDataApiResourceAndDatasourceConfig(resourceName, dsName string, enableDataApi, enableNetworkPeering bool) string {
	// enable_network_peering is omitted when false so those steps also exercise the attribute's default.
	enableNetworkPeeringAttr := ""
	if enableNetworkPeering {
		enableNetworkPeeringAttr = "\n  enable_network_peering = true"
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  enable_data_api = %[7]t%[8]s
}

data "couchbase-capella_data_api" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"

  depends_on = [couchbase-capella_data_api.%[2]s]
}
`, globalProviderBlock, resourceName, dsName, globalOrgId, globalProjectId, globalClusterId, enableDataApi, enableNetworkPeeringAttr)
}

func generateDataApiImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}

func retrieveDataApiStatusFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*data_api.GetDataApiStatusResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		data.HostURL, organizationId, projectId, clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return nil, err
	}
	status := &data_api.GetDataApiStatusResponse{}
	if err := json.Unmarshal(response.Body, status); err != nil {
		return nil, err
	}
	return status, nil
}

func testAccExistsDataApiResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		data := newTestClient(t)
		_, err := retrieveDataApiStatusFromServer(
			data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"],
		)
		return err
	}
}
