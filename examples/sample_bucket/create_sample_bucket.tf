output "new_sample_bucket" {
  value = couchbase-capella_sample_bucket.new_sample_bucket
}

output "samplebucket_id" {
  value = couchbase-capella_sample_bucket.new_sample_bucket.id
}

resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
  name            = var.samplebucket.name
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
