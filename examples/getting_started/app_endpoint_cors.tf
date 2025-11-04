resource "couchbase-capella_app_endpoint_cors" "cors1" {
    organization_id   = var.organization_id
    project_id        = couchbase-capella_project.new_project.id
    cluster_id        = couchbase-capella_cluster.new_cluster.id
    app_service_id    = couchbase-capella_app_service.new_app_service.id
    app_endpoint_name = var.app_endpoint
    
    origin       = ["http://example.com", "http://staging.example.com"]
    login_origin = ["http://example.com"]
    headers      = ["Authorization", "Content-Type"]
    max_age      = 3600
    disabled     = false
    
    depends_on = [couchbase-capella_bucket.new_bucket, couchbase-capella_app_service.new_app_service, couchbase-capella_app_endpoint.endpoint1]
}


