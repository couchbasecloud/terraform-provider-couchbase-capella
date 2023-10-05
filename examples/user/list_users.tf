output "users_list" {
    value = data.capella_users.existing_users
}

data "capella_users" "existing_users" {
  organization_id = var.organization_id
}
