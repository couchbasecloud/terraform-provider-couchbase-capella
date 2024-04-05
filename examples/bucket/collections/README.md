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
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "s1"
      + uid             = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_collection   = {
      + bucket_id       = "YjE="
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
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
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "s1"
      + uid             = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_collection   = {
      + bucket_id       = "YjE="
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + max_ttl         = 200
      + name            = "new_terraform_collection"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "s1"
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
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "max_ttl" = 100
      "name" = "c1"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
}
new_collection = {
  "bucket_id" = "YjE="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "max_ttl" = 200
  "name" = "new_terraform_collection"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
}

```

### Note the Collection name for the new Collection
Command: `terraform output new_collection`

Sample Output:
```
$ terraform output new_collection
{
  "bucket_id" = "YjE="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "max_ttl" = 200
  "name" = "new_terraform_collection"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
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
`terraform import couchbase-capella_collection.new_collection collection_name=c2,scope_name=s1,bucket_id=YjE=,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
$ terraform import couchbase-capella_collection.new_collection collection_name=new_terraform_collection,scope_name=s1,bucket_id=YjE=,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_collection.new_collection: Importing from ID "collection_name=new_terraform_collection,scope_name=s1,bucket_id=YjE=,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
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
### Let us edit the terraform.tfvars file to change the collection configuration settings.

1. Update name for server versions < 7.6.0 (destroy and replace)

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
couchbase-capella_collection.new_collection: Refreshing state...
data.couchbase-capella_collections.existing_collections: Reading...
data.couchbase-capella_collections.existing_collections: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection must be replaced
-/+ resource "couchbase-capella_collection" "new_collection" {
      ~ collection_name = "new_terraform_collection" -> "new_terraform_collection2" # forces replacement
        # (6 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_collection   = {
      ~ collection_name = "new_terraform_collection" -> "new_terraform_collection2"
        # (6 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_collection.new_collection: Destroying...
couchbase-capella_collection.new_collection: Destruction complete after 1s
couchbase-capella_collection.new_collection: Creating...
couchbase-capella_collection.new_collection: Creation complete after 0s

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

collections_list = {
  "bucket_id" = "YjE="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "collection_name" = "new_terraform_collection"
      "max_ttl" = 100
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
}
new_collection = {
  "bucket_id" = "YjE="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "new_terraform_collection2"
  "max_ttl" = 100
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
}
```
2. Update maxTTL for server version < 7.6.0 (Error - Not Supported)

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
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection will be updated in-place
  ~ resource "couchbase-capella_collection" "new_collection" {
      ~ max_ttl         = 100 -> 2000
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ collections_list = {
      ~ data            = null -> [
          + {
              + collection_name = "new_terraform_collection"
              + max_ttl         = 100
            },
        ]
        # (5 unchanged attributes hidden)
    }
  ~ new_collection   = {
      ~ max_ttl         = 100 -> 2000
        # (6 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_collection.new_collection: Modifying...
╷
│ Error: Error updating collection
│ 
│   with couchbase-capella_collection.new_collection,
│   on create_collection.tf line 5, in resource "couchbase-capella_collection" "new_collection":
│    5: resource "couchbase-capella_collection" "new_collection" {
│ 
│ Could not update collection for scope"s1": {"hint":"Returned when attempting to modify a collection but the server version is not supported. This operation is only supported for server version 7.6.0 and
│ above.","message":"Unable to modify the collection. Couchbase Server version '7.2.4' is not supported.","code":11014,"httpStatusCode":422}

```
3. Update maxTTL for server version >= 7.6.0 (Update in-place)

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
data.couchbase-capella_collections.existing_collections: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection will be updated in-place
  ~ resource "couchbase-capella_collection" "new_collection" {
      ~ max_ttl         = 500 -> 10000
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ collections_list = {
      ~ data            = null -> [
          + {
              + collection_name = "new_terraform_collection"
              + max_ttl         = 500
            },
        ]
        # (5 unchanged attributes hidden)
    }
  ~ new_collection   = {
      ~ max_ttl         = 500 -> 10000
        # (6 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_collection.new_collection: Modifying...
couchbase-capella_collection.new_collection: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

collections_list = {
  "bucket_id" = "YjE="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "collection_name" = "new_terraform_collection"
      "max_ttl" = 500
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
}
new_collection = {
  "bucket_id" = "YjE="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "new_terraform_collection"
  "max_ttl" = 10000
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
}

```
4. Update name for server version >= 7.6.0 (destroy and replace)

Sample Output:
```
terraform apply 
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
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_collection.new_collection must be replaced
-/+ resource "couchbase-capella_collection" "new_collection" {
      ~ collection_name = "new_terraform_collection" -> "new_terraform_collection2" # forces replacement
        # (6 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ collections_list = {
      ~ data            = [
          ~ {
              ~ max_ttl         = 100 -> 5000
                # (2 unchanged attributes hidden)
            },
        ]
        # (5 unchanged attributes hidden)
    }
  ~ new_collection   = {
      ~ collection_name = "new_terraform_collection" -> "new_terraform_collection2"
        # (6 unchanged attributes hidden)
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
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "collection_name" = "new_terraform_collection"
      "max_ttl" = 5000
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
}
new_collection = {
  "bucket_id" = "YjE="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "new_terraform_collection2"
  "max_ttl" = 5000
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "s1"
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
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "new_terraform_collection2" -> null
      - max_ttl         = 300 -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "s1" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - collections_list = {
      - bucket_id       = "YjE="
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - data            = [
          - {
              - collection_name = "new_terraform_collection2"
              - max_ttl         = 300
            },
          - {
              - collection_name = "c1"
              - max_ttl         = 100
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - scope_name      = "s1"
    } -> null
  - new_collection   = {
      - bucket_id       = "YjE="
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - collection_name = "new_terraform_collection2"
      - max_ttl         = 300
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - scope_name      = "s1"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_collection.new_collection: Destroying...
couchbase-capella_collection.new_collection: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```