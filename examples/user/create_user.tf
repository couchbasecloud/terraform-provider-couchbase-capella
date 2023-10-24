output "new_user" {
  value = capella_user.new_user
}

resource "capella_user" "new_user" {
  organization_id = var.organization_id

  name  = var.user_name
  email = var.email

  organization_roles = [
    "organizationMember"
  ]

  resources = [
    {
      type = "project"
      id   = var.project_id
      roles = [
        "projectViewer"
      ]
    },
    {
      type = "project"
      id   = var.project_id
      roles = [
        "projectDataReaderWriter"
      ]
    }
  ]
}
