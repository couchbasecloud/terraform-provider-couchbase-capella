package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSampleBucket(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_sample_bucket_")
	resourceReference := "couchbase-capella_sample_bucket." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccSampleBucketWithTravelSampleConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "travel-sample"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
			// Rejected: invalid sample bucket name triggers plan-time validator error
			{
				Config:      testAccWithInvalidSampleInputConfig(resourceName),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
		},
	})
}

func testAccSampleBucketWithTravelSampleConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_sample_bucket" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "travel-sample"
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccWithInvalidSampleInputConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_sample_bucket" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "invalid-sample"
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}
