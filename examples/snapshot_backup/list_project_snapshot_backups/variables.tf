variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
}

variable "project_id" {
  description = "Capella Project ID"
}


variable "page" {
  description = "page"
}

variable "per_page" {
  description = "number per page"
}

variable "sort_by" {
  description = "sort by"
}

variable "sort_direction" {
  description = "sort_direction"
}

variable "existing_cloud_project_snapshot_backups" {
  description = "Existing project backups"
  type = list(object({
    organization_id = string
    project_id      = string
    page            = optional(number)
    per_page        = optional(number)
    sort_by         = optional(string)
    sort_direction  = optional(string)
  }))
  default = []
}
