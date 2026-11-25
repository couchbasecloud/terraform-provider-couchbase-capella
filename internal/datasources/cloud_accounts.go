package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/organization"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*CloudAccounts)(nil)
	_ datasource.DataSourceWithConfigure = (*CloudAccounts)(nil)
)

// CloudAccounts is the data source for Capella's cloud accounts in an organization.
type CloudAccounts struct {
	*providerschema.Data
}

// NewCloudAccounts is a helper function to simplify the provider implementation.
func NewCloudAccounts() datasource.DataSource {
	return &CloudAccounts{}
}

// Metadata returns the data source type name.
func (d *CloudAccounts) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_accounts"
}

// Schema defines the schema for the data source.
func (d *CloudAccounts) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudAccountsSchema()
}

// Configure adds the provider configured client to the data source.
func (d *CloudAccounts) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Data = data
}

// Read refreshes the Terraform state with the organization's cloud accounts.
func (d *CloudAccounts) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.CloudAccountsData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := config.OrganizationId.ValueString()
	cfg := api.EndpointCfg{
		Url:           fmt.Sprintf("%s/v4/organizations/%s/cloudAccounts", d.HostURL, organizationId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cloud Accounts",
			fmt.Sprintf("Could not read the cloud accounts of organization %s: %s", organizationId, api.ParseError(err)),
		)
		return
	}

	var accounts organization.GetCloudAccountsResponse
	if err := json.Unmarshal(response.Body, &accounts); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cloud Accounts",
			fmt.Sprintf("Could not unmarshal the cloud accounts of organization %s: %s", organizationId, err),
		)
		return
	}

	config.AwsCapellaAccount = types.StringValue(accounts.AwsCapellaAccount)
	config.GcpCapellaProject = types.StringValue(accounts.GcpCapellaProject)
	config.AzureCapellaSubscription = types.StringValue(accounts.AzureCapellaSubscription)
	config.AzureCapellaSubscriptionId = types.StringPointerValue(accounts.AzureCapellaSubscriptionId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
