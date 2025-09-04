resource "couchbase-capella_app_endpoint_access_control_function" "acf" {
  organization_id         = var.organization_id
  project_id              = var.project_id
  cluster_id              = var.cluster_id
  app_service_id          = var.app_service_id
  app_endpoint            = var.app_endpoint_name
  scope                   = var.scope
  collection              = var.collection
  access_control_function = var.access_control_function
} 