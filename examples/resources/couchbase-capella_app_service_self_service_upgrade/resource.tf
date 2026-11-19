resource "couchbase-capella_app_service_self_service_upgrade" "new_self_service_upgrade" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
  version_id      = "4.0.1"
  queue_at        = "2027-01-16T02:30:00Z"
}
