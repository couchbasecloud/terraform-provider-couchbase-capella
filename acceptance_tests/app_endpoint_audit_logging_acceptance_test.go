package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
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

// TestAccAppEndpointAuditLoggingClearEventIDs sets two filterable event IDs, then clears them with an empty set, as
// the resource documentation says to. Capella drops an empty event list before it calls Sync Gateway, so if Sync
// Gateway keeps the previous events or applies its defaults, the read-back is not empty and step 2 fails with
// "Provider produced inconsistent result after apply".
func TestAccAppEndpointAuditLoggingClearEventIDs(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_ae_audit_")
	resourceReference := "couchbase-capella_app_endpoint_audit_logging." + name
	prerequisites := appEndpointAuditLoggingPrerequisites(name)

	eventIDs := fmt.Sprintf(`
data "couchbase-capella_app_endpoint_audit_log_event_ids" "%[1]s" {
  organization_id   = "%[2]s"
  project_id        = "%[3]s"
  cluster_id        = "%[4]s"
  app_service_id    = "%[5]s"
  app_endpoint_name = "%[6]s"
  depends_on        = [couchbase-capella_app_endpoint_activation_status.%[1]s]
}

locals {
  filterable_event_ids = [
    for id in slice(sort([
      for e in data.couchbase-capella_app_endpoint_audit_log_event_ids.%[1]s.data : format("%%010d", e.id) if e.filterable
    ]), 0, 2) : tonumber(id)
  ]
}
`, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: prerequisites,
			},
			{
				PreConfig: func() { waitForAppEndpointAuditLogEvents(t) },
				Config: prerequisites + eventIDs + appEndpointAuditLoggingConfig(name, name, `
  enabled_event_ids = local.filterable_event_ids`),
				Check: resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "2"),
			},
			{
				Config: prerequisites + appEndpointAuditLoggingConfig(name, name, `
  enabled_event_ids = []`),
				Check: resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "0"),
			},
		},
	})
}

// waitForAppEndpointAuditLogEvents polls the shared App Endpoint's audit log events until the API answers. Right
// after the endpoint is brought online the API can return 500 for a while, which is not what these tests check.
func waitForAppEndpointAuditLogEvents(t *testing.T) {
	t.Helper()
	data := newTestClient(t)
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/auditLogEvents",
		globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	deadline := time.Now().Add(10 * time.Minute)
	for {
		_, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
		if err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("App Endpoint audit log events still failing after 10 minutes: %v", err)
		}
		t.Logf("App Endpoint audit log events not ready, retrying: %v", err)
		time.Sleep(15 * time.Second)
	}
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
