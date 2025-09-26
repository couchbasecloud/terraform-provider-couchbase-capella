# Capella App Endpoint Example

This example shows how to create and manage App Endpoints in Capella.

This creates a new app endpoint in the selected Capella organization. App endpoints provide sync gateways that enable real-time data synchronization between Couchbase Server and mobile/web applications.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new App Endpoint in Capella as stated in the `create_app_endpoint.tf` file.
2. UPDATE: Update the App Endpoint using Terraform.
3. DELETE: Delete the App Endpoint from Capella.
4. IMPORT: Import an App Endpoint that exists in Capella but not in the terraform state file.

If you check the `terraform.template.tfvars` file - you can see that we need several variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE

Command: `terraform apply`

Sample Output:
```
$ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint.endpoint2 will be created
  + resource "couchbase-capella_app_endpoint" "endpoint2" {
      + admin_url          = (known after apply)
      + app_service_id     = "<app_service_id>"
      + bucket             = "b1"
      + cluster_id         = "<cluster_id>"
      + delta_sync_enabled = false
      + metrics_url        = (known after apply)
      + name               = "test-endpoint-1-cors"
      + organization_id    = "<org_id>"
      + project_id         = "<project_id>"
      + public_url         = (known after apply)
      + require_resync     = (known after apply)
      + scopes             = {
          + "s1" = {
              + collections = {
                  + "c1" = {
                      + access_control_function = (known after apply)
                      + import_filter           = (known after apply)
                    },
                },
            },
        }
      + state              = (known after apply)
      + user_xattr_key     = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint.endpoint2: Creating...
couchbase-capella_app_endpoint.endpoint2: Creation complete after 1s [name=test-endpoint-1-cors]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## UPDATE
### Now update the App Endpoint to add CORS configuration

Command: `terraform apply`

Sample Output:
```
$ terraform apply

couchbase-capella_app_endpoint.endpoint2: Refreshing state... [name=test-endpoint-1-cors]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint.endpoint2 will be updated in-place
  ~ resource "couchbase-capella_app_endpoint" "endpoint2" {
      + cors               = {
          + disabled = false
          + headers  = [
              + "Authorization",
              + "Content-Type",
            ]
          + max_age  = 3600
          + origin   = [
              + "*",
            ]
        }
        name               = "test-endpoint-1-cors"
      + require_resync     = (known after apply)
        # (11 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint.endpoint2: Modifying...
couchbase-capella_app_endpoint.endpoint2: Modifications complete after 1s [name=test-endpoint-1-cors]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy

couchbase-capella_app_endpoint.endpoint2: Refreshing state... [name=test-endpoint-1-cors]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint.endpoint2 will be destroyed
  - resource "couchbase-capella_app_endpoint" "endpoint2" {
      - admin_url          = "https://a44mhfab-t17gzhg.apps.nonprod-project-avengers.com:4985/test-endpoint-1-cors" -> null
      - app_service_id     = <app_service_id> -> null
      - bucket             = "b1" -> null
      - cluster_id         = "<cluster_id>" -> null
      - cors               = {
          - disabled = false -> null
          - headers  = [
              - "Authorization",
              - "Content-Type",
            ] -> null
          - max_age  = 3600 -> null
          - origin   = [
              - "*",
            ] -> null
        } -> null
      - delta_sync_enabled = false -> null
      - metrics_url        = "https://a44mhfab-t17gzhg.apps.nonprod-project-avengers.com:4988/metrics" -> null
      - name               = "test-endpoint-1-cors" -> null
      - organization_id    = "<org_id>" -> null
      - project_id         = "<project_id>" -> null
      - public_url         = "wss://a44mhfab-t17gzhg.apps.nonprod-project-avengers.com:4984/test-endpoint-1-cors" -> null
      - scopes             = {
          - "s1" = {
              - collections = {
                  - "c1" = {
                      - access_control_function = "function (doc, oldDoc, meta) {channel('c1');}" -> null
                      - import_filter           = "" -> null
                    },
                } -> null
            },
        } -> null
      - state              = "Offline" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_endpoint.endpoint2: Destroying...
couchbase-capella_app_endpoint.endpoint2: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```

## IMPORT

Command: `terraform import couchbase-capella_app_endpoint.endpoint1 name=test-endpoint-1,organization_id=<org_id>,project_id=<project_id>,app_service_id=<app_service_id>,cluster_id=<cluster_id>`

Sample Output:
```
$ terraform import couchbase-capella_app_endpoint.endpoint1 name=test-endpoint-1,organization_id=<org_id>,project_id=<project_id>,app_service_id=dffe37cc-07b9-4095-a8e0-70594b8456ea,cluster_id=<cluster_id>
couchbase-capella_app_endpoint.endpoint1: Importing from ID "name=test-endpoint-1,organization_id=<org_id>,project_id=<project_id>,app_service_id=dffe37cc-07b9-4095-a8e0-70594b8456ea,cluster_id=<cluster_id>"...
couchbase-capella_app_endpoint.endpoint1: Import prepared!
  Prepared couchbase-capella_app_endpoint for import
couchbase-capella_app_endpoint.endpoint1: Refreshing state... [name=name=test-endpoint-1,organization_id=<org_id>,project_id=<project_id>,app_service_id=dffe37cc-07b9-4095-a8e0-70594b8456ea,cluster_id=<cluster_id>]
2025-09-09T14:43:55.539-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_app_endpoint.endpoint1 during refresh.
      - .scopes: was null, but now cty.MapVal(map[string]cty.Value{"s1":cty.ObjectVal(map[string]cty.Value{"collections":cty.MapVal(map[string]cty.Value{"c1":cty.ObjectVal(map[string]cty.Value{"access_control_function":cty.StringVal("function (doc, oldDoc, meta) {channel('c1');}"), "import_filter":cty.StringVal("function (doc) {\n  if (doc.type != \"mobile\") {\n    return false;\n  }\n  return true;\n}")}), "c2":cty.ObjectVal(map[string]cty.Value{"access_control_function":cty.StringVal("function (doc, oldDoc, meta) {channel('c2');}"), "import_filter":cty.StringVal("function (doc) {\n  if (doc.type != \"mobile\") {\n    return false;\n  }\n  return true;\n}")})})})})
      - .user_xattr_key: was null, but now cty.StringVal("")
      - .admin_url: was null, but now cty.StringVal("https://a44mhfab-t17gzhg.apps.nonprod-project-avengers.com:4985/test-endpoint-1")
      - .organization_id: was null, but now cty.StringVal("<org_id>")
      - .public_url: was null, but now cty.StringVal("wss://a44mhfab-t17gzhg.apps.nonprod-project-avengers.com:4984/test-endpoint-1")
      - .app_service_id: was null, but now cty.StringVal(<app_service_id>)
      - .cors: was null, but now cty.ObjectVal(map[string]cty.Value{"disabled":cty.False, "headers":cty.SetVal([]cty.Value{cty.StringVal("")}), "login_origin":cty.NullVal(cty.Set(cty.String)), "max_age":cty.NumberIntVal(5), "origin":cty.SetVal([]cty.Value{cty.StringVal("*")})})
      - .state: was null, but now cty.StringVal("Offline")
      - .cluster_id: was null, but now cty.StringVal("<cluster_id>")
      - .name: was cty.StringVal("name=test-endpoint-1,organization_id=<org_id>,project_id=<project_id>,app_service_id=dffe37cc-07b9-4095-a8e0-70594b8456ea,cluster_id=<cluster_id>"), but now cty.StringVal("test-endpoint-1")
      - .project_id: was null, but now cty.StringVal("<project_id>")
      - .bucket: was null, but now cty.StringVal("b1")
      - .delta_sync_enabled: was null, but now cty.False
      - .metrics_url: was null, but now cty.StringVal("https://a44mhfab-t17gzhg.apps.nonprod-project-avengers.com:4988/metrics")

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

