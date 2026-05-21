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

// audit_log_settings is a per-cluster singleton — Create/Update both PUT,
// Delete is a no-op. We exercise: Create -> Read -> Import -> Update.

func TestAccAuditLogSettingsResource(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_audit_cluster_")
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_")
	cidr := generateRandomCIDR()
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	resourceReference := "couchbase-capella_audit_log_settings." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditLogSettingsResourceConfigWithEnterpriseCluster(clusterResourceName, resourceName, cidr, true, []int{20488, 20490, 20491}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAuditLogSettingsResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrPair(resourceReference, "cluster_id", clusterReference, "id"),
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
				Config: testAccAuditLogSettingsResourceConfigWithEnterpriseCluster(clusterResourceName, resourceName, cidr, true, []int{20488}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAuditLogSettingsResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrPair(resourceReference, "cluster_id", clusterReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20488"),
					resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "0"),
				),
			},
			{
				Config: testAccAuditLogSettingsResourceConfigWithEnterpriseCluster(clusterResourceName, resourceName, cidr, false, []int{20488}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAuditLogSettingsResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrPair(resourceReference, "cluster_id", clusterReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "enabled_event_ids.*", "20488"),
				),
			},
		},
	})
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

func testAccAuditLogSettingsResourceConfigWithSupportPlan(clusterResourceName, auditSettingsResourceName, cidr, plan string, auditEnabled bool, enabledEventIDs []int) string {
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
`, globalProviderBlock, clusterResourceName, globalOrgId, globalProjectId, cidr, auditSettingsResourceName, auditEnabled, plan, ids)
}

func testAccAuditLogSettingsResourceConfigWithEnterpriseCluster(clusterResourceName, auditSettingsResourceName, cidr string, auditEnabled bool, enabledEventIDs []int) string {
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
		plan     = "enterprise"
		timezone = "PT"
	}
}

resource "couchbase-capella_audit_log_settings" "%[6]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = couchbase-capella_cluster.%[2]s.id
	audit_enabled     = %[7]t
	enabled_event_ids = %[8]s
	disabled_users    = []
}
`, globalProviderBlock, clusterResourceName, globalOrgId, globalProjectId, cidr, auditSettingsResourceName, auditEnabled, ids)
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
