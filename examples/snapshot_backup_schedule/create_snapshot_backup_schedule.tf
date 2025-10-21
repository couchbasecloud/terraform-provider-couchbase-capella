output "new_cloud_snapshot_backup_schedule" {
  value = couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule
}

resource "couchbase-capella_cloud_snapshot_backup_schedule" "new_cloud_snapshot_backup_schedule" {
  organization_id = var.organization_id
  project_id = var.project_id
  cluster_id = var.cluster_id

  interval   = var.cloud_snapshot_backup_schedule.interval
  retention  = var.cloud_snapshot_backup_schedule.retention
  start_time = var.cloud_snapshot_backup_schedule.start_time
  copy_to_regions = var.cloud_snapshot_backup_schedule.copy_to_regions
}