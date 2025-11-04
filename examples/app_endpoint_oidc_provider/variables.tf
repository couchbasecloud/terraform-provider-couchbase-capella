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
  description = "The URL for the OpenID Connect issuer."
}

variable "client_id" {
  description = "The OpenID Connect provider client ID."
}

variable "discovery_url" {
  description = "The URL for the non-standard discovery endpoint."
  default     = null
}

variable "register" {
  description = "Indicates whether to register a new App Service user account when a user logs in using OpenID Connect."
  default     = null
}

variable "roles_claim" {
  description = "If set, the value(s) of the given OpenID Connect authentication token claim will be added to the user's roles. The value of this claim in the OIDC token must be either a string or an array of strings, any other type will result in an error."
  default     = null
}

variable "user_prefix" {
  description = "Username prefix for all users created for this provider"
  default     = null
}

variable "username_claim" {
  description = "Allows a different OpenID Connect field to be specified instead of the Subject (sub)."
  default     = null
}
