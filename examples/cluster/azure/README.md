# Important notes:

1. Please rename the terraform.premium.tfvars and terraform.ultra.tfvars files to `terraform.template.tfvars` to use the configurations for Premium and Ultra type disks respectively.

2. Then please check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

# TODO: add examples for Premium and Ultra type disks
# Azure Cluster Example

This example shows how to create a Cluster on Azure, with autoexpansion enabled.

For autoexpansion it requires this disk variable in `terraform.template.tfvars`:

```
disk = {
  type = "P6"
  autoexpansion = true
}
```

It should be added to `variables.tf`:

```
type = object({
    size = optional(number)
    type = string
    iops = optional(number)
    autoexpansion = optional(bool)
  })
```

and referenced in `create_cluster.tf`:

```
disk = {
    type    = var.disk.type
    autoexpansion = var.disk.autoexpansion
}
```

Sample Output:

```
terraform apply -var-file terraform.template.tfvars

Terraform will perform the following actions:

  # couchbase-capella_cluster.new_cluster will be created
  + resource "couchbase-capella_cluster" "new_cluster" {
      + app_service_id     = (known after apply)
      + audit              = (known after apply)
      + availability       = {
          + type = "single"
        }
      + cloud_provider     = {
          + cidr   = "10.0.6.0/23"
          + region = "eastus"
          + type   = "azure"
        }
      + configuration_type = (known after apply)
      + couchbase_server   = (known after apply)
      + current_state      = (known after apply)
      + description        = "My first test cluster for multiple services."
      + etag               = (known after apply)
      + id                 = (known after apply)
      + name               = "New Terraform Azure Cluster 6"
      + organization_id    = "7dc36559-a544-4ce6-a132-5da412f350e9"
      + project_id         = "693f1cb1-982b-4c52-8119-3f5dc6828489"
      + service_groups     = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = true
                      + iops          = (known after apply)
                      + storage       = (known after apply)
                      + type          = "P6"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                ]
            },
        ]
      + support            = {
          + plan     = "basic"
          + timezone = "PT"
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + cluster_id  = (known after apply)
  + new_cluster = {
      + app_service_id     = (known after apply)
      + audit              = (known after apply)
      + availability       = {
          + type = "single"
        }
      + cloud_provider     = {
          + cidr   = "10.0.6.0/23"
          + region = "eastus"
          + type   = "azure"
        }
      + configuration_type = (known after apply)
      + couchbase_server   = (known after apply)
      + current_state      = (known after apply)
      + description        = "My first test cluster for multiple services."
      + etag               = (known after apply)
      + id                 = (known after apply)
      + if_match           = null
      + name               = "New Terraform Azure Cluster 6"
      + organization_id    = "7dc36559-a544-4ce6-a132-5da412f350e9"
      + project_id         = "693f1cb1-982b-4c52-8119-3f5dc6828489"
      + service_groups     = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = true
                      + iops          = (known after apply)
                      + storage       = (known after apply)
                      + type          = "P6"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                ]
            },
        ]
      + support            = {
          + plan     = "basic"
          + timezone = "PT"
        }
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


couchbase-capella_cluster.new_cluster: Creating...
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
couchbase-capella_cluster.new_cluster: Still creating... [3m30s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [3m40s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [3m50s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [4m0s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [4m10s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [4m20s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [4m30s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [4m40s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [4m50s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [5m0s elapsed]
couchbase-capella_cluster.new_cluster: Still creating... [5m10s elapsed]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

cluster_id = "10d5c496-502f-430e-99fb-25fa8bd7285a"
new_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2024-01-18 20:48:53.703590843 +0000 UTC"
    "created_by" = "GGZjZ6bOhAiuFJoJ0M1NoFHwm6RHZvkv"
    "modified_at" = "2024-01-18 20:54:04.672833562 +0000 UTC"
    "modified_by" = "GGZjZ6bOhAiuFJoJ0M1NoFHwm6RHZvkv"
    "version" = 3
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "10.0.6.0/23"
    "region" = "eastus"
    "type" = "azure"
  }
  "configuration_type" = "multiNode"
  "couchbase_server" = {
    "version" = "7.2"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "etag" = "Version: 3"
  "id" = "10d5c496-502f-430e-99fb-25fa8bd7285a"
  "if_match" = tostring(null)
  "name" = "New Terraform Azure Cluster 6"
  "organization_id" = "7dc36559-a544-4ce6-a132-5da412f350e9"
  "project_id" = "693f1cb1-982b-4c52-8119-3f5dc6828489"
  "service_groups" = toset([
    {
      "node" = {
        "compute" = {
          "cpu" = 4
          "ram" = 16
        }
        "disk" = {
          "autoexpansion" = true
          "iops" = 240
          "storage" = 64
          "type" = "P6"
        }
      }
      "num_of_nodes" = 3
      "services" = toset([
        "data",
      ])
    },
  ])
  "support" = {
    "plan" = "basic"
    "timezone" = "PT"
  }
}
```