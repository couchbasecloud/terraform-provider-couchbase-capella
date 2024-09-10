# Capella Event Example

This example shows how to retrieve event for Capella.

This fetch the event details based on the event ID and authentication access token.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the event details as stated in the `get_event.tf` file.
2. DELETE: Delete the event data output from terraform state.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## GET
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan

╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-capella
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_event.existing_event: Reading...
data.couchbase-capella_event.existing_event: Read complete after 2s [id=ffffffff-aaaa-1414-eeee-000000000001]

Changes to Outputs:
  + existing_event = {
      + alert_key        = "cluster_deployment_completed"
      + app_service_id   = null
      + app_service_name = null
      + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_name     = "searchBoxColumnarInstance-"
      + id               = "ffffffff-aaaa-1414-eeee-000000000001"
      + image_url        = null
      + incident_ids     = []
      + key              = "cluster_deployment_completed"
      + kv               = "null"
      + occurrence_count = null
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_name     = "!!!!!!!-Shared-Project-!!!!!!!"
      + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + session_id       = null
      + severity         = "info"
      + source           = "cp-jobs"
      + summary          = null
      + timestamp        = "2024-07-08 17:38:42.240367422 +0000 UTC"
      + user_email       = null
      + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + user_name        = "kevin"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to get the events.

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-capella
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_event.existing_event: Reading...
data.couchbase-capella_event.existing_event: Read complete after 3s [id=ffffffff-aaaa-1414-eeee-000000000001]

Changes to Outputs:
  + existing_event = {
      + alert_key        = "cluster_deployment_completed"
      + app_service_id   = null
      + app_service_name = null
      + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_name     = "searchBoxColumnarInstance-"
      + id               = "ffffffff-aaaa-1414-eeee-000000000001"
      + image_url        = null
      + incident_ids     = []
      + key              = "cluster_deployment_completed"
      + kv               = "null"
      + occurrence_count = null
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_name     = "!!!!!!!-Shared-Project-!!!!!!!"
      + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + session_id       = null
      + severity         = "info"
      + source           = "cp-jobs"
      + summary          = null
      + timestamp        = "2024-07-08 17:38:42.240367422 +0000 UTC"
      + user_email       = null
      + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + user_name        = "kevin"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_event = {
  "alert_key" = "cluster_deployment_completed"
  "app_service_id" = tostring(null)
  "app_service_name" = tostring(null)
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_name" = "searchBoxColumnarInstance-"
  "id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "image_url" = tostring(null)
  "incident_ids" = toset([])
  "key" = "cluster_deployment_completed"
  "kv" = "null"
  "occurrence_count" = tonumber(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_name" = "!!!!!!!-Shared-Project-!!!!!!!"
  "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "session_id" = tostring(null)
  "severity" = "info"
  "source" = "cp-jobs"
  "summary" = tostring(null)
  "timestamp" = "2024-07-08 17:38:42.240367422 +0000 UTC"
  "user_email" = tostring(null)
  "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "user_name" = "kevin"
}
```

## DELETE
### Finally, delete the state of the event from terraform outputs

Command: `terraform destroy`

Sample Output
```
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-capella
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_event.existing_event: Reading...
data.couchbase-capella_event.existing_event: Read complete after 2s [id=ffffffff-aaaa-1414-eeee-000000000001]

Changes to Outputs:
  - existing_event = {
      - alert_key        = "cluster_deployment_completed"
      - app_service_id   = null
      - app_service_name = null
      - cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      - cluster_name     = "searchBoxColumnarInstance-"
      - id               = "ffffffff-aaaa-1414-eeee-000000000001"
      - image_url        = null
      - incident_ids     = []
      - key              = "cluster_deployment_completed"
      - kv               = "null"
      - occurrence_count = null
      - organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_name     = "!!!!!!!-Shared-Project-!!!!!!!"
      - request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      - session_id       = null
      - severity         = "info"
      - source           = "cp-jobs"
      - summary          = null
      - timestamp        = "2024-07-08 17:38:42.240367422 +0000 UTC"
      - user_email       = null
      - user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      - user_name        = "kevin"
    } -> null

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes


Destroy complete! Resources: 0 destroyed.
```