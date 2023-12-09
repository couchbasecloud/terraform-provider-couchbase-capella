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

variable "allowlist" {
  description = "Allowlist configuration details useful for creation"

  type = object({
    cidr       = string
    comment    = optional(string)
    expires_at = optional(string)
  })
}
