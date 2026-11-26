# Capella Outbound Private Endpoint Example

This example shows how to create and manage an outbound private endpoint in Capella: an endpoint in a cluster's network that connects to a private link service you own, such as a Cloud Native Gateway in front of a self-managed cluster. Use its `hostname` as the target of an XDCR replication.

The endpoint service must allow Capella's AWS account. Get it from the `couchbase-capella_cloud_accounts` data source. The cluster must run Couchbase Server 8.5 or later.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create an outbound private endpoint as stated in the `create_outbound_private_endpoint.tf` file.
2. ACCEPT: Accept the connection on the endpoint service in your own AWS account, then refresh.
3. IMPORT: Import an outbound private endpoint that exists in Capella but not in the terraform state file.
4. UPDATE: Rename the endpoint. Changing any other field replaces it.
5. DELETE: Delete the endpoint.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access. Leave `hostname` out to have Capella assign one.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint will be created
  + resource "couchbase-capella_outbound_private_endpoint" "new_outbound_private_endpoint" {
      + audit                = (known after apply)
      + cloud_provider       = (known after apply)
      + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
      + endpoint_type        = "interface"
      + failure_reason       = (known after apply)
      + hostname             = "cng.acme.example.com"
      + id                   = (known after apply)
      + ip_addresses         = (known after apply)
      + is_custom_hostname   = (known after apply)
      + name                 = "cng-prod"
      + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + private_dns_enabled  = (known after apply)
      + project_id           = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config      = (known after apply)
      + provider_endpoint_id = (known after apply)
      + provider_status      = (known after apply)
      + service_name         = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
      + status               = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_outbound_private_endpoint = {
      + audit                = (known after apply)
      + cloud_provider       = (known after apply)
      + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
      + endpoint_type        = "interface"
      + failure_reason       = (known after apply)
      + hostname             = "cng.acme.example.com"
      + id                   = (known after apply)
      + ip_addresses         = (known after apply)
      + is_custom_hostname   = (known after apply)
      + name                 = "cng-prod"
      + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + private_dns_enabled  = (known after apply)
      + project_id           = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config      = (known after apply)
      + provider_endpoint_id = (known after apply)
      + provider_status      = (known after apply)
      + service_name         = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
      + status               = (known after apply)
    }
  + outbound_private_endpoint_id  = (known after apply)

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new outbound private endpoint

Create returns once the connection request reaches the endpoint service. `status` stays `pendingAcceptance` until you accept it.

Command: `terraform apply`

Sample Output:
```
$ terraform apply
...
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Creating...
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Creation complete after 0s [id=3f2e1d0c-9b8a-4765-8432-10fedcba9876]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_outbound_private_endpoint = {
  "audit" = {
    "created_at" = "2026-11-24 10:00:00 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2026-11-24 10:00:00 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 1
  }
  "cloud_provider" = {
    "region" = "us-east-1"
    "type" = "aws"
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "endpoint_type" = "interface"
  "failure_reason" = tostring(null)
  "hostname" = "cng.acme.example.com"
  "id" = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
  "ip_addresses" = tolist([])
  "is_custom_hostname" = true
  "name" = "cng-prod"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "private_dns_enabled" = false
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws" = {
      "availability_zones" = tolist([
        "us-east-1a",
        "us-east-1b",
      ])
      "ip_addresses" = tolist([])
      "regional_dns_name" = "vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com"
    }
  }
  "provider_endpoint_id" = "vpce-0a1b2c3d4e5f67890"
  "provider_status" = "pendingAcceptance"
  "service_name" = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
  "status" = "pendingAcceptance"
}
outbound_private_endpoint_id = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint
```

## ACCEPT
### Accept the connection in your AWS account

Accept the endpoint connection on your endpoint service, unless the service accepts connections on its own. The next plan or refresh reads `linked`, the endpoint's addresses and, with your own hostname, private DNS.

Command: `terraform plan`

Sample Output:
```
$ terraform plan
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Refreshing state... [id=3f2e1d0c-9b8a-4765-8432-10fedcba9876]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the
last "terraform apply" which may have affected this plan:

  # couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint has changed
  ~ resource "couchbase-capella_outbound_private_endpoint" "new_outbound_private_endpoint" {
      ~ audit                = {
          ~ modified_at = "2026-11-24 10:00:00 +0000 UTC" -> "2026-11-24 10:12:41 +0000 UTC"
          ~ version     = 1 -> 2
            # (3 unchanged attributes hidden)
        }
        id                   = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
      ~ ip_addresses         = [
          + "10.0.4.21",
          + "10.0.5.37",
        ]
        name                 = "cng-prod"
      ~ private_dns_enabled  = false -> true
      ~ provider_config      = {
          ~ aws = {
              ~ ip_addresses       = [
                  + {
                      + availability_zone = "us-east-1a"
                      + ip_address        = "10.0.4.21"
                    },
                  + {
                      + availability_zone = "us-east-1b"
                      + ip_address        = "10.0.5.37"
                    },
                ]
                # (2 unchanged attributes hidden)
            }
        }
      ~ provider_status      = "pendingAcceptance" -> "available"
      ~ status               = "pendingAcceptance" -> "linked"
        # (9 unchanged attributes hidden)
    }

Unless you have made equivalent changes to your configuration, or ignored the
relevant attributes using ignore_changes, the following plan may include
actions to undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────

Changes to Outputs:
  ~ new_outbound_private_endpoint = {
      ~ audit                = {
          ~ modified_at = "2026-11-24 10:00:00 +0000 UTC" -> "2026-11-24 10:12:41 +0000 UTC"
          ~ version     = 1 -> 2
            # (3 unchanged attributes hidden)
        }
        id                   = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
      ~ ip_addresses         = [
          + "10.0.4.21",
          + "10.0.5.37",
        ]
        name                 = "cng-prod"
      ~ private_dns_enabled  = false -> true
      ~ provider_config      = {
          ~ aws = {
              ~ ip_addresses       = [
                  + {
                      + availability_zone = "us-east-1a"
                      + ip_address        = "10.0.4.21"
                    },
                  + {
                      + availability_zone = "us-east-1b"
                      + ip_address        = "10.0.5.37"
                    },
                ]
                # (2 unchanged attributes hidden)
            }
        }
      ~ provider_status      = "pendingAcceptance" -> "available"
      ~ status               = "pendingAcceptance" -> "linked"
        # (10 unchanged attributes hidden)
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
```

## IMPORT
### Remove the resource `new_outbound_private_endpoint` from the Terraform State file

Command: `terraform state rm couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint`

Sample Output:
```
$ terraform state rm couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint
Removed couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint
Successfully removed 1 resource instance(s).
```

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint id=<outbound_private_endpoint_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
$ terraform import couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint id=3f2e1d0c-9b8a-4765-8432-10fedcba9876,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Importing from ID "id=3f2e1d0c-9b8a-4765-8432-10fedcba9876,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Import prepared!
  Prepared couchbase-capella_outbound_private_endpoint for import
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Refreshing state... [id=id=3f2e1d0c-9b8a-4765-8432-10fedcba9876,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:
```
$ terraform plan
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Refreshing state... [id=3f2e1d0c-9b8a-4765-8432-10fedcba9876]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to rename the endpoint.

Only `name` updates in place. Changing `service_name`, `endpoint_type` or `hostname` replaces the endpoint.

Command: `terraform apply`

Sample Output:
```
$ terraform apply
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Refreshing state... [id=3f2e1d0c-9b8a-4765-8432-10fedcba9876]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint will be updated in-place
  ~ resource "couchbase-capella_outbound_private_endpoint" "new_outbound_private_endpoint" {
      ~ audit                = {
          ~ created_at  = "2026-11-24 10:00:00 +0000 UTC" -> (known after apply)
          ~ created_by  = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
          ~ modified_at = "2026-11-24 10:12:41 +0000 UTC" -> (known after apply)
          ~ modified_by = "ffffffff-aaaa-1414-eeee-000000000000" -> (known after apply)
          ~ version     = 2 -> (known after apply)
        } -> (known after apply)
      + failure_reason       = (known after apply)
        id                   = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
      ~ ip_addresses         = [
          - "10.0.4.21",
          - "10.0.5.37",
        ] -> (known after apply)
      ~ name                 = "cng-prod" -> "cng-prod-east"
      ~ private_dns_enabled  = true -> (known after apply)
      ~ provider_config      = {
          ~ aws = {
              ~ availability_zones = [
                  - "us-east-1a",
                  - "us-east-1b",
                ] -> (known after apply)
              ~ ip_addresses       = [
                  - {
                      - availability_zone = "us-east-1a" -> null
                      - ip_address        = "10.0.4.21" -> null
                    },
                  - {
                      - availability_zone = "us-east-1b" -> null
                      - ip_address        = "10.0.5.37" -> null
                    },
                ] -> (known after apply)
              ~ regional_dns_name  = "vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com" -> (known after apply)
            } -> (known after apply)
        } -> (known after apply)
      ~ provider_status      = "available" -> (known after apply)
      ~ status               = "linked" -> (known after apply)
        # (9 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_outbound_private_endpoint = {
      ~ audit                = {
          - created_at  = "2026-11-24 10:00:00 +0000 UTC"
          - created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
          - modified_at = "2026-11-24 10:12:41 +0000 UTC"
          - modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
          - version     = 2
        } -> (known after apply)
      + failure_reason       = (known after apply)
        id                   = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
      ~ ip_addresses         = [
          - "10.0.4.21",
          - "10.0.5.37",
        ] -> (known after apply)
      ~ name                 = "cng-prod" -> "cng-prod-east"
      ~ private_dns_enabled  = true -> (known after apply)
      ~ provider_config      = {
          - aws = {
              - availability_zones = [
                  - "us-east-1a",
                  - "us-east-1b",
                ]
              - ip_addresses       = [
                  - {
                      - availability_zone = "us-east-1a"
                      - ip_address        = "10.0.4.21"
                    },
                  - {
                      - availability_zone = "us-east-1b"
                      - ip_address        = "10.0.5.37"
                    },
                ]
              - regional_dns_name  = "vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com"
            }
        } -> (known after apply)
      ~ provider_status      = "available" -> (known after apply)
      ~ status               = "linked" -> (known after apply)
        # (9 unchanged attributes hidden)
    }
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Modifying... [id=3f2e1d0c-9b8a-4765-8432-10fedcba9876]
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Modifications complete after 0s [id=3f2e1d0c-9b8a-4765-8432-10fedcba9876]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_outbound_private_endpoint = {
  "audit" = {
    "created_at" = "2026-11-24 10:00:00 +0000 UTC"
    "created_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "modified_at" = "2026-11-24 10:20:03 +0000 UTC"
    "modified_by" = "ffffffff-aaaa-1414-eeee-000000000000"
    "version" = 3
  }
  "cloud_provider" = {
    "region" = "us-east-1"
    "type" = "aws"
  }
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "endpoint_type" = "interface"
  "failure_reason" = tostring(null)
  "hostname" = "cng.acme.example.com"
  "id" = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
  "ip_addresses" = tolist([
    "10.0.4.21",
    "10.0.5.37",
  ])
  "is_custom_hostname" = true
  "name" = "cng-prod-east"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "private_dns_enabled" = true
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "provider_config" = {
    "aws" = {
      "availability_zones" = tolist([
        "us-east-1a",
        "us-east-1b",
      ])
      "ip_addresses" = tolist([
        {
          "availability_zone" = "us-east-1a"
          "ip_address" = "10.0.4.21"
        },
        {
          "availability_zone" = "us-east-1b"
          "ip_address" = "10.0.5.37"
        },
      ])
      "regional_dns_name" = "vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com"
    }
  }
  "provider_endpoint_id" = "vpce-0a1b2c3d4e5f67890"
  "provider_status" = "available"
  "service_name" = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
  "status" = "linked"
}
outbound_private_endpoint_id = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
```

## DESTROY
### Finally, destroy the resources created by Terraform

Destroy waits until the endpoint is gone.

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
...
Plan: 0 to add, 0 to change, 1 to destroy.
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Destroying... [id=3f2e1d0c-9b8a-4765-8432-10fedcba9876]
couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
