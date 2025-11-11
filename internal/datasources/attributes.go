package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Helper functions for common datasource attribute patterns

func requiredString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Required: true,
	}
}

func computedString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Computed: true,
	}
}

func optionalString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Optional: true,
	}
}

func computedInt64() *schema.Int64Attribute {
	return &schema.Int64Attribute{
		Computed: true,
	}
}

func optionalInt64() *schema.Int64Attribute {
	return &schema.Int64Attribute{
		Optional: true,
	}
}

func computedBool() *schema.BoolAttribute {
	return &schema.BoolAttribute{
		Computed: true,
	}
}

func optionalStringSet() *schema.SetAttribute {
	return &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	}
}

func computedStringSet() *schema.SetAttribute {
	return &schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}
}

func computedAudit() *schema.SingleNestedAttribute {
	tempBuilder := capellaschema.NewSchemaBuilder("audit")
	auditAttrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(auditAttrs, "created_at", tempBuilder, &schema.StringAttribute{
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

	computedProjectSnapshot = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The RFC3339 timestamp representing the time at which backup was created.",
			},
			"expiration": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The RFC3339 timestamp representing the time at which the backup will expire.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the snapshot backup resource.",
			},
			"progress": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The progress of the snapshot backup.",
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The status of the snapshot backup.",
					},
					"time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The time of the snapshot backup.",
					},
				},
			},
			"project_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"app_service": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The app service version in the snapshot backup.",
			},
			"cmek": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The CMEK configuration for the snapshot backup.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the CMEK configuration.",
						},
						"provider_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the provider.",
						},
					},
				},
			},

			"cross_region_copies": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The cross region copies of the snapshot backup.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_code": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The region the snapshot backup has been copied to.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the cross region copy.",
						},
						"time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The RFC3339 timestamp representing the time at which the status was last updated.",
						},
					},
				},
			},
			"retention": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The retention time in hours.",
			},
			"server": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version of the server.",
					},
				},
			},
			"database_size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the snapshot backup.",
			},
			"organization_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the snapshot backup.",
			},
		},
	}

	// computedEventAttributes returns a Terraform list nested schema attribute
	// which is configured to be computed and of custom type event.
	computedEventAttributes = schema.ListNestedAttribute{
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "created_by", tempBuilder, &schema.StringAttribute{
		Computed: true,
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_at", tempBuilder, &schema.StringAttribute{
		Computed: true,
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_by", tempBuilder, &schema.StringAttribute{
		Computed: true,
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "version", tempBuilder, &schema.Int64Attribute{
		Computed: true,
	}, "CouchbaseAuditData")

	return &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	}
}
