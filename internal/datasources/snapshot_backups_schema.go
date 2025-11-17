package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotBackupsBuilder = capellaschema.NewSchemaBuilder("snapshotBackups")

func SnapshotBackupsSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", snapshotBackupsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", snapshotBackupsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotBackupsBuilder, requiredString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "created_at", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "expiration", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "id", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "retention", snapshotBackupsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "size", snapshotBackupsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "type", snapshotBackupsBuilder, computedString())

	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(progressAttrs, "time", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "progress", snapshotBackupsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: progressAttrs,
	})

	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(cmekAttrs, "provider_id", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cmek", snapshotBackupsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: cmekAttrs,
		},
	})

	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "server", snapshotBackupsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: serverAttrs,
	})

	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cross_region_copies", snapshotBackupsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: crossRegionCopiesAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "data", snapshotBackupsBuilder, &schema.ListNestedAttribute{
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
		MarkdownDescription: "The snapshot backups data source retrieves snapshot backups associated with a bucket for an operational cluster.",
		Attributes:          attrs,
		Blocks: map[string]schema.Block{
			"filters": schema.SingleNestedBlock{
				Attributes: filterAttrs,
			},
		},
	}
}
