resource "couchbase-capella_app_endpoint" "new_app_endpoint" {
  organization_id    = var.organization_id
  project_id         = var.project_id
  cluster_id         = var.cluster_id
  app_service_id     = var.app_service_id
  bucket             = var.bucket
  name               = var.name
  delta_sync_enabled = var.delta_sync_enabled
  scope = var.scope
  collections        = var.collections
  cors = var.cors
  oidc = var.oidc
}