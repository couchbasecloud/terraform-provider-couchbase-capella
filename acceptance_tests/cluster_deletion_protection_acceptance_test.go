package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

// TestAccClusterDeletionProtectionResource tests the full lifecycle:
// create with protection disabled → update to enabled → import state.
func TestAccClusterDeletionProtectionResource(t *testing.T) {
	disableDeletionProtectionOnCleanup(t)
	resourceName := randomStringWithPrefix("tf_acc_del_prot_")
	resourceReference := "couchbase-capella_cluster_deletion_protection." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Create with deletion_protection = false
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "false"),
				),
			},
			// Step 2: Update from false to true
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "true"),
				),
			},
			// Step 3: ImportState
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateDeletionProtectionImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
			// Step 4: Disable to leave cluster in clean state
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "false"),
				),
			},
		},
	})
}

// TestAccClusterDeletionProtectionInvalidCluster verifies correct error for nonexistent cluster.
func TestAccClusterDeletionProtectionInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_del_prot_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterDeletionProtectionConfig(resourceName, "00000000-0000-0000-0000-000000000000", true),
				ExpectError: regexp.MustCompile(`(?s)Error.*cluster.*(not be found|not found|access.*denied)`),
			},
		},
	})
}

func TestAccClusterDeletionProtectionProtectedDestroyFails(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_cluster_del_prot_")
	deletionProtectionResourceName := randomStringWithPrefix("tf_acc_del_prot_destroy_")
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	deletionProtectionReference := "couchbase-capella_cluster_deletion_protection." + deletionProtectionResourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDeletionProtectionProtectedClusterConfig(clusterResourceName, deletionProtectionResourceName, cidr, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					registerDeletionProtectionCleanupForClusterResource(t, clusterReference),
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(clusterReference, "name", clusterResourceName),
					resource.TestCheckResourceAttrSet(clusterReference, "id"),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "true"),
				),
			},
			{
				Config:      testAccClusterDeletionProtectionProtectedClusterConfig(clusterResourceName, deletionProtectionResourceName, cidr, true),
				Destroy:     true,
				ExpectError: regexp.MustCompile("Cluster deletion protection is enabled"),
			},
			{
				Config: testAccClusterDeletionProtectionProtectedClusterConfig(clusterResourceName, deletionProtectionResourceName, cidr, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "false"),
				),
			},
			{
				Config: testAccClusterDeletionProtectionProtectedClusterConfig(clusterResourceName, deletionProtectionResourceName, cidr, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(clusterReference, "deletion_protection", "false"),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "false"),
				),
			},
		},
	})
}

func TestAccClusterDeletionProtectionProtectedBucketDeleteFails(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_cluster_del_prot_bucket_")
	bucketResourceName := randomStringWithPrefix("tf_acc_bucket_del_prot_")
	deletionProtectionResourceName := randomStringWithPrefix("tf_acc_del_prot_bucket_")
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	bucketReference := "couchbase-capella_bucket." + bucketResourceName
	deletionProtectionReference := "couchbase-capella_cluster_deletion_protection." + deletionProtectionResourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDeletionProtectionProtectedBucketConfig(clusterResourceName, bucketResourceName, deletionProtectionResourceName, cidr, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					registerDeletionProtectionCleanupForClusterResource(t, clusterReference),
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(clusterReference, "name", clusterResourceName),
					resource.TestCheckResourceAttr(bucketReference, "name", bucketResourceName),
					resource.TestCheckResourceAttrSet(bucketReference, "id"),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "true"),
				),
			},
			{
				Config:      testAccClusterDeletionProtectionProtectedBucketConfig(clusterResourceName, bucketResourceName, deletionProtectionResourceName, cidr, true, false),
				ExpectError: regexp.MustCompile(`(?s)Unable to delete bucket.*disable cluster\s+deletion protection`),
			},
			{
				Config: testAccClusterDeletionProtectionProtectedBucketConfig(clusterResourceName, bucketResourceName, deletionProtectionResourceName, cidr, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(bucketReference, "name", bucketResourceName),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "false"),
				),
			},
			{
				Config: testAccClusterDeletionProtectionProtectedBucketConfig(clusterResourceName, bucketResourceName, deletionProtectionResourceName, cidr, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(clusterReference, "deletion_protection", "false"),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "false"),
				),
			},
		},
	})
}

func TestAccClusterDeletionProtectionProtectedBucketFlushFails(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_cluster_del_prot_flush_")
	bucketResourceName := randomStringWithPrefix("tf_acc_bucket_flush_prot_")
	flushResourceName := randomStringWithPrefix("tf_acc_flush_del_prot_")
	deletionProtectionResourceName := randomStringWithPrefix("tf_acc_del_prot_flush_")
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	bucketReference := "couchbase-capella_bucket." + bucketResourceName
	deletionProtectionReference := "couchbase-capella_cluster_deletion_protection." + deletionProtectionResourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDeletionProtectionProtectedBucketFlushConfig(clusterResourceName, bucketResourceName, flushResourceName, deletionProtectionResourceName, cidr, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					registerDeletionProtectionCleanupForClusterResource(t, clusterReference),
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(clusterReference, "name", clusterResourceName),
					resource.TestCheckResourceAttr(bucketReference, "name", bucketResourceName),
					resource.TestCheckResourceAttr(bucketReference, "flush", "true"),
					resource.TestCheckResourceAttrSet(bucketReference, "id"),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "true"),
				),
			},
			{
				Config:      testAccClusterDeletionProtectionProtectedBucketFlushConfig(clusterResourceName, bucketResourceName, flushResourceName, deletionProtectionResourceName, cidr, true, true),
				ExpectError: regexp.MustCompile(`(?is)Error flushing the bucket.*deletion.*protection`),
			},
			{
				Config: testAccClusterDeletionProtectionProtectedBucketFlushConfig(clusterResourceName, bucketResourceName, flushResourceName, deletionProtectionResourceName, cidr, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(bucketReference, "name", bucketResourceName),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "false"),
				),
			},
			{
				Config: testAccClusterDeletionProtectionProtectedBucketFlushConfig(clusterResourceName, bucketResourceName, flushResourceName, deletionProtectionResourceName, cidr, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(clusterReference, "deletion_protection", "false"),
					resource.TestCheckResourceAttr(deletionProtectionReference, "deletion_protection", "false"),
				),
			},
		},
	})
}

// --- Config builders ---

func testAccClusterDeletionProtectionConfig(resourceName, clusterID string, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_deletion_protection" "%[2]s" {
  organization_id     = "%[3]s"
  project_id          = "%[4]s"
  cluster_id          = "%[5]s"
  deletion_protection = %[6]t
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, clusterID, enabled)
}

func testAccClusterDeletionProtectionProtectedClusterConfig(clusterResourceName, deletionProtectionResourceName, cidr string, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	name            = "%[4]s"
	cloud_provider = {
		type   = "aws"
		region = "us-east-1"
		cidr   = "%[6]s"
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
		plan     = "developer pro"
		timezone = "PT"
	}
}

resource "couchbase-capella_cluster_deletion_protection" "%[5]s" {
	organization_id     = "%[2]s"
	project_id          = "%[3]s"
	cluster_id          = couchbase-capella_cluster.%[4]s.id
	deletion_protection = %[7]t

	depends_on = [couchbase-capella_cluster.%[4]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterResourceName, deletionProtectionResourceName, cidr, enabled)
}

func testAccClusterDeletionProtectionProtectedBucketConfig(clusterResourceName, bucketResourceName, deletionProtectionResourceName, cidr string, enabled, includeBucket bool) string {
	bucketConfig := ""
	if includeBucket {
		bucketConfig = fmt.Sprintf(`
resource "couchbase-capella_bucket" %[1]q {
	organization_id = %[2]q
	project_id      = %[3]q
	cluster_id      = couchbase-capella_cluster.%[4]s.id
	name            = %[1]q

	depends_on = [couchbase-capella_cluster.%[4]s]
}
`, bucketResourceName, globalOrgId, globalProjectId, clusterResourceName)
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	name            = "%[4]s"
	cloud_provider = {
		type   = "aws"
		region = "us-east-1"
		cidr   = "%[7]s"
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
		plan     = "developer pro"
		timezone = "PT"
	}
}

%[8]s
resource "couchbase-capella_cluster_deletion_protection" "%[6]s" {
	organization_id     = "%[2]s"
	project_id          = "%[3]s"
	cluster_id          = couchbase-capella_cluster.%[4]s.id
	deletion_protection = %[5]t

	depends_on = [couchbase-capella_cluster.%[4]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterResourceName, enabled, deletionProtectionResourceName, cidr, bucketConfig)
}

func testAccClusterDeletionProtectionProtectedBucketFlushConfig(clusterResourceName, bucketResourceName, flushResourceName, deletionProtectionResourceName, cidr string, enabled, includeFlush bool) string {
	flushConfig := ""
	if includeFlush {
		flushConfig = fmt.Sprintf(`
resource "couchbase-capella_flush" %[1]q {
	organization_id = %[2]q
	project_id      = %[3]q
	cluster_id      = couchbase-capella_cluster.%[4]s.id
	bucket_id       = couchbase-capella_bucket.%[5]s.id

	depends_on = [couchbase-capella_cluster_deletion_protection.%[6]s]
}
`, flushResourceName, globalOrgId, globalProjectId, clusterResourceName, bucketResourceName, deletionProtectionResourceName)
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	name            = "%[4]s"
	cloud_provider = {
		type   = "aws"
		region = "us-east-1"
		cidr   = "%[8]s"
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
		plan     = "developer pro"
		timezone = "PT"
	}
}

resource "couchbase-capella_bucket" "%[5]s" {
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	cluster_id      = couchbase-capella_cluster.%[4]s.id
	name            = "%[5]s"
	flush           = true

	depends_on = [couchbase-capella_cluster.%[4]s]
}

resource "couchbase-capella_cluster_deletion_protection" "%[7]s" {
	organization_id     = "%[2]s"
	project_id          = "%[3]s"
	cluster_id          = couchbase-capella_cluster.%[4]s.id
	deletion_protection = %[6]t

	depends_on = [couchbase-capella_bucket.%[5]s]
}

%[9]s
`, globalProviderBlock, globalOrgId, globalProjectId, clusterResourceName, bucketResourceName, enabled, deletionProtectionResourceName, cidr, flushConfig)
}

// --- Import ID ---

func generateDeletionProtectionImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		return fmt.Sprintf(
			"cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func registerDeletionProtectionCleanupForClusterResource(t *testing.T, clusterReference string) resource.TestCheckFunc {
	t.Helper()
	registered := false

	return func(state *terraform.State) error {
		if registered {
			return nil
		}

		clusterResource, ok := state.RootModule().Resources[clusterReference]
		if !ok || clusterResource.Primary == nil {
			return fmt.Errorf("resource %s not found in state", clusterReference)
		}

		clusterID := clusterResource.Primary.ID
		if clusterID == "" {
			clusterID = clusterResource.Primary.Attributes["id"]
		}
		if clusterID == "" {
			return fmt.Errorf("resource %s has no cluster id in state", clusterReference)
		}

		registered = true
		t.Cleanup(func() {
			disableDeletionProtection(t, clusterID)
		})

		return nil
	}
}

// disableDeletionProtectionOnCleanup registers a t.Cleanup hook that
// unconditionally disables deletion protection on globalClusterId so that
// TestMain teardown can delete the cluster even if the test fails mid-run.
func disableDeletionProtectionOnCleanup(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		disableDeletionProtection(t, globalClusterId)
	})
}

func disableDeletionProtection(t *testing.T, clusterID string) {
	t.Helper()
	ctx := context.Background()
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/deletionProtection",
		globalHost, globalOrgId, globalProjectId, clusterID)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := globalClient.ExecuteWithRetry(
		ctx,
		cfg,
		clusterapi.UpdateDeletionProtectionRequest{DeletionProtection: false},
		globalToken,
		nil,
	)
	if err != nil {
		t.Logf("WARNING: failed to disable deletion protection on cleanup for cluster %s: %v", clusterID, err)
	}
}
