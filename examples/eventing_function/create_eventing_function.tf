resource "couchbase-capella_eventing_function" "new_eventing_function" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  name        = var.eventing_function.name
  description = var.eventing_function.description
  code        = var.eventing_function.code

  state = var.eventing_function.state

  event_source = {
    bucket     = var.eventing_function.event_source.bucket
    scope      = var.eventing_function.event_source.scope
    collection = var.eventing_function.event_source.collection
  }

  event_metadata_storage = {
    bucket     = var.eventing_function.event_metadata_storage.bucket
    scope      = var.eventing_function.event_metadata_storage.scope
    collection = var.eventing_function.event_metadata_storage.collection
  }

  settings = var.eventing_function.settings

  bindings = var.eventing_function.bindings
}
