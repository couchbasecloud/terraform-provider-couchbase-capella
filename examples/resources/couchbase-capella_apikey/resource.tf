resource "couchbase-capella_apikey" "new_apikey" {
  organization_id    = "<organization_id>"
  name               = "Organization Owner API Key"
  description        = "Creates an API key with a Organization Owner role."
  expiry             = 720
  organization_roles = ["organizationOwner"]
  allowed_cidrs      = ["8.8.8.8/32"]
  resources          = []
}