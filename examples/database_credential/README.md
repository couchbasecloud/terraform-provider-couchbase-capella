# Capella Database Credentials Example

This example shows how to create and manage Database Credentials in Capella.

This creates a new database credential in the selected Capella cluster and lists existing database credentials in the cluster. It uses the cluster ID to create and list database credentials.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new database credential in Capella as stated in the `create_database_credential.tf` file.
2. UPDATE: Update the database credential configuration using Terraform.
3. LIST: List existing database credentials in Capella as stated in the `list_database_credentials.tf` file.
4. IMPORT: Import a database credential that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created database credential from Capella.

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_database_credentials.existing_credentials: Reading...
data.capella_database_credentials.existing_credentials: Read complete after 2s

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
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + id              = (known after apply)
      + name            = "test_db_user"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + password        = (sensitive value)
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + database_credentials_list = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = null
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_database_credential   = (sensitive value)

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Database Credential

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
data.capella_database_credentials.existing_credentials: Reading...
data.capella_database_credentials.existing_credentials: Read complete after 2s

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
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + id              = (known after apply)
      + name            = "test_db_user"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + password        = (sensitive value)
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + database_credentials_list = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = null
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_database_credential   = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_database_credential.new_database_credential: Creating...
capella_database_credential.new_database_credential: Creation complete after 3s [id=95591a3b-7031-4257-8d9e-7c4620d14618]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

database_credentials_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist(null) /* of object */
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_database_credential = <sensitive>
```


You can see the password using this command:

Command: `terraform output new_database_credential`

Sample Output:
```
$ terraform output new_database_credential
{
  "access" = tolist([
    {
      "privileges" = tolist([
        "data_reader",
        "data_writer",
      ])
      "resources" = null /* object */
    },
  ])
  "audit" = {
    "created_at" = "2023-10-04 04:58:00.034423938 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-04 04:58:00.034423938 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "id" = "95591a3b-7031-4257-8d9e-7c4620d14618"
  "name" = "test_db_user"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "password" = "Secret12$#"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
```

### Note the Database Credential ID for the new Database Credential
Command: `terraform output new_database_credential`

In this case, the database credential ID for my new database credential is `95591a3b-7031-4257-8d9e-7c4620d14618`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_database_credentials.existing_credentials
couchbase-capella_database_credential.new_database_credential
```

## IMPORT
### Remove the resource `new_database_credential` from the Terraform State file

Command: `terraform state rm couchbase-capella_database_credential.new_database_credential`

Sample Output:
```
$ terraform state rm couchbase-capella_database_credential.new_database_credential
Removed couchbase-capella_database_credential.new_database_credential
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_database_credential.new_database_credential id=<database_credential_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_database_credential.new_database_credential id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d`

Sample Output:
```
$ terraform import capella_database_credential.new_database_credential id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d
capella_database_credential.new_database_credential: Importing from ID "id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d"...
data.capella_database_credentials.existing_credentials: Reading...
capella_database_credential.new_database_credential: Import prepared!
  Prepared capella_database_credential for import
capella_database_credential.new_database_credential: Refreshing state... [id=id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d]
data.capella_database_credentials.existing_credentials: Read complete after 2s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the bucket ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the bucket belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful, do note that the database credential will be updated after we import as the password and access result in an update

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
data.capella_database_credentials.existing_credentials: Reading...
capella_database_credential.new_database_credential: Refreshing state... [id=95591a3b-7031-4257-8d9e-7c4620d14618]
data.capella_database_credentials.existing_credentials: Read complete after 3s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_database_credential.new_database_credential will be updated in-place
  ~ resource "capella_database_credential" "new_database_credential" {
      ~ access          = [
          + {
              + privileges = [
                  + "data_reader",
                  + "data_writer",
                ]
            },
        ]
      ~ audit           = {
          ~ created_at  = "2023-10-04 04:58:00.034423938 +0000 UTC" -> (known after apply)
          ~ created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ modified_at = "2023-10-04 04:58:00.034423938 +0000 UTC" -> (known after apply)
          ~ modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
        id              = "95591a3b-7031-4257-8d9e-7c4620d14618"
        name            = "test_db_user"
      + password        = (sensitive value)
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_database_credential = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_database_credential.new_database_credential: Modifying... [id=95591a3b-7031-4257-8d9e-7c4620d14618]
capella_database_credential.new_database_credential: Modifications complete after 5s [id=95591a3b-7031-4257-8d9e-7c4620d14618]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

database_credentials_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist([
    {
      "access" = tolist([])
      "audit" = {
        "created_at" = "2023-10-04 04:58:00.034423938 +0000 UTC"
        "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
        "modified_at" = "2023-10-04 04:58:00.034423938 +0000 UTC"
        "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
        "version" = 1
      }
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "id" = "95591a3b-7031-4257-8d9e-7c4620d14618"
      "name" = "test_db_user"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_database_credential = <sensitive>
```

## UPDATE
### Let us edit the terraform.tfvars file to change the Database Credential configuration settings.

Command: `terraform apply -var 'access=[{privileges=["data_reader"]}]'`

Sample Output:
```
$ terraform apply -var 'access=[{privileges=["data_reader"]}]'
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_database_credentials.existing_credentials: Reading...
capella_database_credential.new_database_credential: Refreshing state... [id=95591a3b-7031-4257-8d9e-7c4620d14618]
data.capella_database_credentials.existing_credentials: Read complete after 2s

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
          ~ created_at  = "2023-10-04 04:58:00.034423938 +0000 UTC" -> (known after apply)
          ~ created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ modified_at = "2023-10-04 04:58:00.034423938 +0000 UTC" -> (known after apply)
          ~ modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
        id              = "95591a3b-7031-4257-8d9e-7c4620d14618"
        name            = "test_db_user"
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_database_credential = (sensitive value)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_database_credential.new_database_credential: Modifying... [id=95591a3b-7031-4257-8d9e-7c4620d14618]
capella_database_credential.new_database_credential: Modifications complete after 5s [id=95591a3b-7031-4257-8d9e-7c4620d14618]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

database_credentials_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist([
    {
      "access" = tolist([])
      "audit" = {
        "created_at" = "2023-10-04 04:58:00.034423938 +0000 UTC"
        "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
        "modified_at" = "2023-10-04 04:58:00.034423938 +0000 UTC"
        "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
        "version" = 1
      }
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "id" = "95591a3b-7031-4257-8d9e-7c4620d14618"
      "name" = "test_db_user"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_database_credential = <sensitive>
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
data.capella_database_credentials.existing_credentials: Reading...
capella_database_credential.new_database_credential: Refreshing state... [id=95591a3b-7031-4257-8d9e-7c4620d14618]
data.capella_database_credentials.existing_credentials: Read complete after 3s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_database_credential.new_database_credential will be destroyed
  - resource "capella_database_credential" "new_database_credential" {
      - access          = [
          - {
              - privileges = [
                  - "data_reader",
                ] -> null
            },
        ]
      - audit           = {
          - created_at  = "2023-10-04 04:58:00.034423938 +0000 UTC" -> null
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> null
          - modified_at = "2023-10-04 04:58:00.034423938 +0000 UTC" -> null
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> null
          - version     = 1 -> null
        }
      - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6" -> null
      - id              = "95591a3b-7031-4257-8d9e-7c4620d14618" -> null
      - name            = "test_db_user" -> null
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - password        = (sensitive value)
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - database_credentials_list = {
      - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      - data            = [
          - {
              - access          = []
              - audit           = {
                  - created_at  = "2023-10-04 04:58:00.034423938 +0000 UTC"
                  - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
                  - modified_at = "2023-10-04 04:58:00.034423938 +0000 UTC"
                  - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
                  - version     = 1
                }
              - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              - id              = "95591a3b-7031-4257-8d9e-7c4620d14618"
              - name            = "test_db_user"
              - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
            },
        ]
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    } -> null
  - new_database_credential   = (sensitive value)

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_database_credential.new_database_credential: Destroying... [id=95591a3b-7031-4257-8d9e-7c4620d14618]
capella_database_credential.new_database_credential: Destruction complete after 2s

Destroy complete! Resources: 1 destroyed.
```
