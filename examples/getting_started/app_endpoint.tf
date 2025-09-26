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

  oidc = [
    {
      issuer         = "https://example-issuer.com"
      register       = false
      client_id      = "example-client-id"
      user_prefix    = "user_"
      discovery_url  = "https://example-issuer.com/.well-known/openid-configuration"
      username_claim = "sub"
      roles_claim    = "roles"
    }
  ]

  cors = {
    disabled = false
    origin   = ["http://example.com", "http://staging.example.com"]
    last_origin = ["http://example.com"]
    headers  = ["Authorization", "Content-Type"]
    max_age  = 3600
  }
}
