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
16. Store the existing event in an output variable.
17. Store the existing events in an output variable.
18. Store the existing project event in an output variable.
19. Store the existing project events in an output variable.

## Pre-Requisites:

Make sure you have followed the pre-requisite steps from the parent readme that builds the binary and stores the path in a .terraformrc file.

### Variables file

- Copy the `terraform.template.tfvars` file to `terraform.tfvars` file in the same directory
- Create 1 V4 API key in your organization using the Capella UI.
- Replace all placeholders with actual values. Use the above created API Key secret as the value for `auth_token`

## Execution

Command: `terraform plan`

Sample Output:
```
$  terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-capella
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_events.existing_events: Reading...
data.couchbase-capella_organization.existing_organization: Reading...
data.couchbase-capella_organization.existing_organization: Read complete after 0s [name=prod]
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
      + description        = ""
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
      + app_service_id     = (known after apply)
      + audit              = (known after apply)
      + availability       = {
          + type = "multi"
        }
      + cloud_provider     = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type = (known after apply)
      + couchbase_server   = (known after apply)
      + current_state      = (known after apply)
      + description        = "My first test cluster for multiple services."
      + etag               = (known after apply)
      + id                 = (known after apply)
      + name               = "My First Terraform Cluster"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = (known after apply)
      + service_groups     = [
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
      + support            = {
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

  # couchbase-capella_project.new_project will be created
  + resource "couchbase-capella_project" "new_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
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

Plan: 13 to add, 0 to change, 0 to destroy.

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
      + app_service_id     = (known after apply)
      + audit              = (known after apply)
      + availability       = {
          + type = "multi"
        }
      + cloud_provider     = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type = (known after apply)
      + couchbase_server   = (known after apply)
      + current_state      = (known after apply)
      + description        = "My first test cluster for multiple services."
      + etag               = (known after apply)
      + id                 = (known after apply)
      + if_match           = null
      + name               = "My First Terraform Cluster"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = (known after apply)
      + service_groups     = [
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
      + support            = {
          + plan     = "enterprise"
          + timezone = "PT"
        }
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
      + timestamp        = "2024-08-08 20:56:48.543271144 +0000 UTC"
      + user_email       = null
      + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + user_name        = "Kevin"
    }
  + existing_events         = {
      + cluster_ids     = null
      + cursor          = {
          + hrefs = {
              + first    = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=10"
              + last     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2787&perPage=10"
              + next     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=10"
              + previous = ""
            }
          + pages = {
              + last        = 2787
              + next        = 2
              + page        = 1
              + per_page    = 10
              + previous    = 0
              + total_items = 27865
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
              + timestamp        = "2024-08-08 20:56:48.543271144 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 21:08:24.801168066 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 21:28:58.80962237 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_requested"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "obligin"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_requested"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-08-08 21:29:46.255942529 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_completed"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "obligin"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_completed"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = null
              + severity         = "info"
              + source           = "cp-jobs"
              + summary          = null
              + timestamp        = "2024-08-08 21:30:09.67600442 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "logged_out"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_out"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-08-08 21:33:12.906102913 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 21:45:41.055833681 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_requested"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "test"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_requested"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-08-08 21:46:25.085921176 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_completed"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "test"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_completed"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = null
              + severity         = "info"
              + source           = "cp-jobs"
              + summary          = null
              + timestamp        = "2024-08-08 21:48:33.478275112 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 22:09:48.399077024 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
          + created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          + modified_at = "2024-03-26 05:23:05.348085173 +0000 UTC"
          + modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
          + version     = 35
        }
      + description     = ""
      + name            = "prod"
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

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```

Command: `terraform apply`

Sample Output:
```
$  terraform apply   
terraform apply --var-file=terraform.template.tfvars
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-capella
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_organization.existing_organization: Reading...
data.couchbase-capella_events.existing_events: Reading...
data.couchbase-capella_organization.existing_organization: Read complete after 1s [name=prod]
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
      + description        = ""
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
      + app_service_id     = (known after apply)
      + audit              = (known after apply)
      + availability       = {
          + type = "multi"
        }
      + cloud_provider     = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type = (known after apply)
      + couchbase_server   = (known after apply)
      + current_state      = (known after apply)
      + description        = "My first test cluster for multiple services."
      + etag               = (known after apply)
      + id                 = (known after apply)
      + name               = "My First Terraform Cluster"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = (known after apply)
      + service_groups     = [
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
      + support            = {
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

  # couchbase-capella_project.new_project will be created
  + resource "couchbase-capella_project" "new_project" {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + name            = "My First Terraform Project"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
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

Plan: 13 to add, 0 to change, 0 to destroy.

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
      + app_service_id     = (known after apply)
      + audit              = (known after apply)
      + availability       = {
          + type = "multi"
        }
      + cloud_provider     = {
          + cidr   = "10.250.250.0/23"
          + region = "us-east-1"
          + type   = "aws"
        }
      + configuration_type = (known after apply)
      + couchbase_server   = (known after apply)
      + current_state      = (known after apply)
      + description        = "My first test cluster for multiple services."
      + etag               = (known after apply)
      + id                 = (known after apply)
      + if_match           = null
      + name               = "My First Terraform Cluster"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = (known after apply)
      + service_groups     = [
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
      + support            = {
          + plan     = "enterprise"
          + timezone = "PT"
        }
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
      + timestamp        = "2024-08-08 20:56:48.543271144 +0000 UTC"
      + user_email       = null
      + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + user_name        = "Kevin"
    }
  + existing_events         = {
      + cluster_ids     = null
      + cursor          = {
          + hrefs = {
              + first    = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=10"
              + last     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2794&perPage=10"
              + next     = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=10"
              + previous = ""
            }
          + pages = {
              + last        = 2794
              + next        = 2
              + page        = 1
              + per_page    = 10
              + previous    = 0
              + total_items = 27937
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
              + timestamp        = "2024-08-08 20:56:48.543271144 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 21:08:24.801168066 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 21:28:58.80962237 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_requested"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "obliging"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_requested"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-08-08 21:29:46.255942529 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_completed"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "obliging"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_completed"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = null
              + severity         = "info"
              + source           = "cp-jobs"
              + summary          = null
              + timestamp        = "2024-08-08 21:30:09.67600442 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "logged_out"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = null
              + cluster_name     = null
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "logged_out"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = null
              + project_name     = null
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-08-08 21:33:12.906102913 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 21:45:41.055833681 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_requested"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "test"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_requested"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + severity         = "info"
              + source           = "cp-api"
              + summary          = null
              + timestamp        = "2024-08-08 21:46:25.085921176 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
            },
          + {
              + alert_key        = "cluster_deletion_completed"
              + app_service_id   = null
              + app_service_name = null
              + cluster_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name     = "test"
              + id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + image_url        = null
              + incident_ids     = []
              + key              = "cluster_deletion_completed"
              + kv               = "null"
              + occurrence_count = null
              + project_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_name     = "  !!!!!!!-Shared-Project-!!!!!!!"
              + request_id       = "ffffffff-aaaa-1414-eeee-000000000000"
              + session_id       = null
              + severity         = "info"
              + source           = "cp-jobs"
              + summary          = null
              + timestamp        = "2024-08-08 21:48:33.478275112 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
              + timestamp        = "2024-08-08 22:09:48.399077024 +0000 UTC"
              + user_email       = null
              + user_id          = "ffffffff-aaaa-1414-eeee-000000000000"
              + user_name        = "Kevin"
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
          + created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          + modified_at = "2024-03-26 05:23:05.348085173 +0000 UTC"
          + modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
          + version     = 35
        }
      + description     = ""
      + name            = "prod"
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
couchbase-capella_project.new_project: Creation complete after 1s [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_project_events.existing_project_events: Reading...
couchbase-capella_user.new_user: Creating...
couchbase-capella_apikey.new_apikey: Creating...
couchbase-capella_cluster.new_cluster: Creating...
data.couchbase-capella_project_events.existing_project_events: Read complete after 1s
data.couchbase-capella_project_event.existing_project_event: Reading...
couchbase-capella_apikey.new_apikey: Creation complete after 1s [id=U68cAv5ap82Pf7Jj0wQoMs5Pgj30OWif]
data.couchbase-capella_project_event.existing_project_event: Read complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_user.new_user: Creation complete after 1s [id=ffffffff-aaaa-1414-eeee-000000000000]
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
couchbase-capella_cluster.new_cluster: Still creating... [2m20s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [2m30s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [2m40s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [2m50s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [3m0s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [3m10s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [3m20s elapsed]
couchbase-capella_cluster.new_cluster: Creation complete after 3m28s [id=ffffffff-aaaa-1414-eeee-000000000000]
data.couchbase-capella_certificate.existing_certificate: Reading...
data.couchbase-capella_audit_log_event_ids.event_list: Reading...
couchbase-capella_allowlist.new_allowlist: Creating...
couchbase-capella_sample_bucket.new_sample_bucket: Creating...
couchbase-capella_bucket.new_bucket: Creating...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Creating...
couchbase-capella_database_credential.new_database_credential: Creating...
data.couchbase-capella_certificate.existing_certificate: Read complete after 1s
data.couchbase-capella_audit_log_event_ids.event_list: Read complete after 1s
couchbase-capella_audit_log_settings.new_auditlogsettings: Creating...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Creation complete after 1s
couchbase-capella_bucket.new_bucket: Creation complete after 1s [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
couchbase-capella_scope.new_scope: Creating...
couchbase-capella_database_credential.new_database_credential: Creation complete after 1s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_scope.new_scope: Creation complete after 1s
couchbase-capella_collection.new_collection: Creating...
couchbase-capella_collection.new_collection: Creation complete after 1s
couchbase-capella_audit_log_settings.new_auditlogsettings: Creation complete after 5s
couchbase-capella_allowlist.new_allowlist: Creation complete after 8s [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_sample_bucket.new_sample_bucket: Still creating... [10s elapsed]
couchbase-capella_sample_bucket.new_sample_bucket: Still creating... [20s elapsed]
couchbase-capella_sample_bucket.new_sample_bucket: Creation complete after 23s [id=Z2FtZXNpbS1zYW1wbGU=]
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
couchbase-capella_app_service.new_app_service: Still creating... [5m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [5m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [5m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [5m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [5m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [5m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [6m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [6m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [6m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [6m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [6m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [6m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [7m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [7m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [7m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [7m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [7m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [7m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [8m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [8m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [8m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [8m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [8m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [8m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [9m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [9m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [9m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [9m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [9m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [9m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [10m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [10m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [10m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [10m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [10m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [10m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [11m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [11m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [11m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [11m30s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [11m40s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [11m50s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [12m0s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [12m10s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [12m20s elapsed]
couchbase-capella_app_service.new_app_service: Still creating... [12m30s elapsed]
couchbase-capella_app_service.new_app_service: Creation complete after 12m31s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 13 added, 0 changed, 0 destroyed.

Outputs:

apikey = <sensitive>
app_service = {
  "audit" = {
    "created_at" = "2024-08-09 18:42:34.1568404 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000""
    "modified_at" = "2024-08-09 18:55:02.497681001 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000""
    "version" = 8
  }
  "cloud_provider" = "AWS"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "compute" = {
    "cpu" = 2
    "ram" = 4
  }
  "current_state" = "healthy"
  "description" = "My first test app service."
  "etag" = "Version: 8"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "if_match" = tostring(null)
  "name" = "new-terraform-app-service"
  "nodes" = 2
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "version" = "3.1.8-1.0.0"
}
bucket = "new_terraform_bucket"
certificate = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
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
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2024-08-09 18:38:43.700996649 +0000 UTC"
    "created_by" = "LFYpRl0AlDX9sfDXO36ZffBrINTjnS9r"
    "modified_at" = "2024-08-09 18:42:10.299873675 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000""
    "version" = 5
  }
  "availability" = {
    "type" = "multi"
  }
  "cloud_provider" = {
    "cidr" = "10.250.250.0/23"
    "region" = "us-east-1"
    "type" = "aws"
  }
  "configuration_type" = "multiNode"
  "couchbase_server" = {
    "version" = "7.6"
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
  "bucket_id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "new_terraform_collection"
  "max_ttl" = 200
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "new_terraform_scope"
}
database_credential = <sensitive>
existing_event = {
  "alert_key" = "logged_in"
  "app_service_id" = tostring(null)
  "app_service_name" = tostring(null)
  "cluster_id" = tostring(null)
  "cluster_name" = tostring(null)
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "image_url" = tostring(null)
  "incident_ids" = toset([])
  "key" = "logged_in"
  "kv" = "null"
  "occurrence_count" = tonumber(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = tostring(null)
  "project_name" = tostring(null)
  "request_id" = tostring(null)
  "session_id" = tostring(null)
  "severity" = "info"
  "source" = "cp-ns"
  "summary" = tostring(null)
  "timestamp" = "2024-08-08 20:56:48.543271144 +0000 UTC"
  "user_email" = tostring(null)
  "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "user_name" = "Kevin"
}
existing_events = {
  "cluster_ids" = toset(null) /* of string */
  "cursor" = {
    "hrefs" = {
      "first" = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=10"
      "last" = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2794&perPage=10"
      "next" = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/events?page=2&perPage=10"
      "previous" = ""
    }
    "pages" = {
      "last" = 2794
      "next" = 2
      "page" = 1
      "per_page" = 10
      "previous" = 0
      "total_items" = 27937
    }
  }
  "data" = tolist([
    {
      "alert_key" = "logged_in"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "logged_in"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = tostring(null)
      "project_name" = tostring(null)
      "request_id" = tostring(null)
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-ns"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 20:56:48.543271144 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "logged_in"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "logged_in"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = tostring(null)
      "project_name" = tostring(null)
      "request_id" = tostring(null)
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-ns"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:08:24.801168066 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "logged_in"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "logged_in"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = tostring(null)
      "project_name" = tostring(null)
      "request_id" = tostring(null)
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-ns"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:28:58.80962237 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "cluster_deletion_requested"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "obligingmmmusharafhussain"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "cluster_deletion_requested"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "  !!!!!!!-Shared-Project-!!!!!!!"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "severity" = "info"
      "source" = "cp-api"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:29:46.255942529 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "cluster_deletion_completed"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "obligingmmmusharafhussain"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "cluster_deletion_completed"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "  !!!!!!!-Shared-Project-!!!!!!!"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-jobs"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:30:09.67600442 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "logged_out"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "logged_out"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = tostring(null)
      "project_name" = tostring(null)
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "severity" = "info"
      "source" = "cp-api"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:33:12.906102913 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "logged_in"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "logged_in"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = tostring(null)
      "project_name" = tostring(null)
      "request_id" = tostring(null)
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-ns"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:45:41.055833681 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "cluster_deletion_requested"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "test"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "cluster_deletion_requested"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "  !!!!!!!-Shared-Project-!!!!!!!"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "severity" = "info"
      "source" = "cp-api"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:46:25.085921176 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "cluster_deletion_completed"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "test"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "cluster_deletion_completed"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "  !!!!!!!-Shared-Project-!!!!!!!"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-jobs"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 21:48:33.478275112 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
    {
      "alert_key" = "logged_in"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "logged_in"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = tostring(null)
      "project_name" = tostring(null)
      "request_id" = tostring(null)
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-ns"
      "summary" = tostring(null)
      "timestamp" = "2024-08-08 22:09:48.399077024 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = "Kevin"
    },
  ])
  "from" = tostring(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "page" = tonumber(null)
  "per_page" = tonumber(null)
  "project_ids" = toset(null) /* of string */
  "severity_levels" = toset(null) /* of string */
  "sort_by" = tostring(null)
  "sort_direction" = tostring(null)
  "tags" = toset(null) /* of string */
  "to" = tostring(null)
  "user_ids" = toset(null) /* of string */
}
existing_project_event = {
  "alert_key" = "permissions_assigned"
  "app_service_id" = tostring(null)
  "app_service_name" = tostring(null)
  "cluster_id" = tostring(null)
  "cluster_name" = tostring(null)
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "image_url" = tostring(null)
  "incident_ids" = toset([])
  "key" = "permissions_assigned"
  "kv" = "{\"resourceType\":\"project\",\"roles\":\"projectOwner\",\"userAffected\":\"\"}"
  "occurrence_count" = tonumber(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_name" = "My First Terraform Project"
  "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "session_id" = tostring(null)
  "severity" = "info"
  "source" = "cp-open-api"
  "summary" = " was granted projectOwner on project \"prod\" by ."
  "timestamp" = "2024-08-09 18:38:42.795277617 +0000 UTC"
  "user_email" = tostring(null)
  "user_id" = "LFYpRl0AlDX9sfDXO36ZffBrINTjnS9r"
  "user_name" = tostring(null)
}
existing_project_events = {
  "cluster_ids" = toset(null) /* of string */
  "cursor" = {
    "hrefs" = {
      "first" = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=10"
      "last" = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-aaaa-1414-eeee-000000000000/events?page=1&perPage=10"
      "next" = "http://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-aaaa-1414-eeee-000000000000/events?page=0&perPage=10"
      "previous" = ""
    }
    "pages" = {
      "last" = 1
      "next" = 0
      "page" = 1
      "per_page" = 10
      "previous" = 0
      "total_items" = 2
    }
  }
  "data" = tolist([
    {
      "alert_key" = "permissions_assigned"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "permissions_assigned"
      "kv" = "{\"resourceType\":\"project\",\"roles\":\"projectOwner\",\"userAffected\":\"\"}"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "My First Terraform Project"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-open-api"
      "summary" = " was granted projectOwner on project \"prod\" by ."
      "timestamp" = "2024-08-09 18:38:42.795277617 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = tostring(null)
    },
    {
      "alert_key" = "project_created"
      "app_service_id" = tostring(null)
      "app_service_name" = tostring(null)
      "cluster_id" = tostring(null)
      "cluster_name" = tostring(null)
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "image_url" = tostring(null)
      "incident_ids" = toset([])
      "key" = "project_created"
      "kv" = "null"
      "occurrence_count" = tonumber(null)
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_name" = "My First Terraform Project"
      "request_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "session_id" = tostring(null)
      "severity" = "info"
      "source" = "cp-open-api"
      "summary" = tostring(null)
      "timestamp" = "2024-08-09 18:38:42.785906591 +0000 UTC"
      "user_email" = tostring(null)
      "user_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "user_name" = tostring(null)
    },
  ])
  "from" = tostring(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "page" = tonumber(null)
  "per_page" = tonumber(null)
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "severity_levels" = toset(null) /* of string */
  "sort_by" = tostring(null)
  "sort_direction" = tostring(null)
  "tags" = toset(null) /* of string */
  "to" = tostring(null)
  "user_ids" = toset(null) /* of string */
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
    28732,
    28733,
    28734,
    28735,
    28736,
    28737,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
organization = {
  "audit" = {
    "created_at" = "2021-12-03 16:14:45.105347711 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2024-03-26 05:23:05.348085173 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 35
  }
  "description" = ""
  "name" = "prod"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "preferences" = {
    "session_duration" = 7200
  }
}
project = "My First Terraform Project"
sample_bucket = "gamesim-sample"
scope = {
  "bucket_id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections" = toset([])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "new_terraform_scope"
}
user = {
  "audit" = {
    "created_at" = "2023-11-01 10:18:01.680701109 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2023-11-01 10:18:01.680701109 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 1
  }
  "email" = "johndoe@couchbase.com"
  "enable_notifications" = false
  "expires_at" = "2024-01-30T10:18:01.680701749Z"
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
```
