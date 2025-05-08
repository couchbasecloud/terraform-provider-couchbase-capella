resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  name            = "MyAppSyncService"
  description     = "My app sync service."
  nodes           = 2
  compute = {
    cpu = 2
    ram = 4
  }
}