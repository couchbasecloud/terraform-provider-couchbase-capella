# Retrieve a single eventing function by name from a Capella cluster.
data "couchbase-capella_eventing_function" "existing_function" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  name            = "<function_name>"

  # When export is true, read-only fields are omitted from the output.
  export = false
}

# The eventing function exposes sensitive attributes such as code, so the output is marked sensitive.
output "existing_function" {
  value     = data.couchbase-capella_eventing_function.existing_function
  sensitive = true
}
