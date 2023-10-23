auth_token               = "<v4-api-key-secret>"
organization_id          = "<organization_id>"
project_id               = "<project_id>"
host                     = "https://cloudapi.cloud.couchbase.com"
database_credential_name = "test_db_user"
cluster_id               = "<cluster_id>"
password                 = "Secret12$#"
access = [
  {
    privileges = ["data_writer"]
    resources = {
      buckets = [{
        name = "new_terraform_bucket"
        scopes = [
          {
            name        = "_default"
            collections = ["_default"]
          }
        ]
      }]
    }
  },
  {
    privileges = ["data_reader"]
  }
]
