# Stores the App Endpoint CORS configuration in an output variable.
# Can be viewed using `terraform output app_endpoint_cors` command
output "app_endpoint_cors" {
  value = couchbase-capella_app_endpoint_cors.cors
}

resource "couchbase-capella_app_endpoint_cors" "cors" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  app_service_id    = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = couchbase-capella_app_endpoint.new_app_endpoint.name

  origin       = var.app_endpoint_cors.origin
  login_origin = var.app_endpoint_cors.login_origin
  headers      = var.app_endpoint_cors.headers
  max_age      = var.app_endpoint_cors.max_age
  disabled     = var.app_endpoint_cors.disabled
}


