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

variable "cloud_snapshot_backup_schedule" {
  description = "Snapshot backup schedule configuration details useful for creation"

  type = object({
    interval = number
    retention = number
    start_time = optional(string)
    copy_to_regions = optional(list(string))
  })

  validation {
    condition = (var.cloud_snapshot_backup_schedule == null && var.cloud_snapshot_backup_schedule.interval == null && var.cloud_snapshot_backup_schedule.start_time == null) || (var.cloud_snapshot_backup_schedule.retention != null && var.cloud_snapshot_backup_schedule.interval != null)
    error_message = "Either all or none of the snapshot backup schedule attributes must be provided."
  }

  validation {
    condition = anytrue([var.cloud_snapshot_backup_schedule.interval == null, 
                          var.cloud_snapshot_backup_schedule.interval == 1, 
                          var.cloud_snapshot_backup_schedule.interval == 2, 
                          var.cloud_snapshot_backup_schedule.interval == 4, 
                          var.cloud_snapshot_backup_schedule.interval == 6, 
                          var.cloud_snapshot_backup_schedule.interval == 8, 
                          var.cloud_snapshot_backup_schedule.interval == 12, 
                          var.cloud_snapshot_backup_schedule.interval == 24])
    error_message = "Interval must be 1, 2, 4, 6, 8, 12, or 24 hours."
  }

  validation {
    condition = var.cloud_snapshot_backup_schedule.retention == null || (var.cloud_snapshot_backup_schedule.retention >= 24 && var.cloud_snapshot_backup_schedule.retention <= 720)
    error_message = "Retention must be between 24 and 720 hours."
  }

  validation {
    condition = var.cloud_snapshot_backup_schedule.retention == null || var.cloud_snapshot_backup_schedule.retention == floor(var.cloud_snapshot_backup_schedule.retention)
    error_message = "Retention must be an integer."
  }
}
