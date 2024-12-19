# Capella Azure VNET Peering CLI Command Example

This example shows how to retrieve the Azure role assignment command to be run in the Azure CLI that is used to configure a network peer.

Note: Before retrieving this Azure role assignment command, please make sure that the Admin consent granting process has been completed through the Capella UI. 
For more information, please refer to the [steps here](https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/Network-Peers/operation/getAzureVnetPeeringCommand)

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Display the Azure network peer cli command as stated in the `get_network_peer_command.tf` file.

## GET

Command: `terraform apply`

Sample Output:
```
 terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_azure_network_peer_command.azure_network_peer_command: Reading...
data.couchbase-capella_azure_network_peer_command.azure_network_peer_command: Read complete after 0s

Changes to Outputs:
  + azure_network_peer_command = {
      + cluster_id                     = "ffffffff-aaaa-1414-eeee-000000000000"
      + command                        = "az role assignment create --assignee-object-id ffffffff-aaaa-1414-eeee-000000000000 --role \"Network Contributor\" --scope /subscriptions/ffffffff-aaaa-1414-eeee-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/VirtualNetworks/test_vnet --assignee-principal-type ServicePrincipal"
      + organization_id                = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                     = "ffffffff-aaaa-1414-eeee-000000000000"
      + resource_group                 = "test_rg"
      + subscription_id                = "ffffffff-aaaa-1414-eeee-000000000000"
      + tenant_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
      + vnet_id                        = "test_vnet"
      + vnet_peering_service_principal = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

azure_network_peer_command = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "command" = "az role assignment create --assignee-object-id ffffffff-aaaa-1414-eeee-000000000000 --role \"Network Contributor\" --scope /subscriptions/ffffffff-aaaa-1414-eeee-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/VirtualNetworks/test_vnet --assignee-principal-type ServicePrincipal"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "resource_group" = "test_rg"
  "subscription_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "tenant_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "vnet_id" = "test_vnet"
  "vnet_peering_service_principal" = "ffffffff-aaaa-1414-eeee-000000000000"
}

```