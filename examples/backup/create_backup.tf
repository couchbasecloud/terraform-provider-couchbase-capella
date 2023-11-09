output "new_backup" {
  value = capella_backup.new_backup
}

resource "capella_backup" "new_backup" {
  organization_id            = var.organization_id
  project_id                 = var.project_id
  cluster_id                 = var.cluster_id
  bucket_id                  = var.bucket_id
}