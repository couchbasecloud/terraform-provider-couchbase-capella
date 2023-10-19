output "apikeys_list" {
  value = data.capella_apikeys.existing_apikeys
}

data "capella_apikeys" "existing_apikeys" {
  organization_id = var.organization_id
}
