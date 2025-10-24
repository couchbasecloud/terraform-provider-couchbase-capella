output "cloud_snapshot_restores_list" {
  value = data.couchbase-capella_cloud_snapshot_restores.existing_cloud_snapshot_restores
}

data "couchbase-capella_cloud_snapshot_restores" "existing_cloud_snapshot_restores" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  
  filter {
    name = "status"
    values = ["complete"]
  }
}