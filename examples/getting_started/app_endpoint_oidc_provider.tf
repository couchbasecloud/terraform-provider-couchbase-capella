resource "couchbase-capella_app_endpoint_oidc_provider" "oidc1" {
    organization_id   = var.organization_id
    project_id        = couchbase-capella_project.new_project.id
    cluster_id        = couchbase-capella_cluster.new_cluster.id
    app_service_id    = couchbase-capella_app_service.new_app_service.id
    app_endpoint_name = var.app_endpoint
    
    issuer         = var.app_endpoint_oidc.issuer
    client_id      = var.app_endpoint_oidc.client_id
    register       = var.app_endpoint_oidc.register
    username_claim = var.app_endpoint_oidc.username_claim
    roles_claim    = var.app_endpoint_oidc.roles_claim
    user_prefix    = var.app_endpoint_oidc.user_prefix
    
    depends_on = [couchbase-capella_bucket.new_bucket, couchbase-capella_app_service.new_app_service, couchbase-capella_app_endpoint.endpoint1]
}


