# Capella Organizations Example

This example shows how to retrive Organization details in Capella.

This lists the organization details based on the organization ID and authentication access token.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the Capella Organization details as stated in the `get_organization.tf` file.
2. DELETE: Delete the organization data output from terraform state.

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
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 0s

Changes to Outputs:
  + organizations_get = {
      + data            = [
          + {
              + audit       = {
                  + created_at  = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  + created_by  = "fff64e90-e839-4b96-956e-05135e30d35b"
                  + modified_at = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  + modified_by = "fff64e90-e839-4b96-956e-05135e30d35b"
                  + version     = 1
                }
              + description = ""
              + id          = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + name        = "Couchbase"
              + preferences = {
                  + session_duration = null
                }
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
macos:organization $USER$ 
macos:organization $USER$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 1s

Changes to Outputs:
  + organizations_get = {
      + data            = [
          + {
              + audit       = {
                  + created_at  = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  + created_by  = "fff64e90-e839-4b96-956e-05135e30d35b"
                  + modified_at = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  + modified_by = "fff64e90-e839-4b96-956e-05135e30d35b"
                  + version     = 1
                }
              + description = ""
              + id          = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + name        = "Couchbase"
              + preferences = {
                  + session_duration = null
                }
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Bucket

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
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 0s

Changes to Outputs:
  + organizations_get = {
      + data            = [
          + {
              + audit       = {
                  + created_at  = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  + created_by  = "fff64e90-e839-4b96-956e-05135e30d35b"
                  + modified_at = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  + modified_by = "fff64e90-e839-4b96-956e-05135e30d35b"
                  + version     = 1
                }
              + description = ""
              + id          = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + name        = "Couchbase"
              + preferences = {
                  + session_duration = null
                }
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

organizations_get = {
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2022-05-27 04:19:18.057836345 +0000 UTC"
        "created_by" = "fff64e90-e839-4b96-956e-05135e30d35b"
        "modified_at" = "2022-05-27 04:19:18.057836345 +0000 UTC"
        "modified_by" = "fff64e90-e839-4b96-956e-05135e30d35b"
        "version" = 1
      }
      "description" = ""
      "id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "name" = "Couchbase"
      "preferences" = {
        "session_duration" = tonumber(null)
      }
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
}
```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.capella_organization.existing_organization
```

## DESTROY
### Finally, destroy the resources created by Terraform

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
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 0s

Changes to Outputs:
  - organizations_get = {
      - data            = [
          - {
              - audit       = {
                  - created_at  = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  - created_by  = "fff64e90-e839-4b96-956e-05135e30d35b"
                  - modified_at = "2022-05-27 04:19:18.057836345 +0000 UTC"
                  - modified_by = "fff64e90-e839-4b96-956e-05135e30d35b"
                  - version     = 1
                }
              - description = ""
              - id          = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - name        = "Couchbase"
              - preferences = {
                  - session_duration = null
                }
            },
        ]
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    } -> null

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes


Destroy complete! Resources: 0 destroyed.
```
