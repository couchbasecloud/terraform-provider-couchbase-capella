output "organizations_get" {
  value = data.capella_organizations.existing_organizations
}

data "capella_organizations" "existing_organizations" {
  organization_id = var.organization_id
}
