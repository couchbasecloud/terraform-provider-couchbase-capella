# Capella Collections Example

This example shows how to create and manage Collections in Capella.

This creates a new collection in the selected Capella cluster and lists existing collections in a scope of the bucket. It uses the scope name to create and list collections.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new collection in Capella as stated in the `create_collection.tf` file.
2. UPDATE: Update the collection configuration using Terraform.
3. LIST: List existing collections in Capella as stated in the `list_collections.tf` file.
4. IMPORT: Import a collection that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created collection from Capella.

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
data.couchbase-capella_collections.existing_collections: Reading...
data.couchbase-capella_collections.existing_collections: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection will be created
  + resource "couchbase-capella_collection" "new_collection" {
      + bucket_id       = "YjE="
      + cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + scope_name      = "s1"
      + uid             = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_collection   = {
      + bucket_id       = "YjE="
      + cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + scope_name      = "s1"
      + uid             = (known after apply)
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new collection

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
data.couchbase-capella_collections.existing_collections: Reading...
data.couchbase-capella_collections.existing_collections: Read complete after 7s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection will be created
  + resource "couchbase-capella_collection" "new_collection" {
      + bucket_id       = "YjE="
      + cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + scope_name      = "s1"
      + uid             = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_collection   = {
      + bucket_id       = "YjE="
      + cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + scope_name      = "s1"
      + uid             = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_collection.new_collection: Creating...
couchbase-capella_collection.new_collection: Creation complete after 1s [name=new_terraform_collection]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

collections_list = {
  "bucket_id" = "YjE="
  "cluster_id" = "0d9a6dd5-4d55-49be-8137-896f21425beb"
  "data" = tolist([
    {
      "max_ttl" = 100
      "name" = "c1"
      "uid" = "8"
    },
  ])
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "scope_name" = "s1"
}
new_collection = {
  "bucket_id" = "YjE="
  "cluster_id" = "0d9a6dd5-4d55-49be-8137-896f21425beb"
  "max_ttl" = 200
  "name" = "new_terraform_collection"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "scope_name" = "s1"
  "uid" = "a"
}

```

### Note the Collection name for the new Collection
Command: `terraform output new_collection`

Sample Output:
```
$ terraform output new_collection
{
  "bucket_id" = "YjE="
  "cluster_id" = "0d9a6dd5-4d55-49be-8137-896f21425beb"
  "max_ttl" = 200
  "name" = "new_terraform_collection"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "scope_name" = "s1"
  "uid" = "a"
}

```

In this case, the collection name for my new collection is `new_terraform_collection`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$  terraform state list
data.couchbase-capella_collections.existing_collections
couchbase-capella_collection.new_collection
```

## IMPORT
### Remove the resource `new_collection` from the Terraform State file

Command: `terraform state rm couchbase-capella_collection.new_collection`

Sample Output:
```
$ terraform state rm couchbase-capella_collection.new_collection
Removed couchbase-capella_collection.new_collection
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_collection.new_collection collection_name=<collection_name>,scope_name=<scope_name>,bucket_id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_collection.new_collection collection_name=new_terraform_collection,scope_name=s1,bucket_id=YjE=,cluster_id=0d9a6dd5-4d55-49be-8137-896f21425beb,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0`

Sample Output:
```
$ terraform import couchbase-capella_collection.new_collection collection_name=new_terraform_collection,scope_name=s1,bucket_id=YjE=,cluster_id=0d9a6dd5-4d55-49be-8137-896f21425beb,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0
couchbase-capella_collection.new_collection: Importing from ID "collection_name=new_terraform_collection,scope_name=s1,bucket_id=YjE=,cluster_id=0d9a6dd5-4d55-49be-8137-896f21425beb,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0"...
data.couchbase-capella_collections.existing_collections: Reading...
couchbase-capella_collection.new_collection: Import prepared!
  Prepared couchbase-capella_collection for import
couchbase-capella_collection.new_collection: Refreshing state...
data.couchbase-capella_collections.existing_collections: Read complete after 0s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the collection Name i.e. the Name of the resource that we want to import.
The second ID in the string is the scope Name i.e. the ID of the scope to which the collection belongs.
The third ID is the bucket ID i.e. the ID of the bucket to which the scope belongs.
The fourth ID is the cluster ID i.e. the ID of the cluster to which the bucket belongs.
The fifth ID is the project ID i.e. the ID of the project to which the cluster belongs.
The sixth ID is the organization ID i.e. the ID of the organization to which the project belongs.

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
data.couchbase-capella_collections.existing_collections: Reading...
couchbase-capella_collection.new_collection: Refreshing state...
data.couchbase-capella_collections.existing_collections: Read complete after 0s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed
```

## UPDATE
### Let us edit the terraform.tfvars file to change the scope configuration settings.

Command: `terraform apply -var 'collection={collection_name="new_terraform_collection2", max_ttl=300}'`

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
data.couchbase-capella_collections.existing_collections: Reading...
couchbase-capella_collection.new_collection: Refreshing state...
data.couchbase-capella_collections.existing_collections: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection must be replaced
-/+ resource "couchbase-capella_collection" "new_collection" {
      ~ collection_name = "new_terraform_collection" -> "new_terraform_collection2" # forces replacement
      ~ max_ttl         = 200 -> 300 # forces replacement
      ~ uid             = "b" -> (known after apply)
        # (5 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_collection   = {
      ~ collection_name = "new_terraform_collection" -> "new_terraform_collection2"
      ~ max_ttl         = 200 -> 300
      ~ uid             = "b" -> (known after apply)
        # (5 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_collection.new_collection: Destroying...
couchbase-capella_collection.new_collection: Destruction complete after 1s
couchbase-capella_collection.new_collection: Creating...
couchbase-capella_collection.new_collection: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

collections_list = {
  "bucket_id" = "YjE="
  "cluster_id" = "0d9a6dd5-4d55-49be-8137-896f21425beb"
  "data" = tolist([
    {
      "collection_name" = "new_terraform_collection"
      "max_ttl" = 200
      "uid" = "b"
    },
    {
      "collection_name" = "c1"
      "max_ttl" = 100
      "uid" = "8"
    },
  ])
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "scope_name" = "s1"
}
new_collection = {
  "bucket_id" = "YjE="
  "cluster_id" = "0d9a6dd5-4d55-49be-8137-896f21425beb"
  "collection_name" = "new_terraform_collection2"
  "max_ttl" = 300
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "scope_name" = "s1"
  "uid" = "c"
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
data.couchbase-capella_collections.existing_collections: Reading...
couchbase-capella_collection.new_collection: Refreshing state...
data.couchbase-capella_collections.existing_collections: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection will be destroyed
  - resource "couchbase-capella_collection" "new_collection" {
      - bucket_id       = "YjE=" -> null
      - cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb" -> null
      - collection_name = "new_terraform_collection2" -> null
      - max_ttl         = 300 -> null
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0" -> null
      - project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3" -> null
      - scope_name      = "s1" -> null
      - uid             = "c" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - collections_list = {
      - bucket_id       = "YjE="
      - cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      - data            = [
          - {
              - collection_name = "new_terraform_collection2"
              - max_ttl         = 300
              - uid             = "c"
            },
          - {
              - collection_name = "c1"
              - max_ttl         = 100
              - uid             = "8"
            },
        ]
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      - project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      - scope_name      = "s1"
    } -> null
  - new_collection   = {
      - bucket_id       = "YjE="
      - cluster_id      = "0d9a6dd5-4d55-49be-8137-896f21425beb"
      - collection_name = "new_terraform_collection2"
      - max_ttl         = 300
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      - project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      - scope_name      = "s1"
      - uid             = "c"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_collection.new_collection: Destroying...
couchbase-capella_collection.new_collection: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```