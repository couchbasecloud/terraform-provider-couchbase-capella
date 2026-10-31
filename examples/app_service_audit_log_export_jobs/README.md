# Capella App Service Audit Log Export Jobs Example

This example lists the audit log export jobs of a Capella App Service.

It uses the organization, project, cluster and App Service IDs to list the export jobs. Set the optional `export_id` argument in `list_app_service_audit_log_export_jobs.tf` to return only one export job.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

If you check the `terraform.template.tfvars` file, you can see the variables that you need to run the terraform commands.
Copy the file to `terraform.tfvars` and update the values of the variables for your organization.

## LIST
### View the export jobs

Command: `terraform plan`

The export jobs are shown in the `app_service_audit_log_export_jobs` output.
