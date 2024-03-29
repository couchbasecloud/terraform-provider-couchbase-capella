---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_app_services Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  
---

# couchbase-capella_app_services (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `organization_id` (String)

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `audit` (Attributes) (see [below for nested schema](#nestedatt--data--audit))
- `cloud_provider` (String)
- `cluster_id` (String)
- `compute` (Attributes) (see [below for nested schema](#nestedatt--data--compute))
- `current_state` (String)
- `description` (String)
- `id` (String)
- `name` (String)
- `nodes` (Number)
- `organization_id` (String)
- `version` (String)

<a id="nestedatt--data--audit"></a>
### Nested Schema for `data.audit`

Read-Only:

- `created_at` (String)
- `created_by` (String)
- `modified_at` (String)
- `modified_by` (String)
- `version` (Number)


<a id="nestedatt--data--compute"></a>
### Nested Schema for `data.compute`

Read-Only:

- `cpu` (Number)
- `ram` (Number)
