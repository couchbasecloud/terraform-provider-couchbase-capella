output "backups_list" {
  value = data.couchbase-capella_cloud_snapshot_backups.existing_cloud_snapshot_backups
}

data "couchbase-capella_cloud_snapshot_backups" "existing_cloud_snapshot_backups" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  filter {
    name = "status"
    values = ["complete"]
  }
}
