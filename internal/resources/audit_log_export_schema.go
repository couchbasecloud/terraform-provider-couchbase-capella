package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AuditLogExportSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                     stringAttribute(computed, useStateForUnknown),
			"organization_id":        stringAttribute(required),
			"project_id":             stringAttribute(required),
			"cluster_id":             stringAttribute(required),
			"audit_log_download_url": stringAttribute(computed),
			"expiration":             stringAttribute(computed),
			"start":                  stringAttribute(computed, optional),
			"end":                    stringAttribute(computed, optional),
			"created_at":             stringAttribute(computed),
			"status":                 stringAttribute(computed),
		},
	}
}
