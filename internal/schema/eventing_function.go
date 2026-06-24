package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventing_function"
	eventingapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventingfunction"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// EventingFunction is the Terraform data source model for a single eventing function.
type EventingFunction struct {
	OneEventingFunction
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	Export         types.Bool   `tfsdk:"export"`
}

// EventingFunctions is the Terraform data source model for listing eventing functions in a cluster.
// Status is an optional filter that maps to the status query parameter on the list endpoint.
type EventingFunctions struct {
	OrganizationId    types.String          `tfsdk:"organization_id"`
	ProjectId         types.String          `tfsdk:"project_id"`
	ClusterId         types.String          `tfsdk:"cluster_id"`
	Status            types.Set             `tfsdk:"status"`
	EventingFunctions []OneEventingFunction `tfsdk:"eventing_functions"`
}

// OneEventingFunction is a single eventing function entry returned by the list eventing functions
// data source. It mirrors the per-function fields of EventingFunction without the request-scoped
// identifiers and export flag, which live on the parent EventingFunctions model.
type OneEventingFunction struct {
	Description          types.String `tfsdk:"description"`
	Status               types.String `tfsdk:"status"`
	Code                 types.String `tfsdk:"code"`
	Name                 types.String `tfsdk:"name"`
	EventSource          types.Object `tfsdk:"event_source"`
	EventMetadataStorage types.Object `tfsdk:"event_metadata_storage"`
	Settings             types.Object `tfsdk:"settings"`
	Bindings             types.Object `tfsdk:"bindings"`
}

// EventingFunctionResource is the Terraform schema for the eventing function resource.
type EventingFunctionResource struct {
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
	Settings types.Object `tfsdk:"settings"`

	// Bindings holds the bucket, URL and constant bindings.
	Bindings *EventingFunctionBindingsResource `tfsdk:"bindings"`

	// State is the desired terminal activation state, applied via the activationState endpoint.
	// Enum: deployed, undeployed, paused, resumed. It is a write-only control input: the GET
	// response reports the read-only Status, which is mapped back onto State across refreshes.
	State types.String `tfsdk:"state"`
}

// EventingFunctionKeyspace identifies the bucket, scope and collection of an event source or
// metadata store.
type EventingFunctionKeyspace struct {
	Bucket     types.String `tfsdk:"bucket"`
	Scope      types.String `tfsdk:"scope"`
	Collection types.String `tfsdk:"collection"`
}

func (k EventingFunctionKeyspace) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bucket":     types.StringType,
		"scope":      types.StringType,
		"collection": types.StringType,
	}
}

// EventingFunctionSettings holds the runtime settings of an eventing function.
type EventingFunctionSettings struct {
	SqlConsistency        types.String `tfsdk:"sql_consistency"`
	LanguageCompatibility types.String `tfsdk:"language_compatibility"`
	FeedBoundary          types.String `tfsdk:"feed_boundary"`
	WorkerCount           types.Int64  `tfsdk:"worker_count"`
	ScriptTimeout         types.Int64  `tfsdk:"script_timeout"`
	MaxTimerContextSize   types.Int64  `tfsdk:"max_timer_context_size"`
	AllowSyncDocuments    types.Bool   `tfsdk:"allow_sync_documents"`
	CursorAware           types.Bool   `tfsdk:"cursor_aware"`
}

func (s EventingFunctionSettings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"worker_count":           types.Int64Type,
		"script_timeout":         types.Int64Type,
		"sql_consistency":        types.StringType,
		"language_compatibility": types.StringType,
		"feed_boundary":          types.StringType,
		"max_timer_context_size": types.Int64Type,
		"allow_sync_documents":   types.BoolType,
		"cursor_aware":           types.BoolType,
	}
}

// EventingFunctionBindings groups the bucket, URL and constant bindings of an eventing function.
type EventingFunctionBindings struct {
	Buckets   types.List `tfsdk:"buckets"`
	Urls      types.List `tfsdk:"urls"`
	Constants types.List `tfsdk:"constants"`
}

// EventingFunctionBindingsResource groups the bucket, URL and constant bindings.
type EventingFunctionBindingsResource struct {
	Buckets   []EventingFunctionBucketBinding   `tfsdk:"buckets"`
	Urls      []EventingFunctionUrlBinding      `tfsdk:"urls"`
	Constants []EventingFunctionConstantBinding `tfsdk:"constants"`
}

func (b EventingFunctionBindings) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"buckets":   types.ListType{ElemType: types.ObjectType{AttrTypes: EventingFunctionBucketBinding{}.AttributeTypes()}},
		"urls":      types.ListType{ElemType: types.ObjectType{AttrTypes: EventingFunctionUrlBinding{}.AttributeTypes()}},
		"constants": types.ListType{ElemType: types.ObjectType{AttrTypes: EventingFunctionConstantBinding{}.AttributeTypes()}},
	}
}

// EventingFunctionBucketBinding binds a collection to an alias used in the function code.
type EventingFunctionBucketBinding struct {
	Alias      types.String `tfsdk:"alias"`
	Bucket     types.String `tfsdk:"bucket"`
	Scope      types.String `tfsdk:"scope"`
	Collection types.String `tfsdk:"collection"`
	Permission types.String `tfsdk:"permission"`
}

func (b EventingFunctionBucketBinding) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"alias":      types.StringType,
		"bucket":     types.StringType,
		"scope":      types.StringType,
		"collection": types.StringType,
		"permission": types.StringType,
	}
}

// EventingFunctionUrlBinding binds an external URL to an alias used in the function code.
type EventingFunctionUrlBinding struct {
	Alias                  types.String `tfsdk:"alias"`
	Url                    types.String `tfsdk:"url"`
	AllowCookies           types.Bool   `tfsdk:"allow_cookies"`
	ValidateTLSCertificate types.Bool   `tfsdk:"validate_tls_certificate"`
	Authentication         types.Object `tfsdk:"authentication"`
}

func (u EventingFunctionUrlBinding) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"alias":                    types.StringType,
		"url":                      types.StringType,
		"allow_cookies":            types.BoolType,
		"validate_tls_certificate": types.BoolType,
		"authentication":           types.ObjectType{AttrTypes: EventingFunctionURLBindingAuthentication{}.AttributeTypes()},
	}
}

// EventingFunctionURLBindingAuthentication is the flattened representation of the URL binding
// authentication union. Only the fields relevant to the active Type are populated; password and
// bearer_token are returned redacted by the eventing service.
type EventingFunctionURLBindingAuthentication struct {
	Type        types.String `tfsdk:"type"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	BearerToken types.String `tfsdk:"bearer_token"`
}

func (a EventingFunctionURLBindingAuthentication) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":         types.StringType,
		"username":     types.StringType,
		"password":     types.StringType,
		"bearer_token": types.StringType,
	}
}

// AuthenticationFromObject decodes a URL binding authentication object into its struct form,
// returning nil for a null or unknown object.
func AuthenticationFromObject(obj types.Object) *EventingFunctionURLBindingAuthentication {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	var auth EventingFunctionURLBindingAuthentication
	if diags := obj.As(context.Background(), &auth, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil
	}
	return &auth
}

// EventingFunctionConstantBinding binds a constant value to an alias used in the function code.
type EventingFunctionConstantBinding struct {
	Alias types.String `tfsdk:"alias"`
	Value types.String `tfsdk:"value"`
}

func (c EventingFunctionConstantBinding) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"alias": types.StringType,
		"value": types.StringType,
	}
}

func (e EventingFunctionResource) Validate() (map[Attr]string, error) {
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

// NewEventingFunction converts an eventing function API response into the Terraform data source
// model. The export value supplied by the caller is preserved so it round-trips into state.
func NewEventingFunction(
	ctx context.Context,
	function eventing_function.EventingFunction,
	organizationId, projectId, clusterId, name string,
	export types.Bool,
) (EventingFunction, diag.Diagnostics) {
	one, diags := NewOneEventingFunction(ctx, function)

	model := EventingFunction{
		OneEventingFunction: one,
		OrganizationId:      types.StringValue(organizationId),
		ProjectId:           types.StringValue(projectId),
		ClusterId:           types.StringValue(clusterId),
		Export:              export,
	}

	return model, diags
}

// NewOneEventingFunction converts a single eventing function API response into its Terraform model.
// It is shared by the single and list eventing function data sources to build the per-function
// fields and their nested objects.
func NewOneEventingFunction(ctx context.Context, function eventing_function.EventingFunction) (OneEventingFunction, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := OneEventingFunction{
		Name:        types.StringValue(function.Name),
		Description: types.StringPointerValue(function.Description),
		Status:      types.StringPointerValue(function.Status),
		Code:        types.StringPointerValue(function.Code),
	}

	eventSource, d := newEventingFunctionKeyspaceObject(ctx, function.EventSource)
	diags.Append(d...)
	model.EventSource = eventSource

	eventMetadataStorage, d := newEventingFunctionKeyspaceObject(ctx, function.EventMetadataStorage)
	diags.Append(d...)
	model.EventMetadataStorage = eventMetadataStorage

	settings, d := newEventingFunctionSettingsObject(ctx, function.Settings)
	diags.Append(d...)
	model.Settings = settings

	bindings, d := newEventingFunctionBindingsObject(ctx, function.Bindings)
	diags.Append(d...)
	model.Bindings = bindings

	return model, diags
}

func newEventingFunctionKeyspaceObject(ctx context.Context, keyspace eventing_function.Keyspace) (types.Object, diag.Diagnostics) {
	model := EventingFunctionKeyspace{
		Bucket:     types.StringValue(keyspace.Bucket),
		Scope:      types.StringPointerValue(keyspace.Scope),
		Collection: types.StringPointerValue(keyspace.Collection),
	}
	return types.ObjectValueFrom(ctx, model.AttributeTypes(), model)
}

// NewEventingFunctionResource converts an eventing function API response into the Terraform schema.
// prior carries forward values that the GET response does not return: the State action verb
// and any URL binding authentication secrets (matched by alias).
func NewEventingFunctionResource(
	resp *eventingapi.GetEventingFunctionResponse,
	organizationId, projectId, clusterId string,
	prior *EventingFunctionResource,
) (*EventingFunctionResource, error) {
	// The API returns "" for an unset description; preserve the practitioner's planned value
	// (null or "") so the resulting state matches the plan and Terraform does not report an
	// inconsistent result after apply. prior is nil from waitForStatus, where the result is
	// not used as final state, so fall back to the raw API value there.
	description := types.StringPointerValue(resp.Description)
	if prior != nil && description.ValueString() == "" {
		description = prior.Description
	}

	bindings, err := bindingsToSchema(resp.Bindings)
	if err != nil {
		return nil, err
	}

	settings, err := settingsToSchema(resp.Settings)
	if err != nil {
		return nil, err
	}

	fn := &EventingFunctionResource{
		OrganizationId:       types.StringValue(organizationId),
		ProjectId:            types.StringValue(projectId),
		ClusterId:            types.StringValue(clusterId),
		Name:                 types.StringValue(resp.Name),
		Description:          description,
		Code:                 types.StringValue(resp.Code),
		EventSource:          keyspaceToSchema(resp.EventSource),
		EventMetadataStorage: keyspaceToSchema(resp.EventMetadataStorage),
		Settings:             settings,
		Bindings:             bindings,
		State:                types.StringValue(resp.Status),
	}

	if prior != nil {
		carryForwardURLSecrets(fn.Bindings, prior.Bindings)
	}

	return fn, nil
}

func keyspaceToSchema(k eventingapi.Keyspace) *EventingFunctionKeyspace {
	return &EventingFunctionKeyspace{
		Bucket:     types.StringValue(k.Bucket),
		Scope:      types.StringPointerValue(k.Scope),
		Collection: types.StringPointerValue(k.Collection),
	}
}

func settingsToSchema(s eventingapi.Settings) (types.Object, error) {
	obj, d := types.ObjectValue(EventingFunctionSettings{}.AttributeTypes(), map[string]attr.Value{
		"worker_count":           types.Int64PointerValue(s.WorkerCount),
		"script_timeout":         types.Int64PointerValue(s.ScriptTimeout),
		"sql_consistency":        types.StringPointerValue(s.SqlConsistency),
		"language_compatibility": types.StringPointerValue(s.LanguageCompatibility),
		"feed_boundary":          types.StringPointerValue(s.FeedBoundary),
		"max_timer_context_size": types.Int64PointerValue(s.MaxTimerContextSize),
		"allow_sync_documents":   types.BoolPointerValue(s.AllowSyncDocuments),
		"cursor_aware":           types.BoolPointerValue(s.CursorAware),
	})
	if d.HasError() {
		return types.Object{}, fmt.Errorf("failed to convert eventing function settings: %v", d.Errors())
	}
	return obj, nil
}

// bindingsToSchema returns nil when no bindings are present so the optional attribute stays null.
func bindingsToSchema(b eventingapi.Bindings) (*EventingFunctionBindingsResource, error) {
	if len(b.Buckets) == 0 && len(b.Urls) == 0 && len(b.Constants) == 0 {
		return nil, nil
	}

	bindings := &EventingFunctionBindingsResource{}

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
			Authentication:         types.ObjectNull(EventingFunctionURLBindingAuthentication{}.AttributeTypes()),
		}

		if u.Authentication != nil {
			authTypes := EventingFunctionURLBindingAuthentication{}.AttributeTypes()
			auth, d := types.ObjectValue(authTypes, map[string]attr.Value{
				"type":         types.StringValue(u.Authentication.Type),
				"username":     types.StringPointerValue(u.Authentication.Username),
				"password":     types.StringNull(),
				"bearer_token": types.StringNull(),
			})
			if d.HasError() {
				return nil, fmt.Errorf("failed to convert URL binding authentication: %v", d.Errors())
			}
			urlBinding.Authentication = auth
		}
		bindings.Urls = append(bindings.Urls, urlBinding)
	}

	for _, c := range b.Constants {
		bindings.Constants = append(bindings.Constants, EventingFunctionConstantBinding{
			Alias: types.StringValue(c.Alias),
			Value: types.StringValue(c.Value),
		})
	}

	return bindings, nil
}

// carryForwardURLSecrets copies sensitive URL binding authentication values (password,
// bearer token) from the prior state into the refreshed state when the GET response omits them.
//
// why do this ?
//
// eventing API always returns ***** for secrets, so do not set state
// with those values otherwise it will trigger an update. in other
// words, drift detection is not possible with secrets.
func carryForwardURLSecrets(refreshed, prior *EventingFunctionBindingsResource) {
	if refreshed == nil || prior == nil {
		return
	}

	priorByAlias := make(map[string]*EventingFunctionURLBindingAuthentication, len(prior.Urls))
	for i := range prior.Urls {
		if auth := AuthenticationFromObject(prior.Urls[i].Authentication); auth != nil {
			priorByAlias[prior.Urls[i].Alias.ValueString()] = auth
		}
	}

	for i := range refreshed.Urls {
		auth := AuthenticationFromObject(refreshed.Urls[i].Authentication)
		if auth == nil {
			continue
		}

		priorAuth, ok := priorByAlias[refreshed.Urls[i].Alias.ValueString()]
		if !ok {
			continue
		}

		auth.Password = priorAuth.Password
		auth.BearerToken = priorAuth.BearerToken

		obj, d := types.ObjectValueFrom(context.Background(), EventingFunctionURLBindingAuthentication{}.AttributeTypes(), auth)
		if d.HasError() {
			continue
		}
		refreshed.Urls[i].Authentication = obj
	}
}

func newEventingFunctionSettingsObject(ctx context.Context, settings *eventing_function.Settings) (types.Object, diag.Diagnostics) {
	if settings == nil {
		return types.ObjectNull(EventingFunctionSettings{}.AttributeTypes()), nil
	}
	model := EventingFunctionSettings{
		WorkerCount:           types.Int64PointerValue(settings.WorkerCount),
		ScriptTimeout:         types.Int64PointerValue(settings.ScriptTimeout),
		SqlConsistency:        types.StringPointerValue(settings.SqlConsistency),
		LanguageCompatibility: types.StringPointerValue(settings.LanguageCompatibility),
		FeedBoundary:          types.StringPointerValue(settings.FeedBoundary),
		MaxTimerContextSize:   types.Int64PointerValue(settings.MaxTimerContextSize),
		AllowSyncDocuments:    types.BoolPointerValue(settings.AllowSyncDocuments),
		CursorAware:           types.BoolPointerValue(settings.CursorAware),
	}
	return types.ObjectValueFrom(ctx, model.AttributeTypes(), model)
}

func newEventingFunctionBindingsObject(ctx context.Context, bindings *eventing_function.Bindings) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if bindings == nil {
		return types.ObjectNull(EventingFunctionBindings{}.AttributeTypes()), diags
	}

	bucketModels := make([]EventingFunctionBucketBinding, 0, len(bindings.Buckets))
	for _, b := range bindings.Buckets {
		bucketModels = append(bucketModels, EventingFunctionBucketBinding{
			Alias:      types.StringValue(b.Alias),
			Bucket:     types.StringValue(b.Bucket),
			Scope:      types.StringPointerValue(b.Scope),
			Collection: types.StringPointerValue(b.Collection),
			Permission: types.StringPointerValue(b.Permission),
		})
	}
	buckets, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: EventingFunctionBucketBinding{}.AttributeTypes()}, bucketModels)
	diags.Append(d...)

	urlModels := make([]EventingFunctionUrlBinding, 0, len(bindings.Urls))
	for _, u := range bindings.Urls {
		authentication, ad := newEventingFunctionAuthenticationObject(ctx, u.Authentication)
		diags.Append(ad...)
		urlModels = append(urlModels, EventingFunctionUrlBinding{
			Alias:                  types.StringValue(u.Alias),
			Url:                    types.StringValue(u.Url),
			AllowCookies:           types.BoolPointerValue(u.AllowCookies),
			ValidateTLSCertificate: types.BoolPointerValue(u.ValidateTLSCertificate),
			Authentication:         authentication,
		})
	}
	urls, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: EventingFunctionUrlBinding{}.AttributeTypes()}, urlModels)
	diags.Append(d...)

	constantModels := make([]EventingFunctionConstantBinding, 0, len(bindings.Constants))
	for _, c := range bindings.Constants {
		constantModels = append(constantModels, EventingFunctionConstantBinding{
			Alias: types.StringValue(c.Alias),
			Value: types.StringValue(c.Value),
		})
	}
	constants, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: EventingFunctionConstantBinding{}.AttributeTypes()}, constantModels)
	diags.Append(d...)

	model := EventingFunctionBindings{
		Buckets:   buckets,
		Urls:      urls,
		Constants: constants,
	}
	obj, d := types.ObjectValueFrom(ctx, model.AttributeTypes(), model)
	diags.Append(d...)
	return obj, diags
}

func newEventingFunctionAuthenticationObject(ctx context.Context, authentication *eventing_function.Authentication) (types.Object, diag.Diagnostics) {
	if authentication == nil {
		return types.ObjectNull(EventingFunctionURLBindingAuthentication{}.AttributeTypes()), nil
	}
	model := EventingFunctionURLBindingAuthentication{
		Type:        types.StringValue(authentication.Type),
		Username:    types.StringPointerValue(authentication.Username),
		Password:    types.StringPointerValue(authentication.Password),
		BearerToken: types.StringPointerValue(authentication.BearerToken),
	}
	return types.ObjectValueFrom(ctx, model.AttributeTypes(), model)
}
