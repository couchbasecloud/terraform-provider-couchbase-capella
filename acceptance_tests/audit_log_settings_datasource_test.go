package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceAuditLogSettings(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_audit_cluster_for_ds_")
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_for_ds_")
	dsName := randomStringWithPrefix("tf_acc_audit_log_settings_ds_")
	cidr := generateRandomCIDR()
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	dsReference := "data.couchbase-capella_audit_log_settings." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditLogSettingsResourceAndDatasourceConfigWithEnterpriseCluster(clusterResourceName, resourceName, dsName, cidr, true, []int{20488, 20489}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrPair(dsReference, "cluster_id", clusterReference, "id"),
					resource.TestCheckResourceAttr(dsReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttr(dsReference, "enabled_event_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr(dsReference, "enabled_event_ids.*", "20488"),
					resource.TestCheckTypeSetElemAttr(dsReference, "enabled_event_ids.*", "20489"),
					resource.TestCheckResourceAttr(dsReference, "disabled_users.#", "0"),
				),
			},
		},
	})
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

func testAccAuditLogSettingsResourceAndDatasourceConfigWithEnterpriseCluster(clusterResourceName, auditSettingsResourceName, dsName, cidr string, auditEnabled bool, enabledEventIDs []int) string {
	ids := "["
	for i, id := range enabledEventIDs {
		if i > 0 {
			ids += ", "
		}
		ids += fmt.Sprintf("%d", id)
	}
	ids += "]"

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[2]s" {
	organization_id = "%[5]s"
	project_id      = "%[6]s"
	name            = "%[2]s"

	cloud_provider = {
		type   = "aws"
		region = "us-east-1"
		cidr   = "%[4]s"
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
		plan     = "enterprise"
		timezone = "PT"
	}
}

resource "couchbase-capella_audit_log_settings" "%[3]s" {
	organization_id   = "%[5]s"
	project_id        = "%[6]s"
	cluster_id        = couchbase-capella_cluster.%[2]s.id
	audit_enabled     = %[7]t
	enabled_event_ids = %[8]s
  disabled_users    = []
}

data "couchbase-capella_audit_log_settings" "%[9]s" {
	organization_id = "%[5]s"
	project_id      = "%[6]s"
	cluster_id      = couchbase-capella_cluster.%[2]s.id

	depends_on = [couchbase-capella_audit_log_settings.%[3]s]
}
`, globalProviderBlock, clusterResourceName, auditSettingsResourceName, cidr, globalOrgId, globalProjectId, auditEnabled, ids, dsName)
}
