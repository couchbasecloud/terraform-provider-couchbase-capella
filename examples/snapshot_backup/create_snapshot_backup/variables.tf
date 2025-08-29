variable "tenant_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "snapshot_backup" {
  description = "Backup configuration details useful for creation"

  type = object({
    retention = number
  })
}
