output "buckets_list" {
  value = data.couchbase-capella_buckets.existing_buckets
}

data "couchbase-capella_buckets" "existing_buckets" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
