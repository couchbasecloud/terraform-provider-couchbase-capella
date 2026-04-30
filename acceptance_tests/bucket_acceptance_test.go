package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBucketEnumValidators verifies that enum validators on the bucket resource
// reject out-of-enum values at plan time and accept valid values.
func TestAccBucketEnumValidators(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bucket_")
	resourceReference := "couchbase-capella_bucket." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Rejected: invalid type value
			{
				Config:      testAccBucketWithTypeConfig(resourceName, "invalid_type"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Rejected: invalid storage_backend value
			{
				Config:      testAccBucketWithStorageBackendConfig(resourceName, "invalid_backend"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Rejected: invalid bucket_conflict_resolution value
			{
				Config:      testAccBucketWithConflictResolutionConfig(resourceName, "invalid_cr"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Rejected: invalid durability_level value
			{
				Config:      testAccBucketWithDurabilityLevelConfig(resourceName, "invalid_level"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Rejected: invalid replicas value
			{
				Config:      testAccBucketWithReplicasConfig(resourceName, 5),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Rejected: invalid eviction_policy value
			{
				Config:      testAccBucketWithEvictionPolicyConfig(resourceName, "invalid_policy"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Accepted: valid enum values create the bucket successfully
			{
				Config: testAccBucketWithValidEnumValues(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "type", "couchbase"),
					resource.TestCheckResourceAttr(resourceReference, "storage_backend", "couchstore"),
					resource.TestCheckResourceAttr(resourceReference, "bucket_conflict_resolution", "seqno"),
					resource.TestCheckResourceAttr(resourceReference, "durability_level", "none"),
					resource.TestCheckResourceAttr(resourceReference, "replicas", "1"),
					resource.TestCheckResourceAttr(resourceReference, "eviction_policy", "fullEviction"),
				),
			},
		},
	})
}

func testAccBucketWithTypeConfig(resourceName, bucketType string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  type            = "%[6]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, bucketType)
}

func testAccBucketWithStorageBackendConfig(resourceName, storageBackend string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  storage_backend = "%[6]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, storageBackend)
}

func testAccBucketWithConflictResolutionConfig(resourceName, conflictResolution string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id           = "%[3]s"
  project_id                = "%[4]s"
  cluster_id                = "%[5]s"
  name                      = "%[2]s"
  bucket_conflict_resolution = "%[6]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, conflictResolution)
}

func testAccBucketWithDurabilityLevelConfig(resourceName, durabilityLevel string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  durability_level = "%[6]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, durabilityLevel)
}

func testAccBucketWithReplicasConfig(resourceName string, replicas int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  replicas        = %[6]d
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, replicas)
}

func testAccBucketWithEvictionPolicyConfig(resourceName, evictionPolicy string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  eviction_policy = "%[6]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, evictionPolicy)
}

func testAccBucketWithValidEnumValues(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id           = "%[3]s"
  project_id                = "%[4]s"
  cluster_id                = "%[5]s"
  name                      = "%[2]s"
  type                      = "couchbase"
  storage_backend           = "couchstore"
  bucket_conflict_resolution = "seqno"
  durability_level          = "none"
  replicas                  = 1
  eviction_policy           = "fullEviction"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}
