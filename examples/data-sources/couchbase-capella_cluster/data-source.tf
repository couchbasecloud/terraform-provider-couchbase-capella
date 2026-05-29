data "couchbase-capella_cluster" "existing_cluster" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  id              = "<cluster_id>"
}

output "existing_cluster" {
  value = data.couchbase-capella_cluster.existing_cluster
}
