# Capella AWS Command Example

This example shows how to retrieve the AWS command used to configure a VPC endpoint.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Display the AWS command as stated in the `get_config.tf` file.

## GET

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_aws_private_endpoint_command.aws_command: Reading...
data.couchbase-capella_aws_private_endpoint_command.aws_command: Read complete after 1s

Changes to Outputs:
  + aws_command = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + command         = "aws ec2 create-vpc-endpoint --vpc-id vpc-1234 --region us-east-1 --service-name com.amazonaws.vpce.us-east-1.vpce-svc-1234 --vpc-endpoint-type Interface --subnet-ids "
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + subnet_ids      = null
      + vpc_id          = "vpc-1234"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```