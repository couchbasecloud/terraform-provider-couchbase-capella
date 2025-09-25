package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func SnapshotBackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages snapshot backup resource",
		Attributes: map[string]schema.Attribute{
			"app_service":     WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the app service."),
			"id":              WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the snapshot backup."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the cluster."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the project."),
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the organization."),
			"created_at":      WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The RFC3339 timestamp representing the time at which snapshot backup was created."),
			"expiration":      WithDescription(stringAttribute([]string{computed}), "The RFC3339 timestamp representing the time at which snapshot backup will expire."),
			"retention":       WithDescription(int64Attribute(optional, computed), "The retention of the snapshot backup in hours."),
			"regions_to_copy": WithDescription(stringListAttribute(optional), "The regions to copy the snapshot backup to."),
			"cross_region_copies": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_code": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The region the snapshot backup has been copied to."),
						"status":      WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The status of the cross region copy."),
						"time":        WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The RFC3339 timestamp representing the time at which the status was last updated."),
					},
				},
			},
			"progress": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"status": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The status of the snapshot backup. Snapshot backup statuses are 'pending', 'complete', and 'failed'."),
					"time":   WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The RFC3339 timestamp representing the time at which the status was last updated."),
				},
			},
			"cmek": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the snapshot backup."),
						"provider_id": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the provider."),
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
