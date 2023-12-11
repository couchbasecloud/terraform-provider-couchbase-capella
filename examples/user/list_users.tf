output "users_list" {
  value = data.couchbase-capella_users.existing_users
}

data "couchbase-capella_users" "existing_users" {
  organization_id = var.organization_id
}
