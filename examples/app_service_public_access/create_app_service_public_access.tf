output "app_service_public_access" {
  value = couchbase-capella_app_service_public_access.app_service_public_access
}

resource "couchbase-capella_app_service_public_access" "app_service_public_access" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  public_access = var.public_access
}
