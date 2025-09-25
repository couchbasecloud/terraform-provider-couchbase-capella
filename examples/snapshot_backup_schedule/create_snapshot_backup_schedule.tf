output "new_snapshot_backup_schedule" {
  value = couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule
}

resource "couchbase-capella_snapshot_backup_schedule" "new_snapshot_backup_schedule" {
  organization_id = var.organization_id
  project_id = var.project_id
  id = var.cluster_id

  interval   = var.snapshot_backup_schedule.interval
  retention  = var.snapshot_backup_schedule.retention
  start_time = var.snapshot_backup_schedule.start_time
  copy_to_regions = var.snapshot_backup_schedule.copy_to_regions
}