package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// DatabaseCredentialSchema defines the schema for the terraform provider resource - "DatabaseCredential".
// This terraform resource directly maps to the database credential created for a Capella cluster.
// DatabaseCredential resource supports Create, Destroy, Read, Import and List operations.
func DatabaseCredentialSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":            stringAttribute(required),
			"password":        stringAttribute(optional, computed, sensitive),
			"organization_id": stringAttribute(required),
			"project_id":      stringAttribute(required),
			"cluster_id":      stringAttribute(required),
			"audit":           computedAuditAttribute(),
			"access": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"privileges": stringListAttribute(required),
						"resources": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"buckets": schema.ListNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": stringAttribute(required),
											"scopes": schema.ListNestedAttribute{
												Optional: true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"name":        stringAttribute(required),
														"collections": stringListAttribute(optional),
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
			},
		},
	}
}
