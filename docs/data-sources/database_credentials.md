---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_database_credentials Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  The data source to retrieve database credentials for a cluster. Database credentials provide programmatic and application-level access to data on a database.
---

# couchbase-capella_database_credentials (Data Source)

The data source to retrieve database credentials for a cluster. Database credentials provide programmatic and application-level access to data on a database.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_id` (String) The GUID4 ID of the cluster.
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Optional:

- `access` (Attributes List) Describes the access information of the database credential. (see [below for nested schema](#nestedatt--data--access))

Read-Only:

- `audit` (Attributes) Couchbase audit data. (see [below for nested schema](#nestedatt--data--audit))
- `cluster_id` (String) The GUID4 ID of the cluster.
- `id` (String) The ID of the database credential created.
- `name` (String) Name of the database credential created (up to 256 characters).
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.

<a id="nestedatt--data--access"></a>
### Nested Schema for `data.access`

Required:

- `privileges` (List of String) The privileges field in this API represents the privilege level for users.

Optional:

- `resources` (Attributes) The resources for which access will be granted on. Leaving this empty will grant access to all buckets. (see [below for nested schema](#nestedatt--data--access--resources))

<a id="nestedatt--data--access--resources"></a>
### Nested Schema for `data.access.resources`

Optional:

- `buckets` (Attributes List) (see [below for nested schema](#nestedatt--data--access--resources--buckets))

<a id="nestedatt--data--access--resources--buckets"></a>
### Nested Schema for `data.access.resources.buckets`

Required:

- `name` (String) The name of the bucket.

Optional:

- `scopes` (Attributes List) The scopes under a bucket. (see [below for nested schema](#nestedatt--data--access--resources--buckets--scopes))

<a id="nestedatt--data--access--resources--buckets--scopes"></a>
### Nested Schema for `data.access.resources.buckets.scopes`

Required:

- `name` (String) The name of the scope.

Optional:

- `collections` (List of String) The collections under a scope.





<a id="nestedatt--data--audit"></a>
### Nested Schema for `data.audit`

Read-Only:

- `created_at` (String) The RFC3339 timestamp when the resource was created.
- `created_by` (String) The user who created the resource.
- `modified_at` (String) The RFC3339 timestamp when the resource was last modified.
- `modified_by` (String) The user who last modified the resource.
- `version` (Number) The version of the document. This value is incremented each time the resource is modified.
