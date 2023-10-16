# Capella App Service Example

This example shows how to create and manage App Services in Capella.

This creates a new app service in the selected Capella cluster. It uses the organization ID, projectId and clusterId to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Create a new app service entry in an existing Capella cluster as stated in the `create_app_service.tf` file.

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
terraform plan 
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
capella_app_service.new_app_service: Refreshing state... [id=b8f67cbf-b5fd-4111-b291-91c6b729e588]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # capella_app_service.new_app_service has been deleted
  - resource "capella_app_service" "new_app_service" {
      - audit           = {
          - created_at  = "2023-10-16 22:42:56.847307886 +0000 UTC" -> null
          - created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> null
          - modified_at = "2023-10-16 22:47:23.026286895 +0000 UTC" -> null
          - modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> null
          - version     = 7 -> null
        } -> null
      - cloud_provider  = "AWS" -> null
      - cluster_id      = "ebac819f-53af-47df-a302-dd9772f203b0" -> null
      - compute         = {
          - cpu = 2 -> null
          - ram = 4 -> null
        } -> null
      - current_state   = "healthy" -> null
      - description     = "test_nidhi" -> null
      - id              = "b8f67cbf-b5fd-4111-b291-91c6b729e588" -> null
      - name            = "test_appservice" -> null
      - nodes           = 2 -> null
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0" -> null
      - project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e" -> null
      - version         = "3.0.8-1.0.0" -> null
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to
undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_app_service.new_app_service will be created
  + resource "capella_app_service" "new_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "ebac819f-53af-47df-a302-dd9772f203b0"
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "test_nidhi"
      + id              = (known after apply)
      + name            = "test_appservice"
      + nodes           = 2
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + version         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "ebac819f-53af-47df-a302-dd9772f203b0"
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "test_nidhi"
      + id              = (known after apply)
      + name            = "test_appservice"
      + nodes           = 2
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + version         = (known after apply)
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new App Service entry

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
capella_app_service.new_app_service: Refreshing state... [id=b8f67cbf-b5fd-4111-b291-91c6b729e588]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # capella_app_service.new_app_service has been deleted
  - resource "capella_app_service" "new_app_service" {
      - audit           = {
          - created_at  = "2023-10-16 22:42:56.847307886 +0000 UTC" -> null
          - created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> null
          - modified_at = "2023-10-16 22:47:23.026286895 +0000 UTC" -> null
          - modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> null
          - version     = 7 -> null
        } -> null
      - cloud_provider  = "AWS" -> null
      - cluster_id      = "ebac819f-53af-47df-a302-dd9772f203b0" -> null
      - compute         = {
          - cpu = 2 -> null
          - ram = 4 -> null
        } -> null
      - current_state   = "healthy" -> null
      - description     = "test_nidhi" -> null
      - id              = "b8f67cbf-b5fd-4111-b291-91c6b729e588" -> null
      - name            = "test_appservice" -> null
      - nodes           = 2 -> null
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0" -> null
      - project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e" -> null
      - version         = "3.0.8-1.0.0" -> null
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to
undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_app_service.new_app_service will be created
  + resource "capella_app_service" "new_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "ebac819f-53af-47df-a302-dd9772f203b0"
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "test_nidhi"
      + id              = (known after apply)
      + name            = "test_appservice"
      + nodes           = 2
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + version         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "ebac819f-53af-47df-a302-dd9772f203b0"
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "test_nidhi"
      + id              = (known after apply)
      + name            = "test_appservice"
      + nodes           = 2
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + version         = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_app_service.new_app_service: Creating...
capella_app_service.new_app_service: Creation complete after 7s [id=a8cbf864-9cb4-43e2-b942-7c1de0dae988]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_app_service = {
  "audit" = {
    "created_at" = "2023-10-16 22:57:01.906089933 +0000 UTC"
    "created_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "modified_at" = "2023-10-16 22:57:02.278150862 +0000 UTC"
    "modified_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "version" = 2
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "ebac819f-53af-47df-a302-dd9772f203b0"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "deploying"
  "description" = "test_nidhi"
  "id" = "a8cbf864-9cb4-43e2-b942-7c1de0dae988"
  "name" = "test_appservice"
  "nodes" = 2
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
  "version" = "3.0.8-1.0.0"
}

```

