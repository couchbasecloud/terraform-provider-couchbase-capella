---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_network_peers Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  
---

# couchbase-capella_network_peers (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_id` (String)
- `organization_id` (String)
- `project_id` (String)

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `audit` (Attributes) (see [below for nested schema](#nestedatt--data--audit))
- `id` (String)
- `name` (String)
- `provider_config` (Attributes) (see [below for nested schema](#nestedatt--data--provider_config))
- `status` (Attributes) (see [below for nested schema](#nestedatt--data--status))

<a id="nestedatt--data--audit"></a>
### Nested Schema for `data.audit`

Read-Only:

- `created_at` (String)
- `created_by` (String)
- `modified_at` (String)
- `modified_by` (String)
- `version` (Number)


<a id="nestedatt--data--provider_config"></a>
### Nested Schema for `data.provider_config`

Read-Only:

- `aws_config` (Attributes) (see [below for nested schema](#nestedatt--data--provider_config--aws_config))
- `gcp_config` (Attributes) (see [below for nested schema](#nestedatt--data--provider_config--gcp_config))

<a id="nestedatt--data--provider_config--aws_config"></a>
### Nested Schema for `data.provider_config.aws_config`

Read-Only:

- `account_id` (String)
- `cidr` (String)
- `provider_id` (String)
- `region` (String)
- `vpc_id` (String)


<a id="nestedatt--data--provider_config--gcp_config"></a>
### Nested Schema for `data.provider_config.gcp_config`

Read-Only:

- `cidr` (String)
- `network_name` (String)
- `project_id` (String)
- `provider_id` (String)
- `service_account` (String)



<a id="nestedatt--data--status"></a>
### Nested Schema for `data.status`

Read-Only:

- `reasoning` (String)
- `state` (String)