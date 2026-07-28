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

// TestAccNetworkPeerProviderTypeConfigMismatch_AV_131927 verifies that the provider rejects a
// network_peer configuration at plan time when provider_type does not match the provider_config
// block. For example, provider_type = "azure" with provider_config.aws_config should fail
// validation before any API call.
func TestAccNetworkPeerProviderTypeConfigMismatch_AV_131927(t *testing.T) {
	tests := []struct {
		name         string
		providerType string
		configBlock  string
		errorPattern string
	}{
		{
			name:         "azure provider_type with aws_config",
			providerType: "azure",
			configBlock:  "aws_config",
			errorPattern: `provider_type does not match the provider config block`,
		},
		{
			name:         "aws provider_type with azure_config",
			providerType: "aws",
			configBlock:  "azure_config",
			errorPattern: `provider_type does not match the provider config block`,
		},
		{
			name:         "gcp provider_type with aws_config",
			providerType: "gcp",
			configBlock:  "aws_config",
			errorPattern: `provider_type does not match the provider config block`,
		},
		{
			name:         "aws provider_type with gcp_config",
			providerType: "aws",
			configBlock:  "gcp_config",
			errorPattern: `provider_type does not match the provider config block`,
		},
		{
			name:         "gcp provider_type with azure_config",
			providerType: "gcp",
			configBlock:  "azure_config",
			errorPattern: `provider_type does not match the provider config block`,
		},
		{
			name:         "azure provider_type with gcp_config",
			providerType: "azure",
			configBlock:  "gcp_config",
			errorPattern: `provider_type does not match the provider config block`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_np_mismatch_")
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config:      testAccNetworkPeerProviderTypeConfigMismatchConfig(resourceName, tt.providerType, tt.configBlock),
						ExpectError: regexp.MustCompile(tt.errorPattern),
					},
				},
			})
		})
	}
}

func testAccNetworkPeerProviderTypeConfigMismatchConfig(resourceName, providerType, configBlock string) string {
	var config string
	switch configBlock {
	case "aws_config":
		config = `
    aws_config = {
      account_id = "123456789012"
      vpc_id     = "vpc-1234567890abcdef0"
      cidr       = "10.10.0.0/16"
      region     = "us-east-1"
    }`
	case "gcp_config":
		config = `
    gcp_config = {
      cidr            = "10.10.0.0/16"
      network_name    = "test-network"
      project_id      = "test-project"
      service_account = "test@test.iam.gserviceaccount.com"
    }`
	case "azure_config":
		config = `
    azure_config = {
      cidr            = "10.10.0.0/16"
      tenant_id       = "fee88efb-27a4-4ef6-937e-886a970af84b"
      resource_group  = "test-resource-group"
      subscription_id = "7df08e2f-efb1-4ed0-be9c-bb9a9d99ec84"
      vnet_id         = "test-vnet"
    }`
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_network_peer" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  name            = "qe-mismatch-provider-config"
  provider_type   = "%[3]s"
  provider_config = {%[4]s
  }
}
`, globalProviderBlock, resourceName, providerType, config)
}
