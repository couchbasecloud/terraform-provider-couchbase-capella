output "existing_cluster" {
  value = data.couchbase-capella_cluster.existing_cluster
}

data "couchbase-capella_cluster" "existing_cluster" {
  organization_id = var.organization_id
  project_id      = var.project_id
  id              = var.cluster_id
}
