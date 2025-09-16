package datasources

import (
	"context"
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/google/uuid"
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
		MarkdownDescription: "The data source to retrieve information about a Capella organization.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the organization.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description of the organization.",
			},
			"preferences": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Organization-wide preferences.",
				Attributes: map[string]schema.Attribute{
					"session_duration": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The maximum allowed time (in seconds) a users can spend in the organization.",
					},
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

	orgID := state.OrganizationId.ValueString()
	orgUUID, parseErr := uuid.Parse(orgID)
	if parseErr != nil {
		resp.Diagnostics.AddError("Error Reading Capella Organization", "invalid organization_id: "+parseErr.Error())
		return
	}

	res, err := o.ClientV2.GetOrganizationByIDWithResponse(ctx, apigen.OrganizationId(orgUUID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Organization",
			"Could not read organization in cluster "+state.OrganizationId.String()+": "+err.Error(),
		)
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("Error Reading Capella Organization", "unexpected response status: "+res.Status())
		return
	}

	org := *res.JSON200
	audit := providerschema.CouchbaseAuditData{
		CreatedAt:  types.StringValue(org.Audit.CreatedAt.String()),
		CreatedBy:  types.StringValue(org.Audit.CreatedBy),
		ModifiedAt: types.StringValue(org.Audit.ModifiedAt.String()),
		ModifiedBy: types.StringValue(org.Audit.ModifiedBy),
		Version:    types.Int64Value(int64(org.Audit.Version)),
	}
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error while audit conversion",
			"Could not perform audit conversion",
		)
		return
	}

	preferences := providerschema.Preferences{}
	if org.Preferences.SessionDuration != nil {
		preferences.SessionDuration = types.Int64Value(int64(*org.Preferences.SessionDuration))
	} else {
		preferences.SessionDuration = types.Int64Null()
	}
	preferencesObj, diags := types.ObjectValueFrom(ctx, preferences.AttributeTypes(), preferences)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error while preferences conversion",
			"Could not perform preferences conversion",
		)
		return
	}

	state = providerschema.Organization{
		OrganizationId: types.StringValue(org.Id.String()),
		Name:           types.StringValue(org.Name),
		Description:    types.StringValue(org.Description),
		Audit:          auditObj,
		Preferences:    preferencesObj,
	}

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
