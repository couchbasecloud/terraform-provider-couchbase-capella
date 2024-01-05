# Capella Users Example

This example shows how to create and manage Users in Capella.

This creates a new user in the selected Capella cluster and lists existing users in the organization. It uses the organization ID to create and list users.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new user in Capella as stated in the `create_user.tf` file.
2. UPDATE: Update the user configuration using Terraform.
3. LIST: List existing users in Capella as stated in the `list_users.tf` file.
4. IMPORT: Import a user that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created user from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_users.existing_users: Reading...
data.capella_users.existing_users: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John"
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John"
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
                  + created_at  = "2022-05-24 15:48:58.579952813 +0000 UTC"
                  + created_by  = "05bb85da-9a26-4f53-818f-23b975dbb250"
                  + modified_at = "2022-05-24 15:48:58.579952813 +0000 UTC"
                  + modified_by = "05bb85da-9a26-4f53-818f-23b975dbb250"
                  + version     = 1
                }
              + email                = "taylor.swift@couchbase.com"
              + enable_notifications = false
              + expires_at           = "2023-11-12T21:41:25.105558801Z"
              + id                   = "05bb85da-9a26-4f53-818f-23b975dbb250"
              + inactive             = true
              + last_login           = "2023-09-26T22:07:25.867987345Z"
              + name                 = "Taylor Swift"
              + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + region               = ""
              + status               = "verified"
              + time_zone            = ""
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new User

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_users.existing_users: Reading...
data.capella_users.existing_users: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John"
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John"
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
                  + created_at  = "2023-03-01 05:19:26.63541879 +0000 UTC"
                  + created_by  = "17339af7-65b7-42e8-bd47-da749ab2d9bf"
                  + modified_at = "2023-03-01 05:19:26.63541879 +0000 UTC"
                  + modified_by = "17339af7-65b7-42e8-bd47-da749ab2d9bf"
                  + version     = 1
                }
              + email                = "taylor.swift@couchbase.com"
              + enable_notifications = false
              + expires_at           = "2023-12-13T11:12:31.487192911Z"
              + id                   = "17339af7-65b7-42e8-bd47-da749ab2d9bf"
              + inactive             = false
              + last_login           = "2023-09-14T11:12:39.371570481Z"
              + name                 = "Taylor Swift"
              + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + region               = ""
              + status               = "verified"
              + time_zone            = ""
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_user.new_user: Creating...
capella_user.new_user: Creation complete after 2s [id=ba422259-2b1b-429e-b3f1-d000475e8998]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_user = {
  "audit" = {
    "created_at" = "2023-10-04 22:17:06.703077348 +0000 UTC"
    "created_by" = "ba422259-2b1b-429e-b3f1-d000475e8998"
    "modified_at" = "2023-10-04 22:17:06.703077348 +0000 UTC"
    "modified_by" = "ba422259-2b1b-429e-b3f1-d000475e8998"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-02T22:17:06.70307822Z"
  "id" = "ba422259-2b1b-429e-b3f1-d000475e8998"
  "inactive" = true
  "last_login" = ""
  "name" = "John"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "roles" = tolist([
        "projectViewer",
      ])
      "type" = "project"
    },
    {
      "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
        "created_at" = "2023-03-01 05:19:26.63541879 +0000 UTC"
        "created_by" = "17339af7-65b7-42e8-bd47-da749ab2d9bf"
        "modified_at" = "2023-03-01 05:19:26.63541879 +0000 UTC"
        "modified_by" = "17339af7-65b7-42e8-bd47-da749ab2d9bf"
        "version" = 1
      }
      "email" = "taylor.swift@couchbase.com"
      "enable_notifications" = false
      "expires_at" = "2023-12-13T11:12:31.487192911Z"
      "id" = "17339af7-65b7-42e8-bd47-da749ab2d9bf"
      "inactive" = false
      "last_login" = "2023-09-14T11:12:39.371570481Z"
      "name" = "Taylor Swift"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "region" = ""
      "status" = "verified"
      "time_zone" = ""
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
}
```

### Note the User ID for the new User
Command: `terraform output new_user`

Sample Output:
```
$ terraform output new_user
{
  "audit" = {
    "created_at" = "2023-10-04 22:17:06.703077348 +0000 UTC"
    "created_by" = "ba422259-2b1b-429e-b3f1-d000475e8998"
    "modified_at" = "2023-10-04 22:17:06.703077348 +0000 UTC"
    "modified_by" = "ba422259-2b1b-429e-b3f1-d000475e8998"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-02T22:17:06.70307822Z"
  "id" = "ba422259-2b1b-429e-b3f1-d000475e8998"
  "inactive" = true
  "last_login" = ""
  "name" = "John"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "roles" = tolist([
        "projectViewer",
      ])
      "type" = "project"
    },
    {
      "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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

In this case, the user ID for my new user is `ba422259-2b1b-429e-b3f1-d000475e8998`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_users.existing_users
couchbase-capella_user.new_user
```

## IMPORT
### Remove the resource `new_user` from the Terraform State file

Command: `terraform state rm couchbase-capella_user.new_user`

Sample Output:
```
$ terraform state rm couchbase-capella_user.new_user
Removed couchbase-capella_user.new_user
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_user.new_user id=<user_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_user.new_user id=ba422259-2b1b-429e-b3f1-d000475e8998,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d`

Sample Output:
```
$ terraform import couchbase-capella_user.new_user id=ba422259-2b1b-429e-b3f1-d000475e8998,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d
capella_user.new_user: Importing from ID "id=ba422259-2b1b-429e-b3f1-d000475e8998,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d"...
data.capella_users.existing_users: Reading...
capella_user.new_user: Import prepared!
  Prepared capella_user for import
capella_user.new_user: Refreshing state... [id=id=ba422259-2b1b-429e-b3f1-d000475e8998,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d]
data.capella_users.existing_users: Read complete after 1s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the user ID i.e. the ID of the resource that we want to import.
The second ID is the organization ID i.e. the ID of the organization to which the user belongs.


## UPDATE
### Let us edit the terraform.tfvars file to change the user configuration settings.

Command: `terraform apply -var 'user={durability_level="majority",name="new_terraform_user"}'`

Sample Output:
```
$ terraform apply -var user_name="John Doe"
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_users.existing_users: Reading...
capella_user.new_user: Refreshing state... [id=ba422259-2b1b-429e-b3f1-d000475e8998]
data.capella_users.existing_users: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # capella_user.new_user must be replaced
-/+ resource "capella_user" "new_user" {
      ~ audit                = {
          ~ created_at  = "2023-10-04 22:17:06.703077348 +0000 UTC" -> (known after apply)
          ~ created_by  = "ba422259-2b1b-429e-b3f1-d000475e8998" -> (known after apply)
          ~ modified_at = "2023-10-04 22:17:06.703077348 +0000 UTC" -> (known after apply)
          ~ modified_by = "ba422259-2b1b-429e-b3f1-d000475e8998" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
      ~ enable_notifications = false -> (known after apply) # forces replacement
      ~ expires_at           = "2024-01-02T22:17:06.70307822Z" -> (known after apply) # forces replacement
      ~ id                   = "ba422259-2b1b-429e-b3f1-d000475e8998" -> (known after apply)
      ~ inactive             = true -> (known after apply) # forces replacement
      + last_login           = (known after apply) # forces replacement
      ~ name                 = "John" -> "John Doe" # forces replacement
      + region               = (known after apply) # forces replacement
      ~ resources            = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81" # forces replacement
              + roles = [
                  + "projectViewer",
                ] # forces replacement
              + type  = "project" # forces replacement
            },
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81" # forces replacement
              + roles = [
                  + "projectDataReaderWriter",
                ] # forces replacement
              + type  = "project" # forces replacement
            },
        ]
      ~ status               = "not-verified" -> (known after apply) # forces replacement
      + time_zone            = (known after apply) # forces replacement
        # (3 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_user   = {
      ~ audit                = {
          - created_at  = "2023-10-04 22:17:06.703077348 +0000 UTC"
          - created_by  = "ba422259-2b1b-429e-b3f1-d000475e8998"
          - modified_at = "2023-10-04 22:17:06.703077348 +0000 UTC"
          - modified_by = "ba422259-2b1b-429e-b3f1-d000475e8998"
          - version     = 1
        } -> (known after apply)
      ~ enable_notifications = false -> (known after apply)
      ~ expires_at           = "2024-01-02T22:17:06.70307822Z" -> (known after apply)
      ~ id                   = "ba422259-2b1b-429e-b3f1-d000475e8998" -> (known after apply)
      ~ inactive             = true -> (known after apply)
      ~ last_login           = "" -> (known after apply)
      ~ name                 = "John" -> "John Doe"
      ~ region               = "" -> (known after apply)
      ~ resources            = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      ~ status               = "not-verified" -> (known after apply)
      ~ time_zone            = "" -> (known after apply)
        # (3 unchanged elements hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_user.new_user: Destroying... [id=ba422259-2b1b-429e-b3f1-d000475e8998]
capella_user.new_user: Destruction complete after 1s
capella_user.new_user: Creating...
capella_user.new_user: Creation complete after 0s [id=8ef7377a-3030-4a82-92cc-5c5210f9d718]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

new_user = {
  "audit" = {
    "created_at" = "2023-10-04 22:24:31.945169501 +0000 UTC"
    "created_by" = "8ef7377a-3030-4a82-92cc-5c5210f9d718"
    "modified_at" = "2023-10-04 22:24:31.945169501 +0000 UTC"
    "modified_by" = "8ef7377a-3030-4a82-92cc-5c5210f9d718"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-02T22:24:31.945170096Z"
  "id" = "8ef7377a-3030-4a82-92cc-5c5210f9d718"
  "inactive" = true
  "last_login" = ""
  "name" = "John Doe"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "roles" = tolist([
        "projectViewer",
      ])
      "type" = "project"
    },
    {
      "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_user.new_user: Destroying... [id=8ef7377a-3030-4a82-92cc-5c5210f9d718]
capella_user.new_user: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
