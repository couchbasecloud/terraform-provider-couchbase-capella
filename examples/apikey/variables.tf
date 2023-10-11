variable "host" {
  default     = "https://cloudapi.dev.nonprod-project-avengers.com"
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "project_name" {
  default     = "terraform-couchbasecapella-project"
  description = "Project Name for Project Created via Terraform"
}

variable "apikey" {
  description = "ApiKey creation details useful for apikey creation"

  type = object({
    name               = string
    description        = string
    allowed_cidrs      = list(string)
    organization_roles = list(string)
    expiry             = number
  })
}

variable "resource" {
  description = "Resource details useful for apikey creation"

  type = object({
    id    = string
    roles = list(string)
    type  = string
  })
}