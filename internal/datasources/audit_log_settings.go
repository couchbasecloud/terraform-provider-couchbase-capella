package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AuditLogSettings{}
	_ datasource.DataSourceWithConfigure = &AuditLogSettings{}
)

// AuditLogSettings is the data source implementation.
type AuditLogSettings struct {
	*providerschema.Data
}

// NewAuditLogSettings is a helper function to simplify the provider implementation.
func NewAuditLogSettings() datasource.DataSource {
	return &AuditLogSettings{}
}

// Metadata returns the certificates data source type name.
func (a *AuditLogSettings) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auditlogsettings"
}

func (a *AuditLogSettings) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"auditEnabled":    computedBoolAttribute,
			"enabledEventIDs": computedIntListAttribute,
			"disabledUsers": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{Computed: true, Optional: true},
						"name":   schema.StringAttribute{Computed: true, Optional: true},
					},
				},
			},
		},
	}
}

func (a *AuditLogSettings) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ClusterAuditSettings
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Settings",
			"Could not read Audit Log Settings for cluster ID "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/auditLog", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Settings",
			"Could not read audit log settings in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	auditResponse := api.GetClusterAuditSettingsResponse{}
	err = json.Unmarshal(response.Body, &auditResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading audit settings",
			"Could not read audit settings in cluster, unexpected error: "+err.Error(),
		)
		return
	}

	eventIDs := make([]types.Int64, len(auditResponse.EnabledEventIDs))
	for _, e := range auditResponse.EnabledEventIDs {
		eventIDs = append(eventIDs, types.Int64Value(int64(e)))
	}

	disabledUsers := make([]providerschema.AuditSettingsDisabledUser, len(auditResponse.DisabledUsers))
	for _, u := range auditResponse.DisabledUsers {
		disabledUser := providerschema.AuditSettingsDisabledUser{}
		if u.Domain != nil {
			disabledUser.Domain = types.StringValue(*u.Domain)
		}
		if u.Name != nil {
			disabledUser.Name = types.StringValue(*u.Name)
		}
		disabledUsers = append(disabledUsers, disabledUser)
	}

	state.DisabledUsers = disabledUsers
	state.EnabledEventIDs = eventIDs

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (a *AuditLogSettings) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	a.Data = data
}
