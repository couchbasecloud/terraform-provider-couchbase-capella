variable "organization_id" {
  description = "Capella Organization ID"
}

variable "host_url" {
  description = "Capella Host URL"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "user" {
  type = object({
    name               = optional(string)
    email              = string
    organization_roles = list(string)
  })
}

variable "resources" {
  type = list(object({
    type  = optional(string)
    id    = string
    roles = list(string)
  }))
  default = null
}
