package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeTierClusterResourceCloudProvidersSequential(t *testing.T) {
	awsResourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_aws_")
	azureResourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_azure_")
	gcpResourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_gcp_")
	awsCidr := generateRandomCIDR()
	azureCidr := generateRandomCIDR()
	gcpCidr := generateRandomCIDR()

	awsConfig := testAccFreeTierClusterResourceCloudProviderConfig(awsResourceName, "aws", "us-east-2", awsCidr, "Valid free tier AWS cluster.")
	azureConfig := testAccFreeTierClusterResourceCloudProviderConfig(azureResourceName, "azure", "eastus", azureCidr, "Valid free tier Azure cluster.")
	gcpConfig := testAccFreeTierClusterResourceCloudProviderConfig(gcpResourceName, "gcp", "us-central1", gcpCidr, "Valid free tier GCP cluster.")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: awsConfig,
				Check: testAccFreeTierClusterResourceCloudProviderChecks(
					awsResourceName,
					"aws",
					"us-east-2",
					awsCidr,
					"Valid free tier AWS cluster.",
				),
			},
			{
				Config:  awsConfig,
				Destroy: true,
			},
			{
				Config: azureConfig,
				Check: testAccFreeTierClusterResourceCloudProviderChecks(
					azureResourceName,
					"azure",
					"eastus",
					azureCidr,
					"Valid free tier Azure cluster.",
				),
			},
			{
				Config:  azureConfig,
				Destroy: true,
			},
			{
				Config: gcpConfig,
				Check: testAccFreeTierClusterResourceCloudProviderChecks(
					gcpResourceName,
					"gcp",
					"us-central1",
					gcpCidr,
					"Valid free tier GCP cluster.",
				),
			},
			{
				Config:  gcpConfig,
				Destroy: true,
			},
		},
	})
}

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

func testAccFreeTierClusterResourceCloudProviderChecks(resourceName, cloudProviderType, region, cidr, description string) resource.TestCheckFunc {
	resourceReference := "couchbase-capella_free_tier_cluster." + resourceName

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceReference, "id"),
		resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
		resource.TestCheckResourceAttr(resourceReference, "description", description),
		resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", cloudProviderType),
		resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", region),
		resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
		resource.TestCheckResourceAttrSet(resourceReference, "etag"),
	)
}

func testAccFreeTierClusterResourceCloudProviderConfig(resourceName, cloudProviderType, region, cidr, description string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[4]s" {
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	name            = "%[4]s"
	description     = "%[8]s"

	cloud_provider = {
		type   = "%[5]s"
		region = "%[6]s"
		cidr   = "%[7]s"
	}
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cloudProviderType, region, cidr, description)
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
