auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

database_role = {
  name        = "test-role-001"
  description = "A test role with read and write access"
}

access = [
  {
    privileges = ["dataRead"]
    resources = {
      buckets = [{
        name = "travel-sample"
        scopes = [{
          name        = "inventory"
          collections = ["airport", "airline"]
        }]
      }]
    }
  },
  {
    privileges = ["queryManage"]
    resources = {
      buckets = [{
        name = "travel-sample"
        scopes = [{
          name        = "inventory"
          collections = ["hotel", "route"]
        }]
      }]
    }
  }
]

