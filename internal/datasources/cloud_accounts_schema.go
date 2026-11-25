package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var cloudAccountsBuilder = capellaschema.NewSchemaBuilder("cloudAccounts", "CloudAccounts")

// CloudAccountsSchema returns the schema for the cloud_accounts data source.
func CloudAccountsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", cloudAccountsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "aws_capella_account", cloudAccountsBuilder, computedString())
	capellaschema.AddAttr(attrs, "gcp_capella_project", cloudAccountsBuilder, computedString())
	capellaschema.AddAttr(attrs, "azure_capella_subscription", cloudAccountsBuilder, computedString())
	capellaschema.AddAttr(attrs, "azure_capella_subscription_id", cloudAccountsBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve Capella's cloud accounts for an organization, to grant them access in your own cloud account.",
		Attributes:          attrs,
	}
}
