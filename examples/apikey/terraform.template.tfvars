auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
host            = "https://cloudapi.cloud.couchbase.com"

apikey = {
  name               = "New Terraform Api Key"
  description        = "A Capella Api Key"
  allowed_cidrs      = ["10.1.42.0/23", "10.1.42.0/23"]
  organization_roles = ["organizationMember"]
  expiry             = 179
}

resource = {
  id    = "<project_id>"
  roles = ["projectManager", "projectDataReader"]
  type  = "project"
}
