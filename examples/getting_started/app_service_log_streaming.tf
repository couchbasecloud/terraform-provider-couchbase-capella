resource "couchbase-capella_app_service_log_streaming" "app_service_log_streaming" {
  organization_id  = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  app_service_id = couchbase-capella_app_service.new_app_service.id

  output_type = var.app_service_log_streaming.output_type
  credentials = var.app_service_log_streaming.credentials

  depends_on = [couchbase-capella_app_service.new_app_service]
}