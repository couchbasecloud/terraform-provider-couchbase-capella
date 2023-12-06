output "new_apikey" {
  value     = couchbase-capella_apikey.new_apikey
  sensitive = true
}

output "apikey_id" {
  value = couchbase-capella_apikey.new_apikey.id
}

resource "couchbase-capella_apikey" "new_apikey" {
  organization_id    = var.organization_id
  name               = var.apikey.name
  description        = var.apikey.description
  expiry             = var.apikey.expiry
  organization_roles = var.apikey.organization_roles
  allowed_cidrs      = var.apikey.allowed_cidrs
  resources          = var.resources
}

