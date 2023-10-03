---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "capella_organization Data Source - terraform-provider-capella"
subcategory: ""
description: |-
  
---

# capella_organization (Data Source)





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
- `description` (String)
- `id` (String)
- `name` (String)
- `preferences` (Attributes) (see [below for nested schema](#nestedatt--data--preferences))

<a id="nestedatt--data--audit"></a>
### Nested Schema for `data.audit`

Read-Only:

- `created_at` (String)
- `created_by` (String)
- `modified_at` (String)
- `modified_by` (String)
- `version` (Number)


<a id="nestedatt--data--preferences"></a>
### Nested Schema for `data.preferences`

Read-Only:

- `session_duration` (Number)