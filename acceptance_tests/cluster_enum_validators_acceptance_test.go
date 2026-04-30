package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccClusterAvailabilityTypeEnumValidator verifies that the availability.type field
// rejects out-of-enum values and accepts valid ones.
func TestAccClusterAvailabilityTypeEnumValidator(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Rejected: invalid availability type
			{
				Config:      testAccClusterWithAvailabilityType(resourceName, cidr, "invalid_type"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Accepted: valid availability type — plan only, no cluster deployed
			{
				Config:             testAccClusterWithAvailabilityType(resourceName, cidr, "multi"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccClusterSupportPlanEnumValidator verifies that the support.plan field
// rejects out-of-enum values and accepts valid ones.
func TestAccClusterSupportPlanEnumValidator(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Rejected: invalid support plan
			{
				Config:      testAccClusterWithSupportPlan(resourceName, cidr, "invalid_plan"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Accepted: valid support plan — plan only, no cluster deployed
			{
				Config:             testAccClusterWithSupportPlan(resourceName, cidr, "enterprise"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccClusterSupportTimezoneEnumValidator verifies that the support.timezone field
// rejects out-of-enum values and accepts valid ones.
func TestAccClusterSupportTimezoneEnumValidator(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Rejected: invalid timezone
			{
				Config:      testAccClusterWithTimezone(resourceName, cidr, "invalid_tz"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Accepted: valid timezone — plan only, no cluster deployed
			{
				Config:             testAccClusterWithTimezone(resourceName, cidr, "PT"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccClusterServicesEnumValidator verifies that the service_groups[].services set
// rejects out-of-enum values and accepts valid ones.
func TestAccClusterServicesEnumValidator(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Rejected: invalid service name in set
			{
				Config:      testAccClusterWithServices(resourceName, cidr, `["data", "invalid_service"]`),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Accepted: valid services — plan only, no cluster deployed
			{
				Config:             testAccClusterWithServices(resourceName, cidr, `["data", "index", "query"]`),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccClusterWithAvailabilityType(resourceName, cidr, availabilityType string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  name            = "%[2]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = { cpu = 4, ram = 16 }
        disk    = { storage = 50, type = "io2", iops = 3000 }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = { type = "%[6]s" }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr, availabilityType)
}

func testAccClusterWithSupportPlan(resourceName, cidr, plan string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  name            = "%[2]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = { cpu = 4, ram = 16 }
        disk    = { storage = 50, type = "io2", iops = 3000 }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = { type = "multi" }
  support = {
    plan     = "%[6]s"
    timezone = "PT"
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr, plan)
}

func testAccClusterWithTimezone(resourceName, cidr, timezone string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  name            = "%[2]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = { cpu = 4, ram = 16 }
        disk    = { storage = 50, type = "io2", iops = 3000 }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = { type = "multi" }
  support = {
    plan     = "enterprise"
    timezone = "%[6]s"
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr, timezone)
}

func testAccClusterWithServices(resourceName, cidr, services string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  name            = "%[2]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = { cpu = 4, ram = 16 }
        disk    = { storage = 50, type = "io2", iops = 3000 }
      }
      num_of_nodes = 3
      services     = %[6]s
    }
  ]
  availability = { type = "multi" }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr, services)
}
