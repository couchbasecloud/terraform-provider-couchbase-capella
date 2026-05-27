package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBucketResourceRequiredOnly(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_req_")
	resourceReference := "couchbase-capella_bucket." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketResourceConfigRequiredOnly(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "type"),
					resource.TestCheckResourceAttrSet(resourceReference, "storage_backend"),
					resource.TestCheckResourceAttrSet(resourceReference, "memory_allocation_in_mb"),
					resource.TestCheckResourceAttrSet(resourceReference, "durability_level"),
					resource.TestCheckResourceAttrSet(resourceReference, "replicas"),
					resource.TestCheckResourceAttrSet(resourceReference, "eviction_policy"),
					resource.TestCheckResourceAttr(resourceReference, "flush", "false"),
					resource.TestCheckResourceAttr(resourceReference, "bucket_conflict_resolution", "seqno"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.item_count"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.ops_per_second"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.disk_used_in_mib"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.memory_used_in_mib"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateIdFunc: generateBucketImportIdForResource(resourceReference),
			},
		},
	})
}

func TestAccBucketResourceAllOptionalFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_full_")
	resourceReference := "couchbase-capella_bucket." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketResourceConfigAllOptional(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "type", "couchbase"),
					resource.TestCheckResourceAttr(resourceReference, "storage_backend", "couchstore"),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "200"),
					resource.TestCheckResourceAttr(resourceReference, "bucket_conflict_resolution", "seqno"),
					resource.TestCheckResourceAttr(resourceReference, "durability_level", "majority"),
					resource.TestCheckResourceAttr(resourceReference, "replicas", "2"),
					resource.TestCheckResourceAttr(resourceReference, "flush", "true"),
					resource.TestCheckResourceAttr(resourceReference, "time_to_live_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceReference, "eviction_policy", "fullEviction"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

func TestAccBucketResourceUpdate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_upd_")
	resourceReference := "couchbase-capella_bucket." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketResourceConfigUpdatable(resourceName, 200, "none", 1, 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "200"),
					resource.TestCheckResourceAttr(resourceReference, "durability_level", "none"),
					resource.TestCheckResourceAttr(resourceReference, "replicas", "1"),
					resource.TestCheckResourceAttr(resourceReference, "time_to_live_in_seconds", "0"),
				),
			},
			{
				Config: testAccBucketResourceConfigUpdatable(resourceName, 256, "majority", 2, 3600),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "256"),
					resource.TestCheckResourceAttr(resourceReference, "durability_level", "majority"),
					resource.TestCheckResourceAttr(resourceReference, "replicas", "2"),
					resource.TestCheckResourceAttr(resourceReference, "time_to_live_in_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", dmClusterId),
				),
			},
		},
	})
}

func TestAccBucketResourceEphemeralType(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_eph_")
	resourceReference := "couchbase-capella_bucket." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketResourceConfigEphemeral(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "type", "ephemeral"),
					resource.TestCheckResourceAttr(resourceReference, "eviction_policy", "nruEviction"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

func TestAccBucketResourceInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
  name            = "%[2]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error creating bucket|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccBucketResourceInvalidProject(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_bad_project_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "00000000-0000-0000-0000-000000000000"
  cluster_id      = "%[4]s"
  name            = "%[2]s"
}
`, globalProviderBlock, resourceName, globalOrgId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error creating bucket|project.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccBucketResourceInvalidDurabilityLevel(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_bad_dur_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id  = "%[3]s"
  project_id       = "%[4]s"
  cluster_id       = "%[5]s"
  name             = "%[2]s"
  durability_level = "invalidLevel"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error creating bucket|durability|invalid|unrecognized`),
			},
		},
	})
}

func TestAccBucketResourceInvalidReplicaCount(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bkt_bad_rep_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  replicas        = 4
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error creating bucket|replica|invalid|exceeds`),
			},
		},
	})
}

func TestAccScopeResource(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_scope_bkt_")
	scopeName := randomStringWithPrefix("tf_acc_scope_")
	bucketReference := "couchbase-capella_bucket." + bucketName
	scopeReference := "couchbase-capella_scope." + scopeName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccScopeResourceConfig(bucketName, scopeName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(scopeReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(scopeReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(scopeReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttrPair(scopeReference, "bucket_id", bucketReference, "id"),
					resource.TestCheckResourceAttr(scopeReference, "scope_name", scopeName),
					resource.TestCheckResourceAttrSet(scopeReference, "collections.#"),
				),
			},
			{
				ResourceName:                         scopeReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateScopeImportIdForResource(scopeReference),
				ImportStateVerifyIdentifierAttribute: "scope_name",
			},
		},
	})
}

func TestAccScopeResourceInvalidBucket(t *testing.T) {
	scopeName := randomStringWithPrefix("tf_acc_scope_bad_bkt_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_scope" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "nonexistent-bucket-id"
  scope_name      = "%[2]s"
}
`, globalProviderBlock, scopeName, globalOrgId, globalProjectId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error.*scope|bucket.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccScopeResourceInvalidCluster(t *testing.T) {
	scopeName := randomStringWithPrefix("tf_acc_scope_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_scope" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
  bucket_id       = "%[5]s"
  scope_name      = "%[2]s"
}
`, globalProviderBlock, scopeName, globalOrgId, globalProjectId, dmBucketId),
				ExpectError: regexp.MustCompile(`(?s)Error.*scope|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccCollectionResourceNoTTL(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_coll_bkt_nottl_")
	scopeName := randomStringWithPrefix("tf_acc_coll_scope_nottl_")
	collName := randomStringWithPrefix("tf_acc_coll_nottl_")
	collReference := "couchbase-capella_collection." + collName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigNoTTL(bucketName, scopeName, collName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(collReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(collReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttr(collReference, "scope_name", scopeName),
					resource.TestCheckResourceAttr(collReference, "collection_name", collName),
					resource.TestCheckResourceAttrSet(collReference, "bucket_id"),
					resource.TestCheckResourceAttrSet(collReference, "max_ttl"),
				),
			},
			{
				ResourceName:                         collReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateCollectionImportIdForResource(collReference),
				ImportStateVerifyIdentifierAttribute: "collection_name",
			},
		},
	})
}

func TestAccCollectionResourceWithTTL(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_coll_bkt_ttl_")
	scopeName := randomStringWithPrefix("tf_acc_coll_scope_ttl_")
	collName := randomStringWithPrefix("tf_acc_coll_ttl_")
	collReference := "couchbase-capella_collection." + collName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigWithTTL(bucketName, scopeName, collName, 3600),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(collReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(collReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttr(collReference, "scope_name", scopeName),
					resource.TestCheckResourceAttr(collReference, "collection_name", collName),
					resource.TestCheckResourceAttr(collReference, "max_ttl", "3600"),
					resource.TestCheckResourceAttrSet(collReference, "bucket_id"),
				),
			},
			{
				Config: testAccCollectionResourceConfigWithTTL(bucketName, scopeName, collName, 7200),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(collReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(collReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttr(collReference, "scope_name", scopeName),
					resource.TestCheckResourceAttr(collReference, "collection_name", collName),
					resource.TestCheckResourceAttr(collReference, "max_ttl", "7200"),
					resource.TestCheckResourceAttrSet(collReference, "bucket_id"),
				),
			},
		},
	})
}

func TestAccCollectionResourceInvalidScope(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_coll_bad_scope_bkt_")
	collName := randomStringWithPrefix("tf_acc_coll_bad_scope_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccCollectionResourceConfigInvalidScope(bucketName, collName),
				ExpectError: regexp.MustCompile(`(?s)Error.*collection|scope.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccCollectionResourceInvalidBucket(t *testing.T) {
	collName := randomStringWithPrefix("tf_acc_coll_bad_bkt_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_collection" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "00000000-0000-0000-0000-000000000000"
  scope_name      = "_default"
  collection_name = "%[2]s"
}
`, globalProviderBlock, collName, globalOrgId, globalProjectId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error.*collection|bucket.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

// AV-132307: API accepts negative max_ttl; restore ExpectError once fixed.
func TestAccCollectionResourceNegativeTTL(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_coll_neg_ttl_bkt_")
	scopeName := randomStringWithPrefix("tf_acc_coll_neg_ttl_scope_")
	collName := randomStringWithPrefix("tf_acc_coll_neg_ttl_")
	collReference := "couchbase-capella_collection." + collName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigWithTTL(bucketName, scopeName, collName, -1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collReference, "collection_name", collName),
					resource.TestCheckResourceAttr(collReference, "max_ttl", "-1"),
				),
			},
		},
	})
}

func TestAccFlushBucketResource(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_flush_bkt_")
	flushName := randomStringWithPrefix("tf_acc_flush_")
	bucketReference := "couchbase-capella_bucket." + bucketName
	flushReference := "couchbase-capella_flush." + flushName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFlushBucketResourceConfig(bucketName, flushName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(flushReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(flushReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(flushReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttrPair(flushReference, "bucket_id", bucketReference, "id"),
				),
			},
		},
	})
}

func TestAccFlushBucketResourceFlushDisabled(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_flush_disabled_bkt_")
	flushName := randomStringWithPrefix("tf_acc_flush_disabled_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlushBucketResourceConfigFlushDisabled(bucketName, flushName),
				ExpectError: regexp.MustCompile(`(?s)Error flushing the bucket|flush.*not enabled|flush.*disabled|flushEnabled`),
			},
		},
	})
}

func TestAccFlushBucketResourceInvalidBucket(t *testing.T) {
	flushName := randomStringWithPrefix("tf_acc_flush_bad_bkt_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_flush" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "nonexistent-bucket-id"
}
`, globalProviderBlock, flushName, globalOrgId, globalProjectId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error flushing the bucket|bucket.*not found|Not Found`),
			},
		},
	})
}

func TestAccDatasourceBuckets(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_buckets_")
	dsReference := "data.couchbase-capella_buckets." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketsDatasourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", dmClusterId),
					testAccCheckBucketsNonEmpty(dsReference),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.name"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.vbuckets"),
				),
			},
		},
	})
}

func TestAccDatasourceBucketsInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_buckets_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_buckets" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Buckets|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceScopes(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_scopes_")
	dsReference := "data.couchbase-capella_scopes." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccScopesDatasourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttr(dsReference, "bucket_id", dmBucketId),
					testAccCheckScopesNonEmpty(dsReference),
					resource.TestCheckResourceAttrSet(dsReference, "scopes.0.scope_name"),
				),
			},
		},
	})
}

func TestAccDatasourceScopesInvalidBucket(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_scopes_bad_bkt_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_scopes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "nonexistent-bucket-id"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Scopes|bucket.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceCollections(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_collections_")
	dsReference := "data.couchbase-capella_collections." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionsDatasourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttr(dsReference, "bucket_id", dmBucketId),
					resource.TestCheckResourceAttr(dsReference, "scope_name", globalScopeName),
					testAccCheckCollectionsNonEmpty(dsReference),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.collection_name"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.max_ttl"),
				),
			},
		},
	})
}

func TestAccDatasourceCollectionsInvalidScope(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_collections_bad_scope_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_collections" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  scope_name      = "nonexistent-scope"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, dmClusterId, dmBucketId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Collections|scope.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceCollectionsInvalidBucket(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_collections_bad_bkt_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_collections" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "nonexistent-bucket-id"
  scope_name      = "_default"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, dmClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Collections|bucket.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceSampleBuckets(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_sample_buckets_")
	dsReference := "data.couchbase-capella_sample_buckets." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccSampleBucketsDatasourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", dmClusterId),
					testAccCheckSampleBucketsReadable(dsReference),
				),
			},
		},
	})
}

func TestAccDatasourceSampleBucketsInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_ds_sample_buckets_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_sample_buckets" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Sample Buckets in Capella|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func testAccBucketResourceConfigRequiredOnly(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, dmClusterId)
}

func testAccBucketResourceConfigAllOptional(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id            = "%[3]s"
  project_id                 = "%[4]s"
  cluster_id                 = "%[5]s"
  name                       = "%[2]s"
  type                       = "couchbase"
  storage_backend            = "couchstore"
  memory_allocation_in_mb    = 200
  bucket_conflict_resolution = "seqno"
  durability_level           = "majority"
  replicas                   = 2
  flush                      = true
  time_to_live_in_seconds    = 0
  eviction_policy            = "fullEviction"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, dmClusterId)
}

func testAccBucketResourceConfigUpdatable(resourceName string, memMB int, durability string, replicas int, ttl int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id         = "%[3]s"
  project_id              = "%[4]s"
  cluster_id              = "%[5]s"
  name                    = "%[2]s"
  memory_allocation_in_mb = %[6]d
  durability_level        = "%[7]s"
  replicas                = %[8]d
  time_to_live_in_seconds = %[9]d
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, dmClusterId, memMB, durability, replicas, ttl)
}

func testAccBucketResourceConfigEphemeral(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  type            = "ephemeral"
  eviction_policy = "nruEviction"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, dmClusterId)
}

func testAccScopeResourceConfig(bucketName, scopeName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
}

resource "couchbase-capella_scope" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "%[6]s"
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, dmClusterId, scopeName)
}

func testAccCollectionResourceConfigNoTTL(bucketName, scopeName, collName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
}

resource "couchbase-capella_scope" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "%[6]s"
}

resource "couchbase-capella_collection" "%[7]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "%[6]s"
  collection_name = "%[7]s"

  depends_on = [couchbase-capella_scope.%[6]s]
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, dmClusterId, scopeName, collName)
}

func testAccCollectionResourceConfigWithTTL(bucketName, scopeName, collName string, ttl int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
}

resource "couchbase-capella_scope" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "%[6]s"
}

resource "couchbase-capella_collection" "%[7]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  max_ttl         = %[8]d

  depends_on = [couchbase-capella_scope.%[6]s]
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, dmClusterId, scopeName, collName, ttl)
}

func testAccCollectionResourceConfigInvalidScope(bucketName, collName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
}

resource "couchbase-capella_collection" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "nonexistent-scope"
  collection_name = "%[6]s"
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, dmClusterId, collName)
}

func testAccFlushBucketResourceConfig(bucketName, flushName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  flush           = true
}

resource "couchbase-capella_flush" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, dmClusterId, flushName)
}

func testAccFlushBucketResourceConfigFlushDisabled(bucketName, flushName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  flush           = false
}

resource "couchbase-capella_flush" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, dmClusterId, flushName)
}

func testAccBucketsDatasourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_buckets" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, dmClusterId)
}

func testAccScopesDatasourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_scopes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, dmClusterId, dmBucketId)
}

func testAccCollectionsDatasourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_collections" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  scope_name      = "%[7]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, dmClusterId, dmBucketId, globalScopeName)
}

func testAccSampleBucketsDatasourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_sample_buckets" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, dmClusterId)
}

func generateBucketImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func generateScopeImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"scope_name=%s,bucket_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["scope_name"], rawState["bucket_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func generateCollectionImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"collection_name=%s,scope_name=%s,bucket_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["collection_name"], rawState["scope_name"], rawState["bucket_id"],
			rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func testAccCheckBucketsNonEmpty(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		count := ds.Primary.Attributes["data.#"]
		if count == "0" || count == "" {
			return fmt.Errorf("datasource %s returned no buckets", dsReference)
		}
		return nil
	}
}

func testAccCheckScopesNonEmpty(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		count := ds.Primary.Attributes["scopes.#"]
		if count == "0" || count == "" {
			return fmt.Errorf("datasource %s returned no scopes", dsReference)
		}
		return nil
	}
}

func testAccCheckCollectionsNonEmpty(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		count := ds.Primary.Attributes["data.#"]
		if count == "0" || count == "" {
			return fmt.Errorf("datasource %s returned no collections", dsReference)
		}
		return nil
	}
}

func testAccCheckSampleBucketsReadable(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.RootModule().Resources[dsReference] == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		return nil
	}
}
