resource "couchbase-capella_app_endpoint" "endpoint1" {
    organization_id = var.organization_id
    project_id      = couchbase-capella_project.new_project.id
    cluster_id      = couchbase-capella_cluster.new_cluster.id
    app_service_id = couchbase-capella_app_service.new_app_service.id
    name          = var.app_endpoint
    bucket        = var.bucket.name
    
    scopes = {
    (var.scope.scope_name) = {
      collections = {
        (var.collection.collection_name) = {
          access_control_function = "function (doc, oldDoc, meta) {channel('c1');}"
          import_filter           = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }"
        }
      }
    }
    }
    
    oidc = [
        {
            issuer         = var.app_endpoint_oidc.issuer
            client_id      = var.app_endpoint_oidc.client_id
            register       = var.app_endpoint_oidc.register
            username_claim = var.app_endpoint_oidc.username_claim
            roles_claim    = var.app_endpoint_oidc.roles_claim
            user_prefix    = var.app_endpoint_oidc.user_prefix
        }
    ]
    
    cors = {
    disabled = false
    origin   = ["http://example.com", "http://staging.example.com"]
    last_origin = ["http://example.com"]
    headers  = ["Authorization", "Content-Type"]
    max_age  = 3600
    }
    depends_on = [couchbase-capella_bucket.new_bucket, couchbase-capella_app_service.new_app_service]
}
