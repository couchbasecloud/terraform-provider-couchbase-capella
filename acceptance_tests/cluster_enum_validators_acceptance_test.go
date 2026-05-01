package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccClusterEnumValidators_AV_129334 verifies that schema-level enum validators
// reject out-of-range values at plan time for the cluster resource fields:
// availability.type, support.plan, support.timezone, service_groups[].services elements.
func TestAccClusterEnumValidators_AV_129334(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_validators_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Invalid availability.type
			{
				Config:      testAccClusterConfigWithInvalidAvailabilityType(resourceName, cidr),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid support.plan
			{
				Config:      testAccClusterConfigWithInvalidSupportPlan(resourceName, cidr),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid support.timezone
			{
				Config:      testAccClusterConfigWithInvalidSupportTimezone(resourceName, cidr),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid service in services set
			{
				Config:      testAccClusterConfigWithInvalidService(resourceName, cidr),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Valid configuration — accepted values for every enum field
			{
				Config: testAccClusterConfigWithValidEnumFields(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_cluster."+resourceName, "availability.type", "multi"),
					resource.TestCheckResourceAttr("couchbase-capella_cluster."+resourceName, "support.plan", "enterprise"),
					resource.TestCheckResourceAttr("couchbase-capella_cluster."+resourceName, "support.timezone", "PT"),
				),
			},
		},
	})
}

func testAccClusterConfigWithInvalidAvailabilityType(resourceName, cidr string) string {
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
  service_groups = [{
    node = {
      compute = { cpu = 4, ram = 16 }
      disk    = { storage = 50, type = "io2", iops = 3000 }
    }
    num_of_nodes = 3
    services     = ["data", "index", "query"]
  }]
  availability = { type = "invalid_type" }
  support      = { plan = "enterprise", timezone = "PT" }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr)
}

func testAccClusterConfigWithInvalidSupportPlan(resourceName, cidr string) string {
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
  service_groups = [{
    node = {
      compute = { cpu = 4, ram = 16 }
      disk    = { storage = 50, type = "io2", iops = 3000 }
    }
    num_of_nodes = 3
    services     = ["data", "index", "query"]
  }]
  availability = { type = "multi" }
  support      = { plan = "invalid_plan", timezone = "PT" }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr)
}

func testAccClusterConfigWithInvalidSupportTimezone(resourceName, cidr string) string {
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
  service_groups = [{
    node = {
      compute = { cpu = 4, ram = 16 }
      disk    = { storage = 50, type = "io2", iops = 3000 }
    }
    num_of_nodes = 3
    services     = ["data", "index", "query"]
  }]
  availability = { type = "multi" }
  support      = { plan = "enterprise", timezone = "EST" }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr)
}

func testAccClusterConfigWithInvalidService(resourceName, cidr string) string {
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
  service_groups = [{
    node = {
      compute = { cpu = 4, ram = 16 }
      disk    = { storage = 50, type = "io2", iops = 3000 }
    }
    num_of_nodes = 3
    services     = ["data", "invalid_service"]
  }]
  availability = { type = "multi" }
  support      = { plan = "enterprise", timezone = "PT" }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr)
}

func testAccClusterConfigWithValidEnumFields(resourceName, cidr string) string {
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
  service_groups = [{
    node = {
      compute = { cpu = 4, ram = 16 }
      disk    = { storage = 50, type = "io2", iops = 3000 }
    }
    num_of_nodes = 3
    services     = ["data", "index", "query"]
  }]
  availability = { type = "multi" }
  support      = { plan = "enterprise", timezone = "PT" }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr)
}
