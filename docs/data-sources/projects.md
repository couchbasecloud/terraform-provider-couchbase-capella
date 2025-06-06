---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_projects Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  Data source to retrieve project details in an organization.
---

# couchbase-capella_projects (Data Source)

Data source to retrieve project details in an organization.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `organization_id` (String) The GUID4 ID of the organization.

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `audit` (Attributes) Couchbase audit data. (see [below for nested schema](#nestedatt--data--audit))
- `description` (String) The description of a particular project.
- `etag` (String) The ETag header value returned by the server, used for optimistic concurrency control.
- `id` (String) The GUID4 ID of the project.
- `if_match` (String) A precondition header that specifies the entity tag of a resource.
- `name` (String) The name of the project.
- `organization_id` (String) The GUID4 ID of the organization.

<a id="nestedatt--data--audit"></a>
### Nested Schema for `data.audit`

Read-Only:

- `created_at` (String) The RFC3339 timestamp when the resource was created.
- `created_by` (String) The user who created the resource.
- `modified_at` (String) The RFC3339 timestamp when the resource was last modified.
- `modified_by` (String) The user who last modified the resource.
- `version` (Number) The version of the document. This value is incremented each time the resource is modified.
