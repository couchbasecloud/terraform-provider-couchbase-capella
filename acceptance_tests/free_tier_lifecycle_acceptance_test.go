package acceptance_tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccFreeTierClusterLifecycle exercises the free-tier resources and data
// sources that need a live cluster against a SINGLE shared free-tier cluster:
// the cluster resource, on/off, bucket, and buckets data source, plus the app
// service. The free_tier_clusters data source is intentionally not covered here
// because the Capella API has no free-tier cluster list endpoint and the
// general /clusters list does not include free-tier clusters, so it cannot be
// exercised (see the datasource notes).
//
// Capella allows only one free-tier operational cluster per organization
// (docs/index.md), so these cases cannot run in parallel and previously each
// stood up its own cluster - eight cluster creations across the suite. Folding
// them into one ordered flow reuses a single cluster and keeps the org's lone
// free-tier slot occupied only for the duration of this test.
//
// Step ordering is deliberate:
//   - The on/off toggle runs while only the cluster exists, so no bucket or app
//     service is refreshed against a turned-off cluster.
//   - It ends on "on". The on/off resource's Delete is a no-op (the v4 API has
//     no destroy for it), so the cluster keeps whatever state the last apply
//     set; ending "on" guarantees teardown destroys the bucket, app service and
//     cluster against a healthy cluster.
//   - The app service is added last, once the cluster is confirmed healthy.
func TestAccFreeTierClusterLifecycle(t *testing.T) {
	clusterName := randomStringWithPrefix("tf_acc_free_tier_lifecycle_cluster_")
	bucketName := randomStringWithPrefix("tf_acc_free_tier_lifecycle_bucket_")
	appServiceName := randomStringWithPrefix("tf_acc_free_tier_lifecycle_appsvc_")
	onOffName := randomStringWithPrefix("tf_acc_free_tier_lifecycle_on_off_")
	bucketsDsName := randomStringWithPrefix("tf_acc_free_tier_lifecycle_buckets_ds_")
	cidr := generateRandomCIDR()

	clusterRef := "couchbase-capella_free_tier_cluster." + clusterName
	bucketRef := "couchbase-capella_free_tier_bucket." + bucketName
	appServiceRef := "couchbase-capella_free_tier_app_service." + appServiceName
	onOffRef := "couchbase-capella_free_tier_cluster_on_off." + onOffName
	bucketsDsRef := "data.couchbase-capella_free_tier_buckets." + bucketsDsName

	const (
		createDescription     = "Free tier lifecycle cluster."
		updateDescription     = "Updated free tier lifecycle cluster."
		appServiceDescription = "Updated free tier app service lifecycle."
	)

	// The cluster block is byte-identical from the update step onward (only the
	// description changes, and it is updated in place) so the shared cluster is
	// never replaced across steps.
	clusterCreate := testAccFreeTierLifecycleClusterBlock(clusterName, cidr, createDescription)
	cluster := testAccFreeTierLifecycleClusterBlock(clusterName, cidr, updateDescription)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// 1. Create the shared free-tier cluster.
			{
				Config: testAccFreeTierLifecycleConfig(clusterCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(clusterRef, "id"),
					resource.TestCheckResourceAttr(clusterRef, "name", clusterName),
					resource.TestCheckResourceAttr(clusterRef, "description", createDescription),
					resource.TestCheckResourceAttr(clusterRef, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(clusterRef, "cloud_provider.region", "us-east-2"),
					resource.TestCheckResourceAttr(clusterRef, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttrSet(clusterRef, "etag"),
				),
			},
			// 2. Import the cluster.
			{
				ResourceName:      clusterRef,
				ImportStateIdFunc: generateClusterImportIdForResource(clusterRef),
				ImportState:       true,
			},
			// 3. Update the cluster description in place.
			{
				Config: testAccFreeTierLifecycleConfig(cluster),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterRef, "name", clusterName),
					resource.TestCheckResourceAttr(clusterRef, "description", updateDescription),
				),
			},
			// 4. Turn the cluster off (only the cluster exists at this point).
			{
				Config: testAccFreeTierLifecycleConfig(cluster, testAccFreeTierLifecycleOnOffBlock(clusterName, onOffName, "off")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onOffRef, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(onOffRef, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(onOffRef, "cluster_id"),
					resource.TestCheckResourceAttr(onOffRef, "state", "off"),
				),
			},
			// 5. Import the on/off resource.
			{
				ResourceName:                         onOffRef,
				ImportStateIdFunc:                    generateFreeTierClusterOnOffImportId(onOffRef),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
			// 6. Turn the cluster back on so the rest of the flow runs healthy.
			{
				Config: testAccFreeTierLifecycleConfig(cluster, testAccFreeTierLifecycleOnOffBlock(clusterName, onOffName, "on")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onOffRef, "cluster_id"),
					resource.TestCheckResourceAttr(onOffRef, "state", "on"),
				),
			},
			// 7. Create a bucket on the shared cluster.
			{
				Config: testAccFreeTierLifecycleConfig(cluster, testAccFreeTierLifecycleBucketBlock(clusterName, bucketName, 100)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(bucketRef, "id"),
					resource.TestCheckResourceAttr(bucketRef, "name", bucketName),
					resource.TestCheckResourceAttr(bucketRef, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(bucketRef, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(bucketRef, "cluster_id"),
					resource.TestCheckResourceAttr(bucketRef, "memory_allocation_in_mb", "100"),
					resource.TestCheckResourceAttrSet(bucketRef, "type"),
					resource.TestCheckResourceAttrSet(bucketRef, "storage_backend"),
					resource.TestCheckResourceAttrSet(bucketRef, "vbuckets"),
				),
			},
			// 8. Import the bucket.
			{
				ResourceName:      bucketRef,
				ImportStateIdFunc: generateBucketImportIdForResource(bucketRef),
				ImportState:       true,
			},
			// 9. Update the bucket memory allocation.
			{
				Config: testAccFreeTierLifecycleConfig(cluster, testAccFreeTierLifecycleBucketBlock(clusterName, bucketName, 200)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bucketRef, "name", bucketName),
					resource.TestCheckResourceAttr(bucketRef, "memory_allocation_in_mb", "200"),
				),
			},
			// 10. Read the free-tier buckets data source.
			{
				Config: testAccFreeTierLifecycleConfig(
					cluster,
					testAccFreeTierLifecycleBucketBlock(clusterName, bucketName, 200),
					testAccFreeTierLifecycleBucketsDsBlock(clusterName, bucketName, bucketsDsName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bucketRef, "name", bucketName),
					resource.TestCheckResourceAttrSet(bucketRef, "id"),
					resource.TestCheckResourceAttr(bucketsDsRef, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(bucketsDsRef, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(bucketsDsRef, "cluster_id"),
					resource.TestCheckResourceAttrSet(bucketsDsRef, "data.#"),
					testAccCheckListElemNestedAttrs(bucketsDsRef, "data", map[string]string{
						"name":            bucketName,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
					}),
				),
			},
			// 11. Create an app service (required fields only) on the healthy cluster.
			{
				Config: testAccFreeTierLifecycleConfig(
					cluster,
					testAccFreeTierLifecycleBucketBlock(clusterName, bucketName, 200),
					testAccFreeTierLifecycleAppServiceBlock(clusterName, appServiceName, ""),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(appServiceRef, "id"),
					resource.TestCheckResourceAttr(appServiceRef, "name", appServiceName),
					resource.TestCheckResourceAttr(appServiceRef, "description", ""),
					resource.TestCheckResourceAttr(appServiceRef, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(appServiceRef, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(appServiceRef, "cluster_id"),
					resource.TestCheckResourceAttrSet(appServiceRef, "current_state"),
					resource.TestCheckResourceAttrSet(appServiceRef, "etag"),
				),
			},
			// 12. Import the app service.
			{
				ResourceName:      appServiceRef,
				ImportStateIdFunc: generateAppServiceImportId(appServiceRef),
				ImportState:       true,
			},
			// 13. Update the app service description.
			{
				Config: testAccFreeTierLifecycleConfig(
					cluster,
					testAccFreeTierLifecycleBucketBlock(clusterName, bucketName, 200),
					testAccFreeTierLifecycleAppServiceBlock(clusterName, appServiceName, appServiceDescription),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(appServiceRef, "name", appServiceName),
					resource.TestCheckResourceAttr(appServiceRef, "description", appServiceDescription),
				),
			},
		},
	})
}

// testAccFreeTierLifecycleConfig joins the provider block with the given
// resource/data-source blocks into a single Terraform config.
func testAccFreeTierLifecycleConfig(blocks ...string) string {
	return globalProviderBlock + "\n" + strings.Join(blocks, "\n")
}

func testAccFreeTierLifecycleClusterBlock(clusterName, cidr, description string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_free_tier_cluster" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[1]s"
  description     = "%[4]s"

  cloud_provider = {
    type   = "aws"
    region = "us-east-2"
    cidr   = "%[5]s"
  }
}
`, clusterName, globalOrgId, globalProjectId, description, cidr)
}

func testAccFreeTierLifecycleBucketBlock(clusterName, bucketName string, memoryAllocationInMB int64) string {
	return fmt.Sprintf(`
resource "couchbase-capella_free_tier_bucket" "%[1]s" {
  organization_id         = "%[2]s"
  project_id              = "%[3]s"
  cluster_id              = couchbase-capella_free_tier_cluster.%[4]s.id
  name                    = "%[1]s"
  memory_allocation_in_mb = %[5]d
}
`, bucketName, globalOrgId, globalProjectId, clusterName, memoryAllocationInMB)
}

// testAccFreeTierLifecycleAppServiceBlock omits the description attribute when
// description is empty so the required-fields-only case can be exercised.
func testAccFreeTierLifecycleAppServiceBlock(clusterName, appServiceName, description string) string {
	descriptionLine := ""
	if description != "" {
		descriptionLine = fmt.Sprintf("\n  description     = %q", description)
	}
	return fmt.Sprintf(`
resource "couchbase-capella_free_tier_app_service" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id
  name            = "%[1]s"%[5]s
}
`, appServiceName, globalOrgId, globalProjectId, clusterName, descriptionLine)
}

func testAccFreeTierLifecycleOnOffBlock(clusterName, onOffName, state string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_free_tier_cluster_on_off" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id
  state           = "%[5]s"
}
`, onOffName, globalOrgId, globalProjectId, clusterName, state)
}

func testAccFreeTierLifecycleBucketsDsBlock(clusterName, bucketName, dsName string) string {
	return fmt.Sprintf(`
data "couchbase-capella_free_tier_buckets" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id

  depends_on = [couchbase-capella_free_tier_bucket.%[5]s]
}
`, dsName, globalOrgId, globalProjectId, clusterName, bucketName)
}

// generateFreeTierClusterOnOffImportId builds the composite import ID for the
// free-tier cluster on/off resource.
func generateFreeTierClusterOnOffImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, module := range state.Modules {
			if len(module.Resources) > 0 {
				if v, ok := module.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		return fmt.Sprintf(
			"organization_id=%s,project_id=%s,cluster_id=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
		), nil
	}
}

// testAccCheckListElemNestedAttrs asserts that at least one element of a list
// attribute has all of the expected nested attribute values.
func testAccCheckListElemNestedAttrs(resourceName, listAttribute string, expected map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, module := range state.Modules {
			resourceState, ok := module.Resources[resourceName]
			if !ok {
				continue
			}

			attrs := resourceState.Primary.Attributes
			for index := 0; attrs[fmt.Sprintf("%s.%d.name", listAttribute, index)] != ""; index++ {
				matched := true
				for key, expectedValue := range expected {
					if attrs[fmt.Sprintf("%s.%d.%s", listAttribute, index, key)] != expectedValue {
						matched = false
						break
					}
				}
				if matched {
					return nil
				}
			}
		}

		return fmt.Errorf("no %s element matched expected attributes %v for %s", listAttribute, expected, resourceName)
	}
}
