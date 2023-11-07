package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/api/organization"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	providerschema "terraform-provider-capella/internal/schema"
)

var (
	_ datasource.DataSource              = &Organization{}
	_ datasource.DataSourceWithConfigure = &Organization{}
)

type Organization struct {
	*providerschema.Data
}

func NewOrganization() datasource.DataSource {
	return &Organization{}
}

func (o *Organization) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (o *Organization) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"name":            computedStringAttribute,
			"description":     computedStringAttribute,
			"preferences": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"session_duration": computedInt64Attribute,
				},
			},
			"audit": computedAuditAttribute,
		},
	}
}

func (o *Organization) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Organization
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := o.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Organization",
			"Could not read organization in cluster"+state.OrganizationId.String()+": "+err.Error())
		return
	}

	var organizationId = state.OrganizationId.ValueString()

	// Make request to get organization
	response, err := o.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s", o.HostURL, organizationId),
		http.MethodGet,
		nil,
		o.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Organization",
			"Could not read organization in cluster "+state.OrganizationId.String()+": "+api.ParseError(err),
		)
		return
	}

	organizationsResponse := organization.GetOrganizationResponse{}
	err = json.Unmarshal(response.Body, &organizationsResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization",
			"Could not create organization, unexpected error: "+err.Error(),
		)
		return
	}

	audit := providerschema.NewCouchbaseAuditData(organizationsResponse.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error while audit conversion",
			"Could not perform audit conversion",
		)
		return
	}

	var preferences providerschema.Preferences
	if organizationsResponse.Preferences != nil {
		preferences = providerschema.NewPreferences(*organizationsResponse.Preferences)
	}

	preferencesObj, diags := types.ObjectValueFrom(ctx, preferences.AttributeTypes(), preferences)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error while preferences conversion",
			"Could not perform preferences conversion",
		)
		return
	}

	orgState := providerschema.Organization{
		OrganizationId: types.StringValue(organizationsResponse.Id.String()),
		Name:           types.StringValue(organizationsResponse.Name),
		Description:    types.StringValue(*organizationsResponse.Description),
		Audit:          auditObj,
		Preferences:    preferencesObj,
	}
	state = orgState

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (o *Organization) validate(state providerschema.Organization) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	return nil
}

// Configure adds the provider configured client to the organization data source.
func (o *Organization) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	o.Data = data
}
