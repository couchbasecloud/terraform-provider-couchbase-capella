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
)

// Audit logging and audit log streaming are one per App Service, so these tests run sequentially. Streaming stays
// disabled throughout, so no reachable collector is needed; enabling it is not covered.

const (
	auditStreamingGenericHttp = `
  output_type = "generic_http"
  credentials = {
    generic_http = {
      url = "https://example.com/audit"
    }
  }`
	auditStreamingDatadog = `
  output_type = "datadog"
  credentials = {
    datadog = {
      url     = "https://http-intake.logs.datadoghq.com"
      api_key = "tf-acc-test-not-a-real-key"
    }
  }`
)

// TestAccAppServiceAuditLogStreaming creates audit log streaming, reads it back through the state data source,
// imports it, excludes the shared App Endpoint and then includes it again, and checks that destroy removes the
// configuration from Capella.
func TestAccAppServiceAuditLogStreaming(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_as_audit_streaming_")
	resourceReference := "couchbase-capella_app_service_audit_log_streaming." + name
	dataSourceReference := "data.couchbase-capella_app_service_audit_log_streaming_state." + name
	disabledEndpoints := fmt.Sprintf(`
  disabled_app_endpoints = ["%s"]`, globalAppEndpointName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAuditLogStreamingRemoved(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAuditLogStreamingConfig(name, auditStreamingGenericHttp, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "credentials.generic_http.url", "https://example.com/audit"),
					resource.TestCheckNoResourceAttr(resourceReference, "disabled_app_endpoints.#"),
					resource.TestCheckResourceAttrSet(resourceReference, "log_streaming_state"),

					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(dataSourceReference, "streaming_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(dataSourceReference, "disabled_app_endpoints.#", "0"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "log_streaming_state", resourceReference, "log_streaming_state"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppServiceScopedImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_service_id",
				// The API does not return credentials.
				ImportStateVerifyIgnore: []string{"credentials"},
			},
			{
				Config: testAccAppServiceAuditLogStreamingConfig(name, auditStreamingGenericHttp+disabledEndpoints, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "credentials.generic_http.url", "https://example.com/audit"),
					resource.TestCheckResourceAttr(resourceReference, "disabled_app_endpoints.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "disabled_app_endpoints.*", globalAppEndpointName),
					resource.TestCheckResourceAttr(dataSourceReference, "disabled_app_endpoints.#", "1"),
					resource.TestCheckTypeSetElemAttr(dataSourceReference, "disabled_app_endpoints.*", globalAppEndpointName),
				),
			},
			{
				Config: testAccAppServiceAuditLogStreamingConfig(name, auditStreamingGenericHttp, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckNoResourceAttr(resourceReference, "disabled_app_endpoints.#"),
					resource.TestCheckResourceAttr(dataSourceReference, "disabled_app_endpoints.#", "0"),
				),
			},
		},
	})
}

// TestAccAppServiceAuditLogStreamingOutputTypeChange switches output_type twice while streaming stays disabled.
// After each PUT the provider waits only until streaming_enabled matches the plan, then copies the response into
// state. If the API still returned the previous output_type (or none, on create) at that point, the apply would fail
// with "Provider produced inconsistent result after apply" and the resource would be tainted.
func TestAccAppServiceAuditLogStreamingOutputTypeChange(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_as_audit_streaming_")
	resourceReference := "couchbase-capella_app_service_audit_log_streaming." + name

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAuditLogStreamingConfig(name, auditStreamingGenericHttp, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
				),
			},
			{
				Config: testAccAppServiceAuditLogStreamingConfig(name, auditStreamingDatadog, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "output_type", "datadog"),
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
				),
			},
			{
				Config: testAccAppServiceAuditLogStreamingConfig(name, auditStreamingGenericHttp, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
				),
			},
		},
	})
}

// testAccAppServiceAuditLogStreamingConfig returns disabled audit log streaming with the output lines in output, and
// with withDataSource the state data source reading it. The API rejects audit log streaming with 422 "Audit Logging
// Must Be Enabled To Utilize This Feature" unless App Service audit logging is enabled, so the config enables it;
// depends_on also destroys streaming first, because audit logging cannot be disabled while streaming is configured.
func testAccAppServiceAuditLogStreamingConfig(name, output string, withDataSource bool) string {
	config := fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service_audit_logging" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}

resource "couchbase-capella_app_service_audit_log_streaming" "%[2]s" {
  organization_id   = "%[3]s"
  project_id        = "%[4]s"
  cluster_id        = "%[5]s"
  app_service_id    = "%[6]s"
  streaming_enabled = false
%[7]s

  depends_on = [couchbase-capella_app_service_audit_logging.%[2]s]
}
`, globalProviderBlock, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, output)
	if !withDataSource {
		return config
	}
	return config + fmt.Sprintf(`
data "couchbase-capella_app_service_audit_log_streaming_state" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  depends_on      = [couchbase-capella_app_service_audit_log_streaming.%[1]s]
}
`, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

// testAccCheckAppServiceAuditLogStreamingRemoved checks that the App Service has no audit log streaming
// configuration: the API answers 404 or returns no output type.
func testAccCheckAppServiceAuditLogStreamingRemoved(t *testing.T) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		data := newTestClient(t)
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogStreaming",
			globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
		if err != nil {
			if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
				return nil
			}
			return fmt.Errorf("reading audit log streaming: %w", err)
		}
		var streaming struct {
			OutputType *string `json:"outputType"`
		}
		if err := json.Unmarshal(response.Body, &streaming); err != nil {
			return fmt.Errorf("unmarshalling audit log streaming: %w", err)
		}
		if streaming.OutputType != nil && *streaming.OutputType != "" {
			return fmt.Errorf("audit log streaming still configured with output type %q after destroy", *streaming.OutputType)
		}
		return nil
	}
}
