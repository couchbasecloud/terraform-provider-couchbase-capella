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

variable "origin" {
  description = "Set of allowed origins for CORS"
  type        = set(string)
}

variable "login_origin" {
  description = "Set of allowed login origins for CORS"
  type        = set(string)
  default     = []
}

variable "headers" {
  description = "Set of allowed headers for CORS"
  type        = set(string)
  default     = []
}

variable "max_age" {
  description = "Max age in seconds for preflight cache"
  type        = number
  default     = null
}

variable "disabled" {
  description = "Whether CORS is disabled"
  type        = bool
  default     = null
}


