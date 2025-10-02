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
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	snapshotbackupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

func TestAccSnapshotBackupResource(t *testing.T) {

	resourceName := randomStringWithPrefix("tf_acc_snapshot_backup_")
	resourceReference := "couchbase-capella_snapshot_backup." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSnapshotBackupResourceConfig(resourceName, 168),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsSnapshotBackupResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "retention", "168"),
					resource.TestCheckResourceAttr(resourceReference, "progress.status", "complete"),
					resource.TestCheckResourceAttrSet(resourceReference, "app_service"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "expiration"),
					resource.TestCheckResourceAttrSet(resourceReference, "progress.time"),
					resource.TestCheckResourceAttrSet(resourceReference, "server.version"),
					resource.TestCheckResourceAttrSet(resourceReference, "size"),
					resource.TestCheckResourceAttrSet(resourceReference, "type"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateSnapshotBackupImportIdForResource(resourceReference),
				ImportState:       true,
			},

			// Update retention testing
			{
				Config: testAccSnapshotBackupResourceConfig(resourceName, 240),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsSnapshotBackupResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "retention", "240"),
					resource.TestCheckResourceAttr(resourceReference, "progress.status", "complete"),
					resource.TestCheckResourceAttrSet(resourceReference, "app_service"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "expiration"),
					resource.TestCheckResourceAttrSet(resourceReference, "progress.time"),
					resource.TestCheckResourceAttrSet(resourceReference, "server.version"),
					resource.TestCheckResourceAttrSet(resourceReference, "size"),
					resource.TestCheckResourceAttrSet(resourceReference, "type"),
				),
			},
		},
	})
}

func TestAccSnapshotBackupResourceInvalidRetention(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_snapshot_backup_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnapshotBackupResourceConfig(resourceName, 23),
				ExpectError: regexp.MustCompile("There is an error during snapshot backup creation"),
			},
		},
	})
}

func testAccSnapshotBackupResourceConfig(resourceName string, retention int) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_snapshot_backup" "%[2]s" {
		organization_id = "%[3]s"
		project_id = "%[4]s"
		cluster_id = "%[5]s"
		retention = %[6]d
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, retention)
}

func generateSnapshotBackupImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}

func retrieveSnapshotBackupFromServer(data *providerschema.Data, organizationId, projectId, clusterId, id string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackups", data.HostURL, organizationId, projectId, clusterId)
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
	listSnapshotBackupsResp := snapshotbackupapi.ListSnapshotBackupsResponse{}
	err = json.Unmarshal(response.Body, &listSnapshotBackupsResp)
	if err != nil {
		return err
	}

	if len(listSnapshotBackupsResp.Data) > 0 {
		for _, backup := range listSnapshotBackupsResp.Data {
			if backup.ID == id {
				return nil
			}
		}
	}
	return errors.ErrNotFound
}

func testAccExistsSnapshotBackupResource(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		data := newTestClient()
		err := retrieveSnapshotBackupFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["id"])
		if err != nil {
			return err
		}
		return nil
	}
}
