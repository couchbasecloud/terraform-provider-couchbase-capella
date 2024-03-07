# Capella Scopes Example

This example shows how to create and manage Scopes in Capella.

This creates a new scope in the selected Capella cluster and lists existing scopes in the bucket. It uses the bucket id to create and list scopes.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new scope in Capella as stated in the `create_scope.tf` file.
2. UPDATE: Update the scope configuration using Terraform.
3. LIST: List existing scope in Capella as stated in the `list_scopes.tf` file.
4. IMPORT: Import a scope that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created scope from Capella.

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
data.couchbase-capella_scopes.existing_scopes: Reading...
data.couchbase-capella_scopes.existing_scopes: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_scope.new_scope will be created
  + resource "couchbase-capella_scope" "new_scope" {
      + bucket_id       = "YjE="
      + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      + collections     = (known after apply)
      + name            = "new_terraform_scope"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + uid             = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_scope   = {
      + bucket_id       = "YjE="
      + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      + collections     = (known after apply)
      + name            = "new_terraform_scope"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + uid             = (known after apply)
    }
  + scopes_list = {
      + bucket_id       = "YjE="
      + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + scopes          = [
          + {
              + bucket_id       = "YjE="
              + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
              + collections     = [
                  + {
                      + max_ttl = 0
                      + name    = "c21"
                      + uid     = "c"
                    },
                  + {
                      + max_ttl = 50
                      + name    = "c23"
                      + uid     = "b"
                    },
                  + {
                      + max_ttl = 100
                      + name    = "c22"
                      + uid     = "a"
                    },
                ]
              + name            = "s2"
              + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
              + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
              + uid             = "9"
            },
          + {
              + bucket_id       = "YjE="
              + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
              + collections     = [
                  + {
                      + max_ttl = 0
                      + name    = "c12"
                      + uid     = "9"
                    },
                  + {
                      + max_ttl = 0
                      + name    = "c11"
                      + uid     = "8"
                    },
                ]
              + name            = "s1"
              + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
              + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
              + uid             = "8"
            },
          + {
              + bucket_id       = "YjE="
              + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
              + collections     = [
                  + {
                      + max_ttl = 0
                      + name    = "_default"
                      + uid     = "0"
                    },
                ]
              + name            = "_default"
              + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
              + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
              + uid             = "0"
            },
        ]
      + uid             = "8"
    }


─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new scope

Command: `terraform apply`

Sample Output:
```
$terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_scopes.existing_scopes: Reading...
data.couchbase-capella_scopes.existing_scopes: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_scope.new_scope will be created
  + resource "couchbase-capella_scope" "new_scope" {
      + bucket_id       = "YjE="
      + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      + collections     = (known after apply)
      + name            = "new_terraform_scope"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + uid             = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_scope   = {
      + bucket_id       = "YjE="
      + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      + collections     = (known after apply)
      + name            = "new_terraform_scope"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + uid             = (known after apply)
    }
  + scopes_list = {
      + bucket_id       = "YjE="
      + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + scopes          = [
          + {
              + bucket_id       = "YjE="
              + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
              + collections     = [
                  + {
                      + max_ttl = 0
                      + name    = "c21"
                      + uid     = "c"
                    },
                  + {
                      + max_ttl = 50
                      + name    = "c23"
                      + uid     = "b"
                    },
                  + {
                      + max_ttl = 100
                      + name    = "c22"
                      + uid     = "a"
                    },
                ]
              + name            = "s2"
              + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
              + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
              + uid             = "9"
            },
          + {
              + bucket_id       = "YjE="
              + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
              + collections     = [
                  + {
                      + max_ttl = 0
                      + name    = "c12"
                      + uid     = "9"
                    },
                  + {
                      + max_ttl = 0
                      + name    = "c11"
                      + uid     = "8"
                    },
                ]
              + name            = "s1"
              + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
              + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
              + uid             = "8"
            },
          + {
              + bucket_id       = "YjE="
              + cluster_id      = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
              + collections     = [
                  + {
                      + max_ttl = 0
                      + name    = "_default"
                      + uid     = "0"
                    },
                ]
              + name            = "_default"
              + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
              + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
              + uid             = "0"
            },
        ]
      + uid             = "8"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_scope.new_scope: Creating...
couchbase-capella_scope.new_scope: Creation complete after 1s [name=new_terraform_scope]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_scope = {
  "bucket_id" = "YjE="
  "cluster_id" = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
  "collections" = toset([])
  "name" = "new_terraform_scope"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "uid" = "a"
}
scopes_list = {
  "bucket_id" = "YjE="
  "cluster_id" = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "scopes" = tolist([
    {
      "bucket_id" = "YjE="
      "cluster_id" = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      "collections" = tolist([
        {
          "max_ttl" = 0
          "name" = "c21"
          "uid" = "c"
        },
        {
          "max_ttl" = 50
          "name" = "c23"
          "uid" = "b"
        },
        {
          "max_ttl" = 100
          "name" = "c22"
          "uid" = "a"
        },
      ])
      "name" = "s2"
      "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      "uid" = "9"
    },
    {
      "bucket_id" = "YjE="
      "cluster_id" = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      "collections" = tolist([
        {
          "max_ttl" = 0
          "name" = "c12"
          "uid" = "9"
        },
        {
          "max_ttl" = 0
          "name" = "c11"
          "uid" = "8"
        },
      ])
      "name" = "s1"
      "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      "uid" = "8"
    },
    {
      "bucket_id" = "YjE="
      "cluster_id" = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
      "collections" = tolist([
        {
          "max_ttl" = 0
          "name" = "_default"
          "uid" = "0"
        },
      ])
      "name" = "_default"
      "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      "uid" = "0"
    },
  ])
  "uid" = "8"
}

```

### Note the Scope name for the new Scope
Command: `terraform output new_scope`

Sample Output:
```
$ terraform output new_scope
{
  "bucket_id" = "YjE="
  "cluster_id" = "3733066f-ef48-42a4-95a5-2f7e44c09d72"
  "collections" = toset([])
  "name" = "new_terraform_scope"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "uid" = "a"
}
```

In this case, the scope name for my new scope is `new_terraform_scope`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$  terraform state list
data.couchbase-capella_scopes.existing_scopes
couchbase-capella_scope.new_scope
```

## IMPORT
### Remove the resource `new_scope` from the Terraform State file

Command: `terraform state rm couchbase-capella_scope.new_scope`

Sample Output:
```
$ terraform state rm couchbase-capella_scope.new_scope
Removed couchbase-capella_scope.new_scope
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_scope.new_scope scope_name=<scope_name>,bucket_id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_scope.new_scope scope_name=new_terraform_scope,bucket_id=YjE=,cluster_id=0d9a6dd5-4d55-49be-8137-896f21425beb,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0`

Sample Output:
```
$ terraform import couchbase-capella_scope.new_scope scope_name=new_terraform_scope,bucket_id=YjE=,cluster_id=0d9a6dd5-4d55-49be-8137-896f21425beb,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0
couchbase-capella_scope.new_scope: Importing from ID "scope_name=new_terraform_scope,bucket_id=YjE=,cluster_id=0d9a6dd5-4d55-49be-8137-896f21425beb,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0"...
data.couchbase-capella_scopes.existing_scopes: Reading...
couchbase-capella_scope.new_scope: Import prepared!
  Prepared couchbase-capella_scope for import
couchbase-capella_scope.new_scope: Refreshing state...
data.couchbase-capella_scopes.existing_scopes: Read complete after 1s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the scope Name i.e. the Name of the resource that we want to import.
The second ID is the bucket ID i.e. the ID of the bucket to which the scope belongs.
The third ID is the cluster ID i.e. the ID of the cluster to which the bucket belongs.
The fourth ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fifth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/paulomee.de/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_scopes.existing_scopes: Reading...
couchbase-capella_scope.new_scope: Refreshing state... [name=new_terraform_scope]
data.couchbase-capella_scopes.existing_scopes: Read complete after 1s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to change the scope configuration settings.

Command: `terraform apply -var 'scope={name="terraform_scope_updated"}'`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/paulomee.de/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_scopes.existing_scopes: Reading...
couchbase-capella_scope.new_scope: Refreshing state... [name=new_terraform_scope]
data.couchbase-capella_scopes.existing_scopes: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_scope.new_scope must be replaced
-/+ resource "couchbase-capella_scope" "new_scope" {
      ~ collections     = [] -> (known after apply)
      ~ name            = "new_terraform_scope" -> "terraform_scope_updated" # forces replacement
      ~ uid             = "f" -> (known after apply)
        # (4 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_scope   = {
      ~ collections     = [] -> (known after apply)
      ~ name            = "new_terraform_scope" -> "terraform_scope_updated"
      ~ uid             = "f" -> (known after apply)
        # (4 unchanged attributes hidden)
    }
  ~ scopes_list = {
      ~ scopes          = [
          + {
              + collections = []
              + name        = "new_terraform_scope_example1"
              + uid         = "10"
            },
            {
                collections = []
                name        = "new_terraform_scope"
                uid         = "f"
            },
            # (3 unchanged elements hidden)
        ]
      ~ uid             = "f" -> "10"
        # (4 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_scope.new_scope: Destroying... [name=new_terraform_scope]
couchbase-capella_scope.new_scope: Destruction complete after 1s
couchbase-capella_scope.new_scope: Creating...
couchbase-capella_scope.new_scope: Creation complete after 1s [name=terraform_scope_updated]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

new_scope = {
  "bucket_id" = "YjE="
  "cluster_id" = "0d9a6dd5-4d55-49be-8137-896f21425beb"
  "collections" = toset([])
  "name" = "terraform_scope_updated"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "uid" = "11"
}
scopes_list = {
  "bucket_id" = "YjE="
  "cluster_id" = "0d9a6dd5-4d55-49be-8137-896f21425beb"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "scopes" = tolist([
    {
      "collections" = toset([])
      "name" = "new_terraform_scope_example1"
      "uid" = "10"
    },
    {
      "collections" = toset([])
      "name" = "new_terraform_scope"
      "uid" = "f"
    },
    {
      "collections" = toset([])
      "name" = "s3"
      "uid" = "b"
    },
    {
      "collections" = toset([
        {
          "max_ttl" = 100
          "name" = "c1"
          "uid" = "8"
        },
      ])
      "name" = "s1"
      "uid" = "9"
    },
    {
      "collections" = toset([
        {
          "max_ttl" = 0
          "name" = "_default"
          "uid" = "0"
        },
      ])
      "name" = "_default"
      "uid" = "0"
    },
  ])
  "uid" = "10"
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
│  - couchbasecloud/couchbase-capella in /Users/paulomee.de/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_scopes.existing_scopes: Reading...
couchbase-capella_scope.new_scope: Refreshing state... [name=new_terraform_scope]
data.couchbase-capella_scopes.existing_scopes: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_scope.new_scope will be destroyed
  - resource "couchbase-capella_scope" "new_scope" {
      - bucket_id       = "YjE=" -> null
      - cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb" -> null
      - collections     = [
          - {
              - max_ttl = 214748 -> null
              - name    = "abc" -> null
              - uid     = "9" -> null
            },
        ] -> null
      - name            = "new_terraform_scope" -> null
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0" -> null
      - project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3" -> null
      - uid             = "a" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_scope   = {
      - bucket_id       = "YjE="
      - cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      - collections     = [
          - {
              - max_ttl = 214748
              - name    = "abc"
              - uid     = "9"
            },
        ]
      - name            = "new_terraform_scope"
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      - project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      - uid             = "a"
    } -> null
  - scopes_list = {
      - bucket_id       = "YjE="
      - cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      - project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      - scopes          = [
          - {
              - collections = []
              - name        = "s3"
              - uid         = "b"
            },
          - {
              - collections = [
                  - {
                      - max_ttl = 214748
                      - name    = "abc"
                      - uid     = "9"
                    },
                ]
              - name        = "new_terraform_scope"
              - uid         = "a"
            },
          - {
              - collections = [
                  - {
                      - max_ttl = 100
                      - name    = "c1"
                      - uid     = "8"
                    },
                ]
              - name        = "s1"
              - uid         = "9"
            },
          - {
              - collections = [
                  - {
                      - max_ttl = 0
                      - name    = "_default"
                      - uid     = "0"
                    },
                ]
              - name        = "_default"
              - uid         = "0"
            },
        ]
      - uid             = "d"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_scope.new_scope: Destroying... [name=new_terraform_scope]
couchbase-capella_scope.new_scope: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.

```
