(base) mohammedmadi@ABCDEFGH:~/terraform  $ terraform apply                                                                                                                                                                                                                                                                                                                                         1 ↵
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
  + new_allowedcidr = {
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

new_allowedcidr = {
  "admin_url" = "https://ppfa8uazpglnzqp5.apps.nonprod-project-avengers.com:4985/example-app-endpoint"
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
  "metrics_url" = "https://ppfa8uazpglnzqp5.apps.nonprod-project-avengers.com:4988/metrics"
  "name" = "example-app-endpoint"
  "oidc" = toset(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "public_url" = "wss://ppfa8uazpglnzqp5.apps.nonprod-project-avengers.com:4984/example-app-endpoint"
  "require_resync" = tomap(null) /* of list of string */
  "scope" = "scope1"
  "user_xattr_key" = tostring(null)
}
