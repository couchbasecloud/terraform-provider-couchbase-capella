# Set your Capella API key here or through environment variable CAPELLA_AUTHENTICATION_TOKEN
auth_token = ""

import_filter = {
  organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000"
  app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000"
  app_endpoint_name  = "endpoint_1"
  scope              = "scope1"
  collection         = "coll1"
  import_filter      = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }"
}


