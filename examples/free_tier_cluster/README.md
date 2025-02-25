# Capella Free Tier Cluster examples
This example demonstrates how to create and manage a free-tier cluster in Capella.

It provisions a new free-tier cluster within the selected Capella project and tenant while managing the cluster.

To run this example, configure your Couchbase Capella provider as described in the README file located in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new free-ter in Capella as stated in the `create_free_tier_cluster.tf` file.
2. UPDATE: Update the free-tier cluster  configuration using Terraform.
4. IMPORT: Import a free-tier cluster that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created free-tier cluster from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.


## CREATE & LIST
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

  # couchbase-capella_cluster_free_tier.new_cluster_free_tier will be created
  + resource "couchbase-capella_cluster_free_tier" "new_cluster_free_tier" {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "192.168.0.0/20"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "new test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster modifed"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + cluster_free_tierid   = (known after apply)
  + new_free_tier_cluster = {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "192.168.0.0/20"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "new test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster modifed"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

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

  # couchbase-capella_cluster_free_tier.new_cluster_free_tier will be created
  + resource "couchbase-capella_cluster_free_tier" "new_cluster_free_tier" {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "192.168.0.0/20"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "new test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + cluster_free_tierid   = (known after apply)
  + new_free_tier_cluster = {
      + app_service_id                = (known after apply)
      + audit                         = (known after apply)
      + availability                  = (known after apply)
      + cloud_provider                = {
          + cidr   = "192.168.0.0/20"
          + region = "us-east-2"
          + type   = "aws"
        }
      + cmek_id                       = (known after apply)
      + connection_string             = (known after apply)
      + couchbase_server              = (known after apply)
      + current_state                 = (known after apply)
      + description                   = "new test cluster for multiple services"
      + enable_private_dns_resolution = (known after apply)
      + id                            = (known after apply)
      + name                          = "New free tier cluster modifed"
      + organization_id               = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                    = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_groups                = (known after apply)
      + support                       = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster_free_tier.new_cluster_free_tier: Creating...
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [10s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [20s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [30s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [40s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [50s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [1m0s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [1m10s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [1m20s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [1m30s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [1m40s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [1m50s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still creating... [2m0s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Creation complete after 2m3s [id=a915df3b-c195-43f2-84ea-943e6418bf2d]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

cluster_free_tierid = "a915df3b-c195-43f2-84ea-943e6418bf2d"
new_free_tier_cluster = {
  "app_service_id" = tostring(null)
  "audit" = {
    "created_at" = "2025-02-25 09:13:23.970073027 +0000 UTC"
    "created_by" = "ffffffff0000000"
    "modified_at" = "2025-02-25 09:15:19.148350524 +0000 UTC"
    "modified_by" = "apikey-apikey-ffffffffff000000"
    "version" = 5
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "192.168.0.0/20"
    "region" = "us-east-2"
    "type" = "aws"
  }
  "cmek_id" = tostring(null)
  "connection_string" = "cb.nlojldxiegctpevd.sandbox.nonprod-project-avengers.com"
  "couchbase_server" = {
    "version" = "7.6"
  }
  "current_state" = "healthy"
  "description" = "new test cluster for multiple services"
  "enable_private_dns_resolution" = false
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

## UPDATE
### Update the resource `new_cluster_free_tier` in the Terraform State file

Free tier supports updating only name and description, in the below examples the name and description have been modified
```
terraform apply -var-file=terraform.template.tfvars
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Refreshing state... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cluster_free_tier.new_cluster_free_tier will be updated in-place
  ~ resource "couchbase-capella_cluster_free_tier" "new_cluster_free_tier" {
      ~ audit                         = {
          ~ created_at  = "2025-02-18 15:20:04.028587471 +0000 UTC" -> (known after apply)
          ~ created_by  = "9YN7VgzYU0gapElrX5suRY53MFEo0s1v" -> (known after apply)
          ~ modified_at = "2025-02-18 16:02:52.261505397 +0000 UTC" -> (known after apply)
          ~ modified_by = "9YN7VgzYU0gapElrX5suRY53MFEo0s1v" -> (known after apply)
          ~ version     = 7 -> (known after apply)
        } -> (known after apply)
      ~ availability                  = {
          ~ type = "single" -> (known after apply)
        } -> (known after apply)
      ~ connection_string             = "cb.wqd1ieecvgkz4ri.sandbox.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ enable_private_dns_resolution = false -> (known after apply)
        id                            = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ name                          = "New Terraform Cluster name" -> "New Terraform Cluster name new"
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
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_free_tier_cluster = {
      ~ audit                         = {
          - created_at  = "2025-02-18 15:20:04.028587471 +0000 UTC"
          - created_by  = "9YN7VgzYU0gapElrX5suRY53MFEo0s1v"
          - modified_at = "2025-02-18 16:02:52.261505397 +0000 UTC"
          - modified_by = "9YN7VgzYU0gapElrX5suRY53MFEo0s1v"
          - version     = 7
        } -> (known after apply)
      ~ availability                  = {
          - type = "single"
        } -> (known after apply)
      ~ connection_string             = "cb.wqd1ieecvgkz4ri.sandbox.nonprod-project-avengers.com" -> (known after apply)
      ~ current_state                 = "healthy" -> (known after apply)
      ~ enable_private_dns_resolution = false -> (known after apply)
        id                            = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ name                          = "New Terraform Cluster name" -> "New Terraform Cluster name new"
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
        # (5 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster_free_tier.new_cluster_free_tier: Modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 10s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 20s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 30s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 40s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 50s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 1m0s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 1m10s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 1m20s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 1m30s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 1m40s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 1m50s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still modifying... [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e, 2m0s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Modifications complete after 2m2s [id=1b8fd9d2-ab4f-4f88-855e-95a651448d9e]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

cluster_free_tierid = "1b8fd9d2-ab4f-4f88-855e-95a651448d9e"
new_free_tier_cluster = {
  "audit" = {
    "created_at" = "2025-02-18 15:20:04.028587471 +0000 UTC"
    "created_by" = "9YN7VgzYU0gapElrX5suRY53MFEo0s1v"
    "modified_at" = "2025-02-18 16:07:33.895186552 +0000 UTC"
    "modified_by" = "9YN7VgzYU0gapElrX5suRY53MFEo0s1v"
    "version" = 8
  }
  "availability" = {
    "type" = "single"
  }
  "cloud_provider" = {
    "cidr" = "192.168.0.0/20"
    "region" = "us-east-2"
    "type" = "aws"
  }
  "connection_string" = "cb.wqd1ieecvgkz4ri.sandbox.nonprod-project-avengers.com"
  "couchbase_server" = {
    "version" = "7.6"
  }
  "current_state" = "healthy"
  "description" = "new test cluster for multiple services."
  "enable_private_dns_resolution" = false
  "id" = "1b8fd9d2-ab4f-4f88-855e-95a651448d9e"
  "name" = "New Terraform Cluster name new"
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
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Refreshing state... [id=a915df3b-c195-43f2-84ea-943e6418bf2d]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_cluster_free_tier.new_cluster_free_tier will be destroyed
  - resource "couchbase-capella_cluster_free_tier" "new_cluster_free_tier" {
      - audit                         = {
          - created_at  = "2025-02-25 09:13:23.970073027 +0000 UTC" -> null
          - created_by  = "SwfPdQn662WcD0h0PLmaifR2IHSt0fIN" -> null
          - modified_at = "2025-02-25 09:15:19.148350524 +0000 UTC" -> null
          - modified_by = "apikey-apikey-SwfPdQn662WcD0h0PLmaifR2IHSt0fIN" -> null
          - version     = 5 -> null
        } -> null
      - availability                  = {
          - type = "single" -> null
        } -> null
      - cloud_provider                = {
          - cidr   = "192.168.0.0/20" -> null
          - region = "us-east-2" -> null
          - type   = "aws" -> null
        } -> null
      - connection_string             = "cb.nlojldxiegctpevd.sandbox.nonprod-project-avengers.com" -> null
      - couchbase_server              = {
          - version = "7.6" -> null
        } -> null
      - current_state                 = "healthy" -> null
      - description                   = "new test cluster for multiple services modified" -> null
      - enable_private_dns_resolution = false -> null
      - id                            = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - name                          = "New free tier cluster modifed" -> null
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
  - cluster_free_tierid   = "ffffffff-aaaa-1414-eeee-000000000000" -> null
  - new_free_tier_cluster = {
      - app_service_id                = null
      - audit                         = {
          - created_at  = "2025-02-25 09:13:23.970073027 +0000 UTC"
          - created_by  = "SwfPdQn662WcD0h0PLmaifR2IHSt0fIN"
          - modified_at = "2025-02-25 09:15:19.148350524 +0000 UTC"
          - modified_by = "apikey-apikey-SwfPdQn662WcD0h0PLmaifR2IHSt0fIN"
          - version     = 5
        }
      - availability                  = {
          - type = "single"
        }
      - cloud_provider                = {
          - cidr   = "192.168.0.0/20"
          - region = "us-east-2"
          - type   = "aws"
        }
      - cmek_id                       = null
      - connection_string             = "cb.nlojldxiegctpevd.sandbox.nonprod-project-avengers.com"
      - couchbase_server              = {
          - version = "7.6"
        }
      - current_state                 = "healthy"
      - description                   = "new test cluster for multiple services modified"
      - enable_private_dns_resolution = false
      - id                            = "ffffffff-aaaa-1414-eeee-000000000000"
      - name                          = "New free tier cluster modifed"
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

couchbase-capella_cluster_free_tier.new_cluster_free_tier: Destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 10s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 20s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 30s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 40s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 50s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 1m0s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 1m10s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 1m20s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 1m30s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 1m40s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 1m50s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Still destroying... [id=a915df3b-c195-43f2-84ea-943e6418bf2d, 2m0s elapsed]
couchbase-capella_cluster_free_tier.new_cluster_free_tier: Destruction complete after 2m4s

Destroy complete! Resources: 1 destroyed.
```