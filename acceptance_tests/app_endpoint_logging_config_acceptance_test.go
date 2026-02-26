package acceptance_tests

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// testAccAppEndpointLoggingConfigResource provides the steps to test the full lifecycle
// of the app_endpoint_log_streaming_config resource: Create -> Update -> ImportState,
// and that errors are returned if the log_level or log_keys are invalid
func testAccAppEndpointLoggingConfigResource() []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_log_streaming_config_")
	resourceReference := "couchbase-capella_app_endpoint_log_streaming_config." + resourceName

	logKeys := "[" + strings.Join([]string{"\"HTTP\"", "\"Import\""}, ",") + "]"
	updatedLogKeys := "[" + "\"Auth\"" + "]"
	invalidLogKeys := "[" + strings.Join([]string{"\"Test\"", "\"Test2\""}, ",") + "]"

	appServiceLogStreamingResourceName := randomStringWithPrefix("tf_app_service_log_streaming_config_")

	return []resource.TestStep{
		{
			Config: testAccAppEndpointLoggingConfigResourceConfig(appServiceLogStreamingResourceName, resourceName, "info", logKeys),
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccExistsAppEndpointLoggingConfigResource(resourceReference),
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
				resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
				resource.TestCheckResourceAttr(resourceReference, "log_level", "info"),
				resource.TestCheckResourceAttr(resourceReference, "log_keys.0", "HTTP"),
				resource.TestCheckResourceAttr(resourceReference, "log_keys.1", "Import"),
			),
		},
		{
			Config: testAccAppEndpointLoggingConfigResourceConfig(appServiceLogStreamingResourceName, resourceName, "warn", updatedLogKeys),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
				resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
				resource.TestCheckResourceAttr(resourceReference, "log_level", "warn"),
				resource.TestCheckResourceAttr(resourceReference, "log_keys.0", "Auth"),
			),
		},
		{
			ResourceName:      resourceReference,
			ImportStateIdFunc: generateAppEndpointLoggingConfigImportIdForResource(resourceReference),
			ImportState:       true,
		},

		// tests that an error is returned if the log_level is invalid
		{
			Config:      testAccAppEndpointLoggingConfigResourceConfig(appServiceLogStreamingResourceName, resourceName, "test", logKeys),
			ExpectError: regexp.MustCompile("Error executing upsert app endpoint logging config"),
		},

		// tests that an error is returned if any of the log_keys are invalid
		{
			Config:      testAccAppEndpointLoggingConfigResourceConfig(appServiceLogStreamingResourceName, resourceName, "info", invalidLogKeys),
			ExpectError: regexp.MustCompile("Error executing upsert app endpoint logging config"),
		},
	}
}

// testAccAppEndpointLoggingConfigResourceConfig returns the HCL config for a app endpoint log streaming config resource
// and an app service log streaming resource, as app service log streaming is required to be enabled
func testAccAppEndpointLoggingConfigResourceConfig(appServiceLogStreamingResourceName, resourceName, logLevel, logKeys string) string {
	appServiceLogStreamingResource := appServiceLogStreamingConfig(
		appServiceLogStreamingResourceName,
		"https://example.com/logs",
		"test_user",
		"test_password",
	)

	return fmt.Sprintf(`
	%[1]s

	%[2]s

	resource "couchbase-capella_app_endpoint_log_streaming_config" "%[3]s" {
		organization_id = "%[4]s"
		project_id = "%[5]s"
		cluster_id = "%[6]s"
		app_service_id = "%[7]s"
		app_endpoint_name = "%[8]s"

		log_level = "%[9]s"
		log_keys = %[10]s

		depends_on = [
			%[11]s
		]
	}
	`,
		globalProviderBlock,
		appServiceLogStreamingResource,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		globalAppEndpointName,
		logLevel,
		logKeys,
		"couchbase-capella_app_service_log_streaming."+appServiceLogStreamingResourceName,
	)
}

func testAccExistsAppEndpointLoggingConfigResource(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}

		data := newTestClient()

		err := retrieveAppEndpointLoggingConfigFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["app_service_id"], rawState["app_endpoint_name"])
		if err != nil {
			return err
		}
		return nil
	}
}

func retrieveAppEndpointLoggingConfigFromServer(data *providerschema.Data, organizationId, projectId, clusterId, appServiceId, appEndpointName string) error {

	ctx := context.Background()

	organizationUUID, _ := uuid.Parse(organizationId)
	projectUUID, _ := uuid.Parse(projectId)
	clusterUUID, _ := uuid.Parse(clusterId)
	appServiceUUID, _ := uuid.Parse(appServiceId)

	getLoggingConfigResp, err := data.ClientV2.GetAppEndpointLogStreamingConfigWithResponse(
		ctx,
		organizationUUID,
		projectUUID,
		clusterUUID,
		appServiceUUID,
		appEndpointName,
	)
	if err != nil {
		return err
	}

	if getLoggingConfigResp.JSON200 == nil {
		return errors.New("Unexpected status while getting App Endpoint Logging Config: " + string(getLoggingConfigResp.Body))
	}

	return nil
}

func generateAppEndpointLoggingConfigImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("app_endpoint_name=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["app_endpoint_name"], rawState["app_service_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}
