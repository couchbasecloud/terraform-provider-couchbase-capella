# Deletion Protection Example

This example shows how to read and verify the deletion protection status of an existing Couchbase Capella cluster.

The cluster datasource fetches all attributes including `deletion_protection`. A `check` block asserts the current value matches the expected value from `terraform.tfvars`.

To update deletion protection, set `deletion_protection` on a managed `couchbase-capella_cluster` resource — the provider calls the dedicated `PUT /v4/.../deletionProtection` endpoint automatically.

# Example Walkthrough

In this example, we are going to do the following.

1. READ: Fetch the cluster and output the current deletion protection status.
2. VERIFY: Assert the value matches the expected configuration.

If you check the `terraform.template.tfvars` file — copy it to `terraform.tfvars` and update the values with your organization credentials.

## READ
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan

No changes. Your infrastructure matches the configuration.

Changes to Outputs:
  + deletion_protection = true
```

### Apply the Plan

Command: `terraform apply`

Sample Output:
```
$ terraform apply

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

deletion_protection = true
```

### Verification failure

If the cluster's deletion protection does not match the expected value:

```
│ Warning: Check block assertion failed
│
│   on deletion_protection.tf line 13, in check "deletion_protection_matches":
│   13:     condition     = data.couchbase-capella_cluster.existing_cluster.deletion_protection == var.deletion_protection
│
│ Cluster deletion_protection is false, expected true.
```

