package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterOnOffScheduleResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_")
	resourceReference := "couchbase-capella_cluster_onoff_schedule." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterOnOffScheduleResourceConfig(resourceName, "US/Pacific"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "timezone", "US/Pacific"),
					resource.TestCheckResourceAttr(resourceReference, "days.#", "7"),
				),
			},
			{
				Config: testAccClusterOnOffScheduleResourceUpdateConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "timezone", "US/Pacific"),
					resource.TestCheckResourceAttr(resourceReference, "days.#", "7"),
					resource.TestCheckResourceAttr(resourceReference, "days.0.state", "on"),
				),
			},
		},
	})
}

func TestAccClusterOnOffScheduleResourceInvalidTimezone(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_tz_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterOnOffScheduleResourceConfig(resourceName, "Mars/Olympus"),
				ExpectError: regexp.MustCompile(`(?s)timezone|invalid value|Validation Error`),
			},
		},
	})
}

func TestAccClusterOnOffScheduleResourceInvalidDays(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_days_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  timezone        = "US/Pacific"
  days = [
    { day = "monday",  state = "on" },
    { day = "tuesday", state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)7 days|days must|number of days`),
			},
		},
	})
}

func TestAccClusterOnOffScheduleResourceInvalidState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_state_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  timezone        = "US/Pacific"
  days = [
    { day = "monday",    state = "frozen" },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)state|invalid value|Validation Error`),
			},
		},
	})
}

func TestAccClusterOnOffScheduleResourceInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
  timezone        = "US/Pacific"
  days = [
    { day = "monday",    state = "on" },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error.*[Ss]chedule|cluster.*not found|access to the requested resource is denied`),
			},
		},
	})
}

func testAccClusterOnOffScheduleResourceConfig(resourceName, timezone string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  timezone        = "%[6]s"
  days = [
    {
      day   = "monday"
      state = "custom"
      from  = { hour = 0, minute = 0 }
      to    = { hour = 23, minute = 59 }
    },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, timezone)
}

// testAccClusterOnOffScheduleResourceUpdateConfig uses all-on days to exercise
// the Update (PUT) path — distinct from testAccClusterOnOffScheduleResourceConfig
// which uses a custom Monday to test from/to fields on create.
func testAccClusterOnOffScheduleResourceUpdateConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_onoff_schedule" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  timezone        = "US/Pacific"
  days = [
    { day = "monday",    state = "on" },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}
