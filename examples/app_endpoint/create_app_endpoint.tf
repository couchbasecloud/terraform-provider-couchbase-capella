resource "couchbase-capella_app_endpoint" "endpoint2" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  name            = "${var.app_endpoint_name}-cors"
  bucket          = var.bucket_name
  
  scopes = {
    "s1" = {
      collections = {
        "c1" = {
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
