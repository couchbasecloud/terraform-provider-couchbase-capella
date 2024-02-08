variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Project Name for Project Created via Terraform"
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
    size          = optional(number)
    type          = string
    iops          = optional(number)
    autoexpansion = optional(bool)
  })
}

variable "support" {
  description = "Support configuration applicable to the cluster during creation"

  type = object({
    plan     = string
    timezone = string
  })
}