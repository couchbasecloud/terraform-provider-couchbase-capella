package acceptance_tests

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// appServiceLogStreamingActivationStatusSteps provides the steps to test the full lifecycle
// of the app_service_log_streaming_activation_status resource:
// Create (enabled) -> Update (paused) -> ImportState.
func appServiceLogStreamingActivationStatusSteps() []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_log_streaming_activation_")
	resourceReference := "couchbase-capella_app_service_log_streaming_activation_status." + resourceName

	// The activation status resource requires log streaming to be enabled first,
	// so we set up a log streaming resource as a dependency.
	logStreamingResourceName := randomStringWithPrefix("tf_acc_log_streaming_")

	return []resource.TestStep{
		// Create with state="enabled" (log streaming starts as enabled after setup, so this is a no-op API call
		// but validates the create path).
		{
			Config: testAccAppServiceLogStreamingActivationStatusConfig(logStreamingResourceName, resourceName, string(api.GetLogStreamingResponseConfigStateEnabled)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
				resource.TestCheckResourceAttr(resourceReference, "state", string(api.GetLogStreamingResponseConfigStateEnabled)),
			),
		},
		// Update to state="paused"
		{
			Config: testAccAppServiceLogStreamingActivationStatusConfig(logStreamingResourceName, resourceName, string(api.GetLogStreamingResponseConfigStatePaused)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "state", string(api.GetLogStreamingResponseConfigStatePaused)),
			),
		},
		// ImportState
		{
			ResourceName:                         resourceReference,
			ImportStateIdFunc:                    generateAppServiceLogStreamingActivationStatusImportId(resourceReference),
			ImportState:                          true,
			ImportStateVerifyIdentifierAttribute: "app_service_id",
			ImportStateVerify:                    true,
		},
	}
}

// testAccAppServiceLogStreamingActivationStatusConfig returns the HCL config for testing the
// app_service_log_streaming_activation_status resource. It includes a log streaming resource
// as a dependency since log streaming must be enabled for activation status management to work.
func testAccAppServiceLogStreamingActivationStatusConfig(logStreamingResourceName, resourceName, state string) string {
	logStreamingResource := appServiceLogStreamingConfig(
		logStreamingResourceName,
		"https://example.com/logs",
		"test_user",
		"test_password",
	)

	return fmt.Sprintf(`
	%[1]s

	%[2]s

	resource "couchbase-capella_app_service_log_streaming_activation_status" "%[3]s" {
	  organization_id = "%[4]s"
	  project_id      = "%[5]s"
	  cluster_id      = "%[6]s"
	  app_service_id  = "%[7]s"
	  state           = "%[8]s"

	  depends_on = [couchbase-capella_app_service_log_streaming.%[9]s]
	}
	`,
		globalProviderBlock,
		logStreamingResource,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		state,
		logStreamingResourceName,
	)
}

// generateAppServiceLogStreamingActivationStatusImportId generates the import ID string
// from the Terraform state for the app_service_log_streaming_activation_status resource.
func generateAppServiceLogStreamingActivationStatusImportId(resourceReference string) resource.ImportStateIdFunc {
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
