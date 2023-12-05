# Capella Cluster Example

This example shows how to create and manage Cluster in Capella.

This creates a new cluster in the selected Capella cluster and lists existing clusters in the cluster. It uses the project ID to create and list clusters.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new cluster in Capella as stated in the `create_cluster.tf` file.
2. UPDATE: Update the cluster configuration using Terraform.
3. LIST: List existing clusters in Capella as stated in the `list_clusters.tf` file.
4. IMPORT: Import a cluster that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created cluster from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/github.com/couchbasecloud/terraform-provider-couchbase-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_clusters.existing_clusters: Reading...
data.capella_clusters.existing_clusters: Read complete after 0s

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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + clusters_list = {
      + data            = [
          + {
              + app_service_id   = null
              + audit            = {
                  + created_at  = "2023-10-03 21:04:45.895255387 +0000 UTC"
                  + created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  + modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + version     = 5
                }
              + availability     = {
                  + type = "multi"
                }
              + cloud_provider   = {
                  + cidr   = "10.0.0.0/23"
                  + region = "af-south-1"
                  + type   = "aws"
                }
              + couchbase_server = {
                  + version = "7.2.2"
                }
              + current_state    = "healthy"
              + description      = ""
              + id               = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              + name             = "quickrobertekahn"
              + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + service_groups   = [
                  + {
                      + node         = {
                          + compute = {
                              + cpu = 4
                              + ram = 16
                            }
                          + disk    = {
                              + iops    = 3000
                              + storage = 50
                              + type    = "gp3"
                            }
                        }
                      + num_of_nodes = 3
                      + services     = [
                          + "search",
                          + "index",
                          + "data",
                          + "query",
                        ]
                    },
                ]
              + support          = {
                  + plan     = "developer pro"
                  + timezone = "PT"
                }
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_cluster   = {
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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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

### Apply the Plan, in order to create a new Bucket

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/github.com/couchbasecloud/terraform-provider-couchbase-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_clusters.existing_clusters: Reading...
data.capella_clusters.existing_clusters: Read complete after 0s

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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + clusters_list = {
      + data            = [
          + {
              + app_service_id   = null
              + audit            = {
                  + created_at  = "2023-10-03 21:04:45.895255387 +0000 UTC"
                  + created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  + modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + version     = 5
                }
              + availability     = {
                  + type = "multi"
                }
              + cloud_provider   = {
                  + cidr   = "10.0.0.0/23"
                  + region = "af-south-1"
                  + type   = "aws"
                }
              + couchbase_server = {
                  + version = "7.2.2"
                }
              + current_state    = "healthy"
              + description      = ""
              + id               = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              + name             = "quickrobertekahn"
              + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + service_groups   = [
                  + {
                      + node         = {
                          + compute = {
                              + cpu = 4
                              + ram = 16
                            }
                          + disk    = {
                              + iops    = 3000
                              + storage = 50
                              + type    = "gp3"
                            }
                        }
                      + num_of_nodes = 3
                      + services     = [
                          + "search",
                          + "index",
                          + "data",
                          + "query",
                        ]
                    },
                ]
              + support          = {
                  + plan     = "developer pro"
                  + timezone = "PT"
                }
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_cluster   = {
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
      + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
capella_cluster.new_cluster: Creation complete after 3m4s [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

clusters_list = {
  "data" = tolist([
    {
      "app_service_id" = tostring(null)
      "audit" = {
        "created_at" = "2023-10-03 21:04:45.895255387 +0000 UTC"
        "created_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "modified_at" = "2023-10-03 21:08:45.110430897 +0000 UTC"
        "modified_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "version" = 5
      }
      "availability" = {
        "type" = "multi"
      }
      "cloud_provider" = {
        "cidr" = "10.0.0.0/23"
        "region" = "af-south-1"
        "type" = "aws"
      }
      "couchbase_server" = {
        "version" = "7.2.2"
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "name" = "quickrobertekahn"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "service_groups" = tolist([
        {
          "node" = {
            "compute" = {
              "cpu" = 4
              "ram" = 16
            }
            "disk" = {
              "iops" = 3000
              "storage" = 50
              "type" = "gp3"
            }
          }
          "num_of_nodes" = 3
          "services" = tolist([
            "search",
            "index",
            "data",
            "query",
          ])
        },
      ])
      "support" = {
        "plan" = "developer pro"
        "timezone" = "PT"
      }
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-10-04 19:30:24.442533243 +0000 UTC"
    "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "modified_at" = "2023-10-04 19:33:26.507037621 +0000 UTC"
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
  "id" = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
  "if_match" = tostring(null)
  "name" = "New Terraform Cluster"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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


### Note the Bucket ID for the new Bucket
Command: `terraform output new_cluster`

Sample Output:
```
$ terraform output new_cluster
{
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-10-04 19:30:24.442533243 +0000 UTC"
    "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "modified_at" = "2023-10-04 19:33:26.507037621 +0000 UTC"
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
  "id" = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
  "if_match" = tostring(null)
  "name" = "New Terraform Cluster"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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


In this case, the cluster ID for my new cluster is `f90c1d8a-c01f-4faa-860d-71cdcdf454f6`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.capella_clusters.existing_clusters
capella_cluster.new_cluster
```

## IMPORT
### Remove the resource `new_cluster` from the Terraform State file

Command: `terraform state rm capella_cluster.new_cluster`

Sample Output:
```
$ terraform state rm capella_cluster.new_cluster
Removed capella_cluster.new_cluster
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import capella_cluster.new_cluster id=<cluster_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import capella_cluster.new_cluster id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d`

Sample Output:
```
$ terraform import capella_cluster.new_cluster id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d
capella_cluster.new_cluster: Importing from ID "id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d"...
data.capella_clusters.existing_clusters: Reading...
capella_cluster.new_cluster: Import prepared!
  Prepared capella_cluster for import
capella_cluster.new_cluster: Refreshing state... [id=id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d]
data.capella_clusters.existing_clusters: Read complete after 1s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the cluster ID i.e. the ID of the resource that we want to import.
The second ID is the project ID i.e. the ID of the project to which the cluster belongs.
The third ID is the organization ID i.e. the ID of the organization to which the project belongs.


## UPDATE
### Let us edit the terraform.tfvars file to change the cluster configuration settings.

Command: `terraform apply -var 'support={plan="enterprise", timezone="IST"}'`

Sample Output:
```
$ terraform apply -var 'support={plan="enterprise", timezone="IST"}'
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/github.com/couchbasecloud/terraform-provider-couchbase-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_clusters.existing_clusters: Reading...
capella_cluster.new_cluster: Refreshing state... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6]
data.capella_clusters.existing_clusters: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_cluster.new_cluster will be updated in-place
  ~ resource "capella_cluster" "new_cluster" {
      + app_service_id   = (known after apply)
      ~ audit            = {
          ~ created_at  = "2023-10-04 19:30:24.442533243 +0000 UTC" -> (known after apply)
          ~ created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> (known after apply)
          ~ modified_at = "2023-10-04 19:33:26.507037621 +0000 UTC" -> (known after apply)
          ~ modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> (known after apply)
          ~ version     = 5 -> (known after apply)
        }
      ~ current_state    = "healthy" -> (known after apply)
      ~ etag             = "Version: 5" -> (known after apply)
        id               = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
        name             = "New Terraform Cluster"
      ~ service_groups   = [
          ~ {
              ~ services     = [
                  - "index",
                    "data",
                  + "index",
                    "query",
                ]
                # (2 unchanged attributes hidden)
            },
        ]
      ~ support          = {
          ~ plan     = "developer pro" -> "enterprise"
          ~ timezone = "PT" -> "IST"
        }
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ clusters_list = {
      ~ data            = [
          - {
              - app_service_id   = null
              - audit            = {
                  - created_at  = "2023-10-03 21:04:45.895255387 +0000 UTC"
                  - created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  - modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - version     = 5
                }
              - availability     = {
                  - type = "multi"
                }
              - cloud_provider   = {
                  - cidr   = "10.0.0.0/23"
                  - region = "af-south-1"
                  - type   = "aws"
                }
              - couchbase_server = {
                  - version = "7.2.2"
                }
              - current_state    = "healthy"
              - description      = ""
              - id               = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              - name             = "quickrobertekahn"
              - organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              - service_groups   = [
                  - {
                      - node         = {
                          - compute = {
                              - cpu = 4
                              - ram = 16
                            }
                          - disk    = {
                              - iops    = 3000
                              - storage = 50
                              - type    = "gp3"
                            }
                        }
                      - num_of_nodes = 3
                      - services     = [
                          - "search",
                          - "index",
                          - "data",
                          - "query",
                        ]
                    },
                ]
              - support          = {
                  - plan     = "developer pro"
                  - timezone = "PT"
                }
            },
            {
                app_service_id   = null
                audit            = {
                    created_at  = "2023-10-04 19:30:24.442533243 +0000 UTC"
                    created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
                    modified_at = "2023-10-04 19:33:26.507037621 +0000 UTC"
                    modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
                    version     = 5
                }
                availability     = {
                    type = "multi"
                }
                cloud_provider   = {
                    cidr   = "192.168.0.0/20"
                    region = "us-east-1"
                    type   = "aws"
                }
                couchbase_server = {
                    version = "7.1.5"
                }
                current_state    = "healthy"
                description      = "My first test cluster for multiple services."
                id               = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
                name             = "New Terraform Cluster"
                organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
                project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
                service_groups   = [
                    {
                        node         = {
                            compute = {
                                cpu = 4
                                ram = 16
                            }
                            disk    = {
                                iops    = 5000
                                storage = 50
                                type    = "io2"
                            }
                        }
                        num_of_nodes = 3
                        services     = [
                            "index",
                            "data",
                            "query",
                        ]
                    },
                ]
                support          = {
                    plan     = "developer pro"
                    timezone = "PT"
                }
            },
          + {
              + app_service_id   = null
              + audit            = {
                  + created_at  = "2023-10-03 21:04:45.895255387 +0000 UTC"
                  + created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  + modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  + version     = 5
                }
              + availability     = {
                  + type = "multi"
                }
              + cloud_provider   = {
                  + cidr   = "10.0.0.0/23"
                  + region = "af-south-1"
                  + type   = "aws"
                }
              + couchbase_server = {
                  + version = "7.2.2"
                }
              + current_state    = "healthy"
              + description      = ""
              + id               = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              + name             = "quickrobertekahn"
              + organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + service_groups   = [
                  + {
                      + node         = {
                          + compute = {
                              + cpu = 4
                              + ram = 16
                            }
                          + disk    = {
                              + iops    = 3000
                              + storage = 50
                              + type    = "gp3"
                            }
                        }
                      + num_of_nodes = 3
                      + services     = [
                          + "search",
                          + "index",
                          + "data",
                          + "query",
                        ]
                    },
                ]
              + support          = {
                  + plan     = "developer pro"
                  + timezone = "PT"
                }
            },
        ]
        # (2 unchanged elements hidden)
    }
  ~ new_cluster   = {
      ~ app_service_id   = null -> (known after apply)
      ~ audit            = {
          - created_at  = "2023-10-04 19:30:24.442533243 +0000 UTC"
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
          - modified_at = "2023-10-04 19:33:26.507037621 +0000 UTC"
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
          - version     = 5
        } -> (known after apply)
      ~ current_state    = "healthy" -> (known after apply)
      ~ etag             = "Version: 5" -> (known after apply)
        id               = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
        name             = "New Terraform Cluster"
      ~ service_groups   = [
          ~ {
              ~ services     = [
                  - "index",
                    "data",
                  + "index",
                    "query",
                ]
                # (2 unchanged elements hidden)
            },
        ]
      ~ support          = {
          ~ plan     = "developer pro" -> "enterprise"
          ~ timezone = "PT" -> "IST"
        }
        # (7 unchanged elements hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_cluster.new_cluster: Modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 10s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 20s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 30s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 40s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 50s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m0s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m10s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m20s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m30s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m40s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m50s elapsed]
capella_cluster.new_cluster: Still modifying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 2m0s elapsed]
capella_cluster.new_cluster: Modifications complete after 2m7s [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

clusters_list = {
  "data" = tolist([
    {
      "app_service_id" = tostring(null)
      "audit" = {
        "created_at" = "2023-10-04 19:30:24.442533243 +0000 UTC"
        "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
        "modified_at" = "2023-10-04 19:33:26.507037621 +0000 UTC"
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
        "version" = "7.1.5"
      }
      "current_state" = "healthy"
      "description" = "My first test cluster for multiple services."
      "id" = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
      "name" = "New Terraform Cluster"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
            "index",
            "data",
            "query",
          ])
        },
      ])
      "support" = {
        "plan" = "developer pro"
        "timezone" = "PT"
      }
    },
    {
      "app_service_id" = tostring(null)
      "audit" = {
        "created_at" = "2023-10-03 21:04:45.895255387 +0000 UTC"
        "created_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "modified_at" = "2023-10-03 21:08:45.110430897 +0000 UTC"
        "modified_by" = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
        "version" = 5
      }
      "availability" = {
        "type" = "multi"
      }
      "cloud_provider" = {
        "cidr" = "10.0.0.0/23"
        "region" = "af-south-1"
        "type" = "aws"
      }
      "couchbase_server" = {
        "version" = "7.2.2"
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "name" = "quickrobertekahn"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "service_groups" = tolist([
        {
          "node" = {
            "compute" = {
              "cpu" = 4
              "ram" = 16
            }
            "disk" = {
              "iops" = 3000
              "storage" = 50
              "type" = "gp3"
            }
          }
          "num_of_nodes" = 3
          "services" = tolist([
            "search",
            "index",
            "data",
            "query",
          ])
        },
      ])
      "support" = {
        "plan" = "developer pro"
        "timezone" = "PT"
      }
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2023-10-04 19:30:24.442533243 +0000 UTC"
    "created_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "modified_at" = "2023-10-04 19:43:19.445903132 +0000 UTC"
    "modified_by" = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
    "version" = 10
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
    "version" = "7.1.5"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 10"
  "id" = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
  "if_match" = tostring(null)
  "name" = "New Terraform Cluster"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
    "plan" = "enterprise"
    "timezone" = "IST"
  }
}
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
│  - hashicorp.com/couchabasecloud/capella in /Users/talina.shrotriya/workspace/github.com/couchbasecloud/terraform-provider-couchbase-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_clusters.existing_clusters: Reading...
capella_cluster.new_cluster: Refreshing state... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6]
data.capella_clusters.existing_clusters: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_cluster.new_cluster will be destroyed
  - resource "capella_cluster" "new_cluster" {
      - audit            = {
          - created_at  = "2023-10-04 19:30:24.442533243 +0000 UTC" -> null
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - modified_at = "2023-10-04 19:43:19.445903132 +0000 UTC" -> null
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC" -> null
          - version     = 10 -> null
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
          - version = "7.1.5" -> null
        }
      - current_state    = "healthy" -> null
      - description      = "My first test cluster for multiple services." -> null
      - etag             = "Version: 10" -> null
      - id               = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6" -> null
      - name             = "New Terraform Cluster" -> null
      - organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81" -> null
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
          - plan     = "enterprise" -> null
          - timezone = "IST" -> null
        }
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - clusters_list = {
      - data            = [
          - {
              - app_service_id   = null
              - audit            = {
                  - created_at  = "2023-10-04 19:30:24.442533243 +0000 UTC"
                  - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
                  - modified_at = "2023-10-04 19:43:19.445903132 +0000 UTC"
                  - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
                  - version     = 10
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
                  - version = "7.1.5"
                }
              - current_state    = "healthy"
              - description      = "My first test cluster for multiple services."
              - id               = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
              - name             = "New Terraform Cluster"
              - organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
                          - "index",
                          - "data",
                          - "query",
                        ]
                    },
                ]
              - support          = {
                  - plan     = "enterprise"
                  - timezone = "IST"
                }
            },
          - {
              - app_service_id   = null
              - audit            = {
                  - created_at  = "2023-10-03 21:04:45.895255387 +0000 UTC"
                  - created_by  = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  - modified_by = "7dfc7b93-a71b-4a2e-b5a7-4255ab29cab9"
                  - version     = 5
                }
              - availability     = {
                  - type = "multi"
                }
              - cloud_provider   = {
                  - cidr   = "10.0.0.0/23"
                  - region = "af-south-1"
                  - type   = "aws"
                }
              - couchbase_server = {
                  - version = "7.2.2"
                }
              - current_state    = "healthy"
              - description      = ""
              - id               = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              - name             = "quickrobertekahn"
              - organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              - service_groups   = [
                  - {
                      - node         = {
                          - compute = {
                              - cpu = 4
                              - ram = 16
                            }
                          - disk    = {
                              - iops    = 3000
                              - storage = 50
                              - type    = "gp3"
                            }
                        }
                      - num_of_nodes = 3
                      - services     = [
                          - "search",
                          - "index",
                          - "data",
                          - "query",
                        ]
                    },
                ]
              - support          = {
                  - plan     = "developer pro"
                  - timezone = "PT"
                }
            },
        ]
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    } -> null
  - new_cluster   = {
      - app_service_id   = null
      - audit            = {
          - created_at  = "2023-10-04 19:30:24.442533243 +0000 UTC"
          - created_by  = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
          - modified_at = "2023-10-04 19:43:19.445903132 +0000 UTC"
          - modified_by = "FYcHU0XHvw5tuQl4smAwyCKbGPMUZKUC"
          - version     = 10
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
          - version = "7.1.5"
        }
      - current_state    = "healthy"
      - description      = "My first test cluster for multiple services."
      - etag             = "Version: 10"
      - id               = "f90c1d8a-c01f-4faa-860d-71cdcdf454f6"
      - if_match         = null
      - name             = "New Terraform Cluster"
      - organization_id  = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id       = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
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
          - plan     = "enterprise"
          - timezone = "IST"
        }
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_cluster.new_cluster: Destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 1m50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 2m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 2m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 2m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 2m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 2m40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 2m50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 3m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 3m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 3m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=f90c1d8a-c01f-4faa-860d-71cdcdf454f6, 3m30s elapsed]
capella_cluster.new_cluster: Destruction complete after 3m32s

Destroy complete! Resources: 1 destroyed.
```
