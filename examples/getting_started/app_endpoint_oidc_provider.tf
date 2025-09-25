# Stores the App Endpoint OIDC Provider details in an output variable.
# Can be viewed using `terraform output app_endpoint_oidc_provider` command
output "app_endpoint_oidc_provider" {
  value = couchbase-capella_app_endpoint_oidc_provider.oidc
}

resource "couchbase-capella_app_endpoint_oidc_provider" "oidc" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  app_service_id    = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = couchbase-capella_app_endpoint.new_app_endpoint.name

  issuer    = var.app_endpoint_oidc.issuer
  client_id = var.app_endpoint_oidc.client_id

  discovery_url  = var.app_endpoint_oidc.discovery_url
  register       = var.app_endpoint_oidc.register
  roles_claim    = var.app_endpoint_oidc.roles_claim
  user_prefix    = var.app_endpoint_oidc.user_prefix
  username_claim = var.app_endpoint_oidc.username_claim
}


