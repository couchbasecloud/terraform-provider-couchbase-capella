data "couchbase-capella_cluster" "existing_cluster" {
  organization_id = var.organization_id
  project_id      = var.project_id
  id              = var.cluster_id
}

output "deletion_protection" {
  value = data.couchbase-capella_cluster.existing_cluster.deletion_protection
}

check "deletion_protection_matches" {
  assert {
    condition     = data.couchbase-capella_cluster.existing_cluster.deletion_protection == var.deletion_protection
    error_message = "Cluster deletion_protection is ${data.couchbase-capella_cluster.existing_cluster.deletion_protection}, expected ${var.deletion_protection}."
  }
}

