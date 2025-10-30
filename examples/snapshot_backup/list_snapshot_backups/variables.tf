variable "organization_id" {
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

variable "existing_cloud_snapshot_backups" {
  description = "Existing backups"
  type = list(object({
    organization_id = string
    project_id      = string
    cluster_id      = string
  }))
  default = []
}
