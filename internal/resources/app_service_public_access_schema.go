package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePublicAccessBuilder = capellaschema.NewSchemaBuilder("appServicePublicAccess")

// AppServicePublicAccessSchema returns the schema for the app_service_public_access resource.
func AppServicePublicAccessSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePublicAccessBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServicePublicAccessBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePublicAccessBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePublicAccessBuilder, requiredUUIDStringAttribute())

	capellaschema.AddAttr(attrs, "public_access", appServicePublicAccessBuilder, stringAttribute([]string{required}))

	capellaschema.AddAttr(attrs, "public_access_status", appServicePublicAccessBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "Manages the public access mode of an App Service. " +
			"Destroying this resource removes it from the Terraform state only, and does not change public access.",
		Attributes: attrs,
	}
}
