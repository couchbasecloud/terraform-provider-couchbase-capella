output "existing_bucket_backup_export" {
  value     = data.couchbase-capella_bucket_backup_export.existing_bucket_backup_export
  sensitive = true
}

output "download_url" {
  value     = data.couchbase-capella_bucket_backup_export.existing_bucket_backup_export.backup_download_url
  sensitive = true
}

data "couchbase-capella_bucket_backup_export" "existing_bucket_backup_export" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
  backup_id       = var.backup_id
  id              = couchbase-capella_bucket_backup_export.new_bucket_backup_export.id
}
