# Capella Free Tier App Service Example

This example demonstrates how to create a free tier app service using the Couchbase Capella Terraform provider.

This example creates a new free tier app service in the free tier cluster.
It uses the organization ID, project ID and cluster ID to create the app service.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.


# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new free-tier app service entry in an existing Capella cluster as stated in the `create_free_tier_app_service.tf` file.
2. UPDATE: Update the free-tier app service configuration using Terraform.
3. IMPORT: Import a free-tier app service that exists in Capella but not in the terraform state file.
4. DELETE: Delete the newly created free-tier app service from Capella.


## CREATE

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_free_tier_app_service.new_free_tier_app_service will be created
  + resource "couchbase-capella_free_tier_app_service" "new_free_tier_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + compute         = (known after apply)
      + current_state   = (known after apply)
      + description     = "My first test free tier app service"
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "new-terraform-app-service"
      + nodes           = (known after apply)
      + organization_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + plan            = (known after apply)
      + project_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + version         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + appservice_id   = (known after apply)
  + new_app_service = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + compute         = (known after apply)
      + current_state   = (known after apply)
      + description     = "My first test free tier app service"
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "new-terraform-app-service"
      + nodes           = (known after apply)
      + organization_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + plan            = (known after apply)
      + project_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + version         = (known after apply)
    }

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the changes to create the resources

Command: `terraform apply`

Sample Output:
```
$terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_free_tier_app_service.new_free_tier_app_service will be created
  + resource "couchbase-capella_free_tier_app_service" "new_free_tier_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + compute         = (known after apply)
      + current_state   = (known after apply)
      + description     = "My first test free tier app service"
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "new-terraform-free-tier-app-service"
      + nodes           = (known after apply)
      + organization_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + plan            = (known after apply)
      + project_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + version         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + free_tier_app_service_id  = (known after apply)
  + new_free_tier_app_service = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + compute         = (known after apply)
      + current_state   = (known after apply)
      + description     = "My first test free tier app service"
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "new-terraform-free-tier-app-service"
      + nodes           = (known after apply)
      + organization_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + plan            = (known after apply)
      + project_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      + version         = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_free_tier_app_service.new_free_tier_app_service: Creating...
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [30s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [40s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [50s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [1m0s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [1m10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [1m20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [1m30s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [1m40s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [1m50s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [2m0s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [2m10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [2m20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [2m30s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [2m40s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [2m50s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [3m0s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [3m10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [3m20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [3m30s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [3m40s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [3m50s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [4m0s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [4m10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [4m20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [4m30s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [4m40s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [4m50s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [5m0s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still creating... [5m10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Creation complete after 5m16s [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

free_tier_app_service_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
new_free_tier_app_service = {
  "audit" = {
    "created_at" = "2025-04-02 14:15:35.829169232 +0000 UTC"
    "created_by" = "uiQuPitdOw46pKm7R7uIYcDc6mLTiNS0"
    "modified_at" = "2025-04-02 14:20:48.690883036 +0000 UTC"
    "modified_by" = "apikey-uiQuPitdOw46pKm7R7uIYcDc6mLTiNS0"
    "version" = 6
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = "My first test free tier app service"
  "etag" = "Version: 6"
  "id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "name" = "new-terraform-free-tier-app-service"
  "nodes" = 1
  "organization_id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "plan" = "free"
  "project_id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "version" = "3.2.2-1.0.1"
}

```


## UPDATE

### Name and description fields can be updated for free-tier appservice
### In the following example name field is being updated
### Change the name in the terraform configuration

Command: `terraform plan`

Sample Output
```
$terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Refreshing state... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_free_tier_app_service.new_free_tier_app_service will be updated in-place
  ~ resource "couchbase-capella_free_tier_app_service" "new_free_tier_app_service" {
      ~ audit           = {
          ~ created_at  = "2025-03-25 13:53:49.656829741 +0000 UTC" -> (known after apply)
          ~ created_by  = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> (known after apply)
          ~ modified_at = "2025-03-25 14:48:21.181294999 +0000 UTC" -> (known after apply)
          ~ modified_by = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> (known after apply)
          ~ version     = 36 -> (known after apply)
        } -> (known after apply)
      ~ current_state   = "healthy" -> (known after apply)
      ~ etag            = "Version: 36" -> (known after apply)
        id              = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      ~ name            = "new-terraform-app-service" -> "new-terraform-app-service modified"
      ~ plan            = "free" -> (known after apply)
      ~ version         = "3.2.2-1.0.1" -> (known after apply)
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_app_service = {
      ~ audit           = {
          - created_at  = "2025-03-25 13:53:49.656829741 +0000 UTC"
          - created_by  = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
          - modified_at = "2025-03-25 14:48:21.181294999 +0000 UTC"
          - modified_by = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
          - version     = 36
        } -> (known after apply)
      ~ current_state   = "healthy" -> (known after apply)
      ~ etag            = "Version: 36" -> (known after apply)
        id              = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      ~ name            = "new-terraform-app-service" -> "new-terraform-app-service modified"
      ~ plan            = "free" -> (known after apply)
      ~ version         = "3.2.2-1.0.1" -> (known after apply)
        # (7 unchanged attributes hidden)
    }

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the changes to update the resources

Command: `terraform apply`

Sample Output
```
$terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Refreshing state... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_free_tier_app_service.new_free_tier_app_service will be updated in-place
  ~ resource "couchbase-capella_free_tier_app_service" "new_free_tier_app_service" {
      ~ audit           = {
          ~ created_at  = "2025-03-25 13:53:49.656829741 +0000 UTC" -> (known after apply)
          ~ created_by  = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> (known after apply)
          ~ modified_at = "2025-03-25 14:48:21.181294999 +0000 UTC" -> (known after apply)
          ~ modified_by = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> (known after apply)
          ~ version     = 36 -> (known after apply)
        } -> (known after apply)
      ~ current_state   = "healthy" -> (known after apply)
      ~ etag            = "Version: 36" -> (known after apply)
        id              = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      ~ name            = "new-terraform-app-service" -> "new-terraform-app-service modified"
      ~ plan            = "free" -> (known after apply)
      ~ version         = "3.2.2-1.0.1" -> (known after apply)
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_app_service = {
      ~ audit           = {
          - created_at  = "2025-03-25 13:53:49.656829741 +0000 UTC"
          - created_by  = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
          - modified_at = "2025-03-25 14:48:21.181294999 +0000 UTC"
          - modified_by = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
          - version     = 36
        } -> (known after apply)
      ~ current_state   = "healthy" -> (known after apply)
      ~ etag            = "Version: 36" -> (known after apply)
        id              = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      ~ name            = "new-terraform-app-service" -> "new-terraform-app-service modified"
      ~ plan            = "free" -> (known after apply)
      ~ version         = "3.2.2-1.0.1" -> (known after apply)
        # (7 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_free_tier_app_service.new_free_tier_app_service: Modifying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Modifications complete after 5s [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

appservice_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
new_app_service = {
  "audit" = {
    "created_at" = "2025-03-25 13:53:49.656829741 +0000 UTC"
    "created_by" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
    "modified_at" = "2025-03-25 14:50:05.867178045 +0000 UTC"
    "modified_by" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
    "version" = 38
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = ""
  "etag" = "Version: 38"
  "id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "name" = "new-terraform-app-service modified"
  "nodes" = 1
  "organization_id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "plan" = "free"
  "project_id" = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
  "version" = "3.2.2-1.0.1"
}
```


## IMPORT
### Remove the resource `new_free_tier_app_service` from the state file

### First get the resource name from the state
Command: `terraform state list`

Sample Output
```
$terraform state list
couchbase-capella_free_tier_app_service.new_free_tier_app_service
```

### Remove the resource from the state
Command: `terraform state rm <resource_name>`

Sample Output:
```
$terraform state rm couchbase-capella_free_tier_app_service.new_free_tier_app_service
Removed couchbase-capella_free_tier_app_service.new_free_tier_app_service
Successfully removed 1 resource instance(s).
```

### Import the resource to the state file using the following command

Command: `terraform import <resource_name> id=<resource_id>,organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>`

Sample Output:
```
 $terraform import couchbase-capella_free_tier_app_service.new_free_tier_app_service id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,organization_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,project_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,cluster_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Importing from ID "id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,organization_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,project_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,cluster_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd"...
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Import prepared!
  Prepared couchbase-capella_free_tier_app_service for import
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Refreshing state... [id=id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,organization_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,project_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd,cluster_id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

### Verify the state file to check if the resource is imported by using the state list command


## DELETE

### Delete the resources that were created

Command: `terraform destroy`

Sample Output:
```
$terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Refreshing state... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_free_tier_app_service.new_free_tier_app_service will be destroyed
  - resource "couchbase-capella_free_tier_app_service" "new_free_tier_app_service" {
      - audit           = {
          - created_at  = "2025-03-25 13:53:49.656829741 +0000 UTC" -> null
          - created_by  = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> null
          - modified_at = "2025-03-25 14:50:05.867178045 +0000 UTC" -> null
          - modified_by = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> null
          - version     = 38 -> null
        } -> null
      - cloud_provider  = "AWS" -> null
      - cluster_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> null
      - compute         = {
          - cpu = 2 -> null
          - ram = 4 -> null
        } -> null
      - current_state   = "healthy" -> null
      - etag            = "Version: 38" -> null
      - id              = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> null
      - name            = "new-terraform-app-service modified" -> null
      - nodes           = 1 -> null
      - organization_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> null
      - plan            = "free" -> null
      - project_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> null
      - version         = "3.2.2-1.0.1" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - appservice_id   = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd" -> null
  - new_app_service = {
      - audit           = {
          - created_at  = "2025-03-25 13:53:49.656829741 +0000 UTC"
          - created_by  = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
          - modified_at = "2025-03-25 14:50:05.867178045 +0000 UTC"
          - modified_by = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
          - version     = 38
        }
      - cloud_provider  = "AWS"
      - cluster_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      - compute         = {
          - cpu = 2
          - ram = 4
        }
      - current_state   = "healthy"
      - description     = ""
      - etag            = "Version: 38"
      - id              = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      - name            = "new-terraform-app-service modified"
      - nodes           = 1
      - organization_id = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      - plan            = "free"
      - project_id      = "aaaaaaaa-bbbbbbbb-cccccccc-dddddd"
      - version         = "3.2.2-1.0.1"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_free_tier_app_service.new_free_tier_app_service: Destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 30s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 40s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 50s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 1m0s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 1m10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 1m20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 1m30s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 1m40s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 1m50s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 2m0s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 2m10s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Still destroying... [id=aaaaaaaa-bbbbbbbb-cccccccc-dddddd, 2m20s elapsed]
couchbase-capella_free_tier_app_service.new_free_tier_app_service: Destruction complete after 2m23s

Destroy complete! Resources: 1 destroyed.
```
