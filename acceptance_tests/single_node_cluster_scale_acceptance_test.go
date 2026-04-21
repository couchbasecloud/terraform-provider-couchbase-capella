package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSingleNodeClusterScaleWithTravelSample is an acceptance test that validates
// the full lifecycle of:
//  1. Creating a single-node AWS cluster (availability=single, num_of_nodes=1).
//  2. Importing the travel-sample bucket into that cluster.
//  3. Scaling the cluster in-place from 1 node to 3 nodes without destroying/recreating it.
//
// This test ensures that scaling is an in-place update and not a destroy+create.
func TestAccSingleNodeClusterScaleWithTravelSample(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_cluster_")
	clusterResourceRef := "couchbase-capella_cluster." + clusterResourceName

	sampleBucketResourceName := randomStringWithPrefix("tf_acc_sample_")
	sampleBucketResourceRef := "couchbase-capella_sample_bucket." + sampleBucketResourceName

	cidr := generateRandomCIDR()

	// clusterID is captured after the cluster is created and verified unchanged
	// after the scale step to prove scaling is an in-place update (no destroy+create).
	var clusterID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Create a single-node cluster
			{
				Config: testAccSingleNodeClusterConfig(clusterResourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterResourceRef),
					resource.TestCheckResourceAttr(clusterResourceRef, "name", clusterResourceName),
					resource.TestCheckResourceAttr(clusterResourceRef, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(clusterResourceRef, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(clusterResourceRef, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.num_of_nodes", "1"),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.services.#", "1"),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.node.disk.iops", "3000"),
					resource.TestCheckResourceAttr(clusterResourceRef, "availability.type", "single"),
					resource.TestCheckResourceAttr(clusterResourceRef, "support.plan", "basic"),
					resource.TestCheckResourceAttr(clusterResourceRef, "support.timezone", "PT"),
					// Capture the cluster ID for later comparison
					resource.TestCheckResourceAttrWith(clusterResourceRef, "id", func(value string) error {
						clusterID = value
						return nil
					}),
				),
			},
			// Step 2: Import travel-sample into the cluster
			{
				Config: testAccSingleNodeClusterWithTravelSampleConfig(clusterResourceName, sampleBucketResourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Cluster still exists and unchanged
					testAccExistsClusterResource(t, clusterResourceRef),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.num_of_nodes", "1"),
					// travel-sample bucket was created
					resource.TestCheckResourceAttr(sampleBucketResourceRef, "name", "travel-sample"),
					resource.TestCheckResourceAttrSet(sampleBucketResourceRef, "id"),
				),
			},
			// Step 3: Scale cluster from 1 node to 3 nodes (in-place update, no destroy)
			{
				Config: testAccScaledClusterWithTravelSampleConfig(clusterResourceName, sampleBucketResourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterResourceRef),
					resource.TestCheckResourceAttr(clusterResourceRef, "service_groups.0.num_of_nodes", "3"),
					// Immutable fields remain unchanged
					resource.TestCheckResourceAttr(clusterResourceRef, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(clusterResourceRef, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(clusterResourceRef, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(clusterResourceRef, "availability.type", "single"),
					resource.TestCheckResourceAttr(clusterResourceRef, "support.plan", "basic"),
					resource.TestCheckResourceAttr(clusterResourceRef, "support.timezone", "PT"),
					// Verify the cluster ID is unchanged, proving in-place update (no destroy+create)
					resource.TestCheckResourceAttrWith(clusterResourceRef, "id", func(value string) error {
						if value != clusterID {
							return fmt.Errorf("cluster was recreated: ID changed from %s to %s", clusterID, value)
						}
						return nil
					}),
					// travel-sample bucket still exists
					resource.TestCheckResourceAttr(sampleBucketResourceRef, "name", "travel-sample"),
					resource.TestCheckResourceAttrSet(sampleBucketResourceRef, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccSingleNodeClusterConfig returns the Terraform config for a single-node AWS cluster.
func testAccSingleNodeClusterConfig(clusterName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Single node acceptance test cluster"

  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
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
      num_of_nodes = 1
      services     = ["data"]
    }
  ]

  availability = {
    "type" : "single"
  }

  support = {
    plan     = "basic"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, cidr)
}

// testAccSingleNodeClusterWithTravelSampleConfig returns the Terraform config for a
// single-node cluster plus the travel-sample bucket resource.
func testAccSingleNodeClusterWithTravelSampleConfig(clusterName, sampleBucketName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Single node acceptance test cluster"

  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
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
      num_of_nodes = 1
      services     = ["data"]
    }
  ]

  availability = {
    "type" : "single"
  }

  support = {
    plan     = "basic"
    timezone = "PT"
  }
}

resource "couchbase-capella_sample_bucket" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[4]s.id
  name            = "travel-sample"
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, cidr, sampleBucketName)
}

// testAccScaledClusterWithTravelSampleConfig returns the Terraform config with the
// cluster scaled from 1 to 3 nodes. All immutable fields (availability, cloud_provider,
// zones) remain unchanged so Terraform performs an in-place update instead of destroy+create.
func testAccScaledClusterWithTravelSampleConfig(clusterName, sampleBucketName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Single node acceptance test cluster"

  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
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
      services     = ["data"]
    }
  ]

  availability = {
    "type" : "single"
  }

  support = {
    plan     = "basic"
    timezone = "PT"
  }
}

resource "couchbase-capella_sample_bucket" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[4]s.id
  name            = "travel-sample"
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, cidr, sampleBucketName)
}
