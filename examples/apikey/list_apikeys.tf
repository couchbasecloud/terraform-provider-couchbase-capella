output "apikeys_list" {
  value = data.couchbase-capella_apikeys.existing_apikeys
}

data "couchbase-capella_apikeys" "existing_apikeys" {
  organization_id = var.organization_id
}
