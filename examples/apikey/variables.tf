variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "apikey" {
  description = "ApiKey creation details useful for apikey creation"

  type = object({
    name               = string
    description        = optional(string)
    allowed_cidrs      = optional(list(string))
    organization_roles = list(string)
    expiry             = optional(number)
  })
}

variable "resources" {
  description = "Resource details useful for apikey creation"

  type = list(object({
    id    = string
    roles = list(string)
    type  = optional(string)
  }))
  default = []
}
