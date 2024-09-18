package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = (*GSI)(nil)
	_ resource.ResourceWithConfigure      = (*GSI)(nil)
	_ resource.ResourceWithImportState    = (*GSI)(nil)
	_ resource.ResourceWithValidateConfig = (*GSI)(nil)
)

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

// Create will create/drop/build/alter a primary or secondary index.
func (g *GSI) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.GsiDefinition
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ddl string

	// generate build index statement
	if !plan.BuildIndexes.IsNull() {
		elements := plan.BuildIndexes.Elements()
		var indexes string
		for i, index := range elements {
			indexes += index.String()

			if i < len(elements)-1 {
				indexes += ","
			}
		}
		indexes = strings.ReplaceAll(indexes, "\"", "")

		ddl = fmt.Sprintf(
			"BUILD INDEX ON `%s`.`%s`.`%s`(%s)",
			plan.BucketName.ValueString(),
			plan.ScopeName.ValueString(),
			plan.CollectionName.ValueString(),
			indexes,
		)

	} else {
		if plan.IsPrimary.ValueBool() {
			// named primary index
			if plan.IndexName.ValueString() != "" {
				ddl = fmt.Sprintf(
					"CREATE PRIMARY INDEX `%s` ON `%s`.`%s`.`%s`",
					plan.IndexName.ValueString(),
					plan.BucketName.ValueString(),
					plan.ScopeName.ValueString(),
					plan.CollectionName.ValueString(),
				)

				if plan.With != nil {
					var deferBuild bool
					var replicas int64 = 0
					if !plan.With.DeferBuild.IsNull() {
						deferBuild = plan.With.DeferBuild.ValueBool()
					}
					if !plan.With.NumReplica.IsNull() {
						replicas = plan.With.NumReplica.ValueInt64()
					}
					with := fmt.Sprintf(
						" WITH { \"defer_build\": %t,  \"num_replica\": %d }",
						deferBuild,
						replicas,
					)

					ddl += with
				}
			} else {
				// unamed primary index
				ddl = fmt.Sprintf(
					"CREATE PRIMARY INDEX ON `%s`.`%s`.`%s` WITH { \"defer_build\": %t,  \"num_replica\": %d }",
					plan.BucketName.ValueString(),
					plan.ScopeName.ValueString(),
					plan.CollectionName.ValueString(),
				)

				if plan.With != nil {
					var deferBuild bool
					var replicas int64 = 0
					if !plan.With.DeferBuild.IsNull() {
						deferBuild = plan.With.DeferBuild.ValueBool()
					}
					if !plan.With.NumReplica.IsNull() {
						replicas = plan.With.NumReplica.ValueInt64()
					}
					with := fmt.Sprintf(
						" WITH { \"defer_build\": %t,  \"num_replica\": %d }",
						deferBuild,
						replicas,
					)

					ddl += with
				}
			}
		} else {
			// secondary index
			elements := plan.IndexKeys.Elements()
			var keys string
			for i, k := range elements {
				keys += k.String()

				if i < len(elements)-1 {
					keys += ","
				}
			}
			keys = strings.ReplaceAll(keys, "\"", "")

			ddl = fmt.Sprintf(
				"CREATE INDEX `%s` ON `%s`.`%s`.`%s`(%s) ",
				plan.IndexName.ValueString(),
				plan.BucketName.ValueString(),
				plan.ScopeName.ValueString(),
				plan.CollectionName.ValueString(),
				keys,
			)

			if plan.PartitionBy.ValueString() != "" {
				ddl += fmt.Sprintf(" PARTITION BY HASH(%s) ", plan.PartitionBy.String())
			}

			if plan.Where.ValueString() != "" {
				ddl += fmt.Sprintf(" WHERE %s ", plan.Where.String())
			}

			if plan.With != nil {
				var with string
				if !plan.PartitionBy.IsNull() && !plan.With.NumPartitions.IsNull() {
					var deferBuild bool
					var replicas, partitions int64 = 0, 1
					if !plan.With.DeferBuild.IsNull() {
						deferBuild = plan.With.DeferBuild.ValueBool()
					}
					if !plan.With.NumReplica.IsNull() {
						replicas = plan.With.NumReplica.ValueInt64()
					}
					if !plan.With.NumPartitions.IsNull() {
						partitions = plan.With.NumReplica.ValueInt64()
					}
					with = fmt.Sprintf(
						" WITH { \"defer_build\": %t,  \"num_replica\": %d, \"num_partition\": %d } ",
						deferBuild,
						replicas,
						partitions,
					)
				} else {
					var deferBuild bool
					var replicas int64 = 0
					if !plan.With.DeferBuild.IsNull() {
						deferBuild = plan.With.DeferBuild.ValueBool()
					}
					if !plan.With.NumReplica.IsNull() {
						replicas = plan.With.NumReplica.ValueInt64()
					}

					with = fmt.Sprintf(
						" WITH { \"defer_build\": %t,  \"num_replica\": %d } ",
						deferBuild,
						replicas,
					)
				}
				ddl += with
			}
		}
	}

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

func (g *GSI) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.GsiDefinition
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.BuildIndexes.IsNull() {
		return
	}

	attrs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error during GSI state validation",
			"Could not validate state for "+state.IndexName.ValueString()+": "+err.Error(),
		)
		return
	}

	var (
		organizationID = attrs[providerschema.OrganizationId]
		projectID      = attrs[providerschema.ProjectId]
		clusterID      = attrs[providerschema.ClusterId]
	)

	index, err := g.getQueryIndex(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading query index",
			"Could not read query index "+state.IndexName.ValueString()+": "+err.Error(),
		)
		return
	}

	// when importing index, set all state values
	if state.OrganizationId.IsNull() {
		state.OrganizationId = types.StringValue(organizationID)
		state.ProjectId = types.StringValue(projectID)
		state.ClusterId = types.StringValue(clusterID)
		state.BucketName = types.StringValue(index.Bucket)
		state.ScopeName = types.StringValue(index.Scope)
		state.CollectionName = types.StringValue(index.Collection)
		state.IsPrimary = types.BoolValue(index.IsPrimary)
		state.IndexName = types.StringValue(index.IndexName)

		var keys []attr.Value
		for _, key := range index.SecExprs {
			keys = append(keys, types.StringValue(key))
		}

		keyList, newDiags := types.SetValue(types.StringType, keys)
		if newDiags.HasError() {
			resp.Diagnostics.AddError(
				"Error converting index keys to set type",
				"Could not convert index keys to set type for index "+state.IndexName.ValueString()+": "+err.Error(),
			)
			return
		}

		state.IndexKeys = keyList
		state.PartitionBy = types.StringValue(index.PartitionBy)
		state.Where = types.StringValue(index.Where)
		if state.With != nil {
			state.With.NumReplica = types.Int64Value(int64(index.NumReplica))
			state.With.NumPartitions = types.Int64Value(int64(index.NumPartition))
		}

	} else {
		// when reading an index, only update number of replicas
		if state.With != nil {
			state.With.NumReplica = types.Int64Value(int64(index.NumReplica))
		}

	}

	resp.State.Set(ctx, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

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

func (g *GSI) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.GsiDefinition
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.BuildIndexes.IsNull() {
		return
	}

	var ddl string
	// DROP PRIMARY INDEX is only used for unamed primary indexes.
	if state.IsPrimary.ValueBool() && state.IndexName.ValueString() == "" {
		ddl = fmt.Sprintf(
			"DROP PRIMARY INDEX ON `%s`.`%s`.`%s`",
			state.BucketName.ValueString(),
			state.ScopeName.ValueString(),
			state.CollectionName.ValueString(),
		)
	} else {
		//	DROP INDEX is for named primary index or secondary indexes.
		ddl = fmt.Sprintf(
			"DROP INDEX `%s` ON `%s`.`%s`.`%s`",
			state.IndexName.ValueString(),
			state.BucketName.ValueString(),
			state.ScopeName.ValueString(),
			state.CollectionName.ValueString(),
		)
	}

	if err := g.executeGsiDdl(ctx, &state, ddl); err != nil {
		resp.Diagnostics.AddError(
			"An error occurred while executing index DDL",
			"Error during index DDL execution: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

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
			!config.PartitionBy.IsNull() ||
			config.With != nil {

			resp.Diagnostics.AddAttributeError(
				path.Root("build_indexes"),
				"Invalid Attribute Configuration",
				"build_indexes is set so other optional attributes must be null",
			)
			return
		}
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
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexes",
		g.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	ddlRequest := api.IndexDDLRequest{Definition: ddl}
	response, err := g.Client.ExecuteWithRetry(
		ctx,
		cfg,
		ddlRequest,
		g.Token,
		nil,
	)
	if err != nil {
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

func (g *GSI) getQueryIndex(ctx context.Context, plan *providerschema.GsiDefinition) (
	*api.IndexDefinitionResponse, error,
) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexes/%s?bucket=%s&scope=%s&collection=%s",
		g.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.IndexName.ValueString(),
		plan.BucketName.ValueString(),
		plan.ScopeName.ValueString(),
		plan.CollectionName.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
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
