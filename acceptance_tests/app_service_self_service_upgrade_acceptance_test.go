package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_self_service_upgrade"
)

func TestAccDatasourceAppServiceSelfServiceUpgradeAvailableVersions(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_self_service_upgrade_available_versions_")
	dataSourceReference := "data.couchbase-capella_app_service_self_service_upgrade_available_versions." + dataSourceName

	versions, err := listAvailableSelfServiceUpgradeVersions(context.Background())
	if err != nil {
		t.Fatalf("could not list the available versions of app service %s: %v", globalAppServiceId, err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
		resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
		resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
		resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
		resource.TestCheckResourceAttr(dataSourceReference, "data.#", fmt.Sprint(len(versions))),
	}
	// The data source must expose every field of each version as the API returns it.
	for i, v := range versions {
		item := fmt.Sprintf("data.%d.", i)
		checks = append(checks,
			resource.TestCheckResourceAttr(dataSourceReference, item+"id", v.Id),
			resource.TestCheckResourceAttr(dataSourceReference, item+"title", v.Title),
			resource.TestCheckResourceAttr(dataSourceReference, item+"description", v.Description),
			resource.TestCheckResourceAttr(dataSourceReference, item+"priority", v.Priority),
			resource.TestCheckResourceAttr(dataSourceReference, item+"release_notes.#", fmt.Sprint(len(v.ReleaseNotes))),
		)
		for j, note := range v.ReleaseNotes {
			checks = append(checks, resource.TestCheckResourceAttr(dataSourceReference, fmt.Sprintf("%srelease_notes.%d", item, j), note))
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceSelfServiceUpgradeAvailableVersionsConfig(dataSourceName),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccDatasourceAppServiceSelfServiceUpgrades(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_self_service_upgrades_")
	dataSourceReference := "data.couchbase-capella_app_service_self_service_upgrades." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceSelfServiceUpgradesConfig(dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "data.#"),
				),
			},
		},
	})
}

// TestAccAppServiceSelfServiceUpgradeResource schedules an upgrade of the shared App Service days
// ahead and cancels it on destroy, so the App Service is never upgraded. It runs sequentially, not in
// parallel, because TestAccAppServiceMaintenanceSchedule reschedules any pending maintenance of the
// same App Service, which would include this upgrade.
func TestAccAppServiceSelfServiceUpgradeResource(t *testing.T) {
	versions, err := listAvailableSelfServiceUpgradeVersions(context.Background())
	if err != nil {
		t.Fatalf("could not list the available versions of app service %s: %v", globalAppServiceId, err)
	}
	if len(versions) == 0 {
		t.Skip("no self-service upgrade version is available for app service " + globalAppServiceId)
	}
	versionId := versions[0].Id

	resourceName := randomStringWithPrefix("tf_acc_app_service_self_service_upgrade_")
	resourceReference := "couchbase-capella_app_service_self_service_upgrade." + resourceName
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_self_service_upgrades_")
	dataSourceReference := "data.couchbase-capella_app_service_self_service_upgrades." + dataSourceName

	// The API accepts a queue time at most 7 days ahead.
	now := time.Now().UTC().Truncate(time.Minute)
	queueAt := now.Add(5 * 24 * time.Hour).Format(time.RFC3339)
	laterQueueAt := now.Add(6 * 24 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceSelfServiceUpgradeCancelled(resourceReference),
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAppServiceSelfServiceUpgradeResourceConfig(resourceName, versionId, queueAt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsAppServiceSelfServiceUpgradeResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "version_id", versionId),
					resource.TestCheckResourceAttr(resourceReference, "queue_at", queueAt),
					resource.TestCheckResourceAttr(resourceReference, "current_state", "pending"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckNoResourceAttr(resourceReference, "completed_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppServiceSelfServiceUpgradeImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update, with the list data source reading the same upgrade
			{
				Config: testAccAppServiceSelfServiceUpgradeResourceConfig(resourceName, versionId, laterQueueAt) +
					testAccAppServiceSelfServiceUpgradesDataSourceBlock(dataSourceName, resourceReference),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "version_id", versionId),
					resource.TestCheckResourceAttr(resourceReference, "queue_at", laterQueueAt),
					resource.TestCheckResourceAttr(resourceReference, "current_state", "pending"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					testAccCheckSelfServiceUpgradesListContains(dataSourceReference, resourceReference),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppServiceSelfServiceUpgradeResourceConfig(resourceName, versionId, queueAt string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_self_service_upgrade" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  version_id      = "%[7]s"
  queue_at        = "%[8]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, versionId, queueAt)
}

// testAccAppServiceSelfServiceUpgradesDataSourceBlock reads the list after the resource is applied.
func testAccAppServiceSelfServiceUpgradesDataSourceBlock(dataSourceName, resourceReference string) string {
	return fmt.Sprintf(`
data "couchbase-capella_app_service_self_service_upgrades" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"

  depends_on = [%[6]s]
}
`, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, resourceReference)
}

func testAccAppServiceSelfServiceUpgradeAvailableVersionsConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_self_service_upgrade_available_versions" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}
`, globalProviderBlock, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

func testAccAppServiceSelfServiceUpgradesConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_self_service_upgrades" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}
`, globalProviderBlock, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

func generateAppServiceSelfServiceUpgradeImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rawState, err := appServiceSelfServiceUpgradeAttributes(state, resourceReference)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"id=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["app_service_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func appServiceSelfServiceUpgradeAttributes(state *terraform.State, resourceReference string) (map[string]string, error) {
	for _, m := range state.Modules {
		if v, ok := m.Resources[resourceReference]; ok {
			return v.Primary.Attributes, nil
		}
	}
	return nil, fmt.Errorf("resource %s not found in state", resourceReference)
}

func appServiceSelfServiceUpgradesURL() string {
	return fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/selfServiceUpgrades",
		globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

func getAppServiceSelfServiceUpgrade(ctx context.Context, upgradeId string) (*appservice_maintenance.GetMaintenanceResponse, error) {
	cfg := api.EndpointCfg{Url: appServiceSelfServiceUpgradesURL() + "/" + upgradeId, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := globalClient.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		return nil, err
	}
	var upgrade appservice_maintenance.GetMaintenanceResponse
	if err = json.Unmarshal(response.Body, &upgrade); err != nil {
		return nil, err
	}
	return &upgrade, nil
}

func listAvailableSelfServiceUpgradeVersions(ctx context.Context) ([]appservice_self_service_upgrade.AvailableVersion, error) {
	cfg := api.EndpointCfg{Url: appServiceSelfServiceUpgradesURL() + "/availableVersions", Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := globalClient.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		return nil, err
	}
	var versions appservice_self_service_upgrade.ListAvailableVersionsResponse
	if err = json.Unmarshal(response.Body, &versions); err != nil {
		return nil, err
	}
	return versions.Data, nil
}

func testAccExistsAppServiceSelfServiceUpgradeResource(resourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rawState, err := appServiceSelfServiceUpgradeAttributes(state, resourceReference)
		if err != nil {
			return err
		}
		_, err = getAppServiceSelfServiceUpgrade(context.Background(), rawState["id"])
		return err
	}
}

// testAccCheckSelfServiceUpgradesListContains asserts the list data source holds the resource's upgrade.
func testAccCheckSelfServiceUpgradesListContains(dataSourceReference, resourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		upgrade, err := appServiceSelfServiceUpgradeAttributes(state, resourceReference)
		if err != nil {
			return err
		}
		list, err := appServiceSelfServiceUpgradeAttributes(state, dataSourceReference)
		if err != nil {
			return err
		}
		var count int
		if _, err = fmt.Sscan(list["data.#"], &count); err != nil {
			return fmt.Errorf("%s has no data: %w", dataSourceReference, err)
		}
		for i := 0; i < count; i++ {
			item := fmt.Sprintf("data.%d.", i)
			if list[item+"id"] == upgrade["id"] {
				return checkSelfServiceUpgradeListItem(list, item, upgrade)
			}
		}
		return fmt.Errorf("%s does not list upgrade %s", dataSourceReference, upgrade["id"])
	}
}

// checkSelfServiceUpgradeListItem asserts a list item maps the same pending upgrade the resource holds.
func checkSelfServiceUpgradeListItem(list map[string]string, item string, upgrade map[string]string) error {
	expected := map[string]string{
		"kind":          "upgrade",
		"initiated_by":  "selfService",
		"version_id":    upgrade["version_id"],
		"current_state": upgrade["current_state"],
	}
	for field, want := range expected {
		if got := list[item+field]; got != want {
			return fmt.Errorf("%s%s is %q, expected %q", item, field, got, want)
		}
	}

	listQueueAt, err := time.Parse(time.RFC3339, list[item+"queue_at"])
	if err != nil {
		return fmt.Errorf("%squeue_at: %w", item, err)
	}
	resourceQueueAt, err := time.Parse(time.RFC3339, upgrade["queue_at"])
	if err != nil {
		return fmt.Errorf("resource queue_at: %w", err)
	}
	if !listQueueAt.Equal(resourceQueueAt) {
		return fmt.Errorf("%squeue_at is %s, expected %s", item, listQueueAt, resourceQueueAt)
	}

	if completedAt := list[item+"completed_at"]; completedAt != "" {
		return fmt.Errorf("%scompleted_at is %q, expected null for a pending upgrade", item, completedAt)
	}
	for _, field := range []string{"audit.created_at", "audit.created_by", "audit.modified_at", "audit.modified_by", "audit.version"} {
		if list[item+field] == "" {
			return fmt.Errorf("%s%s is not set", item, field)
		}
	}
	return nil
}

// testAccCheckAppServiceSelfServiceUpgradeCancelled asserts destroy cancelled the upgrade, so the
// shared App Service is not upgraded.
func testAccCheckAppServiceSelfServiceUpgradeCancelled(resourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, m := range state.Modules {
			for name, rs := range m.Resources {
				if name != resourceReference {
					continue
				}
				upgrade, err := getAppServiceSelfServiceUpgrade(context.Background(), rs.Primary.Attributes["id"])
				if err != nil {
					if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
						return nil
					}
					return err
				}
				if upgrade.CurrentState != "cancelled" {
					return fmt.Errorf("upgrade %s is %s after destroy, expected cancelled", upgrade.Id, upgrade.CurrentState)
				}
			}
		}
		return nil
	}
}
