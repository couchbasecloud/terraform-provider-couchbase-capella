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

variable "app_endpoint" {
  description = "App Endpoint name"
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

variable "cluster_onoff_schedule" {
  description = "Cluster On Off Schedule configuration details useful for creation"

  type = object({
    timezone = string
  })
}

variable "days" {

  description = "Days configuration useful for cluster on/off schedule creation"

  type = list(object({
    state = string
    day   = string
    from = optional(object({
      hour   = optional(number)
      minute = optional(number)
    }))
    to = optional(object({
      hour   = optional(number)
      minute = optional(number)
    }))
  }))
}

variable "audit_log_settings" {
  description = "configure cluster audit log settings"

  type = object({
    audit_enabled = bool
    disabled_users = list(object({
      name   = string
      domain = string
    }))
  })
}

variable "enabled" {
  description = "Enable or disable private endpoint service"
}

variable "network_peer" {
  description = "Network Peer configuration details useful for creation"

  type = object({
    name          = string
    provider_type = string
  })
}

variable "aws_config" {
  description = "AWS configuration details useful for network peer creation"

  type = object({
    account_id = optional(string)
    vpc_id     = optional(string)
    cidr       = string
    region     = optional(string)
  })
}

variable "index_keys" {
  description = "index keys"
}

variable "where" {
  description = "WHERE clause"
}

variable "with" {
  description = "WITH clause"

  type = object({
    defer_build   = optional(bool)
    num_replica   = optional(number)
    num_partition = optional(number)
  })
}

variable app_services_cidr {
  description = "App Services CIDR"
}

variable "app_endpoint_oidc" {
    description = "App Endpoint OIDC configuration options"
    type = object({
        client_id       = string
        issuer          = string
        discovery_url   = string
        user_prefix    = string
        roles_claim   = string
        username_claim = string
        register       = bool
    })
}

variable "app_service_log_streaming" {
  description = "App Service Log Streaming Config"
  type = object({
    output_type = string
    credentials = object({
      datadog = optional(object({
        url     = string
        api_key = string
      }))
      dynatrace = optional(object({
        url       = string
        api_token = string
      }))
      elastic = optional(object({
        url      = string
        user     = string
        password = string
      }))
      generic_http = optional(object({
        url      = string
        user     = optional(string)
        password = optional(string)
      }))
      loki = optional(object({
        url      = string
        user     = string
        password = string
      }))
      splunk = optional(object({
        url          = string
        splunk_token = string
      }))
      sumologic = optional(object({
        url = string
      }))
    })
  })
}

variable "app_endpoint_log_streaming_config" {
  description = "App Endpoint Logging Config"
  type = object({
    log_level = string
    log_keys = set(string)
  })
}