# Get Existing App Service Log Streaming Data source

This example shows how to get the Log Streaming configuration and state that already exists in Capella for a given App Service. It uses the organization ID, project ID, cluster ID, and App Service ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_app_service_log_streaming.existing_app_service_log_streaming: Reading...
data.couchbase-capella_app_service_log_streaming.existing_app_service_log_streaming: Read complete after 1s

Changes to Outputs:
  + existing_app_service_log_streaming = {
      + app_service_id  = "8554d8b8-029f-449e-a644-215dd9a02129"
      + cluster_id      = "1777a520-94b7-4b99-9022-4b38996a370d"
      + config_state    = "enabled"
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + output_type     = "generic_http"
      + project_id      = "7013dfa2-3dd8-436b-adf9-e2580a406dd0"
      + streaming_state = "unknown"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan

Command: `terraform apply`

Sample Output:

```
$ terraform apply
data.couchbase-capella_app_service_log_streaming.existing_app_service_log_streaming: Reading...
data.couchbase-capella_app_service_log_streaming.existing_app_service_log_streaming: Read complete after 0s

Changes to Outputs:
  + existing_app_service_log_streaming = {
      + app_service_id  = "8554d8b8-029f-449e-a644-215dd9a02129"
      + cluster_id      = "1777a520-94b7-4b99-9022-4b38996a370d"
      + config_state    = "enabled"
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + output_type     = "generic_http"
      + project_id      = "7013dfa2-3dd8-436b-adf9-e2580a406dd0"
      + streaming_state = "unknown"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_app_service_log_streaming = {
  "app_service_id" = "8554d8b8-029f-449e-a644-215dd9a02129"
  "cluster_id" = "1777a520-94b7-4b99-9022-4b38996a370d"
  "config_state" = "enabled"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "output_type" = "generic_http"
  "project_id" = "7013dfa2-3dd8-436b-adf9-e2580a406dd0"
  "streaming_state" = "unhealthy"
}
```
