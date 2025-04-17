# Capella Free Tier Cluster examples
This example demonstrates how to create and manage a free-tier cluster in Capella.

It provisions a new free-tier cluster within the selected Capella project and tenant while managing the cluster.

To run this example, configure your Couchbase Capella provider as described in the README file located in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new free-ter in Capella as stated in the `create_free_tier_cluster.tf` file.
2. LIST: List the free-tier cluster in the given project
3. UPDATE: Update the free-tier cluster configuration using Terraform.
4. IMPORT: Import a free-tier cluster that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created free-tier cluster from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.


## CREATE
### View the plan for the resources that Terraform will create


Command: `terraform plan`

Sample Output:
```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
Terraform will perform the following actions:
  # couchbase-capella_free_tier_cluster.new_free_tier_cluster will be created
  + resource "couchbase-capella_free_tier_cluster" "new_free_tier_cluster" {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "10.1.0.0/24"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "New free tier test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }
Plan: 1 to add, 0 to change, 0 to destroy.
Changes to Outputs:
  + free_tier_cluster_id  = (known after apply)
  + new_free_tier_cluster = {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "10.1.0.0/24"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "New free tier test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }
──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```
### Apply the Plan, in order to create a new free-tier cluster

Command: `terraform apply`

Sample Output:
```
$ terraform apply

╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
Terraform will perform the following actions:
  # couchbase-capella_free_tier_cluster.new_free_tier_cluster will be created
  + resource "couchbase-capella_free_tier_cluster" "new_free_tier_cluster" {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "10.1.0.0/24"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "New free tier test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }
Plan: 1 to add, 0 to change, 0 to destroy.
Changes to Outputs:
  + free_tier_cluster_id  = (known after apply)
  + new_free_tier_cluster = {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "10.1.0.0/24"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "New free tier test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }
Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.
  Enter a value: yes
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Creating...
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [10s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [20s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [30s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [40s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [50s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [1m0s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [1m10s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [1m20s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [1m30s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [1m40s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [1m50s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still creating... [2m0s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Creation complete after 2m2s [id=ffffffff-aaaa-1414-eeee-000000000000]
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
Outputs:
free_tier_cluster_id = "ffffffff-aaaa-1414-eeee-000000000000"
new_free_tier_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2025-03-04 15:22:59.96921537 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2025-03-04 15:24:30.453257937 +0000 UTC"
    "modified_by" = "apikey-apikey-ffffffff-aaaa-1414-eeee-000000000000""
    "version" = 5
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "10.1.0.0/24"
    "region" = "us-east-2"
    "type" = "aws"
  }
  "cmek_id" = tostring(null)
  "connection_string" = "ffffffff-aaaa-1414-eeee-000000000000"
  "couchbase_server" = {
    "version" = "7.6"
  }
  "current_state" = "healthy"
  "description" = "New free tier test cluster for multiple services"
  "enable_private_dns_resolution" = false
  "etag" = "Version: 5"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "New free tier cluster"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "service_groups" = toset([
    {
      "node" = {
        "compute" = {
          "cpu" = 2
          "ram" = 8
        }
        "disk" = {
          "autoexpansion" = tobool(null)
          "iops" = 3000
          "storage" = 10
          "type" = "gp3"
        }
      }
      "num_of_nodes" = 1
      "services" = toset([
        "data",
        "index",
        "query",
        "search",
      ])
    },
  ])
  "support" = {
    "plan" = "free"
    "timezone" = "PT"
  }
}

```

## LIST
## List the free-tier cluster in the given project

Command: `terraform plan`

Sample Output:
```
$terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_free_tier_clusters.existing_free_tier_clusters: Reading...
data.couchbase-capella_free_tier_clusters.existing_free_tier_clusters: Read complete after 1s

Changes to Outputs:
  + free_tier_clusters_list = {
      + data            = [
          + {
              + app_service_id                = null
              + audit                         = {
                  + created_at  = "2025-04-02 06:18:34.966637849 +0000 UTC"
                  + created_by  = "ea36341e-c899-4cfe-b3bc-3d8e60f14654"
                  + modified_at = "2025-04-02 08:54:05.539289368 +0000 UTC"
                  + modified_by = "aaaaa"
                  + version     = 11
                }
              + availability                  = {
                  + type = "single"
                }
              + cloud_provider                = {
                  + cidr   = "10.0.0.0/24"
                  + region = "us-east-2"
                  + type   = "aws"
                }
              + connection_string             = "ffffffff-aaaa-1414-eeee-000000000000"
              + couchbase_server              = {
                  + version = "7.6.5"
                }
              + current_state                 = "healthy"
              + description                   = ""
              + enable_private_dns_resolution = false
              + id                            = "ffffffff-aaaa-1414-eeee-000000000000"
              + name                          = "daringschahramdustdar"
              + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
              + service_groups                = [
                  + {
                      + node         = {
                          + compute = {
                              + cpu = 2
                              + ram = 8
                            }
                          + disk    = {
                              + autoexpansion = null
                              + iops          = 3000
                              + storage       = 10
                              + type          = "gp3"
                            }
                        }
                      + num_of_nodes = 1
                      + services     = [
                          + "search",
                          + "index",
                          + "data",
                          + "query",
                        ]
                    },
                ]
              + support                       = {
                  + plan     = "free"
                  + timezone = "PT"
                }
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

Now apply the plan to list the free-tier cluster in the given project

Command: `terraform apply`

Sample Output
```
$terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_free_tier_clusters.existing_free_tier_clusters: Reading...
data.couchbase-capella_free_tier_clusters.existing_free_tier_clusters: Read complete after 1s

Changes to Outputs:
  + free_tier_clusters_list = {
      + data            = [
          + {
              + app_service_id                = null
              + audit                         = {
                  + created_at  = "2025-04-02 06:18:34.966637849 +0000 UTC"
                  + created_by  = "ea36341e-c899-4cfe-b3bc-3d8e60f14654"
                  + modified_at = "2025-04-02 08:54:05.539289368 +0000 UTC"
                  + modified_by = "CP-INTERNAL-API-SYSTEM-USER"
                  + version     = 11
                }
              + availability                  = {
                  + type = "single"
                }
              + cloud_provider                = {
                  + cidr   = "10.0.0.0/24"
                  + region = "us-east-2"
                  + type   = "aws"
                }
              + connection_string             = "ffffffff-aaaa-1414-eeee-000000000000"
              + couchbase_server              = {
                  + version = "7.6.5"
                }
              + current_state                 = "healthy"
              + description                   = ""
              + enable_private_dns_resolution = false
              + id                            = "ffffffff-aaaa-1414-eeee-000000000000"
              + name                          = "daringschahramdustdar"
              + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
              + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
              + service_groups                = [
                  + {
                      + node         = {
                          + compute = {
                              + cpu = 2
                              + ram = 8
                            }
                          + disk    = {
                              + autoexpansion = null
                              + iops          = 3000
                              + storage       = 10
                              + type          = "gp3"
                            }
                        }
                      + num_of_nodes = 1
                      + services     = [
                          + "search",
                          + "index",
                          + "data",
                          + "query",
                        ]
                    },
                ]
              + support                       = {
                  + plan     = "free"
                  + timezone = "PT"
                }
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

free_tier_clusters_list = {
  "data" = tolist([
    {
      "app_service_id" = tostring(null)
      "audit" = {
        "created_at" = "2025-04-02 06:18:34.966637849 +0000 UTC"
        "created_by" = "ea36341e-c899-4cfe-b3bc-3d8e60f14654"
        "modified_at" = "2025-04-02 08:54:05.539289368 +0000 UTC"
        "modified_by" = "CP-INTERNAL-API-SYSTEM-USER"
        "version" = 11
      }
      "availability" = {
        "type" = "single"
      }
      "cloud_provider" = {
        "cidr" = "10.0.0.0/24"
        "region" = "us-east-2"
        "type" = "aws"
      }
      "connection_string" = "ffffffff-aaaa-1414-eeee-000000000000"
      "couchbase_server" = {
        "version" = "7.6.5"
      }
      "current_state" = "healthy"
      "description" = ""
      "enable_private_dns_resolution" = false
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "name" = "daringschahramdustdar"
      "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "service_groups" = tolist([
        {
          "node" = {
            "compute" = {
              "cpu" = 2
              "ram" = 8
            }
            "disk" = {
              "autoexpansion" = tobool(null)
              "iops" = 3000
              "storage" = 10
              "type" = "gp3"
            }
          }
          "num_of_nodes" = 1
          "services" = tolist([
            "search",
            "index",
            "data",
            "query",
          ])
        },
      ])
      "support" = {
        "plan" = "free"
        "timezone" = "PT"
      }
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```


## UPDATE
### Update the resource `new_cluster_free_tier` in the Terraform State file

Free tier supports updating only name and description, in the below examples the name and description have been modified
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Refreshing state... [id=8877dd0f-40bc-4cf8-8c03-c146f7b50e55]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_free_tier_cluster.new_free_tier_cluster will be updated in-place
  ~ resource "couchbase-capella_free_tier_cluster" "new_free_tier_cluster" {
      + app_service_id                = (known after apply)
      ~ audit                         = {
          ~ created_at  = "2025-03-05 20:29:53.509880483 +0000 UTC" -> (known after apply)
          ~ created_by  = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
          ~ modified_at = "2025-03-05 20:31:17.549230941 +0000 UTC" -> (known after apply)
          ~ modified_by = "apikey-apikey-ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
          ~ version     = 5 -> (known after apply)
        } -> (known after apply)
      ~ availability                  = {
          ~ type = "single" -> (known after apply)
        } -> (known after apply)
      + cmek_id                       = (known after apply)
      ~ connection_string             = "cb.0i0gtdcztbnhomdf.sandbox.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ description                   = "New free tier test cluster for multiple services" -> "New free tier test cluster for multiple services modified"
      ~ enable_private_dns_resolution = false -> (known after apply)
      ~ etag                          = "Version: 5" -> (known after apply)
        id                            = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ name                          = "New free tier cluster" -> "New free tier cluster modified"
      ~ service_groups                = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 2 -> null
                      - ram = 8 -> null
                    } -> null
                  - disk    = {
                      - iops    = 3000 -> null
                      - storage = 10 -> null
                      - type    = "gp3" -> null
                    } -> null
                } -> null
              - num_of_nodes = 1 -> null
              - services     = [
                  - "data",
                  - "index",
                  - "query",
                  - "search",
                ] -> null
            },
        ] -> (known after apply)
      ~ support                       = {
          ~ plan     = "free" -> (known after apply)
          ~ timezone = "PT" -> (known after apply)
        } -> (known after apply)
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_free_tier_cluster = {
      + app_service_id                = (known after apply)
      ~ audit                         = {
          - created_at  = "2025-03-05 20:29:53.509880483 +0000 UTC"
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          - modified_at = "2025-03-05 20:31:17.549230941 +0000 UTC"
          - modified_by = "apikey-apikey-ffffffff-aaaa-1414-eeee-000000000000"
          - version     = 5
        } -> (known after apply)
      ~ availability                  = {
          - type = "single"
        } -> (known after apply)
      + cmek_id                       = (known after apply)
      ~ connection_string             = "cb.0i0gtdcztbnhomdf.sandbox.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ description                   = "New free tier test cluster for multiple services" -> "New free tier test cluster for multiple services modified"
      ~ enable_private_dns_resolution = false -> (known after apply)
      ~ etag                          = "Version: 5" -> (known after apply)
        id                            = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ name                          = "New free tier cluster" -> "New free tier cluster modified"
      ~ service_groups                = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 2
                      - ram = 8
                    }
                  - disk    = {
                      - autoexpansion = null
                      - iops          = 3000
                      - storage       = 10
                      - type          = "gp3"
                    }
                }
              - num_of_nodes = 1
              - services     = [
                  - "data",
                  - "index",
                  - "query",
                  - "search",
                ]
            },
        ] -> (known after apply)
      ~ support                       = {
          - plan     = "free"
          - timezone = "PT"
        } -> (known after apply)
        # (4 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_free_tier_cluster.new_free_tier_cluster: Modifying... [id=8877dd0f-40bc-4cf8-8c03-c146f7b50e55]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Modifications complete after 1s [id=8877dd0f-40bc-4cf8-8c03-c146f7b50e55]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

free_tier_cluster_id = "ffffffff-aaaa-1414-eeee-000000000000"
new_free_tier_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2025-03-05 20:29:53.509880483 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2025-03-05 20:32:37.921970375 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 6
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "10.1.0.0/24"
    "region" = "us-east-2"
    "type" = "aws"
  }
  "cmek_id" = tostring(null)
  "connection_string" = "cb.0i0gtdcztbnhomdf.sandbox.nonprod-project-avengers.com"
  "couchbase_server" = {
    "version" = "7.6"
  }
  "current_state" = "healthy"
  "description" = "New free tier test cluster for multiple services modified"
  "enable_private_dns_resolution" = false
  "etag" = "Version: 6"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "name" = "New free tier cluster modified"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "service_groups" = toset([
    {
      "node" = {
        "compute" = {
          "cpu" = 2
          "ram" = 8
        }
        "disk" = {
          "autoexpansion" = tobool(null)
          "iops" = 3000
          "storage" = 10
          "type" = "gp3"
        }
      }
      "num_of_nodes" = 1
      "services" = toset([
        "data",
        "index",
        "query",
        "search",
      ])
    },
  ])
  "support" = {
    "plan" = "free"
    "timezone" = "PT"
  }
}

```

## IMPORT
### Remove the resource `new_cluster_free_tier` from the Terraform State file

Command: `terraform state rm couchbase-capella_cluster_free_tier.new_cluster_free_tier`

Sample Output
```
$ terraform state rm couchbase-capella_cluster_free_tier.new_cluster_free_tier
Removed couchbase-capella_cluster_free_tier.new_cluster_free_tier
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command:
`terraform import couchbase-capella_cluster_free_tier.new_cluster_free_tier id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Example:
```
$ terraform import couchbase-capella_cluster_free_tier.new_cluster_free_tier id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Import prepared!
  Prepared couchbase-capella_cluster_free_tier for import
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Refreshing state... [id=id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

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
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_free_tier_cluster.new_free_tier_cluster will be destroyed
  - resource "couchbase-capella_free_tier_cluster" "new_free_tier_cluster" {
      - audit                         = {
          - created_at  = "2025-03-05 20:12:16.276387369 +0000 UTC" -> null
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
          - modified_at = "2025-03-05 20:13:41.108694607 +0000 UTC" -> null
          - modified_by = "apikey-apikey-ffffffff-aaaa-1414-eeee-000000000000" -> null
          - version     = 5 -> null
        } -> null
      - availability                  = {
          - type = "single" -> null
        } -> null
      - cloud_provider                = {
          - cidr   = "10.1.0.0/24" -> null
          - region = "us-east-2" -> null
          - type   = "aws" -> null
        } -> null
      - connection_string             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - couchbase_server              = {
          - version = "7.6" -> null
        } -> null
      - current_state                 = "healthy" -> null
      - description                   = "New free tier test cluster for multiple services" -> null
      - enable_private_dns_resolution = false -> null
      - etag                          = "Version: 5" -> null
      - id                            = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name                          = "New free tier cluster" -> null
      - organization_id               = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id                    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - service_groups                = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 2 -> null
                      - ram = 8 -> null
                    } -> null
                  - disk    = {
                      - iops    = 3000 -> null
                      - storage = 10 -> null
                      - type    = "gp3" -> null
                    } -> null
                } -> null
              - num_of_nodes = 1 -> null
              - services     = [
                  - "data",
                  - "index",
                  - "query",
                  - "search",
                ] -> null
            },
        ] -> null
      - support                       = {
          - plan     = "free" -> null
          - timezone = "PT" -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - free_tier_cluster_id  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
  - new_free_tier_cluster = {
      - app_service_id                = null
      - audit                         = {
          - created_at  = "2025-03-05 20:12:16.276387369 +0000 UTC"
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          - modified_at = "2025-03-05 20:13:41.108694607 +0000 UTC"
          - modified_by = "apikey-apikey-ffffffff-aaaa-1414-eeee-000000000000"
          - version     = 5
        }
      - availability                  = {
          - type = "single"
        }
      - cloud_provider                = {
          - cidr   = "10.1.0.0/24"
          - region = "us-east-2"
          - type   = "aws"
        }
      - cmek_id                       = null
      - connection_string             = "ffffffff-aaaa-1414-eeee-000000000000"
      - couchbase_server              = {
          - version = "7.6"
        }
      - current_state                 = "healthy"
      - description                   = "New free tier test cluster for multiple services"
      - enable_private_dns_resolution = false
      - etag                          = "Version: 5"
      - id                            = "ffffffff-aaaa-1414-eeee-000000000000"
      - name                          = "New free tier cluster"
      - organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      - service_groups                = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 2
                      - ram = 8
                    }
                  - disk    = {
                      - autoexpansion = null
                      - iops          = 3000
                      - storage       = 10
                      - type          = "gp3"
                    }
                }
              - num_of_nodes = 1
              - services     = [
                  - "data",
                  - "index",
                  - "query",
                  - "search",
                ]
            },
        ]
      - support                       = {
          - plan     = "free"
          - timezone = "PT"
        }
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_free_tier_cluster.new_free_tier_cluster: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 10s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 20s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 30s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 40s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 50s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m0s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m10s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m20s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m30s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m40s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 1m50s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-000000000000, 2m0s elapsed]
couchbase-capella_free_tier_cluster.new_free_tier_cluster: Destruction complete after 2m2s

Destroy complete! Resources: 1 destroyed.

```