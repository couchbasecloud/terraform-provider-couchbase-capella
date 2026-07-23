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
  # load_balancer_cidr pins the load balancer subnet CIDR (Azure App Services only). Allocated dynamically when omitted.
  load_balancer_cidr = var.app_service.load_balancer_cidr
}
