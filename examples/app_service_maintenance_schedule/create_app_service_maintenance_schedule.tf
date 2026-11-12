output "new_app_service_maintenance_schedule" {
  value = couchbase-capella_app_service_maintenance_schedule.new_app_service_maintenance_schedule
}

resource "couchbase-capella_app_service_maintenance_schedule" "new_app_service_maintenance_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  # The first pending record changes once it starts, which would replace this resource onto another
  # maintenance. Pin a specific ID from the list for anything longer-lived than one apply.
  maintenance_id = local.reschedulable_maintenance_ids[0]
  queue_at       = var.queue_at
}
