package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceCloudSnapshotRestores(t *testing.T) {
	clusterID, err := ensureSnapshotCluster()
	if err != nil {
		t.Fatalf("ensureSnapshotCluster: %v", err)
	}
	backupResourceName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_")
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_restores_ds_")
	dsReference := "data.couchbase-capella_cloud_snapshot_restores." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: create snapshot backup (no restore yet).
			{
				Config: testAccCloudSnapshotRestoreBackupOnlyConfig(clusterID, backupResourceName),
			},
			// Step 2: trigger a restore by setting restore_times, then read the list.
			{
				Config: testAccCloudSnapshotRestoresDatasourceConfig(clusterID, backupResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", clusterID),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.snapshot"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.status"),
				),
			},
		},
	})
}

func TestAccDatasourceCloudSnapshotRestore(t *testing.T) {
	clusterID, err := ensureSnapshotCluster()
	if err != nil {
		t.Fatalf("ensureSnapshotCluster: %v", err)
	}
	backupResourceName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_")
	listDsName := randomStringWithPrefix("tf_acc_cloud_snapshot_restores_ds_")
	singleDsName := randomStringWithPrefix("tf_acc_cloud_snapshot_restore_ds_")
	singleDsReference := "data.couchbase-capella_cloud_snapshot_restore." + singleDsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSnapshotRestoreBackupOnlyConfig(clusterID, backupResourceName),
			},
			{
				Config: testAccCloudSnapshotRestoreSingleDatasourceConfig(clusterID, backupResourceName, listDsName, singleDsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(singleDsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(singleDsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(singleDsReference, "cluster_id", clusterID),
					resource.TestCheckResourceAttrSet(singleDsReference, "id"),
					resource.TestCheckResourceAttrSet(singleDsReference, "snapshot"),
					resource.TestCheckResourceAttrSet(singleDsReference, "status"),
				),
			},
		},
	})
}

func TestAccDatasourceCloudSnapshotRestoreInvalidID(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_restore_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_snapshot_restore" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  id              = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile("Snapshot Restore Not Found|Could not find snapshot restore"),
			},
		},
	})
}

func TestAccDatasourceCloudSnapshotRestoresInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_restores_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_snapshot_restores" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile("Error Reading Capella Snapshot Restores"),
			},
		},
	})
}

func TestAccDatasourceCloudSnapshotRestoresInvalidFilter(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_restores_filter_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_snapshot_restores" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  filter {
    name = "status"
  }
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile("Invalid Filters Configuration|Both 'name' and 'values'"),
			},
		},
	})
}

func testAccCloudSnapshotRestoreBackupOnlyConfig(clusterID, backupResourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cloud_snapshot_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  retention       = 168
}
`, globalProviderBlock, backupResourceName, globalOrgId, globalProjectId, clusterID)
}

func testAccCloudSnapshotRestoresDatasourceConfig(clusterID, backupResourceName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cloud_snapshot_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  retention       = 168
  restore_times   = 1
}

data "couchbase-capella_cloud_snapshot_restores" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  depends_on      = [couchbase-capella_cloud_snapshot_backup.%[2]s]
}
`, globalProviderBlock, backupResourceName, globalOrgId, globalProjectId, clusterID, dsName)
}

func testAccCloudSnapshotRestoreSingleDatasourceConfig(clusterID, backupResourceName, listDsName, singleDsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cloud_snapshot_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  retention       = 168
  restore_times   = 1
}

data "couchbase-capella_cloud_snapshot_restores" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  depends_on      = [couchbase-capella_cloud_snapshot_backup.%[2]s]
}

data "couchbase-capella_cloud_snapshot_restore" "%[7]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  id              = data.couchbase-capella_cloud_snapshot_restores.%[6]s.data[0].id
}
`, globalProviderBlock, backupResourceName, globalOrgId, globalProjectId, clusterID, listDsName, singleDsName)
}
