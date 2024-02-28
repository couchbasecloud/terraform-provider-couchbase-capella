output "new_project" {
  value = couchbase-capella_project.new_project
}

output "project_id" {
  value = couchbase-capella_project.new_project.id
}

resource "couchbase-capella_project" "new_project" {
  organization_id = var.organization_id
  name            = "project2"
  description     = "A Capella Project that will host many Capella clusters."
}

