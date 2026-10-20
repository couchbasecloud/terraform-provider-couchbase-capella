output "list_endpoints" {
  value = data.couchbase-capella_app_service_private_endpoints.list_endpoints
}

data "couchbase-capella_app_service_private_endpoints" "list_endpoints" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}
