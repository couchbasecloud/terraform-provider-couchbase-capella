package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotRestoresBuilder = capellaschema.NewSchemaBuilder("snapshotRestores")

func SnapshotRestoresSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "cluster_id", snapshotRestoresBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", snapshotRestoresBuilder, requiredString())
	capellaschema.AddAttr(attrs, "organization_id", snapshotRestoresBuilder, requiredString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", snapshotRestoresBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "created_at", snapshotRestoresBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "restore_to", snapshotRestoresBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "snapshot", snapshotRestoresBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", snapshotRestoresBuilder, computedString())
	capellaschema.AddAttr(attrs, "data", snapshotRestoresBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	filterAttrs := make(map[string]schema.Attribute)
	filterAttrs["name"] = schema.StringAttribute{
		MarkdownDescription: "The name of the attribute to filter.",
		Optional:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("status"),
		},
	}
	filterAttrs["values"] = schema.SetAttribute{
		MarkdownDescription: "List of values to match against.",
		Optional:            true,
		ElementType:         types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all snapshot restore information for a cluster.",
		Attributes:          attrs,

		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				MarkdownDescription: "Filter criteria for the Cloud Snapshot Restores. Only filtering by Snapshot Restore Status is supported.",
				Attributes:          filterAttrs,
			},
		},
	}
}
