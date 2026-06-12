resource "couchbase-capella_cluster_deletion_protection" "example" {
  organization_id     = "<organization_id>"
  project_id          = "<project_id>"
  cluster_id          = "<cluster_id>"
  deletion_protection = true
}
