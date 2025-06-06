---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_backup Resource - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  Manages backup resource associated with a bucket for an operational Capella cluster.
---

# couchbase-capella_backup (Resource)

Manages backup resource associated with a bucket for an operational Capella cluster.

## Example Usage

```terraform
resource "couchbase-capella_backup" "new_backup" {
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

### Optional

- `restore` (Attributes) Represents restore configuration for the backup. (see [below for nested schema](#nestedatt--restore))
- `restore_before` (String) The RFC3339 timestamp representing the time at which backup will expire.
- `restore_times` (Number) Number of times the backup to be restored.

### Read-Only

- `backup_stats` (Attributes) Represents various backup level data that couchbase provides. (see [below for nested schema](#nestedatt--backup_stats))
- `bucket_name` (String) The name of the bucket for which the backup belongs to.
- `cloud_provider` (String) The Cloud Service Provider where the cluster is hosted.
- `cycle_id` (String) The GUID4 ID of the cycle this backup belongs to.
- `date` (String) The RFC3339 timestamp representing the time at which backup was created.
- `elapsed_time_in_seconds` (Number) The amount of seconds that have elapsed between the creation and completion of the backup.
- `id` (String) The GUID4 ID of the backup resource.
- `method` (String) The mechanism of the backup. It can be either incremental or full.
- `schedule_info` (Attributes) Represents the schedule information of the backup. (see [below for nested schema](#nestedatt--schedule_info))
- `source` (String) The source of the backup. It can be either cluster or bucket.
- `status` (String) The status of the backup. Backup statuses are 'pending', 'ready', and 'failed'.

<a id="nestedatt--restore"></a>
### Nested Schema for `restore`

Required:

- `services` (List of String) Services is a list of services to be restored. It is used to specify the services that should be included in the restore operation.
- `source_cluster_id` (String) The ID of the source cluster the restore is based on.
- `target_cluster_id` (String) The ID of the target cluster to restore to.

Optional:

- `auto_remove_collections` (Boolean) Automatically remove collections that are not present in the backup.
- `exclude_data` (String) Skips restoring the data specified here.
- `filter_keys` (String) Only restore data where the key matches a particular regular expression.
- `filter_values` (String) Only restore data where the value matches a particular regular expression.
- `force_updates` (Boolean) Forces data in the operational cluster to be overwritten even if the data in the cluster is newer.
- `include_data` (String) Restores only the data specified here.
- `map_data` (String) Specified when you want to restore source data into a different location.
- `replace_ttl` (String) Sets a new expiration (time-to-live) value for the specified keys. These values are 'none', 'all', and 'expired'.
- `replace_ttl_with` (String) Updates the expiration for the keys.

Read-Only:

- `status` (String) The status of the restore.


<a id="nestedatt--backup_stats"></a>
### Nested Schema for `backup_stats`

Read-Only:

- `cbas` (Number) The number of analytics entities saved during the backup.
- `event` (Number) The number of event entities saved during the backup.
- `fts` (Number) The number of full text search entities saved during the backup.
- `gsi` (Number) The number of global secondary indexes saved during the backup.
- `items` (Number) The number of items saved during the backup.
- `mutations` (Number) The number of mutations saved during the backup.
- `size_in_mb` (Number) Size in MB is the total size of the backup in megabytes. It represents the amount of data that was backed up during the backup operation.
- `tombstones` (Number) The number of tombstones saved during the backup.


<a id="nestedatt--schedule_info"></a>
### Nested Schema for `schedule_info`

Read-Only:

- `backup_time` (String) Represents the time at which the backup is scheduled to be taken.
- `backup_type` (String) Represents whether the backup is a weekly or daily backup.
- `increment` (Number) Represents interval in hours for incremental backup.
- `retention` (String) Represents retention time in days.

## Import

Import is supported using the following syntax:

```shell
terraform import couchbase-capella_backup.new_backup id=<backup_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>
```
