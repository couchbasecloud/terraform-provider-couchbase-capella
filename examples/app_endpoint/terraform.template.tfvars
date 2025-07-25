auth_token = "<token>"

organization_id = "<orgId>"
project_id = "<projectId>"
cluster_id = "<clusterId>"
app_service_id  = "<appServiceId>"

bucket             = "test"
name               = "example-app-endpoint"
delta_sync_enabled = false

scope = "scope1"

collections = {
  red_collection = {
    access_control_function = "exampleAccessControlFunction_red"
    import_filter           = "exampleImportFilter_red"
  },
  blue_collection = {
    access_control_function = "exampleAccessControlFunction_blue"
    import_filter           = "exampleImportFilter_blue"
  },
}

cors = {
  origin      = ["https://example.com"]
  login_origin = ["https://login.example.com"]
  headers     = ["Authorization", "Content-Type"]
  max_age     = 3600
  disabled    = false
}

oidc = [
  {
    issuer         = "https://oidc.example.com"
    register       = true
    client_id      = "example-client-id"
    user_prefix    = "example-prefix"
    discovery_url  = "https://oidc.example.com/.well-known/openid-configuration"
    username_claim = "preferred_username"
    roles_claim    = "roles"
  }
]
