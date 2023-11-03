output "new_backup" {
  value = capella_backup.new_backup
}

resource "capella_backup" "new_backup" {
  organization_id            = var.organization_id
  project_id                 = var.project_id
  cluster_id                 = var.cluster_id
  bucket_id                  = var.bucket_id

#  ToDo Required for Backup Schedule, tracking under -https://couchbasecloud.atlassian.net/browse/AV-66698
#  type = var.backup.type
#  weekly_schedule = {
#    day_of_week = var.weekly_schedule.day_of_week
#    start_at = var.weekly_schedule.start_at
#    incremental_every = var.weekly_schedule.incremental_every
#    retention_time = var.weekly_schedule.retention_time
#    cost_optimized_retention = var.weekly_schedule.cost_optimized_retention
#  }
}