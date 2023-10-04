output "new_user" {
  value = capella_user.new_user
}

resource "capella_user" "new_user" {
  organization_id = var.organization_id
  name            = var.user_name
  email           = var.user_email
  organization_roles = var.org_roles
  resources = [
    {
      type = "project"
      id   = var.project_id
      roles = var.project_roles
    }
  ]
}
