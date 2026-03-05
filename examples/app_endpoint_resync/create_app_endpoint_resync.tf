output "new_app_endpoint_resync_job" {
  value = couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job
}

resource "couchbase-capella_app_endpoint_resync_job" "new_app_endpoint_resync_job" {
  organization_id  = var.organization_id
  project_id = var.project_id
  cluster_id = var.cluster_id
  app_service_id = var.app_service_id
  app_endpoint_name = var.app_endpoint_name
  scopes = var.new_app_endpoint_resync_job.scopes
}