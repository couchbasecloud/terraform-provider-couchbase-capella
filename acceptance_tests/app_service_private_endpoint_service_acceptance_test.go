package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

func TestAccAppServicePrivateEndpointServiceResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_pe_service_")
	resourceReference := "couchbase-capella_app_service_private_endpoint_service." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServicePrivateEndpointServiceDestroy(t),
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAppServicePrivateEndpointServiceResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "state", "enabled"),
				),
			},
			// ImportState
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppServicePrivateEndpointServiceImportIdForResource(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_service_id",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppServicePrivateEndpointServiceResourceConfig(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_service_private_endpoint_service" "%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

func generateAppServicePrivateEndpointServiceImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["app_service_id"],
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}

// testAccCheckAppServicePrivateEndpointServiceDestroy verifies that the remote state is "disabled" after destroy,
// because destroy disables the private endpoint service on the App Service and does not delete a resource.
func testAccCheckAppServicePrivateEndpointServiceDestroy(t *testing.T) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		data := newTestClient(t)

		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService",
			data.HostURL, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId,
		)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
		if err != nil {
			return fmt.Errorf("failed to get App Service private endpoint service state after destroy: %w", err)
		}

		var stateResp apigen.GetAppServicePrivateEndpointStateResponse
		if err := json.Unmarshal(response.Body, &stateResp); err != nil {
			return fmt.Errorf("failed to unmarshal App Service private endpoint service state: %w", err)
		}

		disabled := apigen.GetAppServicePrivateEndpointStateResponseStateDisabled
		if stateResp.State == nil || *stateResp.State != disabled {
			var actual string
			if stateResp.State != nil {
				actual = string(*stateResp.State)
			}
			return fmt.Errorf("expected state to be %q after destroy, got %q", disabled, actual)
		}

		return nil
	}
}
