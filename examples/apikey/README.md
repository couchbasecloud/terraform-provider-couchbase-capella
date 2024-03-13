# Capella API Keys Example

This example shows how to create and manage API Keys in Capella.

This creates a new apikey in the selected Capella cluster and lists existing apikeys in the cluster. It uses the organization ID to create and list apikeys.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new apikey in Capella as stated in the `create_apikey.tf` file.
2. UPDATE: Update the apikey configuration using Terraform.
3. LIST: List existing apikeys in Capella as stated in the `list_apikeys.tf` file.
4. IMPORT: Import a apikey that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created apikey from Capella.

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
data.capella_apikeys.existing_apikeys: Reading...
data.capella_apikeys.existing_apikeys: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.42.0/23",
        ]
      + audit              = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "New Terraform Api Key"
      + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectManager",
                  + "projectDataReader",
                ]
              + type  = "project"
            },
        ]
      + rotate             = (known after apply)
      + secret             = (sensitive value)
      + token              = (sensitive value)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikeys_list = {
      + data            = [
          + {
              + allowed_cidrs      = [
                  + "0.0.0.0/0",
                ]
              + audit              = {
                  + created_at  = "2023-09-24 22:28:31.406184875 +0000 UTC"
                  + created_by  = "848f9d26-c94f-4149-a43e-325e38a76ef7"
                  + modified_at = "2023-09-24 22:28:31.406184875 +0000 UTC"
                  + modified_by = "848f9d26-c94f-4149-a43e-325e38a76ef7"
                  + version     = 1
                }
              + description        = "test key for terraform provider"
              + expiry             = -1
              + id                 = "892OxHFxMCMHHdBgCBX2YR6gOC7dG2Ah"
              + name               = "priyar-terraform-test"
              + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + organization_roles = [
                  + "organizationOwner",
                ]
              + resources          = null
            },
          + {
              + allowed_cidrs      = [
                  + "64.124.71.194/32",
                ]
              + audit              = {
                  + created_at  = "2023-10-04 19:16:34.399671119 +0000 UTC"
                  + created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + modified_at = "2023-10-04 19:16:34.399671119 +0000 UTC"
                  + modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + version     = 1
                }
              + description        = ""
              + expiry             = -1
              + id                 = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
              + name               = "my-terraform-key"
              + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + organization_roles = [
                  + "organizationOwner",
                ]
              + resources          = null
            },
          + {
              + allowed_cidrs      = [
                  + "0.0.0.0/0",
                ]
              + audit              = {
                  + created_at  = "2023-09-28 19:40:21.054898215 +0000 UTC"
                  + created_by  = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
                  + modified_at = "2023-09-28 19:40:21.054898215 +0000 UTC"
                  + modified_by = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
                  + version     = 1
                }
              + description        = ""
              + expiry             = 180
              + id                 = "b8OpNV4554YeZve6SQHEGs2Nks8FykRk"
              + name               = "test-key"
              + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + organization_roles = [
                  + "organizationOwner",
                ]
              + resources          = null
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }
  + new_apikey   = (sensitive value)

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```


### Apply the Plan, in order to create a new API Key

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
data.capella_apikeys.existing_apikeys: Reading...
data.capella_apikeys.existing_apikeys: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.42.0/23",
        ]
      + audit              = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "New Terraform Api Key"
      + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + roles = [
                  + "projectManager",
                  + "projectDataReader",
                ]
              + type  = "project"
            },
        ]
      + rotate             = (known after apply)
      + secret             = (sensitive value)
      + token              = (sensitive value)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikeys_list = {
      + data            = [
          + {
              + allowed_cidrs      = [
                  + "0.0.0.0/0",
                ]
              + audit              = {
                  + created_at  = "2023-09-24 22:28:31.406184875 +0000 UTC"
                  + created_by  = "848f9d26-c94f-4149-a43e-325e38a76ef7"
                  + modified_at = "2023-09-24 22:28:31.406184875 +0000 UTC"
                  + modified_by = "848f9d26-c94f-4149-a43e-325e38a76ef7"
                  + version     = 1
                }
              + description        = "test key for terraform provider"
              + expiry             = -1
              + id                 = "892OxHFxMCMHHdBgCBX2YR6gOC7dG2Ah"
              + name               = "priyar-terraform-test"
              + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + organization_roles = [
                  + "organizationOwner",
                ]
              + resources          = null
            },
          + {
              + allowed_cidrs      = [
                  + "64.124.71.194/32",
                ]
              + audit              = {
                  + created_at  = "2023-10-04 19:16:34.399671119 +0000 UTC"
                  + created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + modified_at = "2023-10-04 19:16:34.399671119 +0000 UTC"
                  + modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + version     = 1
                }
              + description        = ""
              + expiry             = -1
              + id                 = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
              + name               = "my-terraform-key"
              + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + organization_roles = [
                  + "organizationOwner",
                ]
              + resources          = null
            },
          + {
              + allowed_cidrs      = [
                  + "0.0.0.0/0",
                ]
              + audit              = {
                  + created_at  = "2023-09-28 19:40:21.054898215 +0000 UTC"
                  + created_by  = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
                  + modified_at = "2023-09-28 19:40:21.054898215 +0000 UTC"
                  + modified_by = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
                  + version     = 1
                }
              + description        = ""
              + expiry             = 180
              + id                 = "b8OpNV4554YeZve6SQHEGs2Nks8FykRk"
              + name               = "test-key"
              + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + organization_roles = [
                  + "organizationOwner",
                ]
              + resources          = null
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }
  + new_apikey   = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_apikey.new_apikey: Creating...
capella_apikey.new_apikey: Creation complete after 1s [id=8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

apikeys_list = {
  "data" = tolist([
    {
      "allowed_cidrs" = tolist([
        "0.0.0.0/0",
      ])
      "audit" = {
        "created_at" = "2023-09-24 22:28:31.406184875 +0000 UTC"
        "created_by" = "848f9d26-c94f-4149-a43e-325e38a76ef7"
        "modified_at" = "2023-09-24 22:28:31.406184875 +0000 UTC"
        "modified_by" = "848f9d26-c94f-4149-a43e-325e38a76ef7"
        "version" = 1
      }
      "description" = "test key for terraform provider"
      "expiry" = -1
      "id" = "892OxHFxMCMHHdBgCBX2YR6gOC7dG2Ah"
      "name" = "priyar-terraform-test"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "organization_roles" = tolist([
        "organizationOwner",
      ])
      "resources" = tolist(null) /* of object */
    },
    {
      "allowed_cidrs" = tolist([
        "64.124.71.194/32",
      ])
      "audit" = {
        "created_at" = "2023-10-04 19:16:34.399671119 +0000 UTC"
        "created_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "modified_at" = "2023-10-04 19:16:34.399671119 +0000 UTC"
        "modified_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "version" = 1
      }
      "description" = ""
      "expiry" = -1
      "id" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
      "name" = "my-terraform-key"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "organization_roles" = tolist([
        "organizationOwner",
      ])
      "resources" = tolist(null) /* of object */
    },
    {
      "allowed_cidrs" = tolist([
        "0.0.0.0/0",
      ])
      "audit" = {
        "created_at" = "2023-09-28 19:40:21.054898215 +0000 UTC"
        "created_by" = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
        "modified_at" = "2023-09-28 19:40:21.054898215 +0000 UTC"
        "modified_by" = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
        "version" = 1
      }
      "description" = ""
      "expiry" = 180
      "id" = "b8OpNV4554YeZve6SQHEGs2Nks8FykRk"
      "name" = "test-key"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "organization_roles" = tolist([
        "organizationOwner",
      ])
      "resources" = tolist(null) /* of object */
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
}
new_apikey = <sensitive>
```


### Note the API Key ID for the new API Key
Command: `terraform output new_apikey`

Sample Output:
```
$ terraform output new_apikey
{
  "allowed_cidrs" = tolist([
    "10.1.42.0/23",
    "10.1.42.0/23",
  ])
  "audit" = {
    "created_at" = "2023-10-04 20:19:41.409153281 +0000 UTC"
    "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "modified_at" = "2023-10-04 20:19:41.409153281 +0000 UTC"
    "modified_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "version" = 1
  }
  "description" = ""
  "expiry" = 180
  "id" = "8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH"
  "name" = "New Terraform Api Key"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "resources" = tolist([
    {
      "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "roles" = tolist([
        "projectManager",
        "projectDataReader",
      ])
      "type" = "project"
    },
  ])
  "rotate" = tonumber(null)
  "secret" = tostring(null)
  "token" = "OEtFSjlPRFVRSGNwc2tQNWFuQk9GVHpWOFpleldKYUg6dVpIUElGTldxdyNmOFZmU0tPSXNLRCE2YTgwUjRoSmt1UzQ0cmltMnhXOW5nNXQleG9hZjUleWd1b05WdVF4Qw=="
}
```


In this case, the apikey ID for my new apikey is `bmV3X3RlcnJhZm9ybV9idWNrZXQ=`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_apikeys.existing_apikeys
couchbase-capella_apikey.new_apikey
```

## IMPORT
### Remove the resource `new_apikey` from the Terraform State file

Command: `terraform state rm couchbase-capella_apikey.new_apikey`

Sample Output:
```
$ terraform state rm couchbase-capella_apikey.new_apikey
Removed capella_apikey.new_apikey
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_apikey.new_apikey id=<apikey_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_apikey.new_apikey id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d`

Sample Output:
```
$ terraform import couchbase-capella_apikey.new_apikey id=8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d
capella_apikey.new_apikey: Importing from ID "id=8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d"...
data.capella_apikeys.existing_apikeys: Reading...
capella_apikey.new_apikey: Import prepared!
  Prepared capella_apikey for import
capella_apikey.new_apikey: Refreshing state... [id=id=8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d]
data.capella_apikeys.existing_apikeys: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the apikey ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the apikey belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

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
data.capella_apikeys.existing_apikeys: Reading...
capella_apikey.new_apikey: Refreshing state... [id=8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH]
data.capella_apikeys.existing_apikeys: Read complete after 0s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to change the apikey configuration settings.

Command: `terraform apply -var 'apikey={allowed_cidrs=["10.1.42.0/23", "10.1.43.0/23", "10.1.44.0/23"], name="New Terraform Api Key", description="A Capella Api Key", organization_roles = ["organizationMember"], expiry = 179}'`

Sample Output:
```
$ terraform apply -var 'apikey={allowed_cidrs=["10.1.42.0/23", "10.1.43.0/23", "10.1.44.0/23"], name="New Terraform Api Key", description="A Capella Api Key", organization_roles = ["organizationMember"], expiry             = 179}'
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_apikeys.existing_apikeys: Reading...
capella_apikey.new_apikey: Refreshing state... [id=8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH]
data.capella_apikeys.existing_apikeys: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # capella_apikey.new_apikey must be replaced
-/+ resource "capella_apikey" "new_apikey" {
      ~ allowed_cidrs      = [ # forces replacement
            "10.1.42.0/23",
          - "10.1.42.0/23",
          + "10.1.43.0/23",
          + "10.1.44.0/23",
        ]
      ~ audit              = {
          ~ created_at  = "2023-10-04 20:19:41.409153281 +0000 UTC" -> (known after apply)
          ~ created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> (known after apply)
          ~ modified_at = "2023-10-04 20:19:41.409153281 +0000 UTC" -> (known after apply)
          ~ modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
      + description        = (known after apply)
      ~ expiry             = 180 -> (known after apply)
      ~ id                 = "8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH" -> (known after apply)
        name               = "New Terraform Api Key"
      + rotate             = (known after apply)
      + secret             = (sensitive value)
      + token              = (sensitive value)
        # (3 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_apikey = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_apikey.new_apikey: Destroying... [id=8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH]
capella_apikey.new_apikey: Destruction complete after 0s
capella_apikey.new_apikey: Creating...
capella_apikey.new_apikey: Creation complete after 0s [id=UA72qZ7ddRop6jkM0sTG8UnotBMBmBEW]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

apikeys_list = {
  "data" = tolist([
    {
      "allowed_cidrs" = tolist([
        "0.0.0.0/0",
      ])
      "audit" = {
        "created_at" = "2023-09-24 22:28:31.406184875 +0000 UTC"
        "created_by" = "848f9d26-c94f-4149-a43e-325e38a76ef7"
        "modified_at" = "2023-09-24 22:28:31.406184875 +0000 UTC"
        "modified_by" = "848f9d26-c94f-4149-a43e-325e38a76ef7"
        "version" = 1
      }
      "description" = "test key for terraform provider"
      "expiry" = -1
      "id" = "892OxHFxMCMHHdBgCBX2YR6gOC7dG2Ah"
      "name" = "priyar-terraform-test"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "organization_roles" = tolist([
        "organizationOwner",
      ])
      "resources" = tolist(null) /* of object */
    },
    {
      "allowed_cidrs" = tolist([
        "10.1.42.0/23",
        "10.1.42.0/23",
      ])
      "audit" = {
        "created_at" = "2023-10-04 20:19:41.409153281 +0000 UTC"
        "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
        "modified_at" = "2023-10-04 20:19:41.409153281 +0000 UTC"
        "modified_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
        "version" = 1
      }
      "description" = ""
      "expiry" = 180
      "id" = "8KEJ9ODUQHcpskP5anBOFTzV8ZezWJaH"
      "name" = "New Terraform Api Key"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "organization_roles" = tolist([
        "organizationMember",
      ])
      "resources" = tolist([
        {
          "id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
          "roles" = tolist([
            "projectManager",
            "projectDataReader",
          ])
          "type" = "project"
        },
      ])
    },
    {
      "allowed_cidrs" = tolist([
        "64.124.71.194/32",
      ])
      "audit" = {
        "created_at" = "2023-10-04 19:16:34.399671119 +0000 UTC"
        "created_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "modified_at" = "2023-10-04 19:16:34.399671119 +0000 UTC"
        "modified_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "version" = 1
      }
      "description" = ""
      "expiry" = -1
      "id" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
      "name" = "my-terraform-key"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "organization_roles" = tolist([
        "organizationOwner",
      ])
      "resources" = tolist(null) /* of object */
    },
    {
      "allowed_cidrs" = tolist([
        "0.0.0.0/0",
      ])
      "audit" = {
        "created_at" = "2023-09-28 19:40:21.054898215 +0000 UTC"
        "created_by" = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
        "modified_at" = "2023-09-28 19:40:21.054898215 +0000 UTC"
        "modified_by" = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
        "version" = 1
      }
      "description" = ""
      "expiry" = 180
      "id" = "b8OpNV4554YeZve6SQHEGs2Nks8FykRk"
      "name" = "test-key"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "organization_roles" = tolist([
        "organizationOwner",
      ])
      "resources" = tolist(null) /* of object */
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
}
new_apikey = <sensitive>
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_apikeys.existing_apikeys: Reading...
capella_apikey.new_apikey: Refreshing state... [id=UA72qZ7ddRop6jkM0sTG8UnotBMBmBEW]
data.capella_apikeys.existing_apikeys: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_apikey.new_apikey will be destroyed
  - resource "capella_apikey" "new_apikey" {
      - allowed_cidrs      = [
          - "10.1.42.0/23",
          - "10.1.43.0/23",
          - "10.1.44.0/23",
        ] -> null
      - audit              = {
          - created_at  = "2023-10-04 20:32:13.25431423 +0000 UTC" -> null
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - modified_at = "2023-10-04 20:32:13.25431423 +0000 UTC" -> null
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - version     = 1 -> null
        }
      - expiry             = 180 -> null
      - id                 = "UA72qZ7ddRop6jkM0sTG8UnotBMBmBEW" -> null
      - name               = "New Terraform Api Key" -> null
      - organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - organization_roles = [
          - "organizationMember",
        ] -> null
      - resources          = [
          - {
              - id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81" -> null
              - roles = [
                  - "projectManager",
                  - "projectDataReader",
                ] -> null
              - type  = "project" -> null
            },
        ]
      - token              = (sensitive value)
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - apikeys_list = {
      - data            = [
          - {
              - allowed_cidrs      = [
                  - "0.0.0.0/0",
                ]
              - audit              = {
                  - created_at  = "2023-09-24 22:28:31.406184875 +0000 UTC"
                  - created_by  = "848f9d26-c94f-4149-a43e-325e38a76ef7"
                  - modified_at = "2023-09-24 22:28:31.406184875 +0000 UTC"
                  - modified_by = "848f9d26-c94f-4149-a43e-325e38a76ef7"
                  - version     = 1
                }
              - description        = "test key for terraform provider"
              - expiry             = -1
              - id                 = "892OxHFxMCMHHdBgCBX2YR6gOC7dG2Ah"
              - name               = "priyar-terraform-test"
              - organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - organization_roles = [
                  - "organizationOwner",
                ]
              - resources          = null
            },
          - {
              - allowed_cidrs      = [
                  - "64.124.71.194/32",
                ]
              - audit              = {
                  - created_at  = "2023-10-04 19:16:34.399671119 +0000 UTC"
                  - created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - modified_at = "2023-10-04 19:16:34.399671119 +0000 UTC"
                  - modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - version     = 1
                }
              - description        = ""
              - expiry             = -1
              - id                 = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
              - name               = "my-terraform-key"
              - organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - organization_roles = [
                  - "organizationOwner",
                ]
              - resources          = null
            },
          - {
              - allowed_cidrs      = [
                  - "10.1.42.0/23",
                  - "10.1.43.0/23",
                  - "10.1.44.0/23",
                ]
              - audit              = {
                  - created_at  = "2023-10-04 20:32:13.25431423 +0000 UTC"
                  - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
                  - modified_at = "2023-10-04 20:32:13.25431423 +0000 UTC"
                  - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
                  - version     = 1
                }
              - description        = ""
              - expiry             = 180
              - id                 = "UA72qZ7ddRop6jkM0sTG8UnotBMBmBEW"
              - name               = "New Terraform Api Key"
              - organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - organization_roles = [
                  - "organizationMember",
                ]
              - resources          = [
                  - {
                      - id    = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
                      - roles = [
                          - "projectManager",
                          - "projectDataReader",
                        ]
                      - type  = "project"
                    },
                ]
            },
          - {
              - allowed_cidrs      = [
                  - "0.0.0.0/0",
                ]
              - audit              = {
                  - created_at  = "2023-09-28 19:40:21.054898215 +0000 UTC"
                  - created_by  = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
                  - modified_at = "2023-09-28 19:40:21.054898215 +0000 UTC"
                  - modified_by = "297fd203-d4ec-4f33-8ea1-76e736ed65b1"
                  - version     = 1
                }
              - description        = ""
              - expiry             = 180
              - id                 = "b8OpNV4554YeZve6SQHEGs2Nks8FykRk"
              - name               = "test-key"
              - organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - organization_roles = [
                  - "organizationOwner",
                ]
              - resources          = null
            },
        ]
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    } -> null
  - new_apikey   = (sensitive value)

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_apikey.new_apikey: Destroying... [id=UA72qZ7ddRop6jkM0sTG8UnotBMBmBEW]
capella_apikey.new_apikey: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
