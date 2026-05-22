# Cluster Deletion Protection Example

This example shows how to manage deletion protection on an existing Couchbase Capella cluster using the `couchbase-capella_cluster_deletion_protection` resource.

> **Note:** Do not set `deletion_protection` on the `couchbase-capella_cluster` resource when using this resource. The cluster resource treats `deletion_protection` as read-only (computed). This resource is the sole owner of that field.

## Example Walkthrough

1. **CREATE** — Enable deletion protection on an existing cluster.
2. **UPDATE** — Toggle deletion protection off.
3. **DESTROY** — Remove the resource from state (does not alter the cluster).

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the values with your credentials.

## CREATE

### View the plan

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

### Apply the plan

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

## DESTROY

Removing this resource from state does **not** disable deletion protection on the cluster. This is intentional — protection remains enabled as a safety measure.

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy

  # couchbase-capella_cluster_deletion_protection.cluster will be destroyed (state only)

Destroy complete! Resources: 1 destroyed.
```

If you need to delete the cluster afterwards, first set `deletion_protection = false` and apply, then destroy.
