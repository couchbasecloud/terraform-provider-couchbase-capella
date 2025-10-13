resource "couchbase-capella_app_endpoint_activation_status" "activation1" {
    organization_id   = var.organization_id
    project_id        = couchbase-capella_project.new_project.id
    cluster_id        = couchbase-capella_cluster.new_cluster.id
    app_service_id    = couchbase-capella_app_service.new_app_service.id
    app_endpoint_name = var.app_endpoint
    state             = "Online"
    
    depends_on = [couchbase-capella_bucket.new_bucket, couchbase-capella_app_service.new_app_service, couchbase-capella_app_endpoint.endpoint1]
}


