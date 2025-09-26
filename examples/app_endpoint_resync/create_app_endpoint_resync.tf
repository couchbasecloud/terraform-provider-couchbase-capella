output "app_endpoint_resync" {
  value = couchbase-capella_app_endpoint_resync.this
}

resource "couchbase-capella_app_endpoint_resync" "this" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  app_service_id    = var.app_service_id
  app_endpoint      = var.app_endpoint_name

  # Optional: resync only specific scopes/collections
  # Set to null or omit to resync all
  scopes = var.scopes
}


