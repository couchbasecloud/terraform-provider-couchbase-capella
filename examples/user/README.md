# Capella User Example

This example shows how to create and manage users in Capella.

This creates a new user in the selected Capella project. It uses the organization ID and projectId to do so.

An invitation email is triggered and sent to the user. Upon receiving the invitation email, the user is required to click on a provided URL, which will redirect them to a page with a user interface (UI) where they can set their username and password.

The modification of any personal information related to a user can only be performed by the user through the UI. Similarly, the user can solely conduct password updates through the UI.

The "caller" possessing Organization Owner access rights retains the exclusive user creation capability. They hold the authority to assign roles at the organization and project levels.

At present, our support is limited to the capella resourceType of "project" exclusively.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough
### View the plan for the resources that Terraform will create

Command: `terraform plan`
Sample Output:
```
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/mattymaclean/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
╵
data.capella_users.existing_users: Reading...
data.capella_users.existing_users: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  + create

Terraform will perform the following actions:

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "matty.maclean+2@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
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
  + new_user   = {
      + audit                = (known after apply)
      + email                = "matty.maclean+2@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }
  + users_list = {
      + data            = [
          + {
              + audit                = {
                  + created_at  = "2023-10-06 15:43:02.805868342 +0000 UTC"
                  + created_by  = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                  + modified_at = "2023-10-06 15:43:02.805868342 +0000 UTC"
                  + modified_by = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                  + version     = 1
                }
              + email                = "matty.maclean@couchbase.com"
              + enable_notifications = false
              + expires_at           = "2024-01-04T15:43:02.805868342Z"
              + id                   = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
              + inactive             = false
              + last_login           = "2023-10-06T15:47:57.491646422Z"
              + name                 = "matty.maclean"
              + organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
              + organization_roles   = [
                  + "organizationOwner",
                ]
              + region               = ""
              + resources            = null
              + status               = "verified"
              + time_zone            = ""
            },
        ]
      + organization_id = "93f13778-3d11-43c5-861f-417a4b00ba81"
    }

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run
"terraform apply" now.
```

### Apply the Plan, in order to create a new User in Capella
Command: `terraform apply`
Sample Output:
```
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/mattymaclean/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_users.existing_users: Reading...
data.capella_users.existing_users: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "matty.maclean+2@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
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
  + new_user   = {
      + audit                = (known after apply)
      + email                = "matty.maclean+2@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "Matty"
      + organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
      + organization_roles   = [
          + "projectCreator",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }
  + users_list = {
      + data            = [
          + {
              + audit                = {
                  + created_at  = "2023-10-06 15:43:02.805868342 +0000 UTC"
                  + created_by  = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                  + modified_at = "2023-10-06 15:43:02.805868342 +0000 UTC"
                  + modified_by = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                  + version     = 1
                }
              + email                = "matty.maclean@couchbase.com"
              + enable_notifications = false
              + expires_at           = "2024-01-04T15:43:02.805868342Z"
              + id                   = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
              + inactive             = false
              + last_login           = "2023-10-06T15:47:57.491646422Z"
              + name                 = "matty.maclean"
              + organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
              + organization_roles   = [
                  + "organizationOwner",
                ]
              + region               = ""
              + resources            = null
              + status               = "verified"
              + time_zone            = ""
            },
        ]
      + organization_id = "93f13778-3d11-43c5-861f-417a4b00ba81"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_user.new_user: Creating...
capella_user.new_user: Creation complete after 1s [id=9ddcf5d2-901e-457c-9d62-4709ef0eb46d]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_user = {
  "audit" = {
    "created_at" = "2023-10-06 16:05:06.620419302 +0000 UTC"
    "created_by" = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
    "modified_at" = "2023-10-06 16:05:06.620419302 +0000 UTC"
    "modified_by" = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
    "version" = 1
  }
  "email" = "matty.maclean+2@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-04T16:05:06.620419427Z"
  "id" = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
  "inactive" = true
  "last_login" = ""
  "name" = "Matty"
  "organization_id" = "93f13778-3d11-43c5-861f-417a4b00ba81"
  "organization_roles" = tolist([
    "projectCreator",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
      "roles" = tolist([
        "projectDataReaderWriter",
      ])
      "type" = "project"
    },
  ])
  "status" = "not-verified"
  "time_zone" = ""
}
users_list = {
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2023-10-06 15:43:02.805868342 +0000 UTC"
        "created_by" = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
        "modified_at" = "2023-10-06 15:43:02.805868342 +0000 UTC"
        "modified_by" = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
        "version" = 1
      }
      "email" = "matty.maclean@couchbase.com"
      "enable_notifications" = false
      "expires_at" = "2024-01-04T15:43:02.805868342Z"
      "id" = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
      "inactive" = false
      "last_login" = "2023-10-06T15:47:57.491646422Z"
      "name" = "matty.maclean"
      "organization_id" = "93f13778-3d11-43c5-861f-417a4b00ba81"
      "organization_roles" = tolist([
        "organizationOwner",
      ])
      "region" = ""
      "resources" = tolist(null) /* of object */
      "status" = "verified"
      "time_zone" = ""
    },
  ])
  "organization_id" = "93f13778-3d11-43c5-861f-417a4b00ba81"
}
```
### Note the User ID for the new User
Command: `terraform show`

Sample Output:
```
# data.capella_users.existing_users:
data "capella_users" "existing_users" {
    data            = [
        {
            audit                = {
                created_at  = "2023-10-06 15:43:02.805868342 +0000 UTC"
                created_by  = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                modified_at = "2023-10-06 15:43:02.805868342 +0000 UTC"
                modified_by = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                version     = 1
            }
            email                = "matty.maclean@couchbase.com"
            enable_notifications = false
            expires_at           = "2024-01-04T15:43:02.805868342Z"
            id                   = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
            inactive             = false
            last_login           = "2023-10-06T15:47:57.491646422Z"
            name                 = "matty.maclean"
            organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
            organization_roles   = [
                "organizationOwner",
            ]
            region               = ""
            status               = "verified"
            time_zone            = ""
        },
    ]
    organization_id = "93f13778-3d11-43c5-861f-417a4b00ba81"
}

# capella_user.new_user:
resource "capella_user" "new_user" {
    audit                = {
        created_at  = "2023-10-06 16:05:06.620419302 +0000 UTC"
        created_by  = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
        modified_at = "2023-10-06 16:05:06.620419302 +0000 UTC"
        modified_by = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
        version     = 1
    }
    email                = "matty.maclean+2@couchbase.com"
    enable_notifications = false
    expires_at           = "2024-01-04T16:05:06.620419427Z"
    id                   = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
    inactive             = true
    name                 = "Matty"
    organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
    organization_roles   = [
        "projectCreator",
    ]
    resources            = [
        {
            id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
            roles = [
                "projectDataReaderWriter",
            ]
            type  = "project"
        },
    ]
    status               = "not-verified"
}


Outputs:

new_user = {
    audit                = {
        created_at  = "2023-10-06 16:05:06.620419302 +0000 UTC"
        created_by  = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
        modified_at = "2023-10-06 16:05:06.620419302 +0000 UTC"
        modified_by = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
        version     = 1
    }
    email                = "matty.maclean+2@couchbase.com"
    enable_notifications = false
    expires_at           = "2024-01-04T16:05:06.620419427Z"
    id                   = "9ddcf5d2-901e-457c-9d62-4709ef0eb46d"
    inactive             = true
    last_login           = ""
    name                 = "Matty"
    organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
    organization_roles   = [
        "projectCreator",
    ]
    region               = ""
    resources            = [
        {
            id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
            roles = [
                "projectDataReaderWriter",
            ]
            type  = "project"
        },
    ]
    status               = "not-verified"
    time_zone            = ""
}
users_list = {
    data            = [
        {
            audit                = {
                created_at  = "2023-10-06 15:43:02.805868342 +0000 UTC"
                created_by  = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                modified_at = "2023-10-06 15:43:02.805868342 +0000 UTC"
                modified_by = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                version     = 1
            }
            email                = "matty.maclean@couchbase.com"
            enable_notifications = false
            expires_at           = "2024-01-04T15:43:02.805868342Z"
            id                   = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
            inactive             = false
            last_login           = "2023-10-06T15:47:57.491646422Z"
            name                 = "matty.maclean"
            organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
            organization_roles   = [
                "organizationOwner",
            ]
            region               = ""
            status               = "verified"
            time_zone            = ""
        },
    ]
    organization_id = "93f13778-3d11-43c5-861f-417a4b00ba81"
}
```
### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
data.capella_users.existing_users
capella_user.new_user
```
### Remove the resource `new_user` from the Terraform State file

Command: `terraform state rm capella_user.new_user`

Sample Output:
```
Removed capella_user.new_user
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.
### Now, let's import the resource in Terraform
Command: `terraform import capella_user.new_user id=<user_id>,organization_id=<organization_id>`

In this case, the complete command is:
``
Sample Output:
```
var.user_email
  Email address of the user

  Enter a value: matty.maclean+2@couchbase.com

var.user_name
  Name of the user

  Enter a value: matty

capella_user.new_user: Importing from ID "id=47c321f7-571c-46bb-ac1f-146aa5aec314,organization_id=93f13778-3d11-43c5-861f-417a4b00ba81"...
capella_user.new_user: Import prepared!
  Prepared capella_user for import
data.capella_users.existing_users: Reading...
capella_user.new_user: Refreshing state... [id=id=47c321f7-571c-46bb-ac1f-146aa5aec314,organization_id=93f13778-3d11-43c5-861f-417a4b00ba81]
data.capella_users.existing_users: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```
### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted
Command: `terraform plan`

Sample Output:
```
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/mattymaclean/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_users.existing_users: Reading...
capella_user.new_user: Refreshing state... [id=47c321f7-571c-46bb-ac1f-146aa5aec314]
data.capella_users.existing_users: Read complete after 0s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

### Finally, destroy the resources created by Terraform
Command: `terraform destroy`
Sample Output:
```
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/mattymaclean/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_users.existing_users: Reading...
capella_user.new_user: Refreshing state... [id=47c321f7-571c-46bb-ac1f-146aa5aec314]
data.capella_users.existing_users: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_user.new_user will be destroyed
  - resource "capella_user" "new_user" {
      - audit                = {
          - created_at  = "2023-10-06 16:19:55.734127171 +0000 UTC" -> null
          - created_by  = "47c321f7-571c-46bb-ac1f-146aa5aec314" -> null
          - modified_at = "2023-10-06 16:19:55.734127171 +0000 UTC" -> null
          - modified_by = "47c321f7-571c-46bb-ac1f-146aa5aec314" -> null
          - version     = 1 -> null
        } -> null
      - email                = "matty.maclean+2@couchbase.com" -> null
      - enable_notifications = false -> null
      - expires_at           = "2024-01-04T16:19:55.734127296Z" -> null
      - id                   = "47c321f7-571c-46bb-ac1f-146aa5aec314" -> null
      - inactive             = true -> null
      - name                 = "Matty" -> null
      - organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81" -> null
      - organization_roles   = [
          - "projectCreator",
        ] -> null
      - resources            = [
          - {
              - id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538" -> null
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
  - new_user   = {
      - audit                = {
          - created_at  = "2023-10-06 16:19:55.734127171 +0000 UTC"
          - created_by  = "47c321f7-571c-46bb-ac1f-146aa5aec314"
          - modified_at = "2023-10-06 16:19:55.734127171 +0000 UTC"
          - modified_by = "47c321f7-571c-46bb-ac1f-146aa5aec314"
          - version     = 1
        }
      - email                = "matty.maclean+2@couchbase.com"
      - enable_notifications = false
      - expires_at           = "2024-01-04T16:19:55.734127296Z"
      - id                   = "47c321f7-571c-46bb-ac1f-146aa5aec314"
      - inactive             = true
      - last_login           = ""
      - name                 = "Matty"
      - organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
      - organization_roles   = [
          - "projectCreator",
        ]
      - region               = ""
      - resources            = [
          - {
              - id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
              - roles = [
                  - "projectDataReaderWriter",
                ]
              - type  = "project"
            },
        ]
      - status               = "not-verified"
      - time_zone            = ""
    } -> null
  - users_list = {
      - data            = [
          - {
              - audit                = {
                  - created_at  = "2023-10-06 15:43:02.805868342 +0000 UTC"
                  - created_by  = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                  - modified_at = "2023-10-06 15:43:02.805868342 +0000 UTC"
                  - modified_by = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
                  - version     = 1
                }
              - email                = "matty.maclean@couchbase.com"
              - enable_notifications = false
              - expires_at           = "2024-01-04T15:43:02.805868342Z"
              - id                   = "a1acd4c3-5604-4050-80a5-58d4886e75b6"
              - inactive             = false
              - last_login           = "2023-10-06T15:47:57.491646422Z"
              - name                 = "matty.maclean"
              - organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
              - organization_roles   = [
                  - "organizationOwner",
                ]
              - region               = ""
              - resources            = null
              - status               = "verified"
              - time_zone            = ""
            },
          - {
              - audit                = {
                  - created_at  = "2023-10-06 16:19:55.734127171 +0000 UTC"
                  - created_by  = "47c321f7-571c-46bb-ac1f-146aa5aec314"
                  - modified_at = "2023-10-06 16:19:55.734127171 +0000 UTC"
                  - modified_by = "47c321f7-571c-46bb-ac1f-146aa5aec314"
                  - version     = 1
                }
              - email                = "matty.maclean+2@couchbase.com"
              - enable_notifications = false
              - expires_at           = "2024-01-04T16:19:55.734127296Z"
              - id                   = "47c321f7-571c-46bb-ac1f-146aa5aec314"
              - inactive             = true
              - last_login           = ""
              - name                 = "Matty"
              - organization_id      = "93f13778-3d11-43c5-861f-417a4b00ba81"
              - organization_roles   = [
                  - "projectCreator",
                ]
              - region               = ""
              - resources            = [
                  - {
                      - id    = "81f7bd87-6e62-4c7f-9a7e-be231c74b538"
                      - roles = [
                          - "projectDataReaderWriter",
                        ]
                      - type  = "project"
                    },
                ]
              - status               = "not-verified"
              - time_zone            = ""
            },
        ]
      - organization_id = "93f13778-3d11-43c5-861f-417a4b00ba81"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_user.new_user: Destroying... [id=47c321f7-571c-46bb-ac1f-146aa5aec314]
capella_user.new_user: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
