---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_backups Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  The backups data source retrieves backups associated with a bucket for an operational cluster.
---

# couchbase-capella_backups (Data Source)

The backups data source retrieves backups associated with a bucket for an operational cluster.

## Example Usage

```terraform
data "couchbase-capella_backups" "existing_backups" {
  organization_id = "aaaaa-bbbb-cccc-dddd-eeee"
  project_id      = "aaaaa-bbbb-cccc-dddd-eeee"
  cluster_id      = "aaaaa-bbbb-cccc-dddd-eeee"
  bucket_id       = "aaaaa-bbbb-cccc"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket_id` (String) The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name.
- `cluster_id` (String) The GUID4 ID of the cluster.
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.

### Read-Only

- `data` (Attributes List) Lists the backups associated with a bucket. (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `backup_stats` (Attributes) Lists the backup level data that Couchbase provides. (see [below for nested schema](#nestedatt--data--backup_stats))
- `bucket_id` (String) The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name.
- `bucket_name` (String) The name of the bucket for which the backup belongs to.
- `cloud_provider` (String) The Cloud Service Provider where the cluster is hosted.
- `cluster_id` (String) The GUID4 ID of the cluster.
- `cycle_id` (String) The GUID4 ID of the cycle this backup belongs to.
- `date` (String) The RFC3339 timestamp representing the time at which backup was created.
- `elapsed_time_in_seconds` (Number) The amount of seconds that have elapsed between the creation and completion of the backup.
- `id` (String) The GUID4 ID of the backup resource.
- `method` (String) The mechanism of the backup. It can be either incremental or full.
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.
- `restore_before` (String) The RFC3339 timestamp representing the time at which the backup will expire.
- `schedule_info` (Attributes) The schedule information of the backup. (see [below for nested schema](#nestedatt--data--schedule_info))
- `source` (String) Specifies whether the backup job was initiated manually or by a schedule.
- `status` (String) The status of the backup. Backup statuses are 'pending', 'ready', 'failed'.

<a id="nestedatt--data--backup_stats"></a>
### Nested Schema for `data.backup_stats`

Read-Only:

- `cbas` (Number) The number of analytics entities saved during the backup.
- `event` (Number) The number of eventing entities saved during the backup.
- `fts` (Number) The number of full text search entities saved during the backup.
- `gsi` (Number) The number of global secondary indexes saved during the backup.
- `items` (Number) The number of items saved during the backup.
- `mutations` (Number) The number of mutations saved during the backup.
- `size_in_mb` (Number) Size in MB is the total size of the backup in megabytes. It represents the amount of data that was backed up during the backup operation.
- `tombstones` (Number) The number of tombstones saved during the backup.


<a id="nestedatt--data--schedule_info"></a>
### Nested Schema for `data.schedule_info`

Read-Only:

- `backup_time` (String) The timestamp indicating when the backup was created.
- `backup_type` (String) Specifies if the backup frequency is daily or weekly.
- `increment` (Number) The interval in hours for incremental backups.
- `retention` (String) The retention time in days.
