package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceQueryIndexes(t *testing.T) {
	idxName := randomStringWithPrefix("tf_acc_qi_ds_idx_")
	dsName := randomStringWithPrefix("tf_acc_qi_ds_")
	idxReference := "couchbase-capella_query_indexes." + idxName
	dsReference := "data.couchbase-capella_query_indexes." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccQueryIndexesDatasourceConfig(idxName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(idxReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(idxReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(idxReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(idxReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(idxReference, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(idxReference, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(idxReference, "index_name", idxName),

					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(dsReference, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(dsReference, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceQueryIndexesBucketLevelOnly(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_qi_bkt_ds_")
	dsReference := "data.couchbase-capella_query_indexes." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccQueryIndexesDatasourceBucketOnlyConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceQueryIndexesScopeOnly(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_qi_scope_ds_")
	dsReference := "data.couchbase-capella_query_indexes." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccQueryIndexesDatasourceScopeOnlyConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(dsReference, "scope_name", globalScopeName),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceQueryIndexesInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_qi_bad_cluster_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
  bucket_name     = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalBucketName),
				ExpectError: regexp.MustCompile(`(?s)Error Listing Query Indexes|cluster.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceQueryIndexesInvalidBucket(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_qi_bad_bkt_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "nonexistent-bucket-8675309"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)Error Listing Query Indexes|bucket.*not found|No bucket|not found`),
			},
		},
	})
}

func TestAccDatasourceQueryIndexesNonExistentScope(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_qi_bad_scope_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "nonexistent-scope-xyz"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId, globalBucketName),
				ExpectError: regexp.MustCompile(`(?s)Error Listing Query Indexes|scope.*not found|not found`),
			},
		},
	})
}

func TestAccDatasourceQueryIndexesMissingBucketName(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_qi_no_bkt_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)bucket_name|argument.*required|Missing required argument`),
			},
		},
	})
}

func TestAccDatasourceQueryIndexMonitorReady(t *testing.T) {
	idxName := randomStringWithPrefix("tf_acc_qm_idx_")
	monitorName := randomStringWithPrefix("tf_acc_qm_ds_")
	monitorReference := "data.couchbase-capella_query_index_monitor." + monitorName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccQueryIndexMonitorReadyConfig(idxName, monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(monitorReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(monitorReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(monitorReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(monitorReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(monitorReference, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(monitorReference, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(monitorReference, "indexes.#", "1"),
				),
			},
		},
	})
}

func TestAccDatasourceQueryIndexMonitorMultipleIndexes(t *testing.T) {
	idx1Name := randomStringWithPrefix("tf_acc_qm_multi_idx1_")
	idx2Name := randomStringWithPrefix("tf_acc_qm_multi_idx2_")
	monitorName := randomStringWithPrefix("tf_acc_qm_multi_ds_")
	monitorReference := "data.couchbase-capella_query_index_monitor." + monitorName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccQueryIndexMonitorMultipleConfig(idx1Name, idx2Name, monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(monitorReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(monitorReference, "indexes.#", "2"),
				),
			},
		},
	})
}

func TestAccDatasourceQueryIndexMonitorDeferred(t *testing.T) {
	idxName := randomStringWithPrefix("tf_acc_qm_defer_idx_")
	monitorName := randomStringWithPrefix("tf_acc_qm_defer_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccQueryIndexMonitorDeferredConfig(idxName, monitorName),
			},
		},
	})
}

func TestAccDatasourceQueryIndexMonitorNonexistentIndex(t *testing.T) {
	monitorName := randomStringWithPrefix("tf_acc_qm_noexist_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_index_monitor" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  indexes         = ["nonexistent_index_xyz"]
}
`, globalProviderBlock, monitorName, globalOrgId, globalProjectId, globalClusterId,
					globalBucketName, globalScopeName, globalCollectionName),
			},
		},
	})
}

func TestAccDatasourceQueryIndexMonitorInvalidCluster(t *testing.T) {
	monitorName := randomStringWithPrefix("tf_acc_qm_bad_cluster_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_index_monitor" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  indexes         = ["any_index"]
}
`, globalProviderBlock, monitorName, globalOrgId, globalProjectId,
					globalBucketName, globalScopeName, globalCollectionName),
			},
		},
	})
}

func TestAccDatasourceQueryIndexMonitorMissingIndexes(t *testing.T) {
	monitorName := randomStringWithPrefix("tf_acc_qm_no_idx_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_index_monitor" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  indexes         = []
}
`, globalProviderBlock, monitorName, globalOrgId, globalProjectId, globalClusterId,
					globalBucketName, globalScopeName, globalCollectionName),
				ExpectError: regexp.MustCompile(`(?s)indexes|must not be empty|at least one|Set must contain at least`),
			},
		},
	})
}

func testAccQueryIndexesDatasourceConfig(idxName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  index_name      = "%[2]s"
  index_keys      = ["f1"]
  with = {
    defer_build = false
  }
}

data "couchbase-capella_query_indexes" "%[9]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"

  depends_on = [couchbase-capella_query_indexes.%[2]s]
}
`, globalProviderBlock, idxName, globalOrgId, globalProjectId, globalClusterId,
		globalBucketName, globalScopeName, globalCollectionName, dsName)
}

func testAccQueryIndexesDatasourceBucketOnlyConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId, globalBucketName)
}

func testAccQueryIndexesDatasourceScopeOnlyConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId,
		globalBucketName, globalScopeName)
}

func testAccQueryIndexMonitorReadyConfig(idxName, monitorName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  index_name      = "%[2]s"
  index_keys      = ["monitor_field"]
  with = {
    defer_build = false
  }
}

data "couchbase-capella_query_index_monitor" "%[9]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  indexes         = ["%[2]s"]

  depends_on = [couchbase-capella_query_indexes.%[2]s]
}
`, globalProviderBlock, idxName, globalOrgId, globalProjectId, globalClusterId,
		globalBucketName, globalScopeName, globalCollectionName, monitorName)
}

func testAccQueryIndexMonitorMultipleConfig(idx1Name, idx2Name, monitorName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  index_name      = "%[2]s"
  index_keys      = ["m_field1"]
  with = {
    defer_build = false
  }
}

resource "couchbase-capella_query_indexes" "%[9]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  index_name      = "%[9]s"
  index_keys      = ["m_field2"]
  with = {
    defer_build = false
  }
}

data "couchbase-capella_query_index_monitor" "%[10]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  indexes         = ["%[2]s", "%[9]s"]

  depends_on = [
    couchbase-capella_query_indexes.%[2]s,
    couchbase-capella_query_indexes.%[9]s,
  ]
}
`, globalProviderBlock, idx1Name, globalOrgId, globalProjectId, globalClusterId,
		globalBucketName, globalScopeName, globalCollectionName, idx2Name, monitorName)
}

func testAccQueryIndexMonitorDeferredConfig(idxName, monitorName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  index_name      = "%[2]s"
  index_keys      = ["deferred_field"]
  with = {
    defer_build = true
  }
}

data "couchbase-capella_query_index_monitor" "%[9]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = "%[6]s"
  scope_name      = "%[7]s"
  collection_name = "%[8]s"
  indexes         = ["%[2]s"]

  depends_on = [couchbase-capella_query_indexes.%[2]s]
}
`, globalProviderBlock, idxName, globalOrgId, globalProjectId, globalClusterId,
		globalBucketName, globalScopeName, globalCollectionName, monitorName)
}
