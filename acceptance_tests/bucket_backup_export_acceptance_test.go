package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	backupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAccBucketBackupExportResource creates a bucket, backs it up and exports that backup.
//
// Everything runs on a bucket created per test run rather than the shared globalBucketId, for the
// same reason TestAccBackupResource does it: Capella serialises manual backups per bucket, and only
// one active export can exist per backup cycle, so a leaked export from a previously-killed run on
// a shared bucket would fail this test with a 409.
//
// The export is still pending when the apply returns - create deliberately does not wait for the
// archive - so the assertions cover the attributes the create and its single refresh populate, not
// the archive attributes that only appear once the export completes.
func TestAccBucketBackupExportResource(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_bucket_")
	bucketResourceReference := "couchbase-capella_bucket." + bucketResourceName

	backupResourceName := randomStringWithPrefix("tf_acc_export_backup_")
	backupResourceReference := "couchbase-capella_backup." + backupResourceName

	resourceName := randomStringWithPrefix("tf_acc_bucket_backup_export_")
	resourceReference := "couchbase-capella_bucket_backup_export." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: testAccBucketBackupExportResourceConfig(bucketResourceName, backupResourceName, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBucketBackupExportResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrPair(resourceReference, "bucket_id", bucketResourceReference, "id"),
					resource.TestCheckResourceAttrPair(resourceReference, "backup_id", backupResourceReference, "id"),
					resource.TestCheckResourceAttrPair(resourceReference, "bucket_name", bucketResourceReference, "name"),
					// The export belongs to the cycle of the backup it was created from.
					resource.TestCheckResourceAttrPair(resourceReference, "cycle_id", backupResourceReference, "cycle_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "created_at"),
					// Status is whatever the job has reached by the time of the refresh.
					resource.TestCheckResourceAttrSet(resourceReference, "status"),
				),
			},
			// ImportState testing. The resource does not support update, so Create / Read /
			// Import / Delete is the whole lifecycle.
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateBucketBackupExportImportIdForResource(resourceReference),
				ImportState:       true,
			},
		},
	})
}

// TestAccBucketBackupExportResourceDuplicateCycle asserts that a second export for a cycle that
// already has one is rejected with 409. depends_on forces the two creates to be ordered, so the
// second reliably loses rather than racing the first.
func TestAccBucketBackupExportResourceDuplicateCycle(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_dup_bucket_")
	backupResourceName := randomStringWithPrefix("tf_acc_export_dup_backup_")
	firstName := randomStringWithPrefix("tf_acc_export_dup_first_")
	secondName := randomStringWithPrefix("tf_acc_export_dup_second_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportDuplicateConfig(bucketResourceName, backupResourceName, firstName, secondName),
				// 14061 is "already has a pending export", 14062 "already has a complete export".
				// Which one comes back depends on how far the first export got.
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":1406[12]`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceUnknownBackup asserts that a backup that does not exist is
// reported as 404/5017 rather than being silently accepted.
func TestAccBucketBackupExportResourceUnknownBackup(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_unknown_backup_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
					"00000000-0000-0000-0000-000000000000",
				),
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":5017`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceUnknownCluster asserts the cluster in the path is checked before
// the backup, so a bad cluster is 404/4025 and not a backup-not-found.
func TestAccBucketBackupExportResourceUnknownCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_unknown_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId,
					"00000000-0000-0000-0000-000000000000", globalBucketId,
					"11111111-1111-1111-1111-111111111111",
				),
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":4025`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceWrongProject asserts that a project which does not own the
// cluster is rejected with 422/4031, rather than leaking whether the backup exists.
func TestAccBucketBackupExportResourceWrongProject(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_wrong_project_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, "00000000-0000-0000-0000-000000000000",
					globalClusterId, globalBucketId,
					"11111111-1111-1111-1111-111111111111",
				),
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":4031`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceMalformedBackupId asserts the schema validator rejects a
// non-UUID backup_id at plan time, before any API call is made.
func TestAccBucketBackupExportResourceMalformedBackupId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_bad_backup_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
					"not-a-uuid",
				),
				ExpectError: regexp.MustCompile(`(?s)Attribute backup_id must be a valid UUID`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceEmptyBucketId asserts bucket_id is rejected when empty. It is
// not a UUID - it is the base64 encoding of the bucket name - so it only carries a length
// validator.
func TestAccBucketBackupExportResourceEmptyBucketId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_empty_bucket_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId, globalClusterId, "",
					"11111111-1111-1111-1111-111111111111",
				),
				ExpectError: regexp.MustCompile(`(?s)Attribute bucket_id string length must be at least 1, got: 0`),
			},
		},
	})
}

// testAccBucketBackupExportResourceConfig declares a fresh bucket, a backup of it and an export of
// that backup. Terraform orders the three by their references, and the backup resource waits for
// the backup to reach a final state, so the export is only created once the backup is ready.
func testAccBucketBackupExportResourceConfig(bucketName, backupName, exportName string) string {
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
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, exportName)
}

// testAccBucketBackupExportDuplicateConfig exports the same backup twice. depends_on orders the
// second create after the first so the conflict is deterministic.
func testAccBucketBackupExportDuplicateConfig(bucketName, backupName, firstName, secondName string) string {
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

resource "couchbase-capella_bucket_backup_export" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
  backup_id       = couchbase-capella_backup.%[6]s.id
  depends_on      = [couchbase-capella_bucket_backup_export.%[7]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, firstName, secondName)
}

// testAccBucketBackupExportResourceConfigWithIDs is used by the invalid-input tests, which fail
// before an export is created and so need no bucket or backup of their own.
func testAccBucketBackupExportResourceConfigWithIDs(resourceName, organizationId, projectId, clusterId, bucketId, backupId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket_backup_export" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  backup_id       = "%[7]s"
}
`, globalProviderBlock, resourceName, organizationId, projectId, clusterId, bucketId, backupId)
}

func generateBucketBackupExportImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"id=%s,backup_id=%s,bucket_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["backup_id"], rawState["bucket_id"],
			rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func retrieveBucketBackupExportFromServer(data *providerschema.Data, organizationId, projectId, clusterId, bucketId, backupId, exportId string) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups/%s/exports/%s",
		data.HostURL, organizationId, projectId, clusterId, bucketId, backupId, exportId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return err
	}

	exportResp := backupapi.GetBucketBackupExportResponse{}
	if err := json.Unmarshal(response.Body, &exportResp); err != nil {
		return err
	}
	if exportResp.Id != exportId {
		return errors.ErrNotFound
	}
	return nil
}

func testAccExistsBucketBackupExportResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		data := newTestClient(t)
		return retrieveBucketBackupExportFromServer(
			data,
			rawState["organization_id"], rawState["project_id"], rawState["cluster_id"],
			rawState["bucket_id"], rawState["backup_id"], rawState["id"],
		)
	}
}
