package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
)

func TestAccSampleBucketTestCases(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSampleBucketWithTravelSampleConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.TestAccWait(time.Second*10),
					resource.TestCheckResourceAttr("couchbase-capella_sample_bucket.add_sample_bucket_travel", "name", "travel-sample"),
					resource.TestCheckResourceAttrSet("couchbase-capella_sample_bucket.add_sample_bucket_travel", "id"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccSampleBucketWithAllBucketsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_sample_bucket.add_sample_bucket_travel", "name", "travel-sample"),
					resource.TestCheckResourceAttrSet("couchbase-capella_sample_bucket.add_sample_bucket_travel", "id"),

					resource.TestCheckResourceAttr("couchbase-capella_sample_bucket.add_sample_bucket_gamesim", "name", "gamesim-sample"),
					resource.TestCheckResourceAttrSet("couchbase-capella_sample_bucket.add_sample_bucket_gamesim", "id"),

					resource.TestCheckResourceAttr("couchbase-capella_sample_bucket.add_sample_bucket_beer", "name", "beer-sample"),
					resource.TestCheckResourceAttrSet("couchbase-capella_sample_bucket.add_sample_bucket_beer", "id"),
				),
				ExpectNonEmptyPlan: true,
			},

			//Invalid Sample Data input
			{
				Config:      testAccWithInvalidSampleInputConfig(),
				ExpectError: regexp.MustCompile("Could not load sample bucket"),
			},
		},
	})
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

func testAccSampleBucketWithTravelSampleConfig() string {
	return fmt.Sprintf(`
%[1]s
output "sample_bucket_travel"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_travel
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_travel" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "travel-sample"
}
`, ProviderBlock, OrgId, ProjectId, ClusterId)
}

func testAccSampleBucketWithAllBucketsConfig() string {
	return fmt.Sprintf(`
%[1]s

output "sample_bucket_beer"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_beer
}
output "sample_bucket_gamesim"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_gamesim
}
output "sample_bucket_travel"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_travel
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_travel" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "travel-sample"
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_beer" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "beer-sample"
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_gamesim" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name			  = "gamesim-sample"
}

`, ProviderBlock, OrgId, ProjectId, ClusterId)

}
