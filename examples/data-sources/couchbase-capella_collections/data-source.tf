data "couchbase-capella_collections" "existing_collections" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  bucket_id       = "YjE="
  scope_name      = "terraform_scope"
}
