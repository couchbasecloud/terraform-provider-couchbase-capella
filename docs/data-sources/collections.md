---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_collections Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  
---

# couchbase-capella_collections (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket_id` (String)
- `cluster_id` (String)
- `organization_id` (String)
- `project_id` (String)
- `scope_name` (String)

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `collection_name` (String)
- `max_ttl` (Number)
