resource "couchbase-capella_app_endpoint_default_oidc_provider" "example_default" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  app_service_id    = var.app_service_id
  app_endpoint_name = var.app_endpoint_name

  provider_id = var.provider_id
}
