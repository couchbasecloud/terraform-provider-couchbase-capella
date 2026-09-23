auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

database_credential = {
  database_credential_name = "test_db_user"
  password                 = "Secret12$#"
}

advanced_database_credential = {
  database_credential_name = "test_advanced_db_user"
  password                 = "Secret12$#"
}

user_roles = ["developer", "bucket_admin"]

access = [
  {
    privileges = ["data_writer"]
    resources = {
      buckets = [{
        name = "new_terraform_bucket"
        scopes = [
          {
            name = "_default"
          }
        ]
      }]
    }
  },
  {
    privileges = ["data_reader"]
  }
]
