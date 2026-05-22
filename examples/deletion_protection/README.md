# Deletion Protection Example

This example shows how to manage deletion protection on an existing Couchbase Capella cluster using the dedicated `couchbase-capella_cluster_deletion_protection` resource.

The resource calls `PUT /v4/.../deletionProtection` to set the desired state, then reads the cluster to confirm the value.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Enable deletion protection on an existing cluster.
2. UPDATE: Toggle deletion protection off.
3. DESTROY: Remove the resource from state (does not alter the cluster).

If you check the `terraform.template.tfvars` file — copy it to `terraform.tfvars` and update the values with your organization credentials.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform will perform the following actions:

  # couchbase-capella_cluster_deletion_protection.cluster will be created
  + resource "couchbase-capella_cluster_deletion_protection" "cluster" {
      + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + deletion_protection = true
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the Plan

Command: `terraform apply`

Sample Output:
```
$ terraform apply

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

deletion_protection = true
```

## UPDATE — Disable protection

Set `deletion_protection = false` in `terraform.tfvars`, then apply:

Command: `terraform apply`

Sample Output:
```
$ terraform apply

  # couchbase-capella_cluster_deletion_protection.cluster will be updated in-place
  ~ resource "couchbase-capella_cluster_deletion_protection" "cluster" {
      ~ deletion_protection = true -> false
    }

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

deletion_protection = false
```

