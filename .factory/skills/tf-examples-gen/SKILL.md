---
name: tf-examples-gen
description: generate terraform HCL examples for a feature.
---

# Terraform Examples Generator

## Instructions

- HCL must be in examples/.  for example if the feature is buckets create folder examples/buckets/
- create main.tf

```
terraform {
  required_providers {
    couchbase-capella = {
      source = "couchbasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
  authentication_token = var.auth_token
}
```

- create variables.tf

  this should have variables needed for the resource and datasource.

  all examples must have these variables

```
variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}
```

- create terraform.template.tfvars file

  this should have placeholder values for the variables. for variables like
  organization_id use "<organization_id>".  for auth_token use "<v4-api-key-secret>"

  for other variables use values in ../couchbase-cloud/cmd/cp-open-api/specs/examples


- create a create_<feature>.tf file

  this creates a resource for example

```

resource "couchbase-capella_<feature>" "new_<feature>" {
  // required and optional arguments for the resource
}
```

  the resource should have required and optional arguments derived from the
  schema in internal/resources/<feature>_schema.go

- determine if the feature has a datasource to list resources.  look for a file
  in internal/datasources/ with the plural name of the feature.
  for example if the feature is Buckets then look for buckets.go in internal/datasources/.

  if there is a datasource to list resources then create a list_<feature>.tf file
  it should look like this:

```
data "couchbase-capella_<feature_plural>" "list_<feature_plural>" {
  // required and optional arguments for the datasource
}
```

 the required and optional arguments should be derived from the schema in internal/datasources/<feature_plural>_schema.go

 if there is no datasource to list resources then skip this step.

 look for a datasource to get a specific resource.
 look for a file in internal/datasources/ with the singular name of the feature.
 for example if the feature is Buckets then look for bucket.go in internal/datasources/

 if there is a datasource to get a specific resource then create a get_<feature>.tf file
 it should look like this:

```
data "couchbase-capella_<feature>" "get_<feature>" {
  // required and optional arguments for the datasource
}
```

 the required and optional arguments should be derived from the schema in internal/datasources/<feature>_schema.go

 if there is no datasource to get a specific resource then skip this step.

- run terraform validate to ensure the examples are valid terraform code.
  fix errors until terraform validate passes.

  do not run terraform init