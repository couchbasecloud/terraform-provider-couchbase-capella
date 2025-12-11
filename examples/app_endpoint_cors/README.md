# Capella App Endpoint CORS Example

This example shows how to configure CORS (Cross-Origin Resource Sharing) for an App Endpoint in Capella.

It uses the organization ID, project ID, cluster ID, app service ID, and app endpoint name to manage CORS settings. You can set allowed `origin`, optional `login_origin`, `headers`, `max_age` and `disabled` values.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following:

1. CREATE: Create a new CORS configuration as defined in `create_app_endpoint_cors.tf`.
2. UPDATE: Update the CORS configuration using Terraform.
3. IMPORT: Import an existing CORS configuration into Terraform state.

If you check the `terraform.template.tfvars` file — copy it to `terraform.tfvars` and update the variable values for your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Example output:

```
 $ terraform plan   


Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_cors.example1 will be created
  + resource "couchbase-capella_app_endpoint_cors" "example1" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled          = (known after apply)
      + max_age           = (known after apply)
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + origin            = [
          + "*",
        ]
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the plan to create the CORS configuration

Command: `terraform apply`

Example output:

```
 $ terraform apply


Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_cors.example1 will be created
  + resource "couchbase-capella_app_endpoint_cors" "example1" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled          = (known after apply)
      + max_age           = (known after apply)
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + origin            = [
          + "*",
        ]
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_cors.example1: Creating...
couchbase-capella_app_endpoint_cors.example1: Creation complete after 2s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## UPDATE
### Update the CORS configuration

Modify attributes in `terraform.tfvars` or the resource block, then run:

Command: `terraform apply`

Example output:
```
 $ terraform apply

couchbase-capella_app_endpoint_cors.example1: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_cors.example1 will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_cors" "example1" {
      ~ disabled          = false -> (known after apply)
      + headers           = [
          + "Accept-Encoding",
        ]
      ~ max_age           = 0 -> (known after apply)
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_cors.example1: Modifying...
couchbase-capella_app_endpoint_cors.example1: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## IMPORT
### Import an existing CORS configuration into Terraform state

Ensure the resource block has the correct IDs configured, then:

Command: `terraform import couchbase-capella_app_endpoint_cors.example organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,app_endpoint_name=<name>`

Example output:
```
 $ terraform import  couchbase-capella_app_endpoint_cors.cors app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000

couchbase-capella_app_endpoint_cors.cors: Importing from ID "app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint_cors.cors: Import prepared!
  Prepared couchbase-capella_app_endpoint_cors for import
couchbase-capella_app_endpoint_cors.cors: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

