# Capella Organization Example

This example shows how to manage Organization in Capella.

This fetches the details of an existing Organization. It uses the organization ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Get an existing organization in Capella as stated in the `get_organization.tf` file.

If you check the `terraform.template.tfvars` file - you can see that we need 3 main variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.


### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
terraform plan 
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 1s [name=cbc-dev]

Changes to Outputs:
  + existing_organization = {
      + audit           = {
          + created_at  = "2020-07-22 12:38:57.437248116 +0000 UTC"
          + created_by  = ""
          + modified_at = "2023-07-25 14:33:56.13967014 +0000 UTC"
          + modified_by = "99b8dc97-b8ae-44af-8ccd-897e3802c3cb"
          + version     = 0
        }
      + description     = ""
      + name            = "cbc-dev"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + preferences     = {
          + session_duration = 7200
        }
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.


```

### Apply the Plan, in order to get the organization

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 1s [name=cbc-dev]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_organization = {
  "audit" = {
    "created_at" = "2020-07-22 12:38:57.437248116 +0000 UTC"
    "created_by" = ""
    "modified_at" = "2023-07-25 14:33:56.13967014 +0000 UTC"
    "modified_by" = "99b8dc97-b8ae-44af-8ccd-897e3802c3cb"
    "version" = 0
  }
  "description" = ""
  "name" = "cbc-dev"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "preferences" = {
    "session_duration" = 7200
  }
}
nidhi.kumar@QFXY6XF4V3 organization % terraform plan 
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 1s [name=cbc-dev]

Changes to Outputs:
  + existing_organization = {
      + audit           = {
          + created_at  = "2020-07-22 12:38:57.437248116 +0000 UTC"
          + created_by  = ""
          + modified_at = "2023-07-25 14:33:56.13967014 +0000 UTC"
          + modified_by = "99b8dc97-b8ae-44af-8ccd-897e3802c3cb"
          + version     = 0
        }
      + description     = ""
      + name            = "cbc-dev"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + preferences     = {
          + session_duration = 7200
        }
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
nidhi.kumar@QFXY6XF4V3 organization % terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 1s [name=cbc-dev]

Changes to Outputs:
  + existing_organization = {
      + audit           = {
          + created_at  = "2020-07-22 12:38:57.437248116 +0000 UTC"
          + created_by  = ""
          + modified_at = "2023-07-25 14:33:56.13967014 +0000 UTC"
          + modified_by = "99b8dc97-b8ae-44af-8ccd-897e3802c3cb"
          + version     = 0
        }
      + description     = ""
      + name            = "cbc-dev"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + preferences     = {
          + session_duration = 7200
        }
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_organization = {
  "audit" = {
    "created_at" = "2020-07-22 12:38:57.437248116 +0000 UTC"
    "created_by" = ""
    "modified_at" = "2023-07-25 14:33:56.13967014 +0000 UTC"
    "modified_by" = "99b8dc97-b8ae-44af-8ccd-897e3802c3cb"
    "version" = 0
  }
  "description" = ""
  "name" = "cbc-dev"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "preferences" = {
    "session_duration" = 7200
  }
}

```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.capella_organization.existing_organization
```
