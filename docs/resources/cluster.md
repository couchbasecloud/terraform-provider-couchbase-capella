---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_cluster Resource - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  Manages the operational cluster resource.
---

# couchbase-capella_cluster (Resource)

Manages the operational cluster resource.

## Example Usage

```terraform
resource "couchbase-capella_cluster" "new_cluster" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  name            = "Terraform Test Cluster"
  description     = "Test cluster created with Terraform"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "10.200.250.0/23"
  }
  couchbase_server = {
    version = "7.6"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "io2"
          iops    = 5000
        }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `availability` (Attributes) Availability configuration for the cluster. (see [below for nested schema](#nestedatt--availability))
- `cloud_provider` (Attributes) The Cloud Service Provider where the cluster will be hosted. (see [below for nested schema](#nestedatt--cloud_provider))
- `name` (String) The name of the cluster (up to 256 characters).
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.
- `service_groups` (Attributes Set) Configuration for the Service Groups in the cluster. Each Service Group represents a set of nodes with the same configuration. (see [below for nested schema](#nestedatt--service_groups))
- `support` (Attributes) Support configuration for the cluster. (see [below for nested schema](#nestedatt--support))

### Optional

- `configuration_type` (String, Deprecated) The configuration type of the cluster. This field is deprecated.
- `couchbase_server` (Attributes) Configuration for the Couchbase Server running on the cluster. (see [below for nested schema](#nestedatt--couchbase_server))
- `description` (String) Description of the cluster (up to 1024 characters).
- `enable_private_dns_resolution` (Boolean) EnablePrivateDNSResolution signals that the cluster should have hostnames that are hosted in a public DNS zone that resolve to a private DNS address. This exists to support the use case of customers connecting from their own data centers where it is not possible to make use of a Cloud Service Provider DNS zone.
- `if_match` (String) The If-Match header value used for optimistic concurrency control.
- `zones` (Set of String) The Cloud Services Provider's availability zones for the cluster.For single availability zone clusters, only 1 zone is allowed in list.

### Read-Only

- `app_service_id` (String) The ID of the App Service associated with this cluster.
- `audit` (Attributes) Couchbase audit data. (see [below for nested schema](#nestedatt--audit))
- `connection_string` (String) The connection string to use to connect to the cluster.
- `current_state` (String) The current state of the cluster.
- `etag` (String) The ETag header value returned by the server, used for optimistic concurrency control.
- `id` (String) The ID of the operational cluster.

<a id="nestedatt--availability"></a>
### Nested Schema for `availability`

Required:

- `type` (String) The availability type of the cluster. Can be 'single' for Single Node or 'multi' for Multi Node.


<a id="nestedatt--cloud_provider"></a>
### Nested Schema for `cloud_provider`

Required:

- `cidr` (String) The CIDR block for the cluster's network.
- `region` (String) The region where the cluster will be hosted.
- `type` (String) The Cloud Service Provider type. Currently supporting AWS, GCP and Azure. For Single Node cluster, only the AWS type Cloud Service Provider is allowed.


<a id="nestedatt--service_groups"></a>
### Nested Schema for `service_groups`

Required:

- `node` (Attributes) Node configuration for this Service Group. (see [below for nested schema](#nestedatt--service_groups--node))
- `num_of_nodes` (Number) The number of nodes in this Service Group.
- `services` (Set of String) The list of Couchbase Services to run on the nodes in this Service Group.

<a id="nestedatt--service_groups--node"></a>
### Nested Schema for `service_groups.node`

Required:

- `compute` (Attributes) Compute resources configuration for the nodes. (see [below for nested schema](#nestedatt--service_groups--node--compute))
- `disk` (Attributes) The 'storage' and 'IOPS' fields are required for AWS. For Azure, only the 'disktype' field is required. For the Ultra disk type, you can provide storage, IOPS, and auto-expansion fields. For Premium type, you can only provide the auto-expansion field, others cannot be set. In the case of GCP, only 'pd ssd' disk type is available, and you cannot set the 'IOPS' field. (see [below for nested schema](#nestedatt--service_groups--node--disk))

<a id="nestedatt--service_groups--node--compute"></a>
### Nested Schema for `service_groups.node.compute`

Required:

- `cpu` (Number) The number of CPU cores for the node. The value must be between 1 and 128.
- `ram` (Number) The amount of RAM in MB for the node.


<a id="nestedatt--service_groups--node--disk"></a>
### Nested Schema for `service_groups.node.disk`

Required:

- `type` (String) The type of disk to use. For AWS: gp3 or io2, for Azure: Premium or UltraSSD, for GCP: pd-ssd.

Optional:

- `autoexpansion` (Boolean) Enable or disable automatic disk expansion.
- `iops` (Number) The number of IOPS for the disk. Only applicable for certain disk types.
- `storage` (Number) The size of the disk in GB.




<a id="nestedatt--support"></a>
### Nested Schema for `support`

Required:

- `plan` (String) The support plan options include 'Basic', 'Developer Pro', or 'Enterprise'.

Optional:

- `timezone` (String) The timezone for support coverage.


<a id="nestedatt--couchbase_server"></a>
### Nested Schema for `couchbase_server`

Optional:

- `version` (String) The version of Couchbase Server to run on the cluster.


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
terraform import couchbase-capella_cluster.new_cluster id=test_id,cluster_id=test_id,project_id=test_id,organization_id=test_id
```
