output "app_service_self_service_upgrade_available_versions" {
  value = data.couchbase-capella_app_service_self_service_upgrade_available_versions.available_versions
}

data "couchbase-capella_app_service_self_service_upgrade_available_versions" "available_versions" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}

# Copy one of these IDs into var.version_id.
output "available_version_ids" {
  value = [
    for v in data.couchbase-capella_app_service_self_service_upgrade_available_versions.available_versions.data : v.id
  ]
}
