# Capella App Service GCP Command Example

This example shows how to retrieve the GCP command used to configure a Private Service Connect endpoint connecting to an App Service's (Sync Gateway) private endpoint service.

**AWS-only as of this writing:** Couchbase Capella's private endpoints for App Services are AWS-only; GCP support is not confirmed GA. This data source is included for API/schema parity with the cluster-level equivalent, but has not been validated against a live GCP-backed App Service.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Display the GCP command as stated in the `get_command.tf` file.

## GET

Command: `terraform apply`

Sample Output:
```
data.couchbase-capella_app_service_gcp_private_endpoint_command.gcp_command: Reading...
data.couchbase-capella_app_service_gcp_private_endpoint_command.gcp_command: Read complete after 1s

Changes to Outputs:
  + gcp_command = {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + command         = <<-EOT
            # ... gcloud compute forwarding-rules create / dns managed-zones commands ...
        EOT
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + subnet_ids      = [
          + "subnet-1234",
        ]
      + vpc_network_id  = "vpc-network-1234"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```
