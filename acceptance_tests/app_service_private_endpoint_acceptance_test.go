package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// TestAccAppServicePrivateEndpointResource accepts an existing private endpoint
// connection request on the App Service, then imports it.
//
// Gated on ACC_APP_SERVICE_ENDPOINT_ID: the ID of a CSP private endpoint (for
// example an AWS interface VPC endpoint) that already requests a connection to
// globalAppServiceId, with the App Service private endpoint service enabled.
// It skips otherwise.
func TestAccAppServicePrivateEndpointResource(t *testing.T) {
	endpointId := os.Getenv("ACC_APP_SERVICE_ENDPOINT_ID")
	if endpointId == "" {
		t.Skip("skipping: set ACC_APP_SERVICE_ENDPOINT_ID to a private endpoint that requests a connection to the App Service")
	}

	resourceName := randomStringWithPrefix("tf_acc_app_service_pe_")
	resourceReference := "couchbase-capella_app_service_private_endpoint." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServicePrivateEndpointDestroy(t, endpointId),
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAppServicePrivateEndpointResourceConfig(resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, endpointId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "endpoint_id", endpointId),
					resource.TestCheckResourceAttr(resourceReference, "status", string(apigen.PrivateEndpointStatusLinked)),
				),
			},
			// ImportState
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppServicePrivateEndpointImportIdForResource(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "endpoint_id",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccAppServicePrivateEndpointInvalidEndpointID verifies that an empty
// endpoint_id is rejected at plan time, before any API call.
func TestAccAppServicePrivateEndpointInvalidEndpointID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_pe_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicePrivateEndpointResourceConfig(
					resourceName,
					"00000000-0000-0000-0000-000000000000",
					"11111111-1111-1111-1111-111111111111",
					"22222222-2222-2222-2222-222222222222",
					"33333333-3333-3333-3333-333333333333",
					"",
				),
				ExpectError: regexp.MustCompile(`(?s)endpoint_id.*string length must be at least 1`),
			},
		},
	})
}

// testAccCheckAppServicePrivateEndpointDestroy fails if the endpoint is still
// linked after destroy. An absent, rejected or pending endpoint passes.
func testAccCheckAppServicePrivateEndpointDestroy(t *testing.T, endpointId string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		data := newTestClient(t)
		cfg := api.EndpointCfg{
			Url: fmt.Sprintf(
				"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/endpoints",
				data.HostURL, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId,
			),
			Method:        http.MethodGet,
			SuccessStatus: http.StatusOK,
		}
		response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
		if err != nil {
			return fmt.Errorf("failed to list App Service private endpoints after destroy: %w", err)
		}

		var listResp apigen.GetPrivateEndpointsResponse
		if err := json.Unmarshal(response.Body, &listResp); err != nil {
			return err
		}
		for _, e := range listResp.Endpoints {
			if e.Id == endpointId && e.Status == apigen.PrivateEndpointStatusLinked {
				return fmt.Errorf("App Service private endpoint %s is still linked after destroy", endpointId)
			}
		}
		return nil
	}
}

func testAccAppServicePrivateEndpointResourceConfig(resourceName, orgId, projectId, clusterId, appServiceId, endpointId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_private_endpoint" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  endpoint_id     = "%[7]s"
}
`, globalProviderBlock, resourceName, orgId, projectId, clusterId, appServiceId, endpointId)
}

func generateAppServicePrivateEndpointImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		a := rs.Primary.Attributes
		return fmt.Sprintf(
			"endpoint_id=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			a["endpoint_id"],
			a["app_service_id"],
			a["cluster_id"],
			a["project_id"],
			a["organization_id"],
		), nil
	}
}
