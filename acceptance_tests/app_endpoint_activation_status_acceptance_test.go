package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppEndpointActivationStatus(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_activation_")
	resourceReference := "couchbase-capella_app_endpoint_activation_status." + resourceName

	// Use a stable endpoint name so we can import by name
	endpointName := globalAppEndpointName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// Create endpoint and set state Online
				Config: testAccAppEndpointActivationStatusConfig(resourceName, endpointName, "Online"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", endpointName),
					resource.TestCheckResourceAttr(resourceReference, "state", "Online"),
				),
			},
			{
				// Update to Offline
				Config: testAccAppEndpointActivationStatusConfig(resourceName, endpointName, "Offline"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "state", "Offline"),
				),
			},
			{
				// Import by composite ID (uses endpoint name, not an ID)
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppEndpointActivationStatusImportId(resourceReference),
				ImportState:       true,
			},
		},
	})
}

func testAccAppEndpointActivationStatusConfig(resourceName, endpointName, desiredState string) string {
	// Create the underlying App Endpoint if it does not exist, then manage activation status
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_endpoint_activation_status" "%[8]s" {
	  organization_id   = "%[2]s"
	  project_id        = "%[3]s"
	  cluster_id        = "%[4]s"
	  app_service_id    = "%[5]s"
	  app_endpoint_name = "%[6]s"
	  state             = "%[7]s"
	}
	`,
		globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		endpointName,
		desiredState,
		resourceName,
	)
}

func generateAppEndpointActivationStatusImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		// Import format: organization_id=...,project_id=...,cluster_id=...,app_service_id=...,app_endpoint_name=...
		return fmt.Sprintf(
			"organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s,app_endpoint_name=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
			rawState["app_service_id"],
			rawState["app_endpoint_name"],
		), nil
	}
}
