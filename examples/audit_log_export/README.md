# Capella Audit Log Export Example

This example shows how to create and audit log export job in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1.  CREATE: Create a new audit log export job in Capella per the `create_audit_log_exports.tf` file.
2.  LIST: List existing audit log export jobs in Capella as stated in the `list_audit_log_exports.tf` file.
3.  UPDATE: Show it is not supported
4.  DELETE: Show it is not supported

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_audit_log_export.new_auditlogexport will be created
  + resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      + auditlog_download_url = (known after apply)
      + cluster_id            = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      + project_id            = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
      + start                 = "2024-03-13T04:44:15+00:00"
      + status                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_auditlogexport = {
      + cluster_id      = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      + data            = null
      + organization_id = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      + project_id      = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
    }
  + new_auditlogexport      = {
      + auditlog_download_url = (known after apply)
      + cluster_id            = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      + project_id            = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
      + start                 = "2024-03-13T04:44:15+00:00"
      + status                = (known after apply)
    }

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply"
now.
```

### Apply the Plan, in order to create an audit log export job

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_audit_log_export.new_auditlogexport will be created
  + resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      + auditlog_download_url = (known after apply)
      + cluster_id            = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      + project_id            = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
      + start                 = "2024-03-13T04:44:15+00:00"
      + status                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_auditlogexport = {
      + cluster_id      = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      + data            = null
      + organization_id = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      + project_id      = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
    }
  + new_auditlogexport      = {
      + auditlog_download_url = (known after apply)
      + cluster_id            = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      + project_id            = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
      + start                 = "2024-03-13T04:44:15+00:00"
      + status                = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_audit_log_export.new_auditlogexport: Creating...
couchbase-capella_audit_log_export.new_auditlogexport: Creation complete after 0s [id=b173a42e-3c2c-4245-af81-f8615d63b5ce]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

existing_auditlogexport = {
  "cluster_id" = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
  "data" = toset(null) /* of object */
  "organization_id" = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
  "project_id" = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
}
new_auditlogexport = {
  "auditlog_download_url" = tostring(null)
  "cluster_id" = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
  "created_at" = "2024-03-13 06:48:43.596764031 +0000 UTC"
  "end" = "2024-03-13T06:44:15+00:00"
  "expiration" = tostring(null)
  "id" = "b173a42e-3c2c-4245-af81-f8615d63b5ce"
  "organization_id" = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
  "project_id" = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
  "start" = "2024-03-13T04:44:15+00:00"
  "status" = tostring(null)
}
```

### View the audit log export job

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=b173a42e-3c2c-4245-af81-f8615d63b5ce]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_audit_log_export.new_auditlogexport has changed
  ~ resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      + auditlog_download_url = "https://cb-audit-logs-637cea1d-fce5-40a5-9d48-8d0690e656ee.s3.amazonaws.com/export/cluster-audit-logs-41938b91-66ed-4c84-b70c-55c3f0ae4266-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVIKDGOEO3%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T064846Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEE8aCXVzLWVhc3QtMSJIMEYCIQDM2n%2BpV6cnfeTrSnvjlcl15s%2FoyMcYLVfljG4BsTYyFAIhAO93hNcpnJUKqtiKj%2BFxF7p5amM1IFDCvyki2UYdCgLMKqACCFgQAxoMMjY0MTM4NDY4Mzk0IgxvjWc3jbG5idbrLWQq%2FQGvHYupgMUmqlFHUvNH1OKIn9sZby26SKtlKiV57%2B8%2BiYZDsztbnB%2FgwZKcLYVWwq7zGvGCqeZEk1%2FgUitGJ3Vj0tiSjlTmwnQCYAwkaZqMMUH4m%2BkM66S8jfMtKIA1RAwTTJ1T2%2FoYVBDe%2FsbORrYiPyN4DfhnKRZySeTSy1cXaYedVzjiA%2BYugUwUB2GwDbaofWofE6Ho2%2FEP5dkfVWs0iHwxVFlYSaq7aJFCl%2BKj8hXhlmujNVeMW%2Bd6T6MHcwSMf2%2B429Pnhm%2BYc00eVOQxOEyvjyMA4V%2FuYDOv5coTkTg5y4m0EHVYve8tsro67UpdMzrumwbFoUOrw7WUMMyYxa8GOpwB8SvDoGuX3OERsbyY%2FIelecnQtQW4wvdGBUdZ2OcmS19dMQjmCuKBO4aKvGpWYpeuUouj9zgYcgUzC39xBMUwtrAgGRn0horrc7tGJD2CHbR1UOIbu7pJSw%2FZOvNTw%2Bg7w%2F8tfgMg9f7Iobnfsi2xYbREquAGShEWi64tC0kf5rCOFFacM%2BZm3tq3x8sbT%2FJ87djItmfJk%2FYf2y9W&X-Amz-SignedHeaders=host&X-Amz-Signature=7eb266db82ee7af58b42d8b6c2daa2ceb4f435b07c503fc839edff460b327d4e"
      + expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC"
        id                    = "b173a42e-3c2c-4245-af81-f8615d63b5ce"
        # (6 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may
include actions to undo or respond to these changes.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Changes to Outputs:
  ~ existing_auditlogexport = {
      ~ data            = null -> [
          + {
              + auditlog_download_url = "https://cb-audit-logs-637cea1d-fce5-40a5-9d48-8d0690e656ee.s3.amazonaws.com/export/cluster-audit-logs-41938b91-66ed-4c84-b70c-55c3f0ae4266-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVIKDGOEO3%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T064846Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEE8aCXVzLWVhc3QtMSJIMEYCIQDM2n%2BpV6cnfeTrSnvjlcl15s%2FoyMcYLVfljG4BsTYyFAIhAO93hNcpnJUKqtiKj%2BFxF7p5amM1IFDCvyki2UYdCgLMKqACCFgQAxoMMjY0MTM4NDY4Mzk0IgxvjWc3jbG5idbrLWQq%2FQGvHYupgMUmqlFHUvNH1OKIn9sZby26SKtlKiV57%2B8%2BiYZDsztbnB%2FgwZKcLYVWwq7zGvGCqeZEk1%2FgUitGJ3Vj0tiSjlTmwnQCYAwkaZqMMUH4m%2BkM66S8jfMtKIA1RAwTTJ1T2%2FoYVBDe%2FsbORrYiPyN4DfhnKRZySeTSy1cXaYedVzjiA%2BYugUwUB2GwDbaofWofE6Ho2%2FEP5dkfVWs0iHwxVFlYSaq7aJFCl%2BKj8hXhlmujNVeMW%2Bd6T6MHcwSMf2%2B429Pnhm%2BYc00eVOQxOEyvjyMA4V%2FuYDOv5coTkTg5y4m0EHVYve8tsro67UpdMzrumwbFoUOrw7WUMMyYxa8GOpwB8SvDoGuX3OERsbyY%2FIelecnQtQW4wvdGBUdZ2OcmS19dMQjmCuKBO4aKvGpWYpeuUouj9zgYcgUzC39xBMUwtrAgGRn0horrc7tGJD2CHbR1UOIbu7pJSw%2FZOvNTw%2Bg7w%2F8tfgMg9f7Iobnfsi2xYbREquAGShEWi64tC0kf5rCOFFacM%2BZm3tq3x8sbT%2FJ87djItmfJk%2FYf2y9W&X-Amz-SignedHeaders=host&X-Amz-Signature=7eb266db82ee7af58b42d8b6c2daa2ceb4f435b07c503fc839edff460b327d4e"
              + cluster_id            = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
              + created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC"
              + end                   = "2024-03-13 06:44:15 +0000 UTC"
              + expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC"
              + id                    = "b173a42e-3c2c-4245-af81-f8615d63b5ce"
              + organization_id       = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
              + project_id            = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
              + start                 = "2024-03-13 04:44:15 +0000 UTC"
              + status                = "Completed"
            },
        ]
        # (3 unchanged attributes hidden)
    }
  ~ new_auditlogexport      = {
      ~ auditlog_download_url = null -> "https://cb-audit-logs-637cea1d-fce5-40a5-9d48-8d0690e656ee.s3.amazonaws.com/export/cluster-audit-logs-41938b91-66ed-4c84-b70c-55c3f0ae4266-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVIKDGOEO3%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T064846Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEE8aCXVzLWVhc3QtMSJIMEYCIQDM2n%2BpV6cnfeTrSnvjlcl15s%2FoyMcYLVfljG4BsTYyFAIhAO93hNcpnJUKqtiKj%2BFxF7p5amM1IFDCvyki2UYdCgLMKqACCFgQAxoMMjY0MTM4NDY4Mzk0IgxvjWc3jbG5idbrLWQq%2FQGvHYupgMUmqlFHUvNH1OKIn9sZby26SKtlKiV57%2B8%2BiYZDsztbnB%2FgwZKcLYVWwq7zGvGCqeZEk1%2FgUitGJ3Vj0tiSjlTmwnQCYAwkaZqMMUH4m%2BkM66S8jfMtKIA1RAwTTJ1T2%2FoYVBDe%2FsbORrYiPyN4DfhnKRZySeTSy1cXaYedVzjiA%2BYugUwUB2GwDbaofWofE6Ho2%2FEP5dkfVWs0iHwxVFlYSaq7aJFCl%2BKj8hXhlmujNVeMW%2Bd6T6MHcwSMf2%2B429Pnhm%2BYc00eVOQxOEyvjyMA4V%2FuYDOv5coTkTg5y4m0EHVYve8tsro67UpdMzrumwbFoUOrw7WUMMyYxa8GOpwB8SvDoGuX3OERsbyY%2FIelecnQtQW4wvdGBUdZ2OcmS19dMQjmCuKBO4aKvGpWYpeuUouj9zgYcgUzC39xBMUwtrAgGRn0horrc7tGJD2CHbR1UOIbu7pJSw%2FZOvNTw%2Bg7w%2F8tfgMg9f7Iobnfsi2xYbREquAGShEWi64tC0kf5rCOFFacM%2BZm3tq3x8sbT%2FJ87djItmfJk%2FYf2y9W&X-Amz-SignedHeaders=host&X-Amz-Signature=7eb266db82ee7af58b42d8b6c2daa2ceb4f435b07c503fc839edff460b327d4e"
      ~ expiration            = null -> "2024-03-16 06:48:46.64587117 +0000 UTC"
        id                    = "b173a42e-3c2c-4245-af81-f8615d63b5ce"
        # (7 unchanged attributes hidden)
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply"
now.
```

### Update is not supported
### In this example we try to update the start time and it fails

Command: `terraform apply`

Sample output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=b173a42e-3c2c-4245-af81-f8615d63b5ce]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_audit_log_export.new_auditlogexport will be updated in-place
  ~ resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      ~ auditlog_download_url = "https://cb-audit-logs-637cea1d-fce5-40a5-9d48-8d0690e656ee.s3.amazonaws.com/export/cluster-audit-logs-41938b91-66ed-4c84-b70c-55c3f0ae4266-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVIKDGOEO3%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T064846Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEE8aCXVzLWVhc3QtMSJIMEYCIQDM2n%2BpV6cnfeTrSnvjlcl15s%2FoyMcYLVfljG4BsTYyFAIhAO93hNcpnJUKqtiKj%2BFxF7p5amM1IFDCvyki2UYdCgLMKqACCFgQAxoMMjY0MTM4NDY4Mzk0IgxvjWc3jbG5idbrLWQq%2FQGvHYupgMUmqlFHUvNH1OKIn9sZby26SKtlKiV57%2B8%2BiYZDsztbnB%2FgwZKcLYVWwq7zGvGCqeZEk1%2FgUitGJ3Vj0tiSjlTmwnQCYAwkaZqMMUH4m%2BkM66S8jfMtKIA1RAwTTJ1T2%2FoYVBDe%2FsbORrYiPyN4DfhnKRZySeTSy1cXaYedVzjiA%2BYugUwUB2GwDbaofWofE6Ho2%2FEP5dkfVWs0iHwxVFlYSaq7aJFCl%2BKj8hXhlmujNVeMW%2Bd6T6MHcwSMf2%2B429Pnhm%2BYc00eVOQxOEyvjyMA4V%2FuYDOv5coTkTg5y4m0EHVYve8tsro67UpdMzrumwbFoUOrw7WUMMyYxa8GOpwB8SvDoGuX3OERsbyY%2FIelecnQtQW4wvdGBUdZ2OcmS19dMQjmCuKBO4aKvGpWYpeuUouj9zgYcgUzC39xBMUwtrAgGRn0horrc7tGJD2CHbR1UOIbu7pJSw%2FZOvNTw%2Bg7w%2F8tfgMg9f7Iobnfsi2xYbREquAGShEWi64tC0kf5rCOFFacM%2BZm3tq3x8sbT%2FJ87djItmfJk%2FYf2y9W&X-Amz-SignedHeaders=host&X-Amz-Signature=7eb266db82ee7af58b42d8b6c2daa2ceb4f435b07c503fc839edff460b327d4e" -> (known after apply)
      ~ created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC" -> (known after apply)
      ~ expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC" -> (known after apply)
        id                    = "b173a42e-3c2c-4245-af81-f8615d63b5ce"
      ~ start                 = "2024-03-13T04:44:15+00:00" -> "2024-03-13T02:44:15+00:00"
      + status                = (known after apply)
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_auditlogexport      = {
      ~ auditlog_download_url = "https://cb-audit-logs-637cea1d-fce5-40a5-9d48-8d0690e656ee.s3.amazonaws.com/export/cluster-audit-logs-41938b91-66ed-4c84-b70c-55c3f0ae4266-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVIKDGOEO3%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T064846Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEE8aCXVzLWVhc3QtMSJIMEYCIQDM2n%2BpV6cnfeTrSnvjlcl15s%2FoyMcYLVfljG4BsTYyFAIhAO93hNcpnJUKqtiKj%2BFxF7p5amM1IFDCvyki2UYdCgLMKqACCFgQAxoMMjY0MTM4NDY4Mzk0IgxvjWc3jbG5idbrLWQq%2FQGvHYupgMUmqlFHUvNH1OKIn9sZby26SKtlKiV57%2B8%2BiYZDsztbnB%2FgwZKcLYVWwq7zGvGCqeZEk1%2FgUitGJ3Vj0tiSjlTmwnQCYAwkaZqMMUH4m%2BkM66S8jfMtKIA1RAwTTJ1T2%2FoYVBDe%2FsbORrYiPyN4DfhnKRZySeTSy1cXaYedVzjiA%2BYugUwUB2GwDbaofWofE6Ho2%2FEP5dkfVWs0iHwxVFlYSaq7aJFCl%2BKj8hXhlmujNVeMW%2Bd6T6MHcwSMf2%2B429Pnhm%2BYc00eVOQxOEyvjyMA4V%2FuYDOv5coTkTg5y4m0EHVYve8tsro67UpdMzrumwbFoUOrw7WUMMyYxa8GOpwB8SvDoGuX3OERsbyY%2FIelecnQtQW4wvdGBUdZ2OcmS19dMQjmCuKBO4aKvGpWYpeuUouj9zgYcgUzC39xBMUwtrAgGRn0horrc7tGJD2CHbR1UOIbu7pJSw%2FZOvNTw%2Bg7w%2F8tfgMg9f7Iobnfsi2xYbREquAGShEWi64tC0kf5rCOFFacM%2BZm3tq3x8sbT%2FJ87djItmfJk%2FYf2y9W&X-Amz-SignedHeaders=host&X-Amz-Signature=7eb266db82ee7af58b42d8b6c2daa2ceb4f435b07c503fc839edff460b327d4e" -> (known after apply)
      ~ created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC" -> (known after apply)
      ~ expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC" -> (known after apply)
        id                    = "b173a42e-3c2c-4245-af81-f8615d63b5ce"
      ~ start                 = "2024-03-13T04:44:15+00:00" -> "2024-03-13T02:44:15+00:00"
      + status                = (known after apply)
        # (4 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_audit_log_export.new_auditlogexport: Modifying... [id=b173a42e-3c2c-4245-af81-f8615d63b5ce]
╷
│ Error: Audit Log Export does not support update
│
│   with couchbase-capella_audit_log_export.new_auditlogexport,
│   on create_auditlog_exports.tf line 6, in resource "couchbase-capella_audit_log_export" "new_auditlogexport":
│    6: resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
│
│ Audit Log Export does not support update
╵
```

###  Delete is not supported

Command: `terraform destroy`

Sample output:
```
terraform destroy
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=b173a42e-3c2c-4245-af81-f8615d63b5ce]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_audit_log_export.new_auditlogexport will be destroyed
  - resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      - auditlog_download_url = "https://cb-audit-logs-637cea1d-fce5-40a5-9d48-8d0690e656ee.s3.amazonaws.com/export/cluster-audit-logs-41938b91-66ed-4c84-b70c-55c3f0ae4266-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVIKDGOEO3%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T064846Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEE8aCXVzLWVhc3QtMSJIMEYCIQDM2n%2BpV6cnfeTrSnvjlcl15s%2FoyMcYLVfljG4BsTYyFAIhAO93hNcpnJUKqtiKj%2BFxF7p5amM1IFDCvyki2UYdCgLMKqACCFgQAxoMMjY0MTM4NDY4Mzk0IgxvjWc3jbG5idbrLWQq%2FQGvHYupgMUmqlFHUvNH1OKIn9sZby26SKtlKiV57%2B8%2BiYZDsztbnB%2FgwZKcLYVWwq7zGvGCqeZEk1%2FgUitGJ3Vj0tiSjlTmwnQCYAwkaZqMMUH4m%2BkM66S8jfMtKIA1RAwTTJ1T2%2FoYVBDe%2FsbORrYiPyN4DfhnKRZySeTSy1cXaYedVzjiA%2BYugUwUB2GwDbaofWofE6Ho2%2FEP5dkfVWs0iHwxVFlYSaq7aJFCl%2BKj8hXhlmujNVeMW%2Bd6T6MHcwSMf2%2B429Pnhm%2BYc00eVOQxOEyvjyMA4V%2FuYDOv5coTkTg5y4m0EHVYve8tsro67UpdMzrumwbFoUOrw7WUMMyYxa8GOpwB8SvDoGuX3OERsbyY%2FIelecnQtQW4wvdGBUdZ2OcmS19dMQjmCuKBO4aKvGpWYpeuUouj9zgYcgUzC39xBMUwtrAgGRn0horrc7tGJD2CHbR1UOIbu7pJSw%2FZOvNTw%2Bg7w%2F8tfgMg9f7Iobnfsi2xYbREquAGShEWi64tC0kf5rCOFFacM%2BZm3tq3x8sbT%2FJ87djItmfJk%2FYf2y9W&X-Amz-SignedHeaders=host&X-Amz-Signature=7eb266db82ee7af58b42d8b6c2daa2ceb4f435b07c503fc839edff460b327d4e" -> null
      - cluster_id            = "41938b91-66ed-4c84-b70c-55c3f0ae4266" -> null
      - created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC" -> null
      - end                   = "2024-03-13T06:44:15+00:00" -> null
      - expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC" -> null
      - id                    = "b173a42e-3c2c-4245-af81-f8615d63b5ce" -> null
      - organization_id       = "637cea1d-fce5-40a5-9d48-8d0690e656ee" -> null
      - project_id            = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90" -> null
      - start                 = "2024-03-13T04:44:15+00:00" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - existing_auditlogexport = {
      - cluster_id      = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
      - data            = [
          - {
              - auditlog_download_url = "https://cb-audit-logs-637cea1d-fce5-40a5-9d48-8d0690e656ee.s3.amazonaws.com/export/cluster-audit-logs-41938b91-66ed-4c84-b70c-55c3f0ae4266-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVIKDGOEO3%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T064846Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEE8aCXVzLWVhc3QtMSJIMEYCIQDM2n%2BpV6cnfeTrSnvjlcl15s%2FoyMcYLVfljG4BsTYyFAIhAO93hNcpnJUKqtiKj%2BFxF7p5amM1IFDCvyki2UYdCgLMKqACCFgQAxoMMjY0MTM4NDY4Mzk0IgxvjWc3jbG5idbrLWQq%2FQGvHYupgMUmqlFHUvNH1OKIn9sZby26SKtlKiV57%2B8%2BiYZDsztbnB%2FgwZKcLYVWwq7zGvGCqeZEk1%2FgUitGJ3Vj0tiSjlTmwnQCYAwkaZqMMUH4m%2BkM66S8jfMtKIA1RAwTTJ1T2%2FoYVBDe%2FsbORrYiPyN4DfhnKRZySeTSy1cXaYedVzjiA%2BYugUwUB2GwDbaofWofE6Ho2%2FEP5dkfVWs0iHwxVFlYSaq7aJFCl%2BKj8hXhlmujNVeMW%2Bd6T6MHcwSMf2%2B429Pnhm%2BYc00eVOQxOEyvjyMA4V%2FuYDOv5coTkTg5y4m0EHVYve8tsro67UpdMzrumwbFoUOrw7WUMMyYxa8GOpwB8SvDoGuX3OERsbyY%2FIelecnQtQW4wvdGBUdZ2OcmS19dMQjmCuKBO4aKvGpWYpeuUouj9zgYcgUzC39xBMUwtrAgGRn0horrc7tGJD2CHbR1UOIbu7pJSw%2FZOvNTw%2Bg7w%2F8tfgMg9f7Iobnfsi2xYbREquAGShEWi64tC0kf5rCOFFacM%2BZm3tq3x8sbT%2FJ87djItmfJk%2FYf2y9W&X-Amz-SignedHeaders=host&X-Amz-Signature=7eb266db82ee7af58b42d8b6c2daa2ceb4f435b07c503fc839edff460b327d4e"
              - cluster_id            = "41938b91-66ed-4c84-b70c-55c3f0ae4266"
              - created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC"
              - end                   = "2024-03-13 06:44:15 +0000 UTC"
              - expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC"
              - id                    = "b173a42e-3c2c-4245-af81-f8615d63b5ce"
              - organization_id       = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
              - project_id            = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
              - start                 = "2024-03-13 04:44:15 +0000 UTC"
              - status                = "Completed"
            },
        ]
      - organization_id = "637cea1d-fce5-40a5-9d48-8d0690e656ee"
      - project_id      = "6aa13dbf-69cb-48e3-97af-d89f57ea7f90"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_audit_log_export.new_auditlogexport: Destroying... [id=b173a42e-3c2c-4245-af81-f8615d63b5ce]
╷
│ Error: Audit Log Export does not support delete
│
│ Audit Log Export does not support delete
```