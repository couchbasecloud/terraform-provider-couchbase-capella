output "backups_list" {
  value = data.couchbase-capella_backups.existing_backups
}

data "couchbase-capella_backups" "existing_backups" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
}