variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "bucket_id" {
  description = "Capella Bucket ID"
}

variable "backup_schedule" {
  description = "Backup Schedule configuration details useful for creation"

  type = object({
    type = string
    weekly_schedule = object({
      day_of_week              = string
      start_at                 = number
      incremental_every        = number
      retention_time           = string
      cost_optimized_retention = bool
    })
  })
}