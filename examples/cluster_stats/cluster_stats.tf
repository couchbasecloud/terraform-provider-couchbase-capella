data "couchbase-capella_cluster_stats" "existing_cluster_stats" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}

output "cluster_stats" {
  value = data.couchbase-capella_cluster_stats.existing_cluster_stats
}
