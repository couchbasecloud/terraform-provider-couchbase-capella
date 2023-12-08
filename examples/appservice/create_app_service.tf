output "new_app_service" {
  value = couchbase-capella_app_service.new_app_service
}

output "appservice_id" {
  value = couchbase-capella_app_service.new_app_service.id
}

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = var.app_service.name
  description     = var.app_service.description
  nodes           = var.app_service.nodes
  compute = {
    cpu = var.app_service.compute.cpu
    ram = var.app_service.compute.ram
  }
}
