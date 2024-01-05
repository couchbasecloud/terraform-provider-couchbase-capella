# Capella Clusters Certificate Example

This example shows how to retrive certificate details for a Capella cluster.

This lists the certificate details based on the cluster ID and authentication access token.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the Capella cluster certificate details as stated in the `get_certificate.tf` file.
2. DELETE: Delete the certificate data output from terraform state.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## GET
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_certificate.existing_certificate: Reading...
data.capella_certificate.existing_certificate: Read complete after 6s

Changes to Outputs:
  + existing_certificate = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = [
          + {
              + certificate = <<-EOT
                    -----BEGIN CERTIFICATE-----
                    MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
                    ...<redacted>...
                    A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfvOIi
                    enG4Dp+hJu9asdxEMRmH70hDyMXv5ZjBhbo39a42QwR59y/rC/sahLLQuNwqif85
                    Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
                    WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
                    rL3ooAw/i0uxr/TQSaxi1l8t4iMt4/gU+W52+8Yol0JbXBTFX6itg62ppb/Eugmn
                    ...<redacted>...
                    -----END CERTIFICATE-----
                EOT
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to get the certificate. 

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_certificate.existing_certificate: Reading...
data.capella_certificate.existing_certificate: Read complete after 7s

Changes to Outputs:
  + existing_certificate = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = [
          + {
              + certificate = <<-EOT
                    -----BEGIN CERTIFICATE-----
                    MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
                    ...<redacted>...
                    A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfvOIi
                    enG4Dp+hJu9asdxEMRmH70hDyMXv5ZjBhbo39a42QwR59y/rC/sahLLQuNwqif85
                    Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
                    WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
                    rL3ooAw/i0uxr/TQSaxi1l8t4iMt4/gU+W52+8Yol0JbXBTFX6itg62ppb/Eugmn
                    ...<redacted>...
                    -----END CERTIFICATE-----
                EOT
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_certificate = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist([
    {
      "certificate" = <<-EOT
      -----BEGIN CERTIFICATE-----
      MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
      ...<redacted>...
      A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfvOIi
      enG4Dp+hJu9asdxEMRmH70hDyMXv5ZjBhbo39a42QwR59y/rC/sahLLQuNwqif85
      Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
      WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
      rL3ooAw/i0uxr/TQSaxi1l8t4iMt4/gU+W52+8Yol0JbXBTFX6itg62ppb/Eugmn
      ...<redacted>...
      EOT
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.capella_certificate.existing_certificate
```

## DESTROY
### Finally, delete the state of the certificate from terraform outputs

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_certificate.existing_certificate: Reading...
data.capella_certificate.existing_certificate: Read complete after 7s

Changes to Outputs:
  - existing_certificate = {
      - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      - data            = [
          - {
              - certificate = <<-EOT
                    -----BEGIN CERTIFICATE-----
                    MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
                    ...<redacted>...
                    A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfvOIi
                    enG4Dp+hJu9asdxEMRmH70hDyMXv5ZjBhbo39a42QwR59y/rC/sahLLQuNwqif85
                    Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
                    WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
                    rL3ooAw/i0uxr/TQSaxi1l8t4iMt4/gU+W52+8Yol0JbXBTFX6itg62ppb/Eugmn
                    ...<redacted>...
                EOT
            },
        ]
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    } -> null

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes


Destroy complete! Resources: 0 destroyed.
```
