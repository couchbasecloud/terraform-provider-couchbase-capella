package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/docs"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SchemaBuilder provides methods for building resource and data source schemas with OpenAPI integration.
// Each resource or data source should have its own SchemaBuilder instance with the resource/data source name.
type SchemaBuilder interface {
	// GetResourceName returns the name of the resource or data source (e.g., "project", "bucket", "cluster")
	GetResourceName() string

	// WithOpenAPIDescription sets the MarkdownDescription by looking up the field
	// description from the OpenAPI specification for this resource/data source.
	WithOpenAPIDescription(attr any, fieldName string) any
}

// BaseSchemaBuilder implements SchemaBuilder interface
type BaseSchemaBuilder struct {
	resourceName string
}

// NewSchemaBuilder creates a new SchemaBuilder for a specific resource or data source
func NewSchemaBuilder(resourceName string) SchemaBuilder {
	return &BaseSchemaBuilder{resourceName: resourceName}
}

// GetResourceName returns the resource name
func (b *BaseSchemaBuilder) GetResourceName() string {
	return b.resourceName
}

// WithOpenAPIDescription sets the MarkdownDescription for the provided attribute
// by looking up the field description from the OpenAPI specification.
// It accepts an attribute and the Terraform field name in snake_case.
func (b *BaseSchemaBuilder) WithOpenAPIDescription(attr any, fieldName string) any {
	description := docs.GetOpenAPIDescription(b.resourceName, fieldName)

	switch v := attr.(type) {
	case *schema.StringAttribute:
		v.MarkdownDescription = description
		return v
	case *schema.Int64Attribute:
		v.MarkdownDescription = description
		return v
	case *schema.BoolAttribute:
		v.MarkdownDescription = description
		return v
	case *schema.SetAttribute:
		v.MarkdownDescription = description
		return v
	case *schema.Float64Attribute:
		v.MarkdownDescription = description
		return v
	case *schema.NumberAttribute:
		v.MarkdownDescription = description
		return v
	case *schema.ListAttribute:
		v.MarkdownDescription = description
		return v
	case *schema.SingleNestedAttribute:
		v.MarkdownDescription = description
		return v
	case *schema.ObjectAttribute:
		v.MarkdownDescription = description
		return v
	default:
		return attr
	}
}
