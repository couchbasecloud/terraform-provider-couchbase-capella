variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

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

variable "backup" {
  description = "Backup configuration details useful for creation"

  type = object({
    type                       = optional(string)
  })
}

# ToDo Required for Backup Schedule, tracking under -https://couchbasecloud.atlassian.net/browse/AV-66698
variable "weekly_schedule" {
  description = ""

  type = object({
    day_of_week = optional(string)
    start_at = optional(number)
    incremental_every = optional(number)
    retention_time = optional(string)
    cost_optimized_retention = optional(bool)
  })
}