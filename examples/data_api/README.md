# Capella Data API Example

This example shows how to enable or disable the Data API and network peering on an existing Couchbase Capella operational cluster using the `couchbase-capella_data_api` resource, and how to read the current Data API status using the `couchbase-capella_data_api` data source.

## Example Walkthrough

1. **CREATE** - Enable the Data API on an existing cluster.
2. **UPDATE** - Enable network peering for the Data API.
3. **IMPORT** - Import an existing Data API configuration into Terraform state.
4. **DESTROY** - Remove the resource from state (a no-op operation that does not alter the cluster Data API configuration).

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the values with your credentials.

## CREATE

### View the plan

Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform will perform the following actions:

  # couchbase-capella_data_api.new_data_api will be created
  + resource "couchbase-capella_data_api" "new_data_api" {
      + cluster_id                = "ffffffff-aaaa-1414-eeee-000000000000"
      + connection_string         = (known after apply)
      + enable_data_api           = true
      + enable_network_peering    = false
      + organization_id           = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                = "ffffffff-aaaa-1414-eeee-000000000000"
      + state_for_data_api        = (known after apply)
      + state_for_network_peering = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the plan

Enabling the Data API is an asynchronous operation, so the apply waits (up to 30 minutes) for the Data API to finish transitioning and reach a final state.

Command: `terraform apply`

Sample Output:
```
$ terraform apply

...

couchbase-capella_data_api.new_data_api: Creating...
couchbase-capella_data_api.new_data_api: Still creating... [00m10s elapsed]
...
couchbase-capella_data_api.new_data_api: Still creating... [08m00s elapsed]
couchbase-capella_data_api.new_data_api: Creation complete after 8m1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_data_api = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "connection_string" = "couchbases://cb.abcdefghijklmnop.cloud.couchbase.com"
  "enable_data_api" = true
  "enable_network_peering" = false
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "state_for_data_api" = "enabled"
  "state_for_network_peering" = "disabled"
}
```

## UPDATE - Enable network peering

Set `enable_network_peering = true` in `terraform.tfvars`, then apply:

Command: `terraform apply`

Sample Output:
```
$ terraform apply

...

  # couchbase-capella_data_api.new_data_api will be updated in-place
  ~ resource "couchbase-capella_data_api" "new_data_api" {
      ~ enable_network_peering    = false -> true
      ~ state_for_data_api        = "enabled" -> (known after apply)
      ~ state_for_network_peering = "disabled" -> (known after apply)
    }

...

couchbase-capella_data_api.new_data_api: Modifying...
couchbase-capella_data_api.new_data_api: Still modifying... [00m10s elapsed]
...
couchbase-capella_data_api.new_data_api: Still modifying... [01m00s elapsed]
couchbase-capella_data_api.new_data_api: Modifications complete after 1m0s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## IMPORT

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
terraform import couchbase-capella_data_api.new_data_api cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>
```

## DESTROY

Removing this resource from state does **not** modify the Data API on the cluster. The Data API and network peering are left in their current state. This destroy is purely a no-op.

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy

couchbase-capella_data_api.new_data_api: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

...

couchbase-capella_data_api.new_data_api: Destroying...
couchbase-capella_data_api.new_data_api: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```

If you need to disable the Data API, set `enable_data_api = false` and apply before destroying.
