---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_scope Resource - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  
---

# couchbase-capella_scope (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket_id` (String)
- `cluster_id` (String)
- `organization_id` (String)
- `project_id` (String)
- `scope_name` (String)

### Read-Only

- `collections` (Attributes Set) (see [below for nested schema](#nestedatt--collections))

<a id="nestedatt--collections"></a>
### Nested Schema for `collections`

Read-Only:

- `max_ttl` (Number)
- `name` (String)
