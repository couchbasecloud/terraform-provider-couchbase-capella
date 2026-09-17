data "couchbase-capella_bucket_backup_export" "existing_bucket_backup_export" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  bucket_id       = "<bucket_id>"
  backup_id       = "<backup_id>"
  id              = "<export_id>"
}
