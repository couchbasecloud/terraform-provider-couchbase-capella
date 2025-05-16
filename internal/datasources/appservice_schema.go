package datasources

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func AppServiceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source retrieves information for an App Service in Capella. App Service is a fully managed application backend designed to provide data synchronization between mobile or IoT applications running Couchbase Lite and your Couchbase Capella database.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the App Service created."},
						"organization_id": schema.StringAttribute{Computed: true, MarkdownDescription: "The GUID4 ID of the organization."},
						"cluster_id":      schema.StringAttribute{Computed: true, MarkdownDescription: "The GUID4 ID of the cluster."},
						"name":            schema.StringAttribute{Computed: true, MarkdownDescription: "Name of the App Service (up to 256 characters)."},
						"description":     schema.StringAttribute{Computed: true, MarkdownDescription: "A description of the App Service (up to 1024 characters)."},
						"nodes":           schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of nodes configured for the App Service."},
						"cloud_provider":  schema.StringAttribute{Computed: true, MarkdownDescription: "The Cloud Service Provider for the App Service."},
						"current_state":   schema.StringAttribute{Computed: true, MarkdownDescription: "The current state of the App Service."},
						"compute": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The CPU and RAM configuration of the App Service.",
							Attributes: map[string]schema.Attribute{
								"cpu": schema.Int64Attribute{Computed: true, MarkdownDescription: "CPU units (cores)."},
								"ram": schema.Int64Attribute{Computed: true, MarkdownDescription: "RAM units (GB)."},
							},
						},
						"version": schema.StringAttribute{Computed: true, MarkdownDescription: "The version of the App Service server. If left empty, it will be defaulted to the latest available version."},
						"audit":   computedAuditAttribute,
					},
				},
			},
		},
	}
}
