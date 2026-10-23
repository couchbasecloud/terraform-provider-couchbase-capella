package acceptance_tests

import (
	"fmt"
	"regexp"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/enums"
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAccReadClusterDatasource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_ds_")
	resourceReference := "data.couchbase-capella_cluster." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDatasourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "name"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "id", globalClusterId),
					resource.TestCheckResourceAttrSet(resourceReference, "cloud_provider.type"),
					resource.TestCheckResourceAttrSet(resourceReference, "current_state"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
		},
	})
}

func TestAccReadClusterDatasource_InvalidID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_ds_invalid_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterDatasourceConfigWithID(resourceName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile("Error Reading Capella Cluster|Could not read cluster"),
			},
		},
	})
}

// TestAccReadClusterDatasourceExpressScaling checks that the cluster and clusters data sources
// report express scaling for the global cluster, which has never had it enabled.
func TestAccReadClusterDatasourceExpressScaling(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_ds_express_scaling_")
	clusterReference := "data.couchbase-capella_cluster." + resourceName
	clustersReference := "data.couchbase-capella_clusters." + resourceName
	disabled := string(clusterapi.ExpressScalingDisabled)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { requireExpressScalingDisabled(t, globalClusterId) },
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDatasourceExpressScalingConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterReference, "express_scaling", "false"),
					resource.TestCheckResourceAttrWith(clusterReference, "express_scaling_state", checkExpressScalingEnumValue),
					resource.TestCheckResourceAttr(clusterReference, "express_scaling_state", disabled),
					resource.TestCheckTypeSetElemNestedAttrs(clustersReference, "data.*", map[string]string{
						"id":                    globalClusterId,
						"express_scaling":       "false",
						"express_scaling_state": disabled,
					}),
				),
			},
		},
	})
}

// checkExpressScalingEnumValue fails unless value is one of the express scaling states the
// published OpenAPI spec allows on a cluster read. The spec's expressScaling field is the state,
// which the provider exposes as express_scaling_state.
func checkExpressScalingEnumValue(value string) error {
	def := enums.Lookup(capellaschema.NewSchemaBuilder("cluster"), []string{"GetClusterResponse"}, "express_scaling")
	if def == nil {
		return fmt.Errorf("no generated enum for GetClusterResponse.expressScaling")
	}
	if !slices.Contains(def.Values, value) {
		return fmt.Errorf("express_scaling %q is not one of %v", value, def.Values)
	}
	return nil
}

func testAccClusterDatasourceExpressScalingConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  id              = "%[4]s"
}

data "couchbase-capella_clusters" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccClusterDatasourceConfig(resourceName string) string {
	return testAccClusterDatasourceConfigWithID(resourceName, globalClusterId)
}

func testAccClusterDatasourceConfigWithID(resourceName, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  id              = "%[4]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterID, resourceName)
}
