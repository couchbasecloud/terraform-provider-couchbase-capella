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
    description = "Capella App Service ID"
    type = string
}

variable "app_endpoint_name" {
    description = "Capella App Endpoint Name"
    type = string
}

variable "app_endpoint_log_streaming_config" {
    description = "App endpoint log streaming config"
    type = object({
        organization_id = string
        project_id = string
        cluster_id = string
        app_service_id = string
        app_endpoint_name = string
    })
    default = null
}
