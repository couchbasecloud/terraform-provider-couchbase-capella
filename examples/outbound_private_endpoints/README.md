# List Capella Outbound Private Endpoints Data source

This example shows how to list the outbound private endpoints on a Capella cluster: the endpoints in the cluster's network that connect to private link services you own. It uses the organization ID, project ID and cluster ID to do so.

The list is newest first and has no paging, since a cluster holds at most its endpoint limit. Pending, failed and deleting endpoints are included. `provider_config.aws` is only set on AWS clusters.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## List
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_outbound_private_endpoints.existing_outbound_private_endpoints: Reading...
data.couchbase-capella_outbound_private_endpoints.existing_outbound_private_endpoints: Read complete after 1s

Changes to Outputs:
  + existing_outbound_private_endpoints = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + audit                = {
                  + created_at  = "2026-09-30 10:00:00 +0000 UTC"
                  + created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + modified_at = "2026-09-30 11:20:05 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
                  + version     = 4
                }
              + cloud_provider       = {
                  + region = "us-east-1"
                  + type   = "aws"
                }
              + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + endpoint_type        = "interface"
              + failure_reason       = null
              + hostname             = "cng.acme.example.com"
              + id                   = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
              + ip_addresses         = [
                  + "10.0.4.21",
                ]
              + is_custom_hostname   = true
              + name                 = "cng-prod"
              + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
              + private_dns_enabled  = true
              + project_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + provider_config      = {
                  + aws = {
                      + availability_zones = [
                          + "us-east-1a",
                        ]
                      + ip_addresses       = [
                          + {
                              + availability_zone = "us-east-1a"
                              + ip_address        = "10.0.4.21"
                            },
                        ]
                      + regional_dns_name  = "vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com"
                    }
                }
              + provider_endpoint_id = "vpce-0a1b2c3d4e5f67890"
              + provider_status      = "available"
              + service_name         = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
              + status               = "linked"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.
```

### Apply the Plan

Command: `terraform apply`

Sample Output:

```
$ terraform apply
data.couchbase-capella_outbound_private_endpoints.existing_outbound_private_endpoints: Reading...
data.couchbase-capella_outbound_private_endpoints.existing_outbound_private_endpoints: Read complete after 1s

Changes to Outputs:
  + existing_outbound_private_endpoints = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + audit                = {
                  + created_at  = "2026-09-30 10:00:00 +0000 UTC"
                  + created_by  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + modified_at = "2026-09-30 11:20:05 +0000 UTC"
                  + modified_by = "ffffffff-aaaa-1414-eeee-000000000000"
                  + version     = 4
                }
              + cloud_provider       = {
                  + region = "us-east-1"
                  + type   = "aws"
                }
              + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + endpoint_type        = "interface"
              + failure_reason       = null
              + hostname             = "cng.acme.example.com"
              + id                   = "3f2e1d0c-9b8a-4765-8432-10fedcba9876"
              + ip_addresses         = [
                  + "10.0.4.21",
                ]
              + is_custom_hostname   = true
              + name                 = "cng-prod"
              + organization_id      = "ffffffff-aaaa-1414-eeee-000000000000"
              + private_dns_enabled  = true
              + project_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + provider_config      = {
                  + aws = {
                      + availability_zones = [
                          + "us-east-1a",
                        ]
                      + ip_addresses       = [
                          + {
                              + availability_zone = "us-east-1a"
                              + ip_address        = "10.0.4.21"
                            },
                        ]
                      + regional_dns_name  = "vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com"
                    }
                }
              + provider_endpoint_id = "vpce-0a1b2c3d4e5f67890"
              + provider_status      = "available"
              + service_name         = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
              + status               = "linked"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```
