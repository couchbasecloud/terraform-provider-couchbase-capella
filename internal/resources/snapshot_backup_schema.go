package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotBackupBuilder = capellaschema.NewSchemaBuilder("snapshotBackup", "cloudSnapshotBackup")

func SnapshotBackupSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "cluster_id", snapshotBackupBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", snapshotBackupBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "organization_id", snapshotBackupBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "created_at", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "expiration", snapshotBackupBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "retention", snapshotBackupBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(attrs, "regions_to_copy", snapshotBackupBuilder, stringSetAttribute(optional, useStateForUnknown))
	capellaschema.AddAttr(attrs, "size", snapshotBackupBuilder, int64Attribute(computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "type", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "restore_times", snapshotBackupBuilder, numberAttribute(optional))
	capellaschema.AddAttr(attrs, "cross_region_restore_preference", snapshotBackupBuilder, stringListAttribute(optional))

	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "cross_region_copies", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: crossRegionCopiesAttrs,
		},
	})

	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(progressAttrs, "time", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "progress", snapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: progressAttrs,
	})

	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(cmekAttrs, "provider_id", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "cmek", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: cmekAttrs,
		},
	})

	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", snapshotBackupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "server", snapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: serverAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "Manages snapshot backup resource associated with a Capella cluster.",
		Attributes:          attrs,
	}
}
