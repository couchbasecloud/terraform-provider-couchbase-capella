package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccNetworkPeerAzureInvalidTenantId deploys an Azure cluster and then validates
// that creating a network peer with an invalid Azure tenant ID returns an appropriate
// error from the Capella API.
func TestAccNetworkPeerAzureInvalidTenantId(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_azure_cluster_")
	networkPeerResourceName := "tf_acc_network_peer_invalid_tenant"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	cidr := "10.0.6.0/23"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Deploy the Azure cluster
			{
				Config: testAccAzureClusterConfig(clusterResourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceReference, "cloud_provider.type", "azure"),
				),
			},
			// Step 2: Attempt to create a network peer with an invalid Azure tenant ID
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, cidr, networkPeerResourceName, "ffffffff-aaaa-1414-eeee-000000000000", "ffffffff-aaaa-1414-eeee-000000000000", "10.1.0.0/16", "test-rg", "test-vnet"),
				ExpectError: regexp.MustCompile("The provided Azure tenant ID is invalid"),
			},
		},
	})
}

// TestAccNetworkPeerAzureInvalidSubscriptionId deploys an Azure cluster and then validates
// that creating a network peer with an invalid Azure subscription ID returns an appropriate
// error from the Capella API.
func TestAccNetworkPeerAzureInvalidSubscriptionId(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_azure_cluster_")
	networkPeerResourceName := "tf_acc_network_peer_invalid_sub"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	cidr := "10.0.8.0/23"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Deploy the Azure cluster
			{
				Config: testAccAzureClusterConfig(clusterResourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceReference, "cloud_provider.type", "azure"),
				),
			},
			// Step 2: Attempt to create a network peer with an invalid Azure subscription ID
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, cidr, networkPeerResourceName, "ffffffff-aaaa-1414-eeee-000000000000", "00000000-0000-0000-0000-000000000000", "10.1.0.0/16", "test-rg", "test-vnet"),
				ExpectError: regexp.MustCompile("There is an error during network peer creation"),
			},
		},
	})
}

// TestAccNetworkPeerAzureInvalidResourceGroup deploys an Azure cluster and then validates
// that creating a network peer with an invalid Azure resource group returns an appropriate
// error from the Capella API.
func TestAccNetworkPeerAzureInvalidResourceGroup(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_azure_cluster_")
	networkPeerResourceName := "tf_acc_network_peer_invalid_rg"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	cidr := "10.0.10.0/23"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Deploy the Azure cluster
			{
				Config: testAccAzureClusterConfig(clusterResourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceReference, "cloud_provider.type", "azure"),
				),
			},
			// Step 2: Attempt to create a network peer with an invalid resource group
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, cidr, networkPeerResourceName, "ffffffff-aaaa-1414-eeee-000000000000", "ffffffff-aaaa-1414-eeee-000000000000", "10.1.0.0/16", "nonexistent-rg-000", "test-vnet"),
				ExpectError: regexp.MustCompile("There is an error during network peer creation"),
			},
		},
	})
}

// TestAccNetworkPeerAzureInvalidVnetId deploys an Azure cluster and then validates
// that creating a network peer with an invalid Azure VNet ID returns an appropriate
// error from the Capella API.
func TestAccNetworkPeerAzureInvalidVnetId(t *testing.T) {
	clusterResourceName := randomStringWithPrefix("tf_acc_azure_cluster_")
	networkPeerResourceName := "tf_acc_network_peer_invalid_vnet"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	cidr := "10.0.12.0/23"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Deploy the Azure cluster
			{
				Config: testAccAzureClusterConfig(clusterResourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceReference, "cloud_provider.type", "azure"),
				),
			},
			// Step 2: Attempt to create a network peer with an invalid VNet ID
			{
				Config:      testAccAzureClusterWithNetworkPeerConfig(clusterResourceName, cidr, networkPeerResourceName, "ffffffff-aaaa-1414-eeee-000000000000", "ffffffff-aaaa-1414-eeee-000000000000", "10.1.0.0/16", "test-rg", "nonexistent-vnet-000"),
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
