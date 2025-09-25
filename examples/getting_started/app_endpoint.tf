# Stores the App Endpoint in an output variable.
# Can be viewed using `terraform output app_endpoint` command
output "app_endpoint" {
  value = couchbase-capella_app_endpoint.new_app_endpoint
}

resource "couchbase-capella_app_endpoint" "new_app_endpoint" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  app_service_id  = couchbase-capella_app_service.new_app_service.id

  bucket = var.app_endpoint.bucket
  name   = var.app_endpoint.name

  # Optional: example nested configuration (commented out)
  # cors = {
  #   origin = ["*"]
  # }
}


