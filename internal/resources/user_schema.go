package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UserSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                   stringAttribute(computed),
			"name":                 stringAttribute(optional, computed),
			"status":               stringAttribute(computed),
			"inactive":             boolAttribute(computed),
			"email":                stringAttribute(required, requiresReplace),
			"organization_id":      stringAttribute(required, requiresReplace),
			"organization_roles":   stringListAttribute(required),
			"last_login":           stringAttribute(computed),
			"region":               stringAttribute(computed),
			"time_zone":            stringAttribute(computed),
			"enable_notifications": boolAttribute(computed),
			"expires_at":           stringAttribute(computed),
			"resources": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":  stringAttribute(optional, computed),
						"id":    stringAttribute(required),
						"roles": stringListAttribute(required),
					},
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
