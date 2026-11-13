output "app_service_maintenances" {
  value = data.couchbase-capella_app_service_maintenances.existing_maintenances
}

data "couchbase-capella_app_service_maintenances" "existing_maintenances" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}

# Only pending maintenance can be rescheduled. Copy one of these IDs into var.maintenance_id.
output "reschedulable_maintenance_ids" {
  value = [
    for m in data.couchbase-capella_app_service_maintenances.existing_maintenances.data : m.id
    if m.current_state == "pending"
  ]
}
