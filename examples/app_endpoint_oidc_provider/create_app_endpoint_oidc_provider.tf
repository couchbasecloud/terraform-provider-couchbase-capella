resource "couchbase-capella_app_endpoint_oidc_provider" "example_oidc_provider" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  app_service_id    = var.app_service_id
  app_endpoint_name = var.app_endpoint_name

  issuer     = var.issuer
  client_id  = var.client_id

  discovery_url = var.discovery_url
  register      = var.register
  roles_claim   = var.roles_claim
  user_prefix   = var.user_prefix
  username_claim = var.username_claim
}
