package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testAccProjectReplicationsReference       = "data.couchbase-capella_replications.project"
	testAccSourceClusterReplicationsReference = "data.couchbase-capella_replications.source_cluster"
	testAccTargetClusterReplicationsReference = "data.couchbase-capella_replications.target_cluster"
)

// TestAccReplicationsDataSource lists the replications of a project and of each cluster. A twoWay
// replication shows as two entries in the project list, and each cluster lists only the direction
// that starts on it.
func TestAccReplicationsDataSource(t *testing.T) {
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
				Config: testAccReplicationConfig(buckets, replicationTwoWay),
				Check:  testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, replicationTwoWay) + testAccReplicationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccProjectReplicationsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(testAccProjectReplicationsReference, "project_id", globalProjectId),
					resource.TestCheckNoResourceAttr(testAccProjectReplicationsReference, "cluster_id"),
					resource.TestCheckResourceAttr(testAccSourceClusterReplicationsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(testAccSourceClusterReplicationsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(testAccSourceClusterReplicationsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(testAccTargetClusterReplicationsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(testAccTargetClusterReplicationsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(testAccTargetClusterReplicationsReference, "cluster_id", dmClusterId),
					testAccCheckReplicationsEntry(testAccProjectReplicationsReference, "id", true),
					testAccCheckReplicationsEntry(testAccProjectReplicationsReference, "reverse_replication_id", true),
					testAccCheckReplicationsEntry(testAccSourceClusterReplicationsReference, "id", true),
					testAccCheckReplicationsEntry(testAccSourceClusterReplicationsReference, "reverse_replication_id", false),
					testAccCheckReplicationsEntry(testAccTargetClusterReplicationsReference, "reverse_replication_id", true),
					testAccCheckReplicationsEntry(testAccTargetClusterReplicationsReference, "id", false),
				),
			},
		},
	})
}

// TestAccReplicationsDataSourceInvalidConfig checks the plan-time validation of the IDs.
func TestAccReplicationsDataSourceInvalidConfig(t *testing.T) {
	cases := []struct {
		arguments string
		wantError string
	}{
		{
			arguments: fmt.Sprintf(`project_id = "%s"`, globalProjectId),
			wantError: `The\s+argument\s+"organization_id"\s+is\s+required`,
		},
		{
			arguments: fmt.Sprintf(`organization_id = "%s"`, globalOrgId),
			wantError: `The\s+argument\s+"project_id"\s+is\s+required`,
		},
		{
			arguments: fmt.Sprintf(`
  organization_id = "not-a-uuid"
  project_id      = "%s"`, globalProjectId),
			wantError: `organization_id\s+must\s+be\s+a\s+valid\s+UUID`,
		},
		{
			arguments: fmt.Sprintf(`
  organization_id = "%s"
  project_id      = "not-a-uuid"`, globalOrgId),
			wantError: `project_id\s+must\s+be\s+a\s+valid\s+UUID`,
		},
		{
			arguments: fmt.Sprintf(`
  organization_id = "%s"
  project_id      = "%s"
  cluster_id      = "not-a-uuid"`, globalOrgId, globalProjectId),
			wantError: `cluster_id\s+must\s+be\s+a\s+valid\s+UUID`,
		},
	}

	steps := make([]resource.TestStep, 0, len(cases))
	for _, c := range cases {
		steps = append(steps, resource.TestStep{
			Config: fmt.Sprintf(`
%s

data "couchbase-capella_replications" "invalid" {
  %s
}
`, globalProviderBlock, c.arguments),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(c.wantError),
		})
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

// testAccCheckReplicationsEntry checks whether the list has the replication ID in the idAttribute of the
// replication resource. A listed entry must be there once, with direction twoWay and all fields set.
func testAccCheckReplicationsEntry(dataSourceReference, idAttribute string, listed bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		replication, ok := state.RootModule().Resources[testAccReplicationReference]
		if !ok {
			return fmt.Errorf("resource %q not found in the state", testAccReplicationReference)
		}
		id := replication.Primary.Attributes[idAttribute]
		if id == "" {
			return fmt.Errorf("%s of %s is empty", idAttribute, testAccReplicationReference)
		}

		dataSource, ok := state.RootModule().Resources[dataSourceReference]
		if !ok {
			return fmt.Errorf("data source %q not found in the state", dataSourceReference)
		}
		attrs := dataSource.Primary.Attributes
		count, err := strconv.Atoi(attrs["data.#"])
		if err != nil {
			return fmt.Errorf("%s has no data list: %w", dataSourceReference, err)
		}

		var found []int
		for i := range count {
			if attrs[fmt.Sprintf("data.%d.id", i)] == id {
				found = append(found, i)
			}
		}

		switch {
		case !listed && len(found) > 0:
			return fmt.Errorf("%s lists %s %s, but it starts on a different cluster", dataSourceReference, idAttribute, id)
		case !listed:
			return nil
		case len(found) != 1:
			return fmt.Errorf("%s lists %s %s %d times, want 1", dataSourceReference, idAttribute, id, len(found))
		}

		entry := fmt.Sprintf("data.%d.", found[0])
		if direction := attrs[entry+"direction"]; direction != "twoWay" {
			return fmt.Errorf("%s: entry %s has direction %q, want twoWay", dataSourceReference, id, direction)
		}

		sourceCluster := replication.Primary.Attributes["source_details.cluster.name"]
		targetCluster := replication.Primary.Attributes["target_details.cluster.name"]
		if sourceCluster == "" || targetCluster == "" {
			return fmt.Errorf("%s has no source or target cluster name", testAccReplicationReference)
		}
		if idAttribute == "reverse_replication_id" {
			sourceCluster, targetCluster = targetCluster, sourceCluster
		}
		for field, want := range map[string]string{"source_cluster": sourceCluster, "target_cluster": targetCluster} {
			if got := attrs[entry+field]; got != want {
				return fmt.Errorf("%s: entry %s has %s %q, want %q", dataSourceReference, id, field, got, want)
			}
		}

		for _, field := range []string{"status", "audit.created_by", "audit.created_at"} {
			if attrs[entry+field] == "" {
				return fmt.Errorf("%s: entry %s has no %s", dataSourceReference, id, field)
			}
		}
		return nil
	}
}

func testAccReplicationsDataSourceConfig() string {
	return fmt.Sprintf(`
data "couchbase-capella_replications" "project" {
  organization_id = %[1]s.organization_id
  project_id      = %[1]s.project_id
}

data "couchbase-capella_replications" "source_cluster" {
  organization_id = %[1]s.organization_id
  project_id      = %[1]s.project_id
  cluster_id      = %[1]s.cluster_id
}

data "couchbase-capella_replications" "target_cluster" {
  organization_id = %[1]s.organization_id
  project_id      = %[1]s.project_id
  cluster_id      = %[1]s.target.cluster
}
`, testAccReplicationReference)
}
