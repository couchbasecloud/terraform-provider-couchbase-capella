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
10. Create a new sample bucket in the cluster.
11. Create a new app service in the cluster.
12. Create a new scope in the bucket of the cluster.
13. Create a new collection in the scope of a bucket.
14. Create a new on/off schedule for the cluster.
15. Create a new audit log settings.
16. Enable private endpoint service on the cluster.
17. Create a new network peer.
18. Create an non-deferred secondary index.

## Pre-Requisites:

Make sure you have followed the pre-requisite steps from the parent readme that builds the binary and stores the path in a .terraformrc file.

### Variables file

- Copy the `terraform.template.tfvars` file to `terraform.tfvars` file in the same directory
- Create 1 V4 API key in your organization using the Capella UI.
- Replace all placeholders with actual values. 
- Set the auth token as an environment variable:

```
export TF_VAR_auth_token='<v4-api-key-secret>'
```

## Execution

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_events.existing_events: Reading...
data.couchbase-capella_organization.existing_organization: Reading...
data.couchbase-capella_organization.existing_organization: Read complete after 1s [name=my-org]apella-prod]
data.couchbase-capella_events.existing_events: Read complete after 1s
data.couchbase-capella_event.existing_event: Reading...
data.couchbase-capella_event.existing_event: Read complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.couchbase-capella_audit_log_event_ids.event_list will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_audit_log_event_ids" "event_list" {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # data.couchbase-capella_certificate.existing_certificate will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_certificate" "existing_certificate" {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # data.couchbase-capella_project_event.existing_project_event will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_project_event" "existing_project_event" {
      + alert_key        = (known after apply)
      + app_service_id   = (known after apply)
      + app_service_name = (known after apply)
      + cluster_id       = (known after apply)
      + cluster_name     = (known after apply)
      + id               = (known after apply)
      + image_url        = (known after apply)
      + incident_ids     = (known after apply)
      + key              = (known after apply)
      + kv               = (known after apply)
      + occurrence_count = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = (known after apply)
      + project_name     = (known after apply)
      + request_id       = (known after apply)
      + session_id       = (known after apply)
      + severity         = (known after apply)
      + source           = (known after apply)
      + summary          = (known after apply)
      + timestamp        = (known after apply)
      + user_email       = (known after apply)
      + user_id          = (known after apply)
      + user_name        = (known after apply)
    }

  # data.couchbase-capella_project_events.existing_project_events will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_project_events" "existing_project_events" {
      + cursor          = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # couchbase-capella_allowlist.new_allowlist will be created
  + resource "couchbase-capella_allowlist" "new_allowlist" {
      + audit           = (known after apply)
      + cidr            = "8.8.8.8/32"
      + cluster_id      = (known after apply)
      + comment         = "Allow access from a public IP"
      + expires_at      = "2043-11-30T23:59:59.465Z"
      + id              = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # couchbase-capella_apikey.new_apikey will be created
  + resource "couchbase-capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.43.0/23",
        ]
      + audit              = (known after apply)
      + expiry             = 180
      + id                 = (known after apply)
      + name               = "My First Terraform API Key"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReader",
                  + "projectManager",
                ]
              + type  = "project"
            },
        ]
      + rotate             = (known after apply)
      + secret             = (sensitive value)
      + token              = (sensitive value)
    }

  # couchbase-capella_app_service.new_app_service will be created
  + resource "couchbase-capella_app_service" "new_app_service" {
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
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be created
  + resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      + audit_enabled     = true
      + cluster_id        = (known after apply)
      + disabled_users    = []
      + enabled_event_ids = (known after apply)
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = (known after apply)
    }

  # couchbase-capella_bucket.new_bucket will be created
  + resource "couchbase-capella_bucket" "new_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = (known after apply)
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = (known after apply)
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # couchbase-capella_cluster.new_cluster will be created
  + resource "couchbase-capella_cluster" "new_cluster" {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = {
          + type = "multi"
        }
      + cloud_provider                = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type            = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "My first test cluster for multiple services."
      + enable_private_dns_resolution = false
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + name                          = "My First Terraform Cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = (known after apply)
      + service_groups                = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = (known after apply)
                      + iops          = 5000
                      + storage       = 50
                      + type          = "io2"
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
      + support                       = {
          + plan     = "enterprise"
          + timezone = "PT"
        }
    }

  # couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule will be created
  + resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
      + cluster_id      = (known after apply)
      + days            = [
          + {
              + day   = "monday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 30
                }
            },
          + {
              + day   = "tuesday"
              + from  = {
                  + hour   = 12
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 19
                  + minute = 30
                }
            },
          + {
              + day   = "wednesday"
              + state = "on"
            },
          + {
              + day   = "thursday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
            },
          + {
              + day   = "friday"
              + from  = {
                  + hour   = 0
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 12
                  + minute = 30
                }
            },
          + {
              + day   = "saturday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 0
                }
            },
          + {
              + day   = "sunday"
              + state = "off"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + timezone        = "US/Pacific"
    }

  # couchbase-capella_collection.new_collection will be created
  + resource "couchbase-capella_collection" "new_collection" {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + max_ttl         = 200
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }

  # couchbase-capella_database_credential.new_database_credential will be created
  + resource "couchbase-capella_database_credential" "new_database_credential" {
      + access          = [
          + {
              + privileges = [
                  + "data_reader",
                ]
            },
          + {
              + privileges = [
                  + "data_writer",
                ]
              + resources  = {
                  + buckets = [
                      + {
                          + name   = "new_terraform_bucket"
                          + scopes = [
                              + {
                                  + name = "_default"
                                },
                            ]
                        },
                    ]
                }
            },
        ]
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + id              = (known after apply)
      + name            = "terraform_db_credential"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + password        = (sensitive value)
      + project_id      = (known after apply)
    }

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + provider_config = {
          + aws_config = {
              + account_id  = "123456789123"
              + cidr        = "10.0.0.0/16"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00"
            }
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }

  # couchbase-capella_private_endpoint_service.new_service will be created
  + resource "couchbase-capella_private_endpoint_service" "new_service" {
      + cluster_id      = (known after apply)
      + enabled         = true
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # couchbase-capella_project.new_project will be created
  + resource "couchbase-capella_project" "new_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
    }

  # couchbase-capella_query_indexes.idx will be created
  + resource "couchbase-capella_query_indexes" "idx" {
      + bucket_name     = "new_terraform_bucket"
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + index_keys      = [
          + "id",
          + "age",
          + "name",
        ]
      + index_name      = "idx1"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
      + status          = (known after apply)
      + where           = "dept = 'sales'"
      + with            = {
          + defer_build = false
          + num_replica = 1
        }
    }

  # couchbase-capella_sample_bucket.new_sample_bucket will be created
  + resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
      + bucket_conflict_resolution = (known after apply)
      + cluster_id                 = (known after apply)
      + durability_level           = (known after apply)
      + eviction_policy            = (known after apply)
      + flush                      = (known after apply)
      + id                         = (known after apply)
      + memory_allocation_in_mb    = (known after apply)
      + name                       = "gamesim-sample"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = (known after apply)
      + replicas                   = (known after apply)
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = (known after apply)
      + type                       = (known after apply)
    }

  # couchbase-capella_scope.new_scope will be created
  + resource "couchbase-capella_scope" "new_scope" {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collections     = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }

  # couchbase-capella_user.new_user will be created
  + resource "couchbase-capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                  + "projectViewer",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Plan: 16 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikey                  = (sensitive value)
  + app_service             = {
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
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }
  + bucket                  = "new_terraform_bucket"
  + certificate             = {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }
  + cluster                 = {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = {
          + type = "multi"
        }
      + cloud_provider                = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type            = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "My first test cluster for multiple services."
      + enable_private_dns_resolution = false
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + if_match                      = null
      + name                          = "My First Terraform Cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = (known after apply)
      + service_groups                = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = (known after apply)
                      + iops          = 5000
                      + storage       = 50
                      + type          = "io2"
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
      + support                       = {
          + plan     = "enterprise"
          + timezone = "PT"
        }
      + zones                         = null
    }
  + cluster_onoff_schedule  = {
      + cluster_id      = (known after apply)
      + days            = [
          + {
              + day   = "monday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 30
                }
            },
          + {
              + day   = "tuesday"
              + from  = {
                  + hour   = 12
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 19
                  + minute = 30
                }
            },
          + {
              + day   = "wednesday"
              + from  = null
              + state = "on"
              + to    = null
            },
          + {
              + day   = "thursday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = null
            },
          + {
              + day   = "friday"
              + from  = {
                  + hour   = 0
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 12
                  + minute = 30
                }
            },
          + {
              + day   = "saturday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 0
                }
            },
          + {
              + day   = "sunday"
              + from  = null
              + state = "off"
              + to    = null
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + timezone        = "US/Pacific"
    }
  + collection              = {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + max_ttl         = 200
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }
  + database_credential     = (sensitive value)
  + existing_event          = {
      + alert_key        = "logged_in"
      + app_service_id   = null
      + app_service_name = null
      + cluster_id       = null
      + cluster_name     = null
      + id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + image_url        = null
      + incident_ids     = []
      + key              = "logged_in"
      + kv               = "null"
      + occurrence_count = null
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = null
      + project_name     = null
      + request_id       = null
      + session_id       = null
      + severity         = "info"
      + source           = "cp-ns"
      + summary          = null
      + timestamp        = "2024-12-11 20:33:36.600542622 +0000 UTC"
      + user_email       = null
      + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + user_name        = "SDK QE"
    }
  + existing_events         = {
      + cluster_ids     = null
      + cursor          = {
          + hrefs = {
              + first    = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=10"
              + last     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=37&perPage=10"
              + next     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=10"
              + previous = ""
            }
          + pages = {
              + last        = 37
              + next        = 2
              + page        = 1
              + per_page    = 10
              + previous    = 0
              + total_items = 369
            }
        }
      + data            = [
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:33:36.600542622 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:38:27.417814023 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "7e74fc9b-3d7c-746d-0c36-751a984a5747"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:38:33.290122146 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:38:41.431711495 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "43ce65aa-9d2c-6320-5045-357f8e5a7ec8"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:39:37.974685759 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "1ce5b903-9f45-eea3-a601-2597282a2a09"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:39:42.809020171 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "126012e4-95cb-0002-13be-34b673035f31"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:39:49.230386108 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "9fb9f8bd-787d-6604-5887-b67d34e8be64"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:41:53.158231931 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:41:59.204057247 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:42:05.079807567 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
        ]
      + from            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = null
      + per_page        = null
      + project_ids     = null
      + severity_levels = null
      + sort_by         = null
      + sort_direction  = null
      + tags            = null
      + to              = null
      + user_ids        = null
    }
  + existing_project_event  = {
      + alert_key        = (known after apply)
      + app_service_id   = (known after apply)
      + app_service_name = (known after apply)
      + cluster_id       = (known after apply)
      + cluster_name     = (known after apply)
      + id               = (known after apply)
      + image_url        = (known after apply)
      + incident_ids     = (known after apply)
      + key              = (known after apply)
      + kv               = (known after apply)
      + occurrence_count = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = (known after apply)
      + project_name     = (known after apply)
      + request_id       = (known after apply)
      + session_id       = (known after apply)
      + severity         = (known after apply)
      + source           = (known after apply)
      + summary          = (known after apply)
      + timestamp        = (known after apply)
      + user_email       = (known after apply)
      + user_id          = (known after apply)
      + user_name        = (known after apply)
    }
  + existing_project_events = {
      + cluster_ids     = null
      + cursor          = (known after apply)
      + data            = (known after apply)
      + from            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = null
      + per_page        = null
      + project_id      = (known after apply)
      + severity_levels = null
      + sort_by         = null
      + sort_direction  = null
      + tags            = null
      + to              = null
      + user_ids        = null
    }
  + idx                     = {
      + bucket_name     = "new_terraform_bucket"
      + build_indexes   = null
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + index_keys      = [
          + "id",
          + "age",
          + "name",
        ]
      + index_name      = "idx1"
      + is_primary      = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + partition_by    = null
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
      + status          = (known after apply)
      + where           = "dept = 'sales'"
      + with            = {
          + defer_build   = false
          + num_partition = null
          + num_replica   = 1
        }
    }
  + network_peer            = {
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + provider_config = {
          + aws_config   = {
              + account_id  = "123456789123"
              + cidr        = "10.0.0.0/16"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00"
            }
          + azure_config = null
          + gcp_config   = null
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }
  + new_auditlogsettings    = {
      + audit_enabled     = true
      + cluster_id        = (known after apply)
      + disabled_users    = []
      + enabled_event_ids = (known after apply)
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = (known after apply)
    }
  + organization            = {
      + audit           = {
          + created_at  = "2021-12-03 16:14:45.105347711 +0000 UTC"
          + created_by  = "b1cf2366-0401-4cac-8770-f24e511f6c0a"
          + modified_at = "2024-11-18 15:22:27.933390014 +0000 UTC"
          + modified_by = "ab88be62-e7f5-4701-828d-129e8641d111"
          + version     = 39
        }
      + description     = ""
      + name            = "capella-prod"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + preferences     = {
          + session_duration = 7200
        }
    }
  + project                 = "My First Terraform Project"
  + sample_bucket           = "gamesim-sample"
  + scope                   = {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collections     = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }
  + user                    = {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                  + "projectViewer",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_organization.existing_organization: Reading...
data.couchbase-capella_events.existing_events: Reading...
data.couchbase-capella_organization.existing_organization: Read complete after 0s [name=my-org]apella-prod]
data.couchbase-capella_events.existing_events: Read complete after 0s
data.couchbase-capella_event.existing_event: Reading...
data.couchbase-capella_event.existing_event: Read complete after 1s [id=ffffffff-aaaa-1414-eeee-000000000000]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.couchbase-capella_audit_log_event_ids.event_list will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_audit_log_event_ids" "event_list" {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # data.couchbase-capella_certificate.existing_certificate will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_certificate" "existing_certificate" {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # data.couchbase-capella_project_event.existing_project_event will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_project_event" "existing_project_event" {
      + alert_key        = (known after apply)
      + app_service_id   = (known after apply)
      + app_service_name = (known after apply)
      + cluster_id       = (known after apply)
      + cluster_name     = (known after apply)
      + id               = (known after apply)
      + image_url        = (known after apply)
      + incident_ids     = (known after apply)
      + key              = (known after apply)
      + kv               = (known after apply)
      + occurrence_count = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = (known after apply)
      + project_name     = (known after apply)
      + request_id       = (known after apply)
      + session_id       = (known after apply)
      + severity         = (known after apply)
      + source           = (known after apply)
      + summary          = (known after apply)
      + timestamp        = (known after apply)
      + user_email       = (known after apply)
      + user_id          = (known after apply)
      + user_name        = (known after apply)
    }

  # data.couchbase-capella_project_events.existing_project_events will be read during apply
  # (config refers to values not yet known)
 <= data "couchbase-capella_project_events" "existing_project_events" {
      + cursor          = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # couchbase-capella_allowlist.new_allowlist will be created
  + resource "couchbase-capella_allowlist" "new_allowlist" {
      + audit           = (known after apply)
      + cidr            = "8.8.8.8/32"
      + cluster_id      = (known after apply)
      + comment         = "Allow access from a public IP"
      + expires_at      = "2043-11-30T23:59:59.465Z"
      + id              = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # couchbase-capella_apikey.new_apikey will be created
  + resource "couchbase-capella_apikey" "new_apikey" {
      + allowed_cidrs      = [
          + "10.1.42.0/23",
          + "10.1.43.0/23",
        ]
      + audit              = (known after apply)
      + expiry             = 180
      + id                 = (known after apply)
      + name               = "My First Terraform API Key"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_roles = [
          + "organizationMember",
        ]
      + resources          = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReader",
                  + "projectManager",
                ]
              + type  = "project"
            },
        ]
      + rotate             = (known after apply)
      + secret             = (sensitive value)
      + token              = (sensitive value)
    }

  # couchbase-capella_app_service.new_app_service will be created
  + resource "couchbase-capella_app_service" "new_app_service" {
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
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be created
  + resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      + audit_enabled     = true
      + cluster_id        = (known after apply)
      + disabled_users    = []
      + enabled_event_ids = (known after apply)
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = (known after apply)
    }

  # couchbase-capella_bucket.new_bucket will be created
  + resource "couchbase-capella_bucket" "new_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = (known after apply)
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = (known after apply)
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # couchbase-capella_cluster.new_cluster will be created
  + resource "couchbase-capella_cluster" "new_cluster" {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = {
          + type = "multi"
        }
      + cloud_provider                = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type            = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "My first test cluster for multiple services."
      + enable_private_dns_resolution = false
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + name                          = "My First Terraform Cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = (known after apply)
      + service_groups                = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = (known after apply)
                      + iops          = 5000
                      + storage       = 50
                      + type          = "io2"
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
      + support                       = {
          + plan     = "enterprise"
          + timezone = "PT"
        }
    }

  # couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule will be created
  + resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
      + cluster_id      = (known after apply)
      + days            = [
          + {
              + day   = "monday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 30
                }
            },
          + {
              + day   = "tuesday"
              + from  = {
                  + hour   = 12
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 19
                  + minute = 30
                }
            },
          + {
              + day   = "wednesday"
              + state = "on"
            },
          + {
              + day   = "thursday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
            },
          + {
              + day   = "friday"
              + from  = {
                  + hour   = 0
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 12
                  + minute = 30
                }
            },
          + {
              + day   = "saturday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 0
                }
            },
          + {
              + day   = "sunday"
              + state = "off"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + timezone        = "US/Pacific"
    }

  # couchbase-capella_collection.new_collection will be created
  + resource "couchbase-capella_collection" "new_collection" {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + max_ttl         = 200
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }

  # couchbase-capella_database_credential.new_database_credential will be created
  + resource "couchbase-capella_database_credential" "new_database_credential" {
      + access          = [
          + {
              + privileges = [
                  + "data_reader",
                ]
            },
          + {
              + privileges = [
                  + "data_writer",
                ]
              + resources  = {
                  + buckets = [
                      + {
                          + name   = "new_terraform_bucket"
                          + scopes = [
                              + {
                                  + name = "_default"
                                },
                            ]
                        },
                    ]
                }
            },
        ]
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + id              = (known after apply)
      + name            = "terraform_db_credential"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + password        = (sensitive value)
      + project_id      = (known after apply)
    }

  # couchbase-capella_network_peer.new_network_peer will be created
  + resource "couchbase-capella_network_peer" "new_network_peer" {
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + provider_config = {
          + aws_config = {
              + account_id  = "123456789123"
              + cidr        = "10.0.0.0/16"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00"
            }
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }

  # couchbase-capella_private_endpoint_service.new_service will be created
  + resource "couchbase-capella_private_endpoint_service" "new_service" {
      + cluster_id      = (known after apply)
      + enabled         = true
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }

  # couchbase-capella_project.new_project will be created
  + resource "couchbase-capella_project" "new_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
    }

  # couchbase-capella_query_indexes.idx will be created
  + resource "couchbase-capella_query_indexes" "idx" {
      + bucket_name     = "new_terraform_bucket"
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + index_keys      = [
          + "id",
          + "age",
          + "name",
        ]
      + index_name      = "idx1"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
      + status          = (known after apply)
      + where           = "dept = 'sales'"
      + with            = {
          + defer_build = false
          + num_replica = 1
        }
    }

  # couchbase-capella_sample_bucket.new_sample_bucket will be created
  + resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
      + bucket_conflict_resolution = (known after apply)
      + cluster_id                 = (known after apply)
      + durability_level           = (known after apply)
      + eviction_policy            = (known after apply)
      + flush                      = (known after apply)
      + id                         = (known after apply)
      + memory_allocation_in_mb    = (known after apply)
      + name                       = "gamesim-sample"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = (known after apply)
      + replicas                   = (known after apply)
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = (known after apply)
      + type                       = (known after apply)
    }

  # couchbase-capella_scope.new_scope will be created
  + resource "couchbase-capella_scope" "new_scope" {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collections     = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }

  # couchbase-capella_user.new_user will be created
  + resource "couchbase-capella_user" "new_user" {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                  + "projectViewer",
                ]
              + type  = "project"
            },
        ]
      + status               = (known after apply)
      + time_zone            = (known after apply)
    }

Plan: 16 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + apikey                  = (sensitive value)
  + app_service             = {
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
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + version         = (known after apply)
    }
  + bucket                  = "new_terraform_bucket"
  + certificate             = {
      + cluster_id      = (known after apply)
      + data            = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
    }
  + cluster                 = {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = {
          + type = "multi"
        }
      + cloud_provider                = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type            = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "My first test cluster for multiple services."
      + enable_private_dns_resolution = false
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + if_match                      = null
      + name                          = "My First Terraform Cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = (known after apply)
      + service_groups                = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = (known after apply)
                      + iops          = 5000
                      + storage       = 50
                      + type          = "io2"
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
      + support                       = {
          + plan     = "enterprise"
          + timezone = "PT"
        }
      + zones                         = null
    }
  + cluster_onoff_schedule  = {
      + cluster_id      = (known after apply)
      + days            = [
          + {
              + day   = "monday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 30
                }
            },
          + {
              + day   = "tuesday"
              + from  = {
                  + hour   = 12
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 19
                  + minute = 30
                }
            },
          + {
              + day   = "wednesday"
              + from  = null
              + state = "on"
              + to    = null
            },
          + {
              + day   = "thursday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = null
            },
          + {
              + day   = "friday"
              + from  = {
                  + hour   = 0
                  + minute = 0
                }
              + state = "custom"
              + to    = {
                  + hour   = 12
                  + minute = 30
                }
            },
          + {
              + day   = "saturday"
              + from  = {
                  + hour   = 12
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
                  + minute = 0
                }
            },
          + {
              + day   = "sunday"
              + from  = null
              + state = "off"
              + to    = null
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + timezone        = "US/Pacific"
    }
  + collection              = {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + max_ttl         = 200
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }
  + database_credential     = (sensitive value)
  + existing_event          = {
      + alert_key        = "logged_in"
      + app_service_id   = null
      + app_service_name = null
      + cluster_id       = null
      + cluster_name     = null
      + id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + image_url        = null
      + incident_ids     = []
      + key              = "logged_in"
      + kv               = "null"
      + occurrence_count = null
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = null
      + project_name     = null
      + request_id       = null
      + session_id       = null
      + severity         = "info"
      + source           = "cp-ns"
      + summary          = null
      + timestamp        = "2024-12-11 20:33:36.600542622 +0000 UTC"
      + user_email       = null
      + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + user_name        = "SDK QE"
    }
  + existing_events         = {
      + cluster_ids     = null
      + cursor          = {
          + hrefs = {
              + first    = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=10"
              + last     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=37&perPage=10"
              + next     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=10"
              + previous = ""
            }
          + pages = {
              + last        = 37
              + next        = 2
              + page        = 1
              + per_page    = 10
              + previous    = 0
              + total_items = 369
            }
        }
      + data            = [
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:33:36.600542622 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:38:27.417814023 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "7e74fc9b-3d7c-746d-0c36-751a984a5747"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:38:33.290122146 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:38:41.431711495 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "43ce65aa-9d2c-6320-5045-357f8e5a7ec8"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:39:37.974685759 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "1ce5b903-9f45-eea3-a601-2597282a2a09"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:39:42.809020171 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "126012e4-95cb-0002-13be-34b673035f31"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:39:49.230386108 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "9fb9f8bd-787d-6604-5887-b67d34e8be64"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:41:53.158231931 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:41:59.204057247 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
          + {
              + alert_key        = "logged_in"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_in"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = null
              + session_id       = null
              + severity         = "info"
              + source           = "cp-ns"
              + summary          = null
              + timestamp        = "2024-12-11 20:42:05.079807567 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "SDK QE"
            },
        ]
      + from            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = null
      + per_page        = null
      + project_ids     = null
      + severity_levels = null
      + sort_by         = null
      + sort_direction  = null
      + tags            = null
      + to              = null
      + user_ids        = null
    }
  + existing_project_event  = {
      + alert_key        = (known after apply)
      + app_service_id   = (known after apply)
      + app_service_name = (known after apply)
      + cluster_id       = (known after apply)
      + cluster_name     = (known after apply)
      + id               = (known after apply)
      + image_url        = (known after apply)
      + incident_ids     = (known after apply)
      + key              = (known after apply)
      + kv               = (known after apply)
      + occurrence_count = (known after apply)
      + organization_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id       = (known after apply)
      + project_name     = (known after apply)
      + request_id       = (known after apply)
      + session_id       = (known after apply)
      + severity         = (known after apply)
      + source           = (known after apply)
      + summary          = (known after apply)
      + timestamp        = (known after apply)
      + user_email       = (known after apply)
      + user_id          = (known after apply)
      + user_name        = (known after apply)
    }
  + existing_project_events = {
      + cluster_ids     = null
      + cursor          = (known after apply)
      + data            = (known after apply)
      + from            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = null
      + per_page        = null
      + project_id      = (known after apply)
      + severity_levels = null
      + sort_by         = null
      + sort_direction  = null
      + tags            = null
      + to              = null
      + user_ids        = null
    }
  + idx                     = {
      + bucket_name     = "new_terraform_bucket"
      + build_indexes   = null
      + cluster_id      = (known after apply)
      + collection_name = "new_terraform_collection"
      + index_keys      = [
          + "id",
          + "age",
          + "name",
        ]
      + index_name      = "idx1"
      + is_primary      = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + partition_by    = null
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
      + status          = (known after apply)
      + where           = "dept = 'sales'"
      + with            = {
          + defer_build   = false
          + num_partition = null
          + num_replica   = 1
        }
    }
  + network_peer            = {
      + audit           = (known after apply)
      + cluster_id      = (known after apply)
      + commands        = (known after apply)
      + id              = (known after apply)
      + name            = "VPCPeerTFTestAWS"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + provider_config = {
          + aws_config   = {
              + account_id  = "123456789123"
              + cidr        = "10.0.0.0/16"
              + provider_id = (known after apply)
              + region      = "us-east-1"
              + vpc_id      = "vpc-141f0fffff141aa00"
            }
          + azure_config = null
          + gcp_config   = null
        }
      + provider_type   = "aws"
      + status          = (known after apply)
    }
  + new_auditlogsettings    = {
      + audit_enabled     = true
      + cluster_id        = (known after apply)
      + disabled_users    = []
      + enabled_event_ids = (known after apply)
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = (known after apply)
    }
  + organization            = {
      + audit           = {
          + created_at  = "2021-12-03 16:14:45.105347711 +0000 UTC"
          + created_by  = "b1cf2366-0401-4cac-8770-f24e511f6c0a"
          + modified_at = "2024-11-18 15:22:27.933390014 +0000 UTC"
          + modified_by = "ab88be62-e7f5-4701-828d-129e8641d111"
          + version     = 39
        }
      + description     = ""
      + name            = "capella-prod"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + preferences     = {
          + session_duration = 7200
        }
    }
  + project                 = "My First Terraform Project"
  + sample_bucket           = "gamesim-sample"
  + scope                   = {
      + bucket_id       = (known after apply)
      + cluster_id      = (known after apply)
      + collections     = (known after apply)
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = (known after apply)
      + scope_name      = "new_terraform_scope"
    }
  + user                    = {
      + audit                = (known after apply)
      + email                = "johndoe@couchbase.com"
      + enable_notifications = (known after apply)
      + expires_at           = (known after apply)
      + id                   = (known after apply)
      + inactive             = (known after apply)
      + last_login           = (known after apply)
      + name                 = "John Doe"
      + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_roles   = [
          + "organizationMember",
        ]
      + region               = (known after apply)
      + resources            = [
          + {
              + id    = (known after apply)
              + roles = [
                  + "projectDataReaderWriter",
                  + "projectViewer",
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

couchbase-capella_project.new_project: Creating...
couchbase-capella_project.new_project: Creation complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_user.new_user: Creating...
couchbase-capella_apikey.new_apikey: Creating...
couchbase-capella_cluster.new_cluster: Creating...
couchbase-capella_apikey.new_apikey: Creation complete after 1s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_user.new_user: Creation complete after 3s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_cluster.new_cluster: Still creating... [10s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [20s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [30s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [40s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [50s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [1m0s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [1m10s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [1m20s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [1m30s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [1m40s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [1m50s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [2m0s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [2m10s elapsed]
couchbase-capella_cluster.new_cluster: Creation complete after 2m20s [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_certificate.existing_certificate: Reading...
couchbase-capella_allowlist.new_allowlist: Creating...
couchbase-capella_bucket.new_bucket: Creating...
couchbase-capella_sample_bucket.new_sample_bucket: Creating...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Creating...
couchbase-capella_database_credential.new_database_credential: Creating...
data.couchbase-capella_certificate.existing_certificate: Read complete after 0s
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Creation complete after 1s
couchbase-capella_bucket.new_bucket: Creation complete after 7s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_scope.new_scope: Creating...
couchbase-capella_database_credential.new_database_credential: Creation complete after 7s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_scope.new_scope: Creation complete after 0s
couchbase-capella_collection.new_collection: Creating...
couchbase-capella_allowlist.new_allowlist: Creation complete after 8s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_collection.new_collection: Creation complete after 1s
couchbase-capella_query_indexes.idx: Creating...
couchbase-capella_query_indexes.idx: Creation complete after 3s
couchbase-capella_sample_bucket.new_sample_bucket: Still creating... [10s elapsed]
couchbase-capella_sample_bucket.new_sample_bucket: Creation complete after 17s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_app_service.new_app_service: Creating...
couchbase-capella_app_service.new_app_service: Still creating... [10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [1m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [1m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [1m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [1m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [1m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [1m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [2m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [2m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [2m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [2m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [2m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [2m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [3m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [3m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [3m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [3m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [3m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [3m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [4m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [4m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [4m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [4m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [4m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [4m50s elapsed]
couchbase-capella_app_service.new_app_service: Creation complete after 4m55s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_audit_log_export.new_auditlogexport: Creating...
couchbase-capella_audit_log_export.new_auditlogexport: Creation complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_private_endpoint_service.new_service: Creating...
couchbase-capella_private_endpoint_service.new_service: Still creating... [10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [7m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Creation complete after 7m0s

Apply complete! Resources: 16 added, 0 changed, 0 destroyed.

Outputs:

apikey = <sensitive>
app_service = {
  "audit" = {
    "created_at" = "2024-03-27 22:48:25.514164834 +0000 UTC"
    "created_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "modified_at" = "2024-03-27 22:53:18.418467595 +0000 UTC"
    "modified_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "version" = 7
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 7"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 2
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "version" = "3.1.4-1.0.0"
}
bucket = "new_terraform_bucket"
certificate = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "certificate" = <<-EOT
      -----BEGIN CERTIFICATE-----
      MIIDFTCCAf2gAwIBAgIRANguFcFZ7eVLTF2mnPqkkhYwDQYJKoZIhvcNAQELBQAw
      JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEwMTgx
      NDUzMzRaFw0yOTEwMTgxNTUzMzRaMCQxEjAQBgNVBAoMCUNvdWNoYmFzZTEOMAwG
      A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDMoL2G
      1yR4XKOL5KrAZbgJI11NkcooxqCSqoibr5nSM+GNARlou42XbopRhkLQlSMlmH7U
      ZreI7xq2MqmCaQvP1jdS5al/GwuwAP+2kU2nz4IHzliCVV6YvYqNy0fygNpYky9/
      wjCu32n8Ae0AZuxcsAzPUtJBvIIGHum08WlLYS3gNrYkfyds6LfvZvqMk703RL5X
      Ny/RXWmbbBXAXh0chsavEK7EsDLI4t4WI2Iv8+lwS7Wo7Vh6NnEmJLPAAp7udNK4
      U3nwjkL5p/yINROT7CxUE9x0IB2l2rZwZiJhgHCpee77J8QesDut+jZu38ZYY3le
      PS38S81T6I6bSSgtAgMBAAGjQjBAMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYE
      FLlocLdzgAeibrlCmEO4OH5Buf3vMA4GA1UdDwEB/wQEAwIBhjANBgkqhkiG9w0B
      AQsFAAOCAQEAkoVX5CJ7rGx2ALfzy5C7Z+tmEmrZ6jdHjDtw4XwWNhlrsgMuuboU
      Y9XMinSSm1TVfvIz4ru82MVMRxq4v1tPwPdZabbzKYclHkwSMxK5BkyEKWzF1Hoq
      UcinTaT68lVzkTc0D8T+gkRzwXIqxjML2ZdruD1foHNzCgeGHzKzdsjYqrnHv17b
      J+f5tqoa5CKbnyWl3HP0k7r3HHQP0GQequoqXcL3XlERX3Ne20Chck9mftNnHhKw
      Dby7ylZaP97sphqOZQ/W/gza7x1JYylrLXvjfdv3Nmu7oSMKO/2cDyWwcbVGkpbk
      8JOQtFENWmr9u2S0cQfwoCSYBWaK0ofivA==
      -----END CERTIFICATE-----
      EOT
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2024-03-27 22:45:48.671838974 +0000 UTC"
    "created_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "modified_at" = "2024-03-27 22:48:07.649855589 +0000 UTC"
    "modified_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "version" = 5
  }
  "availability" = {
    "type" = "multi"
  }
  "cloud_provider" = {
    "cidr" = "10.5.30.0/23"
    "region" = "us-east-1"
    "type" = "aws"
  }
  "configuration_type" = "multiNode"
  "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
  "couchbase_server" = {
    "version" = "7.2"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "if_match" = tostring(null)
  "name" = "My First Terraform Cluster"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "service_groups" = toset([
    {
      "node" = {
        "compute" = {
          "cpu" = 4
          "ram" = 16
        }
        "disk" = {
          "autoexpansion" = tobool(null)
          "iops" = 5000
          "storage" = 50
          "type" = "io2"
        }
      }
      "num_of_nodes" = 3
      "services" = toset([
        "data",
        "index",
        "query",
      ])
    },
  ])
  "support" = {
    "plan" = "enterprise"
    "timezone" = "PT"
  }
}
cluster_onoff_schedule = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "days" = tolist([
    {
      "day" = "monday"
      "from" = {
        "hour" = 12
        "minute" = 30
      }
      "state" = "custom"
      "to" = {
        "hour" = 14
        "minute" = 30
      }
    },
    {
      "day" = "tuesday"
      "from" = {
        "hour" = 12
        "minute" = 0
      }
      "state" = "custom"
      "to" = {
        "hour" = 19
        "minute" = 30
      }
    },
    {
      "day" = "wednesday"
      "from" = null /* object */
      "state" = "on"
      "to" = null /* object */
    },
    {
      "day" = "thursday"
      "from" = {
        "hour" = 12
        "minute" = 30
      }
      "state" = "custom"
      "to" = null /* object */
    },
    {
      "day" = "friday"
      "from" = {
        "hour" = 0
        "minute" = 0
      }
      "state" = "custom"
      "to" = {
        "hour" = 12
        "minute" = 30
      }
    },
    {
      "day" = "saturday"
      "from" = {
        "hour" = 12
        "minute" = 30
      }
      "state" = "custom"
      "to" = {
        "hour" = 14
        "minute" = 0
      }
    },
    {
      "day" = "sunday"
      "from" = null /* object */
      "state" = "off"
      "to" = null /* object */
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "timezone" = "US/Pacific"
}
collection = {
  "bucket_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "new_terraform_collection"
  "max_ttl" = 200
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "new_terraform_scope"
}
database_credential = <sensitive>
network_peer = {
  "audit" = {
    "created_at" = "2024-06-27 19:34:00.202552084 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-27 19:34:18.661844385 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
 "commands" = toset([
    "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-12345678912345678",
    "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-141f0fffff141aa00ff,VPCRegion=us-east-1 --region=us-east-1",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestAWS"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = {
      "account_id" = "123456789123"
      "cidr" = "10.1.0.0/23"
      "provider_id" = "pcx-12345678912345678"
      "region" = "us-east-1"
      "vpc_id" = "vpc-141f0fffff141aa00ff"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "aws"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
peer_id = "ffffffff-aaaa-1414-eeee-000000000000"
organization = {
  "audit" = {
    "created_at" = "2020-07-22 12:38:57.437248116 +0000 UTC"
    "created_by" = ""
    "modified_at" = "2024-03-23 20:41:47.693734149 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 0
  }
  "description" = ""
  "name" = "cbc-dev"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "preferences" = {
    "session_duration" = 7200
  }
}
project = "My First Terraform Project"
sample_bucket = "gamesim-sample"
scope = {
  "bucket_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections" = toset([])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "new_terraform_scope"
}
user = {
  "audit" = {
    "created_at" = "2024-03-27 22:45:49.332340465 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2024-03-27 22:45:49.332340465 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-06-25T22:45:49.332341187Z"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "inactive" = true
  "last_login" = ""
  "name" = "John Doe"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = toset([
    {
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "roles" = toset([
        "projectDataReaderWriter",
        "projectViewer",
      ])
      "type" = "project"
    },
  ])
  "status" = "not-verified"
  "time_zone" = ""
}
new_auditlogsettings = {
  "audit_enabled" = true
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = toset([])
  "enabled_event_ids" = toset([
    28672,
    28673,
    28674,
    28675,
    28676,
    28677,
    28678,
    28679,
    28680,
    28681,
    28682,
    28683,
    28684,
    28685,
    28686,
    28687,
    28688,
    28689,
    28690,
    28691,
    28692,
    28693,
    28694,
    28695,
    28697,
    28698,
    28699,
    28700,
    28701,
    28702,
    28704,
    28705,
    28706,
    28707,
    28708,
    28709,
    28710,
    28711,
    28712,
    28713,
    28714,
    28715,
    28716,
    28717,
    28718,
    28719,
    28720,
    28721,
    28722,
    28723,
    28724,
    28725,
    28726,
    28727,
    28728,
    28729,
    28730,
    28731,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}

service_status = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "enabled" = false
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

Finally, you can view the outputs using the `terraform output` command

Sample Output:
```
$ terraform output
apikey = <sensitive>
app_service = {
  "audit" = {
    "created_at" = "2024-03-27 22:48:25.514164834 +0000 UTC"
    "created_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "modified_at" = "2024-03-27 22:53:18.418467595 +0000 UTC"
    "modified_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "version" = 7
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 7"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 2
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "version" = "3.1.4-1.0.0"
}
bucket = "new_terraform_bucket"
certificate = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "certificate" = <<-EOT
      -----BEGIN CERTIFICATE-----
      MIIDFTCCAf2gAwIBAgIRANguFcFZ7eVLTF2mnPqkkhYwDQYJKoZIhvcNAQELBQAw
      JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEwMTgx
      NDUzMzRaFw0yOTEwMTgxNTUzMzRaMCQxEjAQBgNVBAoMCUNvdWNoYmFzZTEOMAwG
      A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDMoL2G
      1yR4XKOL5KrAZbgJI11NkcooxqCSqoibr5nSM+GNARlou42XbopRhkLQlSMlmH7U
      ZreI7xq2MqmCaQvP1jdS5al/GwuwAP+2kU2nz4IHzliCVV6YvYqNy0fygNpYky9/
      wjCu32n8Ae0AZuxcsAzPUtJBvIIGHum08WlLYS3gNrYkfyds6LfvZvqMk703RL5X
      Ny/RXWmbbBXAXh0chsavEK7EsDLI4t4WI2Iv8+lwS7Wo7Vh6NnEmJLPAAp7udNK4
      U3nwjkL5p/yINROT7CxUE9x0IB2l2rZwZiJhgHCpee77J8QesDut+jZu38ZYY3le
      PS38S81T6I6bSSgtAgMBAAGjQjBAMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYE
      FLlocLdzgAeibrlCmEO4OH5Buf3vMA4GA1UdDwEB/wQEAwIBhjANBgkqhkiG9w0B
      AQsFAAOCAQEAkoVX5CJ7rGx2ALfzy5C7Z+tmEmrZ6jdHjDtw4XwWNhlrsgMuuboU
      Y9XMinSSm1TVfvIz4ru82MVMRxq4v1tPwPdZabbzKYclHkwSMxK5BkyEKWzF1Hoq
      UcinTaT68lVzkTc0D8T+gkRzwXIqxjML2ZdruD1foHNzCgeGHzKzdsjYqrnHv17b
      J+f5tqoa5CKbnyWl3HP0k7r3HHQP0GQequoqXcL3XlERX3Ne20Chck9mftNnHhKw
      Dby7ylZaP97sphqOZQ/W/gza7x1JYylrLXvjfdv3Nmu7oSMKO/2cDyWwcbVGkpbk
      8JOQtFENWmr9u2S0cQfwoCSYBWaK0ofivA==
      -----END CERTIFICATE-----
      EOT
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2024-03-27 22:45:48.671838974 +0000 UTC"
    "created_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "modified_at" = "2024-03-27 22:48:07.649855589 +0000 UTC"
    "modified_by" = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
    "version" = 5
  }
  "availability" = {
    "type" = "multi"
  }
  "cloud_provider" = {
    "cidr" = "10.5.30.0/23"
    "region" = "us-east-1"
    "type" = "aws"
  }
  "configuration_type" = "multiNode"
  "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
  "couchbase_server" = {
    "version" = "7.2"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "if_match" = tostring(null)
  "name" = "My First Terraform Cluster"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "service_groups" = toset([
    {
      "node" = {
        "compute" = {
          "cpu" = 4
          "ram" = 16
        }
        "disk" = {
          "autoexpansion" = tobool(null)
          "iops" = 5000
          "storage" = 50
          "type" = "io2"
        }
      }
      "num_of_nodes" = 3
      "services" = toset([
        "data",
        "index",
        "query",
      ])
    },
  ])
  "support" = {
    "plan" = "enterprise"
    "timezone" = "PT"
  }
}
cluster_onoff_schedule = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "days" = tolist([
    {
      "day" = "monday"
      "from" = {
        "hour" = 12
        "minute" = 30
      }
      "state" = "custom"
      "to" = {
        "hour" = 14
        "minute" = 30
      }
    },
    {
      "day" = "tuesday"
      "from" = {
        "hour" = 12
        "minute" = 0
      }
      "state" = "custom"
      "to" = {
        "hour" = 19
        "minute" = 30
      }
    },
    {
      "day" = "wednesday"
      "from" = null /* object */
      "state" = "on"
      "to" = null /* object */
    },
    {
      "day" = "thursday"
      "from" = {
        "hour" = 12
        "minute" = 30
      }
      "state" = "custom"
      "to" = null /* object */
    },
    {
      "day" = "friday"
      "from" = {
        "hour" = 0
        "minute" = 0
      }
      "state" = "custom"
      "to" = {
        "hour" = 12
        "minute" = 30
      }
    },
    {
      "day" = "saturday"
      "from" = {
        "hour" = 12
        "minute" = 30
      }
      "state" = "custom"
      "to" = {
        "hour" = 14
        "minute" = 0
      }
    },
    {
      "day" = "sunday"
      "from" = null /* object */
      "state" = "off"
      "to" = null /* object */
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "timezone" = "US/Pacific"
}
collection = {
  "bucket_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "new_terraform_collection"
  "max_ttl" = 200
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "new_terraform_scope"
}
database_credential = <sensitive>
organization = {
  "audit" = {
    "created_at" = "2020-07-22 12:38:57.437248116 +0000 UTC"
    "created_by" = ""
    "modified_at" = "2024-03-23 20:41:47.693734149 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 0
  }
  "description" = ""
  "name" = "cbc-dev"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "preferences" = {
    "session_duration" = 7200
  }
}
network_peer = {
  "audit" = {
    "created_at" = "2024-06-27 19:34:00.202552084 +0000 UTC"
    "created_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "modified_at" = "2024-06-27 19:34:18.661844385 +0000 UTC"
    "modified_by" = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
    "version" = 2
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "commands" = toset([
    "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-12345678912345678",
    "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-141f0fffff141aa00ff,VPCRegion=us-east-1 --region=us-east-1",
  ])
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "VPCPeerTFTestAWS"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws_config" = {
      "account_id" = "123456789123"
      "cidr" = "10.1.0.0/23"
      "provider_id" = "pcx-12345678912345678"
      "region" = "us-east-1"
      "vpc_id" = "vpc-141f0fffff141aa00ff"
    }
    "gcp_config" = null /* object */
  }
  "provider_type" = "aws"
  "status" = {
    "reasoning" = ""
    "state" = "complete"
  }
}
project = "My First Terraform Project"
sample_bucket = "gamesim-sample"
scope = {
  "bucket_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections" = toset([])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "new_terraform_scope"
}
user = {
  "audit" = {
    "created_at" = "2024-03-27 22:45:49.332340465 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2024-03-27 22:45:49.332340465 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-06-25T22:45:49.332341187Z"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "inactive" = true
  "last_login" = ""
  "name" = "John Doe"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_roles" = tolist([
    "organizationMember",
  ])
  "region" = ""
  "resources" = toset([
    {
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "roles" = toset([
        "projectDataReaderWriter",
        "projectViewer",
      ])
      "type" = "project"
    },
  ])
  "status" = "not-verified"
  "time_zone" = ""
}
new_auditlogsettings = {
  "audit_enabled" = true
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = toset([])
  "enabled_event_ids" = toset([
    28672,
    28673,
    28674,
    28675,
    28676,
    28677,
    28678,
    28679,
    28680,
    28681,
    28682,
    28683,
    28684,
    28685,
    28686,
    28687,
    28688,
    28689,
    28690,
    28691,
    28692,
    28693,
    28694,
    28695,
    28697,
    28698,
    28699,
    28700,
    28701,
    28702,
    28704,
    28705,
    28706,
    28707,
    28708,
    28709,
    28710,
    28711,
    28712,
    28713,
    28714,
    28715,
    28716,
    28717,
    28718,
    28719,
    28720,
    28721,
    28722,
    28723,
    28724,
    28725,
    28726,
    28727,
    28728,
    28729,
    28730,
    28731,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}

service_status = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "enabled" = false
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

All these resources (except audit log settings) can be destroyed using the `terraform destroy` command

Sample Output:
```
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_organization.existing_organization: Reading...
couchbase-capella_project.new_project: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_organization.existing_organization: Read complete after 1s [name=my-org]
couchbase-capella_apikey.new_apikey: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_user.new_user: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_cluster.new_cluster: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_allowlist.new_allowlist: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_database_credential.new_database_credential: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_certificate.existing_certificate: Reading...
couchbase-capella_sample_bucket.new_sample_bucket: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_bucket.new_bucket: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Refreshing state...
data.couchbase-capella_certificate.existing_certificate: Read complete after 0s
couchbase-capella_scope.new_scope: Refreshing state...
couchbase-capella_app_service.new_app_service: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_collection.new_collection: Refreshing state...
data.couchbase-capella_private_endpoint_service.service_status: Reading...
couchbase-capella_private_endpoint_service.new_service: Refreshing state...
data.couchbase-capella_private_endpoint_service.service_status: Read complete after 0s
couchbase-capella_query_indexes.idx: Refreshing state...
couchbase-capella_audit_log_settings.new_auditlogsettings: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_allowlist.new_allowlist will be destroyed
  - resource "couchbase-capella_allowlist" "new_allowlist" {
      - audit           = {
          - created_at  = "2024-03-27 22:48:15.382200967 +0000 UTC" -> null
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - modified_at = "2024-03-27 22:48:15.382200967 +0000 UTC" -> null
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - version     = 1 -> null
        } -> null
      - cidr            = "8.8.8.8/32" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - comment         = "Allow access from a public IP" -> null
      - expires_at      = "2043-11-30T23:59:59.465Z" -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }

  # couchbase-capella_apikey.new_apikey will be destroyed
  - resource "couchbase-capella_apikey" "new_apikey" {
      - allowed_cidrs      = [
          - "10.1.42.0/23",
          - "10.1.43.0/23",
          - "10.5.30.0/23",
        ] -> null
      - audit              = {
          - created_at  = "2024-03-27 22:45:48.449403667 +0000 UTC" -> null
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - modified_at = "2024-03-27 22:45:48.449403667 +0000 UTC" -> null
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - version     = 1 -> null
        } -> null
      - expiry             = 180 -> null
      - id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name               = "My First Terraform API Key" -> null
      - organization_id    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - organization_roles = [
          - "organizationOwner",
        ] -> null
      - resources          = [
          - {
              - id    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
              - roles = [
                  - "projectDataReader",
                  - "projectManager",
                ] -> null
              - type  = "project" -> null
            },
        ] -> null
      - token              = (sensitive value) -> null
    }

  # couchbase-capella_app_service.new_app_service will be destroyed
  - resource "couchbase-capella_app_service" "new_app_service" {
      - audit           = {
          - created_at  = "2024-03-27 22:48:25.514164834 +0000 UTC" -> null
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - modified_at = "2024-03-27 22:53:18.418467595 +0000 UTC" -> null
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - version     = 7 -> null
        } -> null
      - cloud_provider  = "AWS" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - compute         = {
          - cpu = 2 -> null
          - ram = 4 -> null
        } -> null
      - current_state   = "healthy" -> null
      - description     = "My first test app service." -> null
      - etag            = "Version: 7" -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name            = "new-terraform-app-service" -> null
      - nodes           = 2 -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - version         = "3.1.4-1.0.0" -> null
    }

  # couchbase-capella_bucket.new_bucket will be destroyed
  - resource "couchbase-capella_bucket" "new_bucket" {
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - durability_level           = "none" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = false -> null
      - id                         = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - memory_allocation_in_mb    = 100 -> null
      - name                       = "new_terraform_bucket" -> null
      - organization_id            = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - replicas                   = 1 -> null
      - stats                      = {
          - disk_used_in_mib   = 34 -> null
          - item_count         = 0 -> null
          - memory_used_in_mib = 55 -> null
          - ops_per_second     = 0 -> null
        } -> null
      - storage_backend            = "couchstore" -> null
      - time_to_live_in_seconds    = 0 -> null
      - type                       = "couchbase" -> null
    }

  # couchbase-capella_cluster.new_cluster will be destroyed
  - resource "couchbase-capella_cluster" "new_cluster" {
      - audit              = {
          - created_at  = "2024-03-27 22:45:48.671838974 +0000 UTC" -> null
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - modified_at = "2024-03-27 22:48:07.649855589 +0000 UTC" -> null
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - version     = 5 -> null
        } -> null
      - availability       = {
          - type = "multi" -> null
        } -> null
      - cloud_provider     = {
          - cidr   = "10.5.30.0/23" -> null
          - region = "us-east-1" -> null
          - type   = "aws" -> null
        } -> null
      - configuration_type = "multiNode" -> null
      - connection_string  = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com" -> null
      - couchbase_server   = {
          - version = "7.2" -> null
        } -> null
      - current_state      = "healthy" -> null
      - description        = "My first test cluster for multiple services." -> null
      - etag               = "Version: 5" -> null
      - id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name               = "My First Terraform Cluster" -> null
      - organization_id    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id         = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - service_groups     = [
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
      - support            = {
          - plan     = "enterprise" -> null
          - timezone = "PT" -> null
        } -> null
    }

  # couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule will be destroyed
  - resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - days            = [
          - {
              - day   = "monday" -> null
              - from  = {
                  - hour   = 12 -> null
                  - minute = 30 -> null
                } -> null
              - state = "custom" -> null
              - to    = {
                  - hour   = 14 -> null
                  - minute = 30 -> null
                } -> null
            },
          - {
              - day   = "tuesday" -> null
              - from  = {
                  - hour   = 12 -> null
                  - minute = 0 -> null
                } -> null
              - state = "custom" -> null
              - to    = {
                  - hour   = 19 -> null
                  - minute = 30 -> null
                } -> null
            },
          - {
              - day   = "wednesday" -> null
              - state = "on" -> null
            },
          - {
              - day   = "thursday" -> null
              - from  = {
                  - hour   = 12 -> null
                  - minute = 30 -> null
                } -> null
              - state = "custom" -> null
            },
          - {
              - day   = "friday" -> null
              - from  = {
                  - hour   = 0 -> null
                  - minute = 0 -> null
                } -> null
              - state = "custom" -> null
              - to    = {
                  - hour   = 12 -> null
                  - minute = 30 -> null
                } -> null
            },
          - {
              - day   = "saturday" -> null
              - from  = {
                  - hour   = 12 -> null
                  - minute = 30 -> null
                } -> null
              - state = "custom" -> null
              - to    = {
                  - hour   = 14 -> null
                  - minute = 0 -> null
                } -> null
            },
          - {
              - day   = "sunday" -> null
              - state = "off" -> null
            },
        ] -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - timezone        = "US/Pacific" -> null
    }

  # couchbase-capella_collection.new_collection will be destroyed
  - resource "couchbase-capella_collection" "new_collection" {
      - bucket_id       = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "new_terraform_collection" -> null
      - max_ttl         = 200 -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "new_terraform_scope" -> null
    }

  # couchbase-capella_database_credential.new_database_credential will be destroyed
  - resource "couchbase-capella_database_credential" "new_database_credential" {
      - access          = [
          - {
              - privileges = [
                  - "data_reader",
                ] -> null
            },
          - {
              - privileges = [
                  - "data_writer",
                ] -> null
              - resources  = {
                  - buckets = [
                      - {
                          - name   = "new_terraform_bucket" -> null
                          - scopes = [
                              - {
                                  - collections = [
                                      - "_default",
                                    ] -> null
                                  - name        = "_default" -> null
                                },
                            ] -> null
                        },
                    ] -> null
                } -> null
            },
        ] -> null
      - audit           = {
          - created_at  = "2024-03-27 22:48:14.886768871 +0000 UTC" -> null
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - modified_at = "2024-03-27 22:48:14.886768871 +0000 UTC" -> null
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - version     = 1 -> null
        } -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name            = "terraform_db_credential" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - password        = (sensitive value) -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }
 
  Terraform will perform the following actions:

  # couchbase-capella_network_peer.new_network_peer will be destroyed
  - resource "couchbase-capella_network_peer" "new_network_peer" {
      - audit           = {
          - created_at  = "2024-06-29 00:37:45.4338168 +0000 UTC" -> null
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> null
          - modified_at = "2024-06-29 00:38:04.635673378 +0000 UTC" -> null
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi" -> null
          - version     = 2 -> null
        } -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
       - commands        = [
          - "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-1234567891234567",
          - "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-12345678912345678,VPCRegion=us-east-1 --region=us-east-1",
        ] -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name            = "VPCPeerTFTestAWS" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - provider_config = {
          - aws_config = {
              - account_id  = "123456789123" -> null
              - cidr        = "10.2.0.0/23" -> null
              - provider_id = "pcx-1234567891234567" -> null
              - region      = "us-east-1" -> null
              - vpc_id      = "vpc-12345678912345678" -> null
            } -> null
        } -> null
      - provider_type   = "aws" -> null
      - status          = {
          - reasoning = "" -> null
          - state     = "complete" -> null
        } -> null
    }

  # couchbase-capella_project.new_project will be destroyed
  - resource "couchbase-capella_project" "new_project" {
      - audit           = {
          - created_at  = "2024-03-27 22:45:48.038234717 +0000 UTC" -> null
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - modified_at = "2024-03-27 22:45:48.038250937 +0000 UTC" -> null
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT" -> null
          - version     = 1 -> null
        } -> null
      - description     = "A Capella Project that will host many Capella clusters." -> null
      - etag            = "Version: 1" -> null
      - id              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name            = "My First Terraform Project" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }
    
    # couchbase-capella_query_indexes.idx will be destroyed
  - resource "couchbase-capella_query_indexes" "idx" {
      - bucket_name     = "new_terraform_bucket" -> null
      - cluster_id      = "dad8fd4b-7df7-4851-b9b9-7dc7f1fd8545" -> null
      - collection_name = "new_terraform_collection" -> null
      - index_keys      = [
          - "id",
          - "age",
          - "name",
        ] -> null
      - index_name      = "idx1" -> null
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - project_id      = "6717e550-0620-4a16-b472-b4388a0582b3" -> null
      - scope_name      = "new_terraform_scope" -> null
      - where           = "dept = 'sales'" -> null
      - with            = {
          - defer_build = false -> null
          - num_replica = 1 -> null
        } -> null
    }

  # couchbase-capella_sample_bucket.new_sample_bucket will be destroyed
  - resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - durability_level           = "none" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = false -> null
      - id                         = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - memory_allocation_in_mb    = 200 -> null
      - name                       = "gamesim-sample" -> null
      - organization_id            = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - replicas                   = 1 -> null
      - stats                      = {
          - disk_used_in_mib   = 20 -> null
          - item_count         = 586 -> null
          - memory_used_in_mib = 62 -> null
          - ops_per_second     = 2 -> null
        } -> null
      - storage_backend            = "couchstore" -> null
      - time_to_live_in_seconds    = 0 -> null
      - type                       = "couchbase" -> null
    }

  # couchbase-capella_scope.new_scope will be destroyed
  - resource "couchbase-capella_scope" "new_scope" {
      - bucket_id       = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collections     = [
          - {
              - max_ttl = 200 -> null
              - name    = "new_terraform_collection" -> null
            },
        ] -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "new_terraform_scope" -> null
    }

  # couchbase-capella_user.new_user will be destroyed
  - resource "couchbase-capella_user" "new_user" {
      - audit                = {
          - created_at  = "2024-03-27 22:45:49.332340465 +0000 UTC" -> null
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
          - modified_at = "2024-03-27 22:45:49.332340465 +0000 UTC" -> null
          - modified_by = "ffffffff-aaaa-1414-eeee-000000000000" -> null
          - version     = 1 -> null
        } -> null
      - email                = "johndoe@couchbase.com" -> null
      - enable_notifications = false -> null
      - expires_at           = "2024-06-25T22:45:49.332341187Z" -> null
      - id                   = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - inactive             = true -> null
      - name                 = "John Doe" -> null
      - organization_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - organization_roles   = [
          - "organizationMember",
        ] -> null
      - resources            = [
          - {
              - id    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
              - roles = [
                  - "projectDataReaderWriter",
                  - "projectViewer",
                ] -> null
              - type  = "project" -> null
            },
        ] -> null
      - status               = "not-verified" -> null
    }
  # couchbase-capella_audit_log_settings.new_auditlogsettings will be destroyed
  - resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      - audit_enabled     = true -> null
      - cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - disabled_users    = [] -> null
      - enabled_event_ids = [
          - 28672,
          - 28673,
          - 28674,
          - 28675,
          - 28676,
          - 28677,
          - 28678,
          - 28679,
          - 28680,
          - 28681,
          - 28682,
          - 28683,
          - 28684,
          - 28685,
          - 28686,
          - 28687,
          - 28688,
          - 28689,
          - 28690,
          - 28691,
          - 28692,
          - 28693,
          - 28694,
          - 28695,
          - 28697,
          - 28698,
          - 28699,
          - 28700,
          - 28701,
          - 28702,
          - 28704,
          - 28705,
          - 28706,
          - 28707,
          - 28708,
          - 28709,
          - 28710,
          - 28711,
          - 28712,
          - 28713,
          - 28714,
          - 28715,
          - 28716,
          - 28717,
          - 28718,
          - 28719,
          - 28720,
          - 28721,
          - 28722,
          - 28723,
          - 28724,
          - 28725,
          - 28726,
          - 28727,
          - 28728,
          - 28729,
          - 28730,
          - 28731,
        ] -> null
      - organization_id   = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }
    
  # couchbase-capella_private_endpoint_service.new_service will be destroyed
  - resource "couchbase-capella_private_endpoint_service" "new_service" {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - enabled         = true -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }

Plan: 0 to add, 0 to change, 16 to destroy.

Changes to Outputs:
  - apikey                 = (sensitive value) -> null
  - app_service            = {
      - audit           = {
          - created_at  = "2024-03-27 22:48:25.514164834 +0000 UTC"
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
          - modified_at = "2024-03-27 22:53:18.418467595 +0000 UTC"
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
          - version     = 7
        }
      - cloud_provider  = "AWS"
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - compute         = {
          - cpu = 2
          - ram = 4
        }
      - current_state   = "healthy"
      - description     = "My first test app service."
      - etag            = "Version: 7"
      - id              = "ffffffff-aaaa-1414-eeee-000000000000"
      - if_match        = null
      - name            = "new-terraform-app-service"
      - nodes           = 2
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - version         = "3.1.4-1.0.0"
    } -> null
  - bucket                 = "new_terraform_bucket" -> null
  - certificate            = {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - data            = [
          - {
              - certificate = <<-EOT
                    -----BEGIN CERTIFICATE-----
                    MIIDFTCCAf2gAwIBAgIRANguFcFZ7eVLTF2mnPqkkhYwDQYJKoZIhvcNAQELBQAw
                    JDESMBAGA1UECgwJQ291Y2hiYXNlMQ4wDAYDVQQLDAVDbG91ZDAeFw0xOTEwMTgx
                    NDUzMzRaFw0yOTEwMTgxNTUzMzRaMCQxEjAQBgNVBAoMCUNvdWNoYmFzZTEOMAwG
                    A1UECwwFQ2xvdWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDMoL2G
                    1yR4XKOL5KrAZbgJI11NkcooxqCSqoibr5nSM+GNARlou42XbopRhkLQlSMlmH7U
                    ZreI7xq2MqmCaQvP1jdS5al/GwuwAP+2kU2nz4IHzliCVV6YvYqNy0fygNpYky9/
                    wjCu32n8Ae0AZuxcsAzPUtJBvIIGHum08WlLYS3gNrYkfyds6LfvZvqMk703RL5X
                    Ny/RXWmbbBXAXh0chsavEK7EsDLI4t4WI2Iv8+lwS7Wo7Vh6NnEmJLPAAp7udNK4
                    U3nwjkL5p/yINROT7CxUE9x0IB2l2rZwZiJhgHCpee77J8QesDut+jZu38ZYY3le
                    PS38S81T6I6bSSgtAgMBAAGjQjBAMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYE
                    FLlocLdzgAeibrlCmEO4OH5Buf3vMA4GA1UdDwEB/wQEAwIBhjANBgkqhkiG9w0B
                    AQsFAAOCAQEAkoVX5CJ7rGx2ALfzy5C7Z+tmEmrZ6jdHjDtw4XwWNhlrsgMuuboU
                    Y9XMinSSm1TVfvIz4ru82MVMRxq4v1tPwPdZabbzKYclHkwSMxK5BkyEKWzF1Hoq
                    UcinTaT68lVzkTc0D8T+gkRzwXIqxjML2ZdruD1foHNzCgeGHzKzdsjYqrnHv17b
                    J+f5tqoa5CKbnyWl3HP0k7r3HHQP0GQequoqXcL3XlERX3Ne20Chck9mftNnHhKw
                    Dby7ylZaP97sphqOZQ/W/gza7x1JYylrLXvjfdv3Nmu7oSMKO/2cDyWwcbVGkpbk
                    8JOQtFENWmr9u2S0cQfwoCSYBWaK0ofivA==
                    -----END CERTIFICATE-----
                EOT
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null
  - cluster                = {
      - app_service_id     = null
      - audit              = {
          - created_at  = "2024-03-27 22:45:48.671838974 +0000 UTC"
          - created_by  = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
          - modified_at = "2024-03-27 22:48:07.649855589 +0000 UTC"
          - modified_by = "tEsTNwA3pfwAY3MgddFfbHTDs5tEsT"
          - version     = 5
        }
      - availability       = {
          - type = "multi"
        }
      - cloud_provider     = {
          - cidr   = "10.5.30.0/23"
          - region = "us-east-1"
          - type   = "aws"
        }
      - configuration_type = "multiNode"
      - connection_string  = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
      - couchbase_server   = {
          - version = "7.2"
        }
      - current_state      = "healthy"
      - description        = "My first test cluster for multiple services."
      - etag               = "Version: 5"
      - id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      - if_match           = null
      - name               = "My First Terraform Cluster"
      - organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      - service_groups     = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 4
                      - ram = 16
                    }
                  - disk    = {
                      - autoexpansion = null
                      - iops          = 5000
                      - storage       = 50
                      - type          = "io2"
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
      - support            = {
          - plan     = "enterprise"
          - timezone = "PT"
        }
    } -> null
  - cluster_onoff_schedule = {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - days            = [
          - {
              - day   = "monday"
              - from  = {
                  - hour   = 12
                  - minute = 30
                }
              - state = "custom"
              - to    = {
                  - hour   = 14
                  - minute = 30
                }
            },
          - {
              - day   = "tuesday"
              - from  = {
                  - hour   = 12
                  - minute = 0
                }
              - state = "custom"
              - to    = {
                  - hour   = 19
                  - minute = 30
                }
            },
          - {
              - day   = "wednesday"
              - from  = null
              - state = "on"
              - to    = null
            },
          - {
              - day   = "thursday"
              - from  = {
                  - hour   = 12
                  - minute = 30
                }
              - state = "custom"
              - to    = null
            },
          - {
              - day   = "friday"
              - from  = {
                  - hour   = 0
                  - minute = 0
                }
              - state = "custom"
              - to    = {
                  - hour   = 12
                  - minute = 30
                }
            },
          - {
              - day   = "saturday"
              - from  = {
                  - hour   = 12
                  - minute = 30
                }
              - state = "custom"
              - to    = {
                  - hour   = 14
                  - minute = 0
                }
            },
          - {
              - day   = "sunday"
              - from  = null
              - state = "off"
              - to    = null
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - timezone        = "US/Pacific"
    } -> null
  - collection             = {
      - bucket_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - collection_name = "new_terraform_collection"
      - max_ttl         = 200
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - scope_name      = "new_terraform_scope"
    } -> null
  - database_credential    = (sensitive value) -> null
  - network_peer   = {
      - audit           = {
          - created_at  = "2024-06-29 00:37:45.4338168 +0000 UTC"
          - created_by  = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - modified_at = "2024-06-29 00:38:04.635673378 +0000 UTC"
          - modified_by = "OqrQXTvtsD6PjHiP9Tt4ZhggaCzQVpPi"
          - version     = 2
        }
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - commands        = [
          - "aws ec2 accept-vpc-peering-connection --region=us-east-1 --vpc-peering-connection-id=pcx-1234567891234567",
          - "aws route53 associate-vpc-with-hosted-zone --hosted-zone-id=AAAAA000000FFFFFAAAAAA --vpc=VPCId=vpc-12345678912345678,VPCRegion=us-east-1 --region=us-east-1",
        ]
      - id              = "ffffffff-aaaa-1414-eeee-000000000000"
      - name            = "VPCPeerTFTestAWS"
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - provider_config = {
          - aws_config = {
              - account_id  = "123456789123"
              - cidr        = "10.2.0.0/23"
              - provider_id = "pcx-1234567891234567"
              - region      = "us-east-1"
              - vpc_id      = "vpc-12345678912345678"
            }
          - gcp_config = null
        }
      - provider_type   = "aws"
      - status          = {
          - reasoning = ""
          - state     = "complete"
        }
    } -> null
  - peer_id            = "ffffffff-aaaa-1414-eeee-000000000000" -> null
  - organization           = {
      - audit           = {
          - created_at  = "2020-07-22 12:38:57.437248116 +0000 UTC"
          - created_by  = ""
          - modified_at = "2024-03-23 20:41:47.693734149 +0000 UTC"
          - modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
          - version     = 0
        }
      - description     = ""
      - name            = "cbc-dev"
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - preferences     = {
          - session_duration = 7200
        }
    } -> null
  - project                = "My First Terraform Project" -> null
  - sample_bucket          = "gamesim-sample" -> null
  - scope                  = {
      - bucket_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - collections     = [
          - {
              - max_ttl = 200
              - name    = "new_terraform_collection"
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - scope_name      = "new_terraform_scope"
    } -> null
  - user                   = {
      - audit                = {
          - created_at  = "2024-03-27 22:45:49.332340465 +0000 UTC"
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          - modified_at = "2024-03-27 22:45:49.332340465 +0000 UTC"
          - modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
          - version     = 1
        }
      - email                = "johndoe@couchbase.com"
      - enable_notifications = false
      - expires_at           = "2024-06-25T22:45:49.332341187Z"
      - id                   = "ffffffff-aaaa-1414-eeee-000000000000"
      - inactive             = true
      - last_login           = ""
      - name                 = "John Doe"
      - organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - organization_roles   = [
          - "organizationMember",
        ]
      - region               = ""
      - resources            = [
          - {
              - id    = "ffffffff-aaaa-1414-eeee-000000000000"
              - roles = [
                  - "projectDataReaderWriter",
                  - "projectViewer",
                ]
              - type  = "project"
            },
        ]
      - status               = "not-verified"
      - time_zone            = ""
    } -> null
    
    - service_status = {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - enabled         = true
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_allowlist.new_allowlist: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_user.new_user: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_collection.new_collection: Destroying...
couchbase-capella_app_service.new_app_service: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_apikey.new_apikey: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_database_credential.new_database_credential: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Destroying...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Destruction complete after 0s
couchbase-capella_query_indexes.idx: Destruction complete after 1s
couchbase-capella_apikey.new_apikey: Destruction complete after 1s
couchbase-capella_database_credential.new_database_credential: Destruction complete after 1s
couchbase-capella_collection.new_collection: Destruction complete after 1s
couchbase-capella_scope.new_scope: Destroying...
couchbase-capella_scope.new_scope: Destruction complete after 0s
couchbase-capella_user.new_user: Destruction complete after 3s
couchbase-capella_allowlist.new_allowlist: Destruction complete after 8s
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 10s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 20s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 30s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 40s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 50s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m0s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m10s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m20s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m30s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m40s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m50s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m0s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m10s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m20s elapsed]
couchbase-capella_app_service.new_app_service: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m30s elapsed]
couchbase-capella_app_service.new_app_service: Destruction complete after 2m40s
couchbase-capella_sample_bucket.new_sample_bucket: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_bucket.new_bucket: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_bucket.new_bucket: Destruction complete after 1s
couchbase-capella_sample_bucket.new_sample_bucket: Destruction complete after 2s
couchbase-capella_query_indexes.idx: Destroying...
couchbase-capella_query_indexes.idx: Destruction complete after 1s
couchbase-capella_cluster.new_cluster: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 3m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 3m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 3m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 3m30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 3m40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 3m50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 4m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 4m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 4m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 4m30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 4m40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 4m50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 5m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 5m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 5m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 5m30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 5m40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 5m50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 6m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 6m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 6m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 6m30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 6m40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 6m50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 7m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 7m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 7m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 7m30s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 7m40s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 7m50s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 8m0s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 8m10s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 8m20s elapsed]
couchbase-capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 8m30s elapsed]
couchbase-capella_cluster.new_cluster: Destruction complete after 8m40s
couchbase-capella_project.new_project: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_project.new_project: Destruction complete after 2s
couchbase-capella_audit_log_settings.new_auditlogsettings: Destroying...
couchbase-capella_audit_log_settings.new_auditlogsettings: Destruction complete after 0s
couchbase-capella_audit_log_export.new_auditlogexport: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_audit_log_export.new_auditlogexport: Destruction complete after 0s

```
