variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "import_filter" {
  description = "Import Filter configuration"

  type = object({
    organization_id     = string
    project_id          = string
    cluster_id          = string
    app_service_id      = string
    app_endpoint_name   = string
    scope               = optional(string, "_default")
    collection          = optional(string, "_default")
    import_filter       = string
  })
}


