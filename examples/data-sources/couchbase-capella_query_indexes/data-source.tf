data "couchbase-capella_query_indexes" "list" {
  organization_id = <organization_id>
  project_id      = <project_id>
  cluster_id      = <cluster_id>
  bucket_name     = "api"
  scope_name      = "metrics"
  collection_name = "memory"
}