data "couchbase-capella_cloud_accounts" "existing_cloud_accounts" {
  organization_id = var.organization_id
}

output "existing_cloud_accounts" {
  value = data.couchbase-capella_cloud_accounts.existing_cloud_accounts
}
