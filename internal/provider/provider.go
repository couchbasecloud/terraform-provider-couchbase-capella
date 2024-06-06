package provider

import (
	"context"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/datasources"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/version"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &capellaProvider{}

const (
	capellaAuthenticationTokenField = "authentication_token"
	capellaPublicAPIHostField       = "host"
	apiRequestTimeout               = 60 * time.Second
	defaultAPIHostURL               = "https://cloudapi.cloud.couchbase.com"
	providerName                    = "couchbase-capella"
)

// capellaProvider is the provider implementation.
type capellaProvider struct {
	name string
}

// New is a helper function to simplify provider server and testing implementation.
func New() func() provider.Provider {
	return func() provider.Provider {
		return &capellaProvider{
			name: providerName,
		}
	}
}

// Metadata returns the provider type name.
func (p *capellaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = p.name
	resp.Version = version.ProviderVersion
}

// Schema defines the provider-level schema for configuration data.
func (p *capellaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			capellaPublicAPIHostField: schema.StringAttribute{
				Optional:    true,
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

	// if the host URL is not provided, connect to the production Capella API host url by default.
	if config.Host.IsNull() {
		config.Host = types.StringValue(defaultAPIHostURL)
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
		datasources.NewOrganization,
		datasources.NewUsers,
		datasources.NewProjects,
		datasources.NewClusters,
		datasources.NewCertificate,
		datasources.NewAllowLists,
		datasources.NewBuckets,
		datasources.NewDatabaseCredentials,
		datasources.NewApiKeys,
		datasources.NewAppServices,
		datasources.NewBackups,
		datasources.NewScopes,
		datasources.NewCollections,
		datasources.NewSampleBuckets,
		datasources.NewClusterOnOffSchedule,
		datasources.NewAuditLogSettings,
		datasources.NewAuditLogEventIDs,
		datasources.NewAuditLogExport,
		datasources.NewPrivateEndpointService,
		datasources.NewPrivateEndpoints,
	}
}

// Resources defines the resources implemented in the provider.
func (p *capellaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewUser,
		resources.NewProject,
		resources.NewApiKey,
		resources.NewCluster,
		resources.NewAllowList,
		resources.NewDatabaseCredential,
		resources.NewBucket,
		resources.NewAppService,
		resources.NewBackup,
		resources.NewBackupSchedule,
		resources.NewScope,
		resources.NewCollection,
		resources.NewSampleBucket,
		resources.NewClusterOnOffSchedule,
		resources.NewClusterOnOffOnDemand,
		resources.NewAppServiceOnOffOnDemand,
		resources.NewAuditLogSettings,
		resources.NewAuditLogExport,
		resources.NewPrivateEndpointService,
		resources.NewPrivateEndpoint,
	}
}
