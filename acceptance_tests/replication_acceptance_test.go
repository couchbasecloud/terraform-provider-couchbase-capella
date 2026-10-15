package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

// The replication tests replicate from a bucket on the global cluster to a bucket on the DM
// cluster. Each test makes its own buckets, so the tests do not change the shared fixture buckets.
//
// The tests use resource.Test, not resource.ParallelTest: the control plane runs one XDCR change
// at a time between the same two clusters.
//
// All tests use one anchor replication between the two clusters. It keeps
// the network link up for the whole run, so no test has to wait for a new link.
//
// Each test still waits until both clusters are Healthy:
//   - before it makes its buckets,
//   - after it makes its replication,
//   - after it deletes everything.
//
// A new bucket or a new replication can make a cluster busy for a short time, and the next test
// must find the clusters Healthy.

const testAccReplicationReference = "couchbase-capella_replication.replication"

// replicationImportStateVerifyIgnore lists the attributes that import cannot restore, or that can
// change between the create step and the import step.
var replicationImportStateVerifyIgnore = []string{
	"reverse_replication_id",
	"filter.expressions.skip_restream",
	"status",
	"changes_left",
	"error",
	"audit",
}

// TestAccReplicationResource creates a replication with only the required attributes, checks the
// defaults, imports it, and updates it in place.
func TestAccReplicationResource(t *testing.T) {
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckReplicationClustersHealthy(),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					// A create can leave both clusters in the Peering state. Wait before the post-apply refresh.
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttrSet(testAccReplicationReference, "id"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(testAccReplicationReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(testAccReplicationReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrPair(testAccReplicationReference, "source_bucket", "couchbase-capella_bucket.source", "id"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target.cluster", dmClusterId),
					resource.TestCheckResourceAttrPair(testAccReplicationReference, "target.bucket", "couchbase-capella_bucket.target", "id"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target.type", "capella"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "direction", "oneWay"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "priority", "high"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "network_usage_limit", "0"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "reverse_replication_id"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "mappings.#"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "filter.%"),
					resource.TestCheckResourceAttrSet(testAccReplicationReference, "status"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "source_details.cluster.id", globalClusterId),
					resource.TestCheckResourceAttr(testAccReplicationReference, "source_details.bucket.name", buckets.source),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target_details.cluster.id", dmClusterId),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target_details.bucket.name", buckets.target),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target_details.type", "capella"),
					resource.TestCheckResourceAttrSet(testAccReplicationReference, "audit.created_by"),
					resource.TestCheckResourceAttrSet(testAccReplicationReference, "audit.created_at"),
				),
			},
			testAccReplicationImportStep(),
			{
				Config: testAccReplicationConfig(buckets, replicationUpdatedSettings),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccReplicationReference, "priority", "low"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "network_usage_limit", "50"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(testAccReplicationReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(testAccReplicationReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrPair(testAccReplicationReference, "source_bucket", "couchbase-capella_bucket.source", "id"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target.cluster", dmClusterId),
					resource.TestCheckResourceAttrPair(testAccReplicationReference, "target.bucket", "couchbase-capella_bucket.target", "id"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target.type", "capella"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "direction", "oneWay"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "reverse_replication_id"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "mappings.#"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "filter.%"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "source_details.bucket.name", buckets.source),
					resource.TestCheckResourceAttr(testAccReplicationReference, "target_details.bucket.name", buckets.target),
				),
			},
		},
	})
}

// TestAccReplicationResourceMappings creates a replication with scope and collection mappings,
// imports it, and removes the mappings.
func TestAccReplicationResourceMappings(t *testing.T) {
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckReplicationClustersHealthy(),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationMappings),
				Check: resource.ComposeAggregateTestCheckFunc(
					// A create can leave both clusters in the Peering state. Wait before the post-apply refresh.
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationReference, "mappings.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(testAccReplicationReference, "mappings.*", map[string]string{
						"source_scope":                    "sales",
						"target_scope":                    "sales",
						"collections.#":                   "1",
						"collections.0.source_collection": "orders",
						"collections.0.target_collection": "orders",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(testAccReplicationReference, "mappings.*", map[string]string{
						"source_scope": "inventory",
						"target_scope": "inventory",
					}),
				),
			},
			testAccReplicationImportStep(),
			{
				Config: testAccReplicationConfig(buckets, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckNoResourceAttr(testAccReplicationReference, "mappings.#"),
			},
		},
	})
}

// TestAccReplicationResourceFilter creates a replication with a full filter, imports it, keeps only
// the filter expression, and removes the filter.
func TestAccReplicationResourceFilter(t *testing.T) {
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckReplicationClustersHealthy(),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationFullFilter),
				Check: resource.ComposeAggregateTestCheckFunc(
					// A create can leave both clusters in the Peering state. Wait before the post-apply refresh.
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.document_exclude_options.deletion", "true"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.document_exclude_options.expiration", "true"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.document_exclude_options.ttl", "true"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.document_exclude_options.binary", "true"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.expressions.reg_ex", "REGEXP_CONTAINS(META().id, '^order')"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.expressions.skip_restream", "true"),
				),
			},
			testAccReplicationImportStep(),
			{
				Config: testAccReplicationConfig(buckets, replicationExpressionOnlyFilter),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "filter.document_exclude_options.%"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.expressions.reg_ex", "REGEXP_CONTAINS(META().id, '^order')"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "filter.expressions.skip_restream", "false"),
				),
			},
			{
				Config: testAccReplicationConfig(buckets, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckNoResourceAttr(testAccReplicationReference, "filter.%"),
			},
		},
	})
}

// TestAccReplicationResourceTwoWay changes a oneWay replication to twoWay, which replaces it, and
// checks that one delete removes both directions.
func TestAccReplicationResourceTwoWay(t *testing.T) {
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckReplicationClustersHealthy(),
			testAccCheckReplicationDestroyed(t, testAccReplicationReference),
		),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					// A create can leave both clusters in the Peering state. Wait before the post-apply refresh.
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationReference, "direction", "oneWay"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "reverse_replication_id"),
				),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationTwoWay),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationReference, "direction", "twoWay"),
					resource.TestCheckResourceAttrSet(testAccReplicationReference, "reverse_replication_id"),
				),
			},
		},
	})
}

// TestAccReplicationResourceMobileSyncMode creates a replication with mobile sync on, checks that it turns
// on cross cluster versioning for both buckets, imports it, and turns mobile sync off.
func TestAccReplicationResourceMobileSyncMode(t *testing.T) {
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckReplicationClustersHealthy(),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationMobileActive),
				Check: resource.ComposeAggregateTestCheckFunc(
					// A create can leave both clusters in the Peering state. Wait before the post-apply refresh.
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationReference, "mobile_sync_mode", "active"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "conflict_logging.%"),
					testAccCheckReplicationBucketsECCV(t),
				),
			},
			testAccReplicationImportStep(),
			{
				Config: testAccReplicationConfig(buckets, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(testAccReplicationReference, "mobile_sync_mode", "off"),
			},
		},
	})
}

// TestAccReplicationResourceMobileSyncModeUpdateNeedsECCV checks that an update to mobile sync fails
// while cross cluster versioning is off, because only a create turns it on.
func TestAccReplicationResourceMobileSyncModeUpdateNeedsECCV(t *testing.T) {
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckReplicationClustersHealthy(),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					// A create can leave both clusters in the Peering state. Wait before the post-apply refresh.
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationReference, "mobile_sync_mode", "off"),
				),
			},
			{
				Config:      testAccReplicationConfig(buckets, replicationMobileActive),
				ExpectError: regexp.MustCompile("Error updating replication"),
			},
		},
	})
}

// TestAccReplicationResourceConflictLogging creates a replication with conflict logging and every rule
// state, imports it, disables it, and removes it. It needs Couchbase Server 8.0 or later.
func TestAccReplicationResourceConflictLogging(t *testing.T) {
	requireConflictLoggingServerVersion(t)
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(true)
	buckets.log = randomStringWithPrefix("tf_acc_repl_log_")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckReplicationClustersHealthy(),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationMobileActive+fmt.Sprintf(replicationConflictLogging, "false")),
				Check: resource.ComposeAggregateTestCheckFunc(
					// A create can leave both clusters in the Peering state. Wait before the post-apply refresh.
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationReference, "mobile_sync_mode", "active"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "conflict_logging.disabled", "false"),
					resource.TestCheckResourceAttrPair(testAccReplicationReference, "conflict_logging.bucket", "couchbase-capella_bucket.log", "id"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "conflict_logging.collection", "logs.default"),
					resource.TestCheckResourceAttrPair(testAccReplicationReference, "conflict_logging.logging_rules.inventory.bucket", "couchbase-capella_bucket.log", "id"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "conflict_logging.logging_rules.inventory.collection", "logs.inventory"),
					resource.TestCheckNoResourceAttr(testAccReplicationReference, "conflict_logging.logging_rules.sales.bucket"),
				),
			},
			testAccReplicationImportStep(),
			{
				Config: testAccReplicationConfig(buckets, replicationMobileActive+fmt.Sprintf(replicationConflictLogging, "true")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(testAccReplicationReference, "conflict_logging.disabled", "true"),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationMobileActive+replicationConflictLoggingDisabledEmpty),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccReplicationReference, "conflict_logging.disabled", "true"),
					resource.TestCheckResourceAttr(testAccReplicationReference, "conflict_logging.bucket", ""),
					resource.TestCheckResourceAttr(testAccReplicationReference, "conflict_logging.collection", ""),
				),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationMobileActive),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckNoResourceAttr(testAccReplicationReference, "conflict_logging.%"),
			},
		},
	})
}

// TestAccReplicationResourceInvalidConfig checks the plan-time validation. It makes no API call.
func TestAccReplicationResourceInvalidConfig(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccReplicationPlanOnlyConfig(replicationExternalTwoWay),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("must be oneWay"),
			},
			{
				Config:      testAccReplicationPlanOnlyConfig(replicationBucketNameAsId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("source_bucket must be a bucket ID"),
			},
			{
				Config:      testAccReplicationPlanOnlyConfig(replicationConflictLoggingHalfRule),
				PlanOnly:    true,
				// Terraform wraps long error messages, so the pattern allows any white space between the words.
				ExpectError: regexp.MustCompile(`must\s+set\s+both\s+bucket\s+and\s+collection,\s+or\s+neither`),
			},
		},
	})
}

// testAccReplicationImportStep imports the replication and compares the result with the state.
func testAccReplicationImportStep() resource.TestStep {
	return resource.TestStep{
		ResourceName:            testAccReplicationReference,
		ImportState:             true,
		ImportStateIdFunc:       generateReplicationImportIdForResource(testAccReplicationReference),
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: replicationImportStateVerifyIgnore,
	}
}

// testAccCheckReplicationClustersHealthy waits until both clusters are Healthy. It is for Check and
// CheckDestroy, which take a function that returns an error. For PreConfig, use
// testAccPreConfigWaitForReplicationClusters.
func testAccCheckReplicationClustersHealthy() resource.TestCheckFunc {
	return func(*terraform.State) error {
		return waitForReplicationClusters()
	}
}

// testAccPreConfigWaitForReplicationClusters is the same wait for PreConfig. PreConfig takes a
// function with no return value, so this one stops the test with t.Fatalf when the wait fails.
func testAccPreConfigWaitForReplicationClusters(t *testing.T) func() {
	return func() {
		if err := waitForReplicationClusters(); err != nil {
			t.Fatalf("waiting for the replication clusters: %v", err)
		}
	}
}

// requireConflictLoggingServerVersion skips the test when the clusters are older than Couchbase Server 8.0.
func requireConflictLoggingServerVersion(t *testing.T) {
	t.Helper()
	major, _, _ := strings.Cut(globalServerVersion, ".")
	if version, err := strconv.Atoi(major); err != nil || version < 8 {
		t.Skipf("conflict logging needs Couchbase Server 8.0 or later, and TF_VAR_server_version is %q", globalServerVersion)
	}
}

// testAccCheckReplicationBucketsECCV checks that cross cluster versioning is on for the source and target buckets.
func testAccCheckReplicationBucketsECCV(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		data := newTestClient(t)
		for _, resourceReference := range []string{"couchbase-capella_bucket.source", "couchbase-capella_bucket.target"} {
			bucket, ok := state.RootModule().Resources[resourceReference]
			if !ok {
				return fmt.Errorf("resource %q not found in the state", resourceReference)
			}
			attrs := bucket.Primary.Attributes

			cfg := api.EndpointCfg{
				Url: fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s",
					globalHost, attrs["organization_id"], attrs["project_id"], attrs["cluster_id"], attrs["id"]),
				Method:        http.MethodGet,
				SuccessStatus: http.StatusOK,
			}
			response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
			if err != nil {
				return fmt.Errorf("getting bucket %s: %w", attrs["name"], err)
			}

			var got struct {
				EnableCrossClusterVersioning bool `json:"enableCrossClusterVersioning"`
			}
			if err = json.Unmarshal(response.Body, &got); err != nil {
				return err
			}
			if !got.EnableCrossClusterVersioning {
				return fmt.Errorf("cross cluster versioning is off for bucket %s", attrs["name"])
			}
		}
		return nil
	}
}

// testAccCheckReplicationDestroyed checks that both legs are gone after destroy. One DELETE on
// the forward leg removes both legs of a twoWay replication.
func testAccCheckReplicationDestroyed(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		replication, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %q not found in the state captured before destroy", resourceReference)
		}
		attrs := replication.Primary.Attributes

		legs := map[string]string{attrs["cluster_id"]: attrs["id"]}
		if attrs["reverse_replication_id"] != "" {
			legs[attrs["target.cluster"]] = attrs["reverse_replication_id"]
		}

		data := newTestClient(t)
		for clusterId, replicationId := range legs {
			cfg := api.EndpointCfg{
				Url:           replicationapi.ReplicationURL(globalHost, attrs["organization_id"], attrs["project_id"], clusterId, replicationId),
				Method:        http.MethodGet,
				SuccessStatus: http.StatusOK,
			}
			_, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
			if err == nil {
				return fmt.Errorf("replication %s on cluster %s still exists after destroy", replicationId, clusterId)
			}
			if resourceNotFound, errString := api.CheckResourceNotFoundError(err); !resourceNotFound {
				return fmt.Errorf("could not confirm replication %s was destroyed: %s", replicationId, errString)
			}
		}
		return nil
	}
}

// generateReplicationImportIdForResource builds the import ID from the state of the created replication.
func generateReplicationImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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

// replicationTestBuckets holds the names of a test's buckets.
// withScopes indicates whether to include scopes and collections.
// log, when set, is the name of a conflict log bucket on the source cluster.
type replicationTestBuckets struct {
	source     string
	target     string
	log        string
	withScopes bool
}

// newReplicationTestBuckets returns random names for a test's two buckets. Set withScopes for a test
// that needs the scopes and the collection of the mappings test.
func newReplicationTestBuckets(withScopes bool) replicationTestBuckets {
	return replicationTestBuckets{
		source:     randomStringWithPrefix("tf_acc_repl_src_"),
		target:     randomStringWithPrefix("tf_acc_repl_tgt_"),
		withScopes: withScopes,
	}
}

// testAccReplicationBucketsConfig makes a source bucket on the global cluster and a target bucket
// on the DM cluster, with the same settings, and their scopes when the test needs them.
func testAccReplicationBucketsConfig(buckets replicationTestBuckets) string {
	config := fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "source" {
  organization_id            = "%[2]s"
  project_id                 = "%[3]s"
  cluster_id                 = "%[4]s"
  name                       = "%[6]s"
  storage_backend            = "couchstore"
  memory_allocation_in_mb    = 100
  bucket_conflict_resolution = "seqno"
}

resource "couchbase-capella_bucket" "target" {
  organization_id            = "%[2]s"
  project_id                 = "%[3]s"
  cluster_id                 = "%[5]s"
  name                       = "%[7]s"
  storage_backend            = "couchstore"
  memory_allocation_in_mb    = 100
  bucket_conflict_resolution = "seqno"
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, dmClusterId, buckets.source, buckets.target)

	if buckets.withScopes {
		config += testAccReplicationScopesAndCollectionsConfig()
	}
	if buckets.log != "" {
		config += testAccReplicationLogBucketConfig(buckets.log)
	}
	return config
}

// testAccReplicationLogBucketConfig makes a conflict log bucket on the source cluster, with a "logs"
// scope and two collections. A oneWay replication logs on the source cluster only.
func testAccReplicationLogBucketConfig(name string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_bucket" "log" {
  organization_id            = "%[1]s"
  project_id                 = "%[2]s"
  cluster_id                 = "%[3]s"
  name                       = "%[4]s"
  storage_backend            = "couchstore"
  memory_allocation_in_mb    = 100
  bucket_conflict_resolution = "seqno"
}

resource "couchbase-capella_scope" "log_logs" {
  organization_id = couchbase-capella_bucket.log.organization_id
  project_id      = couchbase-capella_bucket.log.project_id
  cluster_id      = couchbase-capella_bucket.log.cluster_id
  bucket_id       = couchbase-capella_bucket.log.id
  scope_name      = "logs"
}

resource "couchbase-capella_collection" "log_default" {
  organization_id = couchbase-capella_scope.log_logs.organization_id
  project_id      = couchbase-capella_scope.log_logs.project_id
  cluster_id      = couchbase-capella_scope.log_logs.cluster_id
  bucket_id       = couchbase-capella_scope.log_logs.bucket_id
  scope_name      = couchbase-capella_scope.log_logs.scope_name
  collection_name = "default"
}

resource "couchbase-capella_collection" "log_inventory" {
  organization_id = couchbase-capella_scope.log_logs.organization_id
  project_id      = couchbase-capella_scope.log_logs.project_id
  cluster_id      = couchbase-capella_scope.log_logs.cluster_id
  bucket_id       = couchbase-capella_scope.log_logs.bucket_id
  scope_name      = couchbase-capella_scope.log_logs.scope_name
  collection_name = "inventory"
}
`, globalOrgId, globalProjectId, globalClusterId, name)
}

// testAccReplicationConfig is testAccReplicationBucketsConfig plus the replication.
func testAccReplicationConfig(buckets replicationTestBuckets, settings string) string {
	return testAccReplicationBucketsConfig(buckets) + fmt.Sprintf(`
resource "couchbase-capella_replication" "replication" {
  organization_id = "%[1]s"
  project_id      = "%[2]s"
  cluster_id      = "%[3]s"
  source_bucket   = couchbase-capella_bucket.source.id

  target = {
    cluster = "%[4]s"
    bucket  = couchbase-capella_bucket.target.id
  }
%[5]s
}
`, globalOrgId, globalProjectId, globalClusterId, dmClusterId, settings)
}

// testAccReplicationScopesAndCollectionsConfig makes the scopes and the collections for the replication test.
func testAccReplicationScopesAndCollectionsConfig() string {
	var config string
	for _, side := range []string{"source", "target"} {
		config += fmt.Sprintf(`
resource "couchbase-capella_scope" "%[1]s_sales" {
  organization_id = couchbase-capella_bucket.%[1]s.organization_id
  project_id      = couchbase-capella_bucket.%[1]s.project_id
  cluster_id      = couchbase-capella_bucket.%[1]s.cluster_id
  bucket_id       = couchbase-capella_bucket.%[1]s.id
  scope_name      = "sales"
}

resource "couchbase-capella_collection" "%[1]s_orders" {
  organization_id = couchbase-capella_scope.%[1]s_sales.organization_id
  project_id      = couchbase-capella_scope.%[1]s_sales.project_id
  cluster_id      = couchbase-capella_scope.%[1]s_sales.cluster_id
  bucket_id       = couchbase-capella_scope.%[1]s_sales.bucket_id
  scope_name      = couchbase-capella_scope.%[1]s_sales.scope_name
  collection_name = "orders"
}

resource "couchbase-capella_scope" "%[1]s_inventory" {
  organization_id = couchbase-capella_bucket.%[1]s.organization_id
  project_id      = couchbase-capella_bucket.%[1]s.project_id
  cluster_id      = couchbase-capella_bucket.%[1]s.cluster_id
  bucket_id       = couchbase-capella_bucket.%[1]s.id
  scope_name      = "inventory"
}
`, side)
	}
	return config
}

// testAccReplicationPlanOnlyConfig is a replication with no bucket resources, for the plan-only
// tests that fail before any API call. body holds the source and the target.
func testAccReplicationPlanOnlyConfig(body string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_replication" "replication" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
%[5]s
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, body)
}

// The settings that the tests put inside the replication block and the
// bodies of the plan-only configs.
const (
	replicationUpdatedSettings = `
  priority            = "low"
  network_usage_limit = 50
`

	replicationMappings = `
  mappings = [
    {
      source_scope = "sales"
      target_scope = "sales"
      collections  = [{ source_collection = "orders", target_collection = "orders" }]
    },
    {
      source_scope = "inventory"
      target_scope = "inventory"
    },
  ]

  depends_on = [
    couchbase-capella_collection.source_orders,
    couchbase-capella_collection.target_orders,
    couchbase-capella_scope.source_inventory,
    couchbase-capella_scope.target_inventory,
  ]
`

	replicationFullFilter = `
  filter = {
    document_exclude_options = {
      deletion   = true
      expiration = true
      ttl        = true
      binary     = true
    }
    expressions = {
      reg_ex        = "REGEXP_CONTAINS(META().id, '^order')"
      skip_restream = true
    }
  }
`

	replicationExpressionOnlyFilter = `
  filter = {
    expressions = {
      reg_ex = "REGEXP_CONTAINS(META().id, '^order')"
    }
  }
`

	replicationTwoWay = `
  direction = "twoWay"
`

	replicationExternalTwoWay = `
  source_bucket = "c291cmNl"
  target = {
    cluster = "my-remote-reference"
    bucket  = "travel-sample"
    type    = "external"
  }
  direction = "twoWay"
`

	replicationBucketNameAsId = `
  source_bucket = "orders"
  target = {
    cluster = "ffffffff-aaaa-1414-eeee-000000000000"
    bucket  = "dGFyZ2V0"
  }
`

	replicationMobileActive = `
  mobile_sync_mode = "active"
`

	replicationConflictLogging = `
  conflict_logging = {
    disabled   = %s
    bucket     = couchbase-capella_bucket.log.id
    collection = "logs.default"

    logging_rules = {
      "inventory"    = { bucket = couchbase-capella_bucket.log.id, collection = "logs.inventory" }
      "sales"        = {}
      "sales.orders" = null
    }
  }

  depends_on = [
    couchbase-capella_collection.source_orders,
    couchbase-capella_collection.target_orders,
    couchbase-capella_scope.source_inventory,
    couchbase-capella_scope.target_inventory,
    couchbase-capella_collection.log_default,
    couchbase-capella_collection.log_inventory,
  ]
`

	replicationConflictLoggingDisabledEmpty = `
  conflict_logging = {
    disabled   = true
    bucket     = ""
    collection = ""
  }
`

	replicationConflictLoggingHalfRule = `
  source_bucket = "c291cmNl"
  target = {
    cluster = "ffffffff-aaaa-1414-eeee-000000000000"
    bucket  = "dGFyZ2V0"
  }
  conflict_logging = {
    disabled      = false
    bucket        = "bG9ncw=="
    collection    = "logs.default"
    logging_rules = { "inventory" = { bucket = "bG9ncw==" } }
  }
`
)
