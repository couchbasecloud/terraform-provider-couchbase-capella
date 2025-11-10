package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Certificate{}
	_ datasource.DataSourceWithConfigure = &Certificate{}
)

// Certificate is the certificate data source implementation.
type Certificate struct {
	*providerschema.Data
}

// NewCertificate is a helper function to simplify the provider implementation.
func NewCertificate() datasource.DataSource {
	return &Certificate{}
}

// Metadata returns the certificates data source type name.
func (c *Certificate) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

// Schema defines the schema for the allowlist data source.
func (c *Certificate) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CertificateSchema()
}

// Read refreshes the Terraform state with the latest data of projects.
func (c *Certificate) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Certificate
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate state is not empty
	err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Certificate",
			"Could not read certificate in cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/certificates", c.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Certificate",
			"Could not read certificate in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	certResp := api.GetCertificateResponse{}
	err = json.Unmarshal(response.Body, &certResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading certificate",
			"Could not read certificate in cluster, unexpected error: "+err.Error(),
		)
		return
	}

	certState := providerschema.OneCertificate{
		Certificate: types.StringValue(certResp.Certificate),
	}

	state.Data = append(state.Data, certState)

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the project data source.
func (c *Certificate) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	c.Data = data
}
