# Capella Database Role Example

This example shows how to create and manage database user roles in Capella.

Database roles define reusable sets of privileges scoped to specific buckets, scopes, and collections within a cluster. They can be assigned to database credentials.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

## Prerequisites

- A Capella organization, project, and cluster
- A valid V4 API key with `organizationOwner` or `projectOwner` permissions
- At least one bucket in the cluster (the example uses `travel-sample`)

## Example Walkthrough

1. CREATE: Create a new database role as defined in `create_database_role.tf`.
2. UPDATE: Update the role's access privileges or description.
3. IMPORT: Import a database role that exists in Capella but not in the Terraform state file.
4. DELETE: Delete the database role from Capella.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the placeholder values.

## CREATE

### View the plan

Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform used the selected providers to generate the following execution plan.
Resource actions are indicated with the following symbols:
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
                                  + collections = [
                                      + "sales",
                                    ]
                                  + name        = "inventory"
                                },
                            ]
                        },
                    ]
                }
            },
        ]
      + audit           = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + description     = "A test role with read and write access"
      + id              = (known after apply)
      + name            = "test-role-001"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the plan

Command: `terraform apply`

Sample Output:
```
$ terraform apply

couchbase-capella_database_role.new_database_role: Creating...
couchbase-capella_database_role.new_database_role: Creation complete after 1s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_database_role = {
  "access" = toset([
    {
      "privileges" = toset([
        "dataRead",
      ])
      "resources" = {
        "buckets" = toset([
          {
            "name" = "travel-sample"
            "scopes" = toset([
              {
                "collections" = toset([
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
        "buckets" = toset([
          {
            "name" = "travel-sample"
            "scopes" = toset([
              {
                "collections" = toset([
                  "sales",
                ])
                "name" = "inventory"
              },
            ])
          },
        ])
      }
    },
  ])
  "audit" = {
    "created_at" = "2024-01-15 10:30:00.000000000 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2024-01-15 10:30:00.000000000 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 1
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "description" = "A test role with read and write access"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "test-role-001"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
database_role_id = "ffffffff-aaaa-1414-eeee-000000000000"
```

### Note the Database Role ID

Command: `terraform output database_role_id`

Sample Output:
```
$ terraform output database_role_id
"ffffffff-aaaa-1414-eeee-000000000000"
```

## IMPORT

### Remove the resource from the Terraform state file

Command: `terraform state rm couchbase-capella_database_role.new_database_role`

Sample Output:
```
$ terraform state rm couchbase-capella_database_role.new_database_role
Removed couchbase-capella_database_role.new_database_role
Successfully removed 1 resource instance(s).
```

This only removes the resource from state. The role still exists in Capella.

### Import the resource back into Terraform

Command: `terraform import couchbase-capella_database_role.new_database_role id=<role_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
$ terraform import couchbase-capella_database_role.new_database_role id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_database_role.new_database_role: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_database_role.new_database_role: Import prepared!
couchbase-capella_database_role.new_database_role: Refreshing state...

Import successful!
```

The import ID is a comma-separated string of `id=<role_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`.

## UPDATE

Edit `terraform.tfvars` to change the role's access privileges or description, then apply.

Command: `terraform apply`

Sample Output:
```
$ terraform apply

couchbase-capella_database_role.new_database_role: Modifying...
couchbase-capella_database_role.new_database_role: Modifications complete after 1s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## DELETE

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy

couchbase-capella_database_role.new_database_role: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_database_role.new_database_role: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```

