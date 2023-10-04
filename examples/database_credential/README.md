# Capella Database Credentials Example

This example shows how to create and manage Database Credentials in Capella.

This creates a new database credential in the selected Capella organization. It uses the organization ID, project ID and cluster ID.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Create a new database credential in Capella as stated in the `create_database_credential.tf` file.
2. View the sensitive field i.e. database credential password after creation.
3. Update the database credential password.

If you check the `terraform.template.tfvars` file - you can see that we need 7 main variables to run the terraform commands.
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
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_database_credential.new_database_credential will be created
  + resource "capella_database_credential" "new_database_credential" {
      + access          = [
          + {
              + privileges = [
                  + "data_reader",
                  + "data_writer",
                ]
            },
        ]
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cluster_id      = "c082af14-c244-40da-b54a-669392738569"
      + id              = (known after apply)
      + name            = "test_db_user"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + password        = (sensitive value)
      + project_id      = "a1d1a971-092e-40d9-a68b-ef705573f3d8"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_database_credential = (sensitive value)

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

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
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_database_credential.new_database_credential will be created
  + resource "capella_database_credential" "new_database_credential" {
      + access          = [
          + {
              + privileges = [
                  + "data_reader",
                  + "data_writer",
                ]
            },
        ]
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cluster_id      = "c082af14-c244-40da-b54a-669392738569"
      + id              = (known after apply)
      + name            = "test_db_user"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + password        = (sensitive value)
      + project_id      = "a1d1a971-092e-40d9-a68b-ef705573f3d8"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_database_credential = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_database_credential.new_database_credential: Creating...
capella_database_credential.new_database_credential: Creation complete after 2s [id=7ef4675e-513f-4358-a583-ae5c23e6fa67]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_database_credential = <sensitive>
```

### View the create db_credential Password
Command: `terraform output new_database_credential`

Sample Output:
```
$ terraform output new_database_credentials
╷
│ Error: Output "new_database_credentials" not found
│ 
│ The output variable requested could not be found in the state file. If you recently added this to your configuration, be sure to run `terraform apply`, since the state
│ won't be updated with new output variables until that command is run.
╵
macos:database_credential talina.shrotriya$ terraform output new_database_credential
{
  "access" = tolist([
    {
      "privileges" = tolist([
        "data_reader",
        "data_writer",
      ])
    },
  ])
  "audit" = {
    "created_at" = "2023-09-28 23:03:39.742677746 +0000 UTC"
    "created_by" = "wTQ5WXpeWsNpfXTVOIz12FzqH8Ye7m2p"
    "modified_at" = "2023-09-28 23:03:39.742677746 +0000 UTC"
    "modified_by" = "wTQ5WXpeWsNpfXTVOIz12FzqH8Ye7m2p"
    "version" = 1
  }
  "cluster_id" = "c082af14-c244-40da-b54a-669392738569"
  "id" = "7ef4675e-513f-4358-a583-ae5c23e6fa67"
  "name" = "test_db_user"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "password" = "Secret12$#"
  "project_id" = "a1d1a971-092e-40d9-a68b-ef705573f3d8"
}
```

### Update the database credential password
- Change the password in the terraform.tfvars file.
- Execute terraform plan

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
capella_database_credential.new_database_credential: Refreshing state... [id=1a92d0cf-6c41-481f-ad10-c843bd7837f1]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_database_credential.new_database_credential will be updated in-place
  ~ resource "capella_database_credential" "new_database_credential" {
      ~ access          = [
          ~ {
              ~ privileges = [
                    "data_reader",
                  - "data_writer",
                ]
            },
        ]
      ~ audit           = {
          ~ created_at  = "2023-10-03 01:12:14.215211005 +0000 UTC" -> (known after apply)
          ~ created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ modified_at = "2023-10-03 01:12:14.215211005 +0000 UTC" -> (known after apply)
          ~ modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
        id              = "1a92d0cf-6c41-481f-ad10-c843bd7837f1"
        name            = "test_db_user"
      ~ password        = (sensitive value)
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_database_credential = (sensitive value)
```

- Execute terraform apply
Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
capella_database_credential.new_database_credential: Refreshing state... [id=1a92d0cf-6c41-481f-ad10-c843bd7837f1]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_database_credential.new_database_credential will be updated in-place
  ~ resource "capella_database_credential" "new_database_credential" {
      ~ access          = [
          ~ {
              ~ privileges = [
                    "data_reader",
                  - "data_writer",
                ]
            },
        ]
      ~ audit           = {
          ~ created_at  = "2023-10-03 01:12:14.215211005 +0000 UTC" -> (known after apply)
          ~ created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ modified_at = "2023-10-03 01:12:14.215211005 +0000 UTC" -> (known after apply)
          ~ modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
        id              = "1a92d0cf-6c41-481f-ad10-c843bd7837f1"
        name            = "test_db_user"
      ~ password        = (sensitive value)
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_database_credential = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_database_credential.new_database_credential: Modifying... [id=1a92d0cf-6c41-481f-ad10-c843bd7837f1]
capella_database_credential.new_database_credential: Modifications complete after 2s [id=1a92d0cf-6c41-481f-ad10-c843bd7837f1]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_database_credential = <sensitive>
```

- Finally, we can confirm if the password was updated by running the terraform output command.

```
$ terraform output new_database_credential
{
  "access" = tolist([
    {
      "privileges" = tolist([
        "data_reader",
      ])
      "resources" = null /* object */
    },
  ])
  "audit" = {
    "created_at" = "2023-10-03 01:12:14.215211005 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-03 01:12:14.215211005 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cluster_id" = "c082af14-c244-40da-b54a-669392738569"
  "id" = "1a92d0cf-6c41-481f-ad10-c843bd7837f1"
  "name" = "test_db_user"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "password" = "NewSecret12$#"
  "project_id" = "a1d1a971-092e-40d9-a68b-ef705573f3d8"
}
```