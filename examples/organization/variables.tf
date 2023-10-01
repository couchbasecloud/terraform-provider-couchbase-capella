variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_name" {
  default     = "terraform-couchbasecapella-project"
  description = "Organization Name for Organization Created via Terraform"
}