variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
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

variable "state" {
  description = "The desired state for the App Service Log Streaming to be in."
  type        = string
  validation {
    condition     = contains(["enabled", "paused"], var.state)
    error_message = "Desired state can only be set to either 'enabled' or 'paused'."
  }
}