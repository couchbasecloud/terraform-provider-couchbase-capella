output "app_service_maintenances" {
  value = data.couchbase-capella_app_service_maintenances.existing_maintenances
}

data "couchbase-capella_app_service_maintenances" "existing_maintenances" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}

locals {
  # Only pending, Couchbase-initiated maintenance can be rescheduled with this resource.
  reschedulable_maintenance_ids = [
    for m in data.couchbase-capella_app_service_maintenances.existing_maintenances.data : m.id
    if m.current_state == "pending" && m.initiated_by == "couchbase"
  ]
}
