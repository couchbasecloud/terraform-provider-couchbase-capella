# Capella Audit Log Export Example

This example shows how to create and audit log export job in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new audit log export job in Capella per the `create_audit_log_exports.tf` file.
2. LIST: List existing audit log export jobs in Capella as stated in the `list_audit_log_exports.tf` file, and
          detecting if resource was updated outside terraform. 
3. UPDATE: Show it is not supported
4. IMPORT:  Import an existing audit log export job
5. DELETE:  It is not supported by API server so it's just a noop (ie removes resource from state file).


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
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
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
      + audit_log_download_url = (known after apply)
      + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + start                 = "2024-03-13T04:44:15+00:00"
      + status                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_auditlogexport = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_auditlogexport      = {
      + audit_log_download_url = (known after apply)
      + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
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
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
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
      + audit_log_download_url = (known after apply)
      + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + start                 = "2024-03-13T04:44:15+00:00"
      + status                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_auditlogexport = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_auditlogexport      = {
      + audit_log_download_url = (known after apply)
      + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + created_at            = (known after apply)
      + end                   = "2024-03-13T06:44:15+00:00"
      + expiration            = (known after apply)
      + id                    = (known after apply)
      + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + start                 = "2024-03-13T04:44:15+00:00"
      + status                = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_audit_log_export.new_auditlogexport: Creating...
couchbase-capella_audit_log_export.new_auditlogexport: Creation complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

existing_auditlogexport = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = toset(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_auditlogexport = {
  "audit_log_download_url" = tostring(null)
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "created_at" = "2024-03-13 06:48:43.596764031 +0000 UTC"
  "end" = "2024-03-13T06:44:15+00:00"
  "expiration" = tostring(null)
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
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
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_audit_log_export.new_auditlogexport has changed
  ~ resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      + audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX"
      + expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC"
        id                    = "ffffffff-aaaa-1414-eeee-000000000000"
        # (6 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may
include actions to undo or respond to these changes.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Changes to Outputs:
  ~ existing_auditlogexport = {
      ~ data            = null -> [
          + {
              + audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX"
              + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
              + created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC"
              + end                   = "2024-03-13 06:44:15 +0000 UTC"
              + expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC"
              + id                    = "ffffffff-aaaa-1414-eeee-000000000000"
              + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
              + start                 = "2024-03-13 04:44:15 +0000 UTC"
              + status                = "Completed"
            },
        ]
        # (3 unchanged attributes hidden)
    }
  ~ new_auditlogexport      = {
      ~ audit_log_download_url = null -> "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX"
      ~ expiration            = null -> "2024-03-16 06:48:46.64587117 +0000 UTC"
        id                    = "ffffffff-aaaa-1414-eeee-000000000000"
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
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_audit_log_export.new_auditlogexport will be updated in-place
  ~ resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      ~ audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX" -> (known after apply)
      ~ created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC" -> (known after apply)
      ~ expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC" -> (known after apply)
        id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ start                 = "2024-03-13T04:44:15+00:00" -> "2024-03-13T02:44:15+00:00"
      + status                = (known after apply)
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_auditlogexport      = {
      ~ audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX" -> (known after apply)
      ~ created_at            = "2024-03-13 06:48:43.596764031 +0000 UTC" -> (known after apply)
      ~ expiration            = "2024-03-16 06:48:46.64587117 +0000 UTC" -> (known after apply)
        id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ start                 = "2024-03-13T04:44:15+00:00" -> "2024-03-13T02:44:15+00:00"
      + status                = (known after apply)
        # (4 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_audit_log_export.new_auditlogexport: Modifying... [id=ffffffff-aaaa-1414-eeee-000000000000]
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

Command:  `terraform import couchbase-capella_audit_log_export.new_auditlogexport id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
terraform import couchbase-capella_audit_log_export.new_auditlogexport id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_audit_log_export.new_auditlogexport: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Import prepared!
  Prepared couchbase-capella_audit_log_export for import
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]

2024-03-13T12:42:14.621-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_audit_log_export.new_auditlogexport during refresh.
      - .audit_log_download_url: was null, but now cty.StringVal("https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX")
      - .created_at: was null, but now cty.StringVal("2024-03-13 19:38:24.355397616 +0000 UTC")
      - .end: was null, but now cty.StringVal("2024-03-13T19:19:39+00:00")
      - .id: was cty.StringVal("id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"), but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .project_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .cluster_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .expiration: was null, but now cty.StringVal("2024-03-16 19:38:29.842400291 +0000 UTC")
      - .organization_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
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
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_audit_log_export.new_auditlogexport will be updated in-place
  ~ resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
      ~ audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX" -> (known after apply)
      ~ created_at            = "2024-03-13 19:38:24.355397616 +0000 UTC" -> (known after apply)
      ~ end                   = "2024-03-13T19:19:39+00:00" -> "2024-03-13T06:44:15+00:00"
      ~ expiration            = "2024-03-16 19:38:29.842400291 +0000 UTC" -> (known after apply)
        id                    = "ffffffff-aaaa-1414-eeee-000000000000"
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
                audit_log_download_url = null
                cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
                created_at            = "2024-03-13 19:28:38.869267928 +0000 UTC"
                end                   = "2024-03-13 06:44:15 +0000 UTC"
                expiration            = null
                id                    = "ffffffff-aaaa-1414-eeee-000000000000"
                organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
                project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
                start                 = "2024-03-13 02:44:15 +0000 UTC"
                status                = "no audit log files exist within the requested time frame"
            },
          + {
              + audit_log_download_url = null
              + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
              + created_at            = "2024-03-13 19:52:37.937857032 +0000 UTC"
              + end                   = "2024-03-12 19:19:39 +0000 UTC"
              + expiration            = null
              + id                    = "3043ce42-4f75-4751-927d-4f8128356ab3"
              + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
              + start                 = "2024-03-12 17:19:39 +0000 UTC"
              + status                = "no audit log files exist within the requested time frame"
            },
        ]
        # (3 unchanged attributes hidden)
    }
  ~ new_auditlogexport      = {
      ~ audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX" -> (known after apply)
      ~ created_at            = "2024-03-13 19:38:24.355397616 +0000 UTC" -> (known after apply)
      ~ end                   = "2024-03-13T19:19:39+00:00" -> "2024-03-13T06:44:15+00:00"
      ~ expiration            = "2024-03-16 19:38:29.842400291 +0000 UTC" -> (known after apply)
        id                    = "ffffffff-aaaa-1414-eeee-000000000000"
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
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Changes to Outputs:
  ~ existing_auditlogexport = {
      ~ data            = [
            # (3 unchanged elements hidden)
            {
                audit_log_download_url = null
                cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
                created_at            = "2024-03-13 19:28:38.869267928 +0000 UTC"
                end                   = "2024-03-13 06:44:15 +0000 UTC"
                expiration            = null
                id                    = "ffffffff-aaaa-1414-eeee-000000000000"
                organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
                project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
                start                 = "2024-03-13 02:44:15 +0000 UTC"
                status                = "no audit log files exist within the requested time frame"
            },
          + {
              + audit_log_download_url = null
              + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000000"
              + created_at            = "2024-03-13 19:52:37.937857032 +0000 UTC"
              + end                   = "2024-03-12 19:19:39 +0000 UTC"
              + expiration            = null
              + id                    = "3043ce42-4f75-4751-927d-4f8128356ab3"
              + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_id            = "ffffffff-aaaa-1414-eeee-000000000000"
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
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = toset([
    {
      "audit_log_download_url" = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX"
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "created_at" = "2024-03-13 19:38:24.355397616 +0000 UTC"
      "end" = "2024-03-13 19:19:39 +0000 UTC"
      "expiration" = "2024-03-16 19:38:29.842400291 +0000 UTC"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "start" = "2024-03-13 17:19:39 +0000 UTC"
      "status" = "Completed"
    },
    {
      "audit_log_download_url" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "created_at" = "2024-03-13 19:15:09.371689973 +0000 UTC"
      "end" = "2024-03-13 06:44:15 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "6ddcf7f7-408d-44ca-8e8f-6d6188d2702b"
      "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "start" = "2024-03-13 04:44:15 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
    {
      "audit_log_download_url" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "created_at" = "2024-03-13 19:20:11.7025117 +0000 UTC"
      "end" = "2024-03-13 19:19:39 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "93d766c7-119a-4ca6-b0e3-9240d22be878"
      "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "start" = "2024-03-13 17:19:39 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
    {
      "audit_log_download_url" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "created_at" = "2024-03-13 19:28:38.869267928 +0000 UTC"
      "end" = "2024-03-13 06:44:15 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "start" = "2024-03-13 02:44:15 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
    {
      "audit_log_download_url" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "created_at" = "2024-03-13 19:52:37.937857032 +0000 UTC"
      "end" = "2024-03-12 19:19:39 +0000 UTC"
      "expiration" = tostring(null)
      "id" = "3043ce42-4f75-4751-927d-4f8128356ab3"
      "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "start" = "2024-03-12 17:19:39 +0000 UTC"
      "status" = "no audit log files exist within the requested time frame"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_auditlogexport = {
  "audit_log_download_url" = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "created_at" = "2024-03-13 19:38:24.355397616 +0000 UTC"
  "end" = "2024-03-13T19:19:39+00:00"
  "expiration" = "2024-03-16 19:38:29.842400291 +0000 UTC"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "start" = "2024-03-13T17:19:39+00:00"
  "status" = tostring(null)
}
```

## DESTROY
### This only removes resource from the state file.

Command: `terraform destroy`

Sample Output:
```
terraform destroy
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_audit_log_export.existing_auditlogexport: Reading...
couchbase-capella_audit_log_export.new_auditlogexport: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_audit_log_export.existing_auditlogexport: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
- destroy

Terraform will perform the following actions:

# couchbase-capella_audit_log_export.new_auditlogexport will be destroyed
- resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
    - audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX" -> null
    - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    - created_at             = "2024-06-10 21:16:50.415705947 +0000 UTC" -> null
    - end                    = "2024-06-10T21:09:11+00:00" -> null
    - expiration             = "2024-06-13 21:16:54.37230845 +0000 UTC" -> null
    - id                     = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    - project_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    - start                  = "2024-06-10T20:09:11+00:00" -> null
      }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
- existing_auditlogexport = {
    - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    - data            = [
        - {
            - audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX"
            - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - created_at             = "2024-06-10 21:16:50.415705947 +0000 UTC"
            - end                    = "2024-06-10 21:09:11 +0000 UTC"
            - expiration             = "2024-06-13 21:16:54.37230845 +0000 UTC"
            - id                     = "ffffffff-aaaa-1414-eeee-000000000000"
            - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
            - project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - start                  = "2024-06-10 20:09:11 +0000 UTC"
            - status                 = "Completed"
              },
        - {
            - audit_log_download_url = null
            - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - created_at             = "2024-06-10 21:03:15.896348051 +0000 UTC"
            - end                    = "2024-05-28 21:52:38 +0000 UTC"
            - expiration             = null
            - id                     = "058091cf-f966-46a7-96ea-4801ad31bfb2"
            - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
            - project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - start                  = "2024-05-28 18:52:38 +0000 UTC"
            - status                 = "no audit log files exist within the requested time frame"
              },
        - {
            - audit_log_download_url = null
            - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - created_at             = "2024-06-10 21:03:50.081716726 +0000 UTC"
            - end                    = "2024-05-28 21:52:38 +0000 UTC"
            - expiration             = null
            - id                     = "26cb94d8-5d86-4964-ab5a-bd36fedbf13f"
            - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
            - project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - start                  = "2024-05-28 18:52:38 +0000 UTC"
            - status                 = "no audit log files exist within the requested time frame"
              },
        - {
            - audit_log_download_url = null
            - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - created_at             = "2024-06-10 21:12:19.656434023 +0000 UTC"
            - end                    = "2024-06-10 21:09:11 +0000 UTC"
            - expiration             = null
            - id                     = "23b12e1a-d84d-48b6-8d09-4df2e5e6b4a7"
            - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
            - project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - start                  = "2024-06-10 20:09:11 +0000 UTC"
            - status                 = "no audit log files exist within the requested time frame"
              },
        - {
            - audit_log_download_url = null
            - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - created_at             = "2024-06-10 21:14:28.070136553 +0000 UTC"
            - end                    = "2024-06-10 21:09:11 +0000 UTC"
            - expiration             = null
            - id                     = "f658e557-45b7-4be7-a25d-d6c9edc9a122"
            - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
            - project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
            - start                  = "2024-06-10 20:09:11 +0000 UTC"
            - status                 = "no audit log files exist within the requested time frame"
              },
              ]
    - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
    - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      } -> null
- new_auditlogexport      = {
    - audit_log_download_url = "https://cb-audit-logs-1234.s3.amazonaws.com/export/cluster-audit-logs-5678-from-2000-01-01-to-2000-01-02.tar.gz?X-Amz-Algorithm=X&X-Amz-Credential=XXX&X-Amz-Date=1999&X-Amz-Expires=1&X-Amz-Security-Token=XXX&X-Amz-SignedHeaders=host&X-Amz-Signature=XXX"
    - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
    - created_at             = "2024-06-10 21:16:50.415705947 +0000 UTC"
    - end                    = "2024-06-10T21:09:11+00:00"
    - expiration             = "2024-06-13 21:16:54.37230845 +0000 UTC"
    - id                     = "ffffffff-aaaa-1414-eeee-000000000000"
    - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    - project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
    - start                  = "2024-06-10T20:09:11+00:00"
    - status                 = null
      } -> null

Do you really want to destroy all resources?
Terraform will destroy all your managed infrastructure, as shown above.
There is no undo. Only 'yes' will be accepted to confirm.

Enter a value: yes

couchbase-capella_audit_log_export.new_auditlogexport: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_audit_log_export.new_auditlogexport: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```