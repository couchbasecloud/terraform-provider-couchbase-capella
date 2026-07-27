package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Live on/off coverage (toggling a real cluster off and on, plus import) lives
// in TestAccFreeTierClusterLifecycle. This validation-only test fails at plan
// time and never touches a cluster.
func TestAccFreeTierClusterOnOffResourceInvalidState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_on_off_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccFreeTierClusterOnOffResourceConfig(resourceName, "frozen"),
				ExpectError: regexp.MustCompile(`(?s)invalid state value|invalid state`),
			},
		},
	})
}

func testAccFreeTierClusterOnOffResourceConfig(resourceName, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster_on_off" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  state           = "%[3]s"
}
`, globalProviderBlock, resourceName, state)
}
