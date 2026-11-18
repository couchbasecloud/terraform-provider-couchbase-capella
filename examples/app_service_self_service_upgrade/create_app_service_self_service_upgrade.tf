output "new_app_service_self_service_upgrade" {
  value = couchbase-capella_app_service_self_service_upgrade.new_app_service_self_service_upgrade
}

resource "couchbase-capella_app_service_self_service_upgrade" "new_app_service_self_service_upgrade" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  version_id      = var.version_id
  queue_at        = var.queue_at
}
