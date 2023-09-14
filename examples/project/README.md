# Capella Projects Example

This example shows how to create and manage Projects in Capella.

This creates a new project in the selected Capella organization and lists existing Projects in the organization. It uses the organization ID to create and list Projects.

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

Once the Project is created, you can use the Project ID to launch a cluster in it.

To remove the Project

```
 terraform destroy
```