package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/provider"
)

// appServiceAuditLogStreamingPlanConfig returns a configuration whose provider points at an unreachable host, so a
// plan only passes or fails on validation.
func appServiceAuditLogStreamingPlanConfig(outputType, credentials string) string {
	const id = "ffffffff-aaaa-1414-eeee-000000000000"
	return fmt.Sprintf(`
provider "couchbase-capella" {
  host                 = "http://127.0.0.1:1"
  authentication_token = "test-token"
}

resource "couchbase-capella_app_service_audit_log_streaming" "streaming" {
  organization_id   = %[1]q
  project_id        = %[1]q
  cluster_id        = %[1]q
  app_service_id    = %[1]q
  streaming_enabled = true
  output_type       = %[2]q
  credentials = {
%[3]s
  }
}
`, id, outputType, credentials)
}

const (
	auditStreamingGenericHttpCredentials = `    generic_http = { url = "https://example.com/audit" }`
	auditStreamingDatadogCredentials     = `    datadog = { url = "https://http-intake.logs.datadoghq.com", api_key = "test-key" }`
)

func capellaProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(provider.New()()),
	}
}

// TestAppServiceAuditLogStreamingCredentialsMatchOutputType checks that credentials for a different output_type are
// rejected at plan time with the same error as couchbase-capella_app_service_log_streaming. Today the resource has no
// ValidateConfig, so the plan passes and the mismatch only fails at apply (AV-153198).
func TestAppServiceAuditLogStreamingCredentialsMatchOutputType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: capellaProviderFactories(),
		Steps: []resource.TestStep{{
			Config:      appServiceAuditLogStreamingPlanConfig("datadog", auditStreamingGenericHttpCredentials),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`Missing Credential Configuration`),
		}},
	})
}

// TestAppServiceAuditLogStreamingOneCredentialsBlock checks that more than one credentials block is rejected at plan
// time. The shared credentials schema already enforces this.
func TestAppServiceAuditLogStreamingOneCredentialsBlock(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: capellaProviderFactories(),
		Steps: []resource.TestStep{{
			Config: appServiceAuditLogStreamingPlanConfig("generic_http",
				auditStreamingGenericHttpCredentials+"\n"+auditStreamingDatadogCredentials),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`Exactly one of \[datadog, dynatrace, elastic, generic_http, loki, splunk,\s+sumologic\] must be specified`),
		}},
	})
}
