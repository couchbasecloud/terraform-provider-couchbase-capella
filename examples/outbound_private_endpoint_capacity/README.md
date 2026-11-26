# Get Outbound Private Endpoint Capacity Data source

This example shows how to get how many outbound private endpoints a cluster has against its limit. Every endpoint counts, including pending, failed and deleting ones. The limit is 3 unless Couchbase Support raises it. It uses the organization ID, project ID and cluster ID to do so.

The cluster must be on. Reading the capacity of a turned-off cluster fails with a 422.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_outbound_private_endpoint_capacity.existing_capacity: Reading...
data.couchbase-capella_outbound_private_endpoint_capacity.existing_capacity: Read complete after 1s

Changes to Outputs:
  + existing_capacity = {
      + cluster_endpoint_limit = {
          + limit = 3
          + used  = 1
        }
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.
```

### Apply the Plan

Command: `terraform apply`

Sample Output:

```
$ terraform apply
data.couchbase-capella_outbound_private_endpoint_capacity.existing_capacity: Reading...
data.couchbase-capella_outbound_private_endpoint_capacity.existing_capacity: Read complete after 1s

Changes to Outputs:
  + existing_capacity = {
      + cluster_endpoint_limit = {
          + limit = 3
          + used  = 1
        }
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_capacity = {
  "cluster_endpoint_limit" = {
    "limit" = 3
    "used" = 1
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```
