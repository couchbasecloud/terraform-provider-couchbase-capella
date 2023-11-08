output "app_services_list" {
  value = data.capella_app_services.existing_app_services
}

data "capella_app_services" "existing_app_services" {
  organization_id = var.organization_id
}