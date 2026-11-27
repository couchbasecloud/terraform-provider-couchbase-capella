package acceptance_tests

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const sparseVectorKey = "`sv` SPARSE VECTOR"

func TestAccSparseVectorIndexDefaultTunables(t *testing.T) {
	requireSparseVectorSupport(t)

	idxName := randomStringWithPrefix("tf_acc_sv_idx_")
	ref := "couchbase-capella_query_indexes." + idxName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccSparseVectorIndexResource(idxName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(ref, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(ref, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(ref, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(ref, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(ref, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(ref, "index_name", idxName),
					resource.TestCheckResourceAttr(ref, "index_keys.0", sparseVectorKey),
					resource.TestCheckResourceAttr(ref, "status", "Ready"),
				),
			},
		},
	})
}

func TestAccSparseVectorIndexMultiple(t *testing.T) {
	requireSparseVectorSupport(t)

	idx1 := randomStringWithPrefix("tf_acc_sv_multi1_")
	idx2 := randomStringWithPrefix("tf_acc_sv_multi2_")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccSparseVectorIndexResource(idx1, "") + testAccSparseVectorIndexResource(idx2, "depends_on = [couchbase-capella_query_indexes."+idx1+"]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes."+idx1, "status", "Ready"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes."+idx2, "status", "Ready"),
				),
			},
		},
	})
}

func TestAccSparseVectorIndexCompositeWithWhere(t *testing.T) {
	requireSparseVectorSupport(t)

	idxName := randomStringWithPrefix("tf_acc_sv_comp_")
	ref := "couchbase-capella_query_indexes." + idxName
	body := `
  index_keys = ["` + "`type`" + `", "` + sparseVectorKey + `"]
  where      = "` + "`type`" + ` = 'sv_seed'"`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccSparseVectorIndexResourceRaw(idxName, body),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "index_keys.#", "2"),
					resource.TestCheckResourceAttr(ref, "index_keys.0", "`type`"),
					resource.TestCheckResourceAttr(ref, "index_keys.1", sparseVectorKey),
					resource.TestCheckResourceAttr(ref, "where", "`type` = 'sv_seed'"),
					resource.TestCheckResourceAttr(ref, "status", "Ready"),
				),
			},
		},
	})
}

func TestAccSparseVectorIndexReplicaUpdate(t *testing.T) {
	requireSparseVectorSupport(t)

	idxName := randomStringWithPrefix("tf_acc_sv_rep_")
	ref := "couchbase-capella_query_indexes." + idxName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccSparseVectorIndexResource(idxName, "with = {\n    num_replica = 1\n  }"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "with.num_replica", "1"),
					resource.TestCheckResourceAttr(ref, "status", "Ready"),
				),
			},
			{
				Config: globalProviderBlock + testAccSparseVectorIndexResource(idxName, "with = {\n    num_replica = 2\n  }"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "with.num_replica", "2"),
					resource.TestCheckResourceAttr(ref, "status", "Ready"),
				),
			},
		},
	})
}

func TestAccSparseVectorIndexPartitioned(t *testing.T) {
	requireSparseVectorSupport(t)

	idxName := randomStringWithPrefix("tf_acc_sv_part_")
	ref := "couchbase-capella_query_indexes." + idxName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccSparseVectorIndexResource(idxName, "partition_by = [\"meta().id\"]\n  with = {\n    num_partition = 4\n  }"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "partition_by.0", "meta().id"),
					resource.TestCheckResourceAttr(ref, "with.num_partition", "4"),
					resource.TestCheckResourceAttr(ref, "status", "Ready"),
				),
			},
		},
	})
}

func TestAccSparseVectorIndexDeferredBuild(t *testing.T) {
	requireSparseVectorSupport(t)

	idxName := randomStringWithPrefix("tf_acc_sv_def_")
	builderName := randomStringWithPrefix("tf_acc_sv_build_")
	ref := "couchbase-capella_query_indexes." + idxName
	indexConfig := globalProviderBlock + testAccSparseVectorIndexResource(idxName, "with = {\n    defer_build = true\n  }")
	builtConfig := indexConfig + fmt.Sprintf(`
resource "couchbase-capella_query_indexes" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  build_indexes   = [couchbase-capella_query_indexes.%[8]s.index_name]
}
`, builderName, globalOrgId, globalProjectId, globalClusterId, globalBucketName, globalScopeName, globalCollectionName, idxName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: indexConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "with.defer_build", "true"),
					resource.TestCheckResourceAttr(ref, "index_keys.0", sparseVectorKey),
				),
			},
			{
				Config: builtConfig,
				Check:  resource.TestCheckResourceAttr("couchbase-capella_query_indexes."+builderName, "build_indexes.#", "1"),
			},
			{
				Config: builtConfig,
				Check:  resource.TestCheckResourceAttr(ref, "status", "Ready"),
			},
		},
	})
}

func TestAccSparseVectorIndexImport(t *testing.T) {
	requireSparseVectorSupport(t)

	idxName := randomStringWithPrefix("tf_acc_sv_imp_")
	ref := "couchbase-capella_query_indexes." + idxName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: globalProviderBlock + testAccSparseVectorIndexResource(idxName, ""),
				Check:  resource.TestCheckResourceAttr(ref, "status", "Ready"),
			},
			{
				ResourceName:      ref,
				ImportStateIdFunc: generateGsiImportIdForResource(ref),
				ImportState:       true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					if got := states[0].Attributes["index_keys.0"]; got != sparseVectorKey {
						return fmt.Errorf("import changed index_keys.0 to %q, want %q", got, sparseVectorKey)
					}
					return nil
				},
			},
		},
	})
}

func TestAccSparseVectorIndexDatasource(t *testing.T) {
	requireSparseVectorSupport(t)

	idxName := randomStringWithPrefix("tf_acc_sv_ds_idx_")
	dsName := randomStringWithPrefix("tf_acc_sv_ds_")
	dsRef := "data.couchbase-capella_query_indexes." + dsName
	config := globalProviderBlock + testAccSparseVectorIndexResource(idxName, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr("couchbase-capella_query_indexes."+idxName, "status", "Ready"),
			},
			{
				PreConfig: func() { time.Sleep(queryIndexCatalogPropagationDelay) },
				Config: config + fmt.Sprintf(`
data "couchbase-capella_query_indexes" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"

  depends_on = [couchbase-capella_query_indexes.%[8]s]
}
`, dsName, globalOrgId, globalProjectId, globalClusterId, globalBucketName, globalScopeName, globalCollectionName, idxName),
				Check: resource.TestCheckTypeSetElemNestedAttrs(dsRef, "data.*", map[string]string{
					"index_name": idxName,
				}),
			},
		},
	})
}

func requireSparseVectorSupport(t *testing.T) {
	t.Helper()
	parts := strings.SplitN(globalServerVersion, ".", 2)
	major, _ := strconv.Atoi(parts[0])
	minor := 0
	if len(parts) == 2 {
		minor, _ = strconv.Atoi(parts[1])
	}
	if major < 8 || (major == 8 && minor < 5) {
		t.Skipf("sparse vector indexes need Couchbase Server 8.5+, got %s; set TF_VAR_server_version=8.5", globalServerVersion)
	}
	if err := ensureSparseVectorSeedDoc(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func testAccSparseVectorIndexResource(name, extra string) string {
	return testAccSparseVectorIndexResourceRaw(name, `
  index_keys = ["`+"`sv`"+` SPARSE VECTOR"]
  `+extra)
}

func testAccSparseVectorIndexResourceRaw(name, body string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_query_indexes" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "%[1]s"
%[8]s
}
`, name, globalOrgId, globalProjectId, globalClusterId, globalBucketName, globalScopeName, globalCollectionName, body)
}
