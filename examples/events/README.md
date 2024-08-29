# Capella Events Example

This example shows how to retrieve events for Capella.

This lists the event details based on the project ID, cluster ID, user ID, severity levels, tags, from, to, page, perPage, sortBy, sortDirection and authentication access token.

Currently, only tags can have multiple values; all other multivalued filters are included for future-proofing.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the events details as stated in the `list_events.tf` file.
2. DELETE: Delete the events data output from terraform state.

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
data.couchbase-capella_events.existing_events: Reading...
data.couchbase-capella_events.existing_events: Read complete after 1s

Changes to Outputs:
  + existing_events = {
      + cluster_ids     = null
      + cursor          = {
          + hrefs = {
              + first    = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=2>"
              + last     = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=12&perPage=2>"
              + next     = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=2>"
              + previous = ""
            }
          + pages = {
              + last        = 12
              + next        = 2
              + page        = 1
              + per_page    = 2
              + previous    = 0
              + total_items = 24
            }
        }
      + data            = [
          + {
              + alert_key        = "cluster_deployment_requested"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "searchBoxColumnarInstance-"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deployment_requested"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "!!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-07-08 17:37:07.116412925 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "kevin"
            },
          + {
              + alert_key        = "cluster_deployment_completed"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "searchBoxColumnarInstance-"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deployment_completed"
              + kv               = "null"
              + occurrence_count = null
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
            },
        ]
      + from            = "2024-07-07T04:19:25Z"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = 1
      + per_page        = 2
      + project_ids     = [
          + "ffffffff-aaaa-1414-eeee-000000000000",
        ]
      + severity_levels = [
          + "info",
        ]
      + sort_by         = "timestamp"
      + sort_direction  = "asc"
      + tags            = [
          + "availability",
        ]
      + to              = "2024-07-30T04:19:25Z"
      + user_ids        = [
          + "ffffffff-aaaa-1414-eeee-000000000000",
        ]
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

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
data.couchbase-capella_events.existing_events: Reading...
data.couchbase-capella_events.existing_events: Read complete after 2s

Changes to Outputs:
  + existing_events = {
      + cluster_ids     = null
      + cursor          = {
          + hrefs = {
              + first    = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=2>"
              + last     = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=12&perPage=2>"
              + next     = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=2>"
              + previous = ""
            }
          + pages = {
              + last        = 12
              + next        = 2
              + page        = 1
              + per_page    = 2
              + previous    = 0
              + total_items = 24
            }
        }
      + data            = [
          + {
              + alert_key        = "cluster_deployment_requested"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "searchBoxColumnarInstance-"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deployment_requested"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "!!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-07-08 17:37:07.116412925 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "kevin"
            },
          + {
              + alert_key        = "cluster_deployment_completed"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "searchBoxColumnarInstance-"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deployment_completed"
              + kv               = "null"
              + occurrence_count = null
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
            },
        ]
      + from            = "2024-07-07T04:19:25Z"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = 1
      + per_page        = 2
      + project_ids     = [
          + "ffffffff-aaaa-1414-eeee-000000000000",
        ]
      + severity_levels = [
          + "info",
        ]
      + sort_by         = "timestamp"
      + sort_direction  = "asc"
      + tags            = [
          + "availability",
        ]
      + to              = "2024-07-30T04:19:25Z"
      + user_ids        = [
          + "ffffffff-aaaa-1414-eeee-000000000000",
        ]
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_events = {
  "cluster_ids" = toset(null) /* of string */
  "cursor" = {
    "hrefs" = {
      "first" = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=2>"
      "last" = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=12&perPage=2>"
      "next" = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=2>"
      "previous" = ""
    }
    "pages" = {
      "last" = 12
      "next" = 2
      "page" = 1
      "per_page" = 2
      "previous" = 0
      "total_items" = 24
    }
  }
  "data" = tolist([
    {
      "alert_key" = "cluster_deployment_requested"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "searchBoxColumnarInstance-"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "cluster_deployment_requested"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "!!!!!!!-Shared-Project-!!!!!!!"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "severity" = "info"
      "source" = "cp-api"
      "summary" = tostring(null)
      "timestamp" = "2024-07-08 17:37:07.116412925 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "kevin"
    },
    {
      "alert_key" = "cluster_deployment_completed"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "searchBoxColumnarInstance-"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "cluster_deployment_completed"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "!!!!!!!-Shared-Project-!!!!!!!"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = null
      "severity" = "info"
      "source" = "cp-jobs"
      "summary" = tostring(null)
      "timestamp" = "2024-07-08 17:38:42.240367422 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "kevin"
    },
  ])
  "from" = "2024-07-07T04:19:25Z"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "page" = 1
  "per_page" = 2
  "project_ids" = toset([
    "ffffffff-aaaa-1414-eeee-000000000000",
  ])
  "severity_levels" = toset([
    "info",
  ])
  "sort_by" = "timestamp"
  "sort_direction" = "asc"
  "tags" = toset([
    "availability",
  ])
  "to" = "2024-07-30T04:19:25Z"
  "user_ids" = toset([
    "ffffffff-aaaa-1414-eeee-000000000000",
  ])
}
```

## DELETE
### Finally, delete the state of the events from terraform outputs

Command: `terraform destroy`

Sample Output:
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
data.couchbase-capella_events.existing_events: Reading...
data.couchbase-capella_events.existing_events: Read complete after 2s

Changes to Outputs:
  - existing_events = {
      - cluster_ids     = null
      - cursor          = {
          - hrefs = {
              - first    = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=2>"
              - last     = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=12&perPage=2>"
              - next     = "<https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=2>"
              - previous = ""
            }
          - pages = {
              - last        = 12
              - next        = 2
              - page        = 1
              - per_page    = 2
              - previous    = 0
              - total_items = 24
            }
        }
      - data            = [
          - {
              - alert_key        = "cluster_deployment_requested"
              - app_service_id   = null
              - app_service_name = null
              - cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              - cluster_name     = "searchBoxColumnarInstance-"
              - id               = "ffffffff-aaaa-1414-eeee-000000000000"
              - image_url        = null
              - incident_ids     = []
              - key              = "cluster_deployment_requested"
              - kv               = "null"
              - occurrence_count = null
              - project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              - project_name     = "!!!!!!!-Shared-Project-!!!!!!!"
              - request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              - session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              - severity         = "info"
              - source           = "cp-api"
              - summary          = null
              - timestamp        = "2024-07-08 17:37:07.116412925 +0000 UTC"
              - user_email       = null
              - user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              - user_name        = "kevin"
            },
          - {
              - alert_key        = "cluster_deployment_completed"
              - app_service_id   = null
              - app_service_name = null
              - cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              - cluster_name     = "searchBoxColumnarInstance-"
              - id               = "ffffffff-aaaa-1414-eeee-000000000000"
              - image_url        = null
              - incident_ids     = []
              - key              = "cluster_deployment_completed"
              - kv               = "null"
              - occurrence_count = null
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
            },
        ]
      - from            = "2024-07-07T04:19:25Z"
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - page            = 1
      - per_page        = 2
      - project_ids     = [
          - "ffffffff-aaaa-1414-eeee-000000000000",
        ]
      - severity_levels = [
          - "info",
        ]
      - sort_by         = "timestamp"
      - sort_direction  = "asc"
      - tags            = [
          - "availability",
        ]
      - to              = "2024-07-30T04:19:25Z"
      - user_ids        = [
          - "ffffffff-aaaa-1414-eeee-000000000000",
        ]
    } -> null

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

Destroy complete! Resources: 0 destroyed.
```
