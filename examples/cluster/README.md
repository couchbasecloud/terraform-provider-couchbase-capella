# Capella Cluster Example

This example shows how to create and manage Clusters in Capella.

This creates a new cluster in the selected Capella project. It uses the organization ID and projectId to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.
1. Create a new cluster with the specified configuration.


### View the plan for the resources that Terraform will create

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

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

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
      + name             = "New Terraform Cluster"
      + organization_id  = "bdb8662c-7157-46ea-956f-ed86f4c75211"
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

  # capella_project.existing_project will be created
  + resource "capella_project" "existing_project" {
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
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_project = {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }
  + new_cluster      = {
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
      + name             = "New Terraform Cluster"
      + organization_id  = "bdb8662c-7157-46ea-956f-ed86f4c75211"
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

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Cluster in Capella

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

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

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
      + name             = "New Terraform Cluster"
      + organization_id  = "bdb8662c-7157-46ea-956f-ed86f4c75211"
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

  # capella_project.existing_project will be created
  + resource "capella_project" "existing_project" {
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
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_project = {
      + audit           = (known after apply)
      + description     = "A Capella Project that will host many Capella clusters."
      + etag            = (known after apply)
      + id              = (known after apply)
      + if_match        = null
      + name            = "terraform-couchbasecapella-project"
      + organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
    }
  + new_cluster      = {
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
      + name             = "New Terraform Cluster"
      + organization_id  = "bdb8662c-7157-46ea-956f-ed86f4c75211"
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

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_project.existing_project: Creating...
capella_project.existing_project: Creation complete after 1s [id=3dac8dfb-69dd-4145-a852-11da0766d8a9]
capella_cluster.new_cluster: Creating...
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
capella_cluster.new_cluster: Still creating... [3m0s elapsed]
capella_cluster.new_cluster: Creation complete after 3m3s [id=da519c25-79b2-477e-aa77-61e2ee7a677a]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

existing_project = {
  "audit" = {
    "created_at" = "2023-09-19 22:50:45.476879515 +0000 UTC"
    "created_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
    "modified_at" = "2023-09-19 22:50:45.476901411 +0000 UTC"
    "modified_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
    "version" = 1
  }
  "description" = "A Capella Project that will host many Capella clusters."
  "etag" = "Version: 1"
  "id" = "3dac8dfb-69dd-4145-a852-11da0766d8a9"
  "if_match" = tostring(null)
  "name" = "terraform-couchbasecapella-project"
  "organization_id" = "bdb8662c-7157-46ea-956f-ed86f4c75211"
}
new_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-09-19 22:50:45.937512888 +0000 UTC"
    "created_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
    "modified_at" = "2023-09-19 22:53:48.367132006 +0000 UTC"
    "modified_by" = "7eEPh2Jdzb3fwRNesFoONpyAkq5nhAfK"
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
  "id" = "da519c25-79b2-477e-aa77-61e2ee7a677a"
  "if_match" = tostring(null)
  "name" = "New Terraform Cluster"
  "organization_id" = "bdb8662c-7157-46ea-956f-ed86f4c75211"
  "project_id" = "3dac8dfb-69dd-4145-a852-11da0766d8a9"
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
```

