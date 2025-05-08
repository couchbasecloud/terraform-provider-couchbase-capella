resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
  organization_id            = "aaaaa-bbbb-cccc-dddd-eeeeeeeee"
  project_id                 = "aaaaa-bbbb-cccc-dddd-eeeeeeeee"
  cluster_id                 = "aaaaa-bbbb-cccc-dddd-eeeeeeeee"
  state                      = "on"
  turn_on_linked_app_service = "true"
}