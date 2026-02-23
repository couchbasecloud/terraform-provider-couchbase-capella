# Read App Endpoint Logging Config

This example shows how to read an App Endpoint Logging Config from Capella. It uses the organization ID, project ID, cluster ID, App Service ID, and App Endpoint name to do so. 

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
data.couchbase-capella_app_endpoint_log_streaming_config.app_endpoint_log_streaming_config: Reading...
data.couchbase-capella_app_endpoint_log_streaming_config.app_endpoint_log_streaming_config: Read complete after 2s

Changes to Outputs:
  + app_endpoint_log_streaming_config = {
      + app_endpoint_name = "test_app_endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + log_keys          = [
          + "CRUD",
          + "Cache",
          + "Changes",
          + "HTTP",
          + "HTTP+",
          + "Query",
        ]
      + log_level         = "info"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }
```

### Apply the Plan

Command: `terraform apply`

```
data.couchbase-capella_app_endpoint_log_streaming_config.app_endpoint_log_streaming_config: Reading...
data.couchbase-capella_app_endpoint_log_streaming_config.app_endpoint_log_streaming_config: Read complete after 0s

Changes to Outputs:
  + app_endpoint_log_streaming_config = {
      + app_endpoint_name = "test_app_endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + log_keys          = [
          + "CRUD",
          + "Cache",
          + "Changes",
          + "HTTP",
          + "HTTP+",
          + "Query",
        ]
      + log_level         = "info"
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

app_endpoint_log_streaming_config = {
  "app_endpoint_name" = "test_app_endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "log_keys" = toset([
    "CRUD",
    "Cache",
    "Changes",
    "HTTP",
    "HTTP+",
    "Query",
  ])
  "log_level" = "info"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```