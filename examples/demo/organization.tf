# Stores the organization name in an output variable.
# Can be viewed using `terraform output organization` command
output "organization" {
  value = data.capella_organization.existing_organization
}

data "capella_organization" "existing_organization" {
  organization_id = var.organization_id
}
