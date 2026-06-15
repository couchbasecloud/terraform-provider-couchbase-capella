package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	eventingapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventingfunction"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// EventingFunction is the Terraform schema for the eventing function resource.
type EventingFunction struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster the eventing function belongs to.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Name is the name of the eventing function. It is the resource identifier and cannot be changed.
	Name types.String `tfsdk:"name"`

	// Description is the eventing function description.
	Description types.String `tfsdk:"description"`

	// Code is the JavaScript code executed in response to document mutations.
	Code types.String `tfsdk:"code"`

	// EventSource is the keyspace on which the function listens for document mutations.
	EventSource *EventingFunctionKeyspace `tfsdk:"event_source"`

	// EventMetadataStorage is the keyspace used to store function metadata.
	EventMetadataStorage *EventingFunctionKeyspace `tfsdk:"event_metadata_storage"`

	// Settings holds the runtime settings that control how the function executes.
	Settings *EventingFunctionSettings `tfsdk:"settings"`

	// Bindings holds the bucket, URL and constant bindings.
	Bindings *EventingFunctionBindings `tfsdk:"bindings"`

	// State is the desired terminal activation state, applied via the activationState endpoint.
	// Enum: deployed, undeployed, paused, resumed. It is a write-only control input: the GET
	// response reports the read-only Status, which is mapped back onto State across refreshes.
	State types.String `tfsdk:"state"`
}

// EventingFunctionKeyspace identifies a bucket, scope and collection.
type EventingFunctionKeyspace struct {
	Bucket     types.String `tfsdk:"bucket"`
	Scope      types.String `tfsdk:"scope"`
	Collection types.String `tfsdk:"collection"`
}

// EventingFunctionSettings holds the runtime settings of an eventing function.
type EventingFunctionSettings struct {
	WorkerCount           types.Int64  `tfsdk:"worker_count"`
	ScriptTimeout         types.Int64  `tfsdk:"script_timeout"`
	SqlConsistency        types.String `tfsdk:"sql_consistency"`
	LanguageCompatibility types.String `tfsdk:"language_compatibility"`
	FeedBoundary          types.String `tfsdk:"feed_boundary"`
	MaxTimerContextSize   types.Int64  `tfsdk:"max_timer_context_size"`
	AllowSyncDocuments    types.Bool   `tfsdk:"allow_sync_documents"`
	CursorAware           types.Bool   `tfsdk:"cursor_aware"`
}

// EventingFunctionBindings groups the bucket, URL and constant bindings.
type EventingFunctionBindings struct {
	Buckets   []EventingFunctionBucketBinding   `tfsdk:"buckets"`
	Urls      []EventingFunctionUrlBinding      `tfsdk:"urls"`
	Constants []EventingFunctionConstantBinding `tfsdk:"constants"`
}

// EventingFunctionBucketBinding gives the function direct access to a collection.
type EventingFunctionBucketBinding struct {
	Alias      types.String `tfsdk:"alias"`
	Bucket     types.String `tfsdk:"bucket"`
	Scope      types.String `tfsdk:"scope"`
	Collection types.String `tfsdk:"collection"`
	Permission types.String `tfsdk:"permission"`
}

// EventingFunctionUrlBinding lets the function access an external resource.
type EventingFunctionUrlBinding struct {
	Alias                  types.String             `tfsdk:"alias"`
	Url                    types.String             `tfsdk:"url"`
	AllowCookies           types.Bool               `tfsdk:"allow_cookies"`
	ValidateTLSCertificate types.Bool               `tfsdk:"validate_tls_certificate"`
	Authentication         *EventingFunctionUrlAuth `tfsdk:"authentication"`
}

// EventingFunctionUrlAuth is the authentication scheme used when calling a URL binding.
type EventingFunctionUrlAuth struct {
	Type        types.String `tfsdk:"type"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	BearerToken types.String `tfsdk:"bearer_token"`
}

// EventingFunctionConstantBinding exposes a fixed value as a global variable.
type EventingFunctionConstantBinding struct {
	Alias types.String `tfsdk:"alias"`
	Value types.String `tfsdk:"value"`
}

// Validate splits a composite import ID and confirms all identifying fields are present.
// The terraform import CLI format is:
// `terraform import capella_eventing_function.fn function_name=<name>,cluster_id=<id>,project_id=<id>,organization_id=<id>`.
func (e EventingFunction) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: e.OrganizationId,
		ProjectId:      e.ProjectId,
		ClusterId:      e.ClusterId,
		FunctionName:   e.Name,
	}

	IDs, err := validateSchemaState(state, FunctionName)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// NewEventingFunction converts an eventing function API response into the Terraform schema.
// prior carries forward values that the GET response does not return: the State action verb
// and any URL binding authentication secrets (matched by alias).
func NewEventingFunction(
	resp *eventingapi.GetEventingFunctionResponse,
	organizationId, projectId, clusterId string,
	prior *EventingFunction,
) *EventingFunction {
	fn := &EventingFunction{
		OrganizationId:       types.StringValue(organizationId),
		ProjectId:            types.StringValue(projectId),
		ClusterId:            types.StringValue(clusterId),
		Name:                 types.StringValue(resp.Name),
		Description:          types.StringPointerValue(resp.Description),
		Code:                 types.StringValue(resp.Code),
		EventSource:          keyspaceToSchema(resp.EventSource),
		EventMetadataStorage: keyspaceToSchema(resp.EventMetadataStorage),
		Settings:             settingsToSchema(resp.Settings),
		Bindings:             bindingsToSchema(resp.Bindings),
	}

	fn.State = types.StringValue(resp.Status)

	if prior != nil {
		carryForwardURLSecrets(fn.Bindings, prior.Bindings)
	}

	return fn
}

func keyspaceToSchema(k eventingapi.Keyspace) *EventingFunctionKeyspace {
	return &EventingFunctionKeyspace{
		Bucket:     types.StringValue(k.Bucket),
		Scope:      types.StringPointerValue(k.Scope),
		Collection: types.StringPointerValue(k.Collection),
	}
}

func settingsToSchema(s eventingapi.Settings) *EventingFunctionSettings {
	return &EventingFunctionSettings{
		WorkerCount:           types.Int64PointerValue(s.WorkerCount),
		ScriptTimeout:         types.Int64PointerValue(s.ScriptTimeout),
		SqlConsistency:        types.StringPointerValue(s.SqlConsistency),
		LanguageCompatibility: types.StringPointerValue(s.LanguageCompatibility),
		FeedBoundary:          types.StringPointerValue(s.FeedBoundary),
		MaxTimerContextSize:   types.Int64PointerValue(s.MaxTimerContextSize),
		AllowSyncDocuments:    types.BoolPointerValue(s.AllowSyncDocuments),
		CursorAware:           types.BoolPointerValue(s.CursorAware),
	}
}

// bindingsToSchema returns nil when no bindings are present so the optional attribute stays null.
func bindingsToSchema(b eventingapi.Bindings) *EventingFunctionBindings {
	if len(b.Buckets) == 0 && len(b.Urls) == 0 && len(b.Constants) == 0 {
		return nil
	}

	bindings := &EventingFunctionBindings{}

	for _, bucket := range b.Buckets {
		bindings.Buckets = append(bindings.Buckets, EventingFunctionBucketBinding{
			Alias:      types.StringValue(bucket.Alias),
			Bucket:     types.StringValue(bucket.Bucket),
			Scope:      types.StringPointerValue(bucket.Scope),
			Collection: types.StringPointerValue(bucket.Collection),
			Permission: types.StringPointerValue(bucket.Permission),
		})
	}

	for _, u := range b.Urls {
		urlBinding := EventingFunctionUrlBinding{
			Alias:                  types.StringValue(u.Alias),
			Url:                    types.StringValue(u.Url),
			AllowCookies:           types.BoolPointerValue(u.AllowCookies),
			ValidateTLSCertificate: types.BoolPointerValue(u.ValidateTLSCertificate),
		}
		if u.Authentication != nil {
			urlBinding.Authentication = &EventingFunctionUrlAuth{
				Type:     types.StringValue(u.Authentication.Type),
				Username: types.StringPointerValue(u.Authentication.Username),
			}
		}
		bindings.Urls = append(bindings.Urls, urlBinding)
	}

	for _, c := range b.Constants {
		bindings.Constants = append(bindings.Constants, EventingFunctionConstantBinding{
			Alias: types.StringValue(c.Alias),
			Value: types.StringValue(c.Value),
		})
	}

	return bindings
}

// carryForwardURLSecrets copies sensitive URL binding authentication values (password,
// bearer token) from the prior state into the refreshed state when the GET response omits them.
//
// why do this ?
//
// eventing API always returns ***** for secrets, so do not set state
// with those values otherwise it will trigger an update. in other
// words, drift detection is not possible with secrets.
func carryForwardURLSecrets(refreshed, prior *EventingFunctionBindings) {
	if refreshed == nil || prior == nil {
		return
	}

	priorByAlias := make(map[string]*EventingFunctionUrlAuth, len(prior.Urls))
	for i := range prior.Urls {
		if prior.Urls[i].Authentication != nil {
			priorByAlias[prior.Urls[i].Alias.ValueString()] = prior.Urls[i].Authentication
		}
	}

	for i := range refreshed.Urls {
		auth := refreshed.Urls[i].Authentication
		if auth == nil {
			continue
		}
		priorAuth, ok := priorByAlias[refreshed.Urls[i].Alias.ValueString()]
		if ok {
			auth.Password = priorAuth.Password
			auth.BearerToken = priorAuth.BearerToken
		}
	}
}
