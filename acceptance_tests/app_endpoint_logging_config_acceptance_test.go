package acceptance_tests

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAccAppEndpointLoggingConfigResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_log_streaming_config_")
	resourceReference := "couchbase-capella_app_endpoint_log_streaming_config." + resourceName

	logKeys := "[" + strings.Join([]string{"\"HTTP\"", "\"Import\""}, ",") + "]"
	updatedLogKeys := "[" + "\"Auth\"" + "]"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointLoggingConfigResourceConfig(resourceName, "info", logKeys),
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
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppEndpointLoggingConfigImportIdForResource(resourceReference),
				ImportState:       true,
			},

			{
				Config: testAccAppEndpointLoggingConfigResourceConfig(resourceName, "warn", updatedLogKeys),
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
		},
	})
}

func TestAccAppEndpointLoggingConfigResourceInvalidLogLevel(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_log_streaming_config_")

	logKeys := "[" + strings.Join([]string{"\"HTTP\"", "\"Import\""}, ",") + "]"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppEndpointLoggingConfigResourceConfig(resourceName, "test", logKeys),
				ExpectError: regexp.MustCompile("Error executing upsert app endpoint logging config"),
			},
		},
	})
}

func TestAccAppEndpointLoggingConfigResourceInvalidLogKeys(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_log_streaming_config_")

	logKeys := "[" + strings.Join([]string{"\"Test\"", "\"Test2\""}, ",") + "]"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppEndpointLoggingConfigResourceConfig(resourceName, "info", logKeys),
				ExpectError: regexp.MustCompile("Error executing upsert app endpoint logging config"),
			},
		},
	})
}

func testAccAppEndpointLoggingConfigResourceConfig(resourceName, logLevel, logKeys string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_endpoint_log_streaming_config" "%[2]s" {
  		organization_id = "%[3]s"
  		project_id = "%[4]s"
    	cluster_id = "%[5]s"
    	app_service_id = "%[6]s"
  		app_endpoint_name = "%[7]s"

  		log_level = "%[8]s"
  		log_keys = %[9]s
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName, logLevel, logKeys)
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

		data, err := newTestClientV2()
		if err != nil {
			return err
		}

		err = retrieveAppEndpointLoggingConfigFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["app_service_id"], rawState["app_endpoint_name"])
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
		return errors.ErrUnexpectedStatusGettingAppEndpointLoggingConfig
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
