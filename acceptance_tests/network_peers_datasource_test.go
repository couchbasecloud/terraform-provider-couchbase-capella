package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceNetworkPeers(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_network_peers_ds_")
	dataSourceReference := "data.couchbase-capella_network_peers." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkPeersDataSourceConfig(dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "data.#", "0"),
				),
			},
		},
	})
}

func testAccNetworkPeersDataSourceConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_network_peers" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dataSourceName, globalOrgId, globalProjectId, globalClusterId)
}


func TestAccDatasourceNetworkPeersMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "organization_id",
			body: fmt.Sprintf(`
  project_id = "%s"
  cluster_id = "%s"
`, globalProjectId, globalClusterId),
		},
		{
			name: "project_id",
			body: fmt.Sprintf(`
  organization_id = "%s"
  cluster_id      = "%s"
`, globalOrgId, globalClusterId),
		},
		{
			name: "cluster_id",
			body: fmt.Sprintf(`
  organization_id = "%s"
  project_id      = "%s"
`, globalOrgId, globalProjectId),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			dataSourceName := randomStringWithPrefix("tf_acc_network_peers_ds_missing_")
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_network_peers" "%[2]s" {%[3]s}
`, globalProviderBlock, dataSourceName, test.body),
						ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, test.name)),
					},
				},
			})
		})
	}
}


func TestAccDatasourceAzureNetworkPeerCommandInvalidCluster(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_azure_np_command_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureNetworkPeerCommandConfig(dataSourceName, "22222222-2222-2222-2222-222222222222"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Azure network peer command|Could not read Azure network peer command`),
			},
		},
	})
}


func TestAccDatasourceAzureNetworkPeerCommandMissingRequiredFields(t *testing.T) {
	for _, attr := range []string{"tenant_id", "subscription_id", "resource_group", "vnet_id", "vnet_peering_service_principal"} {
		attr := attr
		t.Run(attr, func(t *testing.T) {
			dataSourceName := randomStringWithPrefix("tf_acc_azure_np_command_missing_")
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config:      testAccAzureNetworkPeerCommandConfigOmitting(dataSourceName, attr),
						ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, attr)),
					},
				},
			})
		})
	}
}

func testAccAzureNetworkPeerCommandConfig(dataSourceName, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_azure_network_peer_command" "%[2]s" {
  organization_id                = "%[3]s"
  project_id                     = "%[4]s"
  cluster_id                     = "%[5]s"
  tenant_id                      = "33333333-3333-3333-3333-333333333333"
  subscription_id                = "44444444-4444-4444-4444-444444444444"
  resource_group                 = "tf-acc-resource-group"
  vnet_id                        = "tf-acc-vnet"
  vnet_peering_service_principal = "55555555-5555-5555-5555-555555555555"
}
`, globalProviderBlock, dataSourceName, globalOrgId, globalProjectId, clusterID)
}


func testAccAzureNetworkPeerCommandConfigOmitting(dataSourceName, omit string) string {
	attrs := map[string]string{
		"organization_id":                globalOrgId,
		"project_id":                     globalProjectId,
		"cluster_id":                     globalClusterId,
		"tenant_id":                      "33333333-3333-3333-3333-333333333333",
		"subscription_id":                "44444444-4444-4444-4444-444444444444",
		"resource_group":                 "tf-acc-resource-group",
		"vnet_id":                        "tf-acc-vnet",
		"vnet_peering_service_principal": "55555555-5555-5555-5555-555555555555",
	}
	delete(attrs, omit)

	var body string
	// Iterate a fixed key order so the generated config is stable across runs.
	for _, k := range []string{
		"organization_id", "project_id", "cluster_id", "tenant_id",
		"subscription_id", "resource_group", "vnet_id", "vnet_peering_service_principal",
	} {
		if v, ok := attrs[k]; ok {
			body += fmt.Sprintf("  %s = %q\n", k, v)
		}
	}

	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_azure_network_peer_command" "%[2]s" {
%[3]s}
`, globalProviderBlock, dataSourceName, body)
}
