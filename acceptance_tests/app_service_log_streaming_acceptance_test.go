package acceptance_tests

import (
	"context"
	"fmt"
	re "regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// TestAccAppServiceLogStreaming tests the full lifecycle of the app_service_log_streaming resource:
// Create (generic_http) -> Update (credential change) -> ImportState -> Delete (automatic).
func TestAccAppServiceLogStreaming(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_log_streaming_")
	resourceReference := "couchbase-capella_app_service_log_streaming." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceLogStreamingDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccAppServiceLogStreamingConfig(
					resourceName,
					"https://example.com/logs",
					"test_user",
					"test_password",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "config_state", string(api.GetLogStreamingResponseConfigStateEnabled)),
					resource.TestCheckResourceAttrSet(resourceReference, "streaming_state"),
				),
			},
			// Update credentials (in-place, same output_type)
			{
				Config: testAccAppServiceLogStreamingConfig(
					resourceName,
					"https://example.com/logs/updated",
					"updated_user",
					"updated_password",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "config_state", string(api.GetLogStreamingResponseConfigStateEnabled)),
				),
			},
			// ImportState Testing
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppServiceLogStreamingImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "app_service_id",
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"credentials"},
			},
		},
	})
}

// TestAccAppServiceLogStreamingMismatchedCredentials tests that a validation error is returned
// when the output_type does not match the provided credentials block.
func TestAccAppServiceLogStreamingMismatchedCredentials(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_log_streaming_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceLogStreamingMismatchedConfig(resourceName),
				ExpectError: re.MustCompile("credentials.datadog must be configured when output_type is"),
			},
		},
	})
}

// TestAccAppServiceLogStreamingMissingCredentials tests that a validation error is returned
// when the correct credentials block for the output_type is not provided.
func TestAccAppServiceLogStreamingMissingCredentials(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_log_streaming_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceLogStreamingMissingConfig(resourceName),
				ExpectError: re.MustCompile("credentials.generic_http must be configured when output_type is"),
			},
		},
	})
}

// testAccCheckAppServiceLogStreamingDestroy verifies that after Terraform destroys the log streaming resource, the
// remote config_state has transitioned to "disabled". This is because destroying log streaming does not actually
// delete a resource, but instead disables log streaming on the app service.
func testAccCheckAppServiceLogStreamingDestroy(s *terraform.State) error {
	data := newTestClient()

	orgUUID, err := uuid.Parse(globalOrgId)
	if err != nil {
		return fmt.Errorf("failed to parse organization_id: %w", err)
	}
	projUUID, err := uuid.Parse(globalProjectId)
	if err != nil {
		return fmt.Errorf("failed to parse project_id: %w", err)
	}
	clusterUUID, err := uuid.Parse(globalClusterId)
	if err != nil {
		return fmt.Errorf("failed to parse cluster_id: %w", err)
	}
	appServiceUUID, err := uuid.Parse(globalAppServiceId)
	if err != nil {
		return fmt.Errorf("failed to parse app_service_id: %w", err)
	}

	response, err := data.ClientV2.GetAppServiceLogStreamingWithResponse(
		context.Background(),
		orgUUID,
		projUUID,
		clusterUUID,
		appServiceUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to get log streaming state after destroy: %w", err)
	}

	if response.JSON200 == nil {
		return fmt.Errorf("expected JSON200 response body but got nil, status code: %d", response.StatusCode())
	}

	configState := response.JSON200.ConfigState
	if configState == nil || *configState != api.GetLogStreamingResponseConfigStateDisabled {
		var actual string
		if configState != nil {
			actual = string(*configState)
		}
		return fmt.Errorf("expected config_state to be %q after destroy, got %q",
			api.GetLogStreamingResponseConfigStateDisabled, actual)
	}

	return nil
}

// testAccAppServiceLogStreamingConfig returns the HCL config for a generic_http log streaming resource.
func testAccAppServiceLogStreamingConfig(resourceName, url, user, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_log_streaming" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  output_type     = "generic_http"

  credentials = {
    generic_http = {
      url      = "%[7]s"
      user     = "%[8]s"
      password = "%[9]s"
    }
  }
}
`,
		globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		resourceName,
		url,
		user,
		password,
	)
}

// testAccAppServiceLogStreamingMismatchedConfig returns an HCL config where
// output_type is "datadog" but generic_http credentials are provided.
func testAccAppServiceLogStreamingMismatchedConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_log_streaming" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  output_type     = "datadog"

  credentials = {
    generic_http = {
      url      = "https://example.com/logs"
      user     = "test_user"
      password = "test_password"
    }
  }
}
`,
		globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		resourceName,
	)
}

// testAccAppServiceLogStreamingMissingConfig returns an HCL config where
// output_type is "generic_http" but datadog credentials are provided instead.
func testAccAppServiceLogStreamingMissingConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_log_streaming" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  output_type     = "generic_http"

  credentials = {}
}
`,
		globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		resourceName,
	)
}

// generateAppServiceLogStreamingImportId generates the import ID string
// from the Terraform state for the app_service_log_streaming resource.
func generateAppServiceLogStreamingImportId(resourceReference string) resource.ImportStateIdFunc {
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
