package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccNetworkPeerInvalidProviderType verifies that the network peer resource rejects an
// unsupported provider_type at plan time. The OneOf("aws", "gcp", "azure") validator fires
// before any API call, so dummy org/project/cluster IDs are sufficient.
func TestAccNetworkPeerInvalidProviderType(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_network_peer_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkPeerInvalidProviderTypeConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)provider_type.*value must be one of.*aws.*gcp.*azure`),
			},
		},
	})
}

// TestAccNetworkPeerMismatchedProviderConfig verifies that the network peer resource rejects
// a provider_config block that does not match the declared provider_type at plan time (AV-131927).
// Validation fires before any API call, so dummy org/project/cluster IDs are sufficient.
func TestAccNetworkPeerMismatchedProviderConfig(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_network_peer_mismatch_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkPeerMismatchedProviderConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)provider_type "azure" requires.*provider_config\.azure_config.*provider_config\.aws_config.*configured`),
			},
		},
	})
}

// TestAccNetworkPeerMultipleProviderConfigs verifies that the network peer resource rejects a
// provider_config with more than one cloud config block set at plan time (AV-131927).
func TestAccNetworkPeerMultipleProviderConfigs(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_network_peer_multi_cfg_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkPeerMultipleProviderConfigs(resourceName),
				ExpectError: regexp.MustCompile(`(?s)only one of aws_config.*gcp_config.*azure_config.*may be set`),
			},
		},
	})
}

func testAccNetworkPeerMismatchedProviderConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_network_peer" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  name            = "qe-mismatched-provider-config"
  provider_type   = "azure"
  provider_config = {
    aws_config = {
      account_id = "123456789012"
      vpc_id     = "vpc-1234567890abcdef0"
      cidr       = "10.99.0.0/16"
      region     = "us-east-1"
    }
  }
}
`, globalProviderBlock, resourceName)
}

func testAccNetworkPeerMultipleProviderConfigs(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_network_peer" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  name            = "qe-multiple-provider-configs"
  provider_type   = "aws"
  provider_config = {
    aws_config = {
      account_id = "123456789012"
      vpc_id     = "vpc-1234567890abcdef0"
      cidr       = "10.99.0.0/16"
      region     = "us-east-1"
    }
    azure_config = {
      tenant_id       = "fee88efb-27a4-4ef6-937e-886a970af84b"
      cidr            = "10.0.0.0/16"
      resource_group  = "test-rg"
      subscription_id = "7df08e2f-efb1-4ed0-be9c-bb9a9d99ec84"
      vnet_id         = "test-vnet"
    }
  }
}
`, globalProviderBlock, resourceName)
}

func testAccNetworkPeerInvalidProviderTypeConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_network_peer" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  name            = "qe-invalid-provider-type"
  provider_type   = "ibm"
  provider_config = {
    aws_config = {
      account_id = "123456789012"
      vpc_id     = "vpc-1234567890abcdef0"
      cidr       = "10.10.0.0/16"
      region     = "us-east-1"
    }
  }
}
`, globalProviderBlock, resourceName)
}
