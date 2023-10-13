output "certificates_get" {
  value = data.capella_certificates.existing_certificates
}

data "capella_certificates" "existing_certificates" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
