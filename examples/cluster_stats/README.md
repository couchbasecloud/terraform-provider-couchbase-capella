# Capella Cluster Stats Example

This example shows how to retrieve statistics for a specific cluster in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. READ: Retrieve cluster statistics using the `couchbase-capella_cluster_stats` data source.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## READ
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

Plan: 0 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + cluster_stats = {
      + cluster_id         = "ffffffff-aaaa-1414-eeee-00000000000"
      + free_memory_in_mb  = 1871
      + max_replicas       = 2
      + organization_id    = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id         = "ffffffff-aaaa-1414-eeee-00000000000"
      + total_memory_in_mb = 2071
    }
```

### Apply the Plan

Command: `terraform apply`

Sample Output:
```
$ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

Plan: 0 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + cluster_stats = {
      + cluster_id         = "ffffffff-aaaa-1414-eeee-00000000000"
      + free_memory_in_mb  = 1871
      + max_replicas       = 2
      + organization_id    = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id         = "ffffffff-aaaa-1414-eeee-00000000000"
      + total_memory_in_mb = 2071
    }

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```
