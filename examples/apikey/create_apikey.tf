output "existing_project" {
  value = capella_project.existing_project
}

output "new_apikey" {
  value     = capella_apikey.new_apikey
  sensitive = true
}

resource "capella_project" "existing_project" {
  organization_id = var.organization_id
  name            = var.project_name
  description     = "A Capella Project that will host many Capella clusters."
}

resource "capella_apikey" "new_apikey" {
  organization_id    = var.organization_id
  name               = var.apikey.name
  organization_roles = var.apikey.organization_roles
  allowed_cidrs      = var.apikey.allowed_cidrs
  resources = [
    {
      id    = capella_project.existing_project.id
      roles = var.resource.roles
      type  = var.resource.type
    }
  ]
}

