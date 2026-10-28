# Capella App Service Audit Log Resource Example

This example shows how to enable audit logging on a Capella App Service.

The resource uses the organization ID, project ID, cluster ID and App Service ID. When you create the resource, Terraform enables audit logging on the App Service. When you destroy the resource, Terraform disables audit logging on the App Service.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Enable audit logging on the App Service as stated in the `create_app_service_audit_log.tf` file.
2. IMPORT: Import the App Service audit log state into Terraform.
3. DELETE: Disable audit logging on the App Service.

The resource does not support in-place updates. A change to any attribute replaces the resource.

Copy the `terraform.template.tfvars` file to `terraform.tfvars` and update the values of the variables for your organization.

## CREATE

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_audit_log.new_app_service_audit_log will be created
  + resource "couchbase-capella_app_service_audit_log" "new_app_service_audit_log" {
      + app_service_id  = "9a34f4a3-b05e-4e6e-a249-c472f9f915dd"
      + cluster_id      = "669cb38a-ea74-42df-94fe-bdca9d332e1c"
      + organization_id = "a16247ef-8356-4a09-9edb-7fd2507a2881"
      + project_id      = "e5dbe56f-43bb-482b-a161-80bb7430dbbd"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the plan to enable audit logging

Command: `terraform apply`

Sample Output:
```
$ terraform apply
...
couchbase-capella_app_service_audit_log.new_app_service_audit_log: Creating...
couchbase-capella_app_service_audit_log.new_app_service_audit_log: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## IMPORT

### Remove the resource from the Terraform state file

Command: `terraform state rm couchbase-capella_app_service_audit_log.new_app_service_audit_log`

This command removes the resource only from the Terraform state file. Audit logging stays enabled on the App Service in Capella.

### Import the resource into Terraform

Command: `terraform import couchbase-capella_app_service_audit_log.new_app_service_audit_log app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
$ terraform import couchbase-capella_app_service_audit_log.new_app_service_audit_log app_service_id=9a34f4a3-b05e-4e6e-a249-c472f9f915dd,cluster_id=669cb38a-ea74-42df-94fe-bdca9d332e1c,project_id=e5dbe56f-43bb-482b-a161-80bb7430dbbd,organization_id=a16247ef-8356-4a09-9edb-7fd2507a2881
couchbase-capella_app_service_audit_log.new_app_service_audit_log: Importing from ID "app_service_id=9a34f4a3-b05e-4e6e-a249-c472f9f915dd,cluster_id=669cb38a-ea74-42df-94fe-bdca9d332e1c,project_id=e5dbe56f-43bb-482b-a161-80bb7430dbbd,organization_id=a16247ef-8356-4a09-9edb-7fd2507a2881"...
couchbase-capella_app_service_audit_log.new_app_service_audit_log: Import prepared!
  Prepared couchbase-capella_app_service_audit_log for import
couchbase-capella_app_service_audit_log.new_app_service_audit_log: Refreshing state...

Import successful!
```

The import ID is a comma-separated string of the App Service ID, cluster ID, project ID and organization ID.
If audit logging is disabled on the App Service, the import fails because the resource does not exist.

## DESTROY

### Disable audit logging and destroy all resources managed by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
...
couchbase-capella_app_service_audit_log.new_app_service_audit_log: Destroying...
couchbase-capella_app_service_audit_log.new_app_service_audit_log: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
