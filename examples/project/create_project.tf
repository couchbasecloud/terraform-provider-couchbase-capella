output "new_project" {
  value = capella_project.new_project
}

output "project_id" {
  value = capella_project.new_project.id
}

resource "capella_project" "new_project" {
  organization_id = var.organization_id
  name            = var.project_name
  description     = "A Capella Project that will host many Capella clusters."
}

