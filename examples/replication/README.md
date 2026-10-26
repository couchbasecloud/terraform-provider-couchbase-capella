# Capella XDCR Replication Example

This example shows how to create and manage an XDCR replication in Capella.

It replicates the `airline` and `airport` collections of the `inventory` scope from a bucket on a source cluster to a bucket on a target cluster. It uses a filter to replicate only documents whose key starts with `airline_`, and it removes the TTL (expiry time) from the replicated documents.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Before you start

- Both buckets must exist. Use the bucket **ID**, not the bucket name. The bucket ID is the URL-safe base64 encoding of the bucket name, for example the `id` attribute of a `couchbase-capella_bucket` resource.
- The scopes and collections in `mappings` must exist on both buckets.
- Both buckets must use the same conflict resolution type.
- Both clusters must be healthy. If a cluster is not healthy when Terraform creates the replication, for example during a rebalance, or soon after the last replication between the two clusters is deleted, Terraform tries the create again for up to 15 minutes.
- A refresh, for example in `terraform plan`, fails with a 422 error while the source cluster is not healthy. Run the command again when the cluster is healthy.
- For an external target, set `target.type = "external"`, use the cluster reference name as `target.cluster` and the bucket name as `target.bucket`. The cluster reference must already exist. An external target supports only `oneWay`.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new replication in Capella as stated in the `create_replication.tf` file.
2. IMPORT: Import a replication that exists in Capella but not in the terraform state file.
3. UPDATE: Update the replication settings using Terraform.
4. DESTROY: Delete the newly created replication from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_replication.new_replication will be created
  + resource "couchbase-capella_replication" "new_replication" {
      + audit                  = (known after apply)
      + changes_left           = (known after apply)
      + cluster_id             = "<cluster_id>"
      + direction              = "oneWay"
      + error                  = (known after apply)
      + filter                 = {
          + document_exclude_options = {
              + binary     = false
              + deletion   = false
              + expiration = false
              + ttl        = true
            }
          + expressions              = {
              + reg_ex        = "REGEXP_CONTAINS(META().id, '^airline_')"
              + skip_restream = false
            }
        }
      + id                     = (known after apply)
      + mappings               = [
          + {
              + collections  = [
                  + {
                      + source_collection = "airline"
                      + target_collection = "airline"
                    },
                  + {
                      + source_collection = "airport"
                      + target_collection = "airport"
                    },
                ]
              + source_scope = "inventory"
              + target_scope = "inventory"
            },
        ]
      + network_usage_limit    = 0
      + organization_id        = "<organization_id>"
      + priority               = "high"
      + project_id             = "<project_id>"
      + reverse_replication_id = (known after apply)
      + source_bucket          = "<source_bucket_id>"
      + source_details         = (known after apply)
      + status                 = (known after apply)
      + target                 = {
          + bucket  = "<target_bucket_id>"
          + cluster = "<target_cluster_id>"
          + type    = "capella"
        }
      + target_details         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_replication = {
      + audit                  = (known after apply)
      + changes_left           = (known after apply)
      + cluster_id             = "<cluster_id>"
      + direction              = "oneWay"
      + error                  = (known after apply)
      + filter                 = {
          + document_exclude_options = {
              + binary     = false
              + deletion   = false
              + expiration = false
              + ttl        = true
            }
          + expressions              = {
              + reg_ex        = "REGEXP_CONTAINS(META().id, '^airline_')"
              + skip_restream = false
            }
        }
      + id                     = (known after apply)
      + mappings               = [
          + {
              + collections  = [
                  + {
                      + source_collection = "airline"
                      + target_collection = "airline"
                    },
                  + {
                      + source_collection = "airport"
                      + target_collection = "airport"
                    },
                ]
              + source_scope = "inventory"
              + target_scope = "inventory"
            },
        ]
      + network_usage_limit    = 0
      + organization_id        = "<organization_id>"
      + priority               = "high"
      + project_id             = "<project_id>"
      + reverse_replication_id = (known after apply)
      + source_bucket          = "<source_bucket_id>"
      + source_details         = (known after apply)
      + status                 = (known after apply)
      + target                 = {
          + bucket  = "<target_bucket_id>"
          + cluster = "<target_cluster_id>"
          + type    = "capella"
        }
      + target_details         = (known after apply)
    }
  + replication_id  = (known after apply)

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Replication

Command: `terraform apply`

The create call is asynchronous. Terraform waits until Capella has created the replication. The first replication between two clusters takes longer, because Capella builds the network link between them.

If the create fails after the replication was made, the next `terraform apply` fails with a conflict error that shows the replication ID. Import the replication with that ID (see IMPORT).

Sample Output:
```
$ terraform apply
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_replication.new_replication will be created
  + resource "couchbase-capella_replication" "new_replication" {
      + audit                  = (known after apply)
      + changes_left           = (known after apply)
      + cluster_id             = "<cluster_id>"
      + direction              = "oneWay"
      + error                  = (known after apply)
      + filter                 = {
          + document_exclude_options = {
              + binary     = false
              + deletion   = false
              + expiration = false
              + ttl        = true
            }
          + expressions              = {
              + reg_ex        = "REGEXP_CONTAINS(META().id, '^airline_')"
              + skip_restream = false
            }
        }
      + id                     = (known after apply)
      + mappings               = [
          + {
              + collections  = [
                  + {
                      + source_collection = "airline"
                      + target_collection = "airline"
                    },
                  + {
                      + source_collection = "airport"
                      + target_collection = "airport"
                    },
                ]
              + source_scope = "inventory"
              + target_scope = "inventory"
            },
        ]
      + network_usage_limit    = 0
      + organization_id        = "<organization_id>"
      + priority               = "high"
      + project_id             = "<project_id>"
      + reverse_replication_id = (known after apply)
      + source_bucket          = "<source_bucket_id>"
      + source_details         = (known after apply)
      + status                 = (known after apply)
      + target                 = {
          + bucket  = "<target_bucket_id>"
          + cluster = "<target_cluster_id>"
          + type    = "capella"
        }
      + target_details         = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_replication = {
      + audit                  = (known after apply)
      + changes_left           = (known after apply)
      + cluster_id             = "<cluster_id>"
      + direction              = "oneWay"
      + error                  = (known after apply)
      + filter                 = {
          + document_exclude_options = {
              + binary     = false
              + deletion   = false
              + expiration = false
              + ttl        = true
            }
          + expressions              = {
              + reg_ex        = "REGEXP_CONTAINS(META().id, '^airline_')"
              + skip_restream = false
            }
        }
      + id                     = (known after apply)
      + mappings               = [
          + {
              + collections  = [
                  + {
                      + source_collection = "airline"
                      + target_collection = "airline"
                    },
                  + {
                      + source_collection = "airport"
                      + target_collection = "airport"
                    },
                ]
              + source_scope = "inventory"
              + target_scope = "inventory"
            },
        ]
      + network_usage_limit    = 0
      + organization_id        = "<organization_id>"
      + priority               = "high"
      + project_id             = "<project_id>"
      + reverse_replication_id = (known after apply)
      + source_bucket          = "<source_bucket_id>"
      + source_details         = (known after apply)
      + status                 = (known after apply)
      + target                 = {
          + bucket  = "<target_bucket_id>"
          + cluster = "<target_cluster_id>"
          + type    = "capella"
        }
      + target_details         = (known after apply)
    }
  + replication_id  = (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_replication.new_replication: Creating...
couchbase-capella_replication.new_replication: Still creating... [00m10s elapsed]
couchbase-capella_replication.new_replication: Still creating... [00m20s elapsed]
couchbase-capella_replication.new_replication: Still creating... [00m30s elapsed]
couchbase-capella_replication.new_replication: Still creating... [00m40s elapsed]
couchbase-capella_replication.new_replication: Still creating... [00m50s elapsed]
couchbase-capella_replication.new_replication: Still creating... [01m00s elapsed]
couchbase-capella_replication.new_replication: Still creating... [01m10s elapsed]
couchbase-capella_replication.new_replication: Still creating... [01m20s elapsed]
couchbase-capella_replication.new_replication: Still creating... [01m30s elapsed]
couchbase-capella_replication.new_replication: Still creating... [01m40s elapsed]
couchbase-capella_replication.new_replication: Still creating... [01m50s elapsed]
couchbase-capella_replication.new_replication: Creation complete after 2m0s [id=<replication_id>]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_replication = {
  "audit" = {
    "created_at" = "2026-10-13 22:26:18.630284798 +0000 UTC"
    "created_by" = "apikey-<api_key_id>"
  }
  "changes_left" = 0
  "cluster_id" = "<cluster_id>"
  "direction" = "oneWay"
  "error" = tostring(null)
  "filter" = {
    "document_exclude_options" = {
      "binary" = false
      "deletion" = false
      "expiration" = false
      "ttl" = true
    }
    "expressions" = {
      "reg_ex" = "REGEXP_CONTAINS(META().id, '^airline_')"
      "skip_restream" = false
    }
  }
  "id" = "<replication_id>"
  "mappings" = toset([
    {
      "collections" = toset([
        {
          "source_collection" = "airline"
          "target_collection" = "airline"
        },
        {
          "source_collection" = "airport"
          "target_collection" = "airport"
        },
      ])
      "source_scope" = "inventory"
      "target_scope" = "inventory"
    },
  ])
  "network_usage_limit" = 0
  "organization_id" = "<organization_id>"
  "priority" = "high"
  "project_id" = "<project_id>"
  "reverse_replication_id" = tostring(null)
  "source_bucket" = "<source_bucket_id>"
  "source_details" = {
    "bucket" = {
      "conflict_resolution_type" = "seqno"
      "id" = "<source_bucket_id>"
      "name" = "xdcr_src"
    }
    "cluster" = {
      "id" = "<cluster_id>"
      "name" = "cluster-a"
    }
    "project" = {
      "id" = "<project_id>"
      "name" = "my-project"
    }
  }
  "status" = "running"
  "target" = {
    "bucket" = "<target_bucket_id>"
    "cluster" = "<target_cluster_id>"
    "type" = "capella"
  }
  "target_details" = {
    "bucket" = {
      "conflict_resolution_type" = "seqno"
      "id" = "<target_bucket_id>"
      "name" = "xdcr_tgt"
    }
    "cluster" = {
      "id" = "<target_cluster_id>"
      "name" = "cluster-b"
    }
    "project" = {
      "id" = "<project_id>"
      "name" = "my-project"
    }
    "type" = "capella"
  }
}
replication_id = "<replication_id>"
```

### Note the Replication ID for the new Replication

Command: `terraform output new_replication`

Sample Output:
```
$ terraform output new_replication
{
  "audit" = {
    "created_at" = "2026-10-13 22:26:18.630284798 +0000 UTC"
    "created_by" = "apikey-<api_key_id>"
  }
  "changes_left" = 0
  "cluster_id" = "<cluster_id>"
  "direction" = "oneWay"
  "error" = tostring(null)
  "filter" = {
    "document_exclude_options" = {
      "binary" = false
      "deletion" = false
      "expiration" = false
      "ttl" = true
    }
    "expressions" = {
      "reg_ex" = "REGEXP_CONTAINS(META().id, '^airline_')"
      "skip_restream" = false
    }
  }
  "id" = "<replication_id>"
  "mappings" = toset([
    {
      "collections" = toset([
        {
          "source_collection" = "airline"
          "target_collection" = "airline"
        },
        {
          "source_collection" = "airport"
          "target_collection" = "airport"
        },
      ])
      "source_scope" = "inventory"
      "target_scope" = "inventory"
    },
  ])
  "network_usage_limit" = 0
  "organization_id" = "<organization_id>"
  "priority" = "high"
  "project_id" = "<project_id>"
  "reverse_replication_id" = tostring(null)
  "source_bucket" = "<source_bucket_id>"
  "source_details" = {
    "bucket" = {
      "conflict_resolution_type" = "seqno"
      "id" = "<source_bucket_id>"
      "name" = "xdcr_src"
    }
    "cluster" = {
      "id" = "<cluster_id>"
      "name" = "cluster-a"
    }
    "project" = {
      "id" = "<project_id>"
      "name" = "my-project"
    }
  }
  "status" = "running"
  "target" = {
    "bucket" = "<target_bucket_id>"
    "cluster" = "<target_cluster_id>"
    "type" = "capella"
  }
  "target_details" = {
    "bucket" = {
      "conflict_resolution_type" = "seqno"
      "id" = "<target_bucket_id>"
      "name" = "xdcr_tgt"
    }
    "cluster" = {
      "id" = "<target_cluster_id>"
      "name" = "cluster-b"
    }
    "project" = {
      "id" = "<project_id>"
      "name" = "my-project"
    }
    "type" = "capella"
  }
}
```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
couchbase-capella_replication.new_replication
```

## IMPORT
### Remove the resource `new_replication` from the Terraform State file

Command: `terraform state rm couchbase-capella_replication.new_replication`

Sample Output:
```
$ terraform state rm couchbase-capella_replication.new_replication
Removed couchbase-capella_replication.new_replication
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_replication.new_replication id=<replication_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
$ terraform import couchbase-capella_replication.new_replication id=<replication_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>
couchbase-capella_replication.new_replication: Importing from ID "id=<replication_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>"...
couchbase-capella_replication.new_replication: Import prepared!
  Prepared couchbase-capella_replication for import
couchbase-capella_replication.new_replication: Refreshing state... [id=id=<replication_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the replication ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the source cluster of the replication.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

After an import:

- `reverse_replication_id` is null. The API returns it only when the replication is created.
- `filter.expressions.skip_restream` is `false`. The API does not return it.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:
```
$ terraform plan
couchbase-capella_replication.new_replication: Refreshing state... [id=<replication_id>]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Update the replication settings by overriding the `replication` variable.

Command: `terraform apply -var 'replication={priority="low",network_usage_limit=100}'`

`priority`, `network_usage_limit`, `mappings` and `filter` are updated in place. For a `twoWay` replication, the change applies to both directions.

- Removing `mappings` replicates all scopes again.
- Removing `filter`, or a part of it, clears that part on the server.

A change to `source_bucket`, `target` or `direction` replaces the replication. When the last replication between two clusters is deleted, Capella waits up to 10 minutes before it creates a new replication between them. So a replacement of the only replication between two clusters can take up to 10 minutes longer.

Sample Output:
```
$ terraform apply -var 'replication={priority="low",network_usage_limit=100}'
couchbase-capella_replication.new_replication: Refreshing state... [id=<replication_id>]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_replication.new_replication will be updated in-place
  ~ resource "couchbase-capella_replication" "new_replication" {
      ~ changes_left           = 0 -> (known after apply)
      + error                  = (known after apply)
        id                     = "<replication_id>"
      ~ network_usage_limit    = 0 -> 100
      ~ priority               = "high" -> "low"
      + reverse_replication_id = (known after apply)
      ~ status                 = "running" -> (known after apply)
        # (11 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_replication = {
      ~ changes_left           = 0 -> (known after apply)
      + error                  = (known after apply)
        id                     = "<replication_id>"
      ~ network_usage_limit    = 0 -> 100
      ~ priority               = "high" -> "low"
      + reverse_replication_id = (known after apply)
      ~ status                 = "running" -> (known after apply)
        # (11 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_replication.new_replication: Modifying... [id=<replication_id>]
couchbase-capella_replication.new_replication: Still modifying... [id=<replication_id>, 00m10s elapsed]
couchbase-capella_replication.new_replication: Modifications complete after 11s [id=<replication_id>]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_replication = {
  "audit" = {
    "created_at" = "2026-10-13 22:26:18.630284798 +0000 UTC"
    "created_by" = "apikey-<api_key_id>"
  }
  "changes_left" = 0
  "cluster_id" = "<cluster_id>"
  "direction" = "oneWay"
  "error" = tostring(null)
  "filter" = {
    "document_exclude_options" = {
      "binary" = false
      "deletion" = false
      "expiration" = false
      "ttl" = true
    }
    "expressions" = {
      "reg_ex" = "REGEXP_CONTAINS(META().id, '^airline_')"
      "skip_restream" = false
    }
  }
  "id" = "<replication_id>"
  "mappings" = toset([
    {
      "collections" = toset([
        {
          "source_collection" = "airline"
          "target_collection" = "airline"
        },
        {
          "source_collection" = "airport"
          "target_collection" = "airport"
        },
      ])
      "source_scope" = "inventory"
      "target_scope" = "inventory"
    },
  ])
  "network_usage_limit" = 100
  "organization_id" = "<organization_id>"
  "priority" = "low"
  "project_id" = "<project_id>"
  "reverse_replication_id" = tostring(null)
  "source_bucket" = "<source_bucket_id>"
  "source_details" = {
    "bucket" = {
      "conflict_resolution_type" = "seqno"
      "id" = "<source_bucket_id>"
      "name" = "xdcr_src"
    }
    "cluster" = {
      "id" = "<cluster_id>"
      "name" = "cluster-a"
    }
    "project" = {
      "id" = "<project_id>"
      "name" = "my-project"
    }
  }
  "status" = "running"
  "target" = {
    "bucket" = "<target_bucket_id>"
    "cluster" = "<target_cluster_id>"
    "type" = "capella"
  }
  "target_details" = {
    "bucket" = {
      "conflict_resolution_type" = "seqno"
      "id" = "<target_bucket_id>"
      "name" = "xdcr_tgt"
    }
    "cluster" = {
      "id" = "<target_cluster_id>"
      "name" = "cluster-b"
    }
    "project" = {
      "id" = "<project_id>"
      "name" = "my-project"
    }
    "type" = "capella"
  }
}
replication_id = "<replication_id>"
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

For a `twoWay` replication, one delete removes both directions.

Sample Output:
```
$ terraform destroy
couchbase-capella_replication.new_replication: Refreshing state... [id=<replication_id>]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_replication.new_replication will be destroyed
  - resource "couchbase-capella_replication" "new_replication" {
      - audit               = {
          - created_at = "2026-10-13 22:26:18.630284798 +0000 UTC" -> null
          - created_by = "apikey-<api_key_id>" -> null
        } -> null
      - changes_left        = 0 -> null
      - cluster_id          = "<cluster_id>" -> null
      - direction           = "oneWay" -> null
      - filter              = {
          - document_exclude_options = {
              - binary     = false -> null
              - deletion   = false -> null
              - expiration = false -> null
              - ttl        = true -> null
            } -> null
          - expressions              = {
              - reg_ex        = "REGEXP_CONTAINS(META().id, '^airline_')" -> null
              - skip_restream = false -> null
            } -> null
        } -> null
      - id                  = "<replication_id>" -> null
      - mappings            = [
          - {
              - collections  = [
                  - {
                      - source_collection = "airline" -> null
                      - target_collection = "airline" -> null
                    },
                  - {
                      - source_collection = "airport" -> null
                      - target_collection = "airport" -> null
                    },
                ] -> null
              - source_scope = "inventory" -> null
              - target_scope = "inventory" -> null
            },
        ] -> null
      - network_usage_limit = 100 -> null
      - organization_id     = "<organization_id>" -> null
      - priority            = "low" -> null
      - project_id          = "<project_id>" -> null
      - source_bucket       = "<source_bucket_id>" -> null
      - source_details      = {
          - bucket  = {
              - conflict_resolution_type = "seqno" -> null
              - id                       = "<source_bucket_id>" -> null
              - name                     = "xdcr_src" -> null
            } -> null
          - cluster = {
              - id   = "<cluster_id>" -> null
              - name = "cluster-a" -> null
            } -> null
          - project = {
              - id   = "<project_id>" -> null
              - name = "my-project" -> null
            } -> null
        } -> null
      - status              = "running" -> null
      - target              = {
          - bucket  = "<target_bucket_id>" -> null
          - cluster = "<target_cluster_id>" -> null
          - type    = "capella" -> null
        } -> null
      - target_details      = {
          - bucket  = {
              - conflict_resolution_type = "seqno" -> null
              - id                       = "<target_bucket_id>" -> null
              - name                     = "xdcr_tgt" -> null
            } -> null
          - cluster = {
              - id   = "<target_cluster_id>" -> null
              - name = "cluster-b" -> null
            } -> null
          - project = {
              - id   = "<project_id>" -> null
              - name = "my-project" -> null
            } -> null
          - type    = "capella" -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_replication = {
      - audit                  = {
          - created_at = "2026-10-13 22:26:18.630284798 +0000 UTC"
          - created_by = "apikey-<api_key_id>"
        }
      - changes_left           = 0
      - cluster_id             = "<cluster_id>"
      - direction              = "oneWay"
      - error                  = null
      - filter                 = {
          - document_exclude_options = {
              - binary     = false
              - deletion   = false
              - expiration = false
              - ttl        = true
            }
          - expressions              = {
              - reg_ex        = "REGEXP_CONTAINS(META().id, '^airline_')"
              - skip_restream = false
            }
        }
      - id                     = "<replication_id>"
      - mappings               = [
          - {
              - collections  = [
                  - {
                      - source_collection = "airline"
                      - target_collection = "airline"
                    },
                  - {
                      - source_collection = "airport"
                      - target_collection = "airport"
                    },
                ]
              - source_scope = "inventory"
              - target_scope = "inventory"
            },
        ]
      - network_usage_limit    = 100
      - organization_id        = "<organization_id>"
      - priority               = "low"
      - project_id             = "<project_id>"
      - reverse_replication_id = null
      - source_bucket          = "<source_bucket_id>"
      - source_details         = {
          - bucket  = {
              - conflict_resolution_type = "seqno"
              - id                       = "<source_bucket_id>"
              - name                     = "xdcr_src"
            }
          - cluster = {
              - id   = "<cluster_id>"
              - name = "cluster-a"
            }
          - project = {
              - id   = "<project_id>"
              - name = "my-project"
            }
        }
      - status                 = "running"
      - target                 = {
          - bucket  = "<target_bucket_id>"
          - cluster = "<target_cluster_id>"
          - type    = "capella"
        }
      - target_details         = {
          - bucket  = {
              - conflict_resolution_type = "seqno"
              - id                       = "<target_bucket_id>"
              - name                     = "xdcr_tgt"
            }
          - cluster = {
              - id   = "<target_cluster_id>"
              - name = "cluster-b"
            }
          - project = {
              - id   = "<project_id>"
              - name = "my-project"
            }
          - type    = "capella"
        }
    } -> null
  - replication_id  = "<replication_id>" -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_replication.new_replication: Destroying... [id=<replication_id>]
couchbase-capella_replication.new_replication: Still destroying... [id=<replication_id>, 00m10s elapsed]
couchbase-capella_replication.new_replication: Still destroying... [id=<replication_id>, 00m20s elapsed]
couchbase-capella_replication.new_replication: Destruction complete after 27s

Destroy complete! Resources: 1 destroyed.
```
