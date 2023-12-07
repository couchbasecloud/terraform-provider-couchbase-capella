# Stores the certificate details in an output variable.
# Can be viewed using `terraform output certificate` command
output "certificate" {
  value = data.couchbase-capella_certificate.existing_certificate
}

data "couchbase-capella_certificate" "existing_certificate" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
}
