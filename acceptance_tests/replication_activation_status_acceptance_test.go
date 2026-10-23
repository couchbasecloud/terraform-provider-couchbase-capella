package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

const (
	testAccReplicationForwardStatusReference = "couchbase-capella_replication_activation_status.forward"
	testAccReplicationReverseStatusReference = "couchbase-capella_replication_activation_status.reverse"
)

// TestAccReplicationActivationStatusResource pauses and resumes the two directions of a twoWay
// replication, imports the activation status, and checks that destroy does not change the replication.
func TestAccReplicationActivationStatusResource(t *testing.T) {
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
				Config: testAccReplicationActivationStatusConfig(buckets, "paused", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReplicationClustersHealthy(),
					resource.TestCheckResourceAttr(testAccReplicationForwardStatusReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(testAccReplicationForwardStatusReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(testAccReplicationForwardStatusReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrPair(testAccReplicationForwardStatusReference, "replication_id", testAccReplicationReference, "id"),
					resource.TestCheckResourceAttr(testAccReplicationForwardStatusReference, "state", "paused"),
					testAccCheckReplicationDirectionStatus(t, false, "paused"),
					testAccCheckReplicationDirectionStatus(t, true, "running"),
				),
			},
			{
				ResourceName:                         testAccReplicationForwardStatusReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateReplicationActivationStatusImportId(testAccReplicationForwardStatusReference),
				ImportStateVerifyIdentifierAttribute: "replication_id",
				ImportStateVerify:                    true,
			},
			{
				Config: testAccReplicationActivationStatusConfig(buckets, "running", "paused"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationForwardStatusReference, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(testAccReplicationReverseStatusReference, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccReplicationForwardStatusReference, "state", "running"),
					resource.TestCheckResourceAttr(testAccReplicationReverseStatusReference, "cluster_id", dmClusterId),
					resource.TestCheckResourceAttrPair(testAccReplicationReverseStatusReference, "replication_id", testAccReplicationReference, "reverse_replication_id"),
					resource.TestCheckResourceAttr(testAccReplicationReverseStatusReference, "state", "paused"),
					testAccCheckReplicationDirectionStatus(t, false, "running"),
					testAccCheckReplicationDirectionStatus(t, true, "paused"),
				),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationTwoWay),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccReplicationForwardStatusReference, plancheck.ResourceActionDestroy),
						plancheck.ExpectResourceAction(testAccReplicationReverseStatusReference, plancheck.ResourceActionDestroy),
						plancheck.ExpectResourceAction(testAccReplicationReference, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReplicationDirectionStatus(t, false, "running"),
					testAccCheckReplicationDirectionStatus(t, true, "paused"),
				),
			},
		},
	})
}

// TestAccReplicationActivationStatusResourceInvalidState checks the plan-time validation. It makes no API call.
func TestAccReplicationActivationStatusResourceInvalidState(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_replication_activation_status" "forward" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  replication_id  = "cmVwbGljYXRpb24="
  state           = "stopped"
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`value\s+must\s+be\s+one\s+of`),
			},
		},
	})
}

// testAccCheckReplicationDirectionStatus checks the status of one direction of the twoWay
// replication with a GET. reverse selects the direction that starts on the DM cluster.
func testAccCheckReplicationDirectionStatus(t *testing.T, reverse bool, want string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		replication, ok := state.RootModule().Resources[testAccReplicationReference]
		if !ok {
			return fmt.Errorf("resource %q not found in the state", testAccReplicationReference)
		}
		attrs := replication.Primary.Attributes

		clusterId, replicationId := attrs["cluster_id"], attrs["id"]
		if reverse {
			clusterId, replicationId = attrs["target.cluster"], attrs["reverse_replication_id"]
		}

		cfg := api.EndpointCfg{
			Url:           replicationapi.ReplicationURL(globalHost, globalOrgId, globalProjectId, clusterId, replicationId),
			Method:        http.MethodGet,
			SuccessStatus: http.StatusOK,
		}
		response, err := newTestClient(t).ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
		if err != nil {
			return fmt.Errorf("could not get replication %s on cluster %s: %s", replicationId, clusterId, api.ParseError(err))
		}

		var replicationResp replicationapi.GetReplicationResponse
		if err := json.Unmarshal(response.Body, &replicationResp); err != nil {
			return err
		}

		got := replicationResp.Status
		switch got {
		case replicationapi.StatusPending:
			got = replicationapi.StatusRunning
		case replicationapi.StatusPausing:
			got = replicationapi.StatusPaused
		}
		if string(got) != want {
			return fmt.Errorf("replication %s on cluster %s has status %q, want %q", replicationId, clusterId, replicationResp.Status, want)
		}
		return nil
	}
}

// generateReplicationActivationStatusImportId builds the import ID from the state of the activation status.
func generateReplicationActivationStatusImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		activationStatus, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return "", fmt.Errorf("resource %q not found in the state", resourceReference)
		}
		attrs := activationStatus.Primary.Attributes
		return fmt.Sprintf(
			"replication_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			attrs["replication_id"], attrs["cluster_id"], attrs["project_id"], attrs["organization_id"],
		), nil
	}
}

// testAccReplicationActivationStatusConfig is a twoWay replication with an activation status for each
// direction. An empty state leaves out the resource for that direction.
func testAccReplicationActivationStatusConfig(buckets replicationTestBuckets, forward, reverse string) string {
	config := testAccReplicationConfig(buckets, replicationTwoWay)
	if forward != "" {
		config += fmt.Sprintf(`
resource "couchbase-capella_replication_activation_status" "forward" {
  organization_id = couchbase-capella_replication.replication.organization_id
  project_id      = couchbase-capella_replication.replication.project_id
  cluster_id      = couchbase-capella_replication.replication.cluster_id
  replication_id  = couchbase-capella_replication.replication.id
  state           = "%s"
}
`, forward)
	}
	if reverse != "" {
		config += fmt.Sprintf(`
resource "couchbase-capella_replication_activation_status" "reverse" {
  organization_id = couchbase-capella_replication.replication.organization_id
  project_id      = couchbase-capella_replication.replication.project_id
  cluster_id      = couchbase-capella_replication.replication.target.cluster
  replication_id  = couchbase-capella_replication.replication.reverse_replication_id
  state           = "%s"
}
`, reverse)
	}
	return config
}
