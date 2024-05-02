package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AuditLogExportSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                     stringAttribute([]string{computed, useStateForUnknown}),
			"organization_id":        stringAttribute([]string{required, requiresReplace}),
			"project_id":             stringAttribute([]string{required, requiresReplace}),
			"cluster_id":             stringAttribute([]string{required, requiresReplace}),
			"audit_log_download_url": stringAttribute([]string{computed, requiresReplace}),
			"expiration":             stringAttribute([]string{computed, requiresReplace}),
			"start":                  stringAttribute([]string{required, requiresReplace}),
			"end":                    stringAttribute([]string{required, requiresReplace}),
			"created_at":             stringAttribute([]string{computed, requiresReplace}),
			"status":                 stringAttribute([]string{computed, requiresReplace}),
		},
	}
}
