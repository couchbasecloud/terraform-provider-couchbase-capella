# Capella Api Key Example

This example shows how to create and manage Api Key in Capella.

This creates a new api key in the organization. It uses the organization ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.
1. Create a new api key with the specified configuration.


### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/aniketkumar/.gvm/pkgsets/go1.19/global/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.42.0/23",
        ]
      + audit              = (known after apply)
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "New Terraform Api Key"
      + organization_id    = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectManager",
                  + "projectDataReader",
                ]
              + type  = "project"
            },
        ]
      + token              = (sensitive value)
    }

  # capella_project.existing_project will be created
  + resource "capella_project" "existing_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_project = {
      + description     = "A Capella Project that will host many Capella clusters."
      + if_match        = null
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
    }
  + new_apikey       = (sensitive value)

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```

### Apply the Plan, in order to create a new Api Key in Capella

Command: `terraform apply`

Sample Output:
```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/aniketkumar/.gvm/pkgsets/go1.19/global/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.42.0/23",
        ]
      + audit              = (known after apply)
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "New Terraform Api Key"
      + organization_id    = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectManager",
                  + "projectDataReader",
                ]
              + type  = "project"
            },
        ]
      + token              = (sensitive value)
    }

  # capella_project.existing_project will be created
  + resource "capella_project" "existing_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_project = {
      + description     = "A Capella Project that will host many Capella clusters."
      + if_match        = null
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
    }
  + new_apikey       = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_project.existing_project: Creating...
capella_project.existing_project: Creation complete after 1s [id=9f117f7d-f333-4f56-8390-6895819fccd5]
capella_apikey.new_apikey: Creating...
capella_apikey.new_apikey: Creation complete after 1s [id=fkqycwTdfWwO1torAKXi3nyNAmA1jPDJ]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

existing_project = {
  "audit" = {
    "created_at" = "2023-09-29 14:04:35.085325792 +0000 UTC"
    "created_by" = "IRLp8qQwHiF4Ni3IIblH0nPBa4ox0p8I"
    "modified_at" = "2023-09-29 14:04:35.085338407 +0000 UTC"
    "modified_by" = "IRLp8qQwHiF4Ni3IIblH0nPBa4ox0p8I"
    "version" = 1
  }
  "description" = "A Capella Project that will host many Capella clusters."
  "etag" = "Version: 1"
  "id" = "9f117f7d-f333-4f56-8390-6895819fccd5"
  "if_match" = tostring(null)
  "name" = "terraform-couchbasecapella-project"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
}
new_apikey = <sensitive>

```

