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

variable "existing_cloud_snapshot_backup_schedule" {
  description = "Existing snapshot backup schedule"
  type        = object({
    organization_id = string
    project_id = string
    cluster_id = string
  })
  default = null
}
