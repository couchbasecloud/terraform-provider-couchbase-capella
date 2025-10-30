output "new_cloud_snapshot_backup" {
  value = couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup
}

resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  retention       = var.cloud_snapshot_backup.retention
  regions_to_copy = var.cloud_snapshot_backup.regions_to_copy
}