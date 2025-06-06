---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_users Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  The data source to retrieve users in a Capella organization.
---

# couchbase-capella_users (Data Source)

The data source to retrieve users in a Capella organization.



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
- `email` (String) The email of the user.
- `enable_notifications` (Boolean) After enabling email notifications for your account, you will start receiving email notification alerts from all databases in projects you are a part of.
- `expires_at` (String) Time at which the user expires.
- `id` (String) The ID of the user.
- `inactive` (Boolean) Inactive depicts whether the user has accepted the invite for the organization.
- `last_login` (String) Time(UTC) when the user last logged in.
- `name` (String) The name of the user.
- `organization_id` (String) The GUID4 ID of the organization.
- `organization_roles` (List of String) The organization roles associated with the user. They determine the privileges a user possesses in the organization.
- `region` (String) The region of the user.
- `resources` (Attributes List) (see [below for nested schema](#nestedatt--data--resources))
- `status` (String) Status depicts user status whether they are verified or not. It can be one of the following values: verified, not-verified, pending-primary.
- `time_zone` (String) Time zone of the user.

<a id="nestedatt--data--audit"></a>
### Nested Schema for `data.audit`

Read-Only:

- `created_at` (String) The RFC3339 timestamp when the resource was created.
- `created_by` (String) The user who created the resource.
- `modified_at` (String) The RFC3339 timestamp when the resource was last modified.
- `modified_by` (String) The user who last modified the resource.
- `version` (Number) The version of the document. This value is incremented each time the resource is modified.


<a id="nestedatt--data--resources"></a>
### Nested Schema for `data.resources`

Optional:

- `type` (String) Type of the resource.

Read-Only:

- `id` (String) The GUID4 ID of the project.
- `roles` (List of String) Project Roles associated with the user.
