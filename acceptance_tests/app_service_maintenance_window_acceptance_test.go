package acceptance_tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// TestAccAppServiceMaintenanceWindow runs serially: an App Service has at most one window,
// so every step works against the same singleton on the shared App Service.
func TestAccAppServiceMaintenanceWindow(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_maintenance_window_")
	resourceReference := "couchbase-capella_app_service_maintenance_window." + resourceName
	dataSourceReference := "data.couchbase-capella_app_service_maintenance_window." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		PreCheck: func() {
			if err := deleteAppServiceMaintenanceWindow(context.Background()); err != nil {
				t.Fatalf("could not clear the maintenance window of app service %s: %s", globalAppServiceId, err)
			}
		},
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccAppServiceMaintenanceWindowResourceConfig(resourceName, "sunday", 2, 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "preferred_day", "sunday"),
					resource.TestCheckResourceAttr(resourceReference, "preferred_hour", "2"),
					resource.TestCheckResourceAttr(resourceReference, "preferred_minute", "30"),
				),
			},
			{
				Config: globalProviderBlock +
					testAccAppServiceMaintenanceWindowResourceConfig(resourceName, "wednesday", 23, 59) +
					testAccAppServiceMaintenanceWindowDataSourceConfig(resourceName, resourceReference),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "preferred_day", "wednesday"),
					resource.TestCheckResourceAttr(resourceReference, "preferred_hour", "23"),
					resource.TestCheckResourceAttr(resourceReference, "preferred_minute", "59"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "preferred_day", resourceReference, "preferred_day"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "preferred_hour", resourceReference, "preferred_hour"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "preferred_minute", resourceReference, "preferred_minute"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        appServiceMaintenanceWindowImportId(),
				ImportStateVerifyIdentifierAttribute: "app_service_id",
			},
			{
				// A second resource on the same App Service must not overwrite the first window.
				Config: globalProviderBlock +
					testAccAppServiceMaintenanceWindowResourceConfig(resourceName, "wednesday", 23, 59) +
					testAccAppServiceMaintenanceWindowResourceConfig(resourceName+"_duplicate", "monday", 1, 0),
				ExpectError: regexp.MustCompile(`(?s)already exists.*terraform import`),
			},
			{
				Config: globalProviderBlock,
			},
			{
				Config: globalProviderBlock + testAccAppServiceMaintenanceWindowDataSourceConfig(resourceName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceReference, "preferred_day"),
					resource.TestCheckNoResourceAttr(dataSourceReference, "preferred_hour"),
					resource.TestCheckNoResourceAttr(dataSourceReference, "preferred_minute"),
				),
			},
		},
	})
}

func TestAccAppServiceMaintenanceWindowInvalidValues(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_maintenance_window_invalid_")

	cases := []struct {
		day          string
		hour, minute int
		expectError  string
	}{
		{"Sunday", 2, 30, `preferred_day value must be one of`},
		{"sunday", 24, 30, `preferred_hour value must be between 0 and 23`},
		{"sunday", 2, 60, `preferred_minute value must be between 0 and 59`},
	}

	steps := make([]resource.TestStep, 0, len(cases))
	for _, c := range cases {
		steps = append(steps, resource.TestStep{
			Config:      globalProviderBlock + testAccAppServiceMaintenanceWindowResourceConfig(resourceName, c.day, c.hour, c.minute),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(c.expectError),
		})
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

func testAccAppServiceMaintenanceWindowResourceConfig(resourceName, day string, hour, minute int) string {
	return fmt.Sprintf(`
resource "couchbase-capella_app_service_maintenance_window" "%[1]s" {
  organization_id  = "%[2]s"
  project_id       = "%[3]s"
  cluster_id       = "%[4]s"
  app_service_id   = "%[5]s"
  preferred_day    = "%[6]s"
  preferred_hour   = %[7]d
  preferred_minute = %[8]d
}
`, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, day, hour, minute)
}

// testAccAppServiceMaintenanceWindowDataSourceConfig reads the window after dependsOn, when given,
// so the data source sees the resource's latest write.
func testAccAppServiceMaintenanceWindowDataSourceConfig(dataSourceName, dependsOn string) string {
	dependsOnLine := ""
	if dependsOn != "" {
		dependsOnLine = fmt.Sprintf("  depends_on       = [%s]\n", dependsOn)
	}
	return fmt.Sprintf(`
data "couchbase-capella_app_service_maintenance_window" "%[1]s" {
  organization_id  = "%[2]s"
  project_id       = "%[3]s"
  cluster_id       = "%[4]s"
  app_service_id   = "%[5]s"
%[6]s}
`, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, dependsOnLine)
}

func appServiceMaintenanceWindowImportId() string {
	return fmt.Sprintf("organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s",
		globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}

// deleteAppServiceMaintenanceWindow clears any window already set on the shared App Service.
func deleteAppServiceMaintenanceWindow(ctx context.Context) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/maintenanceWindow",
		globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := globalClient.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	var apiErr *api.Error
	if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusNotFound {
		return nil
	}
	return err
}
