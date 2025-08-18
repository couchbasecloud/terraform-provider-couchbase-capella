auth_token         = "v4-api-key-secret"
organization_id    = "<organization_id>"
project_id         = "<project_id>"
cluster_id         = "<cluster_id>"
app_service_id     = "<app_service_id>"
app_endpoint_name  = "<app_endpoint_name>"
scope             = "<scope_name>"
collection        = "<collection_name>"

access_control_function = {
  organization_id         = "<organization_id>"
  project_id              = "<project_id>"
  cluster_id              = "<cluster_id>"
  app_service_id          = "<app_service_id>"
  app_endpoint_name       = "<app_endpoint_name>"
  scope                   = "<scope_name>"
  collection              = "<collection_name>"
  access_control_function = "function(doc, oldDoc) { if (doc.type === 'public') { channel('public'); } else if (doc.owner) { channel('user:' + doc.owner); } }"
}

# Example access control functions:
# 
# Default collection function:
# access_control_function = "function(doc){channel(doc.channels);}"
#
# Named collection function:
# access_control_function = "function(doc){channel('myCollectionName');}"
#
# Custom access control with user-based channels:
# access_control_function = "function(doc, oldDoc) { if (doc.type === 'public') { channel('public'); } else if (doc.owner) { channel('user:' + doc.owner); } }"
#
# Role-based access control:
# access_control_function = "function(doc, oldDoc) { if (doc.roles) { for (var role of doc.roles) { channel('role:' + role); } } }" 