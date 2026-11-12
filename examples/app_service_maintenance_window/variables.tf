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

variable "maintenance_window" {
  description = "Preferred weekly maintenance window of the App Service, in UTC"

  type = object({
    preferred_day    = string
    preferred_hour   = number
    preferred_minute = number
  })
}
