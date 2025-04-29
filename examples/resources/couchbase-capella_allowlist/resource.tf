resource "couchbase-capella_allowlist" "new_allowlist" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  cidr            = "10.0.0.0/16"
  comment         = "Allow access from a public IP"
  expires_at      = "2023-12-30T23:59:59.465Z"
}
