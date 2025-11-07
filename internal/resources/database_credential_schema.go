package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var databaseCredentialBuilder = capellaschema.NewSchemaBuilder("databaseCredential")

// DatabaseCredentialSchema defines the schema for the terraform provider resource - "DatabaseCredential".
// This terraform resource directly maps to the database credential created for a Capella cluster.
// DatabaseCredential resource supports Create, Destroy, Read, Import and List operations.
func DatabaseCredentialSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", databaseCredentialBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "name", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "organization_id", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["password"] = WithDescription(stringAttribute([]string{optional, computed, sensitive, useStateForUnknown}), "A password associated with the database credential.")
	capellaschema.AddAttr(attrs, "audit", databaseCredentialBuilder, computedAuditAttribute())
	attrs["access"] = schema.SetNestedAttribute{
		Required:            true,
		MarkdownDescription: "Describes the access information of the database credential.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"privileges": WithDescription(stringSetAttribute(required), "The privileges field in this API represents the privilege level for users."),
				"resources": schema.SingleNestedAttribute{
					Optional:            true,
					MarkdownDescription: "The resources for which access will be granted on. Leaving this empty will grant access to all buckets.",
					Attributes: map[string]schema.Attribute{
						"buckets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": WithDescription(stringAttribute([]string{required}), "The name of the bucket."),
									"scopes": schema.SetNestedAttribute{
										Optional:            true,
										MarkdownDescription: "The scopes under a bucket.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name":        WithDescription(stringAttribute([]string{required}), "The name of the scope."),
												"collections": WithDescription(stringSetAttribute(optional), "The collections under a scope."),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return schema.Schema{
		MarkdownDescription: "Resource to create and manage a database credential for a cluster. Database credentials provide programmatic and application-level access to data on a database.",
		Attributes:          attrs,
	}
}
