output "new_app_endpoint" {
  value = couchbase-capella_app_endpoint.new_app_endpoint
}

output "app_endpoint_id" {
  value = couchbase-capella_app_endpoint.new_app_endpoint.id
}

resource "couchbase-capella_app_endpoint" "new_app_endpoint" {
  organization_id    = var.app_endpoint.organization_id
  project_id         = var.app_endpoint.project_id
  cluster_id         = var.app_endpoint.cluster_id
  app_service_id     = var.app_endpoint.app_service_id
  bucket             = var.app_endpoint.bucket
  name               = var.app_endpoint.name
  user_xattr_key     = var.app_endpoint.user_xattr_key
  delta_sync_enabled = var.app_endpoint.delta_sync_enabled

  scopes = var.app_endpoint.scopes

  cors = var.app_endpoint.cors

  oidc = var.app_endpoint.oidc
}