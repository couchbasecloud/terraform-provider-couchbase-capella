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

variable "app_endpoint" {
  description = "App Endpoint configuration"

  type = object({
    organization_id    = string
    project_id         = string
    cluster_id         = string
    app_service_id     = string
    bucket             = string
    name               = string
    user_xattr_key     = optional(string)
    delta_sync_enabled = optional(bool)
    scopes = optional(map(object({
      collections = optional(map(object({
        access_control_function = optional(string)
        import_filter           = optional(string)
      })))
    })))
    cors = optional(object({
      origin      = optional(list(string))
      login_origin = optional(list(string))
      headers     = optional(list(string))
      max_age     = optional(number)
      disabled    = optional(bool)
    }))
    oidc = optional(list(object({
      issuer         = string
      register       = optional(bool)
      client_id      = string
      user_prefix    = optional(string)
      discovery_url  = optional(string)
      username_claim = optional(string)
      roles_claim    = optional(string)
      provider_id    = optional(string)
      is_default     = optional(bool)
    })))
  })
}