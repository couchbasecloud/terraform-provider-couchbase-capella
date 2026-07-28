# Retrieve all eventing functions in a Capella cluster, optionally filtered by state.
data "couchbase-capella_eventing_functions" "existing_functions" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  # status is an optional filter. When set, only functions in the listed states are returned.
  # When omitted, eventing functions in every state are returned.
  # status = ["deployed", "deploying"]
}

# Eventing functions expose sensitive attributes such as code, so the output is marked sensitive.
output "existing_functions" {
  value     = data.couchbase-capella_eventing_functions.existing_functions
  sensitive = true
}
