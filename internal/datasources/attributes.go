package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
				Description: "The RFC3339 timestamp when the resource was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the resource.",
			},
			"modified_at": schema.StringAttribute{
				Computed:    true,
				Description: "The RFC3339 timestamp when the resource was last modified.",
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
					"first": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The base URL, endpoint, and path parameters required to fetch the first page of results.",
					},
					"last": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The base URL, endpoint, and path parameters required to fetch the last page of results.",
					},
					"next": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The base URL, endpoint, and path parameters required to fetch the next page of results. Empty if there is no next page.",
					},
					"previous": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The base URL, endpoint, and path parameters required to fetch the previous page of results. Empty if there is no previous page.",
					},
				},
			},
			"pages": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"last": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of the last page of results.",
					},
					"next": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of the next page of results. Not set on the last page.",
					},
					"page": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Current page of results, starting from page 1.",
					},
					"per_page": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of items displayed in each page.",
					},
					"previous": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of the previous page of results. Not set on the first page.",
					},
					"total_items": schema.Int64Attribute{
						Computed:    true,
						Description: "Total number of items across all pages.",
					},
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
				"alert_key": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Populated on demand based on the Event.Key and select labels in KV.",
				},
				"app_service_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "SyncGatewayID this Event refers to.",
				},
				"app_service_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Name of the sync gateway at the time of event emission.",
				},
				"cluster_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The GUID4 ID of the cluster.",
				},
				"cluster_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Name of the cluster at the time of event emission.",
				},
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "UUID for this instance of an event.",
				},
				"image_url": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "URL to a rendered chart image representing the alert event.",
				},
				"incident_ids": schema.SetAttribute{
					ElementType:         types.StringType,
					Computed:            true,
					MarkdownDescription: "Group events related to an alert incident.",
				},
				"key": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Defines the specific kind of event.",
				},
				"kv": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Key-value pairs for additional event data.",
				},
				"occurrence_count": schema.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "Number of times the alert has fired within this \"incident\".",
				},
				"project_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The GUID4 ID of the project.",
				},
				"project_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Name of the project at the time of event emission.",
				},
				"request_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The request ID for an event.",
				},
				"session_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "ID of the session associated with the user that initiated the request for this event.",
				},
				"severity": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Severity of the event.",
				},
				"source": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Identifies the originator of the event.",
				},
				"summary": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Metadata.SummaryTemplate rendered for this event.",
				},
				"timestamp": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The RFC3339 timestamp when the event was emitted.",
				},
				"user_email": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Email of the associated user at the time of event emission.",
				},
				"user_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "User ID that initiated the request for this event.",
				},
				"user_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Name of the associated user at the time of event emission.",
				},
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

func computedMaptAttribute(T attr.Type) schema.MapAttribute {
	return schema.MapAttribute{
		ElementType: T,
		Computed:    true,
	}
}
