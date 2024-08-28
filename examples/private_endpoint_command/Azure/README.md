# Capella Azure Command Example

This example shows how to retrieve the Azure command used to configure a private endpoint.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Display the AWS command as stated in the `get_command.tf` file.

## GET

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_azure_private_endpoint_command.azure_command: Reading...
data.couchbase-capella_azure_private_endpoint_command.azure_command: Read complete after 1s

Changes to Outputs:
  + azure_command = {
      + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + command             = <<-EOT
            echo This script is only compatible with BASH-like shells, not powershell or cmd.exe.
            echo Please ensure AZ CLI is installed and logged in prior to running this script.
            setopt interactivecomments 2>/dev/null
            # Create private endpoint
            az network private-endpoint create -g test-rg -n pl-ffffffff-aaaa-1414-eeee-000000000000 --vnet-name vnet-1 --subnet subnet-1 --private-connection-resource-id 'pl-ffffffff-aaaa-1414-eeee-000000000000.185b4ebd-c568-4625-a1be-03079958fb2c.eastus.azure.privatelinkservice' --connection-name pl-ffffffff-aaaa-1414-eeee-000000000000 -l eastus --manual-request true
            # Create DNS zone
            az network private-dns zone create -g test-rg -n private-endpoint.abcde.nonprod-project-avengers.com
            # Link DNS zone
            az network private-dns link vnet create -g test-rg -n dnslink-ffffffff-aaaa-1414-eeee-000000000000 -z private-endpoint.abcde.nonprod-project-avengers.com -v vnet-1 -e False
            # Create DNS A record
            ## Fetch and unquote NIC and IP address
            NIC=$(basename $(az network private-endpoint show -g test-rg -n pl-ffffffff-aaaa-1414-eeee-000000000000 --query "networkInterfaces[0].id"))
            NIC=${NIC//\"/} # Trim trailing quote
            IPADDRESS=$(az network nic ip-config list -g test-rg  --nic-name $NIC --query "[0].privateIPAddress")
            IPADDRESS=${IPADDRESS//\"/} # Trim leading and trailing quote
            ## Create the record
            az network private-dns record-set a create --resource-group test-rg --zone-name private-endpoint.abcde.nonprod-project-avengers.com --name '@'
            az network private-dns record-set a add-record --resource-group test-rg --zone-name private-endpoint.abcde.nonprod-project-avengers.com --record-set-name '@' -a $IPADDRESS
        EOT
      + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + resource_group_name = "test-rg"
      + virtual_network     = "vnet-1/subnet-1"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

azure_command = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "command" = <<-EOT

  echo This script is only compatible with BASH-like shells, not powershell or cmd.exe.
  echo Please ensure AZ CLI is installed and logged in prior to running this script.
  setopt interactivecomments 2>/dev/null
  # Create private endpoint
  az network private-endpoint create -g test-rg -n pl-ffffffff-aaaa-1414-eeee-000000000000 --vnet-name vnet-1 --subnet subnet-1 --private-connection-resource-id 'pl-ffffffff-aaaa-1414-eeee-000000000000.185b4ebd-c568-4625-a1be-03079958fb2c.eastus.azure.privatelinkservice' --connection-name pl-ffffffff-aaaa-1414-eeee-000000000000 -l eastus --manual-request true
  # Create DNS zone
  az network private-dns zone create -g test-rg -n private-endpoint.abcde.nonprod-project-avengers.com
  # Link DNS zone
  az network private-dns link vnet create -g test-rg -n dnslink-ffffffff-aaaa-1414-eeee-000000000000 -z private-endpoint.abcde.nonprod-project-avengers.com -v vnet-1 -e False
  # Create DNS A record
  ## Fetch and unquote NIC and IP address
  NIC=$(basename $(az network private-endpoint show -g test-rg -n pl-ffffffff-aaaa-1414-eeee-000000000000 --query "networkInterfaces[0].id"))
  NIC=${NIC//\"/} # Trim trailing quote
  IPADDRESS=$(az network nic ip-config list -g test-rg  --nic-name $NIC --query "[0].privateIPAddress")
  IPADDRESS=${IPADDRESS//\"/} # Trim leading and trailing quote
  ## Create the record
  az network private-dns record-set a create --resource-group test-rg --zone-name private-endpoint.abcde.nonprod-project-avengers.com --name '@'
  az network private-dns record-set a add-record --resource-group test-rg --zone-name private-endpoint.abcde.nonprod-project-avengers.com --record-set-name '@' -a $IPADDRESS

  EOT
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "resource_group_name" = "test-rg"
  "virtual_network" = "vnet-1/subnet-1"
}
```