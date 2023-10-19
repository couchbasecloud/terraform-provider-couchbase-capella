output "clusters_list" {
  value = data.capella_clusters.existing_clusters
}

data "capella_clusters" "existing_clusters" {
  organization_id = var.organization_id
  project_id      = var.project_id
}
