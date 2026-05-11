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
	scheduleapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup_schedule"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TODO: legacy bucket backup schedule endpoint returns 404 "Unable to find the
// specified bucket" for a bucket that the bucket endpoint resolves successfully.
// Track via the bug filed for couchbase-capella_backup_schedule.
func TestAccBackupScheduleResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_schedule_")
	resourceReference := "couchbase-capella_backup_schedule." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupScheduleResourceConfig(resourceName, "weekly", "sunday", 10, 4, "30days", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBackupScheduleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "bucket_id", globalBucketId),
					resource.TestCheckResourceAttr(resourceReference, "type", "weekly"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.day_of_week", "sunday"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.start_at", "10"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.incremental_every", "4"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.retention_time", "30days"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.cost_optimized_retention", "false"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateBackupScheduleImportIdForResource(resourceReference),
				ImportState:       true,
			},
			{
				Config: testAccBackupScheduleResourceConfig(resourceName, "weekly", "monday", 14, 6, "60days", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBackupScheduleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "bucket_id", globalBucketId),
					resource.TestCheckResourceAttr(resourceReference, "type", "weekly"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.day_of_week", "monday"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.start_at", "14"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.incremental_every", "6"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.retention_time", "60days"),
					resource.TestCheckResourceAttr(resourceReference, "weekly_schedule.cost_optimized_retention", "true"),
				),
			},
		},
	})
}

func TestAccBackupScheduleResourceInvalidDayOfWeek(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_schedule_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupScheduleResourceConfig(resourceName, "weekly", "funday", 10, 4, "30days", false),
				ExpectError: regexp.MustCompile("There is an error during backup schedule creation"),
			},
		},
	})
}

func TestAccBackupScheduleResourceInvalidRetention(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_schedule_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupScheduleResourceConfig(resourceName, "weekly", "sunday", 10, 4, "9999days", false),
				ExpectError: regexp.MustCompile("There is an error during backup schedule creation"),
			},
		},
	})
}

func TestAccBackupScheduleResourceInvalidType(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_schedule_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupScheduleResourceConfig(resourceName, "yearly", "sunday", 10, 4, "30days", false),
				ExpectError: regexp.MustCompile("There is an error during backup schedule creation"),
			},
		},
	})
}

func TestAccBackupScheduleResourceInvalidStartAt(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_schedule_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupScheduleResourceConfig(resourceName, "weekly", "sunday", 99, 4, "30days", false),
				ExpectError: regexp.MustCompile("There is an error during backup schedule creation"),
			},
		},
	})
}

func TestAccBackupScheduleResourceInvalidBucket(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_schedule_invalid_bkt_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "00000000-0000-0000-0000-000000000000"
  type            = "weekly"
  weekly_schedule = {
    day_of_week              = "sunday"
    start_at                 = 10
    incremental_every        = 4
    retention_time           = "30days"
    cost_optimized_retention = false
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile("There is an error during backup schedule creation"),
			},
		},
	})
}

func testAccBackupScheduleResourceConfig(resourceName, scheduleType, dayOfWeek string, startAt, incrementalEvery int, retentionTime string, costOptimized bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_backup_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  type            = "%[7]s"
  weekly_schedule = {
    day_of_week              = "%[8]s"
    start_at                 = %[9]d
    incremental_every        = %[10]d
    retention_time           = "%[11]s"
    cost_optimized_retention = %[12]t
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
		scheduleType, dayOfWeek, startAt, incrementalEvery, retentionTime, costOptimized)
}

func generateBackupScheduleImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("organization_id=%s,project_id=%s,cluster_id=%s,bucket_id=%s",
			rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["bucket_id"]), nil
	}
}

func retrieveBackupScheduleFromServer(data *providerschema.Data, organizationId, projectId, clusterId, bucketId string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backup/schedules",
		data.HostURL, organizationId, projectId, clusterId, bucketId)
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
	scheduleResp := scheduleapi.GetBackupScheduleResponse{}
	if err := json.Unmarshal(response.Body, &scheduleResp); err != nil {
		return err
	}
	if scheduleResp.WeeklySchedule == nil {
		return errors.ErrNotFound
	}
	return nil
}

func testAccExistsBackupScheduleResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
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
		return retrieveBackupScheduleFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["bucket_id"])
	}
}
