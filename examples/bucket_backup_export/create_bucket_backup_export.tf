output "new_bucket_backup_export" {
  value     = couchbase-capella_bucket_backup_export.new_bucket_backup_export
  sensitive = true
}

resource "couchbase-capella_bucket_backup_export" "new_bucket_backup_export" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
  backup_id       = var.backup_id
}
