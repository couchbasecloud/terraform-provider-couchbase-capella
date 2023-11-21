output "new_user" {
  value = capella_user.new_user
}

output "user_id" {
  value = capella_user.new_user.id
}

resource "capella_user" "new_user" {
  organization_id = var.organization_id

  name  = var.user_name
  email = var.email

  organization_roles = var.organization_roles

  resources = [
    {
      type = "project"
      id   = var.project_id
      roles = [
        "projectViewer",
        "projectDataReaderWriter"
      ]
    }
  ]
}
