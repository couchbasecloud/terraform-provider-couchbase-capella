---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_audit_log_settings Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  
---

# couchbase-capella_audit_log_settings (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_id` (String)
- `organization_id` (String)
- `project_id` (String)

### Read-Only

- `audit_enabled` (Boolean)
- `disabled_users` (Attributes Set) (see [below for nested schema](#nestedatt--disabled_users))
- `enabled_event_ids` (Set of Number)

<a id="nestedatt--disabled_users"></a>
### Nested Schema for `disabled_users`

Read-Only:

- `domain` (String)
- `name` (String)
