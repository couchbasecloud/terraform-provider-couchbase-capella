# Stores the user details in an output variable.
# Can be viewed using `terraform output user` command
output "user" {
  value = couchbase-capella_user.new_user
}

resource "couchbase-capella_user" "new_user" {
  organization_id = var.organization_id

  name  = var.user.name
  email = var.user.email

  organization_roles = [
    "organizationMember"
  ]

  resources = [
    {
      type = "project"
      id   = couchbase-capella_project.new_project.id
      roles = [
        "projectViewer",
        "projectDataReaderWriter"
      ]
    }
  ]
}
