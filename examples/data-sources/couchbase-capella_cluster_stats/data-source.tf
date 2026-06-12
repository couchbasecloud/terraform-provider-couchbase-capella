data "couchbase-capella_cluster_stats" "existing_cluster_stats" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
}

output "cluster_stats" {
  value = data.couchbase-capella_cluster_stats.existing_cluster_stats
}
