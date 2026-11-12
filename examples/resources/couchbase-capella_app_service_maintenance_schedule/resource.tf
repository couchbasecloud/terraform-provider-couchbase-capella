resource "couchbase-capella_app_service_maintenance_schedule" "new_maintenance_schedule" {
  organization_id = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  project_id      = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  cluster_id      = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  app_service_id  = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  maintenance_id  = "<maintenance_id>"
  queue_at        = "2026-09-05T02:30:00Z"
}
