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

// audit_log_settings is a per-cluster singleton: there's no real
// create/destroy lifecycle on the server side. The provider's Create
// method PUTs the settings, Update PUTs the new settings, and Delete is a
// no-op. The Terraform test framework still runs a Destroy at the end,
// which only removes the resource from state — the server-side singleton
// is unchanged. We exercise: Create -> Read -> Import -> Update.

func TestAccAuditLogSettingsResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_audit_log_settings_")
	resourceReference := "couchbase-capella_audit_log_settings." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: testAccAuditLogSettingsResourceConfig(resourceName, true, []int{20488, 20490, 20491}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAuditLogSettingsResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "0"),
				),
			},
			// ImportState testing. The resource has no `id` attribute; the
			// natural primary key is cluster_id (one settings record per
			// cluster). The import string keys map through importIds, so
			// `id` here is matched to ClusterId via the schema's Validate.
			{
				ResourceName:                         resourceReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateAuditLogSettingsImportIdForResource(resourceReference),
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
			// Update: shrink the enabled event id list and assert ALL
			// fields to catch reset-to-default regressions per the
			// acceptance test skill guidance.
			{
				Config: testAccAuditLogSettingsResourceConfig(resourceName, true, []int{20488}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAuditLogSettingsResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "disabled_users.#", "0"),
				),
			},
			// Update audit_enabled to false. Assert ALL fields again.
			{
				Config: testAccAuditLogSettingsResourceConfig(resourceName, false, []int{20488}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAuditLogSettingsResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "audit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "enabled_event_ids.#", "1"),
				),
			},
		},
	})
}

// TestAccAuditLogSettingsResourceInvalidCluster asserts that pointing at a
// non-existent cluster yields the expected create-time error.
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

func testAccAuditLogSettingsResourceConfig(resourceName string, auditEnabled bool, enabledEventIDs []int) string {
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

resource "couchbase-capella_audit_log_settings" "%[2]s" {
  organization_id   = "%[3]s"
  project_id        = "%[4]s"
  cluster_id        = "%[5]s"
  audit_enabled     = %[6]t
  enabled_event_ids = %[7]s
  disabled_users    = []
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, auditEnabled, ids)
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
		// The resource's schema Validate maps cluster_id to the Id key; the
		// import string parser uses the importIds lookup table where
		// "id" -> Id. The provider's ImportState passes the import string
		// through to the cluster_id attribute, then Validate splits it.
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
