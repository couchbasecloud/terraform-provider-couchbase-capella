output "samplebuckets_list" {
  value = data.couchbase-capella_samplebuckets.existing_samplebuckets
}

data "couchbase-capella_samplebuckets" "existing_samplebuckets" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
