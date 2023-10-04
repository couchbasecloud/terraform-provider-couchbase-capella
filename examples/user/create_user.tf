output "new_user" {
  value = capella_user.new_user
}

resource "capella_user" "new_user" {
  organization_id = var.organization_id
  name            = "Matty"
  email           = "matty.maclean+17@couchbase.com"
  organization_roles = [
    "projectCreator"
  ]
  resources = [
    {
      type = "project"
      id   = var.project_id
      roles = [
        "projectDataReaderWriter"
      ]
    }
  ]
}
