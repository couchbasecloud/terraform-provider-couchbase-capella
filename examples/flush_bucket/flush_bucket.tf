output "new_flush" {
  value = couchbase-capella_flush.new_flush
}

resource "couchbase-capella_flush" "new_flush" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
}

