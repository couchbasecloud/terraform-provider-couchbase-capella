package provider

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
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

// TestCapellaProvider_DocumentedDataSourcesAreRegistered fails when a data source has a page in
// docs/data-sources but is missing from DataSources, which makes Terraform reject it with
// "does not support data source" although it is published in the Registry.
func TestCapellaProvider_DocumentedDataSourcesAreRegistered(t *testing.T) {
	p := &capellaProvider{name: providerName}
	ctx := context.Background()

	registered := make(map[string]bool)
	for _, newDataSource := range p.DataSources(ctx) {
		var resp datasource.MetadataResponse
		newDataSource().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: providerName}, &resp)
		registered[resp.TypeName] = true
	}

	docs, err := filepath.Glob(filepath.Join("..", "..", "docs", "data-sources", "*.md"))
	require.NoError(t, err)
	require.NotEmpty(t, docs, "no data source docs found")

	for _, doc := range docs {
		typeName := providerName + "_" + strings.TrimSuffix(filepath.Base(doc), ".md")
		assert.True(t, registered[typeName], "%s is documented but not registered in DataSources", typeName)
	}
}

func TestCapellaProvider_Resources(t *testing.T) {
	p := &capellaProvider{name: providerName}
	ctx := context.Background()

	resources := p.Resources(ctx)

	require.NotNil(t, resources)
	assert.NotEmpty(t, resources, "provider should expose resources")
}
