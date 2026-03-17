package provider

import (
	"context"
	"testing"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("returns non-nil provider", func(t *testing.T) {
		f := New()
		require.NotNil(t, f)
		p := f()
		require.NotNil(t, p)
		assert.Implements(t, (*provider.Provider)(nil), p)
	})

	t.Run("returns provider with correct name", func(t *testing.T) {
		p := New()()
		ctx := context.Background()
		var resp provider.MetadataResponse
		p.Metadata(ctx, provider.MetadataRequest{}, &resp)
		assert.Equal(t, providerName, resp.TypeName)
		assert.NotEmpty(t, resp.Version)
	})
}

func TestCapellaProvider_Metadata(t *testing.T) {
	p := &capellaProvider{name: providerName}
	ctx := context.Background()
	var resp provider.MetadataResponse

	p.Metadata(ctx, provider.MetadataRequest{}, &resp)

	assert.Equal(t, providerName, resp.TypeName)
	assert.NotEmpty(t, resp.Version)
}

func TestCapellaProvider_Schema(t *testing.T) {
	p := &capellaProvider{name: providerName}
	ctx := context.Background()
	var resp provider.SchemaResponse

	p.Schema(ctx, provider.SchemaRequest{}, &resp)

	require.False(t, resp.Diagnostics.HasError(), "schema should not have errors")
	attrs := resp.Schema.Attributes
	require.NotNil(t, attrs)

	t.Run("has host attribute", func(t *testing.T) {
		a, ok := attrs[capellaPublicAPIHostField]
		require.True(t, ok)
		_, ok = a.(schema.StringAttribute)
		assert.True(t, ok, "host should be StringAttribute")
	})

	t.Run("has authentication_token attribute", func(t *testing.T) {
		a, ok := attrs[capellaAuthenticationTokenField]
		require.True(t, ok)
		sa, ok := a.(schema.StringAttribute)
		require.True(t, ok)
		assert.True(t, sa.Sensitive, "authentication_token should be sensitive")
	})

	t.Run("has global_api_request_timeout attribute", func(t *testing.T) {
		a, ok := attrs[capellaGlobalAPIRequestTimeoutField]
		require.True(t, ok)
		_, ok = a.(schema.Int64Attribute)
		assert.True(t, ok, "global_api_request_timeout should be Int64Attribute")
	})
}

func TestCapellaProvider_DataSources(t *testing.T) {
	p := &capellaProvider{name: providerName}
	ctx := context.Background()

	sources := p.DataSources(ctx)

	require.NotNil(t, sources)
	assert.NotEmpty(t, sources, "provider should expose data sources")
}

func TestCapellaProvider_Resources(t *testing.T) {
	p := &capellaProvider{name: providerName}
	ctx := context.Background()

	resources := p.Resources(ctx)

	require.NotNil(t, resources)
	assert.NotEmpty(t, resources, "provider should expose resources")
}

// TestValidateOpenAPIDescriptions iterates through all resources and data sources
// to ensure every field has a description from either OpenAPI or a manual override.
// This is used by 'make build-docs' to provide clean, summarized feedback.
func TestValidateOpenAPIDescriptions(t *testing.T) {
	providerschema.EnableStrictDocValidation()

	p := &capellaProvider{name: providerName}
	ctx := context.Background()

	// Check Resources
	for _, factory := range p.Resources(ctx) {
		r := factory()
		var resp resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &resp)
	}

	// Check DataSources
	for _, factory := range p.DataSources(ctx) {
		ds := factory()
		var resp datasource.SchemaResponse
		ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
	}

	// This will print a nice summary and exit with 1 if there are any errors.
	providerschema.CheckDocValidationErrors()
}
