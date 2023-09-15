output "example_project" {
  value = capella_project.my_new_project
}

output "example_projects" {
  value = data.capella_projects.existing_projects
}
