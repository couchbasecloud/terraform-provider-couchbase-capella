# End to End Demo for Private Preview

In this demo, we will perform the following operations:

1. Fetch organization details based on the `organization_id`
2. Create a new API Key in the organization.
3. Invite a new user to the organization.
4. Create a new project in the organization. 
5. Create a new AWS cluster in the organization and the newly created project.
6. Store the cluster certificate in an output variable. 
7. Create a new database credential in the newly created cluster. 
8. Add a new allowlist to the cluster. 
9. Create a new bucket in the cluster.
10. Create a new app service in the cluster.

## Pre-Requisites:

Make sure you have followed the pre-requisite steps from the parent readme that builds the binary and stores the path in a .terraformrc file.

### Variables file.

- Copy the `terraform.template.tfvars` file to `terraform.tfvars` file in the same directory
- Create 1 V4 API key in your organization using the Capella UI.
- Replace all placeholders with actual values. Use the above created API Key secret as the value for `auth_token`

## Execution

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_organization.existing_organization: Reading...
capella_project.new_project: Refreshing state... [id=7a15caa2-47c3-4cbd-a293-4befde7c284f]
data.capella_organization.existing_organization: Read complete after 1s [name=capella-prod]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # capella_project.new_project has been deleted
  - resource "capella_project" "new_project" {
      - id              = "7a15caa2-47c3-4cbd-a293-4befde7c284f" -> null
      - name            = "My First Terraform Project" -> null
        # (4 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to
undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.capella_certificate.existing_certificate will be read during apply
  # (config refers to values not yet known)
 <= data "capella_certificate" "existing_certificate" {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
    }

  # capella_allowlist.new_allowlist will be created
  + resource "capella_allowlist" "new_allowlist" {
      + audit           = (known after apply)
      + cidr            = "73.222.28.124/32"
      + cluster_id      = (known after apply)
      + comment         = "Allow access from a public IP"
      + expires_at      = "2023-11-30T23:59:59.465Z"
      + id              = (known after apply)
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
    }

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.43.0/23",
          + "73.222.28.124/32",
        ]
      + audit              = (known after apply)
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "My First Terraform API Key"
      + organization_id    = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectManager",
                  + "projectDataReader",
                ]
              + type  = "project"
            },
        ]
      + rotate             = (known after apply)
      + secret             = (sensitive value)
      + token              = (sensitive value)
    }

  # capella_app_service.new_app_service will be created
  + resource "capella_app_service" "new_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = (known after apply)
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "new-terraform-app-service"
      + nodes           = 2
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = (known after apply)
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id                 = (known after apply)
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # capella_cluster.new_cluster will be created
  + resource "capella_cluster" "new_cluster" {
      + app_service_id   = (known after apply)
      + audit            = (known after apply)
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "10.10.30.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + couchbase_server = {
          + version = "7.1"
        }
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + name             = "My First Terraform Cluster"
      + organization_id  = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id       = (known after apply)
      + service_groups   = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + iops    = 5000
                      + storage = 50
                      + type    = "io2"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                  + "index",
                  + "query",
                ]
            },
        ]
      + support          = {
          + plan     = "developer pro"
          + timezone = "PT"
        }
    }

  # capella_database_credential.new_database_credential will be created
  + resource "capella_database_credential" "new_database_credential" {
      + access          = [
          + {
              + privileges = [
                  + "data_reader",
                  + "data_writer",
                ]
            },
        ]
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + id              = (known after apply)
      + name            = "terraform_db_credential"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + password        = (sensitive value)
      + project_id      = (known after apply)
    }

  # capella_project.new_project will be created
  + resource "capella_project" "new_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
    }

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Plan: 8 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikey              = (sensitive value)
  + app_service         = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = (known after apply)
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "new-terraform-app-service"
      + nodes           = 2
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }
  + bucket              = "new_terraform_bucket"
  + certificate         = {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
    }
  + cluster             = {
      + app_service_id   = (known after apply)
      + audit            = (known after apply)
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "10.10.30.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + couchbase_server = {
          + version = "7.1"
        }
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + if_match         = null
      + name             = "My First Terraform Cluster"
      + organization_id  = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id       = (known after apply)
      + service_groups   = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + iops    = 5000
                      + storage = 50
                      + type    = "io2"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                  + "index",
                  + "query",
                ]
            },
        ]
      + support          = {
          + plan     = "developer pro"
          + timezone = "PT"
        }
    }
  + database_credential = (sensitive value)
  + organization        = {
      + audit           = {
          + created_at  = "2021-12-03 16:14:45.105347711 +0000 UTC"
          + created_by  = "b1cf2366-0401-4cac-8770-f24e511f6c0a"
          + modified_at = "2023-07-07 10:17:11.682962228 +0000 UTC"
          + modified_by = "8578bc6d-e67d-44a3-80f0-d72a5cfbebc6"
          + version     = 28
        }
      + description     = ""
      + name            = "capella-prod"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + preferences     = {
          + session_duration = 7200
        }
    }
  + project             = "My First Terraform Project"
  + user                = {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_organization.existing_organization: Reading...
capella_project.new_project: Refreshing state... [id=7a15caa2-47c3-4cbd-a293-4befde7c284f]
data.capella_organization.existing_organization: Read complete after 0s [name=capella-prod]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # capella_project.new_project has been deleted
  - resource "capella_project" "new_project" {
      - id              = "7a15caa2-47c3-4cbd-a293-4befde7c284f" -> null
      - name            = "My First Terraform Project" -> null
        # (4 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to
undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.capella_certificate.existing_certificate will be read during apply
  # (config refers to values not yet known)
 <= data "capella_certificate" "existing_certificate" {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
    }

  # capella_allowlist.new_allowlist will be created
  + resource "capella_allowlist" "new_allowlist" {
      + audit           = (known after apply)
      + cidr            = "73.222.28.124/32"
      + cluster_id      = (known after apply)
      + comment         = "Allow access from a public IP"
      + expires_at      = "2023-11-30T23:59:59.465Z"
      + id              = (known after apply)
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
    }

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.43.0/23",
          + "73.222.28.124/32",
        ]
      + audit              = (known after apply)
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "My First Terraform API Key"
      + organization_id    = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectManager",
                  + "projectDataReader",
                ]
              + type  = "project"
            },
        ]
      + rotate             = (known after apply)
      + secret             = (sensitive value)
      + token              = (sensitive value)
    }

  # capella_app_service.new_app_service will be created
  + resource "capella_app_service" "new_app_service" {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = (known after apply)
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "new-terraform-app-service"
      + nodes           = 2
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = (known after apply)
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id                 = (known after apply)
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # capella_cluster.new_cluster will be created
  + resource "capella_cluster" "new_cluster" {
      + app_service_id   = (known after apply)
      + audit            = (known after apply)
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "10.10.30.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + couchbase_server = {
          + version = "7.1"
        }
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + name             = "My First Terraform Cluster"
      + organization_id  = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id       = (known after apply)
      + service_groups   = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + iops    = 5000
                      + storage = 50
                      + type    = "io2"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                  + "index",
                  + "query",
                ]
            },
        ]
      + support          = {
          + plan     = "developer pro"
          + timezone = "PT"
        }
    }

  # capella_database_credential.new_database_credential will be created
  + resource "capella_database_credential" "new_database_credential" {
      + access          = [
          + {
              + privileges = [
                  + "data_reader",
                  + "data_writer",
                ]
            },
        ]
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + id              = (known after apply)
      + name            = "terraform_db_credential"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + password        = (sensitive value)
      + project_id      = (known after apply)
    }

  # capella_project.new_project will be created
  + resource "capella_project" "new_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
    }

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Plan: 8 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikey              = (sensitive value)
  + app_service         = {
      + audit           = (known after apply)
      + cloud_provider  = (known after apply)
      + cluster_id      = (known after apply)
      + compute         = {
          + cpu = 2
          + ram = 4
        }
      + current_state   = (known after apply)
      + description     = "My first test app service."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "new-terraform-app-service"
      + nodes           = 2
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }
  + bucket              = "new_terraform_bucket"
  + certificate         = {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = (known after apply)
    }
  + cluster             = {
      + app_service_id   = (known after apply)
      + audit            = (known after apply)
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "10.10.30.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + couchbase_server = {
          + version = "7.1"
        }
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + if_match         = null
      + name             = "My First Terraform Cluster"
      + organization_id  = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id       = (known after apply)
      + service_groups   = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + iops    = 5000
                      + storage = 50
                      + type    = "io2"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                  + "index",
                  + "query",
                ]
            },
        ]
      + support          = {
          + plan     = "developer pro"
          + timezone = "PT"
        }
    }
  + database_credential = (sensitive value)
  + organization        = {
      + audit           = {
          + created_at  = "2021-12-03 16:14:45.105347711 +0000 UTC"
          + created_by  = "b1cf2366-0401-4cac-8770-f24e511f6c0a"
          + modified_at = "2023-07-07 10:17:11.682962228 +0000 UTC"
          + modified_by = "8578bc6d-e67d-44a3-80f0-d72a5cfbebc6"
          + version     = 28
        }
      + description     = ""
      + name            = "capella-prod"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + preferences     = {
          + session_duration = 7200
        }
    }
  + project             = "My First Terraform Project"
  + user                = {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectViewer",
                ]
              + type  = "project"
            },
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_project.new_project: Creating...
capella_project.new_project: Creation complete after 6s [id=fda8057c-c0ee-4b81-94e4-f9a3d6ca6449]
capella_user.new_user: Creating...
capella_apikey.new_apikey: Creating...
capella_cluster.new_cluster: Creating...
capella_apikey.new_apikey: Creation complete after 1s [id=iCInIRNFQVGitmFQAjf1JD2FKf8PaLCD]
capella_user.new_user: Creation complete after 8s [id=0f451c57-2a96-4a59-b8b1-d2b144411598]
capella_cluster.new_cluster: Still creating... [10s elapsed]
capella_cluster.new_cluster: Still creating... [20s elapsed]
capella_cluster.new_cluster: Still creating... [30s elapsed]
capella_cluster.new_cluster: Still creating... [40s elapsed]
capella_cluster.new_cluster: Still creating... [50s elapsed]
capella_cluster.new_cluster: Still creating... [1m0s elapsed]
capella_cluster.new_cluster: Still creating... [1m10s elapsed]
capella_cluster.new_cluster: Still creating... [1m20s elapsed]
capella_cluster.new_cluster: Still creating... [1m30s elapsed]
capella_cluster.new_cluster: Still creating... [1m40s elapsed]
capella_cluster.new_cluster: Still creating... [1m50s elapsed]
capella_cluster.new_cluster: Still creating... [2m0s elapsed]
capella_cluster.new_cluster: Still creating... [2m10s elapsed]
capella_cluster.new_cluster: Still creating... [2m20s elapsed]
capella_cluster.new_cluster: Still creating... [2m30s elapsed]
capella_cluster.new_cluster: Still creating... [2m40s elapsed]
capella_cluster.new_cluster: Still creating... [2m50s elapsed]
capella_cluster.new_cluster: Creation complete after 2m54s [id=16f91637-6162-49c8-950f-b53bb533040f]
data.capella_certificate.existing_certificate: Reading...
capella_app_service.new_app_service: Creating...
capella_allowlist.new_allowlist: Creating...
capella_bucket.new_bucket: Creating...
capella_database_credential.new_database_credential: Creating...
data.capella_certificate.existing_certificate: Read complete after 0s
capella_allowlist.new_allowlist: Creation complete after 1s [id=b6cfe2a8-fbb8-4bba-ba38-b26a29974767]
capella_bucket.new_bucket: Creation complete after 6s [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
capella_database_credential.new_database_credential: Creation complete after 8s [id=62ad5f82-8100-4d6e-9216-70e1f82d0dbe]
capella_app_service.new_app_service: Still creating... [10s elapsed]
capella_app_service.new_app_service: Still creating... [20s elapsed]
capella_app_service.new_app_service: Still creating... [30s elapsed]
capella_app_service.new_app_service: Still creating... [40s elapsed]
capella_app_service.new_app_service: Still creating... [50s elapsed]
.
.
.
capella_app_service.new_app_service: Still creating... [23m0s elapsed]
capella_app_service.new_app_service: Still creating... [23m10s elapsed]
capella_app_service.new_app_service: Creation complete after 23m19s [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf]

Apply complete! Resources: 8 added, 0 changed, 0 destroyed.

Outputs:

apikey = <sensitive>
app_service = {
  "audit" = {
    "created_at" = "2023-10-24 21:49:18.368424839 +0000 UTC"
    "created_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "modified_at" = "2023-10-24 22:12:27.892184801 +0000 UTC"
    "modified_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "version" = 7
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "16f91637-6162-49c8-950f-b53bb533040f"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 7"
  "id" = "51b402e6-aa55-40bf-82ec-aae5ccf257cf"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 2
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
  "version" = "3.0.8-1.0.0"
}
bucket = "new_terraform_bucket"
certificate = {
  "cluster_id" = "16f91637-6162-49c8-950f-b53bb533040f"
  "data" = tolist([
    {
      "certificate" = <<-EOT
      -----BEGIN CERTIFICATE-----
      MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
      JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEyMDYy
      MjEyNTlaFw0yOTEyMDYyMzEyNTlaMCQxEjAQBgNVBAoMCUNvdWNoYmFzZTEOMAwG
      A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfvOIi
      enG4Dp+hJu9asdxEMRmH70hDyMXv5ZjBhbo39a42QwR59y/rC/sahLLQuNwqif85
      Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
      WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
      rL3ooAw/i0uxr/TQSaxi1l8t4iMt4/gU+W52+8Yol0JbXBTFX6itg62ppb/Eugmn
      mQRMgL67ccZs7cJ9/A0wlXencX2ohZQOR3mtknfol3FH4+glQFn27Q4xBCzVkY9j
      KQ20T1LgmGSngBInAgMBAAGjQjBAMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYE
      FJQOBPvrkU2In1Sjoxt97Xy8+cKNMA4GA1UdDwEB/wQEAwIBhjANBgkqhkiG9w0B
      AQsFAAOCAQEARgM6XwcXPLSpFdSf0w8PtpNGehmdWijPM3wHb7WZiS47iNen3oq8
      m2mm6V3Z57wbboPpfI+VEzbhiDcFfVnK1CXMC0tkF3fnOG1BDDvwt4jU95vBiNjY
      xdzlTP/Z+qr0cnVbGBSZ+fbXstSiRaaAVcqQyv3BRvBadKBkCyPwo+7svQnScQ5P
      Js7HEHKVms5tZTgKIw1fbmgR2XHleah1AcANB+MAPBCcTgqurqr5G7W2aPSBLLGA
      fRIiVzm7VFLc7kWbp7ENH39HVG6TZzKnfl9zJYeiklo5vQQhGSMhzBsO70z4RRzi
      DPFAN/4qZAgD5q3AFNIq2WWADFQGSwVJhg==
      -----END CERTIFICATE-----
      EOT
    },
  ])
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
}
cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-10-24 21:46:20.62652889 +0000 UTC"
    "created_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "modified_at" = "2023-10-24 21:49:08.817024833 +0000 UTC"
    "modified_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "version" = 5
  }
  "availability" = {
    "type" = "multi"
  }
  "cloud_provider" = {
    "cidr" = "10.10.30.0/23"
    "region" = "us-east-1"
    "type" = "aws"
  }
  "couchbase_server" = {
    "version" = "7.1"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "16f91637-6162-49c8-950f-b53bb533040f"
  "if_match" = tostring(null)
  "name" = "My First Terraform Cluster"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
  "service_groups" = tolist([
    {
      "node" = {
        "compute" = {
          "cpu" = 4
          "ram" = 16
        }
        "disk" = {
          "iops" = 5000
          "storage" = 50
          "type" = "io2"
        }
      }
      "num_of_nodes" = 3
      "services" = tolist([
        "data",
        "index",
        "query",
      ])
    },
  ])
  "support" = {
    "plan" = "developer pro"
    "timezone" = "PT"
  }
}
database_credential = <sensitive>
organization = {
  "audit" = {
    "created_at" = "2021-12-03 16:14:45.105347711 +0000 UTC"
    "created_by" = "b1cf2366-0401-4cac-8770-f24e511f6c0a"
    "modified_at" = "2023-07-07 10:17:11.682962228 +0000 UTC"
    "modified_by" = "8578bc6d-e67d-44a3-80f0-d72a5cfbebc6"
    "version" = 28
  }
  "description" = ""
  "name" = "capella-prod"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "preferences" = {
    "session_duration" = 7200
  }
}
project = "My First Terraform Project"
user = {
  "audit" = {
    "created_at" = "2023-10-24 20:55:30.005741097 +0000 UTC"
    "created_by" = "0f451c57-2a96-4a59-b8b1-d2b144411598"
    "modified_at" = "2023-10-24 20:55:30.005741097 +0000 UTC"
    "modified_by" = "0f451c57-2a96-4a59-b8b1-d2b144411598"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-22T20:55:30.005741367Z"
  "id" = "0f451c57-2a96-4a59-b8b1-d2b144411598"
  "inactive" = true
  "last_login" = ""
  "name" = "John Doe"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
      "roles" = tolist([
        "projectViewer",
      ])
      "type" = "project"
    },
    {
      "id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
      "roles" = tolist([
        "projectDataReaderWriter",
      ])
      "type" = "project"
    },
  ])
  "status" = "not-verified"
  "time_zone" = ""
}
```

Finally, you can view the outputs using the `terraform output` command

Sample Output:
```
$ terraform output
apikey = <sensitive>
app_service = {
  "audit" = {
    "created_at" = "2023-10-24 21:49:18.368424839 +0000 UTC"
    "created_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "modified_at" = "2023-10-24 22:12:27.892184801 +0000 UTC"
    "modified_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "version" = 7
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "16f91637-6162-49c8-950f-b53bb533040f"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 7"
  "id" = "51b402e6-aa55-40bf-82ec-aae5ccf257cf"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 2
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
  "version" = "3.0.8-1.0.0"
}
bucket = "new_terraform_bucket"
certificate = {
  "cluster_id" = "16f91637-6162-49c8-950f-b53bb533040f"
  "data" = tolist([
    {
      "certificate" = <<-EOT
      -----BEGIN CERTIFICATE-----
      MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
      JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEyMDYy
      MjEyNTlaFw0yOTEyMDYyMzEyNTlaMCQxEjAQBgNVBAoMCUNvdWNoYmFzZTEOMAwG
      A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfvOIi
      enG4Dp+hJu9asdxEMRmH70hDyMXv5ZjBhbo39a42QwR59y/rC/sahLLQuNwqif85
      Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
      WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
      rL3ooAw/i0uxr/TQSaxi1l8t4iMt4/gU+W52+8Yol0JbXBTFX6itg62ppb/Eugmn
      mQRMgL67ccZs7cJ9/A0wlXencX2ohZQOR3mtknfol3FH4+glQFn27Q4xBCzVkY9j
      KQ20T1LgmGSngBInAgMBAAGjQjBAMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYE
      FJQOBPvrkU2In1Sjoxt97Xy8+cKNMA4GA1UdDwEB/wQEAwIBhjANBgkqhkiG9w0B
      AQsFAAOCAQEARgM6XwcXPLSpFdSf0w8PtpNGehmdWijPM3wHb7WZiS47iNen3oq8
      m2mm6V3Z57wbboPpfI+VEzbhiDcFfVnK1CXMC0tkF3fnOG1BDDvwt4jU95vBiNjY
      xdzlTP/Z+qr0cnVbGBSZ+fbXstSiRaaAVcqQyv3BRvBadKBkCyPwo+7svQnScQ5P
      Js7HEHKVms5tZTgKIw1fbmgR2XHleah1AcANB+MAPBCcTgqurqr5G7W2aPSBLLGA
      fRIiVzm7VFLc7kWbp7ENH39HVG6TZzKnfl9zJYeiklo5vQQhGSMhzBsO70z4RRzi
      DPFAN/4qZAgD5q3AFNIq2WWADFQGSwVJhg==
      -----END CERTIFICATE-----
      EOT
    },
  ])
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
}
cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-10-24 21:46:20.62652889 +0000 UTC"
    "created_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "modified_at" = "2023-10-24 21:49:08.817024833 +0000 UTC"
    "modified_by" = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
    "version" = 5
  }
  "availability" = {
    "type" = "multi"
  }
  "cloud_provider" = {
    "cidr" = "10.10.30.0/23"
    "region" = "us-east-1"
    "type" = "aws"
  }
  "couchbase_server" = {
    "version" = "7.1"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "16f91637-6162-49c8-950f-b53bb533040f"
  "if_match" = tostring(null)
  "name" = "My First Terraform Cluster"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
  "service_groups" = tolist([
    {
      "node" = {
        "compute" = {
          "cpu" = 4
          "ram" = 16
        }
        "disk" = {
          "iops" = 5000
          "storage" = 50
          "type" = "io2"
        }
      }
      "num_of_nodes" = 3
      "services" = tolist([
        "data",
        "index",
        "query",
      ])
    },
  ])
  "support" = {
    "plan" = "developer pro"
    "timezone" = "PT"
  }
}
database_credential = <sensitive>
organization = {
  "audit" = {
    "created_at" = "2021-12-03 16:14:45.105347711 +0000 UTC"
    "created_by" = "b1cf2366-0401-4cac-8770-f24e511f6c0a"
    "modified_at" = "2023-07-07 10:17:11.682962228 +0000 UTC"
    "modified_by" = "8578bc6d-e67d-44a3-80f0-d72a5cfbebc6"
    "version" = 28
  }
  "description" = ""
  "name" = "capella-prod"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "preferences" = {
    "session_duration" = 7200
  }
}
project = "My First Terraform Project"
user = {
  "audit" = {
    "created_at" = "2023-10-24 20:55:30.005741097 +0000 UTC"
    "created_by" = "0f451c57-2a96-4a59-b8b1-d2b144411598"
    "modified_at" = "2023-10-24 20:55:30.005741097 +0000 UTC"
    "modified_by" = "0f451c57-2a96-4a59-b8b1-d2b144411598"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-22T20:55:30.005741367Z"
  "id" = "0f451c57-2a96-4a59-b8b1-d2b144411598"
  "inactive" = true
  "last_login" = ""
  "name" = "John Doe"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
      "roles" = tolist([
        "projectViewer",
      ])
      "type" = "project"
    },
    {
      "id" = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
      "roles" = tolist([
        "projectDataReaderWriter",
      ])
      "type" = "project"
    },
  ])
  "status" = "not-verified"
  "time_zone" = ""
}
```

All these resources can be destroyed using the `terraform destroy` command

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
data.capella_organization.existing_organization: Reading...
capella_project.new_project: Refreshing state... [id=fda8057c-c0ee-4b81-94e4-f9a3d6ca6449]
data.capella_organization.existing_organization: Read complete after 0s [name=capella-prod]
capella_apikey.new_apikey: Refreshing state... [id=iCInIRNFQVGitmFQAjf1JD2FKf8PaLCD]
capella_cluster.new_cluster: Refreshing state... [id=16f91637-6162-49c8-950f-b53bb533040f]
capella_database_credential.new_database_credential: Refreshing state... [id=62ad5f82-8100-4d6e-9216-70e1f82d0dbe]
data.capella_certificate.existing_certificate: Reading...
capella_allowlist.new_allowlist: Refreshing state... [id=b6cfe2a8-fbb8-4bba-ba38-b26a29974767]
capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
capella_app_service.new_app_service: Refreshing state... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf]
data.capella_certificate.existing_certificate: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_allowlist.new_allowlist will be destroyed
  - resource "capella_allowlist" "new_allowlist" {
      - audit           = {
          - created_at  = "2023-10-24 21:49:12.340054546 +0000 UTC" -> null
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - modified_at = "2023-10-24 21:49:12.340054546 +0000 UTC" -> null
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - version     = 1 -> null
        } -> null
      - cidr            = "73.222.28.124/32" -> null
      - cluster_id      = "16f91637-6162-49c8-950f-b53bb533040f" -> null
      - comment         = "Allow access from a public IP" -> null
      - expires_at      = "2023-11-30T23:59:59.465Z" -> null
      - id              = "b6cfe2a8-fbb8-4bba-ba38-b26a29974767" -> null
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - project_id      = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449" -> null
    }

  # capella_apikey.new_apikey will be destroyed
  - resource "capella_apikey" "new_apikey" {
      - allowed_cidrs      = [
          - "10.1.42.0/23",
          - "10.1.43.0/23",
          - "73.222.28.124/32",
        ] -> null
      - audit              = {
          - created_at  = "2023-10-24 21:46:18.54425008 +0000 UTC" -> null
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - modified_at = "2023-10-24 21:46:18.54425008 +0000 UTC" -> null
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - version     = 1 -> null
        } -> null
      - expiry             = 180 -> null
      - id                 = "iCInIRNFQVGitmFQAjf1JD2FKf8PaLCD" -> null
      - name               = "My First Terraform API Key" -> null
      - organization_id    = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - organization_roles = [
          - "organizationMember",
        ] -> null
      - resources          = [
          - {
              - id    = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449" -> null
              - roles = [
                  - "projectManager",
                  - "projectDataReader",
                ] -> null
              - type  = "project" -> null
            },
        ] -> null
      - token              = (sensitive value) -> null
    }

  # capella_app_service.new_app_service will be destroyed
  - resource "capella_app_service" "new_app_service" {
      - audit           = {
          - created_at  = "2023-10-24 21:49:18.368424839 +0000 UTC" -> null
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - modified_at = "2023-10-24 22:12:27.892184801 +0000 UTC" -> null
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - version     = 7 -> null
        } -> null
      - cloud_provider  = "AWS" -> null
      - cluster_id      = "16f91637-6162-49c8-950f-b53bb533040f" -> null
      - compute         = {
          - cpu = 2 -> null
          - ram = 4 -> null
        } -> null
      - current_state   = "healthy" -> null
      - description     = "My first test app service." -> null
      - etag            = "Version: 7" -> null
      - id              = "51b402e6-aa55-40bf-82ec-aae5ccf257cf" -> null
      - name            = "new-terraform-app-service" -> null
      - nodes           = 2 -> null
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - project_id      = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449" -> null
      - version         = "3.0.8-1.0.0" -> null
    }

  # capella_bucket.new_bucket will be destroyed
  - resource "capella_bucket" "new_bucket" {
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "16f91637-6162-49c8-950f-b53bb533040f" -> null
      - durability_level           = "none" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = false -> null
      - id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> null
      - memory_allocation_in_mb    = 100 -> null
      - name                       = "new_terraform_bucket" -> null
      - organization_id            = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - project_id                 = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449" -> null
      - replicas                   = 1 -> null
      - stats                      = {
          - disk_used_in_mib   = 18 -> null
          - item_count         = 0 -> null
          - memory_used_in_mib = 54 -> null
          - ops_per_second     = 0 -> null
        } -> null
      - storage_backend            = "couchstore" -> null
      - time_to_live_in_seconds    = 0 -> null
      - type                       = "couchbase" -> null
    }

  # capella_cluster.new_cluster will be destroyed
  - resource "capella_cluster" "new_cluster" {
      - audit            = {
          - created_at  = "2023-10-24 21:46:20.62652889 +0000 UTC" -> null
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - modified_at = "2023-10-24 21:49:08.817024833 +0000 UTC" -> null
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - version     = 5 -> null
        } -> null
      - availability     = {
          - type = "multi" -> null
        } -> null
      - cloud_provider   = {
          - cidr   = "10.10.30.0/23" -> null
          - region = "us-east-1" -> null
          - type   = "aws" -> null
        } -> null
      - couchbase_server = {
          - version = "7.1" -> null
        } -> null
      - current_state    = "healthy" -> null
      - description      = "My first test cluster for multiple services." -> null
      - etag             = "Version: 5" -> null
      - id               = "16f91637-6162-49c8-950f-b53bb533040f" -> null
      - name             = "My First Terraform Cluster" -> null
      - organization_id  = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - project_id       = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449" -> null
      - service_groups   = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 4 -> null
                      - ram = 16 -> null
                    } -> null
                  - disk    = {
                      - iops    = 5000 -> null
                      - storage = 50 -> null
                      - type    = "io2" -> null
                    } -> null
                } -> null
              - num_of_nodes = 3 -> null
              - services     = [
                  - "data",
                  - "index",
                  - "query",
                ] -> null
            },
        ] -> null
      - support          = {
          - plan     = "developer pro" -> null
          - timezone = "PT" -> null
        } -> null
    }

  # capella_database_credential.new_database_credential will be destroyed
  - resource "capella_database_credential" "new_database_credential" {
      - access          = [
          - {
              - privileges = [
                  - "data_reader",
                  - "data_writer",
                ] -> null
            },
        ] -> null
      - audit           = {
          - created_at  = "2023-10-24 21:49:18.320964774 +0000 UTC" -> null
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - modified_at = "2023-10-24 21:49:18.320964774 +0000 UTC" -> null
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - version     = 1 -> null
        } -> null
      - cluster_id      = "16f91637-6162-49c8-950f-b53bb533040f" -> null
      - id              = "62ad5f82-8100-4d6e-9216-70e1f82d0dbe" -> null
      - name            = "terraform_db_credential" -> null
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - password        = (sensitive value) -> null
      - project_id      = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449" -> null
    }

  # capella_project.new_project will be destroyed
  - resource "capella_project" "new_project" {
      - audit           = {
          - created_at  = "2023-10-24 21:46:11.822742826 +0000 UTC" -> null
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - modified_at = "2023-10-24 21:46:11.822753235 +0000 UTC" -> null
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV" -> null
          - version     = 1 -> null
        } -> null
      - description     = "A Capella Project that will host many Capella clusters." -> null
      - etag            = "Version: 1" -> null
      - id              = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449" -> null
      - name            = "My First Terraform Project" -> null
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
    }

Plan: 0 to add, 0 to change, 7 to destroy.

Changes to Outputs:
  - apikey              = (sensitive value) -> null
  - app_service         = {
      - audit           = {
          - created_at  = "2023-10-24 21:49:18.368424839 +0000 UTC"
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
          - modified_at = "2023-10-24 22:12:27.892184801 +0000 UTC"
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
          - version     = 7
        }
      - cloud_provider  = "AWS"
      - cluster_id      = "16f91637-6162-49c8-950f-b53bb533040f"
      - compute         = {
          - cpu = 2
          - ram = 4
        }
      - current_state   = "healthy"
      - description     = "My first test app service."
      - etag            = "Version: 7"
      - id              = "51b402e6-aa55-40bf-82ec-aae5ccf257cf"
      - if_match        = null
      - name            = "new-terraform-app-service"
      - nodes           = 2
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      - project_id      = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
      - version         = "3.0.8-1.0.0"
    } -> null
  - bucket              = "new_terraform_bucket" -> null
  - certificate         = {
      - cluster_id      = "16f91637-6162-49c8-950f-b53bb533040f"
      - data            = [
          - {
              - certificate = <<-EOT
                    -----BEGIN CERTIFICATE-----
                    MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
                    JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEyMDYy
                    MjEyNTlaFw0yOTEyMDYyMzEyNTlaMCQxEjAQBgNVBAoMCUNvdWNoYmFzZTEOMAwG
                    A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfvOIi
                    enG4Dp+hJu9asdxEMRmH70hDyMXv5ZjBhbo39a42QwR59y/rC/sahLLQuNwqif85
                    Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
                    WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
                    rL3ooAw/i0uxr/TQSaxi1l8t4iMt4/gU+W52+8Yol0JbXBTFX6itg62ppb/Eugmn
                    mQRMgL67ccZs7cJ9/A0wlXencX2ohZQOR3mtknfol3FH4+glQFn27Q4xBCzVkY9j
                    KQ20T1LgmGSngBInAgMBAAGjQjBAMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYE
                    FJQOBPvrkU2In1Sjoxt97Xy8+cKNMA4GA1UdDwEB/wQEAwIBhjANBgkqhkiG9w0B
                    AQsFAAOCAQEARgM6XwcXPLSpFdSf0w8PtpNGehmdWijPM3wHb7WZiS47iNen3oq8
                    m2mm6V3Z57wbboPpfI+VEzbhiDcFfVnK1CXMC0tkF3fnOG1BDDvwt4jU95vBiNjY
                    xdzlTP/Z+qr0cnVbGBSZ+fbXstSiRaaAVcqQyv3BRvBadKBkCyPwo+7svQnScQ5P
                    Js7HEHKVms5tZTgKIw1fbmgR2XHleah1AcANB+MAPBCcTgqurqr5G7W2aPSBLLGA
                    fRIiVzm7VFLc7kWbp7ENH39HVG6TZzKnfl9zJYeiklo5vQQhGSMhzBsO70z4RRzi
                    DPFAN/4qZAgD5q3AFNIq2WWADFQGSwVJhg==
                    -----END CERTIFICATE-----
                EOT
            },
        ]
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      - project_id      = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
    } -> null
  - cluster             = {
      - app_service_id   = null
      - audit            = {
          - created_at  = "2023-10-24 21:46:20.62652889 +0000 UTC"
          - created_by  = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
          - modified_at = "2023-10-24 21:49:08.817024833 +0000 UTC"
          - modified_by = "M02xhF8sv9CtkCEFLM9m1lDM10b5kYVV"
          - version     = 5
        }
      - availability     = {
          - type = "multi"
        }
      - cloud_provider   = {
          - cidr   = "10.10.30.0/23"
          - region = "us-east-1"
          - type   = "aws"
        }
      - couchbase_server = {
          - version = "7.1"
        }
      - current_state    = "healthy"
      - description      = "My first test cluster for multiple services."
      - etag             = "Version: 5"
      - id               = "16f91637-6162-49c8-950f-b53bb533040f"
      - if_match         = null
      - name             = "My First Terraform Cluster"
      - organization_id  = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      - project_id       = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
      - service_groups   = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 4
                      - ram = 16
                    }
                  - disk    = {
                      - iops    = 5000
                      - storage = 50
                      - type    = "io2"
                    }
                }
              - num_of_nodes = 3
              - services     = [
                  - "data",
                  - "index",
                  - "query",
                ]
            },
        ]
      - support          = {
          - plan     = "developer pro"
          - timezone = "PT"
        }
    } -> null
  - database_credential = (sensitive value) -> null
  - organization        = {
      - audit           = {
          - created_at  = "2021-12-03 16:14:45.105347711 +0000 UTC"
          - created_by  = "b1cf2366-0401-4cac-8770-f24e511f6c0a"
          - modified_at = "2023-07-07 10:17:11.682962228 +0000 UTC"
          - modified_by = "8578bc6d-e67d-44a3-80f0-d72a5cfbebc6"
          - version     = 28
        }
      - description     = ""
      - name            = "capella-prod"
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      - preferences     = {
          - session_duration = 7200
        }
    } -> null
  - project             = "My First Terraform Project" -> null
  - user                = {
      - audit                = {
          - created_at  = "2023-10-24 20:55:30.005741097 +0000 UTC"
          - created_by  = "0f451c57-2a96-4a59-b8b1-d2b144411598"
          - modified_at = "2023-10-24 20:55:30.005741097 +0000 UTC"
          - modified_by = "0f451c57-2a96-4a59-b8b1-d2b144411598"
          - version     = 1
        }
      - email                = "johndoe@couchbase.com"
      - enable_notifications = false
      - expires_at           = "2024-01-22T20:55:30.005741367Z"
      - id                   = "0f451c57-2a96-4a59-b8b1-d2b144411598"
      - inactive             = true
      - last_login           = ""
      - name                 = "John Doe"
      - organization_id      = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      - organization_roles   = [
          - "organizationMember",
        ]
      - region               = ""
      - resources            = [
          - {
              - id    = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
              - roles = [
                  - "projectViewer",
                ]
              - type  = "project"
            },
          - {
              - id    = "fda8057c-c0ee-4b81-94e4-f9a3d6ca6449"
              - roles = [
                  - "projectDataReaderWriter",
                ]
              - type  = "project"
            },
        ]
      - status               = "not-verified"
      - time_zone            = ""
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_allowlist.new_allowlist: Destroying... [id=b6cfe2a8-fbb8-4bba-ba38-b26a29974767]
capella_database_credential.new_database_credential: Destroying... [id=62ad5f82-8100-4d6e-9216-70e1f82d0dbe]
capella_app_service.new_app_service: Destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf]
capella_apikey.new_apikey: Destroying... [id=iCInIRNFQVGitmFQAjf1JD2FKf8PaLCD]
capella_bucket.new_bucket: Destroying... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
capella_database_credential.new_database_credential: Destruction complete after 0s
capella_apikey.new_apikey: Destruction complete after 0s
capella_bucket.new_bucket: Destruction complete after 2s
capella_allowlist.new_allowlist: Destruction complete after 8s
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 10s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 20s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 30s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 40s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 50s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 1m0s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 1m10s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 1m20s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 1m30s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 1m40s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 1m50s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 2m0s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 2m10s elapsed]
capella_app_service.new_app_service: Still destroying... [id=51b402e6-aa55-40bf-82ec-aae5ccf257cf, 2m20s elapsed]
capella_app_service.new_app_service: Destruction complete after 2m23s
capella_cluster.new_cluster: Destroying... [id=16f91637-6162-49c8-950f-b53bb533040f]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 1m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 1m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 1m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 1m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 1m40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 1m50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 2m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 2m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 2m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 2m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 2m40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 2m50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 3m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=16f91637-6162-49c8-950f-b53bb533040f, 3m10s elapsed]
capella_cluster.new_cluster: Destruction complete after 3m10s
capella_project.new_project: Destroying... [id=fda8057c-c0ee-4b81-94e4-f9a3d6ca6449]
capella_project.new_project: Destruction complete after 2s

Destroy complete! Resources: 7 destroyed. 

```