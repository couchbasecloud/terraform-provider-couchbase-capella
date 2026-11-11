package datasources

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*Replications)(nil)
	_ datasource.DataSourceWithConfigure = (*Replications)(nil)
)

// Replications is the replications data source implementation.
type Replications struct {
	*providerschema.Data
}

// NewReplications is a helper function to simplify the provider implementation.
func NewReplications() datasource.DataSource {
	return &Replications{}
}

// Metadata returns the replications data source type name.
func (d *Replications) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replications"
}

// Schema defines the schema for the replications data source.
func (d *Replications) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationsSchema()
}

// Configure adds the provider configured client to the replications data source.
func (d *Replications) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Data = data
}

// Read lists the replications of the cluster, or of the project when cluster_id is not set.
func (d *Replications) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Replications
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := replicationapi.ProjectReplicationsURL(d.HostURL, organizationId, projectId)
	scope := "project " + projectId
	if clusterId != "" {
		url = replicationapi.ClusterReplicationsURL(d.HostURL, organizationId, projectId, clusterId)
		scope = "cluster " + clusterId
	}
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	summaries, err := api.GetPaginated[[]replicationapi.ReplicationSummary](ctx, d.ClientV1, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replications",
			"Could not list the replications of "+scope+": "+api.ParseError(err),
		)
		return
	}

	// An empty list, not null, so that a check on data.# works when there are no replications.
	state.Data = make([]providerschema.ReplicationSummaryData, 0, len(summaries))
	for _, summary := range uniqueReplicationSummaries(summaries) {
		entry, err := providerschema.NewReplicationSummaryData(ctx, summary)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading replications",
				"Could not read replication "+summary.Id+": "+err.Error(),
			)
			return
		}
		state.Data = append(state.Data, entry)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// uniqueReplicationSummaries removes duplicate entries and sorts the entries by ID.
func uniqueReplicationSummaries(summaries []replicationapi.ReplicationSummary) []replicationapi.ReplicationSummary {
	seen := make(map[string]bool, len(summaries))
	unique := make([]replicationapi.ReplicationSummary, 0, len(summaries))
	for _, summary := range summaries {
		if seen[summary.Id] {
			continue
		}
		seen[summary.Id] = true
		unique = append(unique, summary)
	}

	slices.SortFunc(unique, func(a, b replicationapi.ReplicationSummary) int {
		return strings.Compare(a.Id, b.Id)
	})
	return unique
}
