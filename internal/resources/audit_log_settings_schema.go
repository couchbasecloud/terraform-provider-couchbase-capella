package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AuditLogSettingsSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource to manage audit log configuration settings for an operational cluster. These settings control which audit events are logged and which users are excluded from audit logging.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required}),
				"The GUID4 ID of the organization."),
			"project_id": WithDescription(stringAttribute([]string{required}),
				"The GUID4 ID of the project."),
			"cluster_id": WithDescription(stringAttribute([]string{required}),
				"The GUID4 ID of the cluster to configure audit log settings."),
			"audit_enabled": WithDescription(boolAttribute(computed, optional),
				"Determines whether audit logging is enabled or not on the cluster. Set to 'true' to enable audit logging."),
			"enabled_event_ids": schema.SetAttribute{
				Computed:            true,
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of audit event IDs to enable for logging. These IDs correspond to specific types of events that will be recorded in the audit log. Use the audit_log_event_ids data source to get the list of available event IDs.",
			},
			"disabled_users": schema.SetNestedAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "List of users whose actions will be excluded from audit logging.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": WithDescription(stringAttribute([]string{required}),
							"The authentication domain of the user to exclude. Specifies whether the user is local or external."),
						"name": WithDescription(stringAttribute([]string{required}),
							"The username of the user to exclude from audit logging."),
					},
				},
			},
		},
	}
}
