output "new_user" {
  value = couchbase-capella_user.new_user
}

output "user_id" {
  value = couchbase-capella_user.new_user.id
}

resource "couchbase-capella_user" "new_user" {
  organization_id = var.organization_id

  name               = var.user.name
  email              = var.user.email
  organization_roles = var.user.organization_roles
  resources          = var.resources
}
