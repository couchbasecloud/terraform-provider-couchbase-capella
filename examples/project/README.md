# Capella Projects Example

This example shows how to create and manage Projects in Capella.

This creates a new project in the selected Capella organization and lists existing Projects in the organization. It uses the organization ID to create and list Projects.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new project in Capella as stated in the `create_project.tf` file.
2. UPDATE: Update the project using Terraform.
3. LIST: List existing projects in Capella as stated in the `list_projects.tf` file.
4. IMPORT: Import a project that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created project from Capella.

If you check the `terraform.template.tfvars` file - you can see that we need 3 main variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

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
data.capella_projects.existing_projects: Reading...
data.capella_projects.existing_projects: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_project.new_project will be created
  + resource "capella_project" "new_project" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_project   = {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }
  + projects_list = {
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2023-09-19 20:38:55.873822668 +0000 UTC"
                  + created_by  = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                  + modified_at = "2023-09-19 20:38:55.873836582 +0000 UTC"
                  + modified_by = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                  + version     = 1
                }
              + description     = ""
              + etag            = null
              + id              = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
              + if_match        = null
              + name            = "Tacos"
              + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
            },
        ]
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Project

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
data.capella_projects.existing_projects: Reading...
data.capella_projects.existing_projects: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_project.new_project will be created
  + resource "capella_project" "new_project" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_project   = {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }
  + projects_list = {
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2023-09-19 20:38:55.873822668 +0000 UTC"
                  + created_by  = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                  + modified_at = "2023-09-19 20:38:55.873836582 +0000 UTC"
                  + modified_by = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                  + version     = 1
                }
              + description     = ""
              + etag            = null
              + id              = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
              + if_match        = null
              + name            = "Tacos"
              + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
            },
        ]
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_project.new_project: Creating...
capella_project.new_project: Creation complete after 1s [id=95b69ba0-23f8-45bf-8640-8ea99e8860fd]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_project = {
  "audit" = {
    "created_at" = "2023-09-19 20:39:45.392955893 +0000 UTC"
    "created_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
    "modified_at" = "2023-09-19 20:39:45.392987613 +0000 UTC"
    "modified_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
    "version" = 1
  }
  "description" = "A Capella Project that will host many Capella clusters."
  "etag" = "Version: 1"
  "id" = "95b69ba0-23f8-45bf-8640-8ea99e8860fd"
  "if_match" = tostring(null)
  "name" = "terraform-couchbasecapella-project"
  "organization_id" = "bdb8662c-7157-46ea-956f-ed86f4c75211"
}
projects_list = {
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2023-09-19 20:38:55.873822668 +0000 UTC"
        "created_by" = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
        "modified_at" = "2023-09-19 20:38:55.873836582 +0000 UTC"
        "modified_by" = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
        "version" = 1
      }
      "description" = ""
      "etag" = tostring(null)
      "id" = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
      "if_match" = tostring(null)
      "name" = "Tacos"
      "organization_id" = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    },
  ])
  "organization_id" = "bdb8662c-7157-46ea-956f-ed86f4c75211"
}
```

### Note the Project ID for the new Project
Command: `terraform show`

Sample Output:
```
$ terraform show
# capella_project.new_project:
resource "capella_project" "new_project" {
    audit           = {
        created_at  = "2023-09-19 20:39:45.392955893 +0000 UTC"
        created_by  = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
        modified_at = "2023-09-19 20:39:45.392987613 +0000 UTC"
        modified_by = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
        version     = 1
    }
    description     = "A Capella Project that will host many Capella clusters."
    etag            = "Version: 1"
    id              = "95b69ba0-23f8-45bf-8640-8ea99e8860fd"
    name            = "terraform-couchbasecapella-project"
    organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
}

# data.capella_projects.existing_projects:
data "capella_projects" "existing_projects" {
    data            = [
        # (1 unchanged element hidden)
    ]
    organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
}


Outputs:

new_project = {
    audit           = {
        created_at  = "2023-09-19 20:39:45.392955893 +0000 UTC"
        created_by  = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
        modified_at = "2023-09-19 20:39:45.392987613 +0000 UTC"
        modified_by = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
        version     = 1
    }
    description     = "A Capella Project that will host many Capella clusters."
    etag            = "Version: 1"
    id              = "95b69ba0-23f8-45bf-8640-8ea99e8860fd"
    if_match        = null
    name            = "terraform-couchbasecapella-project"
    organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
}
projects_list = {
    data            = [
        {
            audit           = {
                created_at  = "2023-09-19 20:38:55.873822668 +0000 UTC"
                created_by  = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                modified_at = "2023-09-19 20:38:55.873836582 +0000 UTC"
                modified_by = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                version     = 1
            }
            description     = ""
            etag            = null
            id              = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
            if_match        = null
            name            = "Tacos"
            organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
        },
    ]
    organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
}
```

In this case, the Project ID for my new project is `95b69ba0-23f8-45bf-8640-8ea99e8860fd`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_projects.existing_projects
couchbase-capella_project.new_project
```

## IMPORT
### Remove the resource `new_project` from the Terraform State file

Command: `terraform state rm couchbase-capella_project.new_project`

Sample Output:
```
$ terraform state rm couchbase-capella_project.new_project
Removed couchbase-capella_project.new_project
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_project.new_project id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_project.new_project id=95b69ba0-23f8-45bf-8640-8ea99e8860fd,organization_id=bdb8662c-7157-46ea-956f-ed86f4c75211`

Sample Output:
```
$ terraform import couchbase-capella_project.new_project id=95b69ba0-23f8-45bf-8640-8ea99e8860fd,organization_id=bdb8662c-7157-46ea-956f-ed86f4c75211
capella_project.new_project: Importing from ID "id=95b69ba0-23f8-45bf-8640-8ea99e8860fd,organization_id=bdb8662c-7157-46ea-956f-ed86f4c75211"...
data.capella_projects.existing_projects: Reading...
capella_project.new_project: Import prepared!
  Prepared capella_project for import
capella_project.new_project: Refreshing state... [id=id=95b69ba0-23f8-45bf-8640-8ea99e8860fd,organization_id=bdb8662c-7157-46ea-956f-ed86f4c75211]
data.capella_projects.existing_projects: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the project ID i.e. the ID of the resource that we want to import
The second ID is the organization ID i.e. the ID of the organization to which the project belongs to.

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
data.capella_projects.existing_projects: Reading...
capella_project.new_project: Refreshing state... [id=95b69ba0-23f8-45bf-8640-8ea99e8860fd]
data.capella_projects.existing_projects: Read complete after 1s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to change the project name.

var.project_name = "my_edited_project_name"

Command: `terraform apply -var project_name="my_edited_project_name"`

Sample Output:
```
$ terraform apply -var project_name="my_edited_project_name"
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_projects.existing_projects: Reading...
capella_project.new_project: Refreshing state... [id=1d5f4c38-f645-4279-b7c3-5faec80dad0c]
data.capella_projects.existing_projects: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_project.new_project will be updated in-place
  ~ resource "capella_project" "new_project" {
      ~ audit           = {
          ~ created_at  = "2023-10-03 03:33:59.770847849 +0000 UTC" -> (known after apply)
          ~ created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ modified_at = "2023-10-03 03:33:59.770861069 +0000 UTC" -> (known after apply)
          ~ modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
      ~ etag            = "Version: 1" -> (known after apply)
        id              = "1d5f4c38-f645-4279-b7c3-5faec80dad0c"
      ~ name            = "terraform-couchbasecapella-project" -> "my_edited_project_name"
        # (2 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_project   = {
      ~ audit           = {
          - created_at  = "2023-10-03 03:33:59.770847849 +0000 UTC"
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - modified_at = "2023-10-03 03:33:59.770861069 +0000 UTC"
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - version     = 1
        } -> (known after apply)
      ~ etag            = "Version: 1" -> (known after apply)
        id              = "1d5f4c38-f645-4279-b7c3-5faec80dad0c"
      ~ name            = "terraform-couchbasecapella-project" -> "my_edited_project_name"
        # (3 unchanged elements hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_project.new_project: Modifying... [id=1d5f4c38-f645-4279-b7c3-5faec80dad0c]
capella_project.new_project: Modifications complete after 0s [id=1d5f4c38-f645-4279-b7c3-5faec80dad0c]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_project = {
  "audit" = {
    "created_at" = "2023-10-03 03:33:59.770847849 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-03 03:34:36.811342918 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 2
  }
  "description" = "A Capella Project that will host many Capella clusters."
  "etag" = "Version: 2"
  "id" = "1d5f4c38-f645-4279-b7c3-5faec80dad0c"
  "if_match" = tostring(null)
  "name" = "my_edited_project_name"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
}
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
data.capella_projects.existing_projects: Reading...
capella_project.new_project: Refreshing state... [id=95b69ba0-23f8-45bf-8640-8ea99e8860fd]
data.capella_projects.existing_projects: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_project.new_project will be destroyed
  - resource "capella_project" "new_project" {
      - audit           = {
          - created_at  = "2023-09-19 20:39:45.392955893 +0000 UTC" -> null
          - created_by  = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK" -> null
          - modified_at = "2023-09-19 20:39:45.392987613 +0000 UTC" -> null
          - modified_by = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK" -> null
          - version     = 1 -> null
        }
      - description     = "A Capella Project that will host many Capella clusters." -> null
      - etag            = "Version: 1" -> null
      - id              = "95b69ba0-23f8-45bf-8640-8ea99e8860fd" -> null
      - name            = "terraform-couchbasecapella-project" -> null
      - organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_project   = {
      - audit           = {
          - created_at  = "2023-09-19 20:39:45.392955893 +0000 UTC"
          - created_by  = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
          - modified_at = "2023-09-19 20:39:45.392987613 +0000 UTC"
          - modified_by = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
          - version     = 1
        }
      - description     = "A Capella Project that will host many Capella clusters."
      - etag            = "Version: 1"
      - id              = "95b69ba0-23f8-45bf-8640-8ea99e8860fd"
      - if_match        = null
      - name            = "terraform-couchbasecapella-project"
      - organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    } -> null
  - projects_list = {
      - data            = [
          - {
              - audit           = {
                  - created_at  = "2023-09-19 20:39:45.392955893 +0000 UTC"
                  - created_by  = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
                  - modified_at = "2023-09-19 20:39:45.392987613 +0000 UTC"
                  - modified_by = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
                  - version     = 1
                }
              - description     = "A Capella Project that will host many Capella clusters."
              - etag            = null
              - id              = "95b69ba0-23f8-45bf-8640-8ea99e8860fd"
              - if_match        = null
              - name            = "terraform-couchbasecapella-project"
              - organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
            },
          - {
              - audit           = {
                  - created_at  = "2023-09-19 20:38:55.873822668 +0000 UTC"
                  - created_by  = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                  - modified_at = "2023-09-19 20:38:55.873836582 +0000 UTC"
                  - modified_by = "bff4a7f5-33c0-4324-bb40-0890a01a20ae"
                  - version     = 1
                }
              - description     = ""
              - etag            = null
              - id              = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
              - if_match        = null
              - name            = "Tacos"
              - organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
            },
        ]
      - organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_project.new_project: Destroying... [id=95b69ba0-23f8-45bf-8640-8ea99e8860fd]
capella_project.new_project: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```