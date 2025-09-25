resource "couchbase-capella_app_endpoint" "endpoint1" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  app_service_id = var.app_service.id
  name          = var.app_endpoint
  bucket        = var.bucket.name
  
  scopes = {
    (var.scope.scope_name) = {
      collections = {
        (var.collection.collection_name) = {
          access_control_function = "function (doc, oldDoc, meta) {channel('c1');}"
          import_filter           = ""
        }
      }
    }
  }
  
  cors = {
    disabled = false
    origin   = ["*"]
    headers  = ["Authorization", "Content-Type"]
    max_age  = 3600
  }
}
