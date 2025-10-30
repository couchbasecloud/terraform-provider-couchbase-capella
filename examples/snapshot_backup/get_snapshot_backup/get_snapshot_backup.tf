output "cloud_snapshot_backup" {
  value = data.couchbase-capella_cloud_snapshot_backup.cloud_snapshot_backup
}

data "couchbase-capella_cloud_snapshot_backup" "cloud_snapshot_backup" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  id              = var.id
}