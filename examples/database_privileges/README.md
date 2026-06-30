# Capella Database Privileges Example

This example lists the available Capella database privileges for a given cluster.

Database privileges define what operations can be performed on data at the bucket, scope, and collection level. They can be assigned to database credentials and database user roles.

## Prerequisites

- A Capella organization, project, and provisioned cluster.
- An API key with at least **Project Viewer** access.

## Usage

Copy the template variables file and fill in your values:

```bash
cp terraform.template.tfvars terraform.tfvars
```

Edit `terraform.tfvars` with your organization, project, and cluster IDs.

Initialize and apply:

```bash
terraform init
terraform apply
```

## Example Output

```
database_privileges = {
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  data = [
    {
      group = "Data"
      name  = "dataRead"
      resources = {
        buckets = [
          {
            name = "*"
            scopes = [
              {
                collections = ["*"]
                name        = "*"
              }
            ]
          }
        ]
      }
    },
    {
      group     = "Query"
      name      = "queryIndex"
      resources = {
        buckets = [
          {
            name = "*"
            scopes = [
              {
                collections = ["*"]
                name        = "*"
              }
            ]
          }
        ]
      }
    },
  ]
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

