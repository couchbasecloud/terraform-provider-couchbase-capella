package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

const (
	cycleBucketCount       = 5
	cycleIndexesPerBucket  = 5
	cycleClusterWait       = 40 * time.Minute
	cycleClusterWaitPoll   = 20 * time.Second
	cycleBucketMemoryInMiB = 100
)

// TestAccGSIKeyspaceReadinessAcrossClusterOnOff_AV_133415 asserts that indexes can be
// created on freshly created buckets immediately after a cluster resume, without a
// sleep — the case a hibernate/resume cycle makes hardest, because the query nodes come
// back with empty metadata caches and so take longest to see a new keyspace.
//
// It is the escalated customer's automation sequence (CBSE-22110) end to end: turn the
// cluster off, on, create 5 ephemeral buckets with 5 indexes each, delete them, turn
// off, and do the buckets and indexes again. Scale matters here — 25 indexes per round
// across 5 keyspaces is what the customer runs, and a single index would not stress the
// same code path.
//
// The narrower single-apply case, without any cluster transition, is
// TestAccGSI_AV_133415 in gsi_keyspace_readiness_acceptance_test.go, which runs on the
// shared cluster in minutes. This one creates a cluster and hibernates it twice, so it
// costs hours — select it deliberately with -run rather than sweeping the package.
//
// The cluster is declared in the config and created by the first step, following
// TestAccFreeTierClusterLifecycle. Its block is byte-identical in every step so it is
// never replaced, and Terraform destroys it at the end.
//
// The steps do not map one-to-one onto the reported sequence, for three reasons that
// are properties of the bug rather than of the test:
//
//   - The cluster has to exist before it can be turned off, so step 1 creates it and
//     the reported sequence starts at step 2.
//
//   - Bucket creation and index creation have to happen in the SAME apply. That is
//     what puts the index DDL inside the window where the query service has not yet
//     cached the new keyspace. Split across two applies the second runs long after
//     propagation and the race is gone, so "create buckets" and "create indexes" are
//     one step here.
//
//   - A cluster transition has to be its own step. couchbase-capella_cluster_onoff_ondemand
//     treats the API's 202 as done and never polls (AV-138686), so a bucket declared in
//     the same apply as a turn-on races it even with depends_on, and fails with 422
//     "Temporarily unavailable while the Cluster is in the Turning Off state." Each
//     transition step therefore ends with waitForClusterState, which runs after apply
//     and so gates the following step.
//
// Step 8 is the interesting one: the query nodes come back with cold metadata caches
// after the step 6 turn-off, giving the widest propagation window — and the best chance
// of exceeding the ~15s the control plane retries for since AV-133415 shipped in 2.2.283.
func TestAccGSIKeyspaceReadinessAcrossClusterOnOff_AV_133415(t *testing.T) {
	clusterName := randomStringWithPrefix("tf_acc_onoff_cluster_")
	onOffName := randomStringWithPrefix("tf_acc_onoff_state_")
	prefix := randomStringWithPrefix("tf_acc_onoff_")
	cidr := generateRandomCIDR()

	clusterReference := "couchbase-capella_cluster." + clusterName
	onOffReference := "couchbase-capella_cluster_onoff_ondemand." + onOffName

	buckets := cycleBucketNames(prefix)

	// Byte-identical across every step so the cluster is never replaced mid-flow.
	clusterBlock := cycleClusterBlock(clusterName, cidr)

	onOff := func(state string) string { return cycleOnOffBlock(clusterName, onOffName, state) }
	bucketsAndIndexes := cycleBucketAndIndexBlocks(clusterName, prefix)

	// resource.Test rather than ParallelTest: a multi-hour run that hibernates a
	// cluster should not interleave with anything else.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// Create the cluster the flow runs on. couchbase-capella_cluster polls
				// to a final state on create, so it is Healthy when this returns.
				Config: cycleConfig(clusterBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(clusterReference, "id"),
					resource.TestCheckResourceAttr(clusterReference, "name", clusterName),
					waitForClusterState(clusterReference, clusterapi.Healthy),
				),
			},
			{
				// 1. turn the cluster off
				Config: cycleConfig(clusterBlock, onOff("off")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onOffReference, "state", "off"),
					waitForClusterState(clusterReference, clusterapi.TurnedOff),
				),
			},
			{
				// 2. turn it on
				Config: cycleConfig(clusterBlock, onOff("on")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onOffReference, "state", "on"),
					waitForClusterState(clusterReference, clusterapi.Healthy),
				),
			},
			{
				// 3 + 4. create the buckets and, in the same apply, their indexes.
				Config: cycleConfig(clusterBlock, onOff("on"), bucketsAndIndexes),
				Check: resource.ComposeAggregateTestCheckFunc(
					append(
						cycleBucketAndIndexChecks(clusterReference, buckets),
						waitForClusterState(clusterReference, clusterapi.Healthy),
					)...,
				),
			},
			{
				// 5. delete the ephemeral buckets. The indexes go with them, and
				// Terraform drops them first because they depend on the buckets.
				Config: cycleConfig(clusterBlock, onOff("on")),
				Check: resource.ComposeAggregateTestCheckFunc(
					cycleCheckBucketsAbsent(clusterReference, buckets),
					waitForClusterState(clusterReference, clusterapi.Healthy),
				),
			},
			{
				// 6. turn the cluster off
				Config: cycleConfig(clusterBlock, onOff("off")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onOffReference, "state", "off"),
					waitForClusterState(clusterReference, clusterapi.TurnedOff),
				),
			},
			{
				// Bring the cluster back up before recreating the buckets. Creating a
				// bucket on a hibernated cluster fails for an unrelated reason, which
				// would mask the keyspace race this test exists for.
				Config: cycleConfig(clusterBlock, onOff("on")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onOffReference, "state", "on"),
					waitForClusterState(clusterReference, clusterapi.Healthy),
				),
			},
			{
				// 7 + 8. buckets and indexes again, this time against query nodes
				// whose caches are cold from the turn-off.
				Config: cycleConfig(clusterBlock, onOff("on"), bucketsAndIndexes),
				Check: resource.ComposeAggregateTestCheckFunc(
					append(
						cycleBucketAndIndexChecks(clusterReference, buckets),
						waitForClusterState(clusterReference, clusterapi.Healthy),
					)...,
				),
			},
		},
	})
}

func cycleConfig(blocks ...string) string {
	return globalProviderBlock + "\n" + strings.Join(blocks, "\n")
}

// cycleClusterBlock is the cluster the flow runs on: its own, not the shared one, because
// the flow hibernates it twice and every other test assumes its cluster stays up. It
// carries data, index and query — index DDL needs the latter two. The caller passes a
// randomised CIDR to avoid colliding with existing clusters in the organization.
func cycleClusterBlock(clusterName, cidr string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_cluster" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[1]s"
  description     = "Customer flow acceptance test (CBSE-22110)."
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[4]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, clusterName, globalOrgId, globalProjectId, cidr)
}

func cycleOnOffBlock(clusterName, onOffName, state string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_cluster_onoff_ondemand" "%[1]s" {
  organization_id            = "%[2]s"
  project_id                 = "%[3]s"
  cluster_id                 = couchbase-capella_cluster.%[4]s.id
  state                      = "%[5]s"
  turn_on_linked_app_service = false
}
`, onOffName, globalOrgId, globalProjectId, clusterName, state)
}

func cycleBucketNames(prefix string) []string {
	names := make([]string, 0, cycleBucketCount)
	for i := 1; i <= cycleBucketCount; i++ {
		names = append(names, fmt.Sprintf("%s_%d", prefix, i))
	}
	return names
}

func cycleIndexNames(bucket string) []string {
	names := make([]string, 0, cycleIndexesPerBucket)
	for i := 1; i <= cycleIndexesPerBucket; i++ {
		names = append(names, fmt.Sprintf("%s_idx_%d", bucket, i))
	}
	return names
}

// cycleBucketAndIndexBlocks declares the buckets and their indexes together. bucket_name
// reads through the bucket resource rather than repeating the literal name: that
// reference is what orders the DDL after the bucket create, and puts it immediately
// after, with nothing waiting in between.
//
// The blocks are emitted one per resource rather than with for_each, which would be the
// natural way to write this. GSI.ValidateConfig guards only on organization_id being
// unknown and then tests index_name with ValueString() == ""; under for_each, index_name
// is still unknown when validation runs and ValueString() yields "", so it rejects the
// config at plan time with "Expected index_name to be configured but is null" even
// though it is set (AV-138687). Once that is fixed this can collapse to a for_each over
// a map of index specs.
//
// The indexes are created deferred and built per bucket with one BUILD INDEX. Five
// immediate builds on one keyspace get rejected by the indexer as concurrent builds,
// which the provider surfaces as a warning and leaves the resource unpopulated — that
// would fail the attribute checks for reasons unrelated to the keyspace race. The
// CREATE still has to resolve the keyspace, so the race stands.
func cycleBucketAndIndexBlocks(clusterName, prefix string) string {
	var b strings.Builder

	for i, bucket := range cycleBucketNames(prefix) {
		fmt.Fprintf(&b, `
resource "couchbase-capella_bucket" "%[1]s" {
  organization_id         = "%[2]s"
  project_id              = "%[3]s"
  cluster_id              = couchbase-capella_cluster.%[4]s.id
  name                    = "%[5]s"
  type                    = "ephemeral"
  eviction_policy         = "nruEviction"
  memory_allocation_in_mb = %[6]d
}
`, cycleBucketResourceName(i), globalOrgId, globalProjectId, clusterName, bucket, cycleBucketMemoryInMiB)

		for j, index := range cycleIndexNames(bucket) {
			fmt.Fprintf(&b, `
resource "couchbase-capella_query_indexes" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[4]s.id
  bucket_name     = couchbase-capella_bucket.%[5]s.name
  scope_name      = "_default"
  collection_name = "_default"
  index_name      = "%[6]s"
  index_keys      = ["c%[7]d"]

  with = {
    defer_build = true
  }
}
`, cycleIndexResourceName(i, j), globalOrgId, globalProjectId, clusterName,
				cycleBucketResourceName(i), index, j+1)
		}

		fmt.Fprintf(&b, `
resource "couchbase-capella_query_indexes" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[4]s.id
  bucket_name     = couchbase-capella_bucket.%[5]s.name
  scope_name      = "_default"
  collection_name = "_default"
  build_indexes   = %[6]s

  depends_on = [%[7]s]
}
`, cycleBuildResourceName(i), globalOrgId, globalProjectId, clusterName,
			cycleBucketResourceName(i), cycleHCLStringList(cycleIndexNames(bucket)), cycleIndexResourceRefs(i))
	}

	return b.String()
}

func cycleBucketResourceName(bucketIdx int) string {
	return fmt.Sprintf("bkt_%d", bucketIdx+1)
}

func cycleIndexResourceName(bucketIdx, indexIdx int) string {
	return fmt.Sprintf("idx_%d_%d", bucketIdx+1, indexIdx+1)
}

func cycleBuildResourceName(bucketIdx int) string {
	return fmt.Sprintf("build_%d", bucketIdx+1)
}

func cycleHCLStringList(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("%q", value))
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func cycleIndexResourceRefs(bucketIdx int) string {
	refs := make([]string, 0, cycleIndexesPerBucket)
	for j := 0; j < cycleIndexesPerBucket; j++ {
		refs = append(refs, "couchbase-capella_query_indexes."+cycleIndexResourceName(bucketIdx, j))
	}
	return strings.Join(refs, ", ")
}

// resourceAttrFromState reads one attribute of a resource out of Terraform state, for
// checks that need a value only known after apply — here the generated cluster ID.
func resourceAttrFromState(state *terraform.State, resourceRef, attr string) (string, error) {
	for _, module := range state.Modules {
		rs, ok := module.Resources[resourceRef]
		if !ok {
			continue
		}
		value := rs.Primary.Attributes[attr]
		if value == "" {
			return "", fmt.Errorf("%s has no %q in state", resourceRef, attr)
		}
		return value, nil
	}
	return "", fmt.Errorf("%s not found in state", resourceRef)
}

// waitForClusterState polls until the cluster reaches want. It is used as the last
// check of a step so that the following step starts against a settled cluster: the
// on/off resource returns on the API's 202 without waiting (AV-138686), and a bucket
// create against a turningOn/turningOff cluster fails with 422. It also absorbs the
// rebalance each bucket create kicks off.
func waitForClusterState(clusterResourceRef string, want clusterapi.State) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		clusterID, err := resourceAttrFromState(state, clusterResourceRef, "id")
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), cycleClusterWait)
		defer cancel()

		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s",
			globalHost, globalOrgId, globalProjectId, clusterID)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

		var last clusterapi.State
		for {
			response, err := globalClient.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
			if err != nil {
				return fmt.Errorf("reading cluster %s while waiting for %s: %w", clusterID, want, err)
			}

			var cluster clusterapi.GetClusterResponse
			if err := json.Unmarshal(response.Body, &cluster); err != nil {
				return fmt.Errorf("unmarshalling cluster %s: %w", clusterID, err)
			}
			last = cluster.CurrentState

			if last == want {
				return nil
			}
			// Failed transitions never resolve, so give up rather than spend the
			// whole timeout on a cluster that will not move.
			if strings.HasSuffix(string(last), "Failed") {
				return fmt.Errorf("cluster %s entered terminal state %q while waiting for %q",
					clusterID, last, want)
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("timed out after %s waiting for cluster %s to reach %q, last state %q",
					cycleClusterWait, clusterID, want, last)
			case <-time.After(cycleClusterWaitPoll):
			}
		}
	}
}

func cycleBucketAndIndexChecks(clusterReference string, buckets []string) []resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0, len(buckets)*(cycleIndexesPerBucket+1))

	for i, bucket := range buckets {
		bucketReference := "couchbase-capella_bucket." + cycleBucketResourceName(i)
		checks = append(checks,
			resource.TestCheckResourceAttr(bucketReference, "name", bucket),
			resource.TestCheckResourceAttrPair(bucketReference, "cluster_id", clusterReference, "id"),
			resource.TestCheckResourceAttr(bucketReference, "type", "ephemeral"),
		)

		for j, index := range cycleIndexNames(bucket) {
			indexReference := "couchbase-capella_query_indexes." + cycleIndexResourceName(i, j)
			checks = append(checks,
				resource.TestCheckResourceAttr(indexReference, "bucket_name", bucket),
				resource.TestCheckResourceAttr(indexReference, "index_name", index),
				resource.TestCheckResourceAttr(indexReference, "scope_name", "_default"),
				resource.TestCheckResourceAttr(indexReference, "collection_name", "_default"),
			)
		}
	}

	return checks
}

// cycleCheckBucketsAbsent confirms the buckets are gone from Capella, not merely dropped
// from state, so the delete step actually proves something.
func cycleCheckBucketsAbsent(clusterReference string, buckets []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		clusterID, err := resourceAttrFromState(state, clusterReference, "id")
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
			globalHost, globalOrgId, globalProjectId, clusterID)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

		remaining, err := api.GetPaginated[[]struct {
			Name string `json:"name"`
		}](ctx, globalClient, globalToken, cfg, api.SortById)
		if err != nil {
			return fmt.Errorf("listing buckets on %s: %w", clusterID, err)
		}

		present := make(map[string]struct{}, len(remaining))
		for _, bucket := range remaining {
			present[bucket.Name] = struct{}{}
		}

		var leftover []string
		for _, bucket := range buckets {
			if _, ok := present[bucket]; ok {
				leftover = append(leftover, bucket)
			}
		}
		if len(leftover) > 0 {
			return fmt.Errorf("buckets still present after delete: %s", strings.Join(leftover, ", "))
		}
		return nil
	}
}
