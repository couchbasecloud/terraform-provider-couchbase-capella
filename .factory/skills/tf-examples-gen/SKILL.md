---
name: tf-examples-gen
description: Generate terraform HCL examples for a feature.
---

# Terraform Examples Generator

## Instructions

- HCL must be in examples/.  For example if the feature is buckets create folder examples/buckets/
- Create main.tf

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

If main.tf exists and terraform and provider blocks, skip this step.

- Create variables.tf

  This should have variables needed for the resource and datasource.

  All examples must have these variables

```
variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}
```

If variables.tf exists and has the necessary variables, skip this step.

- Create terraform.template.tfvars file

  This should have placeholder values for the variables. For variables like
  organization_id use "<organization_id>".  For auth_token use "<v4-api-key-secret>"

  For other variables use values in ../couchbase-cloud/cmd/cp-open-api/specs/examples

If terraform.template.tfvars exists and has the necessary variables, skip this step.

- Create a create_<feature>.tf file

  This creates a resource for example

```

resource "couchbase-capella_<feature>" "new_<feature>" {
  // required and optional arguments for the resource
}
```

  The resource should have required and optional arguments derived from the
  schema in internal/resources/<feature>_schema.go

If create_<feature>.tf exists and has the resource block, skip this step.

- Determine if the feature has a datasource to list resources.  Look for a file
  in internal/datasources/ with the plural name of the feature.
  For example if the feature is Buckets then look for buckets.go in internal/datasources/.

  If there is a datasource to list resources then create a list_<feature>.tf file
  it should look like this:

```
data "couchbase-capella_<feature_plural>" "list_<feature_plural>" {
  // required and optional arguments for the datasource
}
```

 The required and optional arguments should be derived from the schema in internal/datasources/<feature_plural>_schema.go

 If there is no datasource to list resources then skip this step.


- Look for a datasource to get a specific resource.
  Look for a file in internal/datasources/ with the singular name of the feature.
  For example if the feature is Buckets then look for bucket.go in internal/datasources/

  If there is a datasource to get a specific resource then create a get_<feature>.tf file
  it should look like this:

```
data "couchbase-capella_<feature>" "get_<feature>" {
  // required and optional arguments for the datasource
}
```

 The required and optional arguments should be derived from the schema in internal/datasources/<feature>_schema.go

 If there is no datasource to get a specific resource then skip this step.

- Run terraform validate to ensure the examples are valid terraform code.
  Fix errors until terraform validate passes.

  Do not run terraform init
