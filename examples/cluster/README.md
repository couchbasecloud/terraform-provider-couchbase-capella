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
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
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
      + connection_string = (known after apply)
      + couchbase_server = {
          + version = "7.1"
        }
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + name             = "New Terraform Cluster"
      + organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
                  + created_by  = "ffffffff-aaaa-1414-eeee-00000000000"
                  + modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-00000000000"
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
              + connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
              + couchbase_server = {
                  + version = "7.2.2"
                }
              + current_state    = "healthy"
              + description      = ""
              + id               = "ffffffff-aaaa-1414-eeee-00000000000"
              + name             = "quickrobertekahn"
              + organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
              + project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
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
      + connection_string = (known after apply)
      + couchbase_server = {
          + version = "7.1"
        }
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + if_match         = null
      + name             = "New Terraform Cluster"
      + organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
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
      + connection_string = (known after apply)
      + couchbase_server = {
          + version = "7.1"
        }
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + name             = "New Terraform Cluster"
      + organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
                  + created_by  = "ffffffff-aaaa-1414-eeee-00000000000"
                  + modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-00000000000"
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
              + connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com" 
              + couchbase_server = {
                  + version = "7.2.2"
                }
              + current_state    = "healthy"
              + description      = ""
              + id               = "ffffffff-aaaa-1414-eeee-00000000000"
              + name             = "quickrobertekahn"
              + organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
              + project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
      + organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
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
      + connection_string = (known after apply)
      + current_state    = (known after apply)
      + description      = "My first test cluster for multiple services."
      + etag             = (known after apply)
      + id               = (known after apply)
      + if_match         = null
      + name             = "New Terraform Cluster"
      + organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
      + project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
capella_cluster.new_cluster: Creation complete after 3m4s [id=ffffffff-aaaa-1414-eeee-00000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

clusters_list = {
  "data" = tolist([
    {
      "app_service_id" = tostring(null)
      "audit" = {
        "created_at" = "2023-10-03 21:04:45.895255387 +0000 UTC"
        "created_by" = "ffffffff-aaaa-1414-eeee-00000000000"
        "modified_at" = "2023-10-03 21:08:45.110430897 +0000 UTC"
        "modified_by" = "ffffffff-aaaa-1414-eeee-00000000000"
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
      "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
      "couchbase_server" = {
        "version" = "7.2.2"
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "name" = "quickrobertekahn"
      "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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
  "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
  "couchbase_server" = {
    "version" = "7.1"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "if_match" = tostring(null)
  "name" = "New Terraform Cluster"
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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
  "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
  "couchbase_server" = {
    "version" = "7.1"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 5"
  "id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "if_match" = tostring(null)
  "name" = "New Terraform Cluster"
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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


In this case, the cluster ID for my new cluster is `ffffffff-aaaa-1414-eeee-00000000000`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_clusters.existing_clusters
couchbase-capella_cluster.new_cluster
```

## IMPORT
### Remove the resource `new_cluster` from the Terraform State file

Command: `terraform state rm couchbase-capella_cluster.new_cluster`

Sample Output:
```
$ terraform state rm couchbase-capella_cluster.new_cluster
Removed couchbase-capella_cluster.new_cluster
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_cluster.new_cluster id=<cluster_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_cluster.new_cluster id=ffffffff-aaaa-1414-eeee-00000000000,project_id=ffffffff-aaaa-1414-eeee-00000000000,organization_id=ffffffff-aaaa-1414-eeee-00000000000`

Sample Output:
```
$ terraform import couchbase-capella_cluster.new_cluster id=ffffffff-aaaa-1414-eeee-00000000000,project_id=ffffffff-aaaa-1414-eeee-00000000000,organization_id=ffffffff-aaaa-1414-eeee-00000000000
capella_cluster.new_cluster: Importing from ID "id=ffffffff-aaaa-1414-eeee-00000000000,project_id=ffffffff-aaaa-1414-eeee-00000000000,organization_id=ffffffff-aaaa-1414-eeee-00000000000"...
data.capella_clusters.existing_clusters: Reading...
capella_cluster.new_cluster: Import prepared!
  Prepared capella_cluster for import
capella_cluster.new_cluster: Refreshing state... [id=id=ffffffff-aaaa-1414-eeee-00000000000,project_id=ffffffff-aaaa-1414-eeee-00000000000,organization_id=ffffffff-aaaa-1414-eeee-00000000000]
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
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_clusters.existing_clusters: Reading...
capella_cluster.new_cluster: Refreshing state... [id=ffffffff-aaaa-1414-eeee-00000000000]
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
        id               = "ffffffff-aaaa-1414-eeee-00000000000"
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
                  - created_by  = "ffffffff-aaaa-1414-eeee-00000000000"
                  - modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  - modified_by = "ffffffff-aaaa-1414-eeee-00000000000"
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
              - connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
              - couchbase_server = {
                  - version = "7.2.2"
                }
              - current_state    = "healthy"
              - description      = ""
              - id               = "ffffffff-aaaa-1414-eeee-00000000000"
              - name             = "quickrobertekahn"
              - organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
              - project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
                connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
                couchbase_server = {
                    version = "7.1.5"
                }
                current_state    = "healthy"
                description      = "My first test cluster for multiple services."
                id               = "ffffffff-aaaa-1414-eeee-00000000000"
                name             = "New Terraform Cluster"
                organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
                project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
                  + created_by  = "ffffffff-aaaa-1414-eeee-00000000000"
                  + modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-00000000000"
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
              + connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
              + couchbase_server = {
                  + version = "7.2.2"
                }
              + current_state    = "healthy"
              + description      = ""
              + id               = "ffffffff-aaaa-1414-eeee-00000000000"
              + name             = "quickrobertekahn"
              + organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
              + project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
      ~ connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com" -> (known after apply)
      ~ current_state    = "healthy" -> (known after apply)
      ~ etag             = "Version: 5" -> (known after apply)
        id               = "ffffffff-aaaa-1414-eeee-00000000000"
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

capella_cluster.new_cluster: Modifying... [id=ffffffff-aaaa-1414-eeee-00000000000]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 10s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 20s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 30s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 40s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 50s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m0s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m10s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m20s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m30s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m40s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m50s elapsed]
capella_cluster.new_cluster: Still modifying... [id=ffffffff-aaaa-1414-eeee-00000000000, 2m0s elapsed]
capella_cluster.new_cluster: Modifications complete after 2m7s [id=ffffffff-aaaa-1414-eeee-00000000000]

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
      "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
      "couchbase_server" = {
        "version" = "7.1.5"
      }
      "current_state" = "healthy"
      "description" = "My first test cluster for multiple services."
      "id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "name" = "New Terraform Cluster"
      "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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
        "created_by" = "ffffffff-aaaa-1414-eeee-00000000000"
        "modified_at" = "2023-10-03 21:08:45.110430897 +0000 UTC"
        "modified_by" = "ffffffff-aaaa-1414-eeee-00000000000"
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
      "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
      "couchbase_server" = {
        "version" = "7.2.2"
      }
      "current_state" = "healthy"
      "description" = ""
      "id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "name" = "quickrobertekahn"
      "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
      "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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
  "connection_string" = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
  "couchbase_server" = {
    "version" = "7.1.5"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 10"
  "id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "if_match" = tostring(null)
  "name" = "New Terraform Cluster"
  "organization_id" = "ffffffff-aaaa-1414-eeee-00000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-00000000000"
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
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.capella_clusters.existing_clusters: Reading...
capella_cluster.new_cluster: Refreshing state... [id=ffffffff-aaaa-1414-eeee-00000000000]
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
      - connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com" -> null
      - couchbase_server = {
          - version = "7.1.5" -> null
        }
      - current_state    = "healthy" -> null
      - description      = "My first test cluster for multiple services." -> null
      - etag             = "Version: 10" -> null
      - id               = "ffffffff-aaaa-1414-eeee-00000000000" -> null
      - name             = "New Terraform Cluster" -> null
      - organization_id  = "ffffffff-aaaa-1414-eeee-00000000000" -> null
      - project_id       = "ffffffff-aaaa-1414-eeee-00000000000" -> null
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
              - connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
              - couchbase_server = {
                  - version = "7.1.5"
                }
              - current_state    = "healthy"
              - description      = "My first test cluster for multiple services."
              - id               = "ffffffff-aaaa-1414-eeee-00000000000"
              - name             = "New Terraform Cluster"
              - organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
              - project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
                  - created_by  = "ffffffff-aaaa-1414-eeee-00000000000"
                  - modified_at = "2023-10-03 21:08:45.110430897 +0000 UTC"
                  - modified_by = "ffffffff-aaaa-1414-eeee-00000000000"
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
              - connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
              - couchbase_server = {
                  - version = "7.2.2"
                }
              - current_state    = "healthy"
              - description      = ""
              - id               = "ffffffff-aaaa-1414-eeee-00000000000"
              - name             = "quickrobertekahn"
              - organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
              - project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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
      - organization_id = "ffffffff-aaaa-1414-eeee-00000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-00000000000"
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
      - connection_string = "couchbases://cb.xxxxxxxxxxxxxx.cloud.couchbase.com"
      - couchbase_server = {
          - version = "7.1.5"
        }
      - current_state    = "healthy"
      - description      = "My first test cluster for multiple services."
      - etag             = "Version: 10"
      - id               = "ffffffff-aaaa-1414-eeee-00000000000"
      - if_match         = null
      - name             = "New Terraform Cluster"
      - organization_id  = "ffffffff-aaaa-1414-eeee-00000000000"
      - project_id       = "ffffffff-aaaa-1414-eeee-00000000000"
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

capella_cluster.new_cluster: Destroying... [id=ffffffff-aaaa-1414-eeee-00000000000]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 1m50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 2m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 2m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 2m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 2m30s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 2m40s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 2m50s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 3m0s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 3m10s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 3m20s elapsed]
capella_cluster.new_cluster: Still destroying... [id=ffffffff-aaaa-1414-eeee-00000000000, 3m30s elapsed]
capella_cluster.new_cluster: Destruction complete after 3m32s

Destroy complete! Resources: 1 destroyed.
```
