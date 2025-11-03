output "new_allowlist" {
  value = couchbase-capella_app_services_cidr.new_allowlist
}

output "allowlist_id" {
  value = couchbase-capella_app_services_cidr.new_allowlist.id
}

resource "couchbase-capella_app_services_cidr" "new_allowlist" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  cidr            = var.app_services_cidr.cidr
  comment         = var.app_services_cidr.comment
  expires_at      = var.app_services_cidr.expires_at
}
