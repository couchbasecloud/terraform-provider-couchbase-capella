# Get Existing App Endpoint Audit Logging Config Data source

This example shows how to get the audit logging config of an App Endpoint in Capella. It uses the organization ID, project ID, cluster ID, App Service ID, and App Endpoint name to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_app_endpoint_audit_logging_config.existing_app_endpoint_audit_logging_config: Reading...
data.couchbase-capella_app_endpoint_audit_logging_config.existing_app_endpoint_audit_logging_config: Read complete after 1s

Changes to Outputs:
  + existing_app_endpoint_audit_logging_config = {
      + app_endpoint_name = "my-app-endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit_enabled     = false
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_roles    = []
      + disabled_users    = []
      + enabled_event_ids = [
          + 53280,
          + 53281,
          + 53282,
          + 53283,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

### Apply the Plan

Command: `terraform apply`

Sample Output:

```
$ terraform apply
data.couchbase-capella_app_endpoint_audit_logging_config.existing_app_endpoint_audit_logging_config: Reading...
data.couchbase-capella_app_endpoint_audit_logging_config.existing_app_endpoint_audit_logging_config: Read complete after 0s

Changes to Outputs:
  + existing_app_endpoint_audit_logging_config = {
      + app_endpoint_name = "my-app-endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit_enabled     = false
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_roles    = []
      + disabled_users    = []
      + enabled_event_ids = [
          + 53280,
          + 53281,
          + 53282,
          + 53283,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_app_endpoint_audit_logging_config = {
  "app_endpoint_name" = "my-app-endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "audit_enabled" = false
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_roles" = toset([])
  "disabled_users" = toset([])
  "enabled_event_ids" = toset([
    53280,
    53281,
    53282,
    53283,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```
