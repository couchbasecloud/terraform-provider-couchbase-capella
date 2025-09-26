# Capella App Endpoint OIDC Provider Example

This example shows how to create and manage an OpenID Connect (OIDC) provider for a Capella App Endpoint.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we will do the following:

1. CREATE: Create a new OIDC provider for an App Endpoint as defined in `create_app_endpoint_oidc_provider.tf`.
2. UPDATE: Update the OIDC provider configuration using Terraform.
3. IMPORT: Import an existing OIDC provider into Terraform state.
4. DELETE: Delete the OIDC provider.

If you check the `terraform.template.tfvars` file, copy it to `terraform.tfvars` and update the values as per your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`
Example output:

```
 $ terraform plan 



Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_oidc_provider.oidc_provider will be created
  + resource "couchbase-capella_app_endpoint_oidc_provider" "oidc_provider" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + client_id         = "tf-client-id2"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + is_default        = (known after apply)
      + issuer            = "https://accounts.google.com"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_id       = (known after apply)
      + register          = true
      + username_claim    = "abc,23asd"
    }

Plan: 1 to add, 0 to change, 0 to destroy.



──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```
### Apply the Plan, in order to create an OIDC provider

Command: `terraform apply`

```
 $ terraform apply 



Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_oidc_provider.oidc_provider will be created
  + resource "couchbase-capella_app_endpoint_oidc_provider" "oidc_provider" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + client_id         = "tf-client-id2"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + is_default        = (known after apply)
      + issuer            = "https://accounts.google.com"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_id       = (known after apply)
      + register          = true
      + username_claim    = "abc,23asd"
    }

Plan: 1 to add, 0 to change, 0 to destroy.



Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Creating...
couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

```

## IMPORT
### Remove the resource from the Terraform State file

First remove the resource from the state file. This does not delete the resource in Capella, just removes it from Terraform's state management.

Command: `terraform state rm couchbase-capella_app_endpoint_oidc_provider.example_oidc_provider`
Example output:
```
$ terraform state rm couchbase-capella_app_endpoint_oidc_provider.oidc_provider 
Removed couchbase-capella_app_endpoint_oidc_provider.oidc_provider
Successfully removed 1 resource instance(s).
```

Then we import the resource back into the state file using the `terraform import` command.

Command: `terraform import `
```
 $ terraform import  couchbase-capella_app_endpoint_oidc_provider.imported_oidc_provider app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,provider_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_app_endpoint_oidc_provider.imported_oidc_provider: Importing from ID "app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,provider_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint_oidc_provider.imported_oidc_provider: Import prepared!
  Prepared couchbase-capella_app_endpoint_oidc_provider for import
couchbase-capella_app_endpoint_oidc_provider.imported_oidc_provider: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the identifiers as a single comma-separated string, with the last value being the `provider_id` of the OIDC provider.

## UPDATE
### Edit `terraform.tfvars` to update provider settings and apply
In this example, we will update the `client_id` and `username_claim` values.
Command: `terraform apply`

Example output:
```
$ terraform apply

couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_oidc_provider.oidc_provider will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_oidc_provider" "oidc_provider" {
      ~ client_id         = "tf-client-id3" -> "tf-client-id22"
      ~ is_default        = false -> (known after apply)
      ~ provider_id       = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
      ~ username_claim    = "abc,23asd" -> "user123"
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.


Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Modifying...
couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```
## DELETE
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Example output:
```
$ terraform destroy                

couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_oidc_provider.oidc_provider will be destroyed
  - resource "couchbase-capella_app_endpoint_oidc_provider" "oidc_provider" {
      - app_endpoint_name = "api" -> null
      - app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - client_id         = "tf-client-id22" -> null
      - cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - is_default        = false -> null
      - issuer            = "https://accounts.google.com" -> null
      - organization_id   = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - provider_id       = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - register          = true -> null
      - username_claim    = "user123" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Destroying...
couchbase-capella_app_endpoint_oidc_provider.oidc_provider: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.

```