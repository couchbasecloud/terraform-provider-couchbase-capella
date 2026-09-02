package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

// TestAccCluster_AV_142277 covers the destroy path of the cluster resource.
//
// checkClusterStatus swallowed every error returned while polling the cluster status, including
// the not found error that the remote server returns once the deletion has completed.  As a
// result, destroying a cluster polled for the full 60 minute status timeout and then failed with
// "cluster creation status transition timed out after initiation", leaving the cluster in the
// terraform state.  The destroy performed by the test framework at the end of this test case
// fails if the not found error is swallowed again.
func TestAccCluster_AV_142277(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckClusterDestroyed(t, resourceReference),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceConfigForDestroy(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "current_state", string(clusterapi.Healthy)),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccCheckClusterDestroyed asserts that the cluster recorded in the pre-destroy state has
// been removed from the remote server, so a destroy that reported success cannot have left the
// cluster behind.
func testAccCheckClusterDestroyed(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		clusterResource, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %q not found in the state captured before destroy", resourceReference)
		}
		attrs := clusterResource.Primary.Attributes

		data := newTestClient(t)
		if _, err := retrieveClusterFromServer(data, attrs["organization_id"], attrs["project_id"], attrs["id"]); err != nil {
			if resourceNotFound, errString := api.CheckResourceNotFoundError(err); !resourceNotFound {
				return fmt.Errorf("could not confirm cluster %s was destroyed: %s", attrs["id"], errString)
			}
			return nil
		}

		return fmt.Errorf("cluster %s still exists after destroy", attrs["id"])
	}
}

func testAccClusterResourceConfigForDestroy(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
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
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "single"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cidr)
}
