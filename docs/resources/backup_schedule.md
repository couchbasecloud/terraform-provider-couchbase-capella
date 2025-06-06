---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_backup_schedule Resource - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  Manages the backup schedule resource associated with a bucket for an operational cluster.
---

# couchbase-capella_backup_schedule (Resource)

Manages the backup schedule resource associated with a bucket for an operational cluster.

## Example Usage

```terraform
resource "couchbase-capella_backup_schedule" "new_backup_schedule" {
  organization_id = "aaaaa-bbbbb-ccccc-dddddd"
  project_id      = "aaaaa-bbbbb-ccccc-dddddd"
  cluster_id      = "aaaaa-bbbbb-ccccc-dddddd"
  bucket_id       = "aaaaa-bbbbbbb"
  type            = "weekly"
  weekly_schedule = {
    day_of_week              = "sunday"
    start_at                 = 10
    incremental_every        = 4
    retention_time           = "90days"
    cost_optimized_retention = false
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket_id` (String) The GUID4 ID of the bucket.
- `cluster_id` (String) The GUID4 ID of the cluster.
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.
- `type` (String) Type of the backup schedule.
- `weekly_schedule` (Attributes) Schedule a full backup once a week with regular incrementals. (see [below for nested schema](#nestedatt--weekly_schedule))

<a id="nestedatt--weekly_schedule"></a>
### Nested Schema for `weekly_schedule`

Required:

- `cost_optimized_retention` (Boolean) Optimize backup retention to reduce total cost of ownership (TCO). This gives the option to keep all but the last backup cycle of the month for thirty days; the last cycle will be kept for the defined retention period.
- `day_of_week` (String) Day of the week for the backup. Values can be "sunday", "monday", "tuesday", "wednesday", "thursday", "friday", or "saturday"
- `incremental_every` (Number) Interval in hours for incremental backup. Integer value between 1 and 24.
- `retention_time` (String) Retention time in days. For example: 30days, 1year, 5years.
- `start_at` (Number) The starting hour (in 24-Hour format). Integer value between 0 and 23.

## Import

Import is supported using the following syntax:

```shell
terraform import couchbase-capella_backup.new_backup id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>
```
