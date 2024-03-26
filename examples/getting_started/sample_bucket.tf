output "sample_bucket" {
  value = couchbase-capella_sample_bucket.new_sample_bucket.name
}

resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
  name            = var.sample_bucket.name
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
}
