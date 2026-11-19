output "app_service_self_service_upgrades" {
  value = data.couchbase-capella_app_service_self_service_upgrades.existing_upgrades
}

data "couchbase-capella_app_service_self_service_upgrades" "existing_upgrades" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}
