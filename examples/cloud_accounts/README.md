# Get Capella Cloud Accounts Data source

This example shows how to get Capella's cloud accounts for an organization: the AWS account, the GCP project and the Azure subscription. Grant them access in your own cloud account, such as an endpoint service's allowed principals on AWS or a Private Link service's visibility on Azure. It uses the organization ID to do so.

`azure_capella_subscription_id` is null when the organization has no Azure account.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_cloud_accounts.existing_cloud_accounts: Reading...
data.couchbase-capella_cloud_accounts.existing_cloud_accounts: Read complete after 0s

Changes to Outputs:
  + existing_cloud_accounts = {
      + aws_capella_account           = "123456789012"
      + azure_capella_subscription    = "capella-azure-sub"
      + azure_capella_subscription_id = "0f1e2d3c-4b5a-4978-8a9b-0c1d2e3f4a5b"
      + gcp_capella_project           = "cb-dp-7a9b1c2d"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.
```

### Apply the Plan

Command: `terraform apply`

Sample Output:

```
$ terraform apply
data.couchbase-capella_cloud_accounts.existing_cloud_accounts: Reading...
data.couchbase-capella_cloud_accounts.existing_cloud_accounts: Read complete after 0s

Changes to Outputs:
  + existing_cloud_accounts = {
      + aws_capella_account           = "123456789012"
      + azure_capella_subscription    = "capella-azure-sub"
      + azure_capella_subscription_id = "0f1e2d3c-4b5a-4978-8a9b-0c1d2e3f4a5b"
      + gcp_capella_project           = "cb-dp-7a9b1c2d"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_cloud_accounts = {
  "aws_capella_account" = "123456789012"
  "azure_capella_subscription" = "capella-azure-sub"
  "azure_capella_subscription_id" = "0f1e2d3c-4b5a-4978-8a9b-0c1d2e3f4a5b"
  "gcp_capella_project" = "cb-dp-7a9b1c2d"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```
