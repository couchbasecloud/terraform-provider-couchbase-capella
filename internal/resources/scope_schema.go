package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ScopeSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute(required, requiresReplace),
			"project_id":      stringAttribute(required, requiresReplace),
			"cluster_id":      stringAttribute(required, requiresReplace),
			"bucket_id":       stringAttribute(required, requiresReplace),
			"name":            stringAttribute(required, requiresReplace),
			"uid":             stringAttribute(computed, useStateForUnknown),
			"collections": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"max_ttl": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"uid": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}

}
