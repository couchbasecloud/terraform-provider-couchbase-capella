output "new_user" {
  value = capella_user.new_user
}

output "user_id" {
  value = capella_user.new_user.id
}

resource "capella_user" "new_user" {
  organization_id = var.organization_id

  name = var.user.name
  email = var.user.email
  organization_roles = var.user.organization_roles
#  resources = [
#    {
#      id = var.resource.id
#      roles = var.resource.roles
#      type = var.resource.type
#    }
#  ]
}
