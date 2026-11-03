# Get Existing App Service Audit Log Streaming State Data source

This example shows how to get the audit log streaming state and output type of an App Service in Capella. It uses the organization ID, project ID, cluster ID, and App Service ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_app_service_audit_log_streaming_state.existing_app_service_audit_log_streaming_state: Reading...
data.couchbase-capella_app_service_audit_log_streaming_state.existing_app_service_audit_log_streaming_state: Read complete after 1s

Changes to Outputs:
  + existing_app_service_audit_log_streaming_state = {
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_app_endpoints = []
      + log_streaming_state    = "disabled"
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + output_type            = null
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + streaming_enabled      = false
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

### Apply the Plan

Command: `terraform apply`

Sample Output:

```
$ terraform apply
data.couchbase-capella_app_service_audit_log_streaming_state.existing_app_service_audit_log_streaming_state: Reading...
data.couchbase-capella_app_service_audit_log_streaming_state.existing_app_service_audit_log_streaming_state: Read complete after 0s

Changes to Outputs:
  + existing_app_service_audit_log_streaming_state = {
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_app_endpoints = []
      + log_streaming_state    = "disabled"
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + output_type            = null
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + streaming_enabled      = false
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_app_service_audit_log_streaming_state = {
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_app_endpoints" = toset([])
  "log_streaming_state" = "disabled"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "output_type" = tostring(null)
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "streaming_enabled" = false
}
```
