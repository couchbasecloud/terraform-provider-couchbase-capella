# Capella Database Role Example

This example shows how to create and manage Database Roles in Capella.

This creates a new database role in the selected Capella cluster and lists existing database roles in the cluster. It uses the cluster ID to create and list database roles.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new database role in Capella as stated in the `create_database_role.tf` file.
2. UPDATE: Update the database role configuration using Terraform.
3. LIST: List existing database roles in Capella as stated in the `list_database_roles.tf` file.
4. IMPORT: Import a database role that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created database role from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

Each privilege can only be granted down to a certain level. For example, `dataRead` can be narrowed to collections, but `queryManage` is a scope-level privilege and is rejected if `collections` are set. Use the `couchbase-capella_database_privileges` data source (see `examples/database_privileges`) to check the levels each privilege accepts.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
data.couchbase-capella_database_roles.existing_roles: Reading...
data.couchbase-capella_database_roles.existing_roles: Read complete after 0s

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_database_role.new_database_role will be created
  + resource "couchbase-capella_database_role" "new_database_role" {
      + access          = [
          + {
              + privileges = [
                  + "dataRead",
                ]
              + resources  = {
                  + buckets = [
                      + {
                          + name   = "travel-sample"
                          + scopes = [
                              + {
                                  + collections = [
                                      + "airline",
                                      + "airport",
                                    ]
                                  + name        = "inventory"
                                },
                            ]
                        },
                    ]
                }
            },
          + {
              + privileges = [
                  + "queryManage",
                ]
              + resources  = {
                  + buckets = [
                      + {
                          + name   = "travel-sample"
                          + scopes = [
                              + {
                                  + name = "inventory"
                                },
                            ]
                        },
                    ]
                }
            },
        ]
      + audit           = (known after apply)
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + description     = "A test role with read and write access"
      + id              = (known after apply)
      + name            = "test-role-001"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + database_role_id    = (known after apply)
  + database_roles_list = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = null
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_database_role   = {
      + access          = [
          + {
              + privileges = [
                  + "dataRead",
                ]
              + resources  = {
                  + buckets = [
                      + {
                          + name   = "travel-sample"
                          + scopes = [
                              + {
                                  + collections = [
                                      + "airline",
                                      + "airport",
                                    ]
                                  + name        = "inventory"
                                },
                            ]
                        },
                    ]
                }
            },
          + {
              + privileges = [
                  + "queryManage",
                ]
              + resources  = {
                  + buckets = [
                      + {
                          + name   = "travel-sample"
                          + scopes = [
                              + {
                                  + collections = null
                                  + name        = "inventory"
                                },
                            ]
                        },
                    ]
                }
            },
        ]
      + audit           = (known after apply)
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + description     = "A test role with read and write access"
      + id              = (known after apply)
      + name            = "test-role-001"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Database Role

Command: `terraform apply`

Sample Output:
```
$ terraform apply
data.couchbase-capella_database_roles.existing_roles: Reading...
data.couchbase-capella_database_roles.existing_roles: Read complete after 0s

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_database_role.new_database_role: Creating...
couchbase-capella_database_role.new_database_role: Creation complete after 3s [id=95591a3b-7031-4257-8d9e-7c4620d14618]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

database_role_id = "95591a3b-7031-4257-8d9e-7c4620d14618"
database_roles_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist(null) /* of object */
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_database_role = {
  "access" = tolist([
    {
      "privileges" = toset([
        "dataRead",
      ])
      "resources" = {
        "buckets" = tolist([
          {
            "name" = "travel-sample"
            "scopes" = tolist([
              {
                "collections" = tolist([
                  "airline",
                  "airport",
                ])
                "name" = "inventory"
              },
            ])
          },
        ])
      }
    },
    {
      "privileges" = toset([
        "queryManage",
      ])
      "resources" = {
        "buckets" = tolist([
          {
            "name" = "travel-sample"
            "scopes" = tolist([
              {
                "collections" = tolist(null) /* of string */
                "name" = "inventory"
              },
            ])
          },
        ])
      }
    },
  ])
  "audit" = {
    "created_at" = "2026-10-14 05:04:35.472885632 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2026-10-14 05:04:35.472885632 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "description" = "A test role with read and write access"
  "id" = "95591a3b-7031-4257-8d9e-7c4620d14618"
  "name" = "test-role-001"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
```

### Note the Database Role ID for the new Database Role
Command: `terraform output new_database_role`

In this case, the database role ID for my new database role is `95591a3b-7031-4257-8d9e-7c4620d14618`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_database_roles.existing_roles
couchbase-capella_database_role.new_database_role
```

## IMPORT
### Remove the resource `new_database_role` from the Terraform State file

Command: `terraform state rm couchbase-capella_database_role.new_database_role`

Sample Output:
```
$ terraform state rm couchbase-capella_database_role.new_database_role
Removed couchbase-capella_database_role.new_database_role
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_database_role.new_database_role id=<database_role_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_database_role.new_database_role id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d`

Sample Output:
```
$ terraform import couchbase-capella_database_role.new_database_role id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d
couchbase-capella_database_role.new_database_role: Importing from ID "id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d"...
couchbase-capella_database_role.new_database_role: Import prepared!
  Prepared couchbase-capella_database_role for import
data.couchbase-capella_database_roles.existing_roles: Reading...
couchbase-capella_database_role.new_database_role: Refreshing state... [id=id=95591a3b-7031-4257-8d9e-7c4620d14618,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d]
data.couchbase-capella_database_roles.existing_roles: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the database role ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the role belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform apply to confirm that the import was successful

Command: `terraform apply`

Sample Output:
```
$ terraform apply
data.couchbase-capella_database_roles.existing_roles: Reading...
data.couchbase-capella_database_roles.existing_roles: Read complete after 0s

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_database_role.new_database_role: Creating...
couchbase-capella_database_role.new_database_role: Creation complete after 1s [id=d431493a-aed0-4dcd-b9bd-67d9da667573]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

database_role_id = "d431493a-aed0-4dcd-b9bd-67d9da667573"
database_roles_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist(null) /* of object */
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_database_role = {
  "access" = tolist([
    {
      "privileges" = toset([
        "dataRead",
      ])
      "resources" = {
        "buckets" = tolist([
          {
            "name" = "travel-sample"
            "scopes" = tolist([
              {
                "collections" = tolist([
                  "airline",
                  "airport",
                ])
                "name" = "inventory"
              },
            ])
          },
        ])
      }
    },
    {
      "privileges" = toset([
        "queryManage",
      ])
      "resources" = {
        "buckets" = tolist([
          {
            "name" = "travel-sample"
            "scopes" = tolist([
              {
                "collections" = tolist(null) /* of string */
                "name" = "inventory"
              },
            ])
          },
        ])
      }
    },
  ])
  "audit" = {
    "created_at" = "2026-10-14 05:04:39.262285175 +0000 UTC"
    "created_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "modified_at" = "2026-10-14 05:04:39.262285175 +0000 UTC"
    "modified_by" = "osxKeibDiShFFyyqAVNvqWRaWryXBxBD"
    "version" = 1
  }
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "description" = "A test role with read and write access"
  "id" = "d431493a-aed0-4dcd-b9bd-67d9da667573"
  "name" = "test-role-001"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
```

## UPDATE
### Update the Database Role configuration by overriding the `access` variable.

Command: `terraform apply -var 'access=[{privileges=["dataRead"]}]'`

Sample Output:
```
$ terraform apply -var 'access=[{privileges=["dataRead"]}]'
data.couchbase-capella_database_roles.existing_roles: Reading...
couchbase-capella_database_role.new_database_role: Refreshing state... [id=d431493a-aed0-4dcd-b9bd-67d9da667573]
data.couchbase-capella_database_roles.existing_roles: Read complete after 0s

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_database_role.new_database_role: Modifying... [id=d431493a-aed0-4dcd-b9bd-67d9da667573]
couchbase-capella_database_role.new_database_role: Modifications complete after 1s [id=d431493a-aed0-4dcd-b9bd-67d9da667573]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
data.couchbase-capella_database_roles.existing_roles: Reading...
couchbase-capella_database_role.new_database_role: Refreshing state... [id=d431493a-aed0-4dcd-b9bd-67d9da667573]
data.couchbase-capella_database_roles.existing_roles: Read complete after 0s

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_database_role.new_database_role: Destroying... [id=d431493a-aed0-4dcd-b9bd-67d9da667573]
couchbase-capella_database_role.new_database_role: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
