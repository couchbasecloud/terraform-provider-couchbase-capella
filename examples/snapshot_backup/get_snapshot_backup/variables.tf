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

variable "id" {
  description = "Capella Snapshot Backup ID"
  type        = string
}

variable "cloud_snapshot_backup" {
  description = "Existing backups"
  type = object({
    organization_id = string
    project_id      = string
    cluster_id      = string
    id              = string
  })
  default = null
}
