package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGSI(t *testing.T) {
	const resourceType = "couchbase-capella_query_indexes"
	primaryIndexResourceName := randomStringWithPrefix("tf_acc_gsi_")
	primaryIndexResourceReference := fmt.Sprintf("%s.%s", resourceType, primaryIndexResourceName)
	gsiResourceIdx1 := fmt.Sprintf("%s.%s", resourceType, "idx1")
	gsiResourceIdx2 := fmt.Sprintf("%s.%s", resourceType, "idx2")
	gsiResourceIdx3 := fmt.Sprintf("%s.%s", resourceType, "idx3")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateGSINonDeferredIndexConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(gsiResourceIdx1, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "index_name", "idx1"),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "index_keys.0", "c1"),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "where", "geo.alt > 1000"),

					resource.TestCheckResourceAttr(gsiResourceIdx2, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "index_name", "idx2"),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "index_keys.0", "c2"),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "where", "geo.alt > 2000"),

					resource.TestCheckResourceAttr(gsiResourceIdx3, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "index_name", "idx3"),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "index_keys.0", "c3"),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "where", "geo.alt > 3000"),
				),
			},
			{
				ResourceName:      gsiResourceIdx1,
				ImportStateIdFunc: generateGsiImportIdForResource(gsiResourceIdx1),
				ImportState:       true,
			},
			{
				Config: testAccCreatePrimaryIndexConfig(primaryIndexResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "index_name", "primary_index"),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "with.num_replica", "1"),
				),
			},
		},
	})
}

func testAccCreateGSINonDeferredIndexConfig() string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "idx1" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "idx1"
  index_keys      = ["c1"]
  where = "geo.alt > 1000"
  with = {
        defer_build = false
  }
}

resource "couchbase-capella_query_indexes" "idx2" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "idx2"
  index_keys      = ["c2"]
  where = "geo.alt > 2000"
  with = {
        defer_build = false
  }
}

resource "couchbase-capella_query_indexes" "idx3" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "idx3"
  index_keys      = ["c3"]
  where = "geo.alt > 3000"
  with = {
        defer_build = false
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalBucketName,
		globalScopeName,
		globalCollectionName,
	)
}

func testAccCreatePrimaryIndexConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "primary_index"
  is_primary      = true
  with = {
        num_replica = 1
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalBucketName,
		globalScopeName,
		globalCollectionName,
		resourceName)
}

// TestAccGSIIncludeMissingImport verifies that a query index whose keys carry
// modifiers (INCLUDE MISSING) and a plain multi-key index import cleanly, with
// no diff on the next plan (i.e. no unnecessary forced replacement, AV-134985).
//
// This exercises the control-plane fix that returns secExprs with key modifiers
// preserved; it therefore requires a cluster running that fix. ImportStateVerify
// also covers the null-handling on import (is_primary/where/num_partition), which
// are all omitted here and must import as null to avoid a spurious diff.
//
// The index_keys must be written in the server's canonical, backtick-quoted form
// so the imported secExprs match the configured value.
func TestAccGSIIncludeMissingImport(t *testing.T) {
	const resourceType = "couchbase-capella_query_indexes"
	includeMissingRef := fmt.Sprintf("%s.%s", resourceType, "include_missing")
	compositeRef := fmt.Sprintf("%s.%s", resourceType, "composite")

	// status is server-runtime state and defer_build is a create-only input not
	// returned on read, so neither round-trips on import.
	importVerifyIgnore := []string{"status", "with.defer_build"}

	// the resource has no id attribute, which is what ImportStateVerify pairs
	// old and new state on by default.
	const importVerifyIdentifier = "index_name"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateGSIIncludeMissingConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(includeMissingRef, "index_name", "include_missing_idx"),
					resource.TestCheckResourceAttr(includeMissingRef, "index_keys.0", "`name` INCLUDE MISSING"),
					resource.TestCheckResourceAttr(compositeRef, "index_name", "composite_idx"),
					resource.TestCheckResourceAttr(compositeRef, "index_keys.0", "`c1`"),
					resource.TestCheckResourceAttr(compositeRef, "index_keys.1", "`c2`"),
				),
			},
			{
				ResourceName:                         includeMissingRef,
				ImportStateIdFunc:                    generateGsiImportIdForResource(includeMissingRef),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: importVerifyIdentifier,
				ImportStateVerifyIgnore:              importVerifyIgnore,
			},
			{
				ResourceName:                         compositeRef,
				ImportStateIdFunc:                    generateGsiImportIdForResource(compositeRef),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: importVerifyIdentifier,
				ImportStateVerifyIgnore:              importVerifyIgnore,
			},
			{
				// The reported symptom is a forced replacement on the next plan,
				// so assert the plan is empty rather than only comparing states.
				Config:   testAccCreateGSIIncludeMissingConfig(),
				PlanOnly: true,
			},
		},
	})
}

func testAccCreateGSIIncludeMissingConfig() string {
	// Backticks cannot appear in a Go raw-string literal, so the canonical,
	// backtick-quoted key lists are passed in as parameters.
	includeMissingKeys := "[\"`name` INCLUDE MISSING\"]"
	compositeKeys := "[\"`c1`\", \"`c2`\"]"

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "include_missing" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "include_missing_idx"
  index_keys      = %[8]s
  with = {
        defer_build = false
  }
}

resource "couchbase-capella_query_indexes" "composite" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "composite_idx"
  index_keys      = %[9]s
  with = {
        defer_build = false
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalBucketName,
		globalScopeName,
		globalCollectionName,
		includeMissingKeys,
		compositeKeys,
	)
}

func generateGsiImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"index_name=%s,collection_name=%s,scope_name=%s,bucket_name=%s,cluster_id=%s,organization_id=%s,project_id=%s",
			rawState["index_name"],
			rawState["collection_name"],
			rawState["scope_name"],
			rawState["bucket_name"],
			rawState["cluster_id"],
			rawState["organization_id"],
			rawState["project_id"],
		), nil
	}
}
