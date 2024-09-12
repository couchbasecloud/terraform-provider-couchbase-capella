package resources

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func GsiSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"bucket_name":     stringAttribute([]string{required, requiresReplace}),
			"scope_name": stringDefaultAttribute(
				"_default", optional, computed, useStateForUnknown, requiresReplace,
			),
			"collection_name": stringDefaultAttribute(
				"_default", optional, computed, useStateForUnknown, requiresReplace,
			),
			"index_name":   stringAttribute([]string{optional, useStateForUnknown, requiresReplace}),
			"is_primary":   boolAttribute(optional),
			"index_keys":   stringSetAttribute(optional, useStateForUnknown, requiresReplace),
			"where":        stringAttribute([]string{optional, useStateForUnknown, requiresReplace}),
			"partition_by": stringAttribute([]string{optional, useStateForUnknown, requiresReplace}),
			"with": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"defer_build":    boolDefaultAttribute(false, optional, computed),
					"num_replica":    int64DefaultAttribute(0, optional, computed),
					"num_partitions": int64DefaultAttribute(1, optional, computed),
				},
			},
			"build_indexes": stringSetAttribute(optional),
		},
	}
}
