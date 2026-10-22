package resources

import (
	"context"
	"encoding/base64"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/utils"
)

var (
	_ resource.Resource                   = (*Replication)(nil)
	_ resource.ResourceWithConfigure      = (*Replication)(nil)
	_ resource.ResourceWithImportState    = (*Replication)(nil)
	_ resource.ResourceWithValidateConfig = (*Replication)(nil)
)

const errorMessageWhileReplicationCreation = "There is an error during replication creation. If the replication was " +
	"created anyway, the next apply fails with a conflict error that shows its ID. Import it with " +
	"`terraform import`, unexpected error: "

const errorMessageReplicationAlreadyExists = "\n\nIf this replication was created by a failed apply, import it " +
	"with `terraform import` and the ID in this message."

const errorMessageAfterReplicationCreation = "Replication creation is successful, but encountered an error while " +
	"checking the current state of the replication. Please run `terraform plan` after 1-2 minutes to know the " +
	"current state. Additionally, run `terraform apply --refresh-only` to update the state from remote, " +
	"unexpected error: "

// replicationDeletePollInterval and replicationDeleteTimeout are overridden in tests to avoid
// waiting on real wall-clock time while exercising waitForReplicationDeletion's polling loop.
var (
	replicationDeletePollInterval = 30 * time.Second
	replicationDeleteTimeout      = 15 * time.Minute
)

// replicationCreateRetryInterval and replicationCreateRetryTimeout are overridden in tests, like
// the delete timing above.
var (
	replicationCreateRetryInterval = 30 * time.Second
	replicationCreateRetryTimeout  = 15 * time.Minute
)

// Replication is the XDCR replication resource implementation.
type Replication struct {
	*providerschema.Data
}

func NewReplication() resource.Resource {
	return &Replication{}
}

// Metadata returns the replication resource type name.
func (r *Replication) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication"
}

// Schema defines the schema for the replication resource.
func (r *Replication) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ReplicationSchema()
}

// Configure adds the provider configured client to the replication resource.
func (r *Replication) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// ValidateConfig rejects combinations that the Capella API rejects, before any API call.
func (r *Replication) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var direction, targetType, sourceBucket, targetBucket types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("direction"), &direction)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("target").AtName("type"), &targetType)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("source_bucket"), &sourceBucket)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("target").AtName("bucket"), &targetBucket)...)
	var mappings types.Set
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("mappings"), &mappings)...)
	var conflictLogging types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("conflict_logging"), &conflictLogging)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Defaults do not apply to config values, so a null value means the default.
	effectiveDirection := valueOrDefault(direction, providerschema.ReplicationDirectionOneWay)
	effectiveTargetType := valueOrDefault(targetType, providerschema.ReplicationTargetTypeCapella)

	if effectiveDirection == providerschema.ReplicationDirectionTwoWay && effectiveTargetType == providerschema.ReplicationTargetTypeExternal {
		resp.Diagnostics.AddAttributeError(
			path.Root("direction"),
			"Invalid Replication Configuration",
			"A replication to an external cluster must be oneWay. Set direction to \"oneWay\", or remove it.",
		)
	}

	if isKnown(sourceBucket) && !isBucketId(sourceBucket.ValueString()) {
		resp.Diagnostics.AddAttributeError(path.Root("source_bucket"), "Invalid Bucket ID", bucketIdErrorMessage("source_bucket"))
	}

	// An external target takes the bucket name, so only a Capella target is checked.
	if effectiveTargetType == providerschema.ReplicationTargetTypeCapella && isKnown(targetBucket) && !isBucketId(targetBucket.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("target").AtName("bucket"),
			"Invalid Bucket ID",
			bucketIdErrorMessage("target.bucket")+" For an external target, set target.type = \"external\" and use the bucket name.",
		)
	}

	validateMappingsConfig(ctx, mappings, &resp.Diagnostics)
	validateConflictLoggingConfig(ctx, conflictLogging, &resp.Diagnostics)
}

// validateMappingsConfig checks that each source scope is mapped only once, and that each source
// collection is mapped only once inside its scope.
func validateMappingsConfig(ctx context.Context, mappings types.Set, diags *diag.Diagnostics) {
	if mappings.IsNull() || mappings.IsUnknown() {
		return
	}

	root := path.Root("mappings")
	scopeCounts := make(map[string]int, len(mappings.Elements()))
	for _, element := range mappings.Elements() {
		mappingObject, ok := element.(types.Object)
		if !ok || mappingObject.IsNull() || mappingObject.IsUnknown() {
			continue
		}

		var mapping providerschema.ReplicationScopeMapping
		if conversionDiags := mappingObject.As(ctx, &mapping, basetypes.ObjectAsOptions{}); conversionDiags.HasError() {
			diags.Append(conversionDiags...)
			return
		}

		mappingPath := root.AtSetValue(element)
		scope := mapping.SourceScope.ValueString()
		if isKnown(mapping.SourceScope) {
			scopeCounts[scope]++
			if scopeCounts[scope] == 2 {
				diags.AddAttributeError(
					mappingPath.AtName("source_scope"),
					"Invalid Replication Mappings",
					fmt.Sprintf("source_scope %q is in more than one mapping. Each source scope can have one mapping. "+
						"To replicate more than one collection of a scope, list them in the collections of one mapping.", scope),
				)
			}
		}

		if mapping.Collections.IsNull() || mapping.Collections.IsUnknown() {
			continue
		}
		collectionCounts := make(map[string]int, len(mapping.Collections.Elements()))
		for _, collectionElement := range mapping.Collections.Elements() {
			collectionObject, ok := collectionElement.(types.Object)
			if !ok || collectionObject.IsNull() || collectionObject.IsUnknown() {
				continue
			}

			var collection providerschema.ReplicationCollectionMapping
			if conversionDiags := collectionObject.As(ctx, &collection, basetypes.ObjectAsOptions{}); conversionDiags.HasError() {
				diags.Append(conversionDiags...)
				return
			}
			if isKnown(collection.SourceCollection) {
				name := collection.SourceCollection.ValueString()
				collectionCounts[name]++
				if collectionCounts[name] == 2 {
					diags.AddAttributeError(
						mappingPath.AtName("collections").AtSetValue(collectionElement).AtName("source_collection"),
						"Invalid Replication Mappings",
						fmt.Sprintf("source_collection %q is in more than one collection mapping of source_scope %q. "+
							"Each source collection can have one mapping.", name, scope),
					)
				}
			}
		}
	}
}

// validateConflictLoggingConfig rejects a conflict logging block that the API rejects. A value that
// is not known yet is not checked.
func validateConflictLoggingConfig(ctx context.Context, value types.Object, diags *diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return
	}

	// Local diagnostics, so that errors from the other checks in ValidateConfig do not stop these checks.
	var conversionDiags diag.Diagnostics
	var conflictLogging providerschema.ReplicationConflictLogging
	conversionDiags.Append(value.As(ctx, &conflictLogging, basetypes.ObjectAsOptions{})...)
	if conversionDiags.HasError() {
		diags.Append(conversionDiags...)
		return
	}

	root := path.Root("conflict_logging")
	if !conflictLogging.Disabled.IsUnknown() && !conflictLogging.Disabled.ValueBool() {
		if isKnown(conflictLogging.Bucket) && conflictLogging.Bucket.ValueString() == "" {
			diags.AddAttributeError(root.AtName("bucket"), "Invalid Conflict Logging Configuration",
				"conflict_logging.bucket must not be empty when disabled is false.")
		}
		if isKnown(conflictLogging.Collection) && conflictLogging.Collection.ValueString() == "" {
			diags.AddAttributeError(root.AtName("collection"), "Invalid Conflict Logging Configuration",
				"conflict_logging.collection must not be empty when disabled is false.")
		}
	}

	if isKnown(conflictLogging.Bucket) && conflictLogging.Bucket.ValueString() != "" && !isBucketId(conflictLogging.Bucket.ValueString()) {
		diags.AddAttributeError(root.AtName("bucket"), "Invalid Bucket ID", bucketIdErrorMessage("conflict_logging.bucket"))
	}

	if conflictLogging.LoggingRules.IsNull() || conflictLogging.LoggingRules.IsUnknown() {
		return
	}
	for namespace, element := range conflictLogging.LoggingRules.Elements() {
		ruleObject, ok := element.(types.Object)
		if !ok || ruleObject.IsNull() || ruleObject.IsUnknown() {
			continue
		}

		var rule providerschema.ReplicationLoggingRule
		conversionDiags.Append(ruleObject.As(ctx, &rule, basetypes.ObjectAsOptions{})...)
		if conversionDiags.HasError() {
			diags.Append(conversionDiags...)
			return
		}

		rulePath := root.AtName("logging_rules").AtMapKey(namespace)
		if rule.Bucket.IsNull() != rule.Collection.IsNull() {
			diags.AddAttributeError(rulePath, "Invalid Conflict Logging Rule",
				fmt.Sprintf("The logging rule for %q must set both bucket and collection, or neither. "+
					"Use {} to log to the default bucket and collection, or null to not log.", namespace))
		}
		if isKnown(rule.Bucket) && !isBucketId(rule.Bucket.ValueString()) {
			diags.AddAttributeError(rulePath.AtName("bucket"), "Invalid Bucket ID",
				bucketIdErrorMessage(fmt.Sprintf("conflict_logging.logging_rules[%q].bucket", namespace)))
		}
	}
}

// Create creates a replication. The create call is async: it returns a job ID, and the
// replication ID is known only when the job completes.
func (r *Replication) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Replication
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()

	createRequest, err := createReplicationRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating replication", "Could not build the create request: "+err.Error())
		return
	}

	response, err := r.createReplication(ctx, organizationId, projectId, clusterId, createRequest)
	if err != nil {
		detail := "Could not create replication, unexpected error: " + api.ParseError(err)
		if isConflictError(err) {
			detail += errorMessageReplicationAlreadyExists
		}
		resp.Diagnostics.AddError("Error creating replication", detail)
		return
	}

	createResponse := replicationapi.CreateReplicationAsyncResponse{}
	if err = json.Unmarshal(response.Body, &createResponse); err != nil {
		resp.Diagnostics.AddError(
			"Error creating replication",
			errorMessageWhileReplicationCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}
	if createResponse.JobId == "" {
		resp.Diagnostics.AddError(
			"Error creating replication",
			errorMessageWhileReplicationCreation+"the create response has no job ID",
		)
		return
	}

	job, err := waitForReplicationJob(ctx, r.Data, organizationId, projectId, clusterId, createResponse.JobId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating replication", errorMessageWhileReplicationCreation+api.ParseError(err))
		return
	}

	// Record the IDs before refreshing. A refresh failure is only a warning, and without this the
	// replication would be orphaned. reverseReplicationId exists only in the job response.
	plan.Id = types.StringValue(job.ReplicationId)
	plan.ReverseReplicationId = types.StringNull()
	if job.ReverseReplicationId != "" {
		plan.ReverseReplicationId = types.StringValue(job.ReverseReplicationId)
	}
	diags = resp.State.Set(ctx, initialReplicationState(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.retrieveReplication(ctx, organizationId, projectId, clusterId, job.ReplicationId, &plan)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading replication after creation",
			errorMessageAfterReplicationCreation+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the replication state.
func (r *Replication) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Replication
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication",
			"Could not read replication with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	organizationId := IDs[providerschema.OrganizationId]
	projectId := IDs[providerschema.ProjectId]
	clusterId := IDs[providerschema.ClusterId]
	replicationId := IDs[providerschema.Id]

	refreshedState, err := r.retrieveReplication(ctx, organizationId, projectId, clusterId, replicationId, &state)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading replication",
			"Could not read replication with ID "+state.Id.String()+": "+errString,
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update changes the replication settings. The change applies to both legs of a twoWay replication.
func (r *Replication) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.Replication
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating replication",
			"Could not update replication with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	organizationId := IDs[providerschema.OrganizationId]
	projectId := IDs[providerschema.ProjectId]
	clusterId := IDs[providerschema.ClusterId]
	replicationId := IDs[providerschema.Id]

	updateRequest, err := updateReplicationRequest(ctx, plan, state)
	if err != nil {
		resp.Diagnostics.AddError("Error updating replication", "Could not build the update request: "+err.Error())
		return
	}

	cfg := api.EndpointCfg{
		Url:           replicationapi.ReplicationURL(r.HostURL, organizationId, projectId, clusterId, replicationId),
		Method:        http.MethodPut,
		SuccessStatus: http.StatusNoContent,
	}
	if _, err = r.ClientV1.ExecuteWithRetry(ctx, cfg, updateRequest, r.Token, nil); err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error updating replication",
			"Could not update replication with ID "+replicationId+": "+errString,
		)
		return
	}

	refreshedState, err := r.retrieveReplication(ctx, organizationId, projectId, clusterId, replicationId, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication after update",
			"Replication was updated but could not be read: "+api.ParseError(err),
		)
		return
	}

	diags := resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the replication. One call deletes both legs of a twoWay replication.
func (r *Replication) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.Replication
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting replication",
			"Could not delete replication with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	organizationId := IDs[providerschema.OrganizationId]
	projectId := IDs[providerschema.ProjectId]
	clusterId := IDs[providerschema.ClusterId]
	replicationId := IDs[providerschema.Id]

	cfg := api.EndpointCfg{
		Url:           replicationapi.ReplicationURL(r.HostURL, organizationId, projectId, clusterId, replicationId),
		Method:        http.MethodDelete,
		SuccessStatus: http.StatusNoContent,
	}
	_, err = r.ClientV1.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	switch {
	case err == nil:
		return
	case isNotFoundError(err):
		tflog.Info(ctx, "replication was already deleted", map[string]interface{}{"replication_id": replicationId})
		return
	case isInternalServerError(err):
		// The server watches the delete job for about 60 seconds, then returns 500 while the delete
		// continues in the background.
		tflog.Info(ctx, "replication delete continues in the background, waiting for it to complete", map[string]interface{}{
			"replication_id": replicationId,
		})
		if err = r.waitForReplicationDeletion(ctx, organizationId, projectId, clusterId, replicationId); err != nil {
			resp.Diagnostics.AddError(
				"Error deleting replication",
				"Could not delete replication with ID "+replicationId+": "+api.ParseError(err),
			)
		}
	default:
		resp.Diagnostics.AddError(
			"Error deleting replication",
			"Could not delete replication with ID "+replicationId+": "+api.ParseError(err),
		)
	}
}

// ImportState imports a remote replication that was not created by Terraform.
// example: id=<replication_id>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>
func (r *Replication) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// createReplication sends the create request. It retries while the source or target cluster is not healthy, up to a timeout.
func (r *Replication) createReplication(
	ctx context.Context,
	organizationId, projectId, clusterId string,
	createRequest replicationapi.CreateReplicationRequest,
) (*api.Response, error) {
	cfg := api.EndpointCfg{
		Url:           replicationapi.ClusterReplicationsURL(r.HostURL, organizationId, projectId, clusterId),
		Method:        http.MethodPost,
		SuccessStatus: http.StatusAccepted,
	}

	deadline := time.Now().Add(replicationCreateRetryTimeout)
	for {
		response, err := r.ClientV1.ExecuteWithRetry(ctx, cfg, createRequest, r.Token, nil)
		if err == nil || !isClusterNotHealthyError(err) || time.Now().After(deadline) {
			return response, err
		}

		tflog.Info(ctx, "cluster is not healthy, retrying the replication create", map[string]interface{}{
			"cluster_id": clusterId,
			"error":      api.ParseError(err),
		})
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(replicationCreateRetryInterval):
		}
	}
}

// waitForReplicationDeletion polls the replication until the API reports it is gone.
func (r *Replication) waitForReplicationDeletion(ctx context.Context, organizationId, projectId, clusterId, replicationId string) error {
	ctx, cancel := context.WithTimeout(ctx, replicationDeleteTimeout)
	defer cancel()

	ticker := time.NewTicker(replicationDeletePollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: replication %s: %w", errors.ErrReplicationDeleteTimeout, replicationId, ctx.Err())
		case <-ticker.C:
		}

		_, err := getReplication(ctx, r.Data, organizationId, projectId, clusterId, replicationId)
		switch {
		case err == nil:
			tflog.Info(ctx, "waiting for replication to be deleted", map[string]interface{}{"replication_id": replicationId})
		case isNotFoundError(err):
			return nil
		default:
			tflog.Warn(ctx, "could not get replication while waiting for deletion, retrying", map[string]interface{}{
				"replication_id": replicationId,
				"error":          err.Error(),
			})
		}
	}
}

func (r *Replication) retrieveReplication(
	ctx context.Context,
	organizationId, projectId, clusterId, replicationId string,
	prior *providerschema.Replication,
) (*providerschema.Replication, error) {
	replicationResp, err := getReplication(ctx, r.Data, organizationId, projectId, clusterId, replicationId)
	if err != nil {
		return nil, err
	}

	refreshedState, err := providerschema.NewReplication(ctx, replicationResp, organizationId, projectId, clusterId, prior)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrRefreshingState, err)
	}
	return refreshedState, nil
}

func getReplication(
	ctx context.Context,
	data *providerschema.Data,
	organizationId, projectId, clusterId, replicationId string,
) (*replicationapi.GetReplicationResponse, error) {
	cfg := api.EndpointCfg{
		Url:           replicationapi.ReplicationURL(data.HostURL, organizationId, projectId, clusterId, replicationId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := data.ClientV1.ExecuteWithRetry(ctx, cfg, nil, data.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	replicationResp := replicationapi.GetReplicationResponse{}
	if err = json.Unmarshal(response.Body, &replicationResp); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &replicationResp, nil
}

// initialReplicationState sets the unknown computed values of the plan to null, so the state can
// be saved before the first refresh. The framework rejects a state that has unknown values.
func initialReplicationState(plan providerschema.Replication) providerschema.Replication {
	plan.Status = types.StringNull()
	plan.ChangesLeft = types.Int64Null()
	plan.Error = types.StringNull()
	plan.SourceDetails = types.ObjectNull(providerschema.ReplicationSourceDetailsAttributeTypes())
	plan.TargetDetails = types.ObjectNull(providerschema.ReplicationTargetDetailsAttributeTypes())
	plan.Audit = types.ObjectNull(providerschema.ReplicationAudit{}.AttributeTypes())
	return plan
}

func isNotFoundError(err error) bool {
	notFound, _ := api.CheckResourceNotFoundError(err)
	return notFound
}

func isInternalServerError(err error) bool {
	var apiError *api.Error
	return stderrors.As(err, &apiError) && apiError.HttpStatusCode == http.StatusInternalServerError
}

func isConflictError(err error) bool {
	var apiError *api.Error
	return stderrors.As(err, &apiError) && apiError.HttpStatusCode == http.StatusConflict
}

// isClusterNotHealthyError reports whether err is the 422 that the API returns while the source or
// target cluster is not healthy. The response has no specific error code, so the message is checked.
func isClusterNotHealthyError(err error) bool {
	var apiError *api.Error
	return stderrors.As(err, &apiError) &&
		apiError.HttpStatusCode == http.StatusUnprocessableEntity &&
		strings.Contains(apiError.Message, "must be in a healthy state")
}

// isBucketId reports whether s is a bucket ID: the URL-safe, padded base64 encoding of a bucket name.
func isBucketId(s string) bool {
	decoded, err := base64.URLEncoding.DecodeString(s)
	return err == nil && len(decoded) > 0
}

func bucketIdErrorMessage(attribute string) string {
	return attribute + " must be a bucket ID, not a bucket name. The bucket ID is the URL-safe base64 " +
		"encoding of the bucket name, for example the id attribute of a couchbase-capella_bucket resource."
}

func isKnown(value types.String) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// valueOrDefault returns the value, or defaultValue when it is null. It returns "" when the value
// is unknown, so a check on it does not match.
func valueOrDefault(value types.String, defaultValue string) string {
	switch {
	case value.IsUnknown():
		return ""
	case value.IsNull():
		return defaultValue
	default:
		return value.ValueString()
	}
}

// createReplicationRequest builds the create request from the plan. The request always uses async mode.
func createReplicationRequest(ctx context.Context, plan providerschema.Replication) (replicationapi.CreateReplicationRequest, error) {
	mappings, err := replicationMappingsRequest(ctx, plan.Mappings)
	if err != nil {
		return replicationapi.CreateReplicationRequest{}, err
	}

	request := replicationapi.CreateReplicationRequest{
		SourceBucket:      plan.SourceBucket.ValueString(),
		Mode:              replicationapi.ModeAsync,
		Direction:         utils.StringPointerIfKnown(plan.Direction),
		Priority:          utils.StringPointerIfKnown(plan.Priority),
		NetworkUsageLimit: utils.Int64PointerIfKnown(plan.NetworkUsageLimit),
		Mappings:          mappings,
		Filter:            replicationFilterRequest(plan.Filter),
		ConflictLogging:   conflictLoggingRequest(plan.ConflictLogging),
	}
	if plan.MobileSyncMode.ValueString() == providerschema.ReplicationMobileSyncModeActive {
		request.MobileSyncMode = plan.MobileSyncMode.ValueStringPointer()
	}
	if plan.Target != nil {
		request.Target = replicationapi.CreateReplicationTarget{
			Cluster: plan.Target.Cluster.ValueString(),
			Bucket:  plan.Target.Bucket.ValueString(),
			Type:    utils.StringPointerIfKnown(plan.Target.Type),
		}
	}
	return request, nil
}

// updateReplicationRequest builds the update request from the plan and the prior state.
func updateReplicationRequest(ctx context.Context, plan, state providerschema.Replication) (replicationapi.UpdateReplicationRequest, error) {
	request := replicationapi.UpdateReplicationRequest{
		Priority:          utils.StringPointerIfKnown(plan.Priority),
		NetworkUsageLimit: utils.Int64PointerIfKnown(plan.NetworkUsageLimit),
	}

	if !plan.Mappings.Equal(state.Mappings) {
		if plan.Mappings.IsNull() {
			allScopes := true
			request.AllScopes = &allScopes
		} else {
			mappings, err := replicationMappingsRequest(ctx, plan.Mappings)
			if err != nil {
				return replicationapi.UpdateReplicationRequest{}, err
			}
			request.Mappings = mappings
		}
	}

	if !replicationFiltersEqual(plan.Filter, state.Filter) {
		request.Filter = fullReplicationFilterRequest(plan.Filter)
	}

	if !plan.MobileSyncMode.IsUnknown() && !plan.MobileSyncMode.Equal(state.MobileSyncMode) {
		request.MobileSyncMode = utils.StringPointerIfKnown(plan.MobileSyncMode)
	}

	if !conflictLoggingEqual(plan.ConflictLogging, state.ConflictLogging) {
		request.ConflictLogging = conflictLoggingRequest(plan.ConflictLogging)
		if request.ConflictLogging == nil {
			// An empty object clears the configuration.
			request.ConflictLogging = &replicationapi.ConflictLogging{}
		}
	}

	return request, nil
}

// replicationMappingsRequest converts the mappings set to the API shape. A null set returns nil.
func replicationMappingsRequest(ctx context.Context, set types.Set) ([]replicationapi.ScopeMapping, error) {
	if set.IsNull() || set.IsUnknown() {
		return nil, nil
	}

	var scopes []providerschema.ReplicationScopeMapping
	if diags := set.ElementsAs(ctx, &scopes, false); diags.HasError() {
		return nil, fmt.Errorf("%w: %s", errors.ErrConvertingReplicationMappings, diags.Errors())
	}

	mappings := make([]replicationapi.ScopeMapping, 0, len(scopes))
	for _, scope := range scopes {
		mapping := replicationapi.ScopeMapping{
			SourceScope: scope.SourceScope.ValueString(),
			TargetScope: scope.TargetScope.ValueString(),
		}
		if !scope.Collections.IsNull() && !scope.Collections.IsUnknown() {
			var collections []providerschema.ReplicationCollectionMapping
			if diags := scope.Collections.ElementsAs(ctx, &collections, false); diags.HasError() {
				return nil, fmt.Errorf("%w: %s", errors.ErrConvertingReplicationMappings, diags.Errors())
			}
			for _, collection := range collections {
				mapping.Collections = append(mapping.Collections, replicationapi.CollectionMapping{
					SourceCollection: collection.SourceCollection.ValueString(),
					TargetCollection: collection.TargetCollection.ValueString(),
				})
			}
		}
		mappings = append(mappings, mapping)
	}
	return mappings, nil
}

// replicationFilterRequest converts the filter for a create request. Parts that are not set are not sent.
func replicationFilterRequest(filter *providerschema.ReplicationFilter) *replicationapi.Filter {
	if filter == nil || (filter.DocumentExcludeOptions == nil && filter.Expressions == nil) {
		return nil
	}

	request := &replicationapi.Filter{}
	if filter.DocumentExcludeOptions != nil {
		request.DocumentExcludeOptions = documentExcludeOptionsRequest(filter.DocumentExcludeOptions)
	}
	if filter.Expressions != nil {
		request.Expressions = filterExpressionsRequest(filter.Expressions)
	}
	return request
}

// fullReplicationFilterRequest converts the filter for an update request. Every key is sent, so a
// part that is not set (or a nil filter) clears the matching setting on the server.
func fullReplicationFilterRequest(filter *providerschema.ReplicationFilter) *replicationapi.Filter {
	request := &replicationapi.Filter{
		DocumentExcludeOptions: &replicationapi.DocumentExcludeOptions{},
		Expressions:            &replicationapi.FilterExpressions{},
	}
	if filter == nil {
		return request
	}
	if filter.DocumentExcludeOptions != nil {
		request.DocumentExcludeOptions = documentExcludeOptionsRequest(filter.DocumentExcludeOptions)
	}
	if filter.Expressions != nil {
		request.Expressions = filterExpressionsRequest(filter.Expressions)
	}
	return request
}

func documentExcludeOptionsRequest(options *providerschema.ReplicationDocumentExcludeOptions) *replicationapi.DocumentExcludeOptions {
	return &replicationapi.DocumentExcludeOptions{
		Deletion:   options.Deletion.ValueBool(),
		Expiration: options.Expiration.ValueBool(),
		TTL:        options.TTL.ValueBool(),
		Binary:     options.Binary.ValueBool(),
	}
}

func filterExpressionsRequest(expressions *providerschema.ReplicationFilterExpressions) *replicationapi.FilterExpressions {
	return &replicationapi.FilterExpressions{
		RegEx:        expressions.RegEx.ValueString(),
		SkipRestream: expressions.SkipRestream.ValueBool(),
	}
}

func replicationFiltersEqual(a, b *providerschema.ReplicationFilter) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return documentExcludeOptionsEqual(a.DocumentExcludeOptions, b.DocumentExcludeOptions) &&
		filterExpressionsEqual(a.Expressions, b.Expressions)
}

func documentExcludeOptionsEqual(a, b *providerschema.ReplicationDocumentExcludeOptions) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Deletion.Equal(b.Deletion) && a.Expiration.Equal(b.Expiration) &&
		a.TTL.Equal(b.TTL) && a.Binary.Equal(b.Binary)
}

func filterExpressionsEqual(a, b *providerschema.ReplicationFilterExpressions) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.RegEx.Equal(b.RegEx) && a.SkipRestream.Equal(b.SkipRestream)
}

// conflictLoggingRequest builds the conflictLogging field of a create or update request from the
// conflict_logging block. It returns nil when there is no block.
func conflictLoggingRequest(conflictLogging *providerschema.ReplicationConflictLogging) *replicationapi.ConflictLogging {
	if conflictLogging == nil {
		return nil
	}
	return &replicationapi.ConflictLogging{
		Disabled:     utils.BoolPointerIfKnown(conflictLogging.Disabled),
		Bucket:       utils.StringPointerIfKnown(conflictLogging.Bucket),
		Collection:   utils.StringPointerIfKnown(conflictLogging.Collection),
		LoggingRules: loggingRulesRequest(conflictLogging.LoggingRules),
	}
}

// loggingRulesRequest builds the loggingRules field from the logging_rules map. A null rule is sent
// as null (do not log). A rule with no bucket and no collection is sent as {} (log to the default).
func loggingRulesRequest(rules types.Map) map[string]*replicationapi.ConflictLoggingRule {
	if rules.IsNull() || rules.IsUnknown() {
		return nil
	}

	request := make(map[string]*replicationapi.ConflictLoggingRule, len(rules.Elements()))
	for namespace, element := range rules.Elements() {
		rule, ok := element.(types.Object)
		if !ok || rule.IsNull() || rule.IsUnknown() {
			request[namespace] = nil
			continue
		}

		attributes := rule.Attributes()
		bucket, _ := attributes["bucket"].(types.String)
		collection, _ := attributes["collection"].(types.String)
		request[namespace] = &replicationapi.ConflictLoggingRule{
			Bucket:     utils.StringPointerIfKnown(bucket),
			Collection: utils.StringPointerIfKnown(collection),
		}
	}
	return request
}

// conflictLoggingEqual reports whether two conflict_logging blocks are the same.
func conflictLoggingEqual(a, b *providerschema.ReplicationConflictLogging) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Disabled.Equal(b.Disabled) && a.Bucket.Equal(b.Bucket) &&
		a.Collection.Equal(b.Collection) && a.LoggingRules.Equal(b.LoggingRules)
}
