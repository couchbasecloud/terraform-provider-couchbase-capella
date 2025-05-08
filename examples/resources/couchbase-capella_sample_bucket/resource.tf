resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
  name            = "travel-sample"
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
}
