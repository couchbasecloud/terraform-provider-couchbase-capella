package provider

import (
	"context"
	"testing"

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

	t.Run("has global_http_client_timeout attribute", func(t *testing.T) {
		a, ok := attrs[capellaGlobalHTTPClientTimeoutField]
		require.True(t, ok)
		_, ok = a.(schema.Int64Attribute)
		assert.True(t, ok, "global_http_client_timeout should be Int64Attribute")
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
