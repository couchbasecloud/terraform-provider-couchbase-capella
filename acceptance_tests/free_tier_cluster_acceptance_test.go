package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Live free-tier cluster CRUD/import coverage lives in
// TestAccFreeTierClusterLifecycle, which shares a single cluster across all
// free-tier resources (only one free-tier cluster is allowed per organization).
// The tests below are validation-only: they fail at plan time and never create
// a cluster, so they stay independent and run in parallel.

func TestAccFreeTierClusterResourceInvalidCloudProviderType(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_invalid_cloud_provider_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccFreeTierClusterResourceInvalidCloudProviderTypeConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)cloud_provider.*type.*value must be one of|Attribute.*type.*one of`),
			},
		},
	})
}

func TestAccFreeTierClusterResourceEmptyName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_empty_name_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccFreeTierClusterResourceEmptyNameConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)name.*at least 1|Attribute.*name.*length`),
			},
		},
	})
}

func TestAccFreeTierClusterResourceNameTooLong(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_name_too_long_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccFreeTierClusterResourceNameTooLongConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)name.*256|Attribute.*name.*length`),
			},
		},
	})
}

func testAccFreeTierClusterResourceInvalidCloudProviderTypeConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[2]s" {
	organization_id = "00000000-0000-0000-0000-000000000000"
	project_id      = "11111111-1111-1111-1111-111111111111"
	name            = "%[2]s"
	description     = "Invalid free tier cloud provider."

	cloud_provider = {
		type   = "oracle"
		region = "us-east-1"
		cidr   = "10.30.0.0/20"
	}
}
`, globalProviderBlock, resourceName)
}

func testAccFreeTierClusterResourceNameTooLongConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
	too_long_name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}

resource "couchbase-capella_free_tier_cluster" "%[2]s" {
	organization_id = "00000000-0000-0000-0000-000000000000"
	project_id      = "11111111-1111-1111-1111-111111111111"
	name            = local.too_long_name
	description     = "Invalid free tier long name."

	cloud_provider = {
		type   = "aws"
		region = "us-east-2"
		cidr   = "10.32.0.0/20"
	}
}
`, globalProviderBlock, resourceName)
}

func testAccFreeTierClusterResourceEmptyNameConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[2]s" {
	organization_id = "00000000-0000-0000-0000-000000000000"
	project_id      = "11111111-1111-1111-1111-111111111111"
	name            = ""
	description     = "Invalid free tier empty name."

	cloud_provider = {
		type   = "aws"
		region = "us-east-1"
		cidr   = "10.31.0.0/20"
	}
}
`, globalProviderBlock, resourceName)
}
