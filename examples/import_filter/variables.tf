variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "import_filter" {
  description = "Import Filter configuration"

  type = object({
    organization_id = string
    project_id      = string
    cluster_id      = string
    app_service_id  = string
    keyspace        = string
    import_filter   = string
  })
}


