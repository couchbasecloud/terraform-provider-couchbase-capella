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
)

// App Service audit logging is one setting per App Service, and several tests in this package turn it on and off on
// the shared App Service, so these tests run sequentially.

// TestAccAppServiceAuditLogging creates App Service audit logging, reads it back through the state data source,
// imports it, and checks that destroy turns audit logging off.
func TestAccAppServiceAuditLogging(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_as_audit_logging_")
	resourceReference := "couchbase-capella_app_service_audit_logging." + resourceName
	dataSourceReference := "data.couchbase-capella_app_service_audit_logging_state." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAuditLoggingDisabled(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAuditLoggingConfig(resourceName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					testAccCheckAppServiceAuditLoggingEnabled(t, true),
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(dataSourceReference, "audit_enabled", "true"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppServiceScopedImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_service_id",
			},
		},
	})
}

// TestAccAppServiceAuditLoggingImportDisabled checks that importing fails while audit logging is off, because the
// resource only exists while audit logging is enabled.
func TestAccAppServiceAuditLoggingImportDisabled(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_as_audit_logging_")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:        testAccAppServiceAuditLoggingConfig(resourceName, false),
				ResourceName:  "couchbase-capella_app_service_audit_logging." + resourceName,
				ImportState:   true,
				ImportStateId: fmt.Sprintf("app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", globalAppServiceId, globalClusterId, globalProjectId, globalOrgId),
				ExpectError:   regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

func TestAccAppServiceAuditLoggingInvalidUUIDs(t *testing.T) {
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
resource "couchbase-capella_app_service_audit_logging" "invalid" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
}
`, globalProviderBlock, ids["organization_id"], ids["project_id"], ids["cluster_id"], ids["app_service_id"]),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*` + attribute + `.*must be a valid UUID`),
					},
				},
			})
		})
	}
}

// testAccAppServiceAuditLoggingConfig returns the audit logging resource and, with withDataSource, the state data
// source reading it.
func testAccAppServiceAuditLoggingConfig(resourceName string, withDataSource bool) string {
	config := fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service_audit_logging" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
	if !withDataSource {
		return config
	}
	return config + fmt.Sprintf(`
data "couchbase-capella_app_service_audit_logging_state" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  depends_on      = [couchbase-capella_app_service_audit_logging.%[1]s]
}
`, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

// generateAppServiceScopedImportId returns the import ID for a resource identified by its App Service.
func generateAppServiceScopedImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		a := rs.Primary.Attributes
		return fmt.Sprintf("app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			a["app_service_id"], a["cluster_id"], a["project_id"], a["organization_id"]), nil
	}
}

func testAccCheckAppServiceAuditLoggingEnabled(t *testing.T, want bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		got, err := getAppServiceAuditLoggingEnabled(t)
		if err != nil {
			return err
		}
		if got != want {
			return fmt.Errorf("App Service audit logging enabled = %t in Capella, want %t", got, want)
		}
		return nil
	}
}

func testAccCheckAppServiceAuditLoggingDisabled(t *testing.T) resource.TestCheckFunc {
	return testAccCheckAppServiceAuditLoggingEnabled(t, false)
}

func getAppServiceAuditLoggingEnabled(t *testing.T) (bool, error) {
	data := newTestClient(t)
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLog",
		globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return false, fmt.Errorf("reading App Service audit logging state: %w", err)
	}
	var state struct {
		AuditEnabled bool `json:"auditEnabled"`
	}
	if err := json.Unmarshal(response.Body, &state); err != nil {
		return false, fmt.Errorf("unmarshalling App Service audit logging state: %w", err)
	}
	return state.AuditEnabled, nil
}
