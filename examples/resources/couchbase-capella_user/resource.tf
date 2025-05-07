resource "couchbase-capella_user" "new_user" {
  organization_id = "<organization_id>"

  name               = "John"
  email              = "john.doe@example.com"
  organization_roles = ["organizationMember"]
  resources = [
    {
      id    = "<project_id>"
      roles = ["projectViewer"]
      type  = "project"
    },
    {
      id    = "<project_id>"
      roles = ["projectDataReaderWriter"]
      type  = "project"
    },
  ]
}
