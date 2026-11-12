package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
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

// TestAccAppServiceMaintenanceSchedule needs a pending, Couchbase-initiated maintenance on the shared
// App Service, which the provider cannot create; it skips when there is none.
func TestAccAppServiceMaintenanceSchedule(t *testing.T) {
	maintenance, err := findReschedulableAppServiceMaintenance(context.Background())
	if err != nil {
		t.Skipf("skipping: could not list maintenances of app service %s: %v", globalAppServiceId, err)
	}
	if maintenance == nil {
		t.Skip("no pending, Couchbase-initiated maintenance on app service " + globalAppServiceId)
	}

	resourceName := randomStringWithPrefix("tf_acc_app_service_maintenance_schedule_")
	resourceReference := "couchbase-capella_app_service_maintenance_schedule." + resourceName
	currentQueueAt := maintenance.QueueAt.UTC().Format(time.RFC3339)
	// The API exposes no scheduling-window end, so the step stays small to keep the new time inside it.
	laterQueueAt := maintenance.QueueAt.Add(5 * time.Minute).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceMaintenanceStillExists(maintenance.Id),
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccAppServiceMaintenanceScheduleConfig(resourceName, maintenance.Id, currentQueueAt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "maintenance_id", maintenance.Id),
					resource.TestCheckResourceAttr(resourceReference, "queue_at", currentQueueAt),
					resource.TestCheckResourceAttr(resourceReference, "initiated_by", "couchbase"),
					resource.TestCheckResourceAttrSet(resourceReference, "kind"),
					resource.TestCheckResourceAttrSet(resourceReference, "current_state"),
				),
			},
			{
				Config: globalProviderBlock + testAccAppServiceMaintenanceScheduleConfig(resourceName, maintenance.Id, laterQueueAt),
				Check:  resource.TestCheckResourceAttr(resourceReference, "queue_at", laterQueueAt),
			},
			{
				ResourceName:                         resourceReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        appServiceMaintenanceScheduleImportId(maintenance.Id),
				ImportStateVerifyIdentifierAttribute: "maintenance_id",
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
	for i := range maintenances {
		m := maintenances[i]
		if m.CurrentState == "pending" && m.InitiatedBy == "couchbase" && m.QueueAt != nil {
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
