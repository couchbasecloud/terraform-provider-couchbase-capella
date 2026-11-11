package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

// Fallback descriptions for the replication data source, where the resource text is about the config.
const (
	ReplicationDataMobileSyncModeDescription  = "\n - `active` when App Services can run on both the source and the target, else `off`."
	ReplicationDataConflictLoggingDescription = "\n - The conflict logging configuration. It is null when conflict logging is not set."
	ReplicationDataDisabledDescription        = "\n - True when the configuration is kept, but conflicts are not logged."
	ReplicationDataRuleBucketDescription      = "\n - The ID of the log bucket for this rule."
	ReplicationDataRuleCollectionDescription  = "\n - The log collection for this rule, as `scope.collection`."
)

// ReplicationData is the Terraform state of the replication data source. It has the attributes of the
// replication resource, except reverse_replication_id and filter.expressions.skip_restream, since the
// API does not return them.
type ReplicationData struct {
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`

	SourceBucket      types.String                `tfsdk:"source_bucket"`
	Target            *ReplicationTarget          `tfsdk:"target"`
	Direction         types.String                `tfsdk:"direction"`
	Priority          types.String                `tfsdk:"priority"`
	NetworkUsageLimit types.Int64                 `tfsdk:"network_usage_limit"`
	Mappings          types.Set                   `tfsdk:"mappings"`
	Filter            *ReplicationDataFilter      `tfsdk:"filter"`
	MobileSyncMode    types.String                `tfsdk:"mobile_sync_mode"`
	ConflictLogging   *ReplicationConflictLogging `tfsdk:"conflict_logging"`

	Status        types.String `tfsdk:"status"`
	ChangesLeft   types.Int64  `tfsdk:"changes_left"`
	Error         types.String `tfsdk:"error"`
	SourceDetails types.Object `tfsdk:"source_details"`
	TargetDetails types.Object `tfsdk:"target_details"`
	Audit         types.Object `tfsdk:"audit"`
}

// ReplicationDataFilter is the filter of the replication data source.
type ReplicationDataFilter struct {
	DocumentExcludeOptions *ReplicationDocumentExcludeOptions `tfsdk:"document_exclude_options"`
	Expressions            *ReplicationDataFilterExpressions  `tfsdk:"expressions"`
}

// ReplicationDataFilterExpressions is the filter expression of the replication data source.
type ReplicationDataFilterExpressions struct {
	RegEx types.String `tfsdk:"reg_ex"`
}

// NewReplicationData builds the data source state from a GET response. It uses NewReplication with no
// prior state, so the data source and the resource read the response in the same way.
func NewReplicationData(
	ctx context.Context,
	resp *replicationapi.GetReplicationResponse,
	organizationId, projectId, clusterId string,
) (*ReplicationData, error) {
	replication, err := NewReplication(ctx, resp, organizationId, projectId, clusterId, nil)
	if err != nil {
		return nil, err
	}

	return &ReplicationData{
		Id:                replication.Id,
		OrganizationId:    replication.OrganizationId,
		ProjectId:         replication.ProjectId,
		ClusterId:         replication.ClusterId,
		SourceBucket:      replication.SourceBucket,
		Target:            replication.Target,
		Direction:         replication.Direction,
		Priority:          replication.Priority,
		NetworkUsageLimit: replication.NetworkUsageLimit,
		Mappings:          replication.Mappings,
		Filter:            newReplicationDataFilter(replication.Filter),
		MobileSyncMode:    replication.MobileSyncMode,
		ConflictLogging:   replication.ConflictLogging,
		Status:            replication.Status,
		ChangesLeft:       replication.ChangesLeft,
		Error:             replication.Error,
		SourceDetails:     replication.SourceDetails,
		TargetDetails:     replication.TargetDetails,
		Audit:             replication.Audit,
	}, nil
}

func newReplicationDataFilter(filter *ReplicationFilter) *ReplicationDataFilter {
	if filter == nil {
		return nil
	}

	dataFilter := &ReplicationDataFilter{DocumentExcludeOptions: filter.DocumentExcludeOptions}
	if filter.Expressions != nil {
		dataFilter.Expressions = &ReplicationDataFilterExpressions{RegEx: filter.Expressions.RegEx}
	}
	return dataFilter
}
