# Get Existing App Service Audit Log Export Job Data source

This example shows how to get an audit log export job of an App Service in Capella. It uses the organization ID, project ID, cluster ID, App Service ID, and export job ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_app_service_audit_log_export_job.existing_app_service_audit_log_export_job: Reading...
data.couchbase-capella_app_service_audit_log_export_job.existing_app_service_audit_log_export_job: Read complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Changes to Outputs:
  + existing_app_service_audit_log_export_job = {
      + app_service_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + audit            = {
          + created_at  = "2026-11-02 19:29:27.925083022 +0000 UTC"
          + created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          + modified_at = "2026-11-02 19:29:28.17842374 +0000 UTC"
          + modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
          + version     = 3
        }
      + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + download_expires = null
      + download_id      = ""
      + id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + status           = "No Log Files In Requested Time Range"
      + tenant_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```
