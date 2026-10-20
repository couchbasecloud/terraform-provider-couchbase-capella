# Capella App Service AWS Command Example

This example shows how to retrieve the AWS command used to configure a VPC endpoint connecting to an App Service's (Sync Gateway) private endpoint service.

The data source also exposes a `service_name` attribute, parsed best-effort from the `--service-name` flag in `command` (the App Service private endpoint status API, unlike the operational cluster's, has no dedicated service name field). It can be fed directly into an `aws_vpc_endpoint` resource's `service_name` argument; if Couchbase changes the command's phrasing it is left null rather than failing the read.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Display the AWS command as stated in the `get_command.tf` file.

## GET

Command: `terraform apply`

Sample Output:
```
data.couchbase-capella_app_service_aws_private_endpoint_command.aws_command: Reading...
data.couchbase-capella_app_service_aws_private_endpoint_command.aws_command: Read complete after 1s

Changes to Outputs:
  + aws_command = {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + command         = "aws ec2 create-vpc-endpoint --vpc-id vpc-1234 --region us-east-1 --service-name com.amazonaws.vpce.us-east-1.vpce-svc-1234 --vpc-endpoint-type Interface --subnet-ids subnet-1234"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_name    = "com.amazonaws.vpce.us-east-1.vpce-svc-1234"
      + subnet_ids      = [
          + "subnet-1234",
        ]
      + vpc_id          = "vpc-1234"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

aws_command = {
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "command" = "aws ec2 create-vpc-endpoint --vpc-id vpc-1234 --region us-east-1 --service-name com.amazonaws.vpce.us-east-1.vpce-svc-1234 --vpc-endpoint-type Interface --subnet-ids subnet-1234"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "service_name" = "com.amazonaws.vpce.us-east-1.vpce-svc-1234"
  "subnet_ids" = toset([
    "subnet-1234",
  ])
  "vpc_id" = "vpc-1234"
}
```
