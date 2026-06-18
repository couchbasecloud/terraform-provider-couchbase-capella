# Capella Eventing Function Example

This example shows how to create and manage an eventing function in Capella.

This creates a new eventing function in the selected Capella cluster. It uses the organization ID, project ID and cluster ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new eventing function in an existing Capella cluster as stated in the `create_eventing_function.tf` file.
2. READ: Retrieve an existing eventing function using the `couchbase-capella_eventing_function` data source as stated in the `get_eventing_function.tf` file.
3. LIST: Retrieve all eventing functions in the cluster using the `couchbase-capella_eventing_functions` data source as stated in the `list_eventing_functions.tf` file.
4. UPDATE: Update the eventing function configuration in Capella.
5. DELETE: Delete the newly created eventing function from Capella.
6. IMPORT: Import an eventing function that exists in Capella but not in the terraform state file.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### Create a new eventing function

Command: `terraform apply`

Sample Output:
```
terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_eventing_function.new_eventing_function will be created
  + resource "couchbase-capella_eventing_function" "new_eventing_function" {
      + bindings               = {
          + constants = [
              + {
                  + alias = "maxRetries"
                  + value = "3"
                },
            ]
        }
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + code                   = "function OnUpdate(doc, meta, xattrs){}"
      + description            = "Replicates document mutations to a downstream collection."
      + event_metadata_storage = {
          + bucket     = "metadata"
          + collection = "_default"
          + scope      = "_default"
        }
      + event_source           = {
          + bucket     = "source"
          + collection = "_default"
          + scope      = "_default"
        }
      + name                   = "my_function"
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + settings               = (known after apply)
      + state                  = "deployed"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_eventing_function.new_eventing_function: Creating...
couchbase-capella_eventing_function.new_eventing_function: Creation complete after 3s [name=my_function]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```


## READ
### Retrieve an existing eventing function using the data source

The `couchbase-capella_eventing_function` data source reads a single full eventing function from a cluster by name. This is useful for referencing a function managed elsewhere or for inspecting its current configuration.

Add the following data source block to your configuration, using the name of the function created above:

```
data "couchbase-capella_eventing_function" "existing_function" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = "my_function"
}

output "existing_function_status" {
  value = data.couchbase-capella_eventing_function.existing_function.status
}
```

Command: `terraform apply`

Sample Output:
```
$ terraform apply

data.couchbase-capella_eventing_function.existing_function: Reading...
data.couchbase-capella_eventing_function.existing_function: Read complete after 0s [name=my_function]

Changes to Outputs:
  + existing_function_status = "deployed"

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_function_status = "deployed"
```

The data source also exposes the function's other fields.

## LIST
### Retrieve all eventing functions using the data source

The `couchbase-capella_eventing_functions` data source reads every eventing function in a cluster. You can optionally filter the results by one or more states using the `status` attribute. When `status` is omitted, eventing functions in every state are returned.

Add the following data source block to your configuration:

```
data "couchbase-capella_eventing_functions" "existing_functions" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  # Optional: only return functions in the listed states.
  status = ["deployed", "deploying"]
}

output "existing_functions_count" {
  value = length(data.couchbase-capella_eventing_functions.existing_functions.eventing_functions)
}
```

Command: `terraform apply`

Sample Output:
```
$ terraform apply

data.couchbase-capella_eventing_functions.existing_functions: Reading...
data.couchbase-capella_eventing_functions.existing_functions: Read complete after 1s

Changes to Outputs:
  + existing_functions_count = 1

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_functions_count = 1
```

Each entry in the `eventing_functions` list exposes the same fields as the single eventing function data source.

## UPDATE
### Let us edit the `terraform.tfvars` file to update the eventing function.

To change the function, edit its configuration in `terraform.tfvars`. For example, to undeploy the function, set its `state` to `undeployed`:

```
eventing_function = {
  name        = "my_function"
  description = "Replicates document mutations to a downstream collection."
  state       = "undeployed"

  # ... remaining configuration unchanged
}
```

Command: `terraform apply`

Sample Output:
```
terraform apply

couchbase-capella_eventing_function.new_eventing_function: Refreshing state... [name=my_function]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_eventing_function.new_eventing_function will be updated in-place
  ~ resource "couchbase-capella_eventing_function" "new_eventing_function" {
        name  = "my_function"
      ~ state = "deployed" -> "undeployed"
        # (10 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_eventing_function.new_eventing_function: Modifying... [name=my_function]
couchbase-capella_eventing_function.new_eventing_function: Modifications complete after 2s [name=my_function]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## DESTROY
### Finally, destroy the resources created by Terraform
### Remove the resource from the configuration then apply it.

Command: `terraform apply`

Sample Output:
```
terraform apply

couchbase-capella_eventing_function.new_eventing_function: Refreshing state... [name=my_function]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_eventing_function.new_eventing_function will be destroyed
  - resource "couchbase-capella_eventing_function" "new_eventing_function" {
      - bindings               = {
          - constants = [
              - {
                  - alias = "maxRetries" -> null
                  - value = "3" -> null
                },
            ] -> null
          - urls      = [
              - {
                  - alias                    = "api" -> null
                  - allow_cookies            = true -> null
                  - authentication           = {
                      - type = "bearer" -> null
                    } -> null
                  - url                      = "https://api.example.com" -> null
                  - validate_tls_certificate = true -> null
                },
            ] -> null
        } -> null
      - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - code                   = "function OnUpdate(doc, meta, xattrs){}" -> null
      - description            = "Replicates document mutations to a downstream collection." -> null
      - event_metadata_storage = {
          - bucket     = "metadata" -> null
          - collection = "_default" -> null
          - scope      = "_default" -> null
        } -> null
      - event_source           = {
          - bucket     = "source" -> null
          - collection = "_default" -> null
          - scope      = "_default" -> null
        } -> null
      - name                   = "my_function" -> null
      - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - settings               = {
          - allow_sync_documents   = false -> null
          - cursor_aware           = false -> null
          - feed_boundary          = "from_now" -> null
          - language_compatibility = "7.2.0" -> null
          - max_timer_context_size = 1024 -> null
          - script_timeout         = 60 -> null
          - sql_consistency        = "none" -> null
          - worker_count           = 1 -> null
        } -> null
      - state                  = "undeployed" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_eventing_function.new_eventing_function: Destroying... [name=my_function]
couchbase-capella_eventing_function.new_eventing_function: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```

## IMPORT
### Import an eventing function that exists in Capella but not in the Terraform state file

Command: `terraform import couchbase-capella_eventing_function.new_eventing_function function_name=<function_name>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`


Sample Output:
```
terraform import couchbase-capella_eventing_function.new_eventing_function function_name=f2,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_eventing_function.new_eventing_function: Importing from ID "function_name=f2,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_eventing_function.new_eventing_function: Import prepared!
  Prepared couchbase-capella_eventing_function for import
couchbase-capella_eventing_function.new_eventing_function: Refreshing state... [name=function_name=f2,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

