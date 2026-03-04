variable "organization_id" {
  description = "Capella Organization ID"
  type = string
}

variable "auth_token" {
  description = "Authentication API Key"
  type = string
}

variable "project_id" {
  description = "Capella Project ID"  
  type = string
}

variable "cluster_id" {
  description = "Capella Cluster ID"
  type = string
}

variable "app_service_id" {
    description = "App Service ID"
    type = string
}

variable "app_endpoint_name" {
    description = "App Endpoint Name"
    type = string
}

variable "new_app_endpoint_resync_job" {
    description = "App Endpoint Resync"
    type = object({
      scopes = optional(map(set(string)))
    })
}