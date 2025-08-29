package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func SnapshotBackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages snapshot backup resource",
		Attributes: map[string]schema.Attribute{
			"app_service": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the app service."),
			"backup_id":   WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the snapshot backup."),
			"cluster_id":  WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the cluster."),
			"project_id":  WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the project."),
			"tenant_id":   WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the tenant."),
			"created_at":  WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The RFC3339 timestamp representing the time at which snapshot backup was created."),
			"expiration":  WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The RFC3339 timestamp representing the time at which snapshot backup will expire."),
			"retention":   WithDescription(int64Attribute(required), "The retention of the snapshot backup in hours."),
			"progress": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"status": WithDescription(stringAttribute([]string{computed}), "The status of the snapshot backup. Snapshot backup statuses are 'pending', 'complete', and 'failed'."),
					"time":   WithDescription(stringAttribute([]string{computed}), "The RFC3339 timestamp representing the time at which the status was last updated."),
				},
			},
			"cmek": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the snapshot backup."),
						"provider_id": WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the provider."),
					},
				},
			},

			"server": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The version of the server."),
				},
			},
			"size": WithDescription(int64Attribute(computed, useStateForUnknown), "The size of the snapshot backup in megabytes."),
			"type": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The type of the snapshot backup."),
		},
	}
}
