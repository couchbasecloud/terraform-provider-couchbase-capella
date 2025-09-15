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

variable "issuer" {
  description = "OIDC issuer URL"
}

variable "client_id" {
  description = "OIDC client ID"
}

variable "discovery_url" {
  description = "Non-standard discovery endpoint URL"
  default     = null
}

variable "register" {
  description = "Register new App Service user on first login"
  default     = null
}

variable "roles_claim" {
  description = "Token claim providing roles"
  default     = null
}

variable "user_prefix" {
  description = "Username prefix for users created by this provider"
  default     = null
}

variable "username_claim" {
  description = "Token claim to use for username"
  default     = null
}
