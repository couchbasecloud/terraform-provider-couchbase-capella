resource "couchbase-capella_app_services_cidr" "new_allowed_cidr" {
    organization_id = var.organization_id
    project_id      = couchbase-capella_project.new_project.id
    cluster_id      = couchbase-capella_cluster.new_cluster.id
    app_service_id = couchbase-capella_app_service.new_app_service.id
    cidr            = var.app_services_cidr
    comment         = var.comment
    expires_at      = var.expires_at

    depends_on = [couchbase-capella_app_service.new_app_service]
}
