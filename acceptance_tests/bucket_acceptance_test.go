package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBucketEnumValidators verifies that schema-level enum validators
// reject out-of-range values at plan time for the bucket resource fields:
// type, storage_backend, bucket_conflict_resolution, durability_level, replicas, eviction_policy.
// Jira: AV-129333.
func TestAccBucketEnumValidators(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bucket_validators_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Invalid type
			{
				Config:      testAccBucketConfigWithInvalidType(resourceName),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid storage_backend
			{
				Config:      testAccBucketConfigWithInvalidStorageBackend(resourceName),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid bucket_conflict_resolution
			{
				Config:      testAccBucketConfigWithInvalidConflictResolution(resourceName),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid durability_level
			{
				Config:      testAccBucketConfigWithInvalidDurabilityLevel(resourceName),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid replicas
			{
				Config:      testAccBucketConfigWithInvalidReplicas(resourceName),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Invalid eviction_policy
			{
				Config:      testAccBucketConfigWithInvalidEvictionPolicy(resourceName),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Valid configuration — accepted values for every enum field
			{
				Config: testAccBucketConfigWithValidEnumFields(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_bucket."+resourceName, "type", "couchbase"),
					resource.TestCheckResourceAttr("couchbase-capella_bucket."+resourceName, "storage_backend", "couchstore"),
					resource.TestCheckResourceAttr("couchbase-capella_bucket."+resourceName, "bucket_conflict_resolution", "seqno"),
					resource.TestCheckResourceAttr("couchbase-capella_bucket."+resourceName, "durability_level", "none"),
					resource.TestCheckResourceAttr("couchbase-capella_bucket."+resourceName, "replicas", "1"),
					resource.TestCheckResourceAttr("couchbase-capella_bucket."+resourceName, "eviction_policy", "fullEviction"),
				),
			},
		},
	})
}

func testAccBucketConfigWithInvalidType(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  type            = "invalid_type"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccBucketConfigWithInvalidStorageBackend(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  storage_backend = "invalid_backend"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccBucketConfigWithInvalidConflictResolution(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_bucket" "%[2]s" {
  organization_id           = "%[3]s"
  project_id                = "%[4]s"
  cluster_id                = "%[5]s"
  name                      = "%[2]s"
  bucket_conflict_resolution = "invalid_resolution"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccBucketConfigWithInvalidDurabilityLevel(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  durability_level = "invalid_level"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccBucketConfigWithInvalidReplicas(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  replicas        = 5
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccBucketConfigWithInvalidEvictionPolicy(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  eviction_policy = "invalid_policy"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccBucketConfigWithValidEnumFields(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_bucket" "%[2]s" {
  organization_id            = "%[3]s"
  project_id                 = "%[4]s"
  cluster_id                 = "%[5]s"
  name                       = "%[2]s"
  type                       = "couchbase"
  storage_backend            = "couchstore"
  bucket_conflict_resolution = "seqno"
  durability_level           = "none"
  replicas                   = 1
  eviction_policy            = "fullEviction"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}
