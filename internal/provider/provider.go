package provider

import (
	"context"
	"terraform-provider-capella/internal/datasources"
	"time"

	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/resources"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &capellaProvider{}

const (
	capellaAuthenticationTokenField = "authentication_token"
	capellaPublicAPIHostField       = "host"
	apiRequestTimeout               = 10 * time.Second
)

// capellaProvider is the provider implementation.
type capellaProvider struct {
	name string
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &capellaProvider{
			name:    "capella",
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *capellaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = p.name
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *capellaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			capellaPublicAPIHostField: schema.StringAttribute{
				Required:    true,
				Description: "Capella Public API HTTPS Host URL",
			},
			capellaAuthenticationTokenField: schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Capella API Token that serves as an authentication mechanism.",
			},
		},
	}
}

// Configure configures the Capella client.
func (p *capellaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring the Capella Client")

	// Retrieve provider data from configuration
	var config providerschema.Config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(capellaPublicAPIHostField),
			"Unknown Capella API Host",
			"The provider cannot create the capella API client as there is an unknown configuration value for the capella API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CAPELLA_HOST environment variable.",
		)
	}

	if config.AuthenticationToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(capellaAuthenticationTokenField),
			"Unknown Capella Authentication Token",
			"The provider cannot create the Capella API client as there is an unknown configuration value for the capella authentication token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CAPELLA_AUTHENTICATION_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Set the host and authentication token to be used

	host := config.Host.ValueString()
	authenticationToken := config.AuthenticationToken.ValueString()

	// If any of the expected configurations are missing, return
	// error with provider-specific guidance.
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(capellaPublicAPIHostField),
			"Missing Capella Public API Host",
			"The provider cannot create the Capella API client as there is a missing or empty value for the Capella API host. "+
				"Set the host value in the configuration or use the TF_VAR_host environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if authenticationToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(capellaAuthenticationTokenField),
			"Missing Capella Authentication Token",
			"The provider cannot create the Capella API client as there is a missing or empty value for the capella authentication token. "+
				"Set the password value in the configuration or use the TF_VAR_auth_token environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, capellaPublicAPIHostField, host)
	ctx = tflog.SetField(ctx, capellaAuthenticationTokenField, authenticationToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, capellaAuthenticationTokenField)

	tflog.Debug(ctx, "Creating Capella client")

	// Create a new capella client using the configuration values
	providerData := &providerschema.Data{
		HostURL: host,
		Token:   authenticationToken,
		Client:  api.NewClient(apiRequestTimeout),
	}

	// Make the Capella client available during DataSource and Resource
	// type Configure methods.
	//
	// DataSourceData is provider-defined data, clients, etc. that is passed
	// to [datasource.ConfigureRequest.ProviderData] for each DataSource type
	// that implements the Configure method.
	resp.DataSourceData = providerData
	// ResourceData is provider-defined data, clients, etc. that is passed
	// to [resource.ConfigureRequest.ProviderData] for each Resource type
	// that implements the Configure method.
	resp.ResourceData = providerData

	tflog.Info(ctx, "Configured Capella client", map[string]any{"success": true})

}

// DataSources defines the data sources implemented in the provider.
func (p *capellaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewProject,
		datasources.NewAllowList,
		datasources.NewOrganization,
		datasources.NewCertificate,
	}
}

// Resources defines the resources implemented in the provider.
func (p *capellaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewProject,
		resources.NewCluster,
		resources.NewAllowList,
		resources.NewDatabaseCredential,
		resources.NewApiKey,
		resources.NewUser,
	}
}
