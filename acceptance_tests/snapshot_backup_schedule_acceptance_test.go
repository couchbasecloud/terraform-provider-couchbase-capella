package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	snapshotbackupscheduleapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/snapshot_backup_schedule"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAccSnapshotBackupScheduleResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_snapshot_backup_schedule_")
	resourceReference := "couchbase-capella_snapshot_backup_schedule." + resourceName

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour).Format(time.RFC3339)
	copyToRegions := "[" + strings.Join([]string{"\"eu-west-1\"", "\"ap-southeast-1\""}, ",") + "]"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotBackupScheduleResourceConfigWithCopyToRegions(resourceName, 12, 240, startTime, copyToRegions),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsSnapshotBackupScheduleResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "interval", "12"),
					resource.TestCheckResourceAttr(resourceReference, "retention", "240"),
					resource.TestCheckResourceAttr(resourceReference, "start_time", startTime),
					resource.TestCheckResourceAttr(resourceReference, "copy_to_regions.0", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceReference, "copy_to_regions.1", "ap-southeast-1"),
				),
			},
		},
	})
}

func TestAccSnapshotBackupScheduleResourceInvalidInterval(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_snapshot_backup_schedule_")

	startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnapshotBackupScheduleResourceConfigWithCopyToRegions(resourceName, 0, 240, startTime, "[]"),
				ExpectError: regexp.MustCompile("There is an error during snapshot backup schedule creation"),
			},
		},
	})
}

func TestAccSnapshotBackupScheduleResourceInvalidRetention(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_snapshot_backup_schedule_")

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnapshotBackupScheduleResourceConfigWithCopyToRegions(resourceName, 12, 721, startTime, "[]"),
				ExpectError: regexp.MustCompile("There is an error during snapshot backup schedule creation"),
			},
		},
	})
}

func TestAccSnapshotBackupScheduleResourceInvalidStartTime(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_snapshot_backup_schedule_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnapshotBackupScheduleResourceConfigWithCopyToRegions(resourceName, 12, 240, "invalid_time", "[]"),
				ExpectError: regexp.MustCompile("There is an error during snapshot backup schedule creation"),
			},
		},
	})
}

func TestAccSnapshotBackupScheduleResourceInvalidCopyToRegions(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_snapshot_backup_schedule_")

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour).Format(time.RFC3339)

	copyToRegions := "[" + strings.Join([]string{"\"us-east-1\""}, ",") + "]"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnapshotBackupScheduleResourceConfigWithCopyToRegions(resourceName, 12, 240, startTime, copyToRegions),
				ExpectError: regexp.MustCompile("There is an error during snapshot backup schedule creation"),
			},
		},
	})
}

func testAccSnapshotBackupScheduleResourceConfigWithCopyToRegions(resourceName string, interval, retention int, startTime, copyToRegions string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_snapshot_backup_schedule" "%[2]s" {
		organization_id = "%[3]s"
		project_id = "%[4]s"
		id = "%[5]s"
		interval = %[6]d
		retention = %[7]d
		start_time = "%[8]s"
		copy_to_regions = %[9]s
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, interval, retention, startTime, copyToRegions)
}

func testAccExistsSnapshotBackupScheduleResource(resourceReference string) resource.TestCheckFunc {
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
		err := retrieveSnapshotBackupScheduleFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
		if err != nil {
			return err
		}
		return nil
	}
}

func retrieveSnapshotBackupScheduleFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/cloudsnapshotbackupschedule", data.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	backupScheduleResp, err := data.Client.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)

	if err != nil {
		return err
	}

	snapshotBackupSchedule := snapshotbackupscheduleapi.SnapshotBackupSchedule{}
	err = json.Unmarshal(backupScheduleResp.Body, &snapshotBackupSchedule)
	if err != nil {
		return err
	}
	if snapshotBackupSchedule.Interval == 0 && snapshotBackupSchedule.Retention == 0 && snapshotBackupSchedule.StartTime == "" {
		return errors.ErrNotFound
	}
	return nil
}
