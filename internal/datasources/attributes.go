package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	// computedStringAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedStringAttribute = schema.StringAttribute{
		Computed: true,
	}

	// requiredStringAttribute returns a Terraform schema attribute
	// which is configured to be required.
	requiredStringAttribute = schema.StringAttribute{
		Required: true,
	}

	// optionalStringAttribute returns a Terraform schema attribute
	// which is configured to be optional.
	optionalStringAttribute = schema.StringAttribute{
		Optional: true,
	}

	// computedBoolAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedBoolAttribute = schema.BoolAttribute{
		Computed: true,
	}

	// optionalInt64Attribute returns a Terraform schema attribute
	// which is configured to be optional.
	optionalInt64Attribute = schema.Int64Attribute{
		Optional: true,
	}

	// computedBoolAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedInt64Attribute = schema.Int64Attribute{
		Computed: true,
	}

	// computedBoolAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedFloat64Attribute = schema.Float64Attribute{
		Computed: true,
	}

	// computedListAttribute returns a Terraform list schema attribute
	// which is configured to be computed and of type string.
	computedListAttribute = schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}

	// computedIntSetAttribute returns a Terraform set schema attribute.
	requiredStringSetAttribute = schema.SetAttribute{
		ElementType: types.StringType,
		Required:    true,
	}

	// computedIntSetAttribute returns a Terraform list schema attribute
	// which is configured to be computed and of type int64.
	computedIntSetAttribute = schema.SetAttribute{
		ElementType: types.Int64Type,
		Computed:    true,
	}

	// computedStringSetAttribute returns a Terraform set schema attribute
	// which is configured to be computed and of type string.
	computedStringSetAttribute = schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}

	// optionalStringSetAttribute returns a Terraform set schema attribute
	// which is configured to be optional and of type string.
	optionalStringSetAttribute = schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	}

	// computedAuditAttribute returns a SingleNestedAttribute to
	// represent couchbase audit data using terraform schema types.
	computedAuditAttribute = schema.SingleNestedAttribute{
		Description: "Couchbase audit data.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the resource.",
			},
			"modified_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the resource was last modified.",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who last modified the resource.",
			},
			"version": schema.Int64Attribute{
				Computed: true,
				Description: "The version of the document. " +
					"This value is incremented each time the resource is modified.",
			},
		},
	}

	// computedCursorAttribute returns a Terraform single nested schema attribute
	// which is configured to be computed and of custom type cursor.
	computedCursorAttribute = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"hrefs": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"first":    computedStringAttribute,
					"last":     computedStringAttribute,
					"next":     computedStringAttribute,
					"previous": computedStringAttribute,
				},
			},
			"pages": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"last":        computedInt64Attribute,
					"next":        computedInt64Attribute,
					"page":        computedInt64Attribute,
					"per_page":    computedInt64Attribute,
					"previous":    computedInt64Attribute,
					"total_items": computedInt64Attribute,
				},
			},
		},
	}

	// computedEventAttributes returns a Terraform list nested schema attribute
	// which is configured to be computed and of custom type event.
	computedEventAttributes = schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"alert_key":        computedStringAttribute,
				"app_service_id":   computedStringAttribute,
				"app_service_name": computedStringAttribute,
				"cluster_id":       computedStringAttribute,
				"cluster_name":     computedStringAttribute,
				"id":               computedStringAttribute,
				"image_url":        computedStringAttribute,
				"incident_ids":     computedStringSetAttribute,
				"key":              computedStringAttribute,
				"kv":               computedStringAttribute,
				"occurrence_count": computedInt64Attribute,
				"project_id":       computedStringAttribute,
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
		},
	}

	computedGsiAttributes = schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"index_name": computedStringAttribute,
				"definition": computedStringAttribute,
			},
		},
	}
)
