variable "host" {
  description = "host URL"
}

variable "organization_id" {
  description = "Capella Organization ID"
  type = string
}

variable "auth_token" {
  description = "Authentication API Key"
  type = string
}

variable "project_id" {
  description = "Capella Project ID"  
  type = string
}

variable "cluster_id" {
  description = "Capella Cluster ID"
  type = string
}

variable "snapshot_backup" {
  description = "Backup configuration details useful for creation"

  type = object({
    retention = optional(number)
  })

  validation {
    condition = var.snapshot_backup.retention == null || (var.snapshot_backup.retention >= 24 && var.snapshot_backup.retention <= 720)
    error_message = "Retention must be between 24 and 720 hours."
  }
}
