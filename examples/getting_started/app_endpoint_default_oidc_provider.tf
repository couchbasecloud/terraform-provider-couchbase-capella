# Stores the App Endpoint Default OIDC provider mapping in an output variable.
# Can be viewed using `terraform output app_endpoint_default_oidc_provider` command
output "app_endpoint_default_oidc_provider" {
  value = couchbase-capella_app_endpoint_default_oidc_provider.default
}

resource "couchbase-capella_app_endpoint_default_oidc_provider" "default" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  app_service_id    = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = couchbase-capella_app_endpoint.new_app_endpoint.name

  provider_id = var.app_endpoint_default_oidc.provider_id
}


