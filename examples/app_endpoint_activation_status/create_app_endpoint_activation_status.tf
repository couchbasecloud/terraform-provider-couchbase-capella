resource "couchbase-capella_app_endpoint_activation_status" "app_endpoint_activation" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  app_service_id    = var.app_service_id
  app_endpoint_name = var.app_endpoint_name
  state            = var.state
}


