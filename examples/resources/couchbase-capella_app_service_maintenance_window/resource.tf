resource "couchbase-capella_app_service_maintenance_window" "new_maintenance_window" {
  organization_id  = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  project_id       = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  cluster_id       = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  app_service_id   = "aaaaaa-bbbbb-cccc-dddd-eeeeeeeeeeee"
  preferred_day    = "sunday"
  preferred_hour   = 2
  preferred_minute = 30
}
