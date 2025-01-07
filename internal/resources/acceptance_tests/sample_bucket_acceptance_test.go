package acceptance_tests

import (
	"fmt"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccSampleBucket(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSampleBucketWithTravelSampleConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_sample_bucket.add_sample_bucket_travel", "name", "travel-sample"),
					resource.TestCheckResourceAttrSet("couchbase-capella_sample_bucket.add_sample_bucket_travel", "id"),
				),
			},
			//Invalid Sample Data input
			{
				Config:      testAccWithInvalidSampleInputConfig(),
				ExpectError: regexp.MustCompile("Could not load sample bucket"),
			},
		},
	})
}

func testAccSampleBucketWithTravelSampleConfig() string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_sample_bucket" "add_sample_bucket_travel" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "travel-sample"
}
`, ProviderBlock, OrgId, ProjectId, ClusterId)
}

func testAccWithInvalidSampleInputConfig() string {
	return fmt.Sprintf(`
%[1]s
output "sample_bucket_invalid"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_invalid
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_invalid" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "invalid-sample"
}
`, ProviderBlock, OrgId, ProjectId, ClusterId)
}
