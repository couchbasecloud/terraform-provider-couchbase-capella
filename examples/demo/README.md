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
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 0s [id=0783f698-ac58-4018-84a3-31c3b6ef785d]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.capella_certificate.existing_certificate will be read during apply
  # (config refers to values not yet known)
 <= data "capella_certificate" "existing_certificate" {
      + cluster_id      = (known after apply)
      + data            = [
        ] -> (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = (known after apply)
    }

  # capella_allowlist.new_allowlist will be created
  + resource "capella_allowlist" "new_allowlist" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cidr            = "8.8.8.8/32"
      + cluster_id      = (known after apply)
      + comment         = "Allow access from a public IP"
      + expires_at      = "2023-11-30T23:59:59.465Z"
      + id              = (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = (known after apply)
    }

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.43.0/23",
        ]
      + audit              = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "My First Terraform API Key"
      + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + audit                      = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = (known after apply)
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id                 = (known after apply)
      + replicas                   = 1
      + stats                      = {
          + disk_used_in_mib   = (known after apply)
          + item_count         = (known after apply)
          + memory_used_in_mib = (known after apply)
          + ops_per_second     = (known after apply)
        }
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # capella_cluster.new_cluster will be created
  + resource "capella_cluster" "new_cluster" {
      + app_service_id   = (known after apply)
      + audit            = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "192.168.0.0/20"
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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cluster_id      = (known after apply)
      + id              = (known after apply)
      + name            = "terraform_db_credential"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + password        = (sensitive value)
      + project_id      = (known after apply)
    }

  # capella_project.new_project will be created
  + resource "capella_project" "new_project" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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

Plan: 7 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikey              = (sensitive value)
  + bucket              = "new_terraform_bucket"
  + certificate         = {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = (known after apply)
    }
  + cluster             = {
      + app_service_id   = (known after apply)
      + audit            = (known after apply)
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "192.168.0.0/20"
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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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
          + created_at  = "2022-05-27 04:19:18.057836345 +0000 UTC"
          + created_by  = "fff64e90-e839-4b96-956e-05135e30d35b"
          + modified_at = "2022-05-27 04:19:18.057836345 +0000 UTC"
          + modified_by = "fff64e90-e839-4b96-956e-05135e30d35b"
          + version     = 1
        }
      + description     = ""
      + id              = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + name            = "Couchbase"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + preferences     = null
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
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

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
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_organization.existing_organization: Reading...
data.capella_organization.existing_organization: Read complete after 1s [id=0783f698-ac58-4018-84a3-31c3b6ef785d]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.capella_certificate.existing_certificate will be read during apply
  # (config refers to values not yet known)
 <= data "capella_certificate" "existing_certificate" {
      + cluster_id      = (known after apply)
      + data            = [
        ] -> (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = (known after apply)
    }

  # capella_allowlist.new_allowlist will be created
  + resource "capella_allowlist" "new_allowlist" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cidr            = "8.8.8.8/32"
      + cluster_id      = (known after apply)
      + comment         = "Allow access from a public IP"
      + expires_at      = "2023-11-30T23:59:59.465Z"
      + id              = (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = (known after apply)
    }

  # capella_apikey.new_apikey will be created
  + resource "capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.43.0/23",
        ]
      + audit              = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description        = (known after apply)
      + expiry             = (known after apply)
      + id                 = (known after apply)
      + name               = "My First Terraform API Key"
      + organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + audit                      = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = (known after apply)
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id                 = (known after apply)
      + replicas                   = 1
      + stats                      = {
          + disk_used_in_mib   = (known after apply)
          + item_count         = (known after apply)
          + memory_used_in_mib = (known after apply)
          + ops_per_second     = (known after apply)
        }
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # capella_cluster.new_cluster will be created
  + resource "capella_cluster" "new_cluster" {
      + app_service_id   = (known after apply)
      + audit            = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "192.168.0.0/20"
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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + cluster_id      = (known after apply)
      + id              = (known after apply)
      + name            = "terraform_db_credential"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + password        = (sensitive value)
      + project_id      = (known after apply)
    }

  # capella_project.new_project will be created
  + resource "capella_project" "new_project" {
      + audit           = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
    }

  # capella_user.new_user will be created
  + resource "capella_user" "new_user" {
      + audit                = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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

Plan: 7 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikey              = (sensitive value)
  + bucket              = "new_terraform_bucket"
  + certificate         = {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = (known after apply)
    }
  + cluster             = {
      + app_service_id   = (known after apply)
      + audit            = (known after apply)
      + availability     = {
          + type = "multi"
        }
      + cloud_provider   = {
          + cidr   = "192.168.0.0/20"
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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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
          + created_at  = "2022-05-27 04:19:18.057836345 +0000 UTC"
          + created_by  = "fff64e90-e839-4b96-956e-05135e30d35b"
          + modified_at = "2022-05-27 04:19:18.057836345 +0000 UTC"
          + modified_by = "fff64e90-e839-4b96-956e-05135e30d35b"
          + version     = 1
        }
      + description     = ""
      + id              = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + name            = "Couchbase"
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + preferences     = null
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
      + organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
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
capella_project.new_project: Creation complete after 1s [id=2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6]
capella_apikey.new_apikey: Creating...
capella_user.new_user: Creating...
capella_cluster.new_cluster: Creating...
capella_apikey.new_apikey: Creation complete after 0s [id=J0PoQ1N1oWgxT2MOLsq5oJ8l248AJSTp]
capella_user.new_user: Creation complete after 0s [id=968385b3-758e-43d9-9f1f-c69d6ca4dd8b]
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
capella_cluster.new_cluster: Creation complete after 2m42s [id=cddd2d52-60f7-46b2-8757-765ef13729b3]
data.capella_certificate.existing_certificate: Reading...
capella_allowlist.new_allowlist: Creating...
capella_bucket.new_bucket: Creating...
capella_database_credential.new_database_credential: Creating...
data.capella_certificate.existing_certificate: Read complete after 0s
capella_bucket.new_bucket: Creation complete after 0s [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
capella_allowlist.new_allowlist: Creation complete after 0s [id=cb3fb642-b575-43ae-bd1d-f086e3d6060e]
capella_database_credential.new_database_credential: Creation complete after 1s [id=12d24fcd-43eb-4d91-8303-25adaf06fd4a]

Apply complete! Resources: 7 added, 0 changed, 0 destroyed.

Outputs:

apikey = <sensitive>
bucket = "new_terraform_bucket"
certificate = {
  "cluster_id" = "cddd2d52-60f7-46b2-8757-765ef13729b3"
  "data" = tolist([
    {
      "certificate" = <<-EOT
      -----BEGIN CERTIFICATE-----
      MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
      JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEyMDYy
      ...redacted...
      Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
      WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
      AQsFAAOCAQEARgM6XwcXPLSpFdSf0w8PtpNGehmdWijPM3wHb7WZiS47iNen3oq8
      ...redacted
      fRIiVzm7VFLc7kWbp7ENH39HVG6TZzKnfl9zJYeiklo5vQQhGSMhzBsO70z4RRzi
      DPFAN/4qZAgD5q3AFNIq2WWADFQGSwVJhg==
      -----END CERTIFICATE-----
      EOT
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
}
cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-10-04 23:40:01.464703649 +0000 UTC"
    "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "modified_at" = "2023-10-04 23:42:42.113481767 +0000 UTC"
    "modified_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "version" = 5
  }
  "availability" = {
    "type" = "multi"
  }
  "cloud_provider" = {
    "cidr" = "192.168.0.0/20"
    "region" = "us-east-1"
    "type" = "aws"
  }
  "couchbase_server" = {
    "version" = "7.1"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "cddd2d52-60f7-46b2-8757-765ef13729b3"
  "if_match" = tostring(null)
  "name" = "My First Terraform Cluster"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
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
    "created_at" = "2022-05-27 04:19:18.057836345 +0000 UTC"
    "created_by" = "fff64e90-e839-4b96-956e-05135e30d35b"
    "modified_at" = "2022-05-27 04:19:18.057836345 +0000 UTC"
    "modified_by" = "fff64e90-e839-4b96-956e-05135e30d35b"
    "version" = 1
  }
  "description" = ""
  "id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "name" = "Couchbase"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "preferences" = null /* object */
}
project = "My First Terraform Project"
user = {
  "audit" = {
    "created_at" = "2023-10-04 23:40:01.108967467 +0000 UTC"
    "created_by" = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
    "modified_at" = "2023-10-04 23:40:01.108967467 +0000 UTC"
    "modified_by" = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-02T23:40:01.108968559Z"
  "id" = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
  "inactive" = true
  "last_login" = ""
  "name" = "John Doe"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
      "roles" = tolist([
        "projectViewer",
      ])
      "type" = "project"
    },
    {
      "id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
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
bucket = "new_terraform_bucket"
certificate = {
  "cluster_id" = "cddd2d52-60f7-46b2-8757-765ef13729b3"
  "data" = tolist([
    {
      "certificate" = <<-EOT
      -----BEGIN CERTIFICATE-----
      MIIDFTCCAf2gAwIBAgIRANLVkgOvtaXiQJi0V6qeNtswDQYJKoZIhvcNAQELBQAw
      JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEyMDYy
      ...redacted...
      Fod1DkqgO6Ng3vecSAwyYVkj5NKdycQu5tzsZkghlpSDAyI0xlIPSQjoORA/pCOU
      WOpymA9dOjC1bo6rDyw0yWP2nFAI/KA4Z806XeqLREuB7292UnSsgFs4/5lqeil6
      AQsFAAOCAQEARgM6XwcXPLSpFdSf0w8PtpNGehmdWijPM3wHb7WZiS47iNen3oq8
      ...redacted
      fRIiVzm7VFLc7kWbp7ENH39HVG6TZzKnfl9zJYeiklo5vQQhGSMhzBsO70z4RRzi
      DPFAN/4qZAgD5q3AFNIq2WWADFQGSwVJhg==
      -----END CERTIFICATE-----
      EOT
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
}
cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-10-04 23:40:01.464703649 +0000 UTC"
    "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "modified_at" = "2023-10-04 23:42:42.113481767 +0000 UTC"
    "modified_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "version" = 5
  }
  "availability" = {
    "type" = "multi"
  }
  "cloud_provider" = {
    "cidr" = "192.168.0.0/20"
    "region" = "us-east-1"
    "type" = "aws"
  }
  "couchbase_server" = {
    "version" = "7.1"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "cddd2d52-60f7-46b2-8757-765ef13729b3"
  "if_match" = tostring(null)
  "name" = "My First Terraform Cluster"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
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
    "created_at" = "2022-05-27 04:19:18.057836345 +0000 UTC"
    "created_by" = "fff64e90-e839-4b96-956e-05135e30d35b"
    "modified_at" = "2022-05-27 04:19:18.057836345 +0000 UTC"
    "modified_by" = "fff64e90-e839-4b96-956e-05135e30d35b"
    "version" = 1
  }
  "description" = ""
  "id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "name" = "Couchbase"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "preferences" = null /* object */
}
project = "My First Terraform Project"
user = {
  "audit" = {
    "created_at" = "2023-10-04 23:40:01.108967467 +0000 UTC"
    "created_by" = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
    "modified_at" = "2023-10-04 23:40:01.108967467 +0000 UTC"
    "modified_by" = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-02T23:40:01.108968559Z"
  "id" = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
  "inactive" = true
  "last_login" = ""
  "name" = "John Doe"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = tolist([
    {
      "id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
      "roles" = tolist([
        "projectViewer",
      ])
      "type" = "project"
    },
    {
      "id" = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
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
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_organization.existing_organization: Reading...
capella_project.new_project: Refreshing state... [id=2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6]
data.capella_organization.existing_organization: Read complete after 1s [id=0783f698-ac58-4018-84a3-31c3b6ef785d]
capella_user.new_user: Refreshing state... [id=968385b3-758e-43d9-9f1f-c69d6ca4dd8b]
capella_apikey.new_apikey: Refreshing state... [id=J0PoQ1N1oWgxT2MOLsq5oJ8l248AJSTp]
capella_cluster.new_cluster: Refreshing state... [id=cddd2d52-60f7-46b2-8757-765ef13729b3]
capella_allowlist.new_allowlist: Refreshing state... [id=cb3fb642-b575-43ae-bd1d-f086e3d6060e]
data.capella_certificate.existing_certificate: Reading...
capella_database_credential.new_database_credential: Refreshing state... [id=12d24fcd-43eb-4d91-8303-25adaf06fd4a]
capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
data.capella_certificate.existing_certificate: Read complete after 6s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_allowlist.new_allowlist will be destroyed
  - resource "capella_allowlist" "new_allowlist" {
      - audit           = {
          - created_at  = "2023-10-04 23:42:43.062217825 +0000 UTC" -> null
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - modified_at = "2023-10-04 23:42:43.062217825 +0000 UTC" -> null
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - version     = 1 -> null
        }
      - cidr            = "8.8.8.8/32" -> null
      - cluster_id      = "cddd2d52-60f7-46b2-8757-765ef13729b3" -> null
      - comment         = "Allow access from a public IP" -> null
      - expires_at      = "2023-11-30T23:59:59.465Z" -> null
      - id              = "cb3fb642-b575-43ae-bd1d-f086e3d6060e" -> null
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - project_id      = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
    }

  # capella_apikey.new_apikey will be destroyed
  - resource "capella_apikey" "new_apikey" {
      - allowed_cidrs      = [
          - "10.1.42.0/23",
          - "10.1.43.0/23",
        ] -> null
      - audit              = {
          - created_at  = "2023-10-04 23:40:00.812689485 +0000 UTC" -> null
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - modified_at = "2023-10-04 23:40:00.812689485 +0000 UTC" -> null
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - version     = 1 -> null
        }
      - expiry             = 180 -> null
      - id                 = "J0PoQ1N1oWgxT2MOLsq5oJ8l248AJSTp" -> null
      - name               = "My First Terraform API Key" -> null
      - organization_id    = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - organization_roles = [
          - "organizationMember",
        ] -> null
      - resources          = [
          - {
              - id    = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
              - roles = [
                  - "projectManager",
                  - "projectDataReader",
                ] -> null
              - type  = "project" -> null
            },
        ]
      - token              = (sensitive value)
    }

  # capella_bucket.new_bucket will be destroyed
  - resource "capella_bucket" "new_bucket" {
      - audit                      = {
          - created_at  = "0001-01-01 00:00:00 +0000 UTC" -> null
          - modified_at = "0001-01-01 00:00:00 +0000 UTC" -> null
          - version     = 0 -> null
        }
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "cddd2d52-60f7-46b2-8757-765ef13729b3" -> null
      - durability_level           = "none" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = false -> null
      - id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> null
      - memory_allocation_in_mb    = 100 -> null
      - name                       = "new_terraform_bucket" -> null
      - organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - project_id                 = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
      - replicas                   = 1 -> null
      - stats                      = {
          - disk_used_in_mib   = 18 -> null
          - item_count         = 0 -> null
          - memory_used_in_mib = 54 -> null
          - ops_per_second     = 0 -> null
        }
      - storage_backend            = "couchstore" -> null
      - time_to_live_in_seconds    = 0 -> null
      - type                       = "couchbase" -> null
    }

  # capella_cluster.new_cluster will be destroyed
  - resource "capella_cluster" "new_cluster" {
      - audit            = {
          - created_at  = "2023-10-04 23:40:01.464703649 +0000 UTC" -> null
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - modified_at = "2023-10-04 23:42:42.113481767 +0000 UTC" -> null
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - version     = 5 -> null
        }
      - availability     = {
          - type = "multi" -> null
        }
      - cloud_provider   = {
          - cidr   = "192.168.0.0/20" -> null
          - region = "us-east-1" -> null
          - type   = "aws" -> null
        }
      - couchbase_server = {
          - version = "7.1" -> null
        }
      - current_state    = "healthy" -> null
      - description      = "My first test cluster for multiple services." -> null
      - etag             = "Version: 5" -> null
      - id               = "cddd2d52-60f7-46b2-8757-765ef13729b3" -> null
      - name             = "My First Terraform Cluster" -> null
      - organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - project_id       = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
      - service_groups   = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 4 -> null
                      - ram = 16 -> null
                    }
                  - disk    = {
                      - iops    = 5000 -> null
                      - storage = 50 -> null
                      - type    = "io2" -> null
                    }
                }
              - num_of_nodes = 3 -> null
              - services     = [
                  - "data",
                  - "index",
                  - "query",
                ] -> null
            },
        ]
      - support          = {
          - plan     = "developer pro" -> null
          - timezone = "PT" -> null
        }
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
        ]
      - audit           = {
          - created_at  = "2023-10-04 23:42:42.749664208 +0000 UTC" -> null
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - modified_at = "2023-10-04 23:42:42.749664208 +0000 UTC" -> null
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - version     = 1 -> null
        }
      - cluster_id      = "cddd2d52-60f7-46b2-8757-765ef13729b3" -> null
      - id              = "12d24fcd-43eb-4d91-8303-25adaf06fd4a" -> null
      - name            = "terraform_db_credential" -> null
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - password        = (sensitive value)
      - project_id      = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
    }

  # capella_project.new_project will be destroyed
  - resource "capella_project" "new_project" {
      - audit           = {
          - created_at  = "2023-10-04 23:40:00.476880098 +0000 UTC" -> null
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - modified_at = "2023-10-04 23:40:00.47689516 +0000 UTC" -> null
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - version     = 1 -> null
        }
      - description     = "A Capella Project that will host many Capella clusters." -> null
      - etag            = "Version: 1" -> null
      - id              = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
      - name            = "My First Terraform Project" -> null
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
    }

  # capella_user.new_user will be destroyed
  - resource "capella_user" "new_user" {
      - audit                = {
          - created_at  = "2023-10-04 23:40:01.108967467 +0000 UTC" -> null
          - created_by  = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b" -> null
          - modified_at = "2023-10-04 23:40:01.108967467 +0000 UTC" -> null
          - modified_by = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b" -> null
          - version     = 1 -> null
        }
      - email                = "johndoe@couchbase.com" -> null
      - enable_notifications = false -> null
      - expires_at           = "2024-01-02T23:40:01.108968559Z" -> null
      - id                   = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b" -> null
      - inactive             = true -> null
      - name                 = "John Doe" -> null
      - organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - organization_roles   = [
          - "organizationMember",
        ] -> null
      - resources            = [
          - {
              - id    = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
              - roles = [
                  - "projectViewer",
                ] -> null
              - type  = "project" -> null
            },
          - {
              - id    = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6" -> null
              - roles = [
                  - "projectDataReaderWriter",
                ] -> null
              - type  = "project" -> null
            },
        ]
      - status               = "not-verified" -> null
    }

Plan: 0 to add, 0 to change, 7 to destroy.

Changes to Outputs:
  - apikey              = (sensitive value)
  - bucket              = "new_terraform_bucket" -> null
  - certificate         = {
      - cluster_id      = "cddd2d52-60f7-46b2-8757-765ef13729b3"
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
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id      = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
    } -> null
  - cluster             = {
      - app_service_id   = null
      - audit            = {
          - created_at  = "2023-10-04 23:40:01.464703649 +0000 UTC"
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
          - modified_at = "2023-10-04 23:42:42.113481767 +0000 UTC"
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
          - version     = 5
        }
      - availability     = {
          - type = "multi"
        }
      - cloud_provider   = {
          - cidr   = "192.168.0.0/20"
          - region = "us-east-1"
          - type   = "aws"
        }
      - couchbase_server = {
          - version = "7.1"
        }
      - current_state    = "healthy"
      - description      = "My first test cluster for multiple services."
      - etag             = "Version: 5"
      - id               = "cddd2d52-60f7-46b2-8757-765ef13729b3"
      - if_match         = null
      - name             = "My First Terraform Cluster"
      - organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id       = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
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
  - database_credential = (sensitive value)
  - organization        = {
      - audit           = {
          - created_at  = "2022-05-27 04:19:18.057836345 +0000 UTC"
          - created_by  = "fff64e90-e839-4b96-956e-05135e30d35b"
          - modified_at = "2022-05-27 04:19:18.057836345 +0000 UTC"
          - modified_by = "fff64e90-e839-4b96-956e-05135e30d35b"
          - version     = 1
        }
      - description     = ""
      - id              = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - name            = "Couchbase"
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - preferences     = null
    } -> null
  - project             = "My First Terraform Project" -> null
  - user                = {
      - audit                = {
          - created_at  = "2023-10-04 23:40:01.108967467 +0000 UTC"
          - created_by  = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
          - modified_at = "2023-10-04 23:40:01.108967467 +0000 UTC"
          - modified_by = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
          - version     = 1
        }
      - email                = "johndoe@couchbase.com"
      - enable_notifications = false
      - expires_at           = "2024-01-02T23:40:01.108968559Z"
      - id                   = "968385b3-758e-43d9-9f1f-c69d6ca4dd8b"
      - inactive             = true
      - last_login           = ""
      - name                 = "John Doe"
      - organization_id      = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - organization_roles   = [
          - "organizationMember",
        ]
      - region               = ""
      - resources            = [
          - {
              - id    = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
              - roles = [
                  - "projectViewer",
                ]
              - type  = "project"
            },
          - {
              - id    = "2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6"
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

capella_allowlist.new_allowlist: Destroying... [id=cb3fb642-b575-43ae-bd1d-f086e3d6060e]
capella_bucket.new_bucket: Destroying... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
capella_user.new_user: Destroying... [id=968385b3-758e-43d9-9f1f-c69d6ca4dd8b]
capella_database_credential.new_database_credential: Destroying... [id=12d24fcd-43eb-4d91-8303-25adaf06fd4a]
capella_apikey.new_apikey: Destroying... [id=J0PoQ1N1oWgxT2MOLsq5oJ8l248AJSTp]
capella_database_credential.new_database_credential: Destruction complete after 0s
capella_apikey.new_apikey: Destruction complete after 0s
capella_user.new_user: Destruction complete after 1s
capella_bucket.new_bucket: Destruction complete after 1s
capella_allowlist.new_allowlist: Destruction complete after 1s
capella_cluster.new_cluster: Destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 1m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 1m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 1m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 1m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 1m40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 1m50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 2m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 2m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 2m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 2m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=cddd2d52-60f7-46b2-8757-765ef13729b3, 2m40s elapsed]
capella_cluster.new_cluster: Destruction complete after 2m42s
capella_project.new_project: Destroying... [id=2cc5c949-709b-4fac-ad94-2c9e1a5cd8d6]
capella_project.new_project: Destruction complete after 0s

Destroy complete! Resources: 7 destroyed.
```