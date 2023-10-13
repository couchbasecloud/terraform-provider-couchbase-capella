output "existing_organization" {
  value = data.capella_organization.existing_organization
}

data "capella_organization" "existing_organization" {
  organization_id = var.organization_id
}
