auth_token        = "<v4-api-key-secret>"
organization_id   = "<organization_id>"
project_id        = "<project_id>"
cluster_id        = "<cluster_id>"
app_service_id    = "<app_service_id>"
app_endpoint_name = "<app_endpoint_name>"

issuer    = "https://accounts.example.com"
client_id = "<client_id>"

# Optional OIDC settings
# Set to null to omit

discovery_url  = "https://accounts.example.com/.well-known/openid-configuration"
register       = true
roles_claim    = "roles"
user_prefix    = "oidc_"
username_claim = "preferred_username"
