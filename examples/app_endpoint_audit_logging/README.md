# Capella App Endpoint Audit Logging Resource Example

This example shows how to manage the audit logging configuration of a Capella App Endpoint.

The resource uses the organization ID, project ID, cluster ID, App Service ID and App Endpoint name. When you create or update the resource, Terraform sets the audit logging configuration of the App Endpoint. There is no delete API. When you destroy the resource, Terraform removes it from the state file only, and the configuration in Capella stays in its last state.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Set the audit logging configuration of the App Endpoint as stated in the `create_app_endpoint_audit_logging.tf` file.
2. UPDATE: Change the enabled event IDs and the disabled users.
3. IMPORT: Import the App Endpoint audit logging configuration into Terraform.
4. DELETE: Remove the resource from the Terraform state file.

A change to an ID attribute replaces the resource. A change to `audit_enabled`, `enabled_event_ids`, `disabled_users` or `disabled_roles` updates the resource in place.

Copy the `terraform.template.tfvars` file to `terraform.tfvars` and update the values of the variables for your organization.

## CREATE

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging will be created
  + resource "couchbase-capella_app_endpoint_audit_logging" "new_app_endpoint_audit_logging" {
      + app_endpoint_name = "my-app-endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_roles    = (known after apply)
      + disabled_users    = [
          + {
              + domain = "sgw"
              + name   = "audit-user"
            },
        ]
      + enabled_event_ids = [
          + 53280,
          + 53281,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the plan to set the audit logging configuration

Command: `terraform apply`

Sample Output:
```
$ terraform apply
...
couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging: Creating...
couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging: Creation complete after 3s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## UPDATE

Change `enabled_event_ids` to `[53280, 53290, 54000]` and run `terraform apply`.

Sample Output:
```
$ terraform apply
...
  # couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_audit_logging" "new_app_endpoint_audit_logging" {
      ~ enabled_event_ids = [
          - 53281,
          + 53290,
          + 54000,
            # (1 unchanged element hidden)
        ]
        # (8 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
...
Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## IMPORT

### Remove the resource from the Terraform state file

Command: `terraform state rm couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging`

### Import the resource into Terraform

Command: `terraform import couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging app_endpoint_name=<app_endpoint_name>,app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
$ terraform import couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging app_endpoint_name=my-app-endpoint,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging: Importing from ID "app_endpoint_name=my-app-endpoint,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging: Import prepared!
  Prepared couchbase-capella_app_endpoint_audit_logging for import
couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging: Refreshing state...

Import successful!
```

The import ID is a comma-separated string of the App Endpoint name, App Service ID, cluster ID, project ID and organization ID.

## DESTROY

### Remove the resource from the Terraform state file

Command: `terraform destroy`

This command makes no API call. The audit logging configuration in Capella stays in its last state.

Sample Output:
```
$ terraform destroy
...
couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging: Destroying...
couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
