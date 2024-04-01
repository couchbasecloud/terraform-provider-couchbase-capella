package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UserSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                   stringAttribute([]string{computed}),
			"name":                 stringAttribute([]string{optional, computed}),
			"status":               stringAttribute([]string{computed}),
			"inactive":             boolAttribute(computed),
			"email":                stringAttribute([]string{required, requiresReplace}),
			"organization_id":      stringAttribute([]string{required, requiresReplace}),
			"organization_roles":   stringListAttribute(required),
			"last_login":           stringAttribute([]string{computed}),
			"region":               stringAttribute([]string{computed}),
			"time_zone":            stringAttribute([]string{computed}),
			"enable_notifications": boolAttribute(computed),
			"expires_at":           stringAttribute([]string{computed}),
			"resources": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":  stringDefaultAttribute("project", optional, computed),
						"id":    stringAttribute([]string{required}),
						"roles": stringSetAttribute(required),
					},
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
