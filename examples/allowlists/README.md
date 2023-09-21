# Capella AllowList Example

This example shows how to create and manage Projects in Capella.

This creates a new allowlist in the selected Capella cluster. It uses the organization ID, projectId and clusterId to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

Running the example

For planning phase

```
terraform plan
```

For apply phase

```
terraform apply
```

Alternatively to passing variable inputs with each command, the `terraform.template.tfvars` file can be copied to `terraform.tfvars` and updated.

To remove the allowlist

```
 terraform destroy
```
