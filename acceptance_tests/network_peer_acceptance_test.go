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
