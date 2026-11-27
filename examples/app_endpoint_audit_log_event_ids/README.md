# Get App Endpoint Audit Log Event IDs Data source

This example shows how to get the audit log event IDs of an App Endpoint in Capella. It uses the organization ID, project ID, cluster ID, App Service ID, and App Endpoint name to do so.

Use these event IDs to configure audit logging for the App Endpoint.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_app_endpoint_audit_log_event_ids.existing_app_endpoint_audit_log_event_ids: Reading...
data.couchbase-capella_app_endpoint_audit_log_event_ids.existing_app_endpoint_audit_log_event_ids: Read complete after 0s

Changes to Outputs:
  + existing_app_endpoint_audit_log_event_ids = {
      + app_endpoint_name = "my-app-endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + data              = [
          + {
              + description = "A database was deleted"
              + enabled     = false
              + filterable  = true
              + id          = 54002
              + name        = "Delete database"
            },
          + {
              + description = "A document was deleted"
              + enabled     = false
              + filterable  = true
              + id          = 55003
              + name        = "Delete document"
            },
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
data.couchbase-capella_app_endpoint_audit_log_event_ids.existing_app_endpoint_audit_log_event_ids: Reading...
data.couchbase-capella_app_endpoint_audit_log_event_ids.existing_app_endpoint_audit_log_event_ids: Read complete after 0s

Changes to Outputs:
  + existing_app_endpoint_audit_log_event_ids = {
      + app_endpoint_name = "my-app-endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + data              = [
          + {
              + description = "A database was deleted"
              + enabled     = false
              + filterable  = true
              + id          = 54002
              + name        = "Delete database"
            },
          + {
              + description = "A document was deleted"
              + enabled     = false
              + filterable  = true
              + id          = 55003
              + name        = "Delete document"
            },
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

existing_app_endpoint_audit_log_event_ids = {
  "app_endpoint_name" = "my-app-endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = toset([
    {
      "description" = "A database was deleted"
      "enabled" = false
      "filterable" = true
      "id" = 54002
      "name" = "Delete database"
    },
    {
      "description" = "A document was deleted"
      "enabled" = false
      "filterable" = true
      "id" = 55003
      "name" = "Delete document"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```
