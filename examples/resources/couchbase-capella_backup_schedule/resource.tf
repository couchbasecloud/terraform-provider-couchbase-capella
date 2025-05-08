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