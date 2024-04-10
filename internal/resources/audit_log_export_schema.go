package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AuditLogExportSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                     stringAttribute([]string{required}),
			"organization_id":        stringAttribute([]string{required}),
			"project_id":             stringAttribute([]string{required}),
			"cluster_id":             stringAttribute([]string{required}),
			"audit_log_download_url": stringAttribute([]string{required}),
			"expiration":             stringAttribute([]string{required}),
			"start":                  stringAttribute([]string{required}),
			"end":                    stringAttribute([]string{required}),
			"created_at":             stringAttribute([]string{required}),
			"status":                 stringAttribute([]string{required}),
		},
	}
}
