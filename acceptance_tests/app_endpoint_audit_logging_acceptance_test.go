package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// App Endpoint audit logging is one configuration per App Endpoint, and these tests share the shared App Endpoint,
// so they run sequentially. Each test first applies the prerequisites on their own, then waits until the endpoint
// answers before the steps under test, because the endpoint can return 500 for a while after it is reported Online
// (AV-153193).

// TestAccAppEndpointAuditLogging creates App Endpoint audit logging with event IDs, a disabled user and a disabled
// role, checks the config and event IDs data sources against it, imports it, then updates every attribute.
func TestAccAppEndpointAuditLogging(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_ae_audit_")
	resourceReference := "couchbase-capella_app_endpoint_audit_logging." + name
	configReference := "data.couchbase-capella_app_endpoint_audit_logging_config." + name
	eventIDsReference := "data.couchbase-capella_app_endpoint_audit_log_event_ids." + name
	prerequisites := appEndpointAuditLoggingPrerequisites(name) + appEndpointAuditLogEventIDsConfig(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: appEndpointAuditLoggingPrerequisites(name),
			},
			{
				PreConfig: func() { waitForAppEndpointAuditLogEvents(t) },
				Config: prerequisites + appEndpointAuditLoggingConfig(name, name, true, `
  enabled_event_ids = local.filterable_event_ids
  disabled_users    = [{ domain = "sgw", name = "tf_acc_user" }]
  disabled_roles    = [{ domain = "sgw", name = "tf_acc_role" }]`) + appEndpointAuditLoggingConfigDataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "2"),
					testAccCheckEnabledEventIDsAreFirstFilterable(resourceReference, eventIDsReference, 2),
					resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "disabled_users.*", map[string]string{"domain": "sgw", "name": "tf_acc_user"}),
					resource.TestCheckResourceAttr(resourceReference, "disabled_roles.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "disabled_roles.*", map[string]string{"domain": "sgw", "name": "tf_acc_role"}),

					// The event IDs data source lists the endpoint's events, some of which can be filtered.
					resource.TestCheckResourceAttrWith(eventIDsReference, "data.#", atLeastOne),
					resource.TestCheckTypeSetElemNestedAttrs(eventIDsReference, "data.*", map[string]string{"filterable": "true"}),

					// The config data source reads the same configuration.
					resource.TestCheckResourceAttr(configReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttr(configReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttrPair(configReference, "enabled_event_ids.#", resourceReference, "enabled_event_ids.#"),
					testAccCheckEnabledEventIDsAreFirstFilterable(configReference, eventIDsReference, 2),
					resource.TestCheckResourceAttr(configReference, "disabled_users.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(configReference, "disabled_users.*", map[string]string{"domain": "sgw", "name": "tf_acc_user"}),
					resource.TestCheckResourceAttr(configReference, "disabled_roles.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(configReference, "disabled_roles.*", map[string]string{"domain": "sgw", "name": "tf_acc_role"}),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppEndpointAuditLoggingImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_endpoint_name",
			},
			{
				Config: prerequisites + appEndpointAuditLoggingConfig(name, name, false, `
  enabled_event_ids = [local.filterable_event_ids[0]]
  disabled_users    = []
  disabled_roles    = [{ domain = "sgw", name = "tf_acc_role" }]`) + appEndpointAuditLoggingConfigDataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
					testAccCheckEnabledEventIDsAreFirstFilterable(resourceReference, eventIDsReference, 1),
					resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "0"),
					resource.TestCheckResourceAttr(resourceReference, "disabled_roles.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "disabled_roles.*", map[string]string{"domain": "sgw", "name": "tf_acc_role"}),
					resource.TestCheckResourceAttr(configReference, "audit_enabled", "false"),
					resource.TestCheckResourceAttr(configReference, "enabled_event_ids.#", "1"),
					resource.TestCheckResourceAttr(configReference, "disabled_users.#", "0"),
					resource.TestCheckResourceAttr(configReference, "disabled_roles.#", "1"),
				),
			},
		},
	})
}

// TestAccAppEndpointAuditLoggingMinimal sets only audit_enabled on the shared App Endpoint.
func TestAccAppEndpointAuditLoggingMinimal(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_ae_audit_")
	resourceReference := "couchbase-capella_app_endpoint_audit_logging." + name
	prerequisites := appEndpointAuditLoggingPrerequisites(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: prerequisites,
			},
			{
				PreConfig: func() { waitForAppEndpointAuditLogEvents(t) },
				Config:    prerequisites + appEndpointAuditLoggingConfig(name, name, true, ""),
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
// configuration and Create leaves unconfigured attributes out of it, so they reset to their API defaults (AV-153185).
//
// Step 2 stands in for settings made outside Terraform: destroying this resource only removes it from state, so
// after step 3 the disabled user stays in Capella. Step 4 then adopts the endpoint with only audit_enabled.
// The prerequisites stay in every step, so turning audit logging off cannot clear the endpoint's settings.
func TestAccAppEndpointAuditLoggingAdoptKeepsUnconfigured(t *testing.T) {
	t.Skip("AV-153185: create resets App Endpoint audit settings the configuration leaves out; unskip once Create keeps them")

	name := randomStringWithPrefix("tf_acc_ae_audit_")
	adopted := "couchbase-capella_app_endpoint_audit_logging.adopt_" + name
	prerequisites := appEndpointAuditLoggingPrerequisites(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: prerequisites,
			},
			{
				PreConfig: func() { waitForAppEndpointAuditLogEvents(t) },
				Config: prerequisites + appEndpointAuditLoggingConfig("outside_"+name, name, true, `
  disabled_users = [{ domain = "sgw", name = "tf_acc_ui_user" }]`),
				Check: resource.TestCheckResourceAttr(
					"couchbase-capella_app_endpoint_audit_logging.outside_"+name, "disabled_users.#", "1"),
			},
			{
				Config: prerequisites,
			},
			{
				Config: prerequisites + appEndpointAuditLoggingConfig("adopt_"+name, name, true, ""),
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
// Gateway kept the previous events or applied its defaults, the read-back would not be empty and the last step
// would fail with "Provider produced inconsistent result after apply".
func TestAccAppEndpointAuditLoggingClearEventIDs(t *testing.T) {
	name := randomStringWithPrefix("tf_acc_ae_audit_")
	resourceReference := "couchbase-capella_app_endpoint_audit_logging." + name
	prerequisites := appEndpointAuditLoggingPrerequisites(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: prerequisites,
			},
			{
				PreConfig: func() { waitForAppEndpointAuditLogEvents(t) },
				Config: prerequisites + appEndpointAuditLogEventIDsConfig(name) + appEndpointAuditLoggingConfig(name, name, true, `
  enabled_event_ids = local.filterable_event_ids`),
				Check: resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "2"),
			},
			{
				Config: prerequisites + appEndpointAuditLoggingConfig(name, name, true, `
  enabled_event_ids = []`),
				Check: resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "0"),
			},
		},
	})
}

func TestAccAppEndpointAuditLoggingInvalidUUIDs(t *testing.T) {
	for _, attribute := range []string{"organization_id", "project_id", "cluster_id", "app_service_id"} {
		t.Run(attribute, func(t *testing.T) {
			ids := map[string]string{
				"organization_id": globalOrgId,
				"project_id":      globalProjectId,
				"cluster_id":      globalClusterId,
				"app_service_id":  globalAppServiceId,
			}
			ids[attribute] = "not-a-uuid"
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_endpoint_audit_logging" "invalid" {
  organization_id   = "%[2]s"
  project_id        = "%[3]s"
  cluster_id        = "%[4]s"
  app_service_id    = "%[5]s"
  app_endpoint_name = "%[6]s"
  audit_enabled     = true
}
`, globalProviderBlock, ids["organization_id"], ids["project_id"], ids["cluster_id"], ids["app_service_id"], globalAppEndpointName),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*` + attribute + `.*must be a valid UUID`),
					},
				},
			})
		})
	}
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

// appEndpointAuditLogEventIDsConfig returns the event IDs data source for the shared App Endpoint, and
// local.filterable_event_ids holding its two lowest filterable event IDs in ascending order.
func appEndpointAuditLogEventIDsConfig(name string) string {
	return fmt.Sprintf(`
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
}

// appEndpointAuditLoggingConfig returns the App Endpoint audit logging resource, depending on the prerequisites
// named prerequisitesName, with extra attribute lines in settings.
func appEndpointAuditLoggingConfig(resourceName, prerequisitesName string, auditEnabled bool, settings string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_app_endpoint_audit_logging" "%[1]s" {
  organization_id   = "%[2]s"
  project_id        = "%[3]s"
  cluster_id        = "%[4]s"
  app_service_id    = "%[5]s"
  app_endpoint_name = "%[6]s"
  audit_enabled     = %[7]t
%[8]s
  depends_on = [
    couchbase-capella_app_service_audit_logging.%[9]s,
    couchbase-capella_app_endpoint_activation_status.%[9]s,
  ]
}
`, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName,
		auditEnabled, settings, prerequisitesName)
}

// appEndpointAuditLoggingConfigDataSource returns the config data source reading the resource named name.
func appEndpointAuditLoggingConfigDataSource(name string) string {
	return fmt.Sprintf(`
data "couchbase-capella_app_endpoint_audit_logging_config" "%[1]s" {
  organization_id   = "%[2]s"
  project_id        = "%[3]s"
  cluster_id        = "%[4]s"
  app_service_id    = "%[5]s"
  app_endpoint_name = "%[6]s"
  depends_on        = [couchbase-capella_app_endpoint_audit_logging.%[1]s]
}
`, name, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, globalAppEndpointName)
}

func generateAppEndpointAuditLoggingImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		a := rs.Primary.Attributes
		return fmt.Sprintf("app_endpoint_name=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			a["app_endpoint_name"], a["app_service_id"], a["cluster_id"], a["project_id"], a["organization_id"]), nil
	}
}

// testAccCheckEnabledEventIDsAreFirstFilterable checks that the enabled_event_ids of reference are exactly the count
// lowest filterable event IDs listed by the event IDs data source, which is what local.filterable_event_ids holds.
func testAccCheckEnabledEventIDsAreFirstFilterable(reference, eventIDsReference string, count int) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		events, ok := state.RootModule().Resources[eventIDsReference]
		if !ok {
			return fmt.Errorf("%s not found in state", eventIDsReference)
		}
		// Set elements are flattened as data.<hash>.<field>.
		var filterable []int
		for key, value := range events.Primary.Attributes {
			element, found := strings.CutSuffix(key, ".filterable")
			if !found || !strings.HasPrefix(element, "data.") || value != "true" {
				continue
			}
			id, err := strconv.Atoi(events.Primary.Attributes[element+".id"])
			if err != nil {
				return fmt.Errorf("reading event ID of %s: %w", element, err)
			}
			filterable = append(filterable, id)
		}
		sort.Ints(filterable)
		if len(filterable) < count {
			return fmt.Errorf("%s lists %d filterable events, want at least %d", eventIDsReference, len(filterable), count)
		}
		checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(reference, "enabled_event_ids.#", strconv.Itoa(count))}
		for _, id := range filterable[:count] {
			checks = append(checks, resource.TestCheckTypeSetElemAttr(reference, "enabled_event_ids.*", strconv.Itoa(id)))
		}
		return resource.ComposeAggregateTestCheckFunc(checks...)(state)
	}
}

func atLeastOne(value string) error {
	n, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("got %d, want at least 1", n)
	}
	return nil
}
