package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// audit_log_settings is a per-cluster singleton — Create/Update both PUT and
// Delete disables auditing. TestAccAuditLog runs the tests that need a real
// cluster as sequential subtests against the shared global cluster so they do
// not overwrite each other's settings, and so the suite no longer deploys a
// cluster per audit test.
func TestAccAuditLog(t *testing.T) {
	// Allow this test to run in parallel with other top-level tests, but ensure that the subtests run sequentially
	// This is normally set by resource.ParallelTest
	t.Parallel()

	t.Run("Audit Log Settings", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    auditLogSettingsSteps(t),
		})
	})

	t.Run("Audit Log Settings Datasource", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    auditLogSettingsDatasourceSteps(),
		})
	})

	t.Run("Audit Log Settings Disable On Removal", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    auditLogSettingsDisableOnRemovalSteps(t),
		})
	})
}

// auditLogSettingsSteps provides the steps to test the lifecycle of the
// audit_log_settings resource on the global cluster: Create -> Import ->
// Update (disabled users) -> Update (audit disabled).
func auditLogSettingsSteps(t *testing.T) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_")
	databaseCredentialName := randomStringWithPrefix("tf_acc_audit_db_credential_")
	resourceReference := "couchbase-capella_audit_log_settings." + resourceName

	return []resource.TestStep{
		{
			Config: testAccAuditLogSettingsResourceConfig(resourceName, true, []int{20488, 20490, 20491}),
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccExistsAuditLogSettingsResource(t, resourceReference),
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
				resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "3"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20488"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20490"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20491"),
				resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "0"),
			),
		},
		{
			ResourceName:                         resourceReference,
			ImportState:                          true,
			ImportStateIdFunc:                    generateAuditLogSettingsImportIdForResource(resourceReference),
			ImportStateVerifyIdentifierAttribute: "cluster_id",
		},
		{
			Config: testAccAuditLogSettingsResourceConfigWithDatabaseCredentialAndDisabledUsers(resourceName, databaseCredentialName, true, []int{20488}, `[
					{
						domain = "local"
						name   = "%[1]s"
					}
				]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccExistsAuditLogSettingsResource(t, resourceReference),
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
				resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20488"),
				resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "disabled_users.*", map[string]string{
					"domain": "local",
					"name":   databaseCredentialName,
				}),
			),
		},
		{
			Config: testAccAuditLogSettingsResourceConfigWithDatabaseCredentialAndDisabledUsers(resourceName, databaseCredentialName, false, []int{20488}, "[]"),
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccExistsAuditLogSettingsResource(t, resourceReference),
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "false"),
				resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20488"),
			),
		},
	}
}

// auditLogSettingsDisableOnRemovalSteps provides the steps to verify that removing the
// audit_log_settings resource from the config reverts audit logging on the cluster to
// disabled, which is what allows a later downgrade from the enterprise support plan.
func auditLogSettingsDisableOnRemovalSteps(t *testing.T) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_audit_rm_settings_")
	resourceReference := "couchbase-capella_audit_log_settings." + resourceName

	return []resource.TestStep{
		{
			Config: testAccAuditLogSettingsResourceConfig(resourceName, true, []int{20488, 20490, 20491}),
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccExistsAuditLogSettingsResource(t, resourceReference),
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
				resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "3"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20488"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20490"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20491"),
				resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "0"),
			),
		},
		{
			Config: globalProviderBlock,
			Check:  testAccCheckAuditLogSettingsDisabled(t),
		},
	}
}

func TestAccAuditLogSettingsResourceInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
  organization_id   = "%[3]s"
  project_id        = "%[4]s"
  cluster_id        = "00000000-0000-0000-0000-000000000000"
  audit_enabled     = true
  enabled_event_ids = [20488]
  disabled_users    = []
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)error during audit log settings creation|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccAuditLogSettingsResourceEmptyOrganizationID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_empty_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogSettingsResourceEmptyOrganizationIDConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Attribute organization_id string length must be at least 1, got: 0`),
			},
		},
	})
}

func TestAccAuditLogSettingsResourceEventIDWrongType(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_bad_event_id_type_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogSettingsResourceEventIDWrongTypeConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)(Incorrect attribute value type|Invalid value for input variable|Invalid Attribute Value).*enabled_event_ids`),
			},
		},
	})
}

func TestAccAuditLogSettingsResourceNegativeEventID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_negative_event_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogSettingsResourceNegativeEventIDConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)enabled_event_ids.*value must be at least 1`),
			},
		},
	})
}

func TestAccAuditLogSettingsResourceDisabledUserEmptyDomain(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_empty_domain_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogSettingsResourceDisabledUserEmptyDomainConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Attribute.*disabled_users.*domain string\s+length must be at least 1, got: 0`),
			},
		},
	})
}

func TestAccAuditLogSettingsResourceDisabledUserEmptyName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_empty_name_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogSettingsResourceDisabledUserEmptyNameConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Attribute.*disabled_users.*name string\s+length must be at least 1, got: 0`),
			},
		},
	})
}

func TestAccAuditLogSettingsResourceDisabledUserMissingName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_missing_name_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogSettingsResourceDisabledUserMissingNameConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)(Missing required|Incorrect attribute value type|attribute).*disabled_users.*name.*required`),
			},
		},
	})
}

// TestAccAuditLogSettingsResourcePlanUpgrade is the one audit test that still deploys its
// own cluster: it asserts audit log settings are rejected on a developer pro cluster and
// accepted once the cluster is upgraded to enterprise, so it cannot run against the shared
// global cluster, which is already on the enterprise plan.
func TestAccAuditLogSettingsResourcePlanUpgrade(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_audit_cluster_upgrade_")
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_upgrade_")
	cidr := generateRandomCIDR()
	resourceReference := "couchbase-capella_audit_log_settings." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAuditLogSettingsResourceConfigWithSupportPlan(clusterResourceName, resourceName, cidr, "developer pro", true, []int{20488}),
				ExpectError: regexp.MustCompile(`(?s)error during audit log settings creation|support package`),
			},
			{
				Config: testAccAuditLogSettingsResourceConfigWithSupportPlan(clusterResourceName, resourceName, cidr, "enterprise", true, []int{20488}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20488"),
				),
			},
		},
	})
}

func testAccAuditLogSettingsResourceDisabledUserEmptyDomainConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id   = "00000000-0000-0000-0000-000000000000"
	project_id        = "11111111-1111-1111-1111-111111111111"
	cluster_id        = "22222222-2222-2222-2222-222222222222"
	audit_enabled     = true
	enabled_event_ids = [20488]

	disabled_users = [
		{
			domain = ""
			name   = "audit-exempt-user"
		}
	]
}
`, globalProviderBlock, resourceName)
}

func testAccAuditLogSettingsResourceDisabledUserEmptyNameConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id   = "00000000-0000-0000-0000-000000000000"
	project_id        = "11111111-1111-1111-1111-111111111111"
	cluster_id        = "22222222-2222-2222-2222-222222222222"
	audit_enabled     = true
	enabled_event_ids = [20488]

	disabled_users = [
		{
			domain = "local"
			name   = ""
		}
	]
}
`, globalProviderBlock, resourceName)
}

func testAccAuditLogSettingsResourceDisabledUserMissingNameConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id   = "00000000-0000-0000-0000-000000000000"
	project_id        = "11111111-1111-1111-1111-111111111111"
	cluster_id        = "22222222-2222-2222-2222-222222222222"
	audit_enabled     = true
	enabled_event_ids = [20488]

	disabled_users = [
		{
			domain = "local"
		}
	]
}
`, globalProviderBlock, resourceName)
}

func testAccAuditLogSettingsResourceEmptyOrganizationIDConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id   = ""
	project_id        = "11111111-1111-1111-1111-111111111111"
	cluster_id        = "22222222-2222-2222-2222-222222222222"
	audit_enabled     = true
	enabled_event_ids = [20488]
	disabled_users    = []
}
`, globalProviderBlock, resourceName)
}

func testAccAuditLogSettingsResourceEventIDWrongTypeConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id   = "00000000-0000-0000-0000-000000000000"
	project_id        = "11111111-1111-1111-1111-111111111111"
	cluster_id        = "22222222-2222-2222-2222-222222222222"
	audit_enabled     = true
	enabled_event_ids = ["not-a-number"]
	disabled_users    = []
}
`, globalProviderBlock, resourceName)
}

func testAccAuditLogSettingsResourceNegativeEventIDConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id   = "00000000-0000-0000-0000-000000000000"
	project_id        = "11111111-1111-1111-1111-111111111111"
	cluster_id        = "22222222-2222-2222-2222-222222222222"
	audit_enabled     = true
	enabled_event_ids = [-1]
	disabled_users    = []
}
`, globalProviderBlock, resourceName)
}

func testAccAuditLogSettingsResourceConfigWithSupportPlan(clusterResourceName, auditSettingsResourceName, cidr, plan string, auditEnabled bool, enabledEventIDs []int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	name            = "%[2]s"

	cloud_provider = {
		type   = "aws"
		region = "us-east-1"
		cidr   = "%[5]s"
	}

	service_groups = [
		{
			node = {
				compute = {
					cpu = 4
					ram = 16
				}
				disk = {
					storage = 50
					type    = "io2"
					iops    = 3000
				}
			}
			num_of_nodes = 3
			services     = ["data", "index", "query"]
		}
	]

	availability = {
		type = "multi"
	}

	support = {
		plan     = "%[8]s"
		timezone = "PT"
	}
}

resource "couchbase-capella_audit_log_settings" "%[6]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = couchbase-capella_cluster.%[2]s.id
	audit_enabled     = %[7]t
	enabled_event_ids = %[9]s
	disabled_users    = []
}
`, globalProviderBlock, clusterResourceName, globalOrgId, globalProjectId, cidr, auditSettingsResourceName, auditEnabled, plan, formatEventIDs(enabledEventIDs))
}

// formatEventIDs renders event IDs as an HCL list literal.
func formatEventIDs(enabledEventIDs []int) string {
	ids := "["
	for i, id := range enabledEventIDs {
		if i > 0 {
			ids += ", "
		}
		ids += fmt.Sprintf("%d", id)
	}
	return ids + "]"
}

func testAccAuditLogSettingsResourceConfig(auditSettingsResourceName string, auditEnabled bool, enabledEventIDs []int) string {
	return testAccAuditLogSettingsResourceConfigWithDisabledUsers(auditSettingsResourceName, auditEnabled, enabledEventIDs, "[]", "", "")
}

func testAccAuditLogSettingsResourceConfigWithDatabaseCredentialAndDisabledUsers(auditSettingsResourceName, databaseCredentialName string, auditEnabled bool, enabledEventIDs []int, disabledUsers string) string {
	databaseCredentialResource := fmt.Sprintf(`
resource "couchbase-capella_database_credential" "%[1]s" {
	name            = "%[1]s"
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	cluster_id      = "%[4]s"
	access = [
		{
			privileges = ["data_writer"]
		},
	]
}
`, databaseCredentialName, globalOrgId, globalProjectId, globalClusterId)

	formattedDisabledUsers := disabledUsers
	if disabledUsers != "[]" {
		formattedDisabledUsers = fmt.Sprintf(disabledUsers, databaseCredentialName)
	}

	return testAccAuditLogSettingsResourceConfigWithDisabledUsers(
		auditSettingsResourceName,
		auditEnabled,
		enabledEventIDs,
		formattedDisabledUsers,
		databaseCredentialName,
		databaseCredentialResource,
	)
}

func testAccAuditLogSettingsResourceConfigWithDisabledUsers(auditSettingsResourceName string, auditEnabled bool, enabledEventIDs []int, disabledUsers, databaseCredentialName, databaseCredentialResource string) string {
	dependsOn := ""
	if databaseCredentialResource != "" {
		dependsOn = fmt.Sprintf(`

	depends_on = [couchbase-capella_database_credential.%[1]s]`, databaseCredentialName)
	}

	return fmt.Sprintf(`
%[1]s

%[8]s

resource "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	audit_enabled     = %[6]t
	enabled_event_ids = %[7]s
	disabled_users    = %[9]s
%[10]s
}
`, globalProviderBlock, auditSettingsResourceName, globalOrgId, globalProjectId, globalClusterId, auditEnabled, formatEventIDs(enabledEventIDs), databaseCredentialResource, disabledUsers, dependsOn)
}

func generateAuditLogSettingsImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"id=%s,project_id=%s,organization_id=%s",
			rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func retrieveAuditLogSettingsFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*api.GetClusterAuditSettingsResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/auditLog",
		data.HostURL, organizationId, projectId, clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return nil, err
	}
	settings := &api.GetClusterAuditSettingsResponse{}
	if err := json.Unmarshal(response.Body, settings); err != nil {
		return nil, err
	}
	return settings, nil
}

func testAccExistsAuditLogSettingsResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		data := newTestClient(t)
		_, err := retrieveAuditLogSettingsFromServer(
			data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"],
		)
		return err
	}
}

func testAccCheckAuditLogSettingsDisabled(t *testing.T) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		data := newTestClient(t)
		settings, err := retrieveAuditLogSettingsFromServer(data, globalOrgId, globalProjectId, globalClusterId)
		if err != nil {
			return err
		}
		if settings.AuditEnabled {
			return fmt.Errorf("expected audit logging disabled after resource removal, got auditEnabled=true")
		}
		if len(settings.EnabledEventIDs) != 0 {
			return fmt.Errorf("expected no enabled event ids after removal, got %d", len(settings.EnabledEventIDs))
		}
		return nil
	}
}
