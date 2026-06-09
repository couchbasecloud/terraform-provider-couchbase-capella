package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFreeTierBucketResource(t *testing.T) {
	clusterName := randomStringWithPrefix("tf_acc_free_tier_cluster_bucket_")
	bucketName := randomStringWithPrefix("tf_acc_free_tier_bucket_")
	resourceReference := "couchbase-capella_free_tier_bucket." + bucketName
	cidr := generateRandomCIDR()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeTierBucketResourceConfig(clusterName, bucketName, cidr, 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "name", bucketName),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "100"),
					resource.TestCheckResourceAttrSet(resourceReference, "type"),
					resource.TestCheckResourceAttrSet(resourceReference, "storage_backend"),
					resource.TestCheckResourceAttrSet(resourceReference, "vbuckets"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateBucketImportIdForResource(resourceReference),
				ImportState:       true,
			},
			{
				Config: testAccFreeTierBucketResourceConfig(clusterName, bucketName, cidr, 200),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", bucketName),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "200"),
				),
			},
		},
	})
}

func TestAccDatasourceFreeTierBuckets(t *testing.T) {
	clusterName := randomStringWithPrefix("tf_acc_free_tier_cluster_buckets_ds_")
	bucketName := randomStringWithPrefix("tf_acc_free_tier_bucket_ds_")
	dsName := randomStringWithPrefix("tf_acc_free_tier_buckets_ds_")
	bucketReference := "couchbase-capella_free_tier_bucket." + bucketName
	dsReference := "data.couchbase-capella_free_tier_buckets." + dsName
	cidr := generateRandomCIDR()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeTierBucketsDatasourceConfig(clusterName, bucketName, dsName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bucketReference, "name", bucketName),
					resource.TestCheckResourceAttrSet(bucketReference, "id"),
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(dsReference, "cluster_id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					testAccCheckListElemNestedAttrs(dsReference, "data", map[string]string{
						"name":            bucketName,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
					}),
				),
			},
		},
	})
}

func testAccCheckListElemNestedAttrs(resourceName, listAttribute string, expected map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, module := range state.Modules {
			resourceState, ok := module.Resources[resourceName]
			if !ok {
				continue
			}

			attrs := resourceState.Primary.Attributes
			for index := 0; attrs[fmt.Sprintf("%s.%d.name", listAttribute, index)] != ""; index++ {
				matched := true
				for key, expectedValue := range expected {
					if attrs[fmt.Sprintf("%s.%d.%s", listAttribute, index, key)] != expectedValue {
						matched = false
						break
					}
				}
				if matched {
					return nil
				}
			}
		}

		return fmt.Errorf("no %s element matched expected attributes %v for %s", listAttribute, expected, resourceName)
	}
}

func testAccFreeTierBucketResourceConfig(clusterName, bucketName, cidr string, memoryAllocationInMB int64) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Free tier cluster for bucket acceptance testing."

  cloud_provider = {
    type   = "aws"
    region = "us-east-2"
    cidr   = "%[6]s"
  }
}

resource "couchbase-capella_free_tier_bucket" "%[5]s" {
  organization_id          = "%[2]s"
  project_id               = "%[3]s"
  cluster_id               = couchbase-capella_free_tier_cluster.%[4]s.id
  name                     = "%[5]s"
  memory_allocation_in_mb  = %[7]d
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, bucketName, cidr, memoryAllocationInMB)
}

func testAccFreeTierBucketsDatasourceConfig(clusterName, bucketName, dsName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Free tier cluster for buckets data source acceptance testing."

  cloud_provider = {
    type   = "aws"
    region = "us-east-2"
    cidr   = "%[7]s"
  }
}

resource "couchbase-capella_free_tier_bucket" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id
  name            = "%[5]s"
}

data "couchbase-capella_free_tier_buckets" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id

  depends_on = [couchbase-capella_free_tier_bucket.%[5]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, bucketName, dsName, cidr)
}
