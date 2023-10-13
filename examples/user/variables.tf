variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

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

variable "user_name" {
  description = "Name of the user"
  type        = string
}

variable "user_email" {
  description = "Email address of the user"
  type        = string
}

variable "org_roles" {
  description = "Roles of the user within the organization"
  type        = list(string)
  default     = []
}

variable "project_roles" {
  description = "Roles of the user within the project"
  type        = list(string)
  default     = []
}
