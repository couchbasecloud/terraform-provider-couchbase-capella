package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppServiceAuditLogStreamingOutputTypeChange creates audit log streaming disabled, then switches output_type
// while it stays disabled, so no reachable collector is needed. After each PUT the provider waits only until
// streaming_enabled matches the plan, then copies the response into state. If the API still returns the previous
// output_type (or none, on create) at that point, the apply fails with "Provider produced inconsistent result after
// apply" and the resource is tainted. Audit logging and audit log streaming are one per App Service, so this runs
// sequentially.
func TestAccAppServiceAuditLogStreamingOutputTypeChange(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_as_audit_streaming_")
	resourceReference := "couchbase-capella_app_service_audit_log_streaming." + name

	genericHttp := `
  output_type = "generic_http"
  credentials = {
    generic_http = {
      url = "https://example.com/audit"
    }
  }`
	datadog := `
  output_type = "datadog"
  credentials = {
    datadog = {
      url     = "https://http-intake.logs.datadoghq.com"
      api_key = "tf-acc-test-not-a-real-key"
    }
  }`

	// The API rejects audit log streaming with 422 "Audit Logging Must Be Enabled To Utilize This Feature" unless
	// audit logging is enabled on the App Service. depends_on also destroys streaming before audit logging.
	config := func(output string) string {
		return fmt.Sprintf(`
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
}`, globalProviderBlock, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, output)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: config(genericHttp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
				),
			},
			{
				Config: config(datadog),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "output_type", "datadog"),
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
				),
			},
			{
				Config: config(genericHttp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "output_type", "generic_http"),
					resource.TestCheckResourceAttr(resourceReference, "streaming_enabled", "false"),
				),
			},
		},
	})
}
