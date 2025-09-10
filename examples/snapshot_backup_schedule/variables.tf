
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

variable "snapshot_backup_schedule" {
  description = "Snapshot backup schedule configuration details useful for creation"

  type = object({
    interval = optional(number)
    retention = optional(number)
    start_time = optional(string)
  })

  validation {
    condition = (var.snapshot_backup_schedule == null && var.snapshot_backup_schedule.interval == null && var.snapshot_backup_schedule.start_time == null) || (var.snapshot_backup_schedule != null && var.snapshot_backup_schedule.interval != null && var.snapshot_backup_schedule.start_time != null)
    error_message = "Either all or none of the snapshot backup schedule attributes must be provided."
  }

  validation {
    condition     = var.snapshot_backup_schedule.interval >= 1 && var.snapshot_backup_schedule.interval <= 24
    error_message = "Interval must be between 1 and 24 hours."
  }

  validation {
    condition = var.snapshot_backup_schedule.interval == floor(var.snapshot_backup_schedule.interval)
    error_message = "Interval must be an integer."
  }

  validation {
    condition = var.snapshot_backup_schedule.retention >= 24 && var.snapshot_backup_schedule.retention <= 720
    error_message = "Retention must be between 24 and 720 hours."
  }

  validation {
    condition = var.snapshot_backup_schedule.retention == floor(var.snapshot_backup_schedule.retention)
    error_message = "Retention must be an integer."
  }
}
