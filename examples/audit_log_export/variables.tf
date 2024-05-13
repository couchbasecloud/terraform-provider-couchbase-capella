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

variable "audit_log_export" {
  description = "create audit log export job"

  type = object({
    start = string
    end   = string
  })
}

