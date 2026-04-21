package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccNetworkPeerAzureInvalidConfigs deploys a single Azure cluster and then validates
// that creating network peers with various invalid Azure configurations returns appropriate
// errors from the Capella API. All negative cases share one cluster to reduce runtime and cost.
func TestAccNetworkPeerAzureInvalidConfigs(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_azure_cluster_")
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	clusterCidr := generateRandomCIDR()

	// Use a peer CIDR in the 172.16.0.0/12 range so it can never overlap
	// with the cluster CIDR which is always in 10.0.0.0/8.
	peerCidr := "172.16.0.0/16"

	// Valid Azure peering parameters (used as base; each negative step overrides one field)
	validTenant := "fee88efb-27a4-4ef6-937e-886a970af84b"
	validSubscription := "7df08e2f-efb1-4ed0-be9c-bb9a9d99ec84"
	validRG := "testing-dataapi-rg"
	validVnet := "dataapi-private-endpoint-vn"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Deploy the Azure cluster
			{
				Config: testAccAzureClusterConfig(clusterResourceName, clusterCidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceReference, "cloud_provider.type", "azure"),
				),
			},
			// Step 2: Invalid tenant ID
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, clusterCidr, "tf_acc_np_invalid_tenant", "ffffffff-aaaa-1414-eeee-000000000000", validSubscription, peerCidr, validRG, validVnet),
				ExpectError: regexp.MustCompile("There is an error during network peer creation"),
			},
			// Step 3: Invalid subscription ID
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, clusterCidr, "tf_acc_np_invalid_sub", validTenant, "00000000-0000-0000-0000-000000000000", peerCidr, validRG, validVnet),
				ExpectError: regexp.MustCompile("There is an error during network peer creation"),
			},
			// Step 4: Invalid resource group
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, clusterCidr, "tf_acc_np_invalid_rg", validTenant, validSubscription, peerCidr, "nonexistent-rg-000", validVnet),
				ExpectError: regexp.MustCompile("There is an error during network peer creation"),
			},
			// Step 5: Invalid VNet ID
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, clusterCidr, "tf_acc_np_invalid_vnet", validTenant, validSubscription, peerCidr, validRG, "nonexistent-vnet-000"),
				ExpectError: regexp.MustCompile("There is an error during network peer creation"),
			},
		},
	})
}

// testAccAzureClusterConfig returns a Terraform config that deploys an Azure cluster.
func testAccAzureClusterConfig(clusterResourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Azure cluster for network peer acceptance test"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
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
          type          = "P6"
          autoexpansion = true
        }
      }
      num_of_nodes = 3
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
`, globalProviderBlock, globalOrgId, globalProjectId, clusterResourceName, cidr)
}

// testAccAzureClusterWithNetworkPeerConfig returns a Terraform config that includes both
// an Azure cluster and a network peer resource with the given Azure peering parameters.
func testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, clusterCidr, networkPeerResourceName, tenantId, subscriptionId, peerCidr, resourceGroup, vnetId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Azure cluster for network peer acceptance test"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
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
          type          = "P6"
          autoexpansion = true
        }
      }
      num_of_nodes = 3
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

resource "couchbase-capella_network_peer" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[4]s.id
  name            = "%[6]s"
  provider_type   = "azure"
  provider_config = {
    azure_config = {
      tenant_id       = "%[7]s"
      subscription_id = "%[8]s"
      cidr            = "%[9]s"
      resource_group  = "%[10]s"
      vnet_id         = "%[11]s"
    }
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterResourceName, clusterCidr,
		networkPeerResourceName, tenantId, subscriptionId, peerCidr, resourceGroup, vnetId)
}
