output "new_project" {
  value = capella_project.new_project
}

resource "capella_project" "new_project" {
  organization_id = var.organization_id
  name            = var.project_name
  description     = "A Capella Project that will host many Capella clusters."
}

