package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAccBucketResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bucket_")
	resourceReference := "couchbase-capella_bucket." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccBucketResourceConfig(resourceName, "tf-acc-test-bucket", 100, "none", 1, false, 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBucketResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", "tf-acc-test-bucket"),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "100"),
					resource.TestCheckResourceAttr(resourceReference, "durability_level", "none"),
					resource.TestCheckResourceAttr(resourceReference, "replicas", "1"),
					resource.TestCheckResourceAttr(resourceReference, "flush", "false"),
					resource.TestCheckResourceAttr(resourceReference, "time_to_live_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceReference, "type", "couchbase"),
					resource.TestCheckResourceAttr(resourceReference, "storage_backend", "couchstore"),
					resource.TestCheckResourceAttr(resourceReference, "bucket_conflict_resolution", "seqno"),
					resource.TestCheckResourceAttr(resourceReference, "eviction_policy", "fullEviction"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.item_count"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.ops_per_second"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.disk_used_in_mib"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.memory_used_in_mib"),
				),
			},
			// ImportState
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateBucketImportIdForResource(resourceReference),
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore:              []string{"stats"},
			},
			// Update updatable fields: memory_allocation_in_mb, durability_level, replicas, flush, time_to_live_in_seconds
			{
				Config: testAccBucketResourceConfig(resourceName, "tf-acc-test-bucket", 200, "majority", 2, true, 3600),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBucketResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", "tf-acc-test-bucket"),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "200"),
					resource.TestCheckResourceAttr(resourceReference, "durability_level", "majority"),
					resource.TestCheckResourceAttr(resourceReference, "replicas", "2"),
					resource.TestCheckResourceAttr(resourceReference, "flush", "true"),
					resource.TestCheckResourceAttr(resourceReference, "time_to_live_in_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceReference, "type", "couchbase"),
					resource.TestCheckResourceAttr(resourceReference, "storage_backend", "couchstore"),
					resource.TestCheckResourceAttr(resourceReference, "bucket_conflict_resolution", "seqno"),
					resource.TestCheckResourceAttr(resourceReference, "eviction_policy", "fullEviction"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.item_count"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.ops_per_second"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.disk_used_in_mib"),
					resource.TestCheckResourceAttrSet(resourceReference, "stats.memory_used_in_mib"),
				),
			},
		},
	})
}

func TestAccBucketResourceDefaults(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bucket_")
	resourceReference := "couchbase-capella_bucket." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketResourceMinimalConfig(resourceName, "tf-acc-defaults-bucket"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBucketResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "tf-acc-defaults-bucket"),
					resource.TestCheckResourceAttr(resourceReference, "type", "couchbase"),
					resource.TestCheckResourceAttr(resourceReference, "storage_backend", "couchstore"),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "100"),
					resource.TestCheckResourceAttr(resourceReference, "bucket_conflict_resolution", "seqno"),
					resource.TestCheckResourceAttr(resourceReference, "durability_level", "none"),
					resource.TestCheckResourceAttr(resourceReference, "replicas", "1"),
					resource.TestCheckResourceAttr(resourceReference, "flush", "false"),
					resource.TestCheckResourceAttr(resourceReference, "time_to_live_in_seconds", "0"),
					resource.TestCheckResourceAttr(resourceReference, "eviction_policy", "fullEviction"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

func TestAccBucketResourceInvalidName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketResourceMinimalConfig(resourceName, " leading-space"),
				ExpectError: regexp.MustCompile("name attribute has leading or trailing spaces"),
			},
		},
	})
}

func testAccBucketResourceConfig(resourceName, bucketName string, memoryMB int, durability string, replicas int, flush bool, ttl int) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_bucket" "%[2]s" {
		organization_id        = "%[3]s"
		project_id             = "%[4]s"
		cluster_id             = "%[5]s"
		name                   = "%[6]s"
		memory_allocation_in_mb = %[7]d
		durability_level       = "%[8]s"
		replicas               = %[9]d
		flush                  = %[10]t
		time_to_live_in_seconds = %[11]d
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, bucketName, memoryMB, durability, replicas, flush, ttl)
}

func testAccBucketResourceMinimalConfig(resourceName, bucketName string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_bucket" "%[2]s" {
		organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
		name            = "%[6]s"
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, bucketName)
}

func generateBucketImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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

func retrieveBucketFromServer(data *providerschema.Data, organizationId, projectId, clusterId, bucketId string) (*bucketapi.GetBucketResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", data.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	bucketResp := bucketapi.GetBucketResponse{}
	if err = json.Unmarshal(response.Body, &bucketResp); err != nil {
		return nil, err
	}

	return &bucketResp, nil
}

func testAccExistsBucketResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}

		data := newTestClient(t)
		bucketResp, err := retrieveBucketFromServer(
			data,
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
			rawState["id"],
		)
		if err != nil {
			return fmt.Errorf("failed to retrieve bucket from server: %w", err)
		}

		if bucketResp.Name != rawState["name"] {
			return fmt.Errorf("bucket name mismatch: API returned %q, state has %q", bucketResp.Name, rawState["name"])
		}
		if bucketResp.Id != rawState["id"] {
			return fmt.Errorf("bucket id mismatch: API returned %q, state has %q", bucketResp.Id, rawState["id"])
		}

		return nil
	}
}
