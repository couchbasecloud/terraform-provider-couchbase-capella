output "existing_cloud_snapshot_backup_schedule" {
  value = data.couchbase-capella_cloud_snapshot_backup_schedule.existing_cloud_snapshot_backup_schedule
}

data "couchbase-capella_cloud_snapshot_backup_schedule" "existing_cloud_snapshot_backup_schedule" {
  organization_id = var.organization_id
  project_id = var.project_id
  id = var.cluster_id
}