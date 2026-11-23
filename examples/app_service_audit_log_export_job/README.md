# Capella App Service Audit Log Export Job Example

This example shows how to export the audit logs of an App Service in Capella for a specific time period.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new App Service audit log export job, as stated in the `create_app_service_audit_log_export_job.tf` file.
2. UPDATE: Not supported. A change to any argument destroys the export job and creates a new one.
3. IMPORT: Not supported. The API does not return `start` and `end`, so an imported export job would be replaced on the next apply.
4. DELETE: Not supported by the API. Destroy only removes the export job from the Terraform state.

Copy the `terraform.template.tfvars` file to `terraform.tfvars`, and update the values of the variables for your organization.
The `end` time must be at least 15 minutes in the past.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job will be created
  + resource "couchbase-capella_app_service_audit_log_export_job" "new_app_service_audit_log_export_job" {
      + app_service_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit            = (known after apply)
      + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + download_expires = (known after apply)
      + download_id      = (known after apply)
      + end              = "2026-11-02T19:13:34Z"
      + id               = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + start            = "2026-11-02T17:33:34Z"
      + status           = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service_audit_log_export_job = {
      + app_service_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit            = (known after apply)
      + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + download_expires = (known after apply)
      + download_id      = (known after apply)
      + end              = "2026-11-02T19:13:34Z"
      + id               = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + start            = "2026-11-02T17:33:34Z"
      + status           = (known after apply)
    }

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the plan to create the export job

Command: `terraform apply`

Sample Output:
```
$ terraform apply
Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job will be created
  + resource "couchbase-capella_app_service_audit_log_export_job" "new_app_service_audit_log_export_job" {
      + app_service_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit            = (known after apply)
      + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + download_expires = (known after apply)
      + download_id      = (known after apply)
      + end              = "2026-11-02T19:13:34Z"
      + id               = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + start            = "2026-11-02T17:33:34Z"
      + status           = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service_audit_log_export_job = {
      + app_service_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit            = (known after apply)
      + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + download_expires = (known after apply)
      + download_id      = (known after apply)
      + end              = "2026-11-02T19:13:34Z"
      + id               = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + start            = "2026-11-02T17:33:34Z"
      + status           = (known after apply)
    }
couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job: Creating...
couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job: Creation complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_app_service_audit_log_export_job = {
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "audit" = {
    "created_at" = "2026-11-02 19:33:34.900431506 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2026-11-02 19:33:35.151515683 +0000 UTC"
    "modified_by" = "apikey-ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 3
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "download_expires" = tostring(null)
  "download_id" = ""
  "end" = "2026-11-02T19:13:34Z"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "start" = "2026-11-02T17:33:34Z"
  "status" = "No Log Files In Requested Time Range"
}
```

The `status` and the download fields change while the export runs. Run `terraform apply --refresh-only` to update them in the state.

### Run a plan again to confirm that there are no changes

Command: `terraform plan`

Sample Output:
```
$ terraform plan
couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

## DELETE
### Remove the export job from the Terraform state

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job will be destroyed
  - resource "couchbase-capella_app_service_audit_log_export_job" "new_app_service_audit_log_export_job" {
      - app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - audit           = {
          - created_at  = "2026-11-02 19:33:34.900431506 +0000 UTC" -> null
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
          - modified_at = "2026-11-02 19:33:35.151515683 +0000 UTC" -> null
          - modified_by = "ffffffff-aaaa-1414-eeee-000000000000" -> null
          - version     = 3 -> null
        } -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - download_id     = "" -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - status          = "No Log Files In Requested Time Range" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_app_service_audit_log_export_job = {
      - app_service_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      - audit            = {
          - created_at  = "2026-11-02 19:33:34.900431506 +0000 UTC"
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          - modified_at = "2026-11-02 19:33:35.151515683 +0000 UTC"
          - modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
          - version     = 3
        }
      - cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      - download_expires = null
      - download_id      = ""
      - end              = null
      - id               = "ffffffff-aaaa-1414-eeee-000000000000"
      - organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      - start            = null
      - status           = "No Log Files In Requested Time Range"
    } -> null
couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
