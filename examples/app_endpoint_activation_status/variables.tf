variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "app_service_id" {
  description = "Capella App Service ID"
}

variable "app_endpoint_name" {
  description = "App Endpoint name"
}

variable "state" {
  description = "The current state of the app endpoint. Valid values are `Online` and `Offline`. The provider may set other values on refresh based on remote state."
  type        = state
}


