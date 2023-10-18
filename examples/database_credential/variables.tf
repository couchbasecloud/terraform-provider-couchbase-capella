variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "database_credential_name" {
  description = "Database Credentials Name"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "password" {
  description = "password for database credential"
  sensitive   = true
}

variable "access" {
  type = list(object({
    privileges = list(string)
  }))
}
