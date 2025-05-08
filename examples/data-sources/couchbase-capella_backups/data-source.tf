data "couchbase-capella_backups" "existing_backups" {
  organization_id = "aaaaa-bbbb-cccc-dddd-eeee"
  project_id      = "aaaaa-bbbb-cccc-dddd-eeee"
  cluster_id      = "aaaaa-bbbb-cccc-dddd-eeee"
  bucket_id       = "aaaaa-bbbb-cccc"
}
