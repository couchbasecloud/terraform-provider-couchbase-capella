# Capella Allowlists Example

This example shows how to create and manage Allowlists in Capella.

This creates a new allowlist in the selected Capella cluster and lists existing allowlists in the cluster. It uses the cluster ID to create and list allowlists.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new allowlist in Capella as stated in the `create_allowlist.tf` file.
2. UPDATE: Update the allowlist configuration using Terraform.
3. LIST: List existing allowlists in Capella as stated in the `list_allowlists.tf` file.
4. IMPORT: Import a allowlist that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created allowlist from Capella.

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
data.capella_allowlists.existing_allowlists: Reading...
data.capella_allowlists.existing_allowlists: Read complete after 0s

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
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + allowlists_list = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + modified_at = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + version     = 1
                }
              + cidr            = "23.121.17.137/32"
              + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              + comment         = ""
              + expires_at      = null
              + id              = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
              + if_match        = null
              + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_allowlist   = {
      + audit           = (known after apply)
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + if_match        = null
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Allowlist

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
data.capella_allowlists.existing_allowlists: Reading...
data.capella_allowlists.existing_allowlists: Read complete after 0s

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
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + allowlists_list = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + modified_at = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + version     = 1
                }
              + cidr            = "23.121.17.137/32"
              + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              + comment         = ""
              + expires_at      = null
              + id              = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
              + if_match        = null
              + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_allowlist   = {
      + audit           = (known after apply)
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + if_match        = null
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_allowlist.new_allowlist: Creating...
capella_allowlist.new_allowlist: Creation complete after 8s [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

allowlists_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "created_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "modified_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "modified_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "version" = 1
      }
      "cidr" = "23.121.17.137/32"
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "comment" = ""
      "expires_at" = tostring(null)
      "id" = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
      "if_match" = tostring(null)
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_allowlist = {
  "audit" = {
    "created_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cidr" = "10.0.0.0/16"
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "comment" = "Allow access from another VPC"
  "expires_at" = "2023-11-14T21:49:58.465Z"
  "id" = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2"
  "if_match" = tostring(null)
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
```

### Note the Allowlist ID for the new Allowlist
Command: `terraform output new_allowlist`

Sample Output:
```
$ terraform output new_allowlist
{
  "audit" = {
    "created_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cidr" = "10.0.0.0/16"
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "comment" = "Allow access from another VPC"
  "expires_at" = "2023-11-14T21:49:58.465Z"
  "id" = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2"
  "if_match" = tostring(null)
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
```

In this case, the allowlist ID for my new allowlist is `854cbdf0-8ae3-4a42-9227-59c52a5ab4f2`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_allowlists.existing_allowlists
couchbase-capella_allowlist.new_allowlist
```

## IMPORT
### Remove the resource `new_allowlist` from the Terraform State file

Command: `terraform state rm couchbase-capella_allowlist.new_allowlist`

Sample Output:
```
$ terraform state rm couchbase-capella_allowlist.new_allowlist
Removed couchbase-capella_allowlist.new_allowlist
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_allowlist.new_allowlist id=<allowlist_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_allowlist.new_allowlist id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d`

Sample Output:
```
$ terraform import couchbase-capella_allowlist.new_allowlist id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d
capella_allowlist.new_allowlist: Importing from ID "id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d"...
data.capella_allowlists.existing_allowlists: Reading...
capella_allowlist.new_allowlist: Import prepared!
  Prepared capella_allowlist for import
capella_allowlist.new_allowlist: Refreshing state... [id=id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d]
data.capella_allowlists.existing_allowlists: Read complete after 1s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the allowlist ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the allowlist belongs.
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_allowlists.existing_allowlists: Reading...
capella_allowlist.new_allowlist: Refreshing state... [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]
data.capella_allowlists.existing_allowlists: Read complete after 1s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to change the allowlist configuration settings.

Command: `terraform apply -var 'allowlist={durability_level="majority",name="new_terraform_allowlist"}'`

Sample Output:
```
$ terraform apply -var comment="updated allowlist comment"
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_allowlists.existing_allowlists: Reading...
capella_allowlist.new_allowlist: Refreshing state... [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]
data.capella_allowlists.existing_allowlists: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # capella_allowlist.new_allowlist must be replaced
-/+ resource "capella_allowlist" "new_allowlist" {
      ~ audit           = {
          ~ created_at  = "2023-10-04 02:44:46.31810374 +0000 UTC" -> (known after apply)
          ~ created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ modified_at = "2023-10-04 02:44:46.31810374 +0000 UTC" -> (known after apply)
          ~ modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
      ~ comment         = "Allow access from another VPC" -> "updated allowlist comment" # forces replacement
      ~ id              = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2" -> (known after apply)
        # (5 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_allowlist = {
      ~ audit           = {
          - created_at  = "2023-10-04 02:44:46.31810374 +0000 UTC"
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - modified_at = "2023-10-04 02:44:46.31810374 +0000 UTC"
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - version     = 1
        } -> (known after apply)
      ~ comment         = "Allow access from another VPC" -> "updated allowlist comment"
      ~ id              = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2" -> (known after apply)
        # (6 unchanged elements hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_allowlist.new_allowlist: Destroying... [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]
capella_allowlist.new_allowlist: Destruction complete after 3s
capella_allowlist.new_allowlist: Creating...
capella_allowlist.new_allowlist: Creation complete after 2s [id=358387f9-9780-4419-9ca6-b5e8a3b457dc]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

allowlists_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
        "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
        "modified_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
        "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
        "version" = 1
      }
      "cidr" = "10.0.0.0/16"
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "comment" = "Allow access from another VPC"
      "expires_at" = "2023-11-14T21:49:58.465Z"
      "id" = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2"
      "if_match" = tostring(null)
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    },
    {
      "audit" = {
        "created_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "created_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "modified_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "modified_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "version" = 1
      }
      "cidr" = "23.121.17.137/32"
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "comment" = ""
      "expires_at" = tostring(null)
      "id" = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
      "if_match" = tostring(null)
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_allowlist = {
  "audit" = {
    "created_at" = "2023-10-04 02:49:31.491935751 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-04 02:49:31.491935751 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cidr" = "10.0.0.0/16"
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "comment" = "updated allowlist comment"
  "expires_at" = "2023-11-14T21:49:58.465Z"
  "id" = "358387f9-9780-4419-9ca6-b5e8a3b457dc"
  "if_match" = tostring(null)
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_allowlists.existing_allowlists: Reading...
capella_allowlist.new_allowlist: Refreshing state... [id=358387f9-9780-4419-9ca6-b5e8a3b457dc]
data.capella_allowlists.existing_allowlists: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_allowlist.new_allowlist will be destroyed
  - resource "capella_allowlist" "new_allowlist" {
      - audit           = {
          - created_at  = "2023-10-04 02:49:31.491935751 +0000 UTC" -> null
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> null
          - modified_at = "2023-10-04 02:49:31.491935751 +0000 UTC" -> null
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> null
          - version     = 1 -> null
        }
      - cidr            = "10.0.0.0/16" -> null
      - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6" -> null
      - comment         = "updated allowlist comment" -> null
      - expires_at      = "2023-11-14T21:49:58.465Z" -> null
      - id              = "358387f9-9780-4419-9ca6-b5e8a3b457dc" -> null
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - allowlists_list = {
      - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      - data            = [
          - {
              - audit           = {
                  - created_at  = "2023-10-04 02:49:31.491935751 +0000 UTC"
                  - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
                  - modified_at = "2023-10-04 02:49:31.491935751 +0000 UTC"
                  - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
                  - version     = 1
                }
              - cidr            = "10.0.0.0/16"
              - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              - comment         = "updated allowlist comment"
              - expires_at      = "2023-11-14T21:49:58.465Z"
              - id              = "358387f9-9780-4419-9ca6-b5e8a3b457dc"
              - if_match        = null
              - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
            },
          - {
              - audit           = {
                  - created_at  = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  - created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - modified_at = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  - modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - version     = 1
                }
              - cidr            = "23.121.17.137/32"
              - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              - comment         = ""
              - expires_at      = null
              - id              = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
              - if_match        = null
              - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
            },
        ]
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    } -> null
  - new_allowlist   = {
      - audit           = {
          - created_at  = "2023-10-04 02:49:31.491935751 +0000 UTC"
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - modified_at = "2023-10-04 02:49:31.491935751 +0000 UTC"
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - version     = 1
        }
      - cidr            = "10.0.0.0/16"
      - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      - comment         = "updated allowlist comment"
      - expires_at      = "2023-11-14T21:49:58.465Z"
      - id              = "358387f9-9780-4419-9ca6-b5e8a3b457dc"
      - if_match        = null
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_allowlist.new_allowlist: Destroying... [id=358387f9-9780-4419-9ca6-b5e8a3b457dc]
capella_allowlist.new_allowlist: Destruction complete after 2s

Destroy complete! Resources: 1 destroyed.
```
