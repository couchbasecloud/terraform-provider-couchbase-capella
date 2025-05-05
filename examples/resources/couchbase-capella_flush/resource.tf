resource "couchbase-capella_flush" "new_flush" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  bucket_id       = "<bucket_id>"
}
