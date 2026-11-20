variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "alerts" {
  type = object({
    status          = optional(set(string))
    cluster_ids     = optional(set(string))
    project_ids     = optional(set(string))
    severity_levels = optional(set(string))
    from            = optional(string)
    to              = optional(string)
  })
  default = {}
}
