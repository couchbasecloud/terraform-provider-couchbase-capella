output "app_endpoint_cors" {
  value = couchbase-capella_app_endpoint_cors.this
}

resource "couchbase-capella_app_endpoint_cors" "this" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  app_service_id    = var.app_service_id
  app_endpoint_name = var.app_endpoint_name

  origin       = var.origin
  login_origin = var.login_origin
  headers      = var.headers
  max_age      = var.max_age
  disabled     = var.disabled
}


