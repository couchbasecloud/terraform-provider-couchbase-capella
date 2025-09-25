# Stores the App Endpoint resync state in an output variable.
# Can be viewed using `terraform output app_endpoint_resync` command
output "app_endpoint_resync" {
  value = couchbase-capella_app_endpoint_resync.resync
}

resource "couchbase-capella_app_endpoint_resync" "resync" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  app_service_id  = couchbase-capella_app_service.new_app_service.id
  app_endpoint    = couchbase-capella_app_endpoint.new_app_endpoint.name

  # Optional: restrict resync to these scope/collections
  scopes = var.app_endpoint_resync.scopes
}


