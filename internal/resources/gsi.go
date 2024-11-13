package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	internalerrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = (*GSI)(nil)
	_ resource.ResourceWithConfigure      = (*GSI)(nil)
	_ resource.ResourceWithImportState    = (*GSI)(nil)
	_ resource.ResourceWithValidateConfig = (*GSI)(nil)
)

// rate limit create index requests to 60 req/min
// higher rates will (suprisingly) cause indexer to choke
// do not remove this
var limiter = rate.NewLimiter(rate.Every(1*time.Second), 1)

// GSI is the GSI resource implementation.
type GSI struct {
	*providerschema.Data
}

func NewGSI() resource.Resource {
	return &GSI{}
}

// Metadata returns the query index resource type name.
func (g *GSI) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query_indexes"
}

// Schema defines the schema for the query index resource.
func (g *GSI) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GsiSchema()
}

// Create will send a request to create a primary or secondary index.
func (g *GSI) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	if err := limiter.Wait(ctx); err != nil {
		// do not block if rate limiter fails
		tflog.Error(ctx, "rate limiter error: "+err.Error())
	}

	var plan providerschema.GsiDefinition
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ddl string

	// create build index statement.
	if !plan.BuildIndexes.IsNull() {
		var indexes []string
		diags := plan.BuildIndexes.ElementsAs(ctx, &indexes, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ddl = fmt.Sprintf(
			"BUILD INDEX ON `%s`.`%s`.`%s`(%s)",
			plan.BucketName.ValueString(),
			plan.ScopeName.ValueString(),
			plan.CollectionName.ValueString(),
			strings.Join(indexes, ","),
		)

		monitor := func(cfg api.EndpointCfg) (response *api.Response, err error) {
			return g.Client.ExecuteWithRetry(
				ctx,
				cfg,
				nil,
				g.Token,
				nil,
			)
		}

		err := api.EnsureIndexesAreCreated(
			indexes,
			monitor,
			api.Options{
				Host:       g.HostURL,
				OrgId:      plan.OrganizationId.ValueString(),
				ProjectId:  plan.ProjectId.ValueString(),
				ClusterId:  plan.ClusterId.ValueString(),
				Bucket:     plan.BucketName.ValueString(),
				Scope:      plan.ScopeName.ValueString(),
				Collection: plan.CollectionName.ValueString(),
			},
		)

		if err != nil {
			resp.Diagnostics.AddError(
				"Watch Indexes Failed",
				fmt.Sprintln("Error: ", err.Error()),
			)

			return
		}

	} else {
		// create primary index statement.
		if plan.IsPrimary.ValueBool() {
			var indexName string
			if !plan.IndexName.IsNull() {
				indexName = plan.IndexName.ValueString()
			} else {
				indexName = "#primary"
			}
			ddl = fmt.Sprintf(
				"CREATE PRIMARY INDEX `%s` ON `%s`.`%s`.`%s`  WITH { \"defer_build\": %t,  \"num_replica\": %d }",
				indexName,
				plan.BucketName.ValueString(),
				plan.ScopeName.ValueString(),
				plan.CollectionName.ValueString(),
				plan.With.DeferBuild.ValueBool(),
				plan.With.NumReplica.ValueInt64(),
			)
		} else {
			// create secondary index statement.
			var index_keys []string
			diags := plan.IndexKeys.ElementsAs(ctx, &index_keys, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			ddl = fmt.Sprintf(
				"CREATE INDEX `%s` ON `%s`.`%s`.`%s`(%s) ",
				plan.IndexName.ValueString(),
				plan.BucketName.ValueString(),
				plan.ScopeName.ValueString(),
				plan.CollectionName.ValueString(),
				strings.Join(index_keys, ","),
			)

			if !plan.PartitionBy.IsNull() {
				var partition_keys []string
				diags := plan.PartitionBy.ElementsAs(ctx, &partition_keys, false)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}
				ddl += fmt.Sprintf(" PARTITION BY HASH(%s) ", strings.Join(partition_keys, ","))
			}

			if !plan.Where.IsNull() {
				ddl += fmt.Sprintf(" WHERE %s ", plan.Where.ValueString())
			}

			if plan.With != nil {
				type with struct {
					Defer_build   bool  `json:"defer_build,omitempty"`
					Num_replica   int64 `json:"num_replica,omitempty"`
					Num_partition int64 `json:"num_partition,omitempty"`
				}

				var w with

				if !plan.With.DeferBuild.IsNull() {
					w.Defer_build = plan.With.DeferBuild.ValueBool()
				}
				if !plan.With.NumReplica.IsNull() {
					w.Num_replica = plan.With.NumReplica.ValueInt64()
				}
				if !plan.With.NumPartition.IsNull() {
					w.Num_partition = plan.With.NumPartition.ValueInt64()
				}

				b, err := json.Marshal(w)
				if err != nil {
					resp.Diagnostics.AddError(
						"Could not marshal with clause",
						"Unable to marshal with clause.  Error: "+err.Error(),
					)
				}

				if string(b) != "{}" {
					ddl = ddl + " WITH " + string(b)
				}
			}
		}
	}

	err := g.executeGsiDdl(ctx, &plan, ddl)
	switch err {
	case nil:
	case internalerrors.ErrIndexBuildInProgress:
		resp.Diagnostics.AddError(
			"Index build is currently in progress",
			fmt.Sprintf(
				`Could not build index %s in %s.%s.%s as there is another index build already in progress.
The index build will automatically be retried in the background.  Please run "terraform apply --refresh-only".

It is recommended to use deferred builds.  Please see documentation for details.`,
				plan.IndexName.ValueString(),
				plan.BucketName.ValueString(),
				plan.ScopeName.ValueString(),
				plan.CollectionName.ValueString(),
			),
		)

		// save to state file so that refresh can work.
		// this is preferrable to import as there can be many indexes in this state.
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return

	default:
		resp.Diagnostics.AddError(
			"Failed to execute index DDL",
			fmt.Sprintf(
				"Could not execute index %s\nError: %s",
				ddl,
				err.Error(),
			),
		)

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read will get index properties like index keys, number of replicas, etc.
func (g *GSI) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.GsiDefinition
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// when building deferred indexes, there is nothing to read so noop.
	if !state.BuildIndexes.IsNull() {
		return
	}

	attrs, err := state.GetAttributeValues()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating import",
			"Error validating import: "+err.Error(),
		)
		return
	}

	var (
		organizationID = attrs[providerschema.OrganizationId]
		projectID      = attrs[providerschema.ProjectId]
		clusterID      = attrs[providerschema.ClusterId]
		bucketName     = attrs[providerschema.BucketName]
		scopeName      = attrs[providerschema.ScopeName]
		collectionName = attrs[providerschema.CollectionName]
		indexName      = attrs[providerschema.IndexName]
	)

	index, err := g.getQueryIndex(
		ctx,
		organizationID,
		projectID,
		clusterID,
		bucketName,
		scopeName,
		collectionName,
		indexName,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading query index",
			"Could not read query index "+state.IndexName.ValueString()+": "+errString,
		)
		return
	}

	if !state.OrganizationId.IsNull() {
		// when reading an index, only update number of replicas.
		state.With.NumReplica = types.Int64Value(int64(index.NumReplica))

	} else {
		// when importing index, set all attributes.
		state.OrganizationId = types.StringValue(organizationID)
		state.ProjectId = types.StringValue(projectID)
		state.ClusterId = types.StringValue(clusterID)
		state.BucketName = types.StringValue(bucketName)
		state.ScopeName = types.StringValue(scopeName)
		state.CollectionName = types.StringValue(collectionName)
		state.IsPrimary = types.BoolValue(index.IsPrimary)
		state.IndexName = types.StringValue(indexName)

		var keys []attr.Value
		for _, key := range index.SecExprs {
			keys = append(keys, types.StringValue(key))
		}

		keyList, newDiags := types.ListValue(types.StringType, keys)
		if newDiags.HasError() {
			resp.Diagnostics.AddError(
				"Error converting index keys to set type",
				"Could not convert index keys to set type for index "+state.IndexName.ValueString()+": "+err.Error(),
			)
			return
		}

		state.IndexKeys = keyList
		// TODO:  set partition by.
		state.Where = types.StringValue(index.Where)
		if state.With != nil {
			state.With.NumReplica = types.Int64Value(int64(index.NumReplica))
			state.With.NumPartition = types.Int64Value(int64(index.NumPartition))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update will send a request to alter an index.
func (g *GSI) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.GsiDefinition
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ddl := fmt.Sprintf(
		"ALTER INDEX `%s` ON `%s`.`%s`.`%s` WITH { \"action\": \"replica_count\", \"num_replica\" : %d }",
		plan.IndexName.ValueString(),
		plan.BucketName.ValueString(),
		plan.ScopeName.ValueString(),
		plan.CollectionName.ValueString(),
		plan.With.NumReplica.ValueInt64(),
	)

	if err := g.executeGsiDdl(ctx, &plan, ddl); err != nil {
		resp.Diagnostics.AddError(
			"An error occurred while executing index DDL",
			"Error during index DDL execution: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete will send a request to drop an index.
func (g *GSI) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.GsiDefinition
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// when a resource is used for deferred indexes, just remove it from state file.
	if !state.BuildIndexes.IsNull() {
		return
	}

	var indexName string
	if state.IsPrimary.ValueBool() && state.IndexName.IsNull() {
		indexName = "#primary"
	} else {
		indexName = state.IndexName.ValueString()
	}

	ddl := fmt.Sprintf(
		"DROP INDEX `%s` ON `%s`.`%s`.`%s`",
		indexName,
		state.BucketName.ValueString(),
		state.ScopeName.ValueString(),
		state.CollectionName.ValueString(),
	)

	if err := g.executeGsiDdl(ctx, &state, ddl); err != nil {
		resp.Diagnostics.AddError(
			"An error occurred while executing index DDL",
			"Error during index DDL execution: "+err.Error(),
		)
		return
	}

}

// Importstate is used to import an index on the data plane cluster.
func (g *GSI) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("index_name"), req, resp)
}

func (g *GSI) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}
	g.Data = data
}

// ValidateConfig is used to validate multiple attributes in a resource.
//
// a.	For primary indexes, index_keys, where and parition_by must be null.  index_name and with are optional.
// b.	For secondary indexes, index_name and index_keys must be valued.  where, partition_by and with are optional.
// c.	If build_indexes is provided, all of the other optional properties (except scope_name and collection_name) must be null.
func (g *GSI) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config providerschema.GsiDefinition
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.OrganizationId.IsUnknown() {
		return
	}

	if !config.BuildIndexes.IsNull() {
		if !config.IsPrimary.IsNull() ||
			!config.IndexName.IsNull() ||
			!config.IndexKeys.IsNull() ||
			!config.Where.IsNull() ||
			!config.PartitionBy.IsNull() {

			resp.Diagnostics.AddAttributeError(
				path.Root("build_indexes"),
				"Invalid Attribute Configuration",
				"build_indexes is set so other optional attributes must be null",
			)
		}

		return
	}

	if !config.IsPrimary.ValueBool() {
		if config.IndexName.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("index_name"),
				"Missing Attribute Configuration",
				"Expected index_name to be configured but is null",
			)
			return
		}

		if config.IndexKeys.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("index_keys"),
				"Missing Attribute Configuration",
				"Expected index_keys to be configured but is null",
			)
			return
		}

		if config.PartitionBy.IsNull() {
			if config.With != nil && !config.With.NumPartition.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("with"),
					"Invalid Attribute Configuration",
					"Cannot set num_partition for a non-partitioned index",
				)
				return
			}
		}

	}

	if config.IsPrimary.ValueBool() {
		if !config.IndexKeys.IsNull() || !config.Where.IsNull() || !config.PartitionBy.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("is_primary"),
				"Invalid Attribute Configuration",
				"A primary index cannot have index keys, where clause or partition by clause",
			)
			return
		}
	}
}

func (g *GSI) executeGsiDdl(ctx context.Context, plan *providerschema.GsiDefinition, ddl string) error {
	uri := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexes",
		g.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: uri, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	ddlRequest := api.IndexDDLRequest{Definition: ddl}

	response, err := g.Client.ExecuteWithRetry(
		ctx,
		cfg,
		ddlRequest,
		g.Token,
		nil,
	)
	if err != nil {
		// Indexer returns an error if there is a build already in progress.
		// Index build is resource intensive operation from indexer and KV perspective
		// as indexer will request data the keyspace from the beginning.
		//
		// Indexer will automatically retry in the background.
		if apiError, ok := err.(*api.Error); ok {
			if apiError.HttpStatusCode == http.StatusInternalServerError &&
				strings.Contains(
					strings.ToLower(apiError.Message), "build already in progress",
				) {

				return internalerrors.ErrIndexBuildInProgress
			}
		}

		return err
	}

	ddlResponse := api.IndexDDLResponse{}
	err = json.Unmarshal(response.Body, &ddlResponse)
	if err != nil {
		return err
	}

	//  There are some cases where an operation fails yet query service returns 200 OK.
	//	For example, when an index is not found or already exists.
	//  In this case, query service returns errors attribute.
	//  See MB-62943 for more details.
	if ddlResponse.Errors != nil {
		var message string
		if len(ddlResponse.Errors) > 0 {
			message = ddlResponse.Errors[0].Msg
		}

		return errors.New(message)
	}

	return nil
}

func (g *GSI) getQueryIndex(
	ctx context.Context, organizationID, projectID, clusterID, bucketName, scopeName, collectionName, indexName string,
) (
	*api.IndexDefinitionResponse, error,
) {
	uri := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexes/%s?bucket=%s&scope=%s&collection=%s",
		g.HostURL,
		organizationID,
		projectID,
		clusterID,
		url.QueryEscape(indexName),
		bucketName,
		scopeName,
		collectionName,
	)

	cfg := api.EndpointCfg{Url: uri, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := g.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		g.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var index api.IndexDefinitionResponse
	if err = json.Unmarshal(response.Body, &index); err != nil {
		return nil, err
	}

	return &index, nil
}
