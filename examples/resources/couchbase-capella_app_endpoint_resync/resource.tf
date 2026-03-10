resource "couchbase-capella_app_endpoint_resync_job" "new_app_endpoint_resync_job" {
  organization_id   = "aaaaa-bbbb-cccc-dddd-eeee"
  project_id        = "aaaaa-bbbb-cccc-dddd-eeee"
  cluster_id        = "aaaaa-bbbb-cccc-dddd-eeee"
  app_service_id    = "aaaaa-bbbb-cccc-dddd-eeee"
  app_endpoint_name = "<app_endpoint_name>"
}