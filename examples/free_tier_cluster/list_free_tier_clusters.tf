output "free_tier_clusters_list" {
  value = data.couchbase-capella_free_tier_clusters.existing_free_tier_clusters
}

data "couchbase-capella_free_tier_clusters" "existing_free_tier_clusters" {
  organization_id = var.organization_id
  project_id      = var.project_id
}
