# Capella Audit Log Export Example

This example shows how to create and audit log export job in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1.  CREATE: Create a new audit log export job in Capella per the `create_audit_log_exports.tf` file.
2.  LIST: List existing audit log export jobs in Capella as stated in the `list_audit_log_exports.tf` file, and
          detecting if resource was updated outside terraform. 
3.  UPDATE: Show it is not supported
4.  DELETE: Show it is not supported
5.  IMPORT:  Import an existing audit log export job


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

## IMPORT
### Remove the resource from the Terraform State file

Command:  `terraform state rm couchbase-capella_audit_log_export.new_auditlogexport`

Sample Output:
```
terraform state rm couchbase-capella_audit_log_export.new_auditlogexport
Removed couchbase-capella_audit_log_export.new_auditlogexport
Successfully removed 1 resource instance(s).

```

### Now, let's import the resource in Terraform

Command:  `terraform import couchbase-capella_audit_log_export.new_auditlogexport id=4c7cb3f9-2ec8-4e58-add0-5f866357dcff,cluster_id=e17880b4-5124-43a5-a0e5-0735f771eae5,project_id=53fbca2e-56b7-4951-9912-7baa7396a607,organization_id=fabaa8f8-8571-4f8f-961f-6898c15b80d6`

Sample Output:
```
terraform import couchbase-capella_audit_log_export.new_auditlogexport id=4c7cb3f9-2ec8-4e58-add0-5f866357dcff,cluster_id=e17880b4-5124-43a5-a0e5-0735f771eae5,project_id=53fbca2e-56b7-4951-9912-7baa7396a607,organization_id=fabaa8f8-8571-4f8f-961f-6898c15b80d6
couchbase-capella_audit_log_export.new_auditlogexport: Importing from ID "id=4c7cb3f9-2ec8-4e58-add0-5f866357dcff,cluster_id=e17880b4-5124-43a5-a0e5-0735f771eae5,project_id=53fbca2e-56b7-4951-9912-7baa7396a607,organization_id=fabaa8f8-8571-4f8f-961f-6898c15b80d6"...
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Import prepared!
  Prepared couchbase-capella_audit_log_export for import
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=id=4c7cb3f9-2ec8-4e58-add0-5f866357dcff,cluster_id=e17880b4-5124-43a5-a0e5-0735f771eae5,project_id=53fbca2e-56b7-4951-9912-7baa7396a607,organization_id=fabaa8f8-8571-4f8f-961f-6898c15b80d6]

2024-03-13T12:42:14.621-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_audit_log_export.new_auditlogexport during refresh.
      - .auditlog_download_url: was null, but now cty.StringVal("https://cb-audit-logs-fabaa8f8-8571-4f8f-961f-6898c15b80d6.s3.amazonaws.com/export/cluster-audit-logs-e17880b4-5124-43a5-a0e5-0735f771eae5-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVAHECFHHJ%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T193829Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEFwaCXVzLWVhc3QtMSJHMEUCIQDzq0SS4TWaUJJPTTKwYtogDssCmwufSkoBApiQqHbxYQIgNukuV20pylKGx5kN0ZRftJGV4id3NF8RiPrfrMfKWLsqoAIIZRADGgwyNjQxMzg0NjgzOTQiDPdTLN4GCLNNyZpjnir9AYVC6%2FTJy%2B5Fwl3WI5C90l3cQxkFMd3vVOCXTUBn%2B16VfZimBx26L0R1Zn%2BPROueDRg54U%2FKABUszlLBV2wot0TEvXFNk0QFMXYjCSUl5K3Oj2%2BXSRwX6G6RUJ2%2FGksgl39fwyIP%2F35Jv8KKky%2F4LcKi4hdpJTFWCitI5La4AWBcFV6SRjJ3H%2FaDBEg5XnXgRXHU2XGCzZ6rQd37av5fPIeVMZqGitwYhrAltGqYctDyauteULrouZazafDDu915lx7g7fi%2BOZO4CjYBm820LaU4EJpnx42ETydUEOBD1wt5Bw2fRUCHZ%2Bk2mPFsjXJth%2BQlJxCapUCMhll5%2BCcwsYHIrwY6nQEtJ9%2BEL1WFShM3diIxuMiR6DVFkr5GPa%2F9x0bQ%2B6S4dEohSYagbauSLbgRweKVDu6prGqjcitS3l7IieSecJ1h7ignYGiXRHaQLNIb8PNZp88WiapJZQw0c6cj%2ByS%2FAYXkH2VubuOOynuRuZpiFNlrUADHLT3mP0kzMxPz8%2BeWfj%2ByX0DSJMU2ijimv4zjY4RZzxmUo5NQQDOTinxO&X-Amz-SignedHeaders=host&X-Amz-Signature=8dd866e9520af9e948ac0f1580f14a1ee18e93e6abc22bb3251e2136c16665af")
      - .created_at: was null, but now cty.StringVal("2024-03-13 19:38:24.355397616 +0000 UTC")
      - .end: was null, but now cty.StringVal("2024-03-13T19:19:39+00:00")
      - .id: was cty.StringVal("id=4c7cb3f9-2ec8-4e58-add0-5f866357dcff,cluster_id=e17880b4-5124-43a5-a0e5-0735f771eae5,project_id=53fbca2e-56b7-4951-9912-7baa7396a607,organization_id=fabaa8f8-8571-4f8f-961f-6898c15b80d6"), but now cty.StringVal("4c7cb3f9-2ec8-4e58-add0-5f866357dcff")
      - .project_id: was null, but now cty.StringVal("53fbca2e-56b7-4951-9912-7baa7396a607")
      - .cluster_id: was null, but now cty.StringVal("e17880b4-5124-43a5-a0e5-0735f771eae5")
      - .expiration: was null, but now cty.StringVal("2024-03-16 19:38:29.842400291 +0000 UTC")
      - .organization_id: was null, but now cty.StringVal("fabaa8f8-8571-4f8f-961f-6898c15b80d6")
      - .start: was null, but now cty.StringVal("2024-03-13T17:19:39+00:00")
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the export ID i.e. the ID of the audit log export job.
The second ID is the cluster ID i.e. the ID of the cluster to which the apikey belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## Detecting changes made outside terraform
### First list audit log export jobs
### In this example there was export job with id 3043ce42-4f75-4751-927d-4f8128356ab3 created
### Do not apply the update

Command:  `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=4c7cb3f9-2ec8-4e58-add0-5f866357dcff]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_audit_log_export.new_auditlogexport will be updated in-place
  ~ resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      ~ auditlog_download_url = "https://cb-audit-logs-fabaa8f8-8571-4f8f-961f-6898c15b80d6.s3.amazonaws.com/export/cluster-audit-logs-e17880b4-5124-43a5-a0e5-0735f771eae5-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVAHECFHHJ%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T193829Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEFwaCXVzLWVhc3QtMSJHMEUCIQDzq0SS4TWaUJJPTTKwYtogDssCmwufSkoBApiQqHbxYQIgNukuV20pylKGx5kN0ZRftJGV4id3NF8RiPrfrMfKWLsqoAIIZRADGgwyNjQxMzg0NjgzOTQiDPdTLN4GCLNNyZpjnir9AYVC6%2FTJy%2B5Fwl3WI5C90l3cQxkFMd3vVOCXTUBn%2B16VfZimBx26L0R1Zn%2BPROueDRg54U%2FKABUszlLBV2wot0TEvXFNk0QFMXYjCSUl5K3Oj2%2BXSRwX6G6RUJ2%2FGksgl39fwyIP%2F35Jv8KKky%2F4LcKi4hdpJTFWCitI5La4AWBcFV6SRjJ3H%2FaDBEg5XnXgRXHU2XGCzZ6rQd37av5fPIeVMZqGitwYhrAltGqYctDyauteULrouZazafDDu915lx7g7fi%2BOZO4CjYBm820LaU4EJpnx42ETydUEOBD1wt5Bw2fRUCHZ%2Bk2mPFsjXJth%2BQlJxCapUCMhll5%2BCcwsYHIrwY6nQEtJ9%2BEL1WFShM3diIxuMiR6DVFkr5GPa%2F9x0bQ%2B6S4dEohSYagbauSLbgRweKVDu6prGqjcitS3l7IieSecJ1h7ignYGiXRHaQLNIb8PNZp88WiapJZQw0c6cj%2ByS%2FAYXkH2VubuOOynuRuZpiFNlrUADHLT3mP0kzMxPz8%2BeWfj%2ByX0DSJMU2ijimv4zjY4RZzxmUo5NQQDOTinxO&X-Amz-SignedHeaders=host&X-Amz-Signature=8dd866e9520af9e948ac0f1580f14a1ee18e93e6abc22bb3251e2136c16665af" -> (known after apply)
      ~ created_at            = "2024-03-13 19:38:24.355397616 +0000 UTC" -> (known after apply)
      ~ end                   = "2024-03-13T19:19:39+00:00" -> "2024-03-13T06:44:15+00:00"
      ~ expiration            = "2024-03-16 19:38:29.842400291 +0000 UTC" -> (known after apply)
        id                    = "4c7cb3f9-2ec8-4e58-add0-5f866357dcff"
      ~ start                 = "2024-03-13T17:19:39+00:00" -> "2024-03-13T02:44:15+00:00"
      + status                = (known after apply)
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ existing_auditlogexport = {
      ~ data            = [
            # (3 unchanged elements hidden)
            {
                auditlog_download_url = null
                cluster_id            = "e17880b4-5124-43a5-a0e5-0735f771eae5"
                created_at            = "2024-03-13 19:28:38.869267928 +0000 UTC"
                end                   = "2024-03-13 06:44:15 +0000 UTC"
                expiration            = null
                id                    = "eb39994b-fe95-483f-87c3-41ac4d2ce6d6"
                organization_id       = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
                project_id            = "53fbca2e-56b7-4951-9912-7baa7396a607"
                start                 = "2024-03-13 02:44:15 +0000 UTC"
                status                = "no audit log files exist within the requested time frame"
            },
          + {
              + auditlog_download_url = null
              + cluster_id            = "e17880b4-5124-43a5-a0e5-0735f771eae5"
              + created_at            = "2024-03-13 19:52:37.937857032 +0000 UTC"
              + end                   = "2024-03-12 19:19:39 +0000 UTC"
              + expiration            = null
              + id                    = "3043ce42-4f75-4751-927d-4f8128356ab3"
              + organization_id       = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
              + project_id            = "53fbca2e-56b7-4951-9912-7baa7396a607"
              + start                 = "2024-03-12 17:19:39 +0000 UTC"
              + status                = "no audit log files exist within the requested time frame"
            },
        ]
        # (3 unchanged attributes hidden)
    }
  ~ new_auditlogexport      = {
      ~ auditlog_download_url = "https://cb-audit-logs-fabaa8f8-8571-4f8f-961f-6898c15b80d6.s3.amazonaws.com/export/cluster-audit-logs-e17880b4-5124-43a5-a0e5-0735f771eae5-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVAHECFHHJ%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T193829Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEFwaCXVzLWVhc3QtMSJHMEUCIQDzq0SS4TWaUJJPTTKwYtogDssCmwufSkoBApiQqHbxYQIgNukuV20pylKGx5kN0ZRftJGV4id3NF8RiPrfrMfKWLsqoAIIZRADGgwyNjQxMzg0NjgzOTQiDPdTLN4GCLNNyZpjnir9AYVC6%2FTJy%2B5Fwl3WI5C90l3cQxkFMd3vVOCXTUBn%2B16VfZimBx26L0R1Zn%2BPROueDRg54U%2FKABUszlLBV2wot0TEvXFNk0QFMXYjCSUl5K3Oj2%2BXSRwX6G6RUJ2%2FGksgl39fwyIP%2F35Jv8KKky%2F4LcKi4hdpJTFWCitI5La4AWBcFV6SRjJ3H%2FaDBEg5XnXgRXHU2XGCzZ6rQd37av5fPIeVMZqGitwYhrAltGqYctDyauteULrouZazafDDu915lx7g7fi%2BOZO4CjYBm820LaU4EJpnx42ETydUEOBD1wt5Bw2fRUCHZ%2Bk2mPFsjXJth%2BQlJxCapUCMhll5%2BCcwsYHIrwY6nQEtJ9%2BEL1WFShM3diIxuMiR6DVFkr5GPa%2F9x0bQ%2B6S4dEohSYagbauSLbgRweKVDu6prGqjcitS3l7IieSecJ1h7ignYGiXRHaQLNIb8PNZp88WiapJZQw0c6cj%2ByS%2FAYXkH2VubuOOynuRuZpiFNlrUADHLT3mP0kzMxPz8%2BeWfj%2ByX0DSJMU2ijimv4zjY4RZzxmUo5NQQDOTinxO&X-Amz-SignedHeaders=host&X-Amz-Signature=8dd866e9520af9e948ac0f1580f14a1ee18e93e6abc22bb3251e2136c16665af" -> (known after apply)
      ~ created_at            = "2024-03-13 19:38:24.355397616 +0000 UTC" -> (known after apply)
      ~ end                   = "2024-03-13T19:19:39+00:00" -> "2024-03-13T06:44:15+00:00"
      ~ expiration            = "2024-03-16 19:38:29.842400291 +0000 UTC" -> (known after apply)
        id                    = "4c7cb3f9-2ec8-4e58-add0-5f866357dcff"
      ~ start                 = "2024-03-13T17:19:39+00:00" -> "2024-03-13T02:44:15+00:00"
      + status                = (known after apply)
        # (3 unchanged attributes hidden)
    }
```

### Refresh state

Command:  `terraform apply --refresh-only`

Sample Output:
```
terraform apply --refresh-only
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=4c7cb3f9-2ec8-4e58-add0-5f866357dcff]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Changes to Outputs:
  ~ existing_auditlogexport = {
      ~ data            = [
            # (3 unchanged elements hidden)
            {
                auditlog_download_url = null
                cluster_id            = "e17880b4-5124-43a5-a0e5-0735f771eae5"
                created_at            = "2024-03-13 19:28:38.869267928 +0000 UTC"
                end                   = "2024-03-13 06:44:15 +0000 UTC"
                expiration            = null
                id                    = "eb39994b-fe95-483f-87c3-41ac4d2ce6d6"
                organization_id       = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
                project_id            = "53fbca2e-56b7-4951-9912-7baa7396a607"
                start                 = "2024-03-13 02:44:15 +0000 UTC"
                status                = "no audit log files exist within the requested time frame"
            },
          + {
              + auditlog_download_url = null
              + cluster_id            = "e17880b4-5124-43a5-a0e5-0735f771eae5"
              + created_at            = "2024-03-13 19:52:37.937857032 +0000 UTC"
              + end                   = "2024-03-12 19:19:39 +0000 UTC"
              + expiration            = null
              + id                    = "3043ce42-4f75-4751-927d-4f8128356ab3"
              + organization_id       = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
              + project_id            = "53fbca2e-56b7-4951-9912-7baa7396a607"
              + start                 = "2024-03-12 17:19:39 +0000 UTC"
              + status                = "no audit log files exist within the requested time frame"
            },
        ]
        # (3 unchanged attributes hidden)
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Would you like to update the Terraform state to reflect these detected changes?
  Terraform will write these changes to the state without modifying any real infrastructure.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_auditlogexport = {
  "cluster_id" = "e17880b4-5124-43a5-a0e5-0735f771eae5"
  "data" = toset([
    {
      "auditlog_download_url" = "https://cb-audit-logs-fabaa8f8-8571-4f8f-961f-6898c15b80d6.s3.amazonaws.com/export/cluster-audit-logs-e17880b4-5124-43a5-a0e5-0735f771eae5-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVAHECFHHJ%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T193829Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEFwaCXVzLWVhc3QtMSJHMEUCIQDzq0SS4TWaUJJPTTKwYtogDssCmwufSkoBApiQqHbxYQIgNukuV20pylKGx5kN0ZRftJGV4id3NF8RiPrfrMfKWLsqoAIIZRADGgwyNjQxMzg0NjgzOTQiDPdTLN4GCLNNyZpjnir9AYVC6%2FTJy%2B5Fwl3WI5C90l3cQxkFMd3vVOCXTUBn%2B16VfZimBx26L0R1Zn%2BPROueDRg54U%2FKABUszlLBV2wot0TEvXFNk0QFMXYjCSUl5K3Oj2%2BXSRwX6G6RUJ2%2FGksgl39fwyIP%2F35Jv8KKky%2F4LcKi4hdpJTFWCitI5La4AWBcFV6SRjJ3H%2FaDBEg5XnXgRXHU2XGCzZ6rQd37av5fPIeVMZqGitwYhrAltGqYctDyauteULrouZazafDDu915lx7g7fi%2BOZO4CjYBm820LaU4EJpnx42ETydUEOBD1wt5Bw2fRUCHZ%2Bk2mPFsjXJth%2BQlJxCapUCMhll5%2BCcwsYHIrwY6nQEtJ9%2BEL1WFShM3diIxuMiR6DVFkr5GPa%2F9x0bQ%2B6S4dEohSYagbauSLbgRweKVDu6prGqjcitS3l7IieSecJ1h7ignYGiXRHaQLNIb8PNZp88WiapJZQw0c6cj%2ByS%2FAYXkH2VubuOOynuRuZpiFNlrUADHLT3mP0kzMxPz8%2BeWfj%2ByX0DSJMU2ijimv4zjY4RZzxmUo5NQQDOTinxO&X-Amz-SignedHeaders=host&X-Amz-Signature=8dd866e9520af9e948ac0f1580f14a1ee18e93e6abc22bb3251e2136c16665af"
      "cluster_id" = "e17880b4-5124-43a5-a0e5-0735f771eae5"
      "created_at" = "2024-03-13 19:38:24.355397616 +0000 UTC"
      "end" = "2024-03-13 19:19:39 +0000 UTC"
      "expiration" = "2024-03-16 19:38:29.842400291 +0000 UTC"
      "id" = "4c7cb3f9-2ec8-4e58-add0-5f866357dcff"
      "organization_id" = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
      "project_id" = "53fbca2e-56b7-4951-9912-7baa7396a607"
      "start" = "2024-03-13 17:19:39 +0000 UTC"
      "status" = "Completed"
    },
    {
      "auditlog_download_url" = tostring(null)
      "cluster_id" = "e17880b4-5124-43a5-a0e5-0735f771eae5"
      "created_at" = "2024-03-13 19:15:09.371689973 +0000 UTC"
      "end" = "2024-03-13 06:44:15 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "6ddcf7f7-408d-44ca-8e8f-6d6188d2702b"
      "organization_id" = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
      "project_id" = "53fbca2e-56b7-4951-9912-7baa7396a607"
      "start" = "2024-03-13 04:44:15 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
    {
      "auditlog_download_url" = tostring(null)
      "cluster_id" = "e17880b4-5124-43a5-a0e5-0735f771eae5"
      "created_at" = "2024-03-13 19:20:11.7025117 +0000 UTC"
      "end" = "2024-03-13 19:19:39 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "93d766c7-119a-4ca6-b0e3-9240d22be878"
      "organization_id" = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
      "project_id" = "53fbca2e-56b7-4951-9912-7baa7396a607"
      "start" = "2024-03-13 17:19:39 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
    {
      "auditlog_download_url" = tostring(null)
      "cluster_id" = "e17880b4-5124-43a5-a0e5-0735f771eae5"
      "created_at" = "2024-03-13 19:28:38.869267928 +0000 UTC"
      "end" = "2024-03-13 06:44:15 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "eb39994b-fe95-483f-87c3-41ac4d2ce6d6"
      "organization_id" = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
      "project_id" = "53fbca2e-56b7-4951-9912-7baa7396a607"
      "start" = "2024-03-13 02:44:15 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
    {
      "auditlog_download_url" = tostring(null)
      "cluster_id" = "e17880b4-5124-43a5-a0e5-0735f771eae5"
      "created_at" = "2024-03-13 19:52:37.937857032 +0000 UTC"
      "end" = "2024-03-12 19:19:39 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "3043ce42-4f75-4751-927d-4f8128356ab3"
      "organization_id" = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
      "project_id" = "53fbca2e-56b7-4951-9912-7baa7396a607"
      "start" = "2024-03-12 17:19:39 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
  ])
  "organization_id" = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
  "project_id" = "53fbca2e-56b7-4951-9912-7baa7396a607"
}
new_auditlogexport = {
  "auditlog_download_url" = "https://cb-audit-logs-fabaa8f8-8571-4f8f-961f-6898c15b80d6.s3.amazonaws.com/export/cluster-audit-logs-e17880b4-5124-43a5-a0e5-0735f771eae5-from-2024-03-13-to-2024-03-13.tar.gz?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAT277BFQVAHECFHHJ%2F20240313%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240313T193829Z&X-Amz-Expires=259200&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEFwaCXVzLWVhc3QtMSJHMEUCIQDzq0SS4TWaUJJPTTKwYtogDssCmwufSkoBApiQqHbxYQIgNukuV20pylKGx5kN0ZRftJGV4id3NF8RiPrfrMfKWLsqoAIIZRADGgwyNjQxMzg0NjgzOTQiDPdTLN4GCLNNyZpjnir9AYVC6%2FTJy%2B5Fwl3WI5C90l3cQxkFMd3vVOCXTUBn%2B16VfZimBx26L0R1Zn%2BPROueDRg54U%2FKABUszlLBV2wot0TEvXFNk0QFMXYjCSUl5K3Oj2%2BXSRwX6G6RUJ2%2FGksgl39fwyIP%2F35Jv8KKky%2F4LcKi4hdpJTFWCitI5La4AWBcFV6SRjJ3H%2FaDBEg5XnXgRXHU2XGCzZ6rQd37av5fPIeVMZqGitwYhrAltGqYctDyauteULrouZazafDDu915lx7g7fi%2BOZO4CjYBm820LaU4EJpnx42ETydUEOBD1wt5Bw2fRUCHZ%2Bk2mPFsjXJth%2BQlJxCapUCMhll5%2BCcwsYHIrwY6nQEtJ9%2BEL1WFShM3diIxuMiR6DVFkr5GPa%2F9x0bQ%2B6S4dEohSYagbauSLbgRweKVDu6prGqjcitS3l7IieSecJ1h7ignYGiXRHaQLNIb8PNZp88WiapJZQw0c6cj%2ByS%2FAYXkH2VubuOOynuRuZpiFNlrUADHLT3mP0kzMxPz8%2BeWfj%2ByX0DSJMU2ijimv4zjY4RZzxmUo5NQQDOTinxO&X-Amz-SignedHeaders=host&X-Amz-Signature=8dd866e9520af9e948ac0f1580f14a1ee18e93e6abc22bb3251e2136c16665af"
  "cluster_id" = "e17880b4-5124-43a5-a0e5-0735f771eae5"
  "created_at" = "2024-03-13 19:38:24.355397616 +0000 UTC"
  "end" = "2024-03-13T19:19:39+00:00"
  "expiration" = "2024-03-16 19:38:29.842400291 +0000 UTC"
  "id" = "4c7cb3f9-2ec8-4e58-add0-5f866357dcff"
  "organization_id" = "fabaa8f8-8571-4f8f-961f-6898c15b80d6"
  "project_id" = "53fbca2e-56b7-4951-9912-7baa7396a607"
  "start" = "2024-03-13T17:19:39+00:00"
  "status" = tostring(null)
}
```