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
	resourceName := randomStringWithPrefix("tf_acc_backup_")
	resourceReference := "couchbase-capella_backup." + resourceName
	dsName := randomStringWithPrefix("tf_acc_backups_ds_")
	dsReference := "data.couchbase-capella_backups." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBackupResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "bucket_id", globalBucketId),
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
				Config: testAccBackupWithDatasourceConfig(resourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "bucket_id", globalBucketId),
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
				ExpectError: regexp.MustCompile("restore times must not be set while create backup"),
			},
		},
	})
}

func testAccBackupResourceConfig(resourceName string) string {
	return testAccBackupResourceConfigWithBucketID(resourceName, globalBucketId)
}

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

func testAccBackupWithDatasourceConfig(resourceName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
}

data "couchbase-capella_backups" "%[7]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  depends_on      = [couchbase-capella_backup.%[2]s]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId, dsName)
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
