package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// ProjectEventsSchema returns the schema for the ProjectEvents data source.
func ProjectEventsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      optionalStringAttribute,
			"cluster_ids":     optionalStringSetAttribute,
			"user_ids":        optionalStringSetAttribute,
			"severity_levels": optionalStringSetAttribute,
			"tags":            optionalStringSetAttribute,
			"from":            optionalStringAttribute,
			"to":              optionalStringAttribute,
			"page":            optionalInt64Attribute,
			"per_page":        optionalInt64Attribute,
			"sort_by":         optionalStringAttribute,
			"sort_direction":  optionalStringAttribute,
			"data":            computedEventAttributes,
			"cursor":          computedCursorAttribute,
		},
	}
}
