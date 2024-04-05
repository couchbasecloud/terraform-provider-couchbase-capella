output "new_sample_bucket" {
  value = couchbase-capella_sample_bucket.new_sample_bucket
}

output "sample_bucket_id" {
  value = couchbase-capella_sample_bucket.new_sample_bucket.id
}

resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
  name            = var.sample_bucket_name
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
