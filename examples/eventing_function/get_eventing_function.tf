# Retrieve a single eventing function by name from a Capella cluster.
data "couchbase-capella_eventing_function" "existing_function" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = var.eventing_function.name

  # When export is true, read-only fields are omitted from the output.
  export = false
}

# The eventing function exposes sensitive attributes such as code, so the output is marked sensitive.
output "existing_function" {
  value     = data.couchbase-capella_eventing_function.existing_function
  sensitive = true
}
