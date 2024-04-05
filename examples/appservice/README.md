# Capella App Service Example

This example shows how to create and manage App Services in Capella.

This creates a new app service in the selected Capella cluster. It uses the organization ID, project ID and cluster ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new app service entry in an existing Capella cluster as stated in the `create_app_service.tf` file.
2. UPDATE: Update the app service configuration using Terraform.
3. LIST: List existing app services in Capella as stated in the `list_app_services.tf` file.
4. IMPORT: Import an app services that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created app service from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_app_services.existing_app_services: Reading...
data.capella_app_services.existing_app_services: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_app_service.new_app_service will be created
  + resource "capella_app_service" "new_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe"
      + compute         = {
          + cpu = 4
          + ram = 8
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "new-terraform-app-service"
      + nodes           = 3
      + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + version         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + app_services_list = {
      + data            = [
          + {
              + audit           = {
                  + created_at  = "2023-10-24 14:44:10.59375193 +0000 UTC"
                  + created_by  = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  + modified_at = "2023-10-24 16:37:34.411650647 +0000 UTC"
                  + modified_by = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  + version     = 10
                }
              + cloud_provider  = "AWS"
              + cluster_id      = "b36fefaf-c7a5-42b4-830c-8fa5ef0c3533"
              + compute         = {
                  + cpu = 8
                  + ram = 16
                }
              + current_state   = "turnedOff"
              + description     = ""
              + id              = "80bef88f-d38b-4cb3-9c72-80e34c4da112"
              + name            = "test-appservice"
              + nodes           = 2
              + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              + version         = "3.0.8-1.0.0"
            },
          + {
              + audit           = {
                  + created_at  = "2023-10-16 09:07:20.033905062 +0000 UTC"
                  + created_by  = "78da0474-dfbf-43a2-8704-970d68d985fb"
                  + modified_at = "2023-10-24 16:02:17.52922551 +0000 UTC"
                  + modified_by = "internal-support"
                  + version     = 82
                }
              + cloud_provider  = "AWS"
              + cluster_id      = "59a2a46d-e84f-4541-8f5a-225cee90d302"
              + compute         = {
                  + cpu = 8
                  + ram = 16
                }
              + current_state   = "destroyFailed"
              + description     = ""
              + id              = "b0a70b56-ea7b-4877-af1e-a569571cb058"
              + name            = "test"
              + nodes           = 2
              + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              + version         = "3.0.8-1.0.0"
            },
          + {
              + audit           = {
                  + created_at  = "2023-10-18 10:02:40.466492747 +0000 UTC"
                  + created_by  = "83f6043a-f465-420e-8145-3e9a9264143e"
                  + modified_at = "2023-10-18 18:20:05.677760226 +0000 UTC"
                  + modified_by = "5b372d36-6fe9-4856-ac35-cc0fbc720d9f"
                  + version     = 16
                }
              + cloud_provider  = "AWS"
              + cluster_id      = "9914e3b9-0f70-4d57-83af-707d40b0bee2"
              + compute         = {
                  + cpu = 8
                  + ram = 16
                }
              + current_state   = "healthy"
              + description     = ""
              + id              = "cac46b6a-cc5d-49fb-aa39-d942a2f44603"
              + name            = "test"
              + nodes           = 2
              + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              + version         = "3.0.8-1.0.0"
            },
          + {
              + audit           = {
                  + created_at  = "2023-10-19 07:56:21.845275874 +0000 UTC"
                  + created_by  = "3106bf79-4c39-456d-95cb-5590c1e0458b"
                  + modified_at = "2023-10-19 08:00:02.595921238 +0000 UTC"
                  + modified_by = "3106bf79-4c39-456d-95cb-5590c1e0458b"
                  + version     = 7
                }
              + cloud_provider  = "AWS"
              + cluster_id      = "17618d17-340b-4414-9399-7368836ec250"
              + compute         = {
                  + cpu = 8
                  + ram = 16
                }
              + current_state   = "healthy"
              + description     = ""
              + id              = "d2777afb-5c31-4a35-ab46-38bb9e6be84f"
              + name            = "temp-app-service"
              + nodes           = 2
              + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              + version         = "3.0.8-1.0.0"
            },
        ]
      + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
    }
  + new_app_service   = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe"
      + compute         = {
          + cpu = 4
          + ram = 8
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "new-terraform-app-service"
      + nodes           = 3
      + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
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
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_app_services.existing_app_services: Reading...
data.capella_app_services.existing_app_services: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_app_service.new_app_service will be created
  + resource "capella_app_service" "new_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe"
      + compute         = {
          + cpu = 4
          + ram = 8
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "new-terraform-app-service"
      + nodes           = 3
      + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + version         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service   = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe"
      + compute         = {
          + cpu = 4
          + ram = 8
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "new-terraform-app-service"
      + nodes           = 3
      + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + version         = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_app_service.new_app_service: Creating...
capella_app_service.new_app_service: Still creating... [10s elapsed]
capella_app_service.new_app_service: Still creating... [20s elapsed]
capella_app_service.new_app_service: Still creating... [30s elapsed]
capella_app_service.new_app_service: Still creating... [40s elapsed]
capella_app_service.new_app_service: Still creating... [50s elapsed]
capella_app_service.new_app_service: Still creating... [1m0s elapsed]
capella_app_service.new_app_service: Still creating... [1m10s elapsed]
capella_app_service.new_app_service: Still creating... [1m20s elapsed]
capella_app_service.new_app_service: Still creating... [1m30s elapsed]
capella_app_service.new_app_service: Still creating... [1m40s elapsed]
capella_app_service.new_app_service: Still creating... [1m50s elapsed]
capella_app_service.new_app_service: Still creating... [2m0s elapsed]
capella_app_service.new_app_service: Still creating... [2m10s elapsed]
capella_app_service.new_app_service: Still creating... [2m20s elapsed]
capella_app_service.new_app_service: Still creating... [2m30s elapsed]
capella_app_service.new_app_service: Still creating... [2m40s elapsed]
capella_app_service.new_app_service: Still creating... [2m50s elapsed]
capella_app_service.new_app_service: Still creating... [3m0s elapsed]
capella_app_service.new_app_service: Still creating... [3m10s elapsed]
capella_app_service.new_app_service: Still creating... [3m20s elapsed]
capella_app_service.new_app_service: Still creating... [3m30s elapsed]
capella_app_service.new_app_service: Still creating... [3m40s elapsed]
capella_app_service.new_app_service: Still creating... [3m50s elapsed]
capella_app_service.new_app_service: Creation complete after 3m55s [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

app_services_list = {
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2023-10-24 14:44:10.59375193 +0000 UTC"
        "created_by" = "b6bc2434-721d-42e0-8147-30a2964c9753"
        "modified_at" = "2023-10-24 16:37:34.411650647 +0000 UTC"
        "modified_by" = "b6bc2434-721d-42e0-8147-30a2964c9753"
        "version" = 10
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "b36fefaf-c7a5-42b4-830c-8fa5ef0c3533"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "turnedOff"
      "description" = ""
      "id" = "80bef88f-d38b-4cb3-9c72-80e34c4da112"
      "name" = "test-appservice"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
    {
      "audit" = {
        "created_at" = "2023-10-16 09:07:20.033905062 +0000 UTC"
        "created_by" = "78da0474-dfbf-43a2-8704-970d68d985fb"
        "modified_at" = "2023-10-24 16:02:17.52922551 +0000 UTC"
        "modified_by" = "internal-support"
        "version" = 82
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "59a2a46d-e84f-4541-8f5a-225cee90d302"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "destroyFailed"
      "description" = ""
      "id" = "b0a70b56-ea7b-4877-af1e-a569571cb058"
      "name" = "test"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
    {
      "audit" = {
        "created_at" = "2023-10-18 10:02:40.466492747 +0000 UTC"
        "created_by" = "83f6043a-f465-420e-8145-3e9a9264143e"
        "modified_at" = "2023-10-18 18:20:05.677760226 +0000 UTC"
        "modified_by" = "5b372d36-6fe9-4856-ac35-cc0fbc720d9f"
        "version" = 16
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "9914e3b9-0f70-4d57-83af-707d40b0bee2"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "cac46b6a-cc5d-49fb-aa39-d942a2f44603"
      "name" = "test"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
    {
      "audit" = {
        "created_at" = "2023-10-19 07:56:21.845275874 +0000 UTC"
        "created_by" = "3106bf79-4c39-456d-95cb-5590c1e0458b"
        "modified_at" = "2023-10-19 08:00:02.595921238 +0000 UTC"
        "modified_by" = "3106bf79-4c39-456d-95cb-5590c1e0458b"
        "version" = 7
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "17618d17-340b-4414-9399-7368836ec250"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "d2777afb-5c31-4a35-ab46-38bb9e6be84f"
      "name" = "temp-app-service"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
  ])
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
}
new_app_service = {
  "audit" = {
    "created_at" = "2023-10-24 20:15:49.459612314 +0000 UTC"
    "created_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "modified_at" = "2023-10-24 20:19:36.216609308 +0000 UTC"
    "modified_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "version" = 8
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "b74f5350-f727-427e-8cad-623a691b1cfe"
  "compute" = {
    "cpu" = 4
    "ram" = 8
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 8"
  "id" = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 3
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
  "version" = "3.0.8-1.0.0"
}
```

### Note the App Service ID of the new App Service
Command: `terraform output new_app_service`

Sample Output:
```
{
  "audit" = {
    "created_at" = "2023-10-24 20:15:49.459612314 +0000 UTC"
    "created_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "modified_at" = "2023-10-24 20:19:36.216609308 +0000 UTC"
    "modified_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "version" = 8
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "b74f5350-f727-427e-8cad-623a691b1cfe"
  "compute" = {
    "cpu" = 4
    "ram" = 8
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 8"
  "id" = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 3
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
  "version" = "3.0.8-1.0.0"
}

```
In this case, the app service ID for my new app service is `e57494b4-c791-44b6-8e7a-ee20db89a7f0`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
terraform state list                                  
data.couchbase-capella_app_services.existing_app_services
couchbase-capella_app_service.new_app_service
```

## IMPORT
### Remove the resource `new_app_service` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_service.new_app_service`

Sample Output:
``` 
terraform state rm couchbase-capella_app_service.new_app_service
Removed capella_app_service.new_app_service
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_service.new_app_service id=<appservice_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_app_service.new_app_service id=e57494b4-c791-44b6-8e7a-ee20db89a7f0,cluster_id=b74f5350-f727-427e-8cad-623a691b1cfe,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d`

Sample Output:
``` 
terraform import couchbase-capella_app_service.new_app_service id=e57494b4-c791-44b6-8e7a-ee20db89a7f0,cluster_id=b74f5350-f727-427e-8cad-623a691b1cfe,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d
capella_app_service.new_app_service: Importing from ID "id=e57494b4-c791-44b6-8e7a-ee20db89a7f0,cluster_id=b74f5350-f727-427e-8cad-623a691b1cfe,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d"...
data.capella_app_services.existing_app_services: Reading...
capella_app_service.new_app_service: Import prepared!
  Prepared capella_app_service for import
capella_app_service.new_app_service: Refreshing state... [id=id=e57494b4-c791-44b6-8e7a-ee20db89a7f0,cluster_id=b74f5350-f727-427e-8cad-623a691b1cfe,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d]
data.capella_app_services.existing_app_services: Read complete after 1s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the app service ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which app service belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit the terraform.tfvars file to change the app service configuration settings.

Command: `terraform apply -var-file=terraform.template.tfvars`

``` 
terraform apply -var-file=terraform.template.tfvars
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_app_services.existing_app_services: Reading...
capella_app_service.new_app_service: Refreshing state... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0]
data.capella_app_services.existing_app_services: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_app_service.new_app_service will be updated in-place
  ~ resource "capella_app_service" "new_app_service" {
      ~ audit           = {
          ~ created_at  = "2023-10-24 20:15:49.459612314 +0000 UTC" -> (known after apply)
          ~ created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> (known after apply)
          ~ modified_at = "2023-10-24 20:19:36.216609308 +0000 UTC" -> (known after apply)
          ~ modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> (known after apply)
          ~ version     = 8 -> (known after apply)
        } -> (known after apply)
      ~ cloud_provider  = "AWS" -> (known after apply)
      ~ compute         = {
          ~ cpu = 4 -> 2
          ~ ram = 8 -> 4
        }
      ~ current_state   = "healthy" -> (known after apply)
      ~ etag            = "Version: 8" -> (known after apply)
        id              = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
        name            = "new-terraform-app-service"
      ~ nodes           = 3 -> 2
      ~ version         = "3.0.8-1.0.0" -> (known after apply)
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ app_services_list = {
      ~ data            = [
          + {
              + audit           = {
                  + created_at  = "2023-10-24 14:44:10.59375193 +0000 UTC"
                  + created_by  = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  + modified_at = "2023-10-24 16:37:34.411650647 +0000 UTC"
                  + modified_by = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  + version     = 10
                }
              + cloud_provider  = "AWS"
              + cluster_id      = "b36fefaf-c7a5-42b4-830c-8fa5ef0c3533"
              + compute         = {
                  + cpu = 8
                  + ram = 16
                }
              + current_state   = "turnedOff"
              + description     = ""
              + id              = "80bef88f-d38b-4cb3-9c72-80e34c4da112"
              + name            = "test-appservice"
              + nodes           = 2
              + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              + version         = "3.0.8-1.0.0"
            },
            {
                audit           = {
                    created_at  = "2023-10-16 09:07:20.033905062 +0000 UTC"
                    created_by  = "78da0474-dfbf-43a2-8704-970d68d985fb"
                    modified_at = "2023-10-24 16:02:17.52922551 +0000 UTC"
                    modified_by = "internal-support"
                    version     = 82
                }
                cloud_provider  = "AWS"
                cluster_id      = "59a2a46d-e84f-4541-8f5a-225cee90d302"
                compute         = {
                    cpu = 8
                    ram = 16
                }
                current_state   = "destroyFailed"
                description     = ""
                id              = "b0a70b56-ea7b-4877-af1e-a569571cb058"
                name            = "test"
                nodes           = 2
                organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
                version         = "3.0.8-1.0.0"
            },
            # (2 unchanged elements hidden)
            {
                audit           = {
                    created_at  = "2023-10-24 20:15:49.459612314 +0000 UTC"
                    created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
                    modified_at = "2023-10-24 20:19:36.216609308 +0000 UTC"
                    modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
                    version     = 8
                }
                cloud_provider  = "AWS"
                cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe"
                compute         = {
                    cpu = 4
                    ram = 8
                }
                current_state   = "healthy"
                description     = "My first test app service."
                id              = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
                name            = "new-terraform-app-service"
                nodes           = 3
                organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
                version         = "3.0.8-1.0.0"
            },
          - {
              - audit           = {
                  - created_at  = "2023-10-24 14:44:10.59375193 +0000 UTC"
                  - created_by  = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  - modified_at = "2023-10-24 16:37:34.411650647 +0000 UTC"
                  - modified_by = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  - version     = 10
                }
              - cloud_provider  = "AWS"
              - cluster_id      = "b36fefaf-c7a5-42b4-830c-8fa5ef0c3533"
              - compute         = {
                  - cpu = 8
                  - ram = 16
                }
              - current_state   = "turnedOff"
              - description     = ""
              - id              = "80bef88f-d38b-4cb3-9c72-80e34c4da112"
              - name            = "test-appservice"
              - nodes           = 2
              - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - version         = "3.0.8-1.0.0"
            },
        ]
        # (1 unchanged attribute hidden)
    }
  ~ new_app_service   = {
      ~ audit           = {
          - created_at  = "2023-10-24 20:15:49.459612314 +0000 UTC"
          - created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
          - modified_at = "2023-10-24 20:19:36.216609308 +0000 UTC"
          - modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
          - version     = 8
        } -> (known after apply)
      ~ cloud_provider  = "AWS" -> (known after apply)
      ~ compute         = {
          ~ cpu = 4 -> 2
          ~ ram = 8 -> 4
        }
      ~ current_state   = "healthy" -> (known after apply)
      ~ etag            = "Version: 8" -> (known after apply)
        id              = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
        name            = "new-terraform-app-service"
      ~ nodes           = 3 -> 2
      ~ version         = "3.0.8-1.0.0" -> (known after apply)
        # (5 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_app_service.new_app_service: Modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 10s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 20s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 30s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 40s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 50s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m0s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m10s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m20s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m30s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m40s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m50s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m0s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m10s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m20s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m30s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m40s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m50s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 3m0s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 3m10s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 3m20s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 3m30s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 3m40s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 3m50s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 4m0s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 4m10s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 4m20s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 4m30s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 4m40s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 4m50s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 5m0s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 5m10s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 5m20s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 5m30s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 5m40s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 5m50s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 6m0s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 6m10s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 6m20s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 6m30s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 6m40s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 6m50s elapsed]
capella_app_service.new_app_service: Still modifying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 7m0s elapsed]
capella_app_service.new_app_service: Modifications complete after 7m4s [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

app_services_list = {
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "2023-10-24 14:44:10.59375193 +0000 UTC"
        "created_by" = "b6bc2434-721d-42e0-8147-30a2964c9753"
        "modified_at" = "2023-10-24 16:37:34.411650647 +0000 UTC"
        "modified_by" = "b6bc2434-721d-42e0-8147-30a2964c9753"
        "version" = 10
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "b36fefaf-c7a5-42b4-830c-8fa5ef0c3533"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "turnedOff"
      "description" = ""
      "id" = "80bef88f-d38b-4cb3-9c72-80e34c4da112"
      "name" = "test-appservice"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
    {
      "audit" = {
        "created_at" = "2023-10-16 09:07:20.033905062 +0000 UTC"
        "created_by" = "78da0474-dfbf-43a2-8704-970d68d985fb"
        "modified_at" = "2023-10-24 16:02:17.52922551 +0000 UTC"
        "modified_by" = "internal-support"
        "version" = 82
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "59a2a46d-e84f-4541-8f5a-225cee90d302"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "destroyFailed"
      "description" = ""
      "id" = "b0a70b56-ea7b-4877-af1e-a569571cb058"
      "name" = "test"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
    {
      "audit" = {
        "created_at" = "2023-10-18 10:02:40.466492747 +0000 UTC"
        "created_by" = "83f6043a-f465-420e-8145-3e9a9264143e"
        "modified_at" = "2023-10-18 18:20:05.677760226 +0000 UTC"
        "modified_by" = "5b372d36-6fe9-4856-ac35-cc0fbc720d9f"
        "version" = 16
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "9914e3b9-0f70-4d57-83af-707d40b0bee2"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "cac46b6a-cc5d-49fb-aa39-d942a2f44603"
      "name" = "test"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
    {
      "audit" = {
        "created_at" = "2023-10-19 07:56:21.845275874 +0000 UTC"
        "created_by" = "3106bf79-4c39-456d-95cb-5590c1e0458b"
        "modified_at" = "2023-10-19 08:00:02.595921238 +0000 UTC"
        "modified_by" = "3106bf79-4c39-456d-95cb-5590c1e0458b"
        "version" = 7
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "17618d17-340b-4414-9399-7368836ec250"
      "compute" = {
        "cpu" = 8
        "ram" = 16
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "d2777afb-5c31-4a35-ab46-38bb9e6be84f"
      "name" = "temp-app-service"
      "nodes" = 2
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
    {
      "audit" = {
        "created_at" = "2023-10-24 20:15:49.459612314 +0000 UTC"
        "created_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
        "modified_at" = "2023-10-24 20:19:36.216609308 +0000 UTC"
        "modified_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
        "version" = 8
      }
      "cloud_provider" = "AWS"
      "cluster_id" = "b74f5350-f727-427e-8cad-623a691b1cfe"
      "compute" = {
        "cpu" = 4
        "ram" = 8
      }
      "current_state" = "healthy"
      "description" = "My first test app service."
      "id" = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
      "name" = "new-terraform-app-service"
      "nodes" = 3
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "version" = "3.0.8-1.0.0"
    },
  ])
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
}
new_app_service = {
  "audit" = {
    "created_at" = "2023-10-24 20:15:49.459612314 +0000 UTC"
    "created_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "modified_at" = "2023-10-24 21:20:27.210554955 +0000 UTC"
    "modified_by" = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
    "version" = 14
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "b74f5350-f727-427e-8cad-623a691b1cfe"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 14"
  "id" = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 2
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
  "version" = "3.0.8-1.0.0"
}
```


## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
 
terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_app_services.existing_app_services: Reading...
capella_app_service.new_app_service: Refreshing state... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0]
data.capella_app_services.existing_app_services: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_app_service.new_app_service will be destroyed
  - resource "capella_app_service" "new_app_service" {
      - audit           = {
          - created_at  = "2023-10-24 20:15:49.459612314 +0000 UTC" -> null
          - created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> null
          - modified_at = "2023-10-24 21:20:27.210554955 +0000 UTC" -> null
          - modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu" -> null
          - version     = 14 -> null
        } -> null
      - cloud_provider  = "AWS" -> null
      - cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe" -> null
      - compute         = {
          - cpu = 2 -> null
          - ram = 4 -> null
        } -> null
      - current_state   = "healthy" -> null
      - description     = "My first test app service." -> null
      - etag            = "Version: 14" -> null
      - id              = "e57494b4-c791-44b6-8e7a-ee20db89a7f0" -> null
      - name            = "new-terraform-app-service" -> null
      - nodes           = 2 -> null
      - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d" -> null
      - project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e" -> null
      - version         = "3.0.8-1.0.0" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - app_services_list = {
      - data            = [
          - {
              - audit           = {
                  - created_at  = "2023-10-16 09:07:20.033905062 +0000 UTC"
                  - created_by  = "78da0474-dfbf-43a2-8704-970d68d985fb"
                  - modified_at = "2023-10-24 16:02:17.52922551 +0000 UTC"
                  - modified_by = "internal-support"
                  - version     = 82
                }
              - cloud_provider  = "AWS"
              - cluster_id      = "59a2a46d-e84f-4541-8f5a-225cee90d302"
              - compute         = {
                  - cpu = 8
                  - ram = 16
                }
              - current_state   = "destroyFailed"
              - description     = ""
              - id              = "b0a70b56-ea7b-4877-af1e-a569571cb058"
              - name            = "test"
              - nodes           = 2
              - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - version         = "3.0.8-1.0.0"
            },
          - {
              - audit           = {
                  - created_at  = "2023-10-18 10:02:40.466492747 +0000 UTC"
                  - created_by  = "83f6043a-f465-420e-8145-3e9a9264143e"
                  - modified_at = "2023-10-18 18:20:05.677760226 +0000 UTC"
                  - modified_by = "5b372d36-6fe9-4856-ac35-cc0fbc720d9f"
                  - version     = 16
                }
              - cloud_provider  = "AWS"
              - cluster_id      = "9914e3b9-0f70-4d57-83af-707d40b0bee2"
              - compute         = {
                  - cpu = 8
                  - ram = 16
                }
              - current_state   = "healthy"
              - description     = ""
              - id              = "cac46b6a-cc5d-49fb-aa39-d942a2f44603"
              - name            = "test"
              - nodes           = 2
              - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - version         = "3.0.8-1.0.0"
            },
          - {
              - audit           = {
                  - created_at  = "2023-10-19 07:56:21.845275874 +0000 UTC"
                  - created_by  = "3106bf79-4c39-456d-95cb-5590c1e0458b"
                  - modified_at = "2023-10-19 08:00:02.595921238 +0000 UTC"
                  - modified_by = "3106bf79-4c39-456d-95cb-5590c1e0458b"
                  - version     = 7
                }
              - cloud_provider  = "AWS"
              - cluster_id      = "17618d17-340b-4414-9399-7368836ec250"
              - compute         = {
                  - cpu = 8
                  - ram = 16
                }
              - current_state   = "healthy"
              - description     = ""
              - id              = "d2777afb-5c31-4a35-ab46-38bb9e6be84f"
              - name            = "temp-app-service"
              - nodes           = 2
              - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - version         = "3.0.8-1.0.0"
            },
          - {
              - audit           = {
                  - created_at  = "2023-10-24 20:15:49.459612314 +0000 UTC"
                  - created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
                  - modified_at = "2023-10-24 21:20:27.210554955 +0000 UTC"
                  - modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
                  - version     = 14
                }
              - cloud_provider  = "AWS"
              - cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe"
              - compute         = {
                  - cpu = 2
                  - ram = 4
                }
              - current_state   = "healthy"
              - description     = "My first test app service."
              - id              = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
              - name            = "new-terraform-app-service"
              - nodes           = 2
              - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - version         = "3.0.8-1.0.0"
            },
          - {
              - audit           = {
                  - created_at  = "2023-10-24 14:44:10.59375193 +0000 UTC"
                  - created_by  = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  - modified_at = "2023-10-24 16:37:34.411650647 +0000 UTC"
                  - modified_by = "b6bc2434-721d-42e0-8147-30a2964c9753"
                  - version     = 10
                }
              - cloud_provider  = "AWS"
              - cluster_id      = "b36fefaf-c7a5-42b4-830c-8fa5ef0c3533"
              - compute         = {
                  - cpu = 8
                  - ram = 16
                }
              - current_state   = "turnedOff"
              - description     = ""
              - id              = "80bef88f-d38b-4cb3-9c72-80e34c4da112"
              - name            = "test-appservice"
              - nodes           = 2
              - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - version         = "3.0.8-1.0.0"
            },
        ]
      - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
    } -> null
  - new_app_service   = {
      - audit           = {
          - created_at  = "2023-10-24 20:15:49.459612314 +0000 UTC"
          - created_by  = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
          - modified_at = "2023-10-24 21:20:27.210554955 +0000 UTC"
          - modified_by = "srUTBdMcU5kqzAqibAdSoVf06Szvxbuu"
          - version     = 14
        }
      - cloud_provider  = "AWS"
      - cluster_id      = "b74f5350-f727-427e-8cad-623a691b1cfe"
      - compute         = {
          - cpu = 2
          - ram = 4
        }
      - current_state   = "healthy"
      - description     = "My first test app service."
      - etag            = "Version: 14"
      - id              = "e57494b4-c791-44b6-8e7a-ee20db89a7f0"
      - if_match        = null
      - name            = "new-terraform-app-service"
      - nodes           = 2
      - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
      - project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      - version         = "3.0.8-1.0.0"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_app_service.new_app_service: Destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 10s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 20s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 30s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 40s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 50s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m0s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m10s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m20s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m30s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m40s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 1m50s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m0s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m10s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m20s elapsed]
capella_app_service.new_app_service: Still destroying... [id=e57494b4-c791-44b6-8e7a-ee20db89a7f0, 2m30s elapsed]
capella_app_service.new_app_service: Destruction complete after 2m39s

Destroy complete! Resources: 1 destroyed. 
```
