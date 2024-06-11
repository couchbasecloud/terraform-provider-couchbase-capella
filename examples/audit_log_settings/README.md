# Capella Audit Log Settings Example

This example shows how to create and manage audit log settings in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1.  CREATE: Create a new audit log settings in Capella per the `create_auditlogsettings.tf` file.
2.  UPDATE: Update the audit log settings configuration using Terraform.
3.  Detect the resource was updated from outside of Terraform
4.  IMPORT: Import audit log settings that exists in Capella but not in the terraform state file.
5.  DELETE:  It is not supported by API server so it's just a noop (ie removes resource from state file).

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Reading...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be created
  + resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = []
      + enabled_event_ids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_auditlogsettings = {
      + audit_enabled     = false
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = []
      + enabled_event_ids = [
          + 8243,
          + 8257,
          + 8265,
          + 20480,
          + 20482,
          + 20483,
          + 20492,
          + 20494,
          + 32768,
          + 32769,
          + 32773,
          + 32774,
          + 32775,
          + 32776,
          + 32777,
          + 32778,
          + 32779,
          + 32781,
          + 32784,
          + 32789,
          + 32790,
          + 32791,
          + 32792,
          + 32794,
          + 32797,
          + 36865,
          + 36866,
          + 40960,
          + 40961,
          + 40962,
          + 40964,
          + 40966,
          + 45056,
          + 45058,
          + 45059,
          + 45060,
          + 45062,
          + 45063,
          + 45064,
          + 45065,
          + 45067,
          + 45068,
          + 45069,
          + 45071,
          + 45072,
          + 45073,
          + 45074,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_auditlogsettings      = {
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = []
      + enabled_event_ids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply"
now.
```

### Apply the Plan, in order to create the audit log setting

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Reading...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be created
  + resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = []
      + enabled_event_ids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_auditlogsettings = {
      + audit_enabled     = false
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = []
      + enabled_event_ids = [
          + 8243,
          + 8257,
          + 8265,
          + 20480,
          + 20482,
          + 20483,
          + 20492,
          + 20494,
          + 32768,
          + 32769,
          + 32773,
          + 32774,
          + 32775,
          + 32776,
          + 32777,
          + 32778,
          + 32779,
          + 32781,
          + 32784,
          + 32789,
          + 32790,
          + 32791,
          + 32792,
          + 32794,
          + 32797,
          + 36865,
          + 36866,
          + 40960,
          + 40961,
          + 40962,
          + 40964,
          + 40966,
          + 45056,
          + 45058,
          + 45059,
          + 45060,
          + 45062,
          + 45063,
          + 45064,
          + 45065,
          + 45067,
          + 45068,
          + 45069,
          + 45071,
          + 45072,
          + 45073,
          + 45074,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_auditlogsettings      = {
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = []
      + enabled_event_ids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_audit_log_settings.new_auditlogsettings: Creating...
couchbase-capella_audit_log_settings.new_auditlogsettings: Creation complete after 3s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

existing_auditlogsettings = {
  "audit_enabled" = false
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = toset([])
  "enabled_event_ids" = toset([
    8243,
    8257,
    8265,
    20480,
    20482,
    20483,
    20492,
    20494,
    32768,
    32769,
    32773,
    32774,
    32775,
    32776,
    32777,
    32778,
    32779,
    32781,
    32784,
    32789,
    32790,
    32791,
    32792,
    32794,
    32797,
    36865,
    36866,
    40960,
    40961,
    40962,
    40964,
    40966,
    45056,
    45058,
    45059,
    45060,
    45062,
    45063,
    45064,
    45065,
    45067,
    45068,
    45069,
    45071,
    45072,
    45073,
    45074,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_auditlogsettings = {
  "audit_enabled" = true
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = tolist([])
  "enabled_event_ids" = tolist([
    20488,
    20490,
    20491,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}

```

### Update resource by setting audit_enabled to false

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Reading...
couchbase-capella_audit_log_settings.new_auditlogsettings: Refreshing state...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be updated in-place
  ~ resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      ~ audit_enabled     = true -> false
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ existing_auditlogsettings = {
      ~ audit_enabled     = false -> true
      ~ enabled_event_ids = [
          - 8243,
          - 8257,
          - 8265,
          - 20480,
          - 20482,
          - 20483,
          - 20492,
          - 20494,
          - 32768,
          - 32769,
          - 32773,
          - 32774,
          - 32775,
          - 32776,
          - 32777,
          - 32778,
          - 32779,
          - 32781,
          - 32784,
          - 32789,
          - 32790,
          - 32791,
          - 32792,
          - 32794,
          - 32797,
          - 36865,
          - 36866,
          - 40960,
          - 40961,
          - 40962,
          - 40964,
          - 40966,
          - 45056,
          - 45058,
          - 45059,
          - 45060,
          - 45062,
          - 45063,
          - 45064,
          - 45065,
          - 45067,
          - 45068,
          - 45069,
          - 45071,
          - 45072,
          - 45073,
          - 45074,
          + 20488,
          + 20490,
          + 20491,
        ]
        # (4 unchanged attributes hidden)
    }
  ~ new_auditlogsettings      = {
      ~ audit_enabled     = true -> false
        # (5 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_audit_log_settings.new_auditlogsettings: Modifying...
couchbase-capella_audit_log_settings.new_auditlogsettings: Modifications complete after 2s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

existing_auditlogsettings = {
  "audit_enabled" = true
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = toset([])
  "enabled_event_ids" = toset([
    20488,
    20490,
    20491,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_auditlogsettings = {
  "audit_enabled" = false
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = tolist([])
  "enabled_event_ids" = tolist([
    20488,
    20490,
    20491,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

### Detect that the resource was modified outside of Terraform
### In this example audit_enabled was set to true outside of Terraform

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Reading...
couchbase-capella_audit_log_settings.new_auditlogsettings: Refreshing state...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Read complete after 1s
2024-03-12T20:26:23.721-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_audit_log_settings.new_auditlogsettings during refresh.
      - .audit_enabled: was cty.False, but now cty.True

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be updated in-place
  ~ resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      ~ audit_enabled     = true -> false
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

## IMPORT
### Remove the resource from the Terraform State file

Command: `terraform state rm couchbase-capella_audit_log_settings.new_auditlogsettings`

Sample Output:
```
terraform state rm couchbase-capella_audit_log_settings.new_auditlogsettings
Removed couchbase-capella_audit_log_settings.new_auditlogsettings
Successfully removed 1 resource instance(s).
```

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_audit_log_settings.new_auditlogsettings id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
terraform import couchbase-capella_audit_log_settings.new_auditlogsettings id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_audit_log_settings.new_auditlogsettings: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Reading...
couchbase-capella_audit_log_settings.new_auditlogsettings: Import prepared!
  Prepared couchbase-capella_audit_log_settings for import
couchbase-capella_audit_log_settings.new_auditlogsettings: Refreshing state...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Read complete after 1s
2024-03-12T20:35:56.149-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_audit_log_settings.new_auditlogsettings during refresh.
      - .enabled_event_ids: was null, but now cty.ListVal([]cty.Value{cty.NumberIntVal(20488), cty.NumberIntVal(20490), cty.NumberIntVal(20491)})
      - .organization_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .project_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .audit_enabled: was null, but now cty.True
      - .cluster_id: was cty.StringVal("id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"), but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .disabled_users: was null, but now cty.ListValEmpty(cty.Object(map[string]cty.Type{"domain":cty.String, "name":cty.String}))

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the cluster ID.
The second ID is the project ID i.e. the ID of the project to which the cluster belongs.
The third ID is the organization ID i.e. the ID of the organization to which the project belongs.

### This example shows how to configure a disabled user
### We use an internal user @eventing in local domain
### In terraform.tfvars set disabled_users = [{"domain": "local", "name": "@eventing"}]

Command: `terraform apply`

Sample Output:
```

terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Reading...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be created
  + resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = [
          + {
              + domain = "local"
              + name   = "@eventing"
            },
        ]
      + enabled_event_ids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + existing_auditlogsettings = {
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = []
      + enabled_event_ids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }
  + new_auditlogsettings      = {
      + audit_enabled     = true
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + disabled_users    = [
          + {
              + domain = "local"
              + name   = "@eventing"
            },
        ]
      + enabled_event_ids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_audit_log_settings.new_auditlogsettings: Creating...
couchbase-capella_audit_log_settings.new_auditlogsettings: Creation complete after 3s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

existing_auditlogsettings = {
  "audit_enabled" = true
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = toset([])
  "enabled_event_ids" = toset([
    20488,
    20490,
    20491,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
new_auditlogsettings = {
  "audit_enabled" = true
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "disabled_users" = toset([
    {
      "domain" = "local"
      "name" = "@eventing"
    },
  ])
  "enabled_event_ids" = toset([
    20488,
    20490,
    20491,
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

## DESTROY
### This only removes resource from the state file.

Command: `terraform destroy`

Sample Output:
```
terraform destroy
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Reading...
couchbase-capella_audit_log_settings.new_auditlogsettings: Refreshing state...
data.couchbase-capella_audit_log_settings.existing_auditlogsettings: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_audit_log_settings.new_auditlogsettings will be destroyed
  - resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
      - audit_enabled     = true -> null
      - cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - disabled_users    = [] -> null
      - enabled_event_ids = [
          - 20488,
          - 20490,
          - 20491,
        ] -> null
      - organization_id   = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - existing_auditlogsettings = {
      - audit_enabled     = true
      - cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      - disabled_users    = []
      - enabled_event_ids = [
          - 20488,
          - 20490,
          - 20491,
        ]
      - organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_audit_log_settings.new_auditlogsettings: Destroying...
couchbase-capella_audit_log_settings.new_auditlogsettings: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```