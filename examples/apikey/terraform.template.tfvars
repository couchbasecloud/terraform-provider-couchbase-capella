auth_token      = "v4-api-key-secret"
organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
host            = "https://cloudapi.cloud.couchbase.com"

apikey = {
  name               = "New Terraform Api Key"
  description        = "A Capella Api Key"
  allowed_cidrs      = ["10.1.42.0/23", "10.1.42.0/23"]
  organization_roles = ["organizationMember"]
  expiry             = 179
}

resource = {
  id    = "resource id"
  roles = ["projectManager", "projectDataReader"]
  type  = "project"
}
