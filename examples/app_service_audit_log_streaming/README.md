# Capella App Service Audit Log Streaming Resource Example

This example shows how to set up and manage audit log streaming for a Capella App Service.

Audit logging must be enabled on the App Service before you configure audit log streaming. Otherwise, the API returns the error "Audit Logging Must Be Enabled To Utilize This Feature".

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Configure audit log streaming for an App Service as stated in the `create_app_service_audit_log_streaming.tf` file.
2. UPDATE: Pause audit log streaming, then change the excluded App Endpoints.
3. IMPORT: Import audit log streaming that is configured for an App Service but is not in the Terraform state file.
4. DELETE: Disable audit log streaming and remove its configuration.

Copy the `terraform.template.tfvars` file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### Apply the plan to configure audit log streaming

Command: `terraform apply`

Sample Output:
```
$ terraform apply
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming will be created
  + resource "couchbase-capella_app_service_audit_log_streaming" "new_app_service_audit_log_streaming" {
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + credentials            = (sensitive value)
      + disabled_app_endpoints = [
          + "<app_endpoint_name>",
        ]
      + log_streaming_state    = (known after apply)
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + output_type            = "datadog"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + streaming_enabled      = true
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service_audit_log_streaming = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming: Creating...
couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming: Creation complete after 7s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_app_service_audit_log_streaming = <sensitive>
```

## UPDATE
### Pause audit log streaming

Set `streaming_enabled = false` in `terraform.tfvars`. When only `streaming_enabled` changes, the provider sends a PATCH request, and Capella keeps the collector configuration.

Command: `terraform apply`

Sample Output:
```
$ terraform apply
  # couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming will be updated in-place
  ~ resource "couchbase-capella_app_service_audit_log_streaming" "new_app_service_audit_log_streaming" {
      ~ log_streaming_state    = "enabled" -> (known after apply)
      ~ streaming_enabled      = true -> false
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
...
Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

### Change other attributes

A change to any other attribute, for example `disabled_app_endpoints`, `output_type` or `credentials`, sends the full configuration in a PUT request.

## IMPORT
### Remove the resource from the Terraform state file

Command: `terraform state rm couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming`

This command removes the resource from the state file only. Audit log streaming stays configured in Capella.

### Import the resource

Command: `terraform import couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
$ terraform import couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming: Importing from ID "app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming: Import prepared!
  Prepared couchbase-capella_app_service_audit_log_streaming for import
couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming: Refreshing state...

Import successful!
```

The API does not return credentials, so `credentials` is null in the state after an import. The next `terraform apply` sends the configured credentials again in a PUT request.

## DELETE
### Disable audit log streaming

Command: `terraform destroy`

The provider sends a PUT request with only `streamingEnabled: false`. This disables audit log streaming and removes the collector configuration, including the credentials.

Sample Output:
```
$ terraform destroy
...
Plan: 0 to add, 0 to change, 1 to destroy.
...
couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming: Destroying...
couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming: Destruction complete after 1m2s

Destroy complete! Resources: 1 destroyed.
```
