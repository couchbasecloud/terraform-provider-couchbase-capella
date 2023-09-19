variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
}

variable "project_name" {
  default     = "terraform-couchbasecapella-project"
  description = "Project Name for Project Created via Terraform"
}

variable "import_project_id" {
  default     = "7f30652a-849d-4041-8862-8a09259bf341"
  description = "Project ID for a project that already exists and is imported to Terraform"
}

variable "imported_project_name" {
  default     = "existing-project-in-capella"
  description = "Project Name for a project that already exists and is imported to Terraform"
}