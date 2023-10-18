# Stores the certificate details in an output variable.
# Can be viewed using `terraform output certificate` command
output "certificate" {
  value = data.capella_certificate.existing_certificate
}

data "capella_certificate" "existing_certificate" {
  organization_id = data.capella_organization.existing_organization.id
  project_id      = capella_project.new_project.id
  cluster_id      = capella_cluster.new_cluster.id
}
