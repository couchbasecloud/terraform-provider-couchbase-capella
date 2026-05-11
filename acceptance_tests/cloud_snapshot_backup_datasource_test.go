package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudSnapshotBackupDatasource(t *testing.T) {
	clusterID, err := ensureSnapshotCluster()
	if err != nil {
		t.Fatalf("ensureSnapshotCluster: %v", err)
	}
	backupResourceName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_")
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_ds_")
	dsReference := "data.couchbase-capella_cloud_snapshot_backup." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSnapshotBackupSingleDatasourceConfig(clusterID, backupResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", clusterID),
					resource.TestCheckResourceAttrSet(dsReference, "id"),
					resource.TestCheckResourceAttrSet(dsReference, "created_at"),
					resource.TestCheckResourceAttrSet(dsReference, "expiration"),
					resource.TestCheckResourceAttr(dsReference, "retention", "168"),
					resource.TestCheckResourceAttrSet(dsReference, "size"),
					resource.TestCheckResourceAttrSet(dsReference, "type"),
					resource.TestCheckResourceAttrSet(dsReference, "progress.status"),
				),
			},
		},
	})
}

func TestAccCloudSnapshotBackupsDatasource(t *testing.T) {
	clusterID, err := ensureSnapshotCluster()
	if err != nil {
		t.Fatalf("ensureSnapshotCluster: %v", err)
	}
	backupResourceName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_")
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_backups_ds_")
	dsReference := "data.couchbase-capella_cloud_snapshot_backups." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSnapshotBackupsListDatasourceConfig(clusterID, backupResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", clusterID),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.created_at"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.retention"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.type"),
				),
			},
		},
	})
}

func TestAccCloudProjectSnapshotBackupsDatasource(t *testing.T) {
	clusterID, err := ensureSnapshotCluster()
	if err != nil {
		t.Fatalf("ensureSnapshotCluster: %v", err)
	}
	backupResourceName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_")
	dsName := randomStringWithPrefix("tf_acc_cloud_project_snapshot_backups_ds_")
	dsReference := "data.couchbase-capella_cloud_project_snapshot_backups." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectSnapshotBackupsDatasourceConfig(clusterID, backupResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.created_at"),
				),
			},
		},
	})
}

func TestAccCloudSnapshotBackupDatasourceInvalidID(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_backup_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_snapshot_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  id              = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile("Snapshot Backup Not Found|Error Reading Capella Snapshot"),
			},
		},
	})
}

func TestAccCloudSnapshotBackupsDatasourceInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_backups_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_snapshot_backups" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile("Error Reading Capella Snapshot Backups"),
			},
		},
	})
}

func TestAccCloudSnapshotBackupsDatasourceInvalidFilter(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_snapshot_backups_ds_filter_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_snapshot_backups" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  filters {
    name = "status"
  }
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile("Invalid Filters Configuration|Both 'name' and 'values'"),
			},
		},
	})
}

func TestAccCloudProjectSnapshotBackupsDatasourceInvalidProject(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cloud_project_snapshot_backups_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_project_snapshot_backups" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId),
				ExpectError: regexp.MustCompile("Error Reading Capella Project Snapshot Backups"),
			},
		},
	})
}

func testAccCloudSnapshotBackupSingleDatasourceConfig(clusterID, backupResourceName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cloud_snapshot_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  retention       = 168
}

data "couchbase-capella_cloud_snapshot_backup" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  id              = couchbase-capella_cloud_snapshot_backup.%[2]s.id
}
`, globalProviderBlock, backupResourceName, globalOrgId, globalProjectId, clusterID, dsName)
}

func testAccCloudSnapshotBackupsListDatasourceConfig(clusterID, backupResourceName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cloud_snapshot_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  retention       = 168
}

data "couchbase-capella_cloud_snapshot_backups" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  depends_on      = [couchbase-capella_cloud_snapshot_backup.%[2]s]
}
`, globalProviderBlock, backupResourceName, globalOrgId, globalProjectId, clusterID, dsName)
}

func testAccCloudProjectSnapshotBackupsDatasourceConfig(clusterID, backupResourceName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cloud_snapshot_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  retention       = 168
}

data "couchbase-capella_cloud_project_snapshot_backups" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  depends_on      = [couchbase-capella_cloud_snapshot_backup.%[2]s]
}
`, globalProviderBlock, backupResourceName, globalOrgId, globalProjectId, clusterID, dsName)
}
