package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppEndpointAuditLoggingMinimal sets only audit_enabled on the shared App Endpoint.
func TestAccAppEndpointAuditLoggingMinimal(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_ae_audit_")
	resourceReference := "couchbase-capella_app_endpoint_audit_logging." + name

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: appEndpointAuditLoggingPrerequisites(name) + appEndpointAuditLoggingConfig(name, name, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
				),
			},
		},
	})
}

// TestAccAppEndpointAuditLoggingAdoptKeepsUnconfigured checks that creating the resource for an App Endpoint that
// already has audit settings keeps the settings the configuration does not mention. The PUT replaces the whole
// configuration and Create leaves unconfigured attributes out of it, so they reset to their API defaults.
//
// Step 1 stands in for settings made outside Terraform: destroying this resource only removes it from state, so
// after step 2 the disabled user stays in Capella. Step 3 then adopts the endpoint with only audit_enabled.
// The prerequisites stay in every step, so turning audit logging off cannot clear the endpoint's settings.
func TestAccAppEndpointAuditLoggingAdoptKeepsUnconfigured(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_ae_audit_")
	adopted := "couchbase-capella_app_endpoint_audit_logging.adopt_" + name
	prerequisites := appEndpointAuditLoggingPrerequisites(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: prerequisites + appEndpointAuditLoggingConfig("outside_"+name, name, `
  disabled_users = [{ domain = "sgw", name = "tf_acc_ui_user" }]`),
				Check: resource.TestCheckResourceAttr(
					"couchbase-capella_app_endpoint_audit_logging.outside_"+name, "disabled_users.#", "1"),
			},
			{
				Config: prerequisites,
			},
			{
				Config: prerequisites + appEndpointAuditLoggingConfig("adopt_"+name, name, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(adopted, "disabled_users.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(adopted, "disabled_users.*", map[string]string{
						"domain": "sgw",
						"name":   "tf_acc_ui_user",
					}),
				),
			},
		},
	})
}

// appEndpointAuditLoggingPrerequisites enables App Service audit logging and brings the shared App Endpoint online.
// The shared App Endpoint is created Offline, and an audit logging PUT against it hangs until the request times out.
func appEndpointAuditLoggingPrerequisites(name string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service_audit_logging" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}

resource "couchbase-capella_app_endpoint_activation_status" "%[2]s" {
  organization_id   = "%[3]s"
  project_id        = "%[4]s"
  cluster_id        = "%[5]s"
  app_service_id    = "%[6]s"
  app_endpoint_name = "%[7]s"
  state             = "Online"
}
`, globalProviderBlock, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName)
}

// appEndpointAuditLoggingConfig returns the App Endpoint audit logging resource, depending on the prerequisites
// named prerequisitesName, with extra attribute lines in settings.
func appEndpointAuditLoggingConfig(resourceName, prerequisitesName, settings string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_app_endpoint_audit_logging" "%[1]s" {
  organization_id   = "%[2]s"
  project_id        = "%[3]s"
  cluster_id        = "%[4]s"
  app_service_id    = "%[5]s"
  app_endpoint_name = "%[6]s"
  audit_enabled     = true
%[7]s
  depends_on = [
    couchbase-capella_app_service_audit_logging.%[8]s,
    couchbase-capella_app_endpoint_activation_status.%[8]s,
  ]
}
`, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName,
		settings, prerequisitesName)
}
