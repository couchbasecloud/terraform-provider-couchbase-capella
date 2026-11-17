package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
)

func TestAccAppServiceMaintenancesDataSource(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_maintenances_")
	dataSourceReference := "data.couchbase-capella_app_service_maintenances." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccAppServiceMaintenancesDataSourceConfig(dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "data.#"),
				),
			},
		},
	})
}

// TestAccAppServiceMaintenanceSchedule needs a pending maintenance on the shared App Service, which
// the provider cannot create; it skips when there is none. It restores the original queue time.
func TestAccAppServiceMaintenanceSchedule(t *testing.T) {
	maintenance, err := findReschedulableAppServiceMaintenance(context.Background())
	if err != nil {
		t.Skipf("skipping: could not list maintenances of app service %s: %v", globalAppServiceId, err)
	}
	if maintenance == nil {
		t.Skip("no pending maintenance queued at least 30 minutes ahead, with room before reschedule_by, on app service " + globalAppServiceId)
	}

	resourceName := randomStringWithPrefix("tf_acc_app_service_maintenance_schedule_")
	resourceReference := "couchbase-capella_app_service_maintenance_schedule." + resourceName
	currentQueueAt := maintenance.QueueAt.UTC().Format(time.RFC3339)
	// Halfway to reschedule_by keeps the new time inside the scheduling window.
	laterQueueAt := maintenance.QueueAt.Add(maintenance.RescheduleBy.Sub(*maintenance.QueueAt) / 2).UTC().Truncate(time.Second).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceMaintenanceStillExists(maintenance.Id),
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccAppServiceMaintenanceScheduleConfig(resourceName, maintenance.Id, currentQueueAt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "maintenance_id", maintenance.Id),
					resource.TestCheckResourceAttr(resourceReference, "queue_at", currentQueueAt),
					resource.TestCheckResourceAttr(resourceReference, "initiated_by", maintenance.InitiatedBy),
					resource.TestCheckResourceAttrSet(resourceReference, "kind"),
					resource.TestCheckResourceAttrSet(resourceReference, "reschedule_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "current_state"),
				),
			},
			{
				Config: globalProviderBlock + testAccAppServiceMaintenanceScheduleConfig(resourceName, maintenance.Id, laterQueueAt),
				Check:  resource.TestCheckResourceAttr(resourceReference, "queue_at", laterQueueAt),
			},
			{
				// The maintenance is shared, so put it back where the test found it.
				Config: globalProviderBlock + testAccAppServiceMaintenanceScheduleConfig(resourceName, maintenance.Id, currentQueueAt),
				Check:  resource.TestCheckResourceAttr(resourceReference, "queue_at", currentQueueAt),
			},
			{
				ResourceName:                         resourceReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        appServiceMaintenanceScheduleImportId(maintenance.Id),
				ImportStateVerifyIdentifierAttribute: "maintenance_id",
				// A self-service maintenance moves reschedule_by on every read.
				ImportStateVerifyIgnore: []string{"reschedule_by"},
			},
		},
	})
}

func testAccAppServiceMaintenancesDataSourceConfig(dataSourceName string) string {
	return fmt.Sprintf(`
data "couchbase-capella_app_service_maintenances" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
}
`, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

func testAccAppServiceMaintenanceScheduleConfig(resourceName, maintenanceId, queueAt string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_app_service_maintenance_schedule" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  maintenance_id  = "%[6]s"
  queue_at        = "%[7]s"
}
`, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, maintenanceId, queueAt)
}

func appServiceMaintenanceScheduleImportId(maintenanceId string) string {
	return fmt.Sprintf("organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s,maintenance_id=%s",
		globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, maintenanceId)
}

func appServiceMaintenancesURL() string {
	return fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/maintenances",
		globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

func findReschedulableAppServiceMaintenance(ctx context.Context) (*appservice_maintenance.GetMaintenanceResponse, error) {
	cfg := api.EndpointCfg{Url: appServiceMaintenancesURL(), Method: http.MethodGet, SuccessStatus: http.StatusOK}
	maintenances, err := api.GetPaginated[[]appservice_maintenance.GetMaintenanceResponse](ctx, globalClient, globalToken, cfg, "")
	if err != nil {
		return nil, err
	}
	// The queue time must stay in the future while the test runs, and the window must leave room to
	// move it later. Only the get endpoint returns rescheduleBy, so each candidate is read on its own.
	earliest := time.Now().Add(30 * time.Minute)
	for _, listed := range maintenances {
		if listed.CurrentState != "pending" || listed.QueueAt == nil || !listed.QueueAt.After(earliest) {
			continue
		}
		cfg := api.EndpointCfg{Url: appServiceMaintenancesURL() + "/" + url.PathEscape(listed.Id), Method: http.MethodGet, SuccessStatus: http.StatusOK}
		response, err := globalClient.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err != nil {
			return nil, err
		}
		var m appservice_maintenance.GetMaintenanceResponse
		if err = json.Unmarshal(response.Body, &m); err != nil {
			return nil, err
		}
		if m.RescheduleBy != nil && m.QueueAt != nil && m.RescheduleBy.Sub(*m.QueueAt) >= 2*time.Minute {
			return &m, nil
		}
	}
	return nil, nil
}

// testAccCheckAppServiceMaintenanceStillExists asserts destroy left the maintenance record in place.
func testAccCheckAppServiceMaintenanceStillExists(maintenanceId string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		cfg := api.EndpointCfg{Url: appServiceMaintenancesURL() + "/" + maintenanceId, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		if _, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil); err != nil {
			return fmt.Errorf("maintenance %s must survive destroy: %w", maintenanceId, err)
		}
		return nil
	}
}
