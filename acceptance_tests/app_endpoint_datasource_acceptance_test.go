package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppEndpointDataSource verifies the couchbase-capella_app_endpoint data
// source (single endpoint read) against the common pre-created endpoint.
func TestAccAppEndpointDataSource(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_ds_app_endpoint_")
	dataSourceReference := "data.couchbase-capella_app_endpoint." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointDataSourceConfig(dataSourceName, globalAppEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(dataSourceReference, "name", globalAppEndpointName),
					resource.TestCheckResourceAttrSet(dataSourceReference, "bucket"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "state"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "delta_sync_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "scopes.%"),
				),
			},
		},
	})
}

// TestAccAppEndpointActivationStatusDataSource verifies the
// couchbase-capella_app_endpoint_activation_status data source returns the
// endpoint state for the common pre-created endpoint.
func TestAccAppEndpointActivationStatusDataSource(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_ds_activation_status_")
	dataSourceReference := "data.couchbase-capella_app_endpoint_activation_status." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointActivationStatusDataSourceConfig(dataSourceName, globalAppEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttrSet(dataSourceReference, "state"),
				),
			},
		},
	})
}

// TestAccAppEndpointsDataSource verifies the couchbase-capella_app_endpoints
// (list) data source returns at least the common pre-created endpoint.
// Skipped: the unfiltered list call returns 500 "couldn't load database: bucket
// not found" when parallel tests have endpoints in mid-creation. Use
// TestAccAppEndpointsDataSourceFiltered for coverage until the App Service list
// API is stable under concurrent load.
func TestAccAppEndpointsDataSource(t *testing.T) {
	t.Skip("AV-130079: unfiltered app_endpoints list returns 500 when other endpoints are mid-creation; use TestAccAppEndpointsDataSourceFiltered for coverage")

	dataSourceName := randomStringWithPrefix("tf_acc_ds_app_endpoints_")
	dataSourceReference := "data.couchbase-capella_app_endpoints." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointsDataSourceConfig(dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "app_endpoints.#"),
				),
			},
		},
	})
}

// TestAccAppEndpointsDataSourceFiltered verifies that the filter block on the
// couchbase-capella_app_endpoints data source returns only the named endpoint.
func TestAccAppEndpointsDataSourceFiltered(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_ds_app_endpoints_filtered_")
	dataSourceReference := "data.couchbase-capella_app_endpoints." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointsDataSourceFilteredConfig(dataSourceName, globalAppEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_endpoints.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceReference, "app_endpoints.*", map[string]string{
						"bucket":             globalBucketName,
						"name":               globalAppEndpointName,
						"delta_sync_enabled": "true",
					}),
				),
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccAppEndpointDataSourceConfig(dataSourceName, endpointName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_endpoint" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  name            = "%[7]s"
}
`,
		globalProviderBlock,
		dataSourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		endpointName,
	)
}

func testAccAppEndpointActivationStatusDataSourceConfig(dataSourceName, endpointName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_endpoint_activation_status" "%[2]s" {
  organization_id   = "%[3]s"
  project_id        = "%[4]s"
  cluster_id        = "%[5]s"
  app_service_id    = "%[6]s"
  app_endpoint_name = "%[7]s"
}
`,
		globalProviderBlock,
		dataSourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		endpointName,
	)
}

func testAccAppEndpointsDataSourceConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_endpoints" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}
`,
		globalProviderBlock,
		dataSourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
	)
}

func testAccAppEndpointsDataSourceFilteredConfig(dataSourceName, endpointName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_endpoints" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"

  filter {
    name   = "name"
    values = ["%[7]s"]
  }
}
`,
		globalProviderBlock,
		dataSourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		endpointName,
	)
}
