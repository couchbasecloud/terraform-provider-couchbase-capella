package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceCloudSnapshotBackupSchedule(t *testing.T) {
	scheduleResourceName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_schedule_")
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_schedule_ds_")
	dsReference := "data.couchbase-capella_cloud_snapshot_backup_schedule." + dsName

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSnapshotBackupScheduleDatasourceConfig(scheduleResourceName, dsName, 12, 240, startTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "interval", "12"),
					resource.TestCheckResourceAttr(dsReference, "retention", "240"),
					resource.TestCheckResourceAttr(dsReference, "start_time", startTime),
					resource.TestCheckResourceAttr(dsReference, "copy_to_regions.#", "0"),
				),
			},
		},
	})
}

func TestAccDatasourceCloudSnapshotBackupScheduleInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_schedule_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_snapshot_backup_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile("Error Reading Capella Snapshot Backup Schedule"),
			},
		},
	})
}

func testAccCloudSnapshotBackupScheduleDatasourceConfig(scheduleResourceName, dsName string, interval, retention int, startTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cloud_snapshot_backup_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  interval        = %[6]d
  retention       = %[7]d
  start_time      = "%[8]s"
  copy_to_regions = null
}

data "couchbase-capella_cloud_snapshot_backup_schedule" "%[9]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  depends_on      = [couchbase-capella_cloud_snapshot_backup_schedule.%[2]s]
}
`, globalProviderBlock, scheduleResourceName, globalOrgId, globalProjectId, globalClusterId,
		interval, retention, startTime, dsName)
}
