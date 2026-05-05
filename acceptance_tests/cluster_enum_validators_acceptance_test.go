package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var enumValidatorError = regexp.MustCompile(`value must be one of:`)

// TestAccClusterEnumValidators verifies that schema-level enum validators
// reject out-of-range values at plan time for the cluster resource fields:
// availability.type, support.plan, support.timezone, service_groups[].services elements.
func TestAccClusterEnumValidators(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_validators_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Invalid availability.type
			{
				Config:      testAccClusterEnumConfig(resourceName, cidr, "invalid_type", "enterprise", "PT", `"data", "index", "query"`),
				ExpectError: enumValidatorError,
			},
			// Invalid support.plan
			{
				Config:      testAccClusterEnumConfig(resourceName, cidr, "multi", "invalid_plan", "PT", `"data", "index", "query"`),
				ExpectError: enumValidatorError,
			},
			// Invalid support.timezone
			{
				Config:      testAccClusterEnumConfig(resourceName, cidr, "multi", "enterprise", "EST", `"data", "index", "query"`),
				ExpectError: enumValidatorError,
			},
			// Invalid service in services set
			{
				Config:      testAccClusterEnumConfig(resourceName, cidr, "multi", "enterprise", "PT", `"data", "invalid_service"`),
				ExpectError: enumValidatorError,
			},
			// Valid configuration — accepted values for every enum field
			{
				Config: testAccClusterEnumConfig(resourceName, cidr, "multi", "enterprise", "PT", `"data", "index", "query"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("couchbase-capella_cluster."+resourceName, "id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_cluster."+resourceName, "etag"),
					resource.TestCheckResourceAttr("couchbase-capella_cluster."+resourceName, "availability.type", "multi"),
					resource.TestCheckResourceAttr("couchbase-capella_cluster."+resourceName, "support.plan", "enterprise"),
					resource.TestCheckResourceAttr("couchbase-capella_cluster."+resourceName, "support.timezone", "PT"),
				),
			},
		},
	})
}

func testAccClusterEnumConfig(resourceName, cidr, availabilityType, supportPlan, supportTimezone, services string) string {
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
    services     = [%[6]s]
  }]
  availability = { type = "%[7]s" }
  support      = { plan = "%[8]s", timezone = "%[9]s" }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr, services, availabilityType, supportPlan, supportTimezone)
}
