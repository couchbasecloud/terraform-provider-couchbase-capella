# Capella App Services Allowed CIDR Example

This example shows how to create and manage allowed CIDR blocks for an App Service in Capella.

This creates a new allowed CIDR block in the selected Capella App Service and lists its currently allowed CIDR blocks. It uses the App Service ID to create and list allowed CIDR blocks.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In the following example, we will:

1. CREATE: Add a new CIDR block to the App Service's list of allowed IP addresses, as stated in the `create_allowedcidr.tf` file.
2. LIST: List existing allowed CIDR blocks in Capella as stated in the `list_allowedcidr.tf` file.
3. IMPORT: Import an allowed CIDR block that exists in Capella but not in the terraform state file.
4. DELETE: Delete the newly created allowed CIDR block from Capella.

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
data.capella_app_services_cidr.existing_allowedcidr: Reading...
data.capella_app_services_cidr.existing_allowedcidr: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_app_services_cidr.new_allowedcidr will be created
  + resource "capella_app_services_cidr" "new_allowedcidr" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + allowedcidr_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + modified_at = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
                  + version     = 1
                }
              + cidr            = "23.121.17.137/32"
              + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
              + comment         = ""
              + expires_at      = null
              + id              = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_allowedcidr   = {
      + audit           = (known after apply)
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + if_match        = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new allowedcidr

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
data.capella_app_services_cidr.existing_allowedcidr: Reading...
data.capella_app_services_cidr.existing_allowedcidr: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_app_services_cidr.new_allowedcidr will be created
  + resource "capella_app_services_cidr" "new_allowedcidr" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + allowedcidr_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + modified_at = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
                  + version     = 1
                }
              + cidr            = "23.121.17.137/32"
              + comment         = ""
              + expires_at      = null
              + id              = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_allowedcidr   = {
      + audit           = (known after apply)
      + cidr            = "10.0.0.0/16"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + comment         = "Allow access from another VPC"
      + expires_at      = "2023-11-14T21:49:58.465Z"
      + id              = (known after apply)
      + if_match        = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_app_services_cidr.new_allowedcidr: Creating...
capella_app_services_cidr.new_allowedcidr: Creation complete after 8s [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

allowedcidr_list = {
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
        "modified_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
        "version" = 1
      }
      "cidr" = "23.121.17.137/32"
      "comment" = ""
      "expires_at" = tostring(null)
      "id" = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_allowedcidr = {
  "audit" = {
    "created_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cidr" = "10.0.0.0/16"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "comment" = "Allow access from another VPC"
  "expires_at" = "2023-11-14T21:49:58.465Z"
  "id" = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2"
  "if_match" = tostring(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

### Note the allowedcidr ID for the new allowedcidr
Command: `terraform output new_allowedcidr`

Sample Output:
```
$ terraform output new_allowedcidr
{
  "audit" = {
    "created_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-04 02:44:46.31810374 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cidr" = "10.0.0.0/16"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "comment" = "Allow access from another VPC"
  "expires_at" = "2023-11-14T21:49:58.465Z"
  "id" = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2"
  "if_match" = tostring(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

In this case, the allowedcidr ID for my new allowedcidr is `854cbdf0-8ae3-4a42-9227-59c52a5ab4f2`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_app_services_cidr.existing_allowedcidr
couchbase-capella_app_services_cidr.new_allowedcidr
```

## IMPORT
### Remove the resource `new_allowedcidr` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_services_cidr.new_allowedcidr`

Sample Output:
```
$ terraform state rm couchbase-capella_app_services_cidr.new_allowedcidr
Removed couchbase-capella_app_services_cidr.new_allowedcidr
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file. The resource will still exist in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_services_cidr.new_allowedcidr id=<allowedcidr_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_app_services_cidr.new_allowedcidr id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
$ terraform import couchbase-capella_app_services_cidr.new_allowedcidr id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000
capella_app_services_cidr.new_allowedcidr: Importing from ID "id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000",,app_service_id=ffffffff-aaaa-1414-eeee-000000000000...
data.capella_app_services_cidr.existing_allowedcidr: Reading...
capella_app_services_cidr.new_allowedcidr: Import prepared!
  Prepared capella_app_services_cidr for import
capella_app_services_cidr.new_allowedcidr: Refreshing state... [id=id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000]
data.capella_app_services_cidr.existing_allowedcidr: Read complete after 1s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

In this section, we pass the IDs as a single comma-separated string.
The first ID in the string is the a CIDR block ID, which identifies the resource that we want to import.
The second ID is the organization ID, which identifies the organization with the project.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the cluster ID, which identifies the cluster associated with the App Service.
The fifth ID is the App Service ID, which identifies the specific App Service we are working with.

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
data.capella_app_services_cidr.existing_allowedcidr: Reading...
capella_app_services_cidr.new_allowedcidr: Refreshing state... [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]
data.capella_app_services_cidr.existing_allowedcidr: Read complete after 1s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to change the allowedcidr configuration settings.

Command: `terraform apply -var 'allowedcidr={durability_level="majority",name="new_terraform_allowedcidr"}'`

Sample Output:
```
$ terraform apply -var comment="updated allowedcidr comment"
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_app_services_cidr.existing_allowedcidr: Reading...
capella_app_services_cidr.new_allowedcidr: Refreshing state... [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]
data.capella_app_services_cidr.existing_allowedcidr: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # capella_app_services_cidr.new_allowedcidr must be replaced
-/+ resource "capella_app_services_cidr" "new_allowedcidr" {
      ~ audit           = {
          ~ created_at  = "2023-10-04 02:44:46.31810374 +0000 UTC" -> (known after apply)
          ~ created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ modified_at = "2023-10-04 02:44:46.31810374 +0000 UTC" -> (known after apply)
          ~ modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD" -> (known after apply)
          ~ version     = 1 -> (known after apply)
        }
      ~ comment         = "Allow access from another VPC" -> "updated allowedcidr comment" # forces replacement
      ~ id              = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2" -> (known after apply)
        # (5 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_allowedcidr = {
      ~ audit           = {
          - created_at  = "2023-10-04 02:44:46.31810374 +0000 UTC"
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - modified_at = "2023-10-04 02:44:46.31810374 +0000 UTC"
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - version     = 1
        } -> (known after apply)
      ~ comment         = "Allow access from another VPC" -> "updated allowedcidr comment"
      ~ id              = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2" -> (known after apply)
        # (6 unchanged elements hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_app_services_cidr.new_allowedcidr: Destroying... [id=854cbdf0-8ae3-4a42-9227-59c52a5ab4f2]
capella_app_services_cidr.new_allowedcidr: Destruction complete after 3s
capella_app_services_cidr.new_allowedcidr: Creating...
capella_app_services_cidr.new_allowedcidr: Creation complete after 2s [id=358387f9-9780-4419-9ca6-b5e8a3b457dc]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

allowedcidr_list = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
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
      "comment" = "Allow access from another VPC"
      "expires_at" = "2023-11-14T21:49:58.465Z"
      "id" = "854cbdf0-8ae3-4a42-9227-59c52a5ab4f2"
    },
    {
      "audit" = {
        "created_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
        "modified_at" = "2023-10-04 02:44:17.216362615 +0000 UTC"
        "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
        "version" = 1
      }
      "cidr" = "23.121.17.137/32"
      "comment" = ""
      "expires_at" = tostring(null)
      "id" = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_allowedcidr = {
  "audit" = {
    "created_at" = "2023-10-04 02:49:31.491935751 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2023-10-04 02:49:31.491935751 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cidr" = "10.0.0.0/16"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "comment" = "updated allowedcidr comment"
  "expires_at" = "2023-11-14T21:49:58.465Z"
  "id" = "358387f9-9780-4419-9ca6-b5e8a3b457dc"
  "if_match" = tostring(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
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
data.capella_app_services_cidr.existing_allowedcidr: Reading...
capella_app_services_cidr.new_allowedcidr: Refreshing state... [id=358387f9-9780-4419-9ca6-b5e8a3b457dc]
data.capella_app_services_cidr.existing_allowedcidr: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_app_services_cidr.new_allowedcidr will be destroyed
  - resource "capella_app_services_cidr" "new_allowedcidr" {
      - audit           = {
          - created_at  = "2023-10-04 02:49:31.491935751 +0000 UTC"
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - modified_at = "2023-10-04 02:49:31.491935751 +0000 UTC"
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - version     = 1
        }
      - cidr            = "10.0.0.0/16"
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - comment         = "updated allowedcidr comment"
      - expires_at      = "2023-11-14T21:49:58.465Z"
      - id              = "358387f9-9780-4419-9ca6-b5e8a3b457dc"
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - allowedcidr_list = {
      - app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
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
              - comment         = "updated allowedcidr comment"
              - expires_at      = "2023-11-14T21:49:58.465Z"
              - id              = "358387f9-9780-4419-9ca6-b5e8a3b457dc"
            },
          - {
              - audit           = {
                  - created_at  = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  - created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
                  - modified_at = "2023-10-04 02:44:17.216362615 +0000 UTC"
                  - modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
                  - version     = 1
                }
              - cidr            = "23.121.17.137/32"
              - comment         = ""
              - expires_at      = null
              - id              = "bbaf68d3-6e8a-433e-aa78-d12a79da4911"
              - if_match        = null
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null
  - new_allowedcidr   = {
      - audit           = {
          - created_at  = "2023-10-04 02:49:31.491935751 +0000 UTC"
          - created_by  = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - modified_at = "2023-10-04 02:49:31.491935751 +0000 UTC"
          - modified_by = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
          - version     = 1
        }
      - cidr            = "10.0.0.0/16"
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - comment         = "updated allowedcidr comment"
      - expires_at      = "2023-11-14T21:49:58.465Z"
      - id              = "358387f9-9780-4419-9ca6-b5e8a3b457dc"
      - if_match        = null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_app_services_cidr.new_allowedcidr: Destroying... [id=358387f9-9780-4419-9ca6-b5e8a3b457dc]
capella_app_services_cidr.new_allowedcidr: Destruction complete after 2s

Destroy complete! Resources: 1 destroyed.
```
