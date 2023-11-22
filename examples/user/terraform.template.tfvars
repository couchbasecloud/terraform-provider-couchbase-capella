auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
host            = "https://cloudapi.cloud.couchbase.com"


user = {
  name = "John ABC"
  email = "john.doe@couchbase.com"
  organization_roles = ["organizationMember"]
}

resource = {
  type = "project"
  id   = "<project_id>"
  roles = [
    "projectViewer",
    "projectDataReaderWriter"
  ]
}
