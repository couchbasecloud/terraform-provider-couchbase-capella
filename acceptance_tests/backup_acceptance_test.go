package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	backupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAccBackupResource creates a backup, verifies its attributes, then
// layers the couchbase-capella_backups data source on top to confirm the
// backups list returns the just-created backup. The data source step
// intentionally reuses the same backup resource (no new backup is created)
// because Capella's legacy bucket backup endpoint serialises manual
// backups per bucket and a back-to-back second backup gets stuck in a
// non-terminal state until the per-bucket spacing window elapses.
func TestAccBackupResource(t *testing.T) {
	// Run on a fresh bucket created per test run rather than the shared
	// globalBucketId. Capella's legacy bucket-backup endpoint serialises
	// manual backups per bucket; on the shared CI tenant a leaked pending
	// backup from a previously-killed run can queue this test's POST
	// indefinitely. A dedicated bucket guarantees zero accumulated state.
	bucketResourceName := randomStringWithPrefix("tf_acc_backup_bucket_")
	bucketResourceReference := "couchbase-capella_bucket." + bucketResourceName

	resourceName := randomStringWithPrefix("tf_acc_backup_")
	resourceReference := "couchbase-capella_backup." + resourceName
	dsName := randomStringWithPrefix("tf_acc_backups_ds_")
	dsReference := "data.couchbase-capella_backups." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupOnIsolatedBucketConfig(bucketResourceName, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBackupResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrPair(resourceReference, "bucket_id", bucketResourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "cycle_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "date"),
					resource.TestCheckResourceAttrSet(resourceReference, "status"),
					resource.TestCheckResourceAttrSet(resourceReference, "method"),
					resource.TestCheckResourceAttrSet(resourceReference, "bucket_name"),
					resource.TestCheckResourceAttrSet(resourceReference, "source"),
					resource.TestCheckResourceAttrSet(resourceReference, "cloud_provider"),
				),
			},
			{
				Config: testAccBackupOnIsolatedBucketWithDatasourceConfig(bucketResourceName, resourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrPair(dsReference, "bucket_id", bucketResourceReference, "id"),
					testAccCheckDataSourceContainsBackup(dsReference, resourceReference),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateBackupImportIdForResource(resourceReference),
				ImportState:       true,
			},
		},
	})
}

func TestAccBackupResourceInvalidBucket(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupResourceConfigWithBucketID(resourceName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile("Error getting latest bucket backup|There is an error during backup creation"),
			},
		},
	})
}

func TestAccBackupResourceInvalidProject(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_invalid_proj_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "00000000-0000-0000-0000-000000000000"
  cluster_id      = "%[4]s"
  bucket_id       = "%[5]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalClusterId, globalBucketId),
				ExpectError: regexp.MustCompile("Error getting latest bucket backup|There is an error during backup creation"),
			},
		},
	})
}

func TestAccBackupResourceRestoreTimesOnCreate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_invalid_rt_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  restore_times   = 1
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId),
				ExpectError: regexp.MustCompile("restore times must not be set during backup creation"),
			},
		},
	})
}

func TestAccBackupResourceRestoreOnCreate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_restore_on_create_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupResourceRestoreOnCreateConfig(resourceName),
				ExpectError: regexp.MustCompile(`restore must not be set during backup creation`),
			},
		},
	})
}

func TestAccBackupResourceRestoreInvalidReplaceTTL(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_invalid_replace_ttl_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupResourceInvalidRestoreReplaceTTLConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)replace_ttl.*(none|all|expired).*forever|forever.*(none|all|expired)`),
			},
		},
	})
}

func testAccBackupResourceRestoreOnCreateConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup" "%[2]s" {
	organization_id = "00000000-0000-0000-0000-000000000000"
	project_id      = "11111111-1111-1111-1111-111111111111"
	cluster_id      = "22222222-2222-2222-2222-222222222222"
	bucket_id       = "default"

	restore = {
		target_cluster_id       = "33333333-3333-3333-3333-333333333333"
		source_cluster_id       = "22222222-2222-2222-2222-222222222222"
		services                = ["data"]
		force_updates           = true
		auto_remove_collections = false
		replace_ttl             = "none"
		replace_ttl_with        = "0"
		include_data            = "bucket.scope.collection"
		exclude_data            = null
		filter_keys             = null
		filter_values           = null
		map_data                = null
	}
}
`, globalProviderBlock, resourceName)
}

// testAccBackupResourceConfigWithBucketID is used by the invalid-input tests,
// which short-circuit before any backup is actually created and therefore have
// no per-bucket-queue dependency — they continue to point at the shared bucket.
func testAccBackupResourceConfigWithBucketID(resourceName, bucketID string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, bucketID)
}

func testAccBackupResourceInvalidRestoreReplaceTTLConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup" "%[2]s" {
	organization_id = "00000000-0000-0000-0000-000000000000"
	project_id      = "11111111-1111-1111-1111-111111111111"
	cluster_id      = "22222222-2222-2222-2222-222222222222"
	bucket_id       = "default"

	restore = {
		target_cluster_id       = "33333333-3333-3333-3333-333333333333"
		source_cluster_id       = "22222222-2222-2222-2222-222222222222"
		services                = ["data"]
		force_updates           = true
		auto_remove_collections = false
		replace_ttl             = "forever"
		replace_ttl_with        = "0"
		include_data            = "bucket.scope.collection"
		exclude_data            = null
		filter_keys             = null
		filter_values           = null
		map_data                = null
	}
}
`, globalProviderBlock, resourceName)
}

// testAccBackupOnIsolatedBucketConfig declares a fresh bucket alongside the
// backup resource so the backup runs on a bucket with zero accumulated state.
// Terraform's destroy step cleans up both at end of test.
func testAccBackupOnIsolatedBucketConfig(bucketName, backupName string) string {
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
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName)
}

func testAccBackupOnIsolatedBucketWithDatasourceConfig(bucketName, backupName, dsName string) string {
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

data "couchbase-capella_backups" "%[7]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
  depends_on      = [couchbase-capella_backup.%[6]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, dsName)
}

func generateBackupImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}

func retrieveBackupFromServer(data *providerschema.Data, organizationId, projectId, clusterId, backupId string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/backups/%s",
		data.HostURL, organizationId, projectId, clusterId, backupId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return err
	}
	backupResp := backupapi.GetBackupResponse{}
	if err := json.Unmarshal(response.Body, &backupResp); err != nil {
		return err
	}
	if backupResp.Id == "" {
		return errors.ErrNotFound
	}
	return nil
}

func testAccExistsBackupResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
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
		return retrieveBackupFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["id"])
	}
}

// testAccCheckDataSourceContainsBackup verifies that the backups data source
// contains the specific backup created by resourceReference, regardless of
// its position in the list (older backups may appear at lower indices).
func testAccCheckDataSourceContainsBackup(dsReference, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		res := s.RootModule().Resources[resourceReference]
		if res == nil {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		expectedID := res.Primary.Attributes["id"]
		count, _ := strconv.Atoi(ds.Primary.Attributes["data.#"])
		for i := 0; i < count; i++ {
			if ds.Primary.Attributes[fmt.Sprintf("data.%d.id", i)] == expectedID {
				return nil
			}
		}
		return fmt.Errorf("datasource %s does not contain backup id=%s (checked %d items)", dsReference, expectedID, count)
	}
}
