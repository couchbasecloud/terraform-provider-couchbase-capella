output "new_app_service_maintenance_schedule" {
  value = couchbase-capella_app_service_maintenance_schedule.new_app_service_maintenance_schedule
}

resource "couchbase-capella_app_service_maintenance_schedule" "new_app_service_maintenance_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  maintenance_id  = local.reschedulable_maintenance_ids[0]
  queue_at        = var.queue_at
}
