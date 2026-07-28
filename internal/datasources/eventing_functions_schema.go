package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// eventingFunctionStatuses are the states an eventing function can be in. They double as the
// accepted values for the status filter on the list eventing functions data source.
var eventingFunctionStatuses = []string{
	"deployed",
	"deploying",
	"undeployed",
	"undeploying",
	"paused",
	"pausing",
}

// getOneEventingFunctionAttrs returns the computed attributes describing a single eventing function
// entry in the list data source. The nested object attributes are shared with the single eventing
// function data source.
func getOneEventingFunctionAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "name", eventingFunctionBuilder, computedString())
	capellaschema.AddAttr(attrs, "description", eventingFunctionBuilder, computedString())
	capellaschema.AddAttr(attrs, "status", eventingFunctionBuilder, computedString())
	capellaschema.AddAttr(attrs, "code", eventingFunctionBuilder, &schema.StringAttribute{
		Computed:  true,
		Sensitive: true,
	})

	capellaschema.AddAttr(attrs, "event_source", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionKeyspaceAttrs(),
	})
	capellaschema.AddAttr(attrs, "event_metadata_storage", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionKeyspaceAttrs(),
	})
	capellaschema.AddAttr(attrs, "settings", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionSettingsAttrs(),
	})
	capellaschema.AddAttr(attrs, "bindings", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionBindingsAttrs(),
	})
	return attrs
}

// EventingFunctionsSchema returns the schema for the list eventing functions data source.
func EventingFunctionsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", eventingFunctionBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", eventingFunctionBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", eventingFunctionBuilder, requiredUUIDString())

	// status is an optional filter passed through to the list endpoint's status query parameter.
	// When omitted, eventing functions in every state are returned.
	capellaschema.AddAttr(attrs, "status", eventingFunctionBuilder, &schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
			setvalidator.ValueStringsAre(stringvalidator.OneOf(eventingFunctionStatuses...)),
		},
	})

	capellaschema.AddAttr(attrs, "eventing_functions", eventingFunctionBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getOneEventingFunctionAttrs(),
		},
	})

	return schema.Schema{
		MarkdownDescription: "The eventing functions data source retrieves the eventing functions in a cluster, optionally filtered by one or more states.",
		Attributes:          attrs,
	}
}
