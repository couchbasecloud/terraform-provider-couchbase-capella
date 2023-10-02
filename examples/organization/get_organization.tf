output "organization_get" {
  value = data.capella_organizations.existing_organization
}

data "capella_organizations" "existing_organization" {
  organization_id = var.organization_id
}
