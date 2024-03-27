variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Project Name for Project Created via Terraform"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "cluster_onoff_schedule" {
  description = "Cluster On Off Schedule configuration details useful for creation"

  type = object({
    timezone = string
  })
}

variable "days" {

  description = "Days configuration useful for cluster on/off schedule creation"

  type = list(object({
    state = string
    day   = string
    from = optional(object({
      hour   = optional(number)
      minute = optional(number)
    }))
    to = optional(object({
      hour   = optional(number)
      minute = optional(number)
    }))
  }))
}


