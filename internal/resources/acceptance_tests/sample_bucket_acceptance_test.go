package acceptance_tests

import (
	"fmt"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccSampleBucket(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_sample_bucket_")
	resourceReference := "couchbase-capella_sample_bucket." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSampleBucketWithTravelSampleConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "travel-sample"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
			//Invalid Sample Data input
			{
				Config:      testAccWithInvalidSampleInputConfig(resourceName),
				ExpectError: regexp.MustCompile("Could not load sample bucket"),
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
`, ProviderBlock, OrgId, ProjectId, ClusterId, resourceName)
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
`, ProviderBlock, OrgId, ProjectId, ClusterId, resourceName)
}
