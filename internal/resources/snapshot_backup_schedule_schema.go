package resources

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func SnapshotBackupScheduleSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages snapshot backup schedule resource",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The GUID4 ID of the cluster."),
			"interval":        WithDescription(int64Attribute(required), "The interval of the snapshot backup schedule."),
			"retention":       WithDescription(int64Attribute(required), "The retention of the snapshot backup schedule."),
			"start_time":      WithDescription(stringDefaultAttribute(time.Now().Truncate(time.Hour).Format(time.RFC3339), optional, computed), "The start time of the snapshot backup schedule."),
			"copy_to_regions": WithDescription(stringSetAttribute(optional, computed, useStateForUnknown), "The regions to copy the snapshot backup to."),
		},
	}
}
