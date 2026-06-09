package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
				// Update step: change only the timezone. A larger payload diff
				// (e.g. rewriting the whole days array) makes the singleton
				// schedule PUT return HTTP 500 (code 10000) from Capella for
				// minutes while it propagates. Keep the change minimal so the
				// Update path is exercised without straining the backend.
				Config: testAccClusterOnOffScheduleResourceConfig(resourceName, "US/Eastern"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "timezone", "US/Eastern"),
					resource.TestCheckResourceAttr(resourceReference, "days.#", "7"),
				),
			},
			{
				// The schedule resource has no simple "id" attribute; it's keyed
				// by compound organization_id+project_id+cluster_id, with
				// cluster_id as the natural primary key (one schedule per
				// cluster). Tell the framework to verify import using
				// cluster_id instead of the default "id".
				ResourceName:                         resourceReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    generateClusterOnOffScheduleImportId(resourceReference),
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
		},
	})
}

func generateClusterOnOffScheduleImportId(resourceReference string) resource.ImportStateIdFunc {
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
			"organization_id=%s,project_id=%s,cluster_id=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
		), nil
	}
}

func TestAccClusterOnOffScheduleResourceInvalidTimezone(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_tz_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterOnOffScheduleResourceConfig(resourceName, "Mars/Olympus"),
				ExpectError: regexp.MustCompile(`must be one of`),
			},
		},
	})
}

// TestAccClusterOnOffScheduleEmptyTimezone verifies that an empty timezone is
// rejected by local config validation (the enum OneOf validator) at terraform
// validate instead of being sent to the API.
func TestAccClusterOnOffScheduleEmptyTimezone(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_empty_tz_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterOnOffScheduleResourceConfig(resourceName, ""),
				ExpectError: regexp.MustCompile(`must be one of`),
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

// TestAccClusterOnOffScheduleAllDaysOff verifies that a schedule with every day
// set to "off" is rejected by local config validation instead of being sent to the API
func TestAccClusterOnOffScheduleAllDaysOff(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_all_off_")

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
    { day = "monday",    state = "off" },
    { day = "tuesday",   state = "off" },
    { day = "wednesday", state = "off" },
    { day = "thursday",  state = "off" },
    { day = "friday",    state = "off" },
    { day = "saturday",  state = "off" },
    { day = "sunday",    state = "off" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				// Subset from the start of the message so Terraform's diagnostic
				// word-wrapping cannot split the match.
				ExpectError: regexp.MustCompile(`Clusters cannot be scheduled`),
			},
		},
	})
}

// TestAccClusterOnOffScheduleOutOfOrderDays verifies that a schedule
// whose days are not in Monday-to-Sunday order is rejected by local config validation.
func TestAccClusterOnOffScheduleOutOfOrderDays(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_out_of_order_")

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
    { day = "sunday",    state = "on" },
    { day = "monday",    state = "on" },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				// Subset from the start of the message so Terraform's diagnostic
				// word-wrapping cannot split the match.
				ExpectError: regexp.MustCompile(`must be in sequence`),
			},
		},
	})
}

// TestAccClusterOnOffScheduleCustomWithoutFrom verifies that a custom
// day without the required from time boundary is rejected by local config
// validation instead of being sent to the API
func TestAccClusterOnOffScheduleCustomWithoutFrom(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_custom_no_from_")

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
    {
      day   = "monday"
      state = "custom"
      to    = { hour = 18, minute = 30 }
    },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`from time boundary is required`),
			},
		},
	})
}

// TestAccClusterOnOffScheduleBoundaryOnNonCustomDay verifies that a
// day with state "on" or "off" cannot contain from/to time boundaries
func TestAccClusterOnOffScheduleBoundaryOnNonCustomDay(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_boundary_on_")

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
    { day = "monday",    state = "on", from = { hour = 8, minute = 0 } },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`cannot contain from/to`),
			},
		},
	})
}

// TestAccClusterOnOffScheduleFromAfterTo verifies that a custom day
// whose from time boundary is later than its to time boundary is rejected by
// local config validation
func TestAccClusterOnOffScheduleFromAfterTo(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_from_after_to_")

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
    {
      day   = "monday"
      state = "custom"
      from  = { hour = 18, minute = 30 }
      to    = { hour = 8, minute = 0 }
    },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`must not be later than`),
			},
		},
	})
}

// TestAccClusterOnOffScheduleInvalidBoundaryHour verifies that a time
// boundary hour outside the valid 0-23 range is rejected by schema validation
func TestAccClusterOnOffScheduleInvalidBoundaryHour(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_hour_")

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
    {
      day   = "monday"
      state = "custom"
      from  = { hour = 24, minute = 0 }
    },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`must be between 0 and 23`),
			},
		},
	})
}

// TestAccClusterOnOffScheduleInvalidBoundaryMinute verifies that a
// time boundary minute other than the valid values 0 and 30 is rejected by
// schema validation
func TestAccClusterOnOffScheduleInvalidBoundaryMinute(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_onoff_schedule_invalid_minute_")

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
    {
      day   = "monday"
      state = "custom"
      from  = { hour = 8, minute = 15 }
    },
    { day = "tuesday",   state = "on" },
    { day = "wednesday", state = "on" },
    { day = "thursday",  state = "on" },
    { day = "friday",    state = "on" },
    { day = "saturday",  state = "on" },
    { day = "sunday",    state = "on" },
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`must be one of`),
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
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, timezone)
}
