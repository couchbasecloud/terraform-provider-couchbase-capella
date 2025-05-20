resource "couchbase-capella_collection" "new_collection" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  bucket_id       = "YjE="
  scope_name      = "terraform_scope"
  collection_name = "new_terraform_collection"
  max_ttl         = 200
}