# Capella App Services Allowed CIDR Example

This example shows how to create and manage App Endpoints on an App Service in Capella.

This creates a new App Endpoint in the selected Capella App Service. It uses the App Service ID and given App Endpoint config to create and list App Endpoints.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In the following example, we will:

1. CREATE: Create a new App Endpoint, as stated in the `create_app_endpoint.tf` file.
2. LIST: List existing App Endpoints in Capella as stated in the `list_app_endpoints.tf` file.
3. IMPORT: Import an App Endpoint that exists in Capella but not in the terraform state file.
4. DELETE: Delete an App Endpoint from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`
```bash
 $ terraform plan               
2025-07-31T15:58:05.657+0100 [INFO]  Terraform version: 1.12.1
2025-07-31T15:58:05.658+0100 [INFO]  Go runtime version: go1.24.2
2025-07-31T15:58:05.658+0100 [INFO]  CLI args: []string{"terraform", "plan"}
2025-07-31T15:58:05.658+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-07-31T15:58:05.658+0100 [INFO]  CLI command args: []string{"plan"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2025-07-31T15:58:05.662+0100 [INFO]  backend/local: starting Plan operation
2025-07-31T15:58:05.665+0100 [INFO]  provider: configuring client automatic mTLS
2025-07-31T15:58:05.688+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-07-31T15:58:05.687+0100"
2025-07-31T15:58:05.710+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=51556
2025-07-31T15:58:05.711+0100 [INFO]  provider: configuring client automatic mTLS
2025-07-31T15:58:05.725+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-07-31T15:58:05.725+0100"
2025-07-31T15:58:05.739+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=51557
2025-07-31T15:58:05.739+0100 [INFO]  backend/local: plan calling Plan
2025-07-31T15:58:05.739+0100 [INFO]  provider: configuring client automatic mTLS
2025-07-31T15:58:05.752+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-07-31T15:58:05.752+0100"
2025-07-31T15:58:05.763+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 timestamp="2025-07-31T15:58:05.763+0100"
2025-07-31T15:58:05.764+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: @module=couchbase_capella host=http://cloudapi.dev.nonprod-project-avengers.com success=true tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 authentication_token="***" tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp="2025-07-31T15:58:05.764+0100"
2025-07-31T15:58:05.768+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=51558
2025-07-31T15:58:05.769+0100 [INFO]  backend/local: plan operation completed

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint.example_app_endpoint will be created
  + resource "couchbase-capella_app_endpoint" "example_app_endpoint" {
      + admin_url          = (known after apply)
      + app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + bucket             = "example-bucket"
      + cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections        = {
          + "collection1" = {},
        }
      + delta_sync_enabled = false
      + metrics_url        = (known after apply)
      + name               = "example-app-endpoint"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + public_url         = (known after apply)
      + require_resync     = (known after apply)
      + scope              = "scope1"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_endpoint = {
      + admin_url          = (known after apply)
      + app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + bucket             = "example-bucket"
      + cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections        = {
          + collection1 = {
              + access_control_function = null
              + import_filter           = null
            }
        }
      + cors               = null
      + delta_sync_enabled = false
      + metrics_url        = (known after apply)
      + name               = "example-app-endpoint"
      + oidc               = null
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + public_url         = (known after apply)
      + require_resync     = (known after apply)
      + scope              = "scope1"
      + user_xattr_key     = null
    }

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new App Endpoint

Command: `terraform apply`

Sample Output:
```bash
$ terraform apply                                                                                                                                                                                                                                                                                                                                         1 ↵
2025-07-31T15:47:09.851+0100 [INFO]  Terraform version: 1.12.1
2025-07-31T15:47:09.851+0100 [INFO]  Go runtime version: go1.24.2
2025-07-31T15:47:09.851+0100 [INFO]  CLI args: []string{"terraform", "apply"}
2025-07-31T15:47:09.852+0100 [INFO]  Loading CLI configuration from ./.terraformrc
2025-07-31T15:47:09.852+0100 [INFO]  CLI command args: []string{"apply"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2025-07-31T15:47:09.857+0100 [INFO]  backend/local: starting Apply operation
2025-07-31T15:47:09.862+0100 [INFO]  provider: configuring client automatic mTLS
2025-07-31T15:47:10.412+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-07-31T15:47:10.411+0100"
2025-07-31T15:47:10.435+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=50352
2025-07-31T15:47:10.436+0100 [INFO]  provider: configuring client automatic mTLS
2025-07-31T15:47:10.705+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-07-31T15:47:10.705+0100"
2025-07-31T15:47:10.728+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=50355
2025-07-31T15:47:10.728+0100 [INFO]  backend/local: apply calling Plan
2025-07-31T15:47:10.729+0100 [INFO]  provider: configuring client automatic mTLS
2025-07-31T15:47:10.746+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-07-31T15:47:10.745+0100"
2025-07-31T15:47:10.761+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 timestamp="2025-07-31T15:47:10.761+0100"
2025-07-31T15:47:10.761+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella success=true tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp="2025-07-31T15:47:10.761+0100"
2025-07-31T15:47:10.767+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=50356

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint.example_app_endpoint will be created
  + resource "couchbase-capella_app_endpoint" "example_app_endpoint" {
      + admin_url          = (known after apply)
      + app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + bucket             = "example-bucket"
      + cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections        = {
          + "collection1" = {},
        }
      + delta_sync_enabled = false
      + metrics_url        = (known after apply)
      + name               = "example-app-endpoint"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + public_url         = (known after apply)
      + require_resync     = (known after apply)
      + scope              = "scope1"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_endpoint = {
      + admin_url          = (known after apply)
      + app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + bucket             = "example-bucket"
      + cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections        = {
          + collection1 = {
              + access_control_function = null
              + import_filter           = null
            }
        }
      + cors               = null
      + delta_sync_enabled = false
      + metrics_url        = (known after apply)
      + name               = "example-app-endpoint"
      + oidc               = null
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + public_url         = (known after apply)
      + require_resync     = (known after apply)
      + scope              = "scope1"
      + user_xattr_key     = null
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

2025-07-31T15:47:12.119+0100 [INFO]  backend/local: apply calling Apply
2025-07-31T15:47:12.121+0100 [INFO]  provider: configuring client automatic mTLS
2025-07-31T15:47:12.146+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-07-31T15:47:12.145+0100"
2025-07-31T15:47:12.164+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp="2025-07-31T15:47:12.164+0100"
2025-07-31T15:47:12.165+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com timestamp="2025-07-31T15:47:12.165+0100"
couchbase-capella_app_endpoint.example_app_endpoint: Creating...
2025-07-31T15:47:12.169+0100 [INFO]  Starting apply for couchbase-capella_app_endpoint.example_app_endpoint
couchbase-capella_app_endpoint.example_app_endpoint: Creation complete after 1s [name=example-app-endpoint]
2025-07-31T15:47:13.286+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=50358

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_app_endpoint = {
  "admin_url" = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4985/example-app-endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "bucket" = "example-bucket"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections" = tomap({
    "collection1" = {
      "access_control_function" = tostring(null)
      "import_filter" = tostring(null)
    }
  })
  "cors" = null /* object */
  "delta_sync_enabled" = false
  "metrics_url" = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4988/metrics"
  "name" = "example-app-endpoint"
  "oidc" = toset(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "public_url" = "wss://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4984/example-app-endpoint"
  "require_resync" = tomap(null) /* of list of string */
  "scope" = "scope1"
  "user_xattr_key" = tostring(null)
}
```


### Note the App Endoint read only fields populated after `terraform apply`
Command: `terraform output new_app_endpoint`

Sample Output:
```bash
$ terraform output new_app_endpoint
2025-07-31T17:20:06.039+0100 [INFO]  Terraform version: 1.12.1
2025-07-31T17:20:06.039+0100 [INFO]  Go runtime version: go1.24.2
2025-07-31T17:20:06.039+0100 [INFO]  CLI args: []string{"terraform", "output", "new_app_endpoint"}
2025-07-31T17:20:06.039+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-07-31T17:20:06.040+0100 [INFO]  CLI command args: []string{"output", "new_app_endpoint"}
{
  "admin_url" = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4985/example-app-endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "bucket" = "example-bucket"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections" = tomap({
    "collection1" = {
      "access_control_function" = tostring(null)
      "import_filter" = tostring(null)
    }
  })
  "cors" = null /* object */
  "delta_sync_enabled" = false
  "metrics_url" = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4988/metrics"
  "name" = "example-app-endpoint"
  "oidc" = toset(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "public_url" = "wss://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4984/example-app-endpoint"
  "require_resync" = tomap(null) /* of list of string */
  "scope" = "scope1"
  "user_xattr_key" = tostring(null)
}
```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```bash
 $ terraform state list
2025-07-31T17:20:20.841+0100 [INFO]  Terraform version: 1.12.1
2025-07-31T17:20:20.841+0100 [INFO]  Go runtime version: go1.24.2
2025-07-31T17:20:20.841+0100 [INFO]  CLI args: []string{"terraform", "state", "list"}
2025-07-31T17:20:20.841+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-07-31T17:20:20.842+0100 [INFO]  CLI command args: []string{"state", "list"}
couchbase-capella_app_endpoint.example_app_endpoint
```

### Remove the resource `new_app_endpoint` from the Terraform State file
 
Command: `terraform state rm couchbase-capella_app_endpoint.new_app_endpoint`

Sample Output:
```bash
$ terraform state rm couchbase-capella_app_endpoint.example_app_endpoint                                                                                                                                                                                                                                                                                  1 ↵
2025-07-31T17:46:52.487+0100 [INFO]  Terraform version: 1.12.1
2025-07-31T17:46:52.487+0100 [INFO]  Go runtime version: go1.24.2
2025-07-31T17:46:52.487+0100 [INFO]  CLI args: []string{"terraform", "state", "rm", "couchbase-capella_app_endpoint.example_app_endpoint"}
2025-07-31T17:46:52.487+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-07-31T17:46:52.488+0100 [INFO]  CLI command args: []string{"state", "rm", "couchbase-capella_app_endpoint.example_app_endpoint"}
Removed couchbase-capella_app_endpoint.example_app_endpoint
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file. The resource will still exist in Capella.

## IMPORT

Command:  `terraform import couchbase-capella_app_endpoint.example_app_endpoint name=example-app-endpoint,organization_id=<organization_id>,project_id=<project_id>,app_service_id=<app_service_id>,cluster_id=<cluster_id>`


```aiignore
 $ terraform import  couchbase-capella_app_endpoint.example_app_endpoint name=example-app-endpoint,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000                                                                                                                                                                                                                                       1 ↵
couchbase-capella_app_endpoint.example_app_endpoint: Importing from ID "name=example-app-endpoint,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint.example_app_endpoint: Import prepared!
  Prepared couchbase-capella_app_endpoint for import
couchbase-capella_app_endpoint.example_app_endpoint: Refreshing state... [name=name=example-app-endpoint,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.


```
## Update

Command: `terraform apply -var delta_sync_enabled=true'`

Sample output:
```bash
 $ terraform apply -var delta_sync_enabled=true                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   1 ↵
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_app_endpoint.example_app_endpoint: Refreshing state... [name=example-app-endpoint]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

example_app_endpoint = {
  "admin_url" = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4985/example-app-endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "bucket" = "example-bucket"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections" = tomap({
    "collection1" = {
      "access_control_function" = "function (doc, oldDoc, meta) {channel('collection1');}"
      "import_filter" = tostring(null)
    }
  })
  "cors" = {
    "disabled" = false
    "headers" = tolist(null) /* of string */
    "login_origin" = tolist([
      "https://login.example.com",
    ])
    "max_age" = 0
    "origin" = tolist([
      "https://example.com",
    ])
  }
  "delta_sync_enabled" = true
  "metrics_url" = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4988/metrics"
  "name" = "example-app-endpoint"
  "oidc" = tolist([
    {
      "client_id" = "example-client"
      "discovery_url" = "https://accounts.google.com/.well-known/openid-configuration"
      "is_default" = true
      "issuer" = "https://accounts.google.com"
      "provider_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "register" = false
      "roles_claim" = ""
      "user_prefix" = "example-prefix2"
      "username_claim" = ""
    },
    {
      "client_id" = "example-client"
      "discovery_url" = "https://accounts.google.com/.well-known/openid-configuration"
      "is_default" = false
      "issuer" = "https://accounts.google.com"
      "provider_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "register" = false
      "roles_claim" = ""
      "user_prefix" = "example-prefix22"
      "username_claim" = ""
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "public_url" = "wss://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4984/example-app-endpoint"
  "require_resync" = tomap(null) /* of list of string */
  "scope" = "scope1"
  "state" = "Offline"
  "user_xattr_key" = "user_xattr"
}

```

## DELETE
### Delete the App Endpoint from Capella
Command: `terraform destroy`
```bash
$ terraform destroy
2025-08-01T11:10:53.788+0100 [INFO]  Terraform version: 1.12.1
2025-08-01T11:10:53.789+0100 [INFO]  Go runtime version: go1.24.2
2025-08-01T11:10:53.789+0100 [INFO]  CLI args: []string{"terraform", "destroy"}
2025-08-01T11:10:53.789+0100 [INFO]  Loading CLI configuration from /Users/mohammedmadi/.terraformrc
2025-08-01T11:10:53.789+0100 [INFO]  CLI command args: []string{"destroy"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2025-08-01T11:10:53.794+0100 [INFO]  backend/local: starting Apply operation
2025-08-01T11:10:53.800+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-01T11:10:53.818+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-01T11:10:53.818+0100"
2025-08-01T11:10:53.840+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=57647
2025-08-01T11:10:53.841+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-01T11:10:53.852+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-01T11:10:53.852+0100"
2025-08-01T11:10:53.866+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=57648
2025-08-01T11:10:53.866+0100 [INFO]  backend/local: apply calling Plan
2025-08-01T11:10:53.866+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-01T11:10:53.880+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-01T11:10:53.879+0100"
2025-08-01T11:10:53.891+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 timestamp="2025-08-01T11:10:53.891+0100"
2025-08-01T11:10:53.891+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 success=true tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 @module=couchbase_capella authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com timestamp="2025-08-01T11:10:53.891+0100"
couchbase-capella_app_endpoint.example_app_endpoint: Refreshing state... [name=example-app-endpoint]
2025-08-01T11:10:54.533+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=57649
2025-08-01T11:10:54.534+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-01T11:10:54.554+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-01T11:10:54.554+0100"
2025-08-01T11:10:54.569+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 timestamp="2025-08-01T11:10:54.569+0100"
2025-08-01T11:10:54.569+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 host=http://cloudapi.dev.nonprod-project-avengers.com success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella @module=couchbase_capella authentication_token="***" timestamp="2025-08-01T11:10:54.569+0100"
2025-08-01T11:10:54.573+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=57652

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint.example_app_endpoint will be destroyed
  - resource "couchbase-capella_app_endpoint" "example_app_endpoint" {
      - admin_url          = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4985/example-app-endpoint" -> null
      - app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - bucket             = "example-bucket" -> null
      - cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collections        = {
          - "collection1" = {},
        } -> null
      - delta_sync_enabled = false -> null
      - metrics_url        = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4988/metrics" -> null
      - name               = "example-app-endpoint" -> null
      - organization_id    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id         = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - public_url         = "wss://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4984/example-app-endpoint" -> null
      - scope              = "scope1" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - example_app_endpoint = {
      - admin_url          = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4985/example-app-endpoint"
      - app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      - bucket             = "example-bucket"
      - cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      - collections        = {
          - collection1 = {
              - access_control_function = null
              - import_filter           = null
            }
        }
      - cors               = null
      - delta_sync_enabled = false
      - metrics_url        = "https://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4988/metrics"
      - name               = "example-app-endpoint"
      - oidc               = null
      - organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      - public_url         = "wss://aaaabbbbccccdddd.apps.nonprod-project-avengers.com:4984/example-app-endpoint"
      - require_resync     = null
      - scope              = "scope1"
      - user_xattr_key     = null
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

2025-08-01T11:10:56.627+0100 [INFO]  backend/local: apply calling Apply
2025-08-01T11:10:56.629+0100 [INFO]  provider: configuring client automatic mTLS
2025-08-01T11:10:56.653+0100 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp="2025-08-01T11:10:56.652+0100"
2025-08-01T11:10:56.670+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:75 @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider timestamp="2025-08-01T11:10:56.670+0100"
2025-08-01T11:10:56.670+0100 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: authentication_token="***" host=http://cloudapi.dev.nonprod-project-avengers.com @module=couchbase_capella success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=ffffffff-aaaa-1414-eeee-000000000000 tf_rpc=ConfigureProvider @caller=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/internal/provider/provider.go:171 timestamp="2025-08-01T11:10:56.670+0100"
couchbase-capella_app_endpoint.example_app_endpoint: Destroying... [name=example-app-endpoint]
2025-08-01T11:10:56.672+0100 [INFO]  Starting apply for couchbase-capella_app_endpoint.example_app_endpoint
couchbase-capella_app_endpoint.example_app_endpoint: Destruction complete after 0s
2025-08-01T11:10:57.327+0100 [INFO]  provider: plugin process exited: plugin=/Applications/gh_2.14.6_macOS_amd64/bin/terraform-provider-couchbase-capella/bin/terraform-provider-couchbase-capella id=57653

```

