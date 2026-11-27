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

	publicAccess := stringAttribute([]string{required})
	publicAccess.MarkdownDescription = "- The public access mode to apply: `restricted` or `unrestricted`."
	capellaschema.AddAttr(attrs, "public_access", appServicePublicAccessBuilder, publicAccess)

	capellaschema.AddAttr(attrs, "public_access_status", appServicePublicAccessBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "Manages the public access mode of an App Service. " +
			"Destroying this resource does not change public access. " +
			"To open public access, set `public_access` to `unrestricted` and apply first. " +
			"Changing to `unrestricted` deletes all public allowed CIDRs.",
		Attributes: attrs,
	}
}
