output "sample_buckets_list" {
  value = data.couchbase-capella_sample_buckets.existing_sample_buckets
}

data "couchbase-capella_sample_buckets" "existing_sample_buckets" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
