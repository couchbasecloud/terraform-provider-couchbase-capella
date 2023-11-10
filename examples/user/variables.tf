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
  description = "Capella User Name"
}

variable "email" {
  description = "Capella Email Address"
}

variable "organization_roles" {
  description = "Capella Organization Roles"
  type = list(string)
}
