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
			"email":                stringAttribute(required),
			"organization_id":      stringAttribute(required),
			"organization_roles":   stringSetAttribute(required),
			"last_login":           stringAttribute(computed),
			"region":               stringAttribute(computed),
			"time_zone":            stringAttribute(computed),
			"enable_notifications": boolAttribute(computed),
			"expires_at":           stringAttribute(computed),
			"resources": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":  stringDefaultAttribute("project", optional, computed),
						"id":    stringAttribute(required),
						"roles": stringSetAttribute(required),
					},
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
