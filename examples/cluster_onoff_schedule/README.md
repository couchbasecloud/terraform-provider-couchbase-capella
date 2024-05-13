# Capella Cluster OnOff Schedule Example

This example shows how to create and manage Cluster OnOff Schedule in Capella.

This creates a new cluster on/off schedule in the selected Capella cluster and lists existing cluster on/off schedules in the cluster. It uses the cluster ID to create and list cluster on/off schedules.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new cluster on/off schedule in Capella as stated in the `create_cluster.tf` file.
2. UPDATE: Update the cluster on/off schedule configuration using Terraform.
3. LIST: List existing cluster on/off schedules in Capella as stated in the `list_clusters.tf` file.
4. IMPORT: Import a cluster on/off schedule that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created cluster on/off schedule from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$  terraform plan 
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Reading...
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule will be created
  + resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
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
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
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
              + state = "off"
            },
          + {
              + day   = "sunday"
              + state = "off"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + timezone        = "US/Pacific"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

```

### Apply the Plan, in order to create a new Cluster On/Off schedule

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
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Reading...
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule will be created
  + resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
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
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
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
              + state = "off"
            },
          + {
              + day   = "sunday"
              + state = "off"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + timezone        = "US/Pacific"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Creating...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

existing_cluster_onoff_schedule = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "days" = tolist(null) /* of object */
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "timezone" = tostring(null)
}
new_cluster_onoff_schedule = {
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
        "minute" = 30
      }
      "state" = "custom"
      "to" = {
        "hour" = 14
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
      "from" = null /* object */
      "state" = "off"
      "to" = null /* object */
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

```

## UPDATE
### Let us edit the terraform.tfvars file to change the cluster on/off schedule configuration settings.

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
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Reading...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Refreshing state...
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule will be updated in-place
  ~ resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
      ~ days            = [
          ~ {
              ~ to    = {
                  ~ minute = 30 -> 0
                    # (1 unchanged attribute hidden)
                }
                # (3 unchanged attributes hidden)
            },
          ~ {
              ~ from  = {
                  ~ minute = 30 -> 0
                    # (1 unchanged attribute hidden)
                }
              ~ to    = {
                  ~ hour   = 14 -> 19
                    # (1 unchanged attribute hidden)
                }
                # (2 unchanged attributes hidden)
            },
            # (5 unchanged elements hidden)
        ]
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ existing_cluster_onoff_schedule = {
      ~ days            = null -> [
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
                  + minute = 30
                }
              + state = "custom"
              + to    = {
                  + hour   = 14
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
              + from  = null
              + state = "off"
              + to    = null
            },
          + {
              + day   = "sunday"
              + from  = null
              + state = "off"
              + to    = null
            },
        ]
      ~ timezone        = null -> "US/Pacific"
        # (3 unchanged attributes hidden)
    }
  ~ new_cluster_onoff_schedule      = {
      ~ days            = [
          ~ {
              ~ to    = {
                  ~ minute = 30 -> 0
                    # (1 unchanged attribute hidden)
                }
                # (3 unchanged attributes hidden)
            },
          ~ {
              ~ from  = {
                  ~ minute = 30 -> 0
                    # (1 unchanged attribute hidden)
                }
              ~ to    = {
                  ~ hour   = 14 -> 19
                    # (1 unchanged attribute hidden)
                }
                # (2 unchanged attributes hidden)
            },
            {
                day   = "wednesday"
                from  = null
                state = "on"
                to    = null
            },
            # (4 unchanged elements hidden)
        ]
        # (4 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Modifying...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Modifications complete after 0s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

existing_cluster_onoff_schedule = {
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
        "minute" = 30
      }
      "state" = "custom"
      "to" = {
        "hour" = 14
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
      "from" = null /* object */
      "state" = "off"
      "to" = null /* object */
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
new_cluster_onoff_schedule = {
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
        "minute" = 0
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
      "from" = null /* object */
      "state" = "off"
      "to" = null /* object */
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
```

### Note the output for the new Cluster On/Off schedule
Command: `terraform output new_cluster_onoff_schedule`

Sample Output:
```
$ terraform output new_cluster_onoff_schedule
{
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
        "minute" = 0
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
      "from" = null /* object */
      "state" = "off"
      "to" = null /* object */
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

```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule
```

## IMPORT
### Remove the resource `new_cluster_onoff_schedule` from the Terraform State file

Command: `terraform state rm couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule`

Sample Output:
```
$ terraform state rm couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule
Removed couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule
Successfully removed 1 resource instance(s).

```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule cluster_id=<cluster_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
$ terraform import couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Importing from ID "cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Reading...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Import prepared!
  Prepared couchbase-capella_cluster_onoff_schedule for import
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Refreshing state...
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Read complete after 0s
2024-03-27T01:00:01.981-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule during refresh.
      - .cluster_id: was cty.StringVal("cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"), but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .days: was null, but now cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"day":cty.StringVal("monday"), "from":cty.ObjectVal(map[string]cty.Value{"hour":cty.NumberIntVal(12), "minute":cty.NumberIntVal(30)}), "state":cty.StringVal("custom"), "to":cty.ObjectVal(map[string]cty.Value{"hour":cty.NumberIntVal(14), "minute":cty.NumberIntVal(0)})}), cty.ObjectVal(map[string]cty.Value{"day":cty.StringVal("tuesday"), "from":cty.ObjectVal(map[string]cty.Value{"hour":cty.NumberIntVal(12), "minute":cty.NumberIntVal(0)}), "state":cty.StringVal("custom"), "to":cty.ObjectVal(map[string]cty.Value{"hour":cty.NumberIntVal(19), "minute":cty.NumberIntVal(30)})}), cty.ObjectVal(map[string]cty.Value{"day":cty.StringVal("wednesday"), "from":cty.NullVal(cty.Object(map[string]cty.Type{"hour":cty.Number, "minute":cty.Number})), "state":cty.StringVal("on"), "to":cty.NullVal(cty.Object(map[string]cty.Type{"hour":cty.Number, "minute":cty.Number}))}), cty.ObjectVal(map[string]cty.Value{"day":cty.StringVal("thursday"), "from":cty.ObjectVal(map[string]cty.Value{"hour":cty.NumberIntVal(12), "minute":cty.NumberIntVal(30)}), "state":cty.StringVal("custom"), "to":cty.NullVal(cty.Object(map[string]cty.Type{"hour":cty.Number, "minute":cty.Number}))}), cty.ObjectVal(map[string]cty.Value{"day":cty.StringVal("friday"), "from":cty.ObjectVal(map[string]cty.Value{"hour":cty.NumberIntVal(0), "minute":cty.NumberIntVal(0)}), "state":cty.StringVal("custom"), "to":cty.ObjectVal(map[string]cty.Value{"hour":cty.NumberIntVal(12), "minute":cty.NumberIntVal(30)})}), cty.ObjectVal(map[string]cty.Value{"day":cty.StringVal("saturday"), "from":cty.NullVal(cty.Object(map[string]cty.Type{"hour":cty.Number, "minute":cty.Number})), "state":cty.StringVal("off"), "to":cty.NullVal(cty.Object(map[string]cty.Type{"hour":cty.Number, "minute":cty.Number}))}), cty.ObjectVal(map[string]cty.Value{"day":cty.StringVal("sunday"), "from":cty.NullVal(cty.Object(map[string]cty.Type{"hour":cty.Number, "minute":cty.Number})), "state":cty.StringVal("off"), "to":cty.NullVal(cty.Object(map[string]cty.Type{"hour":cty.Number, "minute":cty.Number}))})})
      - .organization_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .project_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .timezone: was null, but now cty.StringVal("US/Pacific")

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.


```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the cluster ID i.e. the ID of the cluster schedule that we want to import.
The second ID is the project ID i.e. the ID of the project to which the cluster belongs.
The third ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

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
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Reading...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Refreshing state...
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Read complete after 1s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

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
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Reading...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Refreshing state...
data.couchbase-capella_cluster_onoff_schedule.existing_cluster_onoff_schedule: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

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
                  - minute = 0 -> null
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
              - state = "off" -> null
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

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - existing_cluster_onoff_schedule = {
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
                  - minute = 0
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
              - from  = null
              - state = "off"
              - to    = null
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
  - new_cluster_onoff_schedule      = {
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
                  - minute = 0
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
              - from  = null
              - state = "off"
              - to    = null
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

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Destroying...
couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.

````