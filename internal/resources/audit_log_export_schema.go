package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AuditLogExportSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                    stringAttribute(computed, useStateForUnknown),
			"organization_id":       stringAttribute(required, requiresReplace),
			"project_id":            stringAttribute(required, requiresReplace),
			"cluster_id":            stringAttribute(required, requiresReplace),
			"auditlog_download_url": stringAttribute(computed, requiresReplace),
			"expiration":            stringAttribute(computed, requiresReplace),
			"start":                 stringAttribute(computed, optional, requiresReplace),
			"end":                   stringAttribute(computed, optional, requiresReplace),
			"created_at":            stringAttribute(computed, requiresReplace),
			"status":                stringAttribute(computed, requiresReplace),
		},
	}
}
