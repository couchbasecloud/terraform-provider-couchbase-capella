package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceBucketBackupExport creates an export and reads it back through the data source,
// asserting the two agree field for field.
//
// It runs on a bucket and backup created per test run for the same reason the resource test does:
// only one active export can exist per backup cycle, so a leaked export on a shared bucket would
// fail the create with a 409.
//
// The archive attributes - size_in_bytes, sha256_checksum, expiration, backup_download_url - are
// deliberately not asserted. The export is still pending when the apply returns, so they are null
// at this point; asserting them would make the test depend on how fast the backup infrastructure
// happens to be.
func TestAccDatasourceBucketBackupExport(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_ds_bucket_")
	backupResourceName := randomStringWithPrefix("tf_acc_export_ds_backup_")
	exportResourceName := randomStringWithPrefix("tf_acc_export_ds_export_")
	dsName := randomStringWithPrefix("tf_acc_export_ds_")

	exportReference := "couchbase-capella_bucket_backup_export." + exportResourceName
	dsReference := "data.couchbase-capella_bucket_backup_export." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportDatasourceConfig(bucketResourceName, backupResourceName, exportResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					// The data source must resolve to the same export the resource created.
					resource.TestCheckResourceAttrPair(dsReference, "id", exportReference, "id"),
					resource.TestCheckResourceAttrPair(dsReference, "bucket_id", exportReference, "bucket_id"),
					resource.TestCheckResourceAttrPair(dsReference, "backup_id", exportReference, "backup_id"),
					resource.TestCheckResourceAttrPair(dsReference, "cycle_id", exportReference, "cycle_id"),
					resource.TestCheckResourceAttrPair(dsReference, "bucket_name", exportReference, "bucket_name"),
					resource.TestCheckResourceAttrPair(dsReference, "created_at", exportReference, "created_at"),
					resource.TestCheckResourceAttrSet(dsReference, "status"),
				),
			},
		},
	})
}

// TestAccDatasourceBucketBackupExportUnknownExport asserts an export that does not exist is
// reported as 404/5019.
func TestAccDatasourceBucketBackupExportUnknownExport(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_export_ds_unknown_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportDatasourceConfigWithIDs(
					dsName, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
					"00000000-0000-0000-0000-000000000000",
					"11111111-1111-1111-1111-111111111111",
				),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Bucket Backup Export.*"code":5019`),
			},
		},
	})
}

// TestAccDatasourceBucketBackupExportWrongProject asserts a project that does not own the cluster
// is rejected with 422/4031 rather than leaking whether the export exists.
func TestAccDatasourceBucketBackupExportWrongProject(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_export_ds_wrong_project_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportDatasourceConfigWithIDs(
					dsName, globalOrgId, "00000000-0000-0000-0000-000000000000",
					globalClusterId, globalBucketId,
					"11111111-1111-1111-1111-111111111111",
					"22222222-2222-2222-2222-222222222222",
				),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Bucket Backup Export.*"code":4031`),
			},
		},
	})
}

// TestAccDatasourceBucketBackupExportMalformedId asserts the schema validator rejects a non-UUID
// id at plan time, before any API call is made.
func TestAccDatasourceBucketBackupExportMalformedId(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_export_ds_bad_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportDatasourceConfigWithIDs(
					dsName, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
					"11111111-1111-1111-1111-111111111111",
					"not-a-uuid",
				),
				ExpectError: regexp.MustCompile(`(?s)Attribute id must be a valid UUID`),
			},
		},
	})
}

// testAccBucketBackupExportDatasourceConfig declares a bucket, a backup, an export of that backup
// and a data source reading the export back. Terraform orders them by their references.
func testAccBucketBackupExportDatasourceConfig(bucketName, backupName, exportName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
}

resource "couchbase-capella_backup" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
}

resource "couchbase-capella_bucket_backup_export" "%[7]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
  backup_id       = couchbase-capella_backup.%[6]s.id
}

data "couchbase-capella_bucket_backup_export" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
  backup_id       = couchbase-capella_backup.%[6]s.id
  id              = couchbase-capella_bucket_backup_export.%[7]s.id
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, exportName, dsName)
}

// testAccBucketBackupExportDatasourceConfigWithIDs is used by the invalid-input tests, which fail
// on the read and so need no bucket, backup or export of their own.
func testAccBucketBackupExportDatasourceConfigWithIDs(dsName, organizationId, projectId, clusterId, bucketId, backupId, exportId string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_bucket_backup_export" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  backup_id       = "%[7]s"
  id              = "%[8]s"
}
`, globalProviderBlock, dsName, organizationId, projectId, clusterId, bucketId, backupId, exportId)
}
