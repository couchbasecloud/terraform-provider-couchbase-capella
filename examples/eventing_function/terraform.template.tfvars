auth_token = "<v4-api-key-secret>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

eventing_function = {
  name        = "my_function"
  description = "Replicates document mutations to a downstream collection."
  state       = "deployed"

  code = <<-EOT
    function OnUpdate(doc, meta, xattrs) {
      log("Doc created/updated", meta.id);
    }

    function OnDelete(meta, options) {
      log("Doc deleted/expired", meta.id);
    }
  EOT

  event_source = {
    bucket     = "travel-sample"
    scope      = "inventory"
    collection = "airline"
  }

  event_metadata_storage = {
    bucket     = "metadata"
    scope      = "_default"
    collection = "_default"
  }

  settings = {
    worker_count    = 1
    script_timeout  = 60
    sql_consistency = "none"
    feed_boundary   = "from_now"
  }

  bindings = {
    constants = [
      {
        alias = "maxRetries"
        value = "3"
      }
    ]
  }
}
