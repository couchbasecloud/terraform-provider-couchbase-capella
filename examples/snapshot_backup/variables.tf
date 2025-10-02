variable "host" {
  description = "Capella Host URL"
  type        = string
}

variable "organization_id" {
  description = "Capella Organization ID"
  type        = string
}

variable "auth_token" {
  description = "Authentication API Key"
  type        = string
}

variable "project_id" {
  description = "Capella Project ID"
  type        = string
}

variable "cluster_id" {
  description = "Capella Cluster ID"
  type        = string
}

variable "cloud_snapshot_backup" {
  description = "Backup configuration details useful for creation"

  type = object({
    retention = optional(number)
    regions_to_copy = optional(list(string))
    restore_times = optional(number)
    cross_region_restore_preference = optional(list(string))
      })

  validation {
    condition = var.cloud_snapshot_backup.retention == null || (var.cloud_snapshot_backup.retention >= 24 && var.cloud_snapshot_backup.retention <= 720)
    error_message = "Retention must be between 24 and 720 hours."
  }

  validation {
    condition = var.cloud_snapshot_backup.retention == null || var.cloud_snapshot_backup.retention == floor(var.cloud_snapshot_backup.retention)
    error_message = "Retention must be an integer."
  }

  validation {
    condition = var.cloud_snapshot_backup.restore_times == null || var.cloud_snapshot_backup.restore_times == floor(var.cloud_snapshot_backup.restore_times)
    error_message = "Restore times must be an integer."
  }
}
