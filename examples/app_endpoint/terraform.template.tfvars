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
  collection1 = {
    access_control_function = "function(doc){channel(doc.channels);}"
    import_filter           = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }"
  },
  collection2 = {
    access_control_function = "function(doc){channel(doc.channels);}"
    import_filter           = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }"
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
