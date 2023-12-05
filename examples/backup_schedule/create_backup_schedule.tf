output "new_backup_schedule" {
  value = capella_backup_schedule.new_backup_schedule
}

resource "capella_backup_schedule" "new_backup_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
  type            = var.backup_schedule.type
  weekly_schedule = {
    day_of_week              = var.backup_schedule.weekly_schedule.day_of_week
    start_at                 = var.backup_schedule.weekly_schedule.start_at
    incremental_every        = var.backup_schedule.weekly_schedule.incremental_every
    retention_time           = var.backup_schedule.weekly_schedule.retention_time
    cost_optimized_retention = var.backup_schedule.weekly_schedule.cost_optimized_retention
  }
}