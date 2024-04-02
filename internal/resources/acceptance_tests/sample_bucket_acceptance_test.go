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
	resourceName := "new_cluster"
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("aws")
	fmt.Println(cidr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	testCfg := acctest.ProjectCfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Creating cluster to check the sample bucket configs
			{
				Config: testAccCreateCluster(&testCfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
				),
			},
			{
				Config: testAccSampleBucketWithTravelSample(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.TestAccWait(time.Second*10),
					resource.TestCheckResourceAttr("couchbase-capella_sample_bucket.add_sample_bucket_travel", "name", "travel-sample"),
					resource.TestCheckResourceAttrSet("couchbase-capella_sample_bucket.add_sample_bucket_travel", "id"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccSampleBucketWithAllBuckets(&testCfg),
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
				Config:      testAccWithInvalidSampleInput(testCfg),
				ExpectError: regexp.MustCompile("Could not load sample bucket"),
			},

			{
				Config: testCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDeleteClusterResource(resourceReference),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccWithInvalidSampleInput(cfg string) string {
	cfg = fmt.Sprintf(`
%[1]s
output "sample_bucket_invalid"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_invalid
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_invalid" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name			  = "invalid-sample"
}
`, cfg)
	return cfg
}

func testAccSampleBucketWithTravelSample(cfg string) string {
	cfg = fmt.Sprintf(`
%[1]s
output "sample_bucket_travel"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_travel
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_travel" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name			  = "travel-sample"
}
`, cfg)
	return cfg
}

func testAccSampleBucketWithAllBuckets(cfg *string) string {
	*cfg = fmt.Sprintf(`
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
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name			  = "travel-sample"
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_beer" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name			  = "beer-sample"
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_gamesim" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name			  = "gamesim-sample"
}

`, *cfg)
	return *cfg
}
