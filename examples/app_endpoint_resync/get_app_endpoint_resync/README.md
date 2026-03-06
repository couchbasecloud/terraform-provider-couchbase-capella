# Read App Endpoint Resync

This example shows how to read an App Endpoint Resync from Capella. It uses the organization ID, project ID, cluster ID, App Service ID, and App Endpoint name to do so. 

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
data.couchbase-capella_app_endpoint_resync.app_endpoint_resync: Reading...
data.couchbase-capella_app_endpoint_resync.app_endpoint_resync: Read complete after 2s

Changes to Outputs:
  + app_endpoint_resync = {
      + app_endpoint_name      = "test_app_endpoint"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = {
          + inventory = [
              + "airline",
              + "airport",
              + "hotel",
              + "landmark",
              + "route",
            ]
        }
      + docs_changed           = 0
      + docs_processed         = 0
      + last_error             = ""
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + start_time             = "2026-03-06T12:28:24Z"
      + state                  = "completed"
    }
```

### Apply the Plan

Command: `terraform apply`

```
data.couchbase-capella_app_endpoint_resync.app_endpoint_resync: Reading...
data.couchbase-capella_app_endpoint_resync.app_endpoint_resync: Read complete after 1s

Changes to Outputs:
  + app_endpoint_resync = {
      + app_endpoint_name      = "test_app_endpoint"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = {
          + inventory = [
              + "airline",
              + "airport",
              + "hotel",
              + "landmark",
              + "route",
            ]
        }
      + docs_changed           = 0
      + docs_processed         = 0
      + last_error             = ""
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + start_time             = "2026-03-06T12:28:24Z"
      + state                  = "completed"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

app_endpoint_resync = {
  "app_endpoint_name" = "test_app_endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections_processing" = tomap({
    "inventory" = toset([
      "airline",
      "airport",
      "hotel",
      "landmark",
      "route",
    ])
  })
  "docs_changed" = 0
  "docs_processed" = 0
  "last_error" = ""
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "start_time" = "2026-03-06T12:28:24Z"
  "state" = "completed"
}
```