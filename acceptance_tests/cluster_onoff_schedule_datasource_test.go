package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceClusterOnOffSchedule(t *testing.T) {
	scheduleResourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_")
	dsName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_ds_")
	dsReference := "data.couchbase-capella_cluster_onoff_schedule." + dsName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterOnOffScheduleResourceAndDatasourceConfig(scheduleResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "timezone", "US/Pacific"),
					resource.TestCheckResourceAttr(dsReference, "days.#", "7"),
					resource.TestCheckResourceAttrSet(dsReference, "days.0.day"),
					resource.TestCheckResourceAttr(dsReference, "days.0.state", "on"),
				),
			},
		},
	})
}

func TestAccDatasourceClusterOnOffScheduleInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error.*[Ss]chedule|cluster.*not found|access to the requested resource is denied|on/off schedule`),
			},
		},
	})
}

func testAccClusterOnOffScheduleResourceAndDatasourceConfig(resName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  timezone        = "US/Pacific"
  days = [
    { day = "monday",    state = "on" },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}

data "couchbase-capella_cluster_onoff_schedule" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"

  depends_on = [couchbase-capella_cluster_onoff_schedule.%[2]s]
}
`, globalProviderBlock, resName, dsName, globalOrgId, globalProjectId, globalClusterId)
}
