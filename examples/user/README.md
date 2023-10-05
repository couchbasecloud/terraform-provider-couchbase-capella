# Capella User Example

This example shows how to create and manage users in Capella.

This creates a new user in the selected Capella project. It uses the organization ID and projectId to do so.

An invitation email is triggered and sent to the user. Upon receiving the invitation email, the user is required to click on a provided URL, which will redirect them to a page with a user interface (UI) where they can set their username and password.

The modification of any personal information related to a user can only be performed by the user through the UI. Similarly, the user can solely conduct password updates through the UI.

The "caller" possessing Organization Owner access rights retains the exclusive user creation capability. They hold the authority to assign roles at the organization and project levels.

At present, our support is limited to the capella resourceType of "project" exclusively.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Create a new user in Capella as stated in the `create_user.tf` file.
2. List existing users in Capella as stated in the `list_users.tf` file.
3. Import a user that exists in Capella but not in the terraform state file.
4. Delete the newly created user from Capella.

If you check the `terraform.template.tfvars` file - you can see that we need 3 main variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

### View the plan for the resources that Terraform will create

Command: `terraform plan`
Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/mattymaclean/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "matty.maclean+4@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "1a3c4544-772e-449e-9996-1203e7020b96"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5"
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_user = {
      + audit                = (known after apply)
      + email                = "matty.maclean+4@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "1a3c4544-772e-449e-9996-1203e7020b96"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5"
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new User in Capella
Command: `terraform apply`
Sample Output:
```
$ terraform apply
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/mattymaclean/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "matty.maclean+4@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "1a3c4544-772e-449e-9996-1203e7020b96"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5"
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_user = {
      + audit                = (known after apply)
      + email                = "matty.maclean+4@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "1a3c4544-772e-449e-9996-1203e7020b96"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5"
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_user.new_user: Creating...
capella_user.new_user: Creation complete after 1s [id=ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_user = {
  "audit" = {
    "created_at" = "2023-10-05 11:09:52.803123107 +0000 UTC"
    "created_by" = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec"
    "modified_at" = "2023-10-05 11:09:52.803123107 +0000 UTC"
    "modified_by" = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec"
    "version" = 1
  }
  "email" = "matty.maclean+4@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-03T11:09:52.80312439Z"
  "id" = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec"
  "inactive" = true
  "last_login" = ""
  "name" = "Matty"
  "organization_id" = "1a3c4544-772e-449e-9996-1203e7020b96"
  "organization_roles" = tolist([
    "projectCreator",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5"
      "roles" = tolist([
        "projectDataReaderWriter",
      ])
      "type" = "project"
    },
  ])
  "status" = "not-verified"
  "time_zone" = ""
}
```
### Note the User ID for the new User
Command: `terraform show`

Sample Output:
```
```
### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
```
### Remove the resource `new_user` from the Terraform State file

Command: `terraform state rm capella_project.new_user`

Sample Output:
```
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.
### Now, let's import the resource in Terraform
Command: `terraform import capella_user.new_user id=<user_id>,organization_id=<organization_id>`

In this case, the complete command is:
``
Sample Output:
```
```
### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted
Command: `terraform plan`

Sample Output:
```
```

### Finally, destroy the resources created by Terraform
Command: `terraform destroy`
Sample Output:
```
$ terraform destroy
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/mattymaclean/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
capella_user.new_user: Refreshing state... [id=ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_user.new_user will be destroyed
  - resource "capella_user" "new_user" {
      - audit                = {
          - created_at  = "2023-10-05 11:09:52.803123107 +0000 UTC" -> null
          - created_by  = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec" -> null
          - modified_at = "2023-10-05 11:09:52.803123107 +0000 UTC" -> null
          - modified_by = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec" -> null
          - version     = 1 -> null
        } -> null
      - email                = "matty.maclean+4@couchbase.com" -> null
      - enable_notifications = false -> null
      - expires_at           = "2024-01-03T11:09:52.80312439Z" -> null
      - id                   = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec" -> null
      - inactive             = true -> null
      - name                 = "Matty" -> null
      - organization_id      = "1a3c4544-772e-449e-9996-1203e7020b96" -> null
      - organization_roles   = [
          - "projectCreator",
        ] -> null
      - resources            = [
          - {
              - id    = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5" -> null
              - roles = [
                  - "projectDataReaderWriter",
                ] -> null
              - type  = "project" -> null
            },
        ] -> null
      - status               = "not-verified" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_user = {
      - audit                = {
          - created_at  = "2023-10-05 11:09:52.803123107 +0000 UTC"
          - created_by  = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec"
          - modified_at = "2023-10-05 11:09:52.803123107 +0000 UTC"
          - modified_by = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec"
          - version     = 1
        }
      - email                = "matty.maclean+4@couchbase.com"
      - enable_notifications = false
      - expires_at           = "2024-01-03T11:09:52.80312439Z"
      - id                   = "ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec"
      - inactive             = true
      - last_login           = ""
      - name                 = "Matty"
      - organization_id      = "1a3c4544-772e-449e-9996-1203e7020b96"
      - organization_roles   = [
          - "projectCreator",
        ]
      - region               = ""
      - resources            = [
          - {
              - id    = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5"
              - roles = [
                  - "projectDataReaderWriter",
                ]
              - type  = "project"
            },
        ]
      - status               = "not-verified"
      - time_zone            = ""
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_user.new_user: Destroying... [id=ea22e4e8-d59d-4a31-aa0e-27cc33ca67ec]
capella_user.new_user: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
