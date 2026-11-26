# Verify an Outbound Private Endpoint Service Data source

This example shows how to check a private link service before creating an outbound private endpoint to it. Capella checks the service from the cluster's cloud account and region: the service name's format and region, whether Capella can see the service, its zones, whether the cluster already has an endpoint to it, and on AWS its private DNS name. Nothing is created.

A failed check doesn't fail the read. It shows in `checks` with `ready` set to false, so a configuration can gate the endpoint on `ready`. `provider_config` is only set once the service is found, and only AWS reports fields today.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
$ terraform plan
data.couchbase-capella_outbound_private_endpoint_service_verification.service_check: Reading...
data.couchbase-capella_outbound_private_endpoint_service_verification.service_check: Read complete after 1s

Changes to Outputs:
  + service_check = {
      + checks          = [
          + {
              + message = null
              + name    = "format"
              + result  = "passed"
            },
          + {
              + message = null
              + name    = "region"
              + result  = "passed"
            },
          + {
              + message = null
              + name    = "accessGranted"
              + result  = "passed"
            },
          + {
              + message = null
              + name    = "serviceZones"
              + result  = "passed"
            },
          + {
              + message = null
              + name    = "notAlreadyConnected"
              + result  = "passed"
            },
          + {
              + message = null
              + name    = "privateDnsName"
              + result  = "passed"
            },
        ]
      + cloud_provider  = {
          + region = "us-east-1"
          + type   = "aws"
        }
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + provider_config = {
          + aws = {
              + acceptance_required           = true
              + availability_zones            = [
                  + "us-east-1a",
                  + "us-east-1b",
                  + "us-east-1c",
                ]
              + private_dns_name              = "*.acme.example.com"
              + private_dns_name_verification = "verified"
              + service_owner                 = "210987654321"
            }
        }
      + ready           = true
      + service_name    = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
    }

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.
```

### A service that isn't ready

On Azure, an alias the Capella subscription can't see fails `accessGranted`. The read still succeeds:

```
      + checks          = [
          ...
          + {
              + message = "Add the Capella subscription to the service's visibility."
              + name    = "accessGranted"
              + result  = "failed"
            },
          ...
        ]
      + provider_config = null
      + ready           = false
```
