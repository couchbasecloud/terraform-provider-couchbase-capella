---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_app_service Resource - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  This resource allows you to create and manage an App Service in Capella. App Service is a fully managed application backend designed to provide data synchronization between mobile or IoT applications running Couchbase Lite and your Couchbase Capella database.
---

# couchbase-capella_app_service (Resource)

This resource allows you to create and manage an App Service in Capella. App Service is a fully managed application backend designed to provide data synchronization between mobile or IoT applications running Couchbase Lite and your Couchbase Capella database.

## Example Usage

```terraform
resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  name            = "MyAppSyncService"
  description     = "My app sync service."
  nodes           = 2
  compute = {
    cpu = 2
    ram = 4
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_id` (String) The GUID4 ID of the cluster.
- `compute` (Attributes) The CPU and RAM configuration of the App Service. (see [below for nested schema](#nestedatt--compute))
- `name` (String) Name of the cluster (up to 256 characters).
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.

### Optional

- `cloud_provider` (String) The Cloud Service Provider for the App Service.
- `description` (String) A short description of the App Service.
- `if_match` (String) A precondition header that specifies the entity tag of a resource.
- `nodes` (Number) Number of nodes configured for the App Service. Number of nodes configured for the App Service. The number of nodes can range from 2 to 12.

### Read-Only

- `audit` (Attributes) Couchbase audit data. (see [below for nested schema](#nestedatt--audit))
- `current_state` (String) The current state of the App Service.
- `etag` (String)
- `id` (String) The ID of the App Service created.
- `version` (String) The version of the App Service server. If left empty, it will be defaulted to the latest available version.

<a id="nestedatt--compute"></a>
### Nested Schema for `compute`

Required:

- `cpu` (Number) CPU units (cores).
- `ram` (Number) RAM units (GB).


<a id="nestedatt--audit"></a>
### Nested Schema for `audit`

Read-Only:

- `created_at` (String) The RFC3339 timestamp when the resource was created.
- `created_by` (String) The user who created the resource.
- `modified_at` (String) The RFC3339 timestamp when the resource was last modified.
- `modified_by` (String) The user who last modified the resource.
- `version` (Number) The version of the document. This value is incremented each time the resource is modified.

## Import

Import is supported using the following syntax:

```shell
terraform import couchbase-capella_app_service.new_app_service id=<appservice_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>
```
