# Capella AllowList Example

This example shows how to create and manage Projects in Capella.

This creates a new allowlist in the selected Capella cluster. It uses the organization ID, projectId and clusterId to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Create a new allowlist entry in an existing Capella cluster as stated in the `create_allowlist.tf` file.

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_allowlist.new_allowlist will be created
  + resource "capella_allowlist" "new_allowlist" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "f3818c88-3016-4c01-b3db-233173d8e4fd"
      + comment         = "Allow access from any ip address"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
      + project_id      = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + example_allowlist = {
      + audit           = (known after apply)
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "f3818c88-3016-4c01-b3db-233173d8e4fd"
      + comment         = "Allow access from any ip address"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + if_match        = null
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
      + project_id      = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
    }

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Allowlist entry

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_allowlist.new_allowlist will be created
  + resource "capella_allowlist" "new_allowlist" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "f3818c88-3016-4c01-b3db-233173d8e4fd"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
      + project_id      = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_allowlist = {
      + audit           = (known after apply)
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "f3818c88-3016-4c01-b3db-233173d8e4fd"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + if_match        = null
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
      + project_id      = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_allowlist.new_allowlist: Creating...
capella_allowlist.new_allowlist: Creation complete after 1s [id=08b1221f-33cf-42cd-a4d5-a35f6aa0763e]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_allowlist = {
  "audit" = {
    "created_at" = "2023-09-19 21:57:04.032017652 +0000 UTC"
    "created_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
    "modified_at" = "2023-09-19 21:57:04.032017652 +0000 UTC"
    "modified_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
    "version" = 1
  }
  "cidr" = "10.0.0.0/16"
  "cluster_id" = "f3818c88-3016-4c01-b3db-233173d8e4fd"
  "comment" = "Allow access from another VPC"
  "expires_at" = "2023-11-14T21:49:58.465Z"
  "id" = "08b1221f-33cf-42cd-a4d5-a35f6aa0763e"
  "if_match" = tostring(null)
  "organization_id" = "bdb8662c-7157-46ea-956f-ed86f4c75211"
  "project_id" = "e912ed02-8ac4-403c-a0c5-67c57284a5a4"
}
```

