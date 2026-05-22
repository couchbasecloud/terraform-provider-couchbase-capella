package acceptance_tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkPeerAzureInvalidResourceGroup(t *testing.T) {
	preCheckAzureNetworkPeer(t)

	clusterResourceName := randomStringWithPrefix("tf_acc_network_peer_azure_cluster_")
	clusterCIDR := generateRandomCIDR()
	resourceName := randomStringWithPrefix("tf_acc_network_peer_failed_")
	peerName := randomStringWithPrefix("tf-acc-network-peer-failed-")
	resourceReference := "couchbase-capella_network_peer." + resourceName
	mismatchResourceName := randomStringWithPrefix("tf_acc_network_peer_mismatch_")
	mismatchPeerName := randomStringWithPrefix("tf-acc-network-peer-mismatch-")

	steps := []resource.TestStep{
		{
			Config: testAccNetworkPeerAzureFailedConfig(clusterResourceName, resourceName, peerName, clusterCIDR, azureNetworkPeerCIDR(), azureNetworkPeerVNetID()),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
				resource.TestCheckResourceAttr(resourceReference, "name", peerName),
				resource.TestCheckResourceAttr(resourceReference, "provider_type", "azure"),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.tenant_id", os.Getenv("TF_VAR_azure_tenant_id")),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.subscription_id", os.Getenv("TF_VAR_azure_subscription_id")),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.resource_group", azureNetworkPeerInvalidResourceGroup()),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.vnet_id", azureNetworkPeerVNetID()),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.cidr", azureNetworkPeerCIDR()),
				resource.TestCheckResourceAttr(resourceReference, "status.state", "failed"),
				resource.TestCheckResourceAttrSet(resourceReference, "id"),
			),
		},
		{
			RefreshState: true,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "provider_type", "azure"),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.resource_group", azureNetworkPeerInvalidResourceGroup()),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.vnet_id", azureNetworkPeerVNetID()),
				resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.cidr", azureNetworkPeerCIDR()),
				resource.TestCheckResourceAttr(resourceReference, "status.state", "failed"),
			),
		},
		{
			ResourceName:      resourceReference,
			ImportStateIdFunc: generateNetworkPeerImportIdForResource(resourceReference),
			ImportState:       true,
			// ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"commands",
			},
		},
		// {
		// 	RefreshState: true,
		// 	Check: resource.ComposeAggregateTestCheckFunc(
		// 		resource.TestCheckResourceAttr(resourceReference, "provider_type", "azure"),
		// 		resource.TestCheckResourceAttr(resourceReference, "provider_config.azure_config.resource_group", azureNetworkPeerInvalidResourceGroup()),
		// 	),
		// },
	}

	if os.Getenv("TF_ACC_NETWORK_PEER_AZURE_AMBIGUOUS_NAME") != "" {
		sameNameResourceName := randomStringWithPrefix("tf_acc_network_peer_same_name_")
		sameNamePeerName := randomStringWithPrefix("tf-acc-network-peer-same-name-")
		sameNameResourceReference := "couchbase-capella_network_peer." + sameNameResourceName

		steps = append(steps,
			resource.TestStep{
				Config: testAccNetworkPeerAzureFailedConfig(clusterResourceName, sameNameResourceName, sameNamePeerName, clusterCIDR, "10.10.0.0/16", "tf-acc-old-vnet"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(sameNameResourceReference, "id"),
					resource.TestCheckResourceAttr(sameNameResourceReference, "provider_config.azure_config.cidr", "10.10.0.0/16"),
					removeNetworkPeerFromState(sameNameResourceReference),
				),
			},
			resource.TestStep{
				Config:      testAccNetworkPeerAzureFailedConfig(clusterResourceName, sameNameResourceName, sameNamePeerName, clusterCIDR, "10.11.0.0/16", azureNetworkPeerVNetID()),
				ExpectError: regexp.MustCompile(`(?s)ambiguous|multiple|same name`),
			},
		)
	}

	steps = append(steps, resource.TestStep{
		Config: testAccNetworkPeerAzureInvalidResourceGroupWithMismatchConfig(
			clusterResourceName,
			resourceName,
			peerName,
			mismatchResourceName,
			mismatchPeerName,
			clusterCIDR,
			azureNetworkPeerCIDR(),
			azureNetworkPeerVNetID(),
		),
		ExpectError: regexp.MustCompile(`(?s)provider_type.*azure.*azure_config|provider_type.*match|provider config`),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

func TestAccDatasourceNetworkPeersReadsExistingPeers(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_network_peers_ds_")
	dsReference := "data.couchbase-capella_network_peers." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkPeersDataSourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
				),
			},
		},
	})
}

func preCheckAzureNetworkPeer(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC_NETWORK_PEER_AZURE") == "" {
		t.Skip("Skipping Azure network peer acceptance test. Set TF_ACC_NETWORK_PEER_AZURE=1 after completing Azure consent in Capella UI.")
	}

	required := map[string]string{
		"TF_VAR_project_id":            globalProjectId,
		"TF_VAR_azure_tenant_id":       os.Getenv("TF_VAR_azure_tenant_id"),
		"TF_VAR_azure_subscription_id": os.Getenv("TF_VAR_azure_subscription_id"),
		"TF_VAR_azure_vnet_id":         os.Getenv("TF_VAR_azure_vnet_id"),
		"TF_VAR_azure_vnet_cidr":       os.Getenv("TF_VAR_azure_vnet_cidr"),
	}

	for name, value := range required {
		if value == "" {
			t.Skipf("Skipping Azure network peer acceptance test: %s is not set. Azure consent must be granted through Capella UI before running this test.", name)
		}
	}
}

func testAccNetworkPeerAzureFailedConfig(clusterResourceName, resourceName, peerName, clusterCIDR, cidr, vnetID string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s
`, globalProviderBlock, testAccNetworkPeerAzureClusterConfig(clusterResourceName, clusterCIDR), testAccNetworkPeerAzureInvalidResourceGroupConfig(clusterResourceName, resourceName, peerName, cidr, vnetID))
}

func testAccNetworkPeerAzureInvalidResourceGroupWithMismatchConfig(clusterResourceName, invalidResourceName, invalidPeerName, mismatchResourceName, mismatchPeerName, clusterCIDR, cidr, vnetID string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s
`, globalProviderBlock,
		testAccNetworkPeerAzureClusterConfig(clusterResourceName, clusterCIDR),
		testAccNetworkPeerAzureInvalidResourceGroupConfig(clusterResourceName, invalidResourceName, invalidPeerName, cidr, vnetID),
		testAccNetworkPeerProviderTypeMismatchConfig(clusterResourceName, mismatchResourceName, mismatchPeerName),
	)
}

func testAccNetworkPeerAzureClusterConfig(clusterResourceName, clusterCIDR string) string {
	return fmt.Sprintf(`resource "couchbase-capella_cluster" "%[1]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	name            = "%[1]s"
	description     = "Terraform Acceptance Test Azure network peer cluster"
	cloud_provider = {
		type   = "azure"
		region = "eastus"
		cidr   = "%[2]s"
	}
	service_groups = [
		{
			node = {
				compute = {
					cpu = 4
					ram = 16
				}
				disk = {
					type = "P6"
				}
			}
			num_of_nodes = 3
			services     = ["data"]
		}
	]
	availability = {
		"type" : "multi"
	}
	support = {
		plan     = "enterprise"
		timezone = "PT"
	}
}`, clusterResourceName, clusterCIDR, globalOrgId, globalProjectId)
}

func testAccNetworkPeerAzureInvalidResourceGroupConfig(clusterResourceName, resourceName, peerName, cidr, vnetID string) string {
	return fmt.Sprintf(`resource "couchbase-capella_network_peer" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = couchbase-capella_cluster.%[1]s.id

	name          = "%[5]s"
	provider_type = "azure"

	provider_config = {
		azure_config = {
			tenant_id       = "%[6]s"
			subscription_id = "%[7]s"
			resource_group  = "%[8]s"
			vnet_id         = "%[9]s"
			cidr            = "%[10]s"
		}
	}
}
`, clusterResourceName, resourceName, globalOrgId, globalProjectId, peerName,
		os.Getenv("TF_VAR_azure_tenant_id"), os.Getenv("TF_VAR_azure_subscription_id"), azureNetworkPeerInvalidResourceGroup(), vnetID, cidr)
}

func testAccNetworkPeerProviderTypeMismatchConfig(clusterResourceName, resourceName, peerName string) string {
	return fmt.Sprintf(`resource "couchbase-capella_network_peer" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = couchbase-capella_cluster.%[1]s.id

	name          = "%[5]s"
	provider_type = "azure"

	provider_config = {
		aws_config = {
			account_id = "123456789012"
			vpc_id     = "vpc-1234567890abcdef0"
			region     = "us-east-1"
			cidr       = "10.99.0.0/16"
		}
	}
}
`, clusterResourceName, resourceName, globalOrgId, globalProjectId, peerName)
}

func testAccNetworkPeersDataSourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_network_peers" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId)
}

func generateNetworkPeerImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, module := range state.Modules {
			if len(module.Resources) == 0 {
				continue
			}
			if value, ok := module.Resources[resourceReference]; ok {
				rawState := value.Primary.Attributes
				return fmt.Sprintf(
					"id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
					rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
				), nil
			}
		}
		return "", fmt.Errorf("resource %s not found in state", resourceReference)
	}
}

func removeNetworkPeerFromState(resourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, module := range state.Modules {
			if _, ok := module.Resources[resourceReference]; ok {
				delete(module.Resources, resourceReference)
				return nil
			}
		}
		return fmt.Errorf("resource %s not found in state", resourceReference)
	}
}

func azureNetworkPeerCIDR() string {
	return os.Getenv("TF_VAR_azure_vnet_cidr")
}

func azureNetworkPeerVNetID() string {
	return os.Getenv("TF_VAR_azure_vnet_id")
}

func azureNetworkPeerInvalidResourceGroup() string {
	return "tf-acc-intentionally-invalid-resource-group"
}
