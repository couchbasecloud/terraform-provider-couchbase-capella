package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/docs"
)

// SchemaAttribute is a type constraint for supported attribute types
type SchemaAttribute interface {
	*schema.StringAttribute | *schema.Int64Attribute | *schema.BoolAttribute |
		*schema.SetAttribute | *schema.Float64Attribute | *schema.NumberAttribute |
		*schema.ListAttribute | *schema.SingleNestedAttribute | *schema.ObjectAttribute
}

// SchemaBuilder provides methods for building resource and data source schemas with OpenAPI integration.
// Each resource or data source should have its own SchemaBuilder instance.
type SchemaBuilder struct {
	resourceName string
}

// NewSchemaBuilder creates a new SchemaBuilder for a specific resource or data source
func NewSchemaBuilder(resourceName string) *SchemaBuilder {
	return &SchemaBuilder{resourceName: resourceName}
}

// GetResourceName returns the resource name
func (b *SchemaBuilder) GetResourceName() string {
	return b.resourceName
}

// WithOpenAPIDescription sets the MarkdownDescription for the provided attribute
// by looking up the field description from the OpenAPI specification.
// It accepts an attribute and the Terraform field name in snake_case.
// Returns the same attribute with the description set, preserving the type.
func WithOpenAPIDescription[T SchemaAttribute](b *SchemaBuilder, attr T, fieldName string) T {
	description := docs.GetOpenAPIDescription(b.resourceName, fieldName)

	switch v := any(attr).(type) {
	case *schema.StringAttribute:
		v.MarkdownDescription = description
	case *schema.Int64Attribute:
		v.MarkdownDescription = description
	case *schema.BoolAttribute:
		v.MarkdownDescription = description
	case *schema.SetAttribute:
		v.MarkdownDescription = description
	case *schema.Float64Attribute:
		v.MarkdownDescription = description
	case *schema.NumberAttribute:
		v.MarkdownDescription = description
	case *schema.ListAttribute:
		v.MarkdownDescription = description
	case *schema.SingleNestedAttribute:
		v.MarkdownDescription = description
	case *schema.ObjectAttribute:
		v.MarkdownDescription = description
	}

	return attr
}

// AddAttr adds an attribute with automatic description to the attributes map.
// This eliminates the duplication of the field name.
//
// Description is automatically loaded from the OpenAPI spec:
// 1. Path parameters (organization_id, project_id, etc.) from components.parameters
// 2. Header parameters (if_match) from components.parameters
// 3. Response headers (etag) from components.headers
// 4. Schema references (audit) from components.schemas
// 5. Schema properties (name, description, etc.) from request/response schemas
//
// Example:
//
//	attrs := make(map[string]schema.Attribute)
//	capellaschema.AddAttr(attrs, "name", projectBuilder, stringAttribute([]string{required}))
//	capellaschema.AddAttr(attrs, "organization_id", projectBuilder, stringAttribute([]string{required}))
func AddAttr[T SchemaAttribute](
	attrs map[string]schema.Attribute,
	fieldName string,
	builder *SchemaBuilder,
	attr T,
) {
	// Get description from OpenAPI spec
	description := docs.GetOpenAPIDescription(builder.resourceName, fieldName)

	// Set the description
	switch v := any(attr).(type) {
	case *schema.StringAttribute:
		v.MarkdownDescription = description
	case *schema.Int64Attribute:
		v.MarkdownDescription = description
	case *schema.BoolAttribute:
		v.MarkdownDescription = description
	case *schema.SetAttribute:
		v.MarkdownDescription = description
	case *schema.Float64Attribute:
		v.MarkdownDescription = description
	case *schema.NumberAttribute:
		v.MarkdownDescription = description
	case *schema.ListAttribute:
		v.MarkdownDescription = description
	case *schema.SingleNestedAttribute:
		v.MarkdownDescription = description
	case *schema.ObjectAttribute:
		v.MarkdownDescription = description
	}

	// Convert to schema.Attribute interface
	attrs[fieldName] = any(attr).(schema.Attribute)
}
