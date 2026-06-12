package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccEndpointActivationStatus uses sequential subtests to ensure that resync tests
// do not occur while the app endpoint is online.
func TestAccAppEndpointActivationStatus(t *testing.T) {
	ensureActivationEndpoint(t)

	// Allow this test to run in parallel with other top-level tests, but ensure that the subtests run sequentially
	// This is normally set by resource.ParallelTest
	t.Parallel()

	t.Run("App Endpoint Activation Status", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    testAccAppEndpointActivationStatus(),
		})
	})

	t.Run("App Endpoint Resync", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    testAccAppEndpointResync(t),
		})
	})
}

func testAccAppEndpointActivationStatus() []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_activation_")
	resourceReference := "couchbase-capella_app_endpoint_activation_status." + resourceName

	// Use a stable endpoint name so we can import by name.
	endpointName := appEndpointActivationEndpointName

	return []resource.TestStep{
		{
			// Create endpoint and set state Online
			Config: testAccAppEndpointActivationStatusConfig(resourceName, endpointName, "Online"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", appEndpointClusterId),
				resource.TestCheckResourceAttr(resourceReference, "app_service_id", appEndpointAppServiceId),
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
	}
}

func TestAccAppEndpointActivationStatusInvalidState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_activation_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppEndpointActivationStatusInvalidStateConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*Archived.*Online.*Offline`),
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
		appEndpointClusterId,
		appEndpointAppServiceId,
		endpointName,
		desiredState,
		resourceName,
	)
}

func testAccAppEndpointActivationStatusInvalidStateConfig(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_endpoint_activation_status" "%[2]s" {
	  organization_id   = "00000000-0000-0000-0000-000000000000"
	  project_id        = "11111111-1111-1111-1111-111111111111"
	  cluster_id        = "22222222-2222-2222-2222-222222222222"
	  app_service_id    = "33333333-3333-3333-3333-333333333333"
	  app_endpoint_name = "qe-endpoint"
	  state             = "Archived"
	}
	`,
		globalProviderBlock,
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
