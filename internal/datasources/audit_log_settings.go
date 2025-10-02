package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.TypeName = req.ProviderTypeName + "_audit_log_settings"
}

func (a *AuditLogSettings) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The data source to retrieve audit log configuration settings for an operational cluster. These settings control which events are logged and which users are excluded.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"audit_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies whether audit logging is enabled for this cluster.",
			},
			"enabled_event_ids": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of audit event IDs that are currently enabled for logging. These IDs correspond to specific types of events that will be recorded in the audit log.",
			},
			"disabled_users": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of users whose actions are excluded from audit logging.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The authentication domain of the excluded user.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The username of the excluded user.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest audit log settings.
func (a *AuditLogSettings) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ClusterAuditSettings
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validate(state)
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
	response, err := a.ClientV1.ExecuteWithRetry(
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
			"Error reading audit log settings",
			"Could not read audit log settings in cluster, unexpected error: "+err.Error(),
		)
		return
	}

	eventIDs := make([]types.Int64, len(auditResponse.EnabledEventIDs))
	for i, e := range auditResponse.EnabledEventIDs {
		eventIDs[i] = types.Int64Value(int64(e))
	}

	disabledUsers := make([]providerschema.AuditSettingsDisabledUser, len(auditResponse.DisabledUsers))
	for i, u := range auditResponse.DisabledUsers {
		disabledUser := providerschema.AuditSettingsDisabledUser{}
		if u.Domain != nil {
			disabledUser.Domain = types.StringValue(*u.Domain)
		}
		if u.Name != nil {
			disabledUser.Name = types.StringValue(*u.Name)
		}
		disabledUsers[i] = disabledUser
	}

	state.AuditEnabled = types.BoolValue(auditResponse.AuditEnabled)
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

func (a *AuditLogSettings) validate(state providerschema.ClusterAuditSettings) error {
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
