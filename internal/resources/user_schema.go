package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UserSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                   stringAttribute(computed, requiresReplace),
			"name":                 stringAttribute(computed, requiresReplace),
			"status":               stringAttribute(computed, requiresReplace),
			"inactive":             boolAttribute(computed, requiresReplace),
			"email":                stringAttribute(required, requiresReplace),
			"organization_id":      stringAttribute(required, requiresReplace),
			"organization_roles":   stringListAttribute(required, requiresReplace),
			"last_login":           stringAttribute(computed, requiresReplace),
			"region":               stringAttribute(computed, requiresReplace),
			"time_zone":            stringAttribute(computed, requiresReplace),
			"enable_notifications": boolAttribute(computed, requiresReplace),
			"expires_at":           stringAttribute(computed, requiresReplace),
			"resources": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":  stringAttribute(optional, requiresReplace),
						"id":    stringAttribute(required, requiresReplace),
						"roles": stringListAttribute(required, requiresReplace),
					},
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
