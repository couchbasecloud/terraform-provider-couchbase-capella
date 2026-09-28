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

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "database_role" {
  description = "Database role configuration"
  type = object({
    name        = string
    description = optional(string)
  })
}

variable "access" {
  description = "Access privileges and resource scopes for the database role"
  type = list(object({
    privileges = set(string)
    resources = optional(object({
      buckets = set(object({
        name = string
        scopes = optional(set(object({
          name        = string
          collections = optional(set(string))
        })))
      }))
    }))
  }))
}

