auth_token               = "bGZrdkYxMXRnSXlORXF1NmdKZWU3YVlXSVVVWVdNQm06dHRkdGdpaFRGaVpiTW9CUlpkMFUhSGlsIVQlS2E4a2lTQEtrTG9tWUBhcktOSnNjWkpRNnh0dTlBZ1JTTVdSMg=="
organization_id          = "1a3c4544-772e-449e-9996-1203e7020b96"
project_id               = "73a26cf0-2c4a-43ab-904f-9d86e595bbb5"
host                     = "https://cloudapi.dev.nonprod-project-avengers.com"
database_credential_name = "test_db_user"
cluster_id               = "7d726313-cde1-4439-88c4-6b49e5dd49ac"
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
