package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// ProjectEventSchema returns the schema for the ProjectEvent data source.
func ProjectEventSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":               requiredStringAttribute,
			"organization_id":  requiredStringAttribute,
			"project_id":       requiredStringAttribute,
			"alert_key":        computedStringAttribute,
			"app_service_id":   computedStringAttribute,
			"app_service_name": computedStringAttribute,
			"cluster_id":       computedStringAttribute,
			"cluster_name":     computedStringAttribute,
			"image_url":        computedStringAttribute,
			"incident_ids":     computedStringSetAttribute,
			"key":              computedStringAttribute,
			"kv":               computedStringAttribute,
			"occurrence_count": computedInt64Attribute,
			"project_name":     computedStringAttribute,
			"request_id":       computedStringAttribute,
			"session_id":       computedStringAttribute,
			"severity":         computedStringAttribute,
			"source":           computedStringAttribute,
			"summary":          computedStringAttribute,
			"timestamp":        computedStringAttribute,
			"user_email":       computedStringAttribute,
			"user_id":          computedStringAttribute,
			"user_name":        computedStringAttribute,
		},
	}
}
