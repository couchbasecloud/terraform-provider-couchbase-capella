package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// auditLogSettingsDatasourceSteps provides the steps to read the audit log settings of the
// global cluster through the datasource after they have been set by the resource.
func auditLogSettingsDatasourceSteps() []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_for_ds_")
	dsName := randomStringWithPrefix("tf_acc_audit_log_settings_ds_")
	dsReference := "data.couchbase-capella_audit_log_settings." + dsName

	return []resource.TestStep{
		{
			Config: testAccAuditLogSettingsResourceAndDatasourceConfig(resourceName, dsName, true, []int{20488, 20489}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(dsReference, "audit_enabled", "true"),
				resource.TestCheckResourceAttr(dsReference, "enabled_event_ids.#", "2"),
				resource.TestCheckTypeSetElemAttr(dsReference, "enabled_event_ids.*", "20488"),
				resource.TestCheckTypeSetElemAttr(dsReference, "enabled_event_ids.*", "20489"),
				resource.TestCheckResourceAttr(dsReference, "disabled_users.#", "0"),
			),
		},
	}
}

func TestAccDatasourceAuditLogSettingsInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_audit_log_settings_ds_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_audit_log_settings" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Audit Log Settings|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceAuditLogSettingsMissingCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_audit_log_settings_ds_missing_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_audit_log_settings" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)cluster_id|argument.*required`),
			},
		},
	})
}

func testAccAuditLogSettingsResourceAndDatasourceConfig(auditSettingsResourceName, dsName string, auditEnabled bool, enabledEventIDs []int) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_audit_log_settings" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"

	depends_on = [couchbase-capella_audit_log_settings.%[6]s]
}
`,
		testAccAuditLogSettingsResourceConfig(auditSettingsResourceName, auditEnabled, enabledEventIDs),
		dsName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		auditSettingsResourceName,
	)
}
