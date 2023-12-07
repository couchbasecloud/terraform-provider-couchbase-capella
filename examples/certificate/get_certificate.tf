output "existing_certificate" {
  value = data.couchbase-capella_certificate.existing_certificate
}

data "couchbase-capella_certificate" "existing_certificate" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
