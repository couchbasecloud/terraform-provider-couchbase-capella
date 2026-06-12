package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventing_function"
)

// EventingFunction is the Terraform data source model for a single eventing function.
type EventingFunction struct {
	Description          types.String `tfsdk:"description"`
	Status               types.String `tfsdk:"status"`
	Code                 types.String `tfsdk:"code"`
	OrganizationId       types.String `tfsdk:"organization_id"`
	ProjectId            types.String `tfsdk:"project_id"`
	ClusterId            types.String `tfsdk:"cluster_id"`
	Name                 types.String `tfsdk:"name"`
	EventSource          types.Object `tfsdk:"event_source"`
	EventMetadataStorage types.Object `tfsdk:"event_metadata_storage"`
	Settings             types.Object `tfsdk:"settings"`
	Bindings             types.Object `tfsdk:"bindings"`
	Export               types.Bool   `tfsdk:"export"`
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
	Authentication         types.Object `tfsdk:"authentication"`
	Alias                  types.String `tfsdk:"alias"`
	Url                    types.String `tfsdk:"url"`
	AllowCookies           types.Bool   `tfsdk:"allow_cookies"`
	ValidateTLSCertificate types.Bool   `tfsdk:"validate_tls_certificate"`
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

// NewEventingFunction converts an eventing function API response into the Terraform data source
// model. The export value supplied by the caller is preserved so it round-trips into state.
func NewEventingFunction(
	ctx context.Context,
	function eventing_function.EventingFunction,
	organizationId, projectId, clusterId, name string,
	export types.Bool,
) (EventingFunction, diag.Diagnostics) {
	var diags diag.Diagnostics

	// name is echoed from the configured value, not the API response, so the required attribute
	// always matches config in the read result.
	model := EventingFunction{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Name:           types.StringValue(name),
		Export:         export,
		Description:    types.StringPointerValue(function.Description),
		Status:         types.StringPointerValue(function.Status),
		Code:           types.StringPointerValue(function.Code),
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
