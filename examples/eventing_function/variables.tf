variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "eventing_function" {
  description = "Eventing function configuration details useful for creation"

  type = object({
    name        = string
    description = optional(string)
    code        = optional(string)
    state       = optional(string)

    event_source = object({
      bucket     = string
      scope      = optional(string)
      collection = optional(string)
    })

    event_metadata_storage = object({
      bucket     = string
      scope      = optional(string)
      collection = optional(string)
    })

    settings = optional(object({
      worker_count           = optional(number)
      script_timeout         = optional(number)
      sql_consistency        = optional(string)
      language_compatibility = optional(string)
      feed_boundary          = optional(string)
      max_timer_context_size = optional(number)
      allow_sync_documents   = optional(bool)
      cursor_aware           = optional(bool)
    }))

    bindings = optional(object({
      buckets = optional(list(object({
        alias      = string
        bucket     = string
        scope      = optional(string)
        collection = optional(string)
        permission = optional(string)
      })))
      urls = optional(list(object({
        alias                    = string
        url                      = string
        allow_cookies            = optional(bool)
        validate_tls_certificate = optional(bool)
        authentication = optional(object({
          type         = string
          username     = optional(string)
          password     = optional(string)
          bearer_token = optional(string)
        }))
      })))
      constants = optional(list(object({
        alias = string
        value = string
      })))
    }))
  })
}
