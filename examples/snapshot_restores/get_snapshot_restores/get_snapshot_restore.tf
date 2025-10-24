output "cloud_snapshot_restore" {
  value = data.couchbase-capella_cloud_snapshot_restore.cloud_snapshot_restore
}

data "couchbase-capella_cloud_snapshot_restore" "cloud_snapshot_restore" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  id              = var.id
}