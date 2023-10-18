# Stores the project name in an output variable.
# Can be viewed using `terraform output project` command
output "project" {
  value = capella_project.new_project.name
}

resource "capella_project" "new_project" {
  organization_id = var.organization_id
  name            = var.project_name
  description     = "A Capella Project that will host many Capella clusters."
}

