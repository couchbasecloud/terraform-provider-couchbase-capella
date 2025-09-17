# Capella App Endpoint Default OIDC Provider Example

This example shows how to set the default OpenID Connect (OIDC) provider for a Capella App Endpoint.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we will do the following:

1. CREATE: Set the default OIDC provider for an App Endpoint as defined in `set_default_oidc_provider.tf`.
2. UPDATE: Update the default provider by changing `provider_id`.
3. IMPORT: Import the default OIDC provider state into Terraform.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the values as per your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Example output:

```
 $ terraform plan

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

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```
 $ terraform plan
### Apply the Plan, in order to set the default OIDC provider

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

## UPDATE
### Edit `terraform.tfvars` to change `provider_id` and apply

Command: `terraform apply`
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
╷
│ Warning: Value for undeclared variable
│ 
│ The root module does not declare a variable named "access_control_function" but a value was found in file "terraform.tfvars". If you meant to use this value, add a "variable" block to the configuration.
│ 
│ To silence these warnings, use TF_VAR_... environment variables to provide certain "global" settings to all configurations in your organization. To reduce the verbosity of these warnings, use the -compact-warnings option.
╵

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_default_oidc_provider.default: Modifying...
couchbase-capella_app_endpoint_default_oidc_provider.default: Modifications complete after 7s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

```

## IMPORT
### Remove the resource from the Terraform State file

Command: `terraform state rm couchbase-capella_app_endpoint_default_oidc_provider.example_default`
Example output:

```
 $ terraform state rm couchbase-capella_app_endpoint_default_oidc_provider.default              
Removed couchbase-capella_app_endpoint_default_oidc_provider.default
Successfully removed 1 resource instance(s).
```
### Now, import the resource into Terraform

Command: `terraform import couchbase-capella_app_endpoint_default_oidc_provider.example_default organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>,app_endpoint_name=<app_endpoint_name>`

Here, we pass the identifiers as a single comma-separated string. The current default provider is derived by the provider during Read.

Example output:

```
 $ terraform import  couchbase-capella_app_endpoint_default_oidc_provider.default app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000                                           1 ↵
couchbase-capella_app_endpoint_default_oidc_provider.default: Importing from ID "app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint_default_oidc_provider.default: Import prepared!
  Prepared couchbase-capella_app_endpoint_default_oidc_provider for import
couchbase-capella_app_endpoint_default_oidc_provider.default: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

## DELETE
This resource represents a selection; the underlying API does not support unsetting the default provider. Deleting this resource only removes it from Terraform state.
