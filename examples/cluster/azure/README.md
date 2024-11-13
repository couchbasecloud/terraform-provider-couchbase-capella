# Important notes:

1. Please rename the `terraform.premium.tfvars` and `terraform.ultra.tfvars` files to `terraform.template.tfvars` to use the configurations for Premium and Ultra type disks respectively.

2. To run `terraform apply` directly - please check the `terraform.template.tfvars` file, make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

# Azure Premium type Cluster Example

This example shows how to create a Cluster on Azure, with Premium type disks and autoexpansion enabled. 

Note that you cannot provide the iops and storage value for Premium type disks.

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

Create -

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

Update example-

Updating Timezone for an existing Azure cluster on developer pro plan - 
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cluster.new_cluster will be updated in-place
  ~ resource "couchbase-capella_cluster" "new_cluster" {
      + app_service_id                = (known after apply)
      ~ audit                         = {
          ~ created_at  = "2024-11-13 00:42:55.42572338 +0000 UTC" -> (known after apply)
          ~ created_by  = "CCaCemREbxLhyNfRg61d4ukDW75nWn1I" -> (known after apply)
          ~ modified_at = "2024-11-13 01:14:26.968616922 +0000 UTC" -> (known after apply)
          ~ modified_by = "apikey-CCaCemREbxLhyNfRg61d4ukDW75nWn1I" -> (known after apply)
          ~ version     = 7 -> (known after apply)
        } -> (known after apply)
      ~ connection_string             = "cb.epdsnt2cthimgrpn.aws-guardians.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ etag                          = "Version: 7" -> (known after apply)
        id                            = "0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f"
        name                          = "TF Azure P6 iops"
      ~ support                       = {
          ~ timezone = "GMT" -> "PT"
            # (1 unchanged attribute hidden)
        }
        # (9 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_cluster = {
      + app_service_id                = (known after apply)
      ~ audit                         = {
          - created_at  = "2024-11-13 00:42:55.42572338 +0000 UTC"
          - created_by  = "CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
          - modified_at = "2024-11-13 00:48:04.283065342 +0000 UTC"
          - modified_by = "apikey-CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
          - version     = 3
        } -> (known after apply)
      ~ connection_string             = "cb.epdsnt2cthimgrpn.aws-guardians.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ etag                          = "Version: 3" -> (known after apply)
        id                            = "0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f"
        name                          = "TF Azure P6 iops"
        # (11 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster.new_cluster: Modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 10s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 20s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 30s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 40s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 50s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 1m0s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 1m10s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 1m20s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 1m30s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 1m40s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 1m50s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f, 2m0s elapsed]
couchbase-capella_cluster.new_cluster: Modifications complete after 2m2s [id=0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

cluster_id = "0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f"
new_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2024-11-13 00:42:55.42572338 +0000 UTC"
    "created_by" = "CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
    "modified_at" = "2024-11-13 01:17:36.317030551 +0000 UTC"
    "modified_by" = "apikey-CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
    "version" = 11
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "10.0.0.0/23"
    "region" = "eastus"
    "type" = "azure"
  }
  "configuration_type" = "multiNode"
  "connection_string" = "cb.epdsnt2cthimgrpn.aws-guardians.nonprod-project-avengers.com"
  "couchbase_server" = {
    "version" = "7.6"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "enable_private_dns_resolution" = false
  "etag" = "Version: 11"
  "id" = "0952a9c5-fb5c-425c-8e1e-ca1ba4e7781f"
  "if_match" = tostring(null)
  "name" = "TF Azure P6 iops"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "c09035e8-3971-4b79-a6ec-bccd9056041f"
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
    "plan" = "developer pro"
    "timezone" = "PT"
  }
}

```



# Azure Ultra type cluster example-

This example shows how to create and update a Premium type Cluster on Azure.

For adding the storage, iops and autoexpansion, make sure you have these disk variable in `terraform.template.tfvars`:

```
disk = {
  type          = "Ultra"
  size          = 128
  iops          = 4000
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
          storage       = var.disk.size
          type          = var.disk.type
          iops          = var.disk.iops
          autoexpansion = var.disk.autoexpansion
        }
```

Create- 

Sample Output:

```
terraform apply


Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_cluster.new_cluster will be created
  + resource "couchbase-capella_cluster" "new_cluster" {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = {
          + type = "single"
        }
      + cloud_provider                = {
          + cidr   = "10.0.0.0/23"
          + region = "eastus"
          + type   = "azure"
        }
      + configuration_type            = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "My first test cluster for multiple services."
      + enable_private_dns_resolution = false
      + etag                          = (known after apply)
      + id                            = (known after apply)
      + name                          = "TF Azure Ultra"
      + organization_id               = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id                    = "c09035e8-3971-4b79-a6ec-bccd9056041f"
      + service_groups                = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = true
                      + iops          = 3000
                      + storage       = 64
                      + type          = "Ultra"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                ]
            },
        ]
      + support                       = {
          + plan     = "developer pro"
          + timezone = "PT"
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + cluster_id  = (known after apply)
  + new_cluster = {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = {
          + type = "single"
        }
      + cloud_provider                = {
          + cidr   = "10.0.0.0/23"
          + region = "eastus"
          + type   = "azure"
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
      + name                          = "TF Azure Ultra"
      + organization_id               = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id                    = "c09035e8-3971-4b79-a6ec-bccd9056041f"
      + service_groups                = [
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = true
                      + iops          = 3000
                      + storage       = 64
                      + type          = "Ultra"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                ]
            },
        ]
      + support                       = {
          + plan     = "developer pro"
          + timezone = "PT"
        }
    }
2024-11-13T10:49:30.808-0800 [DEBUG] command: asking for input: "\nDo you want to perform these actions?"

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster.new_cluster: Creating...
2024-11-13T10:49:32.348-0800 [INFO]  Starting apply for couchbase-capella_cluster.new_cluster
2024-11-13T10:49:32.348-0800 [DEBUG] skipping FixUpBlockAttrs
2024-11-13T10:49:32.348-0800 [DEBUG] couchbase-capella_cluster.new_cluster: applying the planned Create change
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


Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

cluster_id = "c142f957-edac-4bc1-aff4-90d7cde7b007"
new_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2024-11-13 18:49:32.43634667 +0000 UTC"
    "created_by" = "CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
    "modified_at" = "2024-11-13 18:55:26.536605251 +0000 UTC"
    "modified_by" = "apikey-CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
    "version" = 3
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "10.0.0.0/23"
    "region" = "eastus"
    "type" = "azure"
  }
  "configuration_type" = "multiNode"
  "connection_string" = "cb.sstd0a1ssoq-qz.aws-guardians.nonprod-project-avengers.com"
  "couchbase_server" = {
    "version" = "7.6"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "enable_private_dns_resolution" = false
  "etag" = "Version: 3"
  "id" = "c142f957-edac-4bc1-aff4-90d7cde7b007"
  "if_match" = tostring(null)
  "name" = "TF Azure Ultra"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "c09035e8-3971-4b79-a6ec-bccd9056041f"
  "service_groups" = toset([
    {
      "node" = {
        "compute" = {
          "cpu" = 4
          "ram" = 16
        }
        "disk" = {
          "autoexpansion" = true
          "iops" = 3000
          "storage" = 64
          "type" = "Ultra"
        }
      }
      "num_of_nodes" = 3
      "services" = toset([
        "data",
      ])
    },
  ])
  "support" = {
    "plan" = "developer pro"
    "timezone" = "PT"
  }
}

```

Update example-

```
terraform apply-

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cluster.new_cluster will be updated in-place
  ~ resource "couchbase-capella_cluster" "new_cluster" {
      + app_service_id                = (known after apply)
      ~ audit                         = {
          ~ created_at  = "2024-11-13 18:49:32.43634667 +0000 UTC" -> (known after apply)
          ~ created_by  = "CCaCemREbxLhyNfRg61d4ukDW75nWn1I" -> (known after apply)
          ~ modified_at = "2024-11-13 18:55:26.536605251 +0000 UTC" -> (known after apply)
          ~ modified_by = "apikey-CCaCemREbxLhyNfRg61d4ukDW75nWn1I" -> (known after apply)
          ~ version     = 3 -> (known after apply)
        } -> (known after apply)
      ~ connection_string             = "cb.sstd0a1ssoq-qz.aws-guardians.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ etag                          = "Version: 3" -> (known after apply)
        id                            = "c142f957-edac-4bc1-aff4-90d7cde7b007"
        name                          = "TF Azure Ultra"
      ~ service_groups                = [
          - {
              - node         = {
                  - compute = {
                      - cpu = 4 -> null
                      - ram = 16 -> null
                    } -> null
                  - disk    = {
                      - autoexpansion = true -> null
                      - iops          = 3000 -> null
                      - storage       = 64 -> null
                      - type          = "Ultra" -> null
                    } -> null
                } -> null
              - num_of_nodes = 3 -> null
              - services     = [
                  - "data",
                ] -> null
            },
          + {
              + node         = {
                  + compute = {
                      + cpu = 4
                      + ram = 16
                    }
                  + disk    = {
                      + autoexpansion = true
                      + iops          = 4000
                      + storage       = 128
                      + type          = "Ultra"
                    }
                }
              + num_of_nodes = 3
              + services     = [
                  + "data",
                ]
            },
        ]
      ~ support                       = {
          ~ timezone = "PT" -> "GMT"
            # (1 unchanged attribute hidden)
        }
        # (8 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_cluster = {
      + app_service_id                = (known after apply)
      ~ audit                         = {
          - created_at  = "2024-11-13 18:49:32.43634667 +0000 UTC"
          - created_by  = "CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
          - modified_at = "2024-11-13 18:55:26.536605251 +0000 UTC"
          - modified_by = "apikey-CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
          - version     = 3
        } -> (known after apply)
      ~ connection_string             = "cb.sstd0a1ssoq-qz.aws-guardians.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ etag                          = "Version: 3" -> (known after apply)
        id                            = "c142f957-edac-4bc1-aff4-90d7cde7b007"
        name                          = "TF Azure Ultra"
      ~ service_groups                = [
          ~ {
              ~ node         = {
                  ~ disk    = {
                      ~ iops          = 3000 -> 4000
                      ~ storage       = 64 -> 128
                        # (2 unchanged attributes hidden)
                    }
                    # (1 unchanged attribute hidden)
                }
                # (2 unchanged attributes hidden)
            },
        ]
      ~ support                       = {
          ~ timezone = "PT" -> "GMT"
            # (1 unchanged attribute hidden)
        }
        # (9 unchanged attributes hidden)
    }
2024-11-13T11:06:35.775-0800 [DEBUG] command: asking for input: "\nDo you want to perform these actions?"

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster.new_cluster: Modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007]
2024-11-13T11:06:37.595-0800 [INFO]  Starting apply for couchbase-capella_cluster.new_cluster
2024-11-13T11:06:37.595-0800 [DEBUG] skipping FixUpBlockAttrs
2024-11-13T11:06:37.595-0800 [DEBUG] couchbase-capella_cluster.new_cluster: applying the planned Update change
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 10s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 20s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 30s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 40s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 50s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 1m0s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 1m10s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 1m20s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 1m30s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 1m40s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 1m50s elapsed]
couchbase-capella_cluster.new_cluster: Still modifying... [id=c142f957-edac-4bc1-aff4-90d7cde7b007, 2m0s elapsed]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

cluster_id = "c142f957-edac-4bc1-aff4-90d7cde7b007"
new_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2024-11-13 18:49:32.43634667 +0000 UTC"
    "created_by" = "CCaCemREbxLhyNfRg61d4ukDW75nWn1I"
    "modified_at" = "2024-11-13 19:20:06.276514922 +0000 UTC"
    "modified_by" = "apikey-Bbmh0bvuZJWriP7lUIafeVh4RS7i8HHz"
    "version" = 7
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "10.0.0.0/23"
    "region" = "eastus"
    "type" = "azure"
  }
  "configuration_type" = "multiNode"
  "connection_string" = "cb.sstd0a1ssoq-qz.aws-guardians.nonprod-project-avengers.com"
  "couchbase_server" = {
    "version" = "7.6"
  }
  "current_state" = "healthy"
  "description" = "My first test cluster for multiple services."
  "enable_private_dns_resolution" = false
  "etag" = "Version: 7"
  "id" = "c142f957-edac-4bc1-aff4-90d7cde7b007"
  "if_match" = tostring(null)
  "name" = "TF Azure Ultra"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "c09035e8-3971-4b79-a6ec-bccd9056041f"
  "service_groups" = toset([
    {
      "node" = {
        "compute" = {
          "cpu" = 4
          "ram" = 16
        }
        "disk" = {
          "autoexpansion" = true
          "iops" = 4000
          "storage" = 128
          "type" = "Ultra"
        }
      }
      "num_of_nodes" = 3
      "services" = toset([
        "data",
      ])
    },
  ])
  "support" = {
    "plan" = "developer pro"
    "timezone" = "GMT"
  }
}


```