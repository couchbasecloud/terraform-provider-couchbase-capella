output "organizations_get" {
  value = data.couchbase-capella_organization.existing_organization
}

data "couchbase-capella_organization" "existing_organization" {
  organization_id = var.organization_id
}
