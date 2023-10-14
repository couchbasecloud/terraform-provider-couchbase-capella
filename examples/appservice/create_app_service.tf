output "new_app_service" {
    value = capella_app_service.new_app_service
}

resource "capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id = var.cluster_id
  name = var.app_service.name
  description = var.app_service.description
  nodes = var.app_service.nodes
  compute = {
     cpu = var.compute.cpu
     ram = var.compute.ram
  }
  version = var.app_service.version
}