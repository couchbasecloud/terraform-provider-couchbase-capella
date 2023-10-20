package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UserSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                   stringAttribute(computed),
			"name":                 stringAttribute(optional),
			"status":               stringAttribute(computed),
			"inactive":             boolAttribute(computed),
			"email":                stringAttribute(required),
			"organization_id":      stringAttribute(required),
			"organization_roles":   stringListAttribute(required),
			"last_login":           stringAttribute(computed),
			"region":               stringAttribute(computed),
			"time_zone":            stringAttribute(computed),
			"enable_notifications": boolAttribute(computed),
			"expires_at":           stringAttribute(computed),
			"resources": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":  stringAttribute(optional),
						"id":    stringAttribute(required),
						"roles": stringListAttribute(required),
					},
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
