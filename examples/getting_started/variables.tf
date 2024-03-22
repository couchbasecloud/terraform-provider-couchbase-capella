variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_name" {
  description = "Capella Project Name"
}

variable "cloud_provider" {
  description = "Cloud Provider details useful for cluster creation"

  type = object({
    name   = string
    region = string
  })
}

variable "cluster" {
  description = "Cluster configuration details useful for creation"

  type = object({
    name               = string
    cidr               = string
    node_count         = number
    couchbase_services = list(string)
    availability_zone  = string
  })
}

variable "couchbase_server" {
  type = object({
    version = string
  })
  default = null
}

variable "compute" {
  description = "All cluster node compute configuration"

  type = object({
    cpu = number
    ram = number
  })
}

variable "disk" {
  description = "All nodes' disk configuration"

  type = object({
    size = optional(number)
    type = string
    iops = optional(number)
  })
}

variable "support" {
  description = "Support configuration applicable to the cluster during creation"

  type = object({
    plan     = string
    timezone = string
  })
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

variable "comment" {
  description = "comment describing the allowlist details"
}

variable "cidr" {
  description = "cidr in the allowlist that will have access to the Capella cluster"
}

variable "expires_at" {
  description = "timestamp when the allowlist expires"
}

variable "database_credential_name" {
  description = "Database Credentials Name"
}

variable "password" {
  description = "password for database credential"
  sensitive   = true
}

variable "access" {
  type = list(object({
    privileges = list(string)
    resources = optional(object({
      buckets = list(object({
        name = string
        scopes = optional(list(object({
          name        = string
          collections = optional(list(string))
        })))
      }))
    }))
  }))
}

variable "bucket" {
  description = "Bucket configuration details useful for creation"

  type = object({
    name                       = string
    type                       = optional(string)
    storage_backend            = optional(string)
    memory_allocation_in_mb    = optional(number)
    bucket_conflict_resolution = optional(string)
    durability_level           = optional(string)
    replicas                   = optional(number)
    flush                      = optional(bool)
    time_to_live_in_seconds    = optional(number)
    eviction_policy            = optional(string)
  })
}

variable "sample_bucket" {
  description = "Bucket configuration details useful for creation"

  type = object({
    name = string
  })
}

variable "user" {
  description = "User details useful for creation"

  type = object({
    name  = string
    email = string
  })
}

variable "app_service" {
  description = "App Service configuration details useful for creation"

  type = object({
    name        = string
    description = optional(string)
    nodes       = optional(number)
    compute = object({
      cpu = number
      ram = number
    })
  })
}

variable "scope" {
  description = "Scope configuration details useful for creation"

  type = object({
    scope_name = string
  })
}


variable "collection" {
  description = "Collection configuration details useful for creation"

  type = object({
    collection_name = string
    max_ttl         = optional(number)
  })
}
