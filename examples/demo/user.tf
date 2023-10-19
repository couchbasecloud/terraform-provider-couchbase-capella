# Stores the user details in an output variable.
# Can be viewed using `terraform output user` command
output "user" {
  value = capella_user.new_user
}

resource "capella_user" "new_user" {
  organization_id = var.organization_id

  name  = var.user.name
  email = var.user.email

  organization_roles = [
    "organizationMember"
  ]

  resources = [
    {
      type = "project"
      id   = capella_project.new_project.id
      roles = [
        "projectViewer"
      ]
    },
    {
      type = "project"
      id   = capella_project.new_project.id
      roles = [
        "projectDataReaderWriter"
      ]
    }
  ]
}
