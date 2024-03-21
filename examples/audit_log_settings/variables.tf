variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "auditlogsettings" {
  description = "configure cluster audit log settings"

  type = object({
    audit_enabled     = bool
    enabled_event_ids = list(number)
    disabled_users = list(object({
      name   = string
      domain = string
    }))
  })
}