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

variable "existing_cloud_project_snapshot_backups" {
  description = "Existing project backups"
  type        = list(object({
    organization_id = string
    project_id      = string
    page            = optional(number)
    per_page        = optional(number)
    sort_by         = optional(string)
    sort_direction  = optional(string)
  }))
  default = []
}
