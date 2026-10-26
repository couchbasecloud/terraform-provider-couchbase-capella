data "couchbase-capella_app_service_private_endpoints" "existing_app_service_private_endpoints" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  # Omit status to return private endpoints in every status.
  status = var.status
}

output "existing_app_service_private_endpoints" {
  value = data.couchbase-capella_app_service_private_endpoints.existing_app_service_private_endpoints
}
