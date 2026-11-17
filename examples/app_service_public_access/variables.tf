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

variable "public_access" {
  description = "The public access mode of the App Service."
  type        = string
  validation {
    condition     = contains(["restricted", "unrestricted"], var.public_access)
    error_message = "Public access can only be set to either 'restricted' or 'unrestricted'."
  }
}
