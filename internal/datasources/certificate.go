package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"
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
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

// Schema defines the schema for the allowlist data source.
func (c *Certificate) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
			},
			"certificate": schema.StringAttribute{
				Computed: true,
			},
		},
	}
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
	err := c.validate(state)
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

	response, err := c.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/certificates", c.HostURL, organizationId, projectId, clusterId),
		http.MethodGet,
		nil,
		c.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Certificate",
				"Could not read certificate in cluster "+state.ClusterId.String()+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Certificate",
			"Could not read certificate in cluster "+state.ClusterId.String()+": "+err.Error(),
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

	state.Certificate = types.StringValue(certResp.Certificate)

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

// validate is used to verify that all the fields in the datasource
// have been populated.
func (c *Certificate) validate(state providerschema.Certificate) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if state.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}
