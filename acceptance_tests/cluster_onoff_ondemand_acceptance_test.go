package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterOnOffOnDemandResourceInvalidState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_ondemand_invalid_state_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterOnOffOnDemandResourceConfig(resourceName, globalClusterId, "frozen"),
				ExpectError: regexp.MustCompile(`(?s)state|invalid value|Validation Error|on or off`),
			},
		},
	})
}

func TestAccClusterOnOffOnDemandResourceInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_ondemand_invalid_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterOnOffOnDemandResourceConfig(resourceName, "00000000-0000-0000-0000-000000000000", "on"),
				ExpectError: regexp.MustCompile(`(?s)Error.*[Cc]luster|cluster.*not found|access to the requested resource is denied|switching on/off`),
			},
		},
	})
}

func TestAccClusterOnOffOnDemandResourceMissingCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_ondemand_missing_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_ondemand" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  state           = "on"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)cluster_id|argument.*required`),
			},
		},
	})
}

func TestAccClusterOnOffOnDemandResourceMissingState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_ondemand_missing_state_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_ondemand" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)state|argument.*required`),
			},
		},
	})
}

func testAccClusterOnOffOnDemandResourceConfig(resourceName, clusterID, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_ondemand" "%[2]s" {
  organization_id            = "%[3]s"
  project_id                 = "%[4]s"
  cluster_id                 = "%[5]s"
  state                      = "%[6]s"
  turn_on_linked_app_service = false
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, clusterID, state)
}
