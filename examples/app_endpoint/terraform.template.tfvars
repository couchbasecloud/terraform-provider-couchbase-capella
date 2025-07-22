auth_token      = "v4-api-key-secret"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"
app_service_id  = "<app_service_id>"

app_endpoint = {
  organization_id    = "<organization_id>"
  project_id         = "<project_id>"
  cluster_id         = "<cluster_id>"
  app_service_id     = "<app_service_id>"
  bucket             = "example-bucket"
  name               = "example-app-endpoint"
  user_xattr_key     = null
  delta_sync_enabled = false
  scopes = {
    scope1 = {
      collections = {
        collection1 = {
          access_control_function = "exampleAccessControlFunction"
          import_filter           = "exampleImportFilter"
        }
      }
    }
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
      provider_id    = null
      is_default     = true
    }
  ]
}