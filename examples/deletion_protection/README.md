# Deletion Protection Example

This example demonstrates how to use the `deletion_protection` attribute on a Couchbase Capella cluster.

## What it does

- Creates a cluster with `deletion_protection = true`
- When protection is enabled, `terraform destroy` will fail with:
  ```
  Error: Cluster deletion protection is enabled. Set deletion_protection = false before destroying.
  ```

## How to disable protection for destruction

1. Set `deletion_protection = false` in your configuration
2. Run `terraform apply` to update the protection setting
3. Run `terraform destroy` to delete the cluster

## Usage

```bash
cp terraform.template.tfvars terraform.tfvars
# Edit terraform.tfvars with your values
terraform init
terraform plan
terraform apply
```

## Toggling protection

To disable protection on an existing cluster, change the attribute and apply:

```hcl
resource "couchbase-capella_cluster" "protected_cluster" {
  # ...
  deletion_protection = false
}
```

```bash
terraform apply
```

