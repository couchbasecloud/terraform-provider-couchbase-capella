output "backups_list" {
  value = data.couchbase-capella_backups.existing_backups
}

data "couchbase-capella_buckets" "existing_buckets" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}

data "couchbase-capella_backups" "existing_backups" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = element(data.couchbase-capella_buckets.existing_buckets.data, 0).id
}
