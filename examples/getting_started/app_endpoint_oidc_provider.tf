resource "couchbase-capella_app_endpoint_oidc_provider" "oidc1" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  app_service_id    = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = var.app_endpoint

  issuer        = "https://accounts.google.com"
  client_id     = "example-client-id"
  discovery_url = "https://accounts.google.com/.well-known/openid-configuration"
  register      = false
  username_claim = "sub"
  roles_claim    = "roles"
  user_prefix    = "user_"
}


