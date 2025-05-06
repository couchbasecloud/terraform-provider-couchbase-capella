resource "couchbase-capella_backup_schedule" "new_backup_schedule" {
  organization_id = "organization_id"
  project_id      = "project_id"
  cluster_id      = "cluster_id"
  bucket_id       = "bucket_id"
  type            = "backup_schedule_type"
  weekly_schedule = {
    day_of_week              = "day_of_week"
    start_at                 = "start_time"
    incremental_every        = "incremental_hours"
    retention_time           = "retention_time"
    cost_optimized_retention = "cost_optimized_retention(bool)"
  }
}