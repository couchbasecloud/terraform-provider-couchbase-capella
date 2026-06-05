package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceClusterOnOffSchedule(t *testing.T) {
	scheduleResourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_")
	dsName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_ds_")
	dsReference := "data.couchbase-capella_cluster_onoff_schedule." + dsName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterOnOffScheduleResourceAndDatasourceConfig(scheduleResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "timezone", "US/Pacific"),
					resource.TestCheckResourceAttr(dsReference, "days.#", "7"),
					resource.TestCheckResourceAttrSet(dsReference, "days.0.day"),
					resource.TestCheckResourceAttrSet(dsReference, "days.0.state"),
				),
			},
		},
	})
}

func TestAccDatasourceClusterOnOffScheduleInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error.*[Ss]chedule|cluster.*not found|access to the requested resource is denied|on/off schedule`),
			},
		},
	})
}

// TestAccDatasourceClusterOnOffScheduleEmptyClusterId verifies that an empty
// cluster_id is rejected by local schema validation instead of being sent to
// the API.
func TestAccDatasourceClusterOnOffScheduleEmptyClusterId(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_empty_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterOnOffScheduleDatasourceIDsConfig(dsName, globalOrgId, globalProjectId, ""),
				ExpectError: regexp.MustCompile(
					`(?s)Invalid Attribute Value Length.*cluster_id.*string length must be at least 1`),
			},
		},
	})
}

// TestAccDatasourceClusterOnOffScheduleInvalidUUIDs verifies that each ID
// attribute rejects non-UUID values via local schema validation, one attribute
// per subtest.
func TestAccDatasourceClusterOnOffScheduleInvalidUUIDs(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		projectID      string
		clusterID      string
	}{
		{
			name:           "organization_id",
			organizationID: "not-a-uuid",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "project_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "not-a-uuid",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "cluster_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "not-a-uuid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dsName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_non_uuid_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccClusterOnOffScheduleDatasourceIDsConfig(
							dsName, test.organizationID, test.projectID, test.clusterID),
						ExpectError: regexp.MustCompile(
							`(?s)Invalid Attribute Value Match.*` + test.name + `.*must be a valid UUID`),
					},
				},
			})
		})
	}
}

func testAccClusterOnOffScheduleDatasourceIDsConfig(dsName, organizationID, projectID, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, organizationID, projectID, clusterID)
}

func testAccClusterOnOffScheduleResourceAndDatasourceConfig(resName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  timezone        = "US/Pacific"
  days = [
    {
      day   = "monday"
      state = "custom"
      from  = { hour = 0, minute = 0 }
      to    = { hour = 23, minute = 30 }
    },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}

data "couchbase-capella_cluster_onoff_schedule" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"

  depends_on = [couchbase-capella_cluster_onoff_schedule.%[2]s]
}
`, globalProviderBlock, resName, dsName, globalOrgId, globalProjectId, globalClusterId)
}
