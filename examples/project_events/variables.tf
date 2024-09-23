variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "events" {
  type = object({
    cluster_ids     = optional(set(string))
    user_ids        = optional(set(string))
    severity_levels = optional(set(string))
    tags            = optional(set(string))
    from            = optional(string)
    to              = optional(string)
    page            = optional(number)
    per_page        = optional(number)
    sort_by         = optional(string)
    sort_direction  = optional(string)
  })
  default = {}
}