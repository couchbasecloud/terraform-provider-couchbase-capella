resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
  organization_id            = "organization_id"
  project_id                 = "project_id"
  cluster_id                 = "cluster_id"
  state                      = "off"
  turn_on_linked_app_service = "true"
}