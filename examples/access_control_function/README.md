# Capella Access Control Function Example

This example shows how to create and manage access control and validation functions for App Endpoints in Capella.

This creates a new access control function for a specific collection within an App Endpoint. Access control functions are JavaScript functions that specify access control policies applied to documents in collections. Every document update is processed by this function.

The default access control function is `function(doc){channel(doc.channels);}` for the default collection and `function(doc){channel(collectionName);}` for named collections.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

1. CREATE: Create a new access control function in an existing Capella App Endpoint as stated in the `create_access_control_function.tf` file.
2. UPDATE: Update the app service configuration using Terraform.
3. LIST: List existing app services in Capella as stated in the `list_access_functions.tf` file.
4. IMPORT: Import an app services that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created app service from Capella.

## Example Walkthrough

### 1. CREATE: View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan                                                                                                                                                                                                                   1 ↵
2025-08-27T14:52:08.917+0100 [INFO]  Terraform version: 1.12.1
2025-08-27T14:52:08.917+0100 [INFO]  Go runtime version: go1.24.2
2025-08-27T14:52:08.917+0100 [INFO]  CLI args: []string{"terraform", "plan"}
2025-08-27T14:52:08.917+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-08-27T14:52:08.918+0100 [INFO]  CLI command args: []string{"plan"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2025-08-27T14:52:08.922+0100 [INFO]  backend/local: starting Plan operation
2025-08-27T14:52:08.925+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:52:08.945+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:52:08.945+0100"
2025-08-27T14:52:08.966+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33634
2025-08-27T14:52:08.967+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:52:08.982+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:52:08.981+0100"
2025-08-27T14:52:08.995+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33635
2025-08-27T14:52:08.995+0100 [INFO]  backend/local: plan calling Plan
2025-08-27T14:52:08.996+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:52:09.007+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:52:09.007+0100"
2025-08-27T14:52:09.019+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=1aca2759-ae4a-2826-92fc-05f53713f177 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 @module=couchbase_capella timestamp="2025-08-27T14:52:09.019+0100"
2025-08-27T14:52:09.019+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: @module=couchbase_capella host=http://cloudapi.dev.nonprod-project-avengers.com tf_req_id=1aca2759-ae4a-2826-92fc-05f53713f177 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 authentication_token="***" success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp="2025-08-27T14:52:09.019+0100"
couchbase-capella_access_control_function.example_access_function: Refreshing state...
2025-08-27T14:52:09.022+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout=""
2025-08-27T14:52:09.637+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33636
2025-08-27T14:52:09.639+0100 [INFO]  backend/local: plan operation completed

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_access_control_function.example_access_function will be updated in-place
  ~ resource "couchbase-capella_access_control_function" "example_access_function" {
      + access_control_function = "function(doc){channel(doc.channels);}"
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Access Control Function in Capella

Command: `terraform apply`

Sample Output:
```
 $ terraform apply                                                                                                                                                                                                                  1 ↵
2025-08-27T14:51:31.820+0100 [INFO]  Terraform version: 1.12.1
2025-08-27T14:51:31.820+0100 [INFO]  Go runtime version: go1.24.2
2025-08-27T14:51:31.820+0100 [INFO]  CLI args: []string{"terraform", "apply"}
2025-08-27T14:51:31.820+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-08-27T14:51:31.821+0100 [INFO]  CLI command args: []string{"apply"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2025-08-27T14:51:31.825+0100 [INFO]  backend/local: starting Apply operation
2025-08-27T14:51:31.828+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:51:31.854+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:51:31.853+0100"
2025-08-27T14:51:31.876+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33563
2025-08-27T14:51:31.876+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:51:31.889+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:51:31.888+0100"
2025-08-27T14:51:31.903+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33564
2025-08-27T14:51:31.903+0100 [INFO]  backend/local: apply calling Plan
2025-08-27T14:51:31.903+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:51:31.918+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:51:31.917+0100"
2025-08-27T14:51:31.930+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=e89b3187-2f63-d8ed-4057-dbd7ee473e9e tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 timestamp="2025-08-27T14:51:31.930+0100"
2025-08-27T14:51:31.930+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=e89b3187-2f63-d8ed-4057-dbd7ee473e9e tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella host=http://cloudapi.dev.nonprod-project-avengers.com authentication_token="***" timestamp="2025-08-27T14:51:31.930+0100"
2025-08-27T14:51:31.933+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33565

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_access_control_function.example_access_function will be created
  + resource "couchbase-capella_access_control_function" "example_access_function" {
      + access_control_function = "function(doc){channel(doc.channels);}"
      + app_endpoint_name       = "api"
      + access_control_function_id          = "9c8fba79-c79b-458c-b833-1a99373277de"
      + cluster_id              = "439971e0-afde-4a11-bc0d-c5524e839680"
      + collection              = "_default"
      + organization_id         = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id              = "edc84e13-420f-4a8d-95a9-7939ff573a01"
      + scope                   = "_default"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

2025-08-27T14:51:32.944+0100 [INFO]  backend/local: apply calling Apply
2025-08-27T14:51:32.952+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:51:32.979+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:51:32.978+0100"
2025-08-27T14:51:32.995+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=4d607e06-a8bd-701d-a8d6-469544be52a6 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 timestamp="2025-08-27T14:51:32.995+0100"
2025-08-27T14:51:32.995+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella success=true tf_req_id=4d607e06-a8bd-701d-a8d6-469544be52a6 timestamp="2025-08-27T14:51:32.995+0100"
couchbase-capella_access_control_function.example_access_function: Creating...
2025-08-27T14:51:33.001+0100 [INFO]  Starting apply for couchbase-capella_access_control_function.example_access_function
2025-08-27T14:51:33.003+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout="function(doc){channel(doc.channels);}"
2025-08-27T14:51:34.013+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout=""
couchbase-capella_access_control_function.example_access_function: Creation complete after 1s
2025-08-27T14:51:34.381+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33568

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

### View the current list of resources that are present in Terraform State

Command: `terraform show`

Sample Output:
```
 $ terraform show
2025-08-27T14:52:59.464+0100 [INFO]  Terraform version: 1.12.1
2025-08-27T14:52:59.465+0100 [INFO]  Go runtime version: go1.24.2
2025-08-27T14:52:59.465+0100 [INFO]  CLI args: []string{"terraform", "show"}
2025-08-27T14:52:59.465+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-08-27T14:52:59.466+0100 [INFO]  CLI command args: []string{"show"}
2025-08-27T14:52:59.477+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:52:59.522+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:52:59.521+0100"
2025-08-27T14:52:59.552+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=33675
# couchbase-capella_access_control_function.example_access_function:
resource "couchbase-capella_access_control_function" "example_access_function" {
    access_control_function = "function(doc){channel(doc.channels);}"
    app_endpoint_name       = "api"
    access_control_function_id          = "9c8fba79-c79b-458c-b833-1a99373277de"
    cluster_id              = "439971e0-afde-4a11-bc0d-c5524e839680"
    collection              = "_default"
    organization_id         = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
    project_id              = "edc84e13-420f-4a8d-95a9-7939ff573a01"
    scope                   = "_default"
}
```


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
 $ terraform state list                                                                
2025-08-27T14:56:48.271+0100 [INFO]  Terraform version: 1.12.1
2025-08-27T14:56:48.272+0100 [INFO]  Go runtime version: go1.24.2
2025-08-27T14:56:48.272+0100 [INFO]  CLI args: []string{"terraform", "state", "list"}
2025-08-27T14:56:48.272+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-08-27T14:56:48.274+0100 [INFO]  CLI command args: []string{"state", "list"}

```


### 2. Update the access control function

Make changes to the `access_control_function` attribute in `create_access_control_function.tf` and run:

Command: `terraform apply`

Sample Output:
```
 $ terraform apply                                                                                                                                                                                                                  1 ↵
2025-08-27T14:54:33.273+0100 [INFO]  Terraform version: 1.12.1
2025-08-27T14:54:33.273+0100 [INFO]  Go runtime version: go1.24.2
2025-08-27T14:54:33.273+0100 [INFO]  CLI args: []string{"terraform", "apply"}
2025-08-27T14:54:33.273+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-08-27T14:54:33.275+0100 [INFO]  CLI command args: []string{"apply"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2025-08-27T14:54:33.286+0100 [INFO]  backend/local: starting Apply operation
2025-08-27T14:54:33.292+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:54:33.854+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:54:33.853+0100"
2025-08-27T14:54:33.885+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34342
2025-08-27T14:54:33.891+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:54:34.162+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:54:34.161+0100"
2025-08-27T14:54:34.180+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34345
2025-08-27T14:54:34.180+0100 [INFO]  backend/local: apply calling Plan
2025-08-27T14:54:34.181+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:54:34.197+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:54:34.197+0100"
2025-08-27T14:54:34.213+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=c25d584d-cd08-5c6f-c55d-02c62d4d63b1 timestamp="2025-08-27T14:54:34.213+0100"
2025-08-27T14:54:34.213+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: @module=couchbase_capella authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com success=true tf_req_id=c25d584d-cd08-5c6f-c55d-02c62d4d63b1 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp="2025-08-27T14:54:34.213+0100"
couchbase-capella_access_control_function.example_access_function: Refreshing state...
2025-08-27T14:54:34.217+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout=""
2025-08-27T14:54:34.832+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34346

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_access_control_function.example_access_function will be updated in-place
  ~ resource "couchbase-capella_access_control_function" "example_access_function" {
      + access_control_function = "function(doc){channel(doc.channels);console.log(doc);}"
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

2025-08-27T14:54:35.754+0100 [INFO]  backend/local: apply calling Apply
2025-08-27T14:54:35.758+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:54:35.787+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:54:35.786+0100"
2025-08-27T14:54:35.805+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_req_id=86c57f01-3879-1edf-da19-a32fc0b4cf05 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp="2025-08-27T14:54:35.805+0100"
2025-08-27T14:54:35.805+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella authentication_token="***" success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella host=http://cloudapi.dev.nonprod-project-avengers.com tf_req_id=86c57f01-3879-1edf-da19-a32fc0b4cf05 tf_rpc=ConfigureProvider timestamp="2025-08-27T14:54:35.805+0100"
couchbase-capella_access_control_function.example_access_function: Modifying...
2025-08-27T14:54:35.809+0100 [INFO]  Starting apply for couchbase-capella_access_control_function.example_access_function
2025-08-27T14:54:35.810+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout="function(doc){channel(doc.channels);console.log(doc);}"
2025-08-27T14:54:36.613+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout=""
couchbase-capella_access_control_function.example_access_function: Modifications complete after 1s
2025-08-27T14:54:36.976+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34347

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```


## IMPORT
### Remove the resource `example_access_function` from the Terraform State file

Command: `terraform state rm couchbase-capella_access_control_function.example_access_function`

Sample Output:
``` 
terraform state rm couchbase-capella_access_control_function.example_access_function
Removed capella_access_control_function.example_access_function
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_access_control_function.example_access_function id=<appservice_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_access_control_function.example_access_function id=e57494b4-c791-44b6-8e7a-ee20db89a7f0,cluster_id=b74f5350-f727-427e-8cad-623a691b1cfe,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d`

Sample Output:
``` 
terraform import 
``` 

### 5. Delete the resources that Terraform manages

Command: `terraform destroy`

Sample Output:
```
 $ terraform destroy
2025-08-27T14:55:36.262+0100 [INFO]  Terraform version: 1.12.1
2025-08-27T14:55:36.263+0100 [INFO]  Go runtime version: go1.24.2
2025-08-27T14:55:36.263+0100 [INFO]  CLI args: []string{"terraform", "destroy"}
2025-08-27T14:55:36.263+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-08-27T14:55:36.264+0100 [INFO]  CLI command args: []string{"destroy"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2025-08-27T14:55:36.274+0100 [INFO]  backend/local: starting Apply operation
2025-08-27T14:55:36.280+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:55:36.336+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:55:36.336+0100"
2025-08-27T14:55:36.367+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34379
2025-08-27T14:55:36.372+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:55:36.386+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:55:36.386+0100"
2025-08-27T14:55:36.403+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34381
2025-08-27T14:55:36.403+0100 [INFO]  backend/local: apply calling Plan
2025-08-27T14:55:36.405+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:55:36.418+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:55:36.417+0100"
2025-08-27T14:55:36.430+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=0b2bd863-1cef-edb6-b0cf-d928244204ba tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 timestamp="2025-08-27T14:55:36.430+0100"
2025-08-27T14:55:36.431+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella authentication_token="***" tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella host=http://cloudapi.dev.nonprod-project-avengers.com success=true tf_req_id=0b2bd863-1cef-edb6-b0cf-d928244204ba timestamp="2025-08-27T14:55:36.430+0100"
couchbase-capella_access_control_function.example_access_function: Refreshing state...
2025-08-27T14:55:36.435+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout=""
2025-08-27T14:55:37.111+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34382
2025-08-27T14:55:37.112+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:55:37.142+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:55:37.141+0100"
2025-08-27T14:55:37.157+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=8d05a9a6-6357-426c-0864-f58434f256a3 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 @module=couchbase_capella timestamp="2025-08-27T14:55:37.157+0100"
2025-08-27T14:55:37.157+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com tf_req_id=8d05a9a6-6357-426c-0864-f58434f256a3 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp="2025-08-27T14:55:37.157+0100"
2025-08-27T14:55:37.160+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34385

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_access_control_function.example_access_function will be destroyed
  - resource "couchbase-capella_access_control_function" "example_access_function" {
      - access_control_function = "function(doc){channel(doc.channels);console.log(doc);}" -> null
      - app_endpoint_name       = "api" -> null
      - access_control_function_id          = "9c8fba79-c79b-458c-b833-1a99373277de" -> null
      - cluster_id              = "439971e0-afde-4a11-bc0d-c5524e839680" -> null
      - collection              = "_default" -> null
      - organization_id         = "6af08c0a-8cab-4c1c-b257-b521575c16d0" -> null
      - project_id              = "edc84e13-420f-4a8d-95a9-7939ff573a01" -> null
      - scope                   = "_default" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

2025-08-27T14:55:38.821+0100 [INFO]  backend/local: apply calling Apply
2025-08-27T14:55:38.826+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-27T14:55:38.854+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-27T14:55:38.853+0100"
2025-08-27T14:55:38.871+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=9cc1b329-a63f-836a-f497-b5ceccb0ec1f tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 @module=couchbase_capella timestamp="2025-08-27T14:55:38.871+0100"
2025-08-27T14:55:38.871+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com tf_req_id=9cc1b329-a63f-836a-f497-b5ceccb0ec1f tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella @module=couchbase_capella timestamp="2025-08-27T14:55:38.871+0100"
couchbase-capella_access_control_function.example_access_function: Destroying...
2025-08-27T14:55:38.873+0100 [INFO]  Starting apply for couchbase-capella_access_control_function.example_access_function
2025-08-27T14:55:38.874+0100 [WARN]  unexpected data: registry.terraform.io/couchbasecloud/couchbase-capella:stdout=""
couchbase-capella_access_control_function.example_access_function: Destruction complete after 1s
2025-08-27T14:55:39.633+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=34387

Destroy complete! Resources: 1 destroyed.
```


## Prerequisites

- Couchbase Capella organization with appropriate permissions
- Existing App Service and App Endpoint
- Valid scope and collection within the App Endpoint

## Important Notes

- Access control functions are applied at the collection level
- Functions must be valid JavaScript
- Changes to access control functions affect all documents in the collection
- The resource requires `app_endpoint_name` rather than ID for identification 