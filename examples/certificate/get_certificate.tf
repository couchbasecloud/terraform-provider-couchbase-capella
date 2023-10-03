output "existing_certificate" {
  value = data.capella_certificates.existing_certificate
}

data "capella_certificates" "existing_certificate" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
