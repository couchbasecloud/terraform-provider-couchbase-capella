
# Capella App Endpoint Default OIDC Provider Example

This example shows how to set and manage the default OIDC provider for an App Endpoint in Capella.

It uses the organization ID, project ID, cluster ID, app service ID, app endpoint name, and the `provider_id` to configure the default provider.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following:

1. CREATE: Set the default OIDC provider as defined in `create_app_endpoint_default_oidc_provider.tf`.
2. READ: Fetch current default OIDC provider.
3. UPDATE: Change the default OIDC provider.
4. IMPORT: Import an existing default OIDC provider configuration into Terraform state.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the variable values for your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Example output:
```
$ terraform plan                                                                                                                                                                                                                                               1 ↵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_default_oidc_provider.default will be created
  + resource "couchbase-capella_app_endpoint_default_oidc_provider" "default" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_id       = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the plan to set the default OIDC provider

Command: `terraform apply`

Example output:
``` 
$ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_default_oidc_provider.default will be created
  + resource "couchbase-capella_app_endpoint_default_oidc_provider" "default" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_id       = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_default_oidc_provider.default: Creating...
couchbase-capella_app_endpoint_default_oidc_provider.default: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

```
## READ
### Show current state

Command: `terraform show`

Example output:
```
$ terraform show 
# couchbase-capella_app_endpoint_default_oidc_provider.default:
resource "couchbase-capella_app_endpoint_default_oidc_provider" "default" {
    app_endpoint_name = "api"
    app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
    cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
    project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    provider_id       = "ffffffff-aaaa-1414-eeee-000000000000"
}
```
## UPDATE
### Change the default provider ID

Edit `default_provider_id` in `terraform.tfvars` or update the resource block, then run:

Command: `terraform apply`

Example output:
```
$ terraform apply

couchbase-capella_app_endpoint_default_oidc_provider.default: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_default_oidc_provider.default will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_default_oidc_provider" "default" {
      ~ provider_id       = "ffffffff-aaaa-1414-eeee-000000000000" -> "ffffffff-aaaa-1414-eeee-000000000000"
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_default_oidc_provider.default: Modifying...
couchbase-capella_app_endpoint_default_oidc_provider.default: Modifications complete after 8s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```
## IMPORT
### Import an existing default OIDC provider configuration into Terraform state

Ensure the resource block has the correct IDs configured, then run:

Command: `terraform import couchbase-capella_app_endpoint_default_oidc_provider.example organization_id=<org_id>,project_id=<proj_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>,app_endpoint_name=<app_endpoint_name>`

Example output:
```
 $ terraform import  couchbase-capella_app_endpoint_default_oidc_provider.default app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_app_endpoint_default_oidc_provider.default: Importing from ID "app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint_default_oidc_provider.default: Import prepared!
  Prepared couchbase-capella_app_endpoint_default_oidc_provider for import
couchbase-capella_app_endpoint_default_oidc_provider.default: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```
