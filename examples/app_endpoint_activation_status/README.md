# Capella App Endpoint Activation Status Example

This example shows how to activate (bring online) or deactivate (take offline) an App Endpoint on demand in Capella.

It uses the organization ID, project ID, cluster ID, app service ID, and app endpoint name to toggle the activation status. The resource maps to:
- POST: bring an App Endpoint online
- DELETE: take an App Endpoint offline
- GET: read the App Endpoint and derive its current state

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Set the activation status for an App Endpoint as defined in the `create_app_endpoint_activation_status.tf` file.
2. UPDATE: Update the current App Endpoint status to online/offline.
3. IMPORT: Import an existing App Endpoint activation status into Terraform state.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Example output:
```
 $ terraform plan   


Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_activation_status.example_activation_status will be created
  + resource "couchbase-capella_app_endpoint_activation_status" "example_activation_status" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + online            = true
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.



──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```


### Apply the Plan, in order to set the App Endpoint online/offline

Command: `terraform apply`

Example output:

```
 $ terraform apply


Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_activation_status.example_activation_status will be created
  + resource "couchbase-capella_app_endpoint_activation_status" "example_activation_status" {
      + app_endpoint_name = "api"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + online            = true
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.



Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_activation_status.example_activation_status: Creating...
couchbase-capella_app_endpoint_activation_status.example_activation_status: Creation complete after 3s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## IMPORT
### Remove the resource `app_endpoint_activation` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_endpoint_activation_status.example_activation_status`

Example output:

```
$ terraform state rm couchbase-capella_app_endpoint_activation_status.example_activation_status
Removed couchbase-capella_app_endpoint_activation_status.example_activation_status
Successfully removed 1 resource instance(s).
```

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_endpoint_activation_status.app_endpoint_activation app_endpoint_name=<app_endpoint_name>,app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Here, we pass the identifiers as a single comma-separated string.
The first identifier is the `app_endpoint_name` (the App Endpoint name, not an ID).
The second is the `app_service_id` (the ID of the App Service hosting the endpoint).
The third is the `cluster_id` (the ID of the cluster the App Service belongs to).
The fourth is the `project_id` (the ID of the project the cluster belongs to).
The fifth is the `organization_id` (the ID of the organization the project belongs to).

Notes:
- `app_endpoint_name` is the App Endpoint name (not an ID).
- `online` is derived during Read by querying the App Endpoint.

Example output:

```
 $ terraform import couchbase-capella_app_endpoint_activation_status.example_activation_status organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,app_endpoint_name=api

couchbase-capella_app_endpoint_activation_status.example_activation_status: Importing from ID "organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,app_endpoint_name=api"...
couchbase-capella_app_endpoint_activation_status.example_activation_status: Import prepared!
  Prepared couchbase-capella_app_endpoint_activation_status for import
couchbase-capella_app_endpoint_activation_status.example_activation_status: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

## UPDATE
### Let us edit terraform.tfvars file to change the App Endpoint status

Command: `terraform apply`
Change the `online` value to `false` to turn off the App Endpoint

Example output:

```
 $ terraform apply                                                                              

couchbase-capella_app_endpoint_activation_status.example_activation_status: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_activation_status.example_activation_status will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_activation_status" "example_activation_status" {
      ~ online            = true -> false
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.


Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_activation_status.example_activation_status: Modifying...
couchbase-capella_app_endpoint_activation_status.example_activation_status: Modifications complete after 1s
```