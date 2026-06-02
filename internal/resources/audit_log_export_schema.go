package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogExportBuilder = capellaschema.NewSchemaBuilder("auditLogExport", "clusterAuditLogExport")

type rfc3339TimestampValidator struct {
	attributeName string
}

func (v rfc3339TimestampValidator) Description(_ context.Context) string {
	return fmt.Sprintf("%s must be a valid RFC3339 timestamp", v.attributeName)
}

func (v rfc3339TimestampValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v rfc3339TimestampValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if _, err := time.Parse(time.RFC3339, req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid RFC3339 Timestamp",
			fmt.Sprintf("%s must be a valid RFC3339 timestamp", v.attributeName),
		)
	}
}

func AuditLogExportSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", auditLogExportBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", auditLogExportBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", auditLogExportBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", auditLogExportBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "audit_log_download_url", auditLogExportBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "expiration", auditLogExportBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "start", auditLogExportBuilder, stringAttribute([]string{required, requiresReplace}, rfc3339TimestampValidator{attributeName: "start"}))
	capellaschema.AddAttr(attrs, "end", auditLogExportBuilder, stringAttribute([]string{required, requiresReplace}, rfc3339TimestampValidator{attributeName: "end"}))
	capellaschema.AddAttr(attrs, "created_at", auditLogExportBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "status", auditLogExportBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage audit log exports for an operational cluster. This allows you to export audit logs for a specific time period and download them for analysis. Audit Logs for the last 30 days can be requested, otherwise they are purged. A pre-signed URL to a s3 bucket location is returned, which is used to download these audit logs.",
		Attributes:          attrs,
	}
}
