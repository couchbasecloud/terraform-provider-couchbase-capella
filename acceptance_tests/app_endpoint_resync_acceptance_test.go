package acceptance_tests

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAccAppEndpointResyncResource provides the steps to test the full lifecycle
// of the app_endpoint_resync_job resource: Create -> ImportState
func TestAccAppEndpointResyncResource(t *testing.T) {

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_resync_job_")
	resourceReference := "couchbase-capella_app_endpoint_resync_job." + resourceName

	scopes := "{\n_default = [\"_default\"]\n}"
	invalidScopes := "{\ntest = [\"test\"]\n}"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{

			// tests that an error is returned if the scopes are invalid
			{
				Config:      testAccAppEndpointResyncConfig(resourceName, invalidScopes),
				ExpectError: regexp.MustCompile("Unexpected status while starting App Endpoint Resync"),
			},

			// Create and Read testing
			{
				Config: testAccAppEndpointResyncConfig(resourceName, scopes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAppEndpointResyncResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "scopes._default.0", "_default"),
					resource.TestCheckResourceAttr(resourceReference, "collections_processing._default.0", "_default"),
					resource.TestCheckResourceAttr(resourceReference, "last_error", ""),
					resource.TestCheckResourceAttrSet(resourceReference, "docs_changed"),
					resource.TestCheckResourceAttrSet(resourceReference, "docs_processed"),
					resource.TestCheckResourceAttrSet(resourceReference, "start_time"),
					resource.TestCheckResourceAttrSet(resourceReference, "state"),
				),
			},

			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppEndpointResyncImportIdForResource(resourceReference),
				ImportState:       true,
			},
		},
	})

}

// testAccAppEndpointResyncConfig returns the HCL config for an app endpoint resync resource
func testAccAppEndpointResyncConfig(resourceName, scopes string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_endpoint_resync_job" "%[2]s" {
		organization_id = "%[3]s"
		project_id = "%[4]s"
		cluster_id = "%[5]s"
		app_service_id = "%[6]s"
		app_endpoint_name = "%[7]s"

		scopes = %[8]s
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName, scopes)
}

func testAccExistsAppEndpointResyncResource(resourceReference string) resource.TestCheckFunc {
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

		err := retrieveAppEndpointResyncFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["app_service_id"], rawState["app_endpoint_name"])
		if err != nil {
			return err
		}
		return nil
	}
}

func retrieveAppEndpointResyncFromServer(data *providerschema.Data, organizationId, projectId, clusterId, appServiceId, appEndpointName string) error {

	ctx := context.Background()

	organizationUUID, _ := uuid.Parse(organizationId)
	projectUUID, _ := uuid.Parse(projectId)
	clusterUUID, _ := uuid.Parse(clusterId)
	appServiceUUID, _ := uuid.Parse(appServiceId)

	getResyncStatusResp, err := data.ClientV2.GetAppEndpointResyncWithResponse(ctx, organizationUUID, projectUUID, clusterUUID, appServiceUUID, appEndpointName)
	if err != nil {
		return err
	}

	if getResyncStatusResp.JSON200 == nil {
		return errors.New("Unexpected status while getting App Endpoint Resync: " + string(getResyncStatusResp.Body))
	}

	return nil
}

func generateAppEndpointResyncImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
