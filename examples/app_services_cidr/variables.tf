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

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "app_services_cidr" {
  description = "Allowed CIDR configuration options"

  type = object({
    cidr       = string
    comment    = optional(string)
    expires_at = optional(string)
  })
}
