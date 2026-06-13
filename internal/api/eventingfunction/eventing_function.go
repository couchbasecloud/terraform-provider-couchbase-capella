package eventingfunction

// Keyspace identifies a bucket, scope and collection that an eventing function
// either listens on (eventSource) or stores its metadata in (eventMetadataStorage).
type Keyspace struct {
	// Bucket is the name of the bucket. It is required.
	Bucket string `json:"bucket"`

	// Scope is the scope within the bucket. Defaults to "_default" server-side.
	Scope *string `json:"scope,omitempty"`

	// Collection is the collection within the scope. Defaults to "_default" server-side.
	Collection *string `json:"collection,omitempty"`
}

// Settings holds the runtime settings that control how an eventing function executes.
// Omitted fields are populated with server-side defaults on create.
type Settings struct {
	// WorkerCount is the number of worker processes the function uses on each Eventing node.
	WorkerCount *int64 `json:"workerCount,omitempty"`

	// ScriptTimeout is the maximum duration, in seconds, a single invocation may run.
	ScriptTimeout *int64 `json:"scriptTimeout,omitempty"`

	// SqlConsistency is the consistency level used by SQL++ (N1QL) statements. Enum: none, request.
	SqlConsistency *string `json:"sqlConsistency,omitempty"`

	// LanguageCompatibility is the language version for backward compatibility. Enum: 6.0.0, 6.5.0, 6.6.2, 7.2.0.
	LanguageCompatibility *string `json:"languageCompatibility,omitempty"`

	// FeedBoundary is the position the function starts processing mutations from on first deployment. Enum: everything, from_now.
	FeedBoundary *string `json:"feedBoundary,omitempty"`

	// MaxTimerContextSize is the maximum size, in bytes, of the context object attached to a Timer.
	MaxTimerContextSize *int64 `json:"maxTimerContextSize,omitempty"`

	// AllowSyncDocuments determines whether the function processes documents managed by App Services.
	AllowSyncDocuments *bool `json:"allowSyncDocuments,omitempty"`

	// CursorAware, when true, suppresses potential duplicate mutations originating from App Services.
	CursorAware *bool `json:"cursorAware,omitempty"`
}

// BucketBinding gives the eventing function direct access to documents in a collection.
type BucketBinding struct {
	// Alias is the name used to refer to the bound collection from the eventing function code.
	Alias string `json:"alias"`

	// Bucket is the name of the bucket this alias resolves to.
	Bucket string `json:"bucket"`

	// Scope is the scope within the bucket. A wildcard "*" gives access to all scopes.
	Scope *string `json:"scope,omitempty"`

	// Collection is the collection within the scope. A wildcard "*" gives access to all collections.
	Collection *string `json:"collection,omitempty"`

	// Permission is the access level on the bound collection. Enum: read, readWrite.
	Permission *string `json:"permission,omitempty"`
}

// URLBindingAuthentication is the authentication scheme used when calling a URL binding.
// The shape is determined by the Type discriminator (none, basic, bearer, digest); a single
// flat struct covers every variant via omitempty.
type URLBindingAuthentication struct {
	// Type is the authentication scheme. Enum: none, basic, bearer, digest.
	Type string `json:"type"`

	// Username is used for basic and digest authentication.
	Username *string `json:"username,omitempty"`

	// Password is used for basic and digest authentication.
	Password *string `json:"password,omitempty"`

	// BearerToken is used for bearer authentication.
	BearerToken *string `json:"bearerToken,omitempty"`
}

// UrlBinding lets the eventing function access an external resource via the curl() function.
type UrlBinding struct {
	// Alias is the name used to refer to the bound endpoint from the eventing function code.
	Alias string `json:"alias"`

	// Url is the fully qualified endpoint URL, including the HTTP or HTTPS scheme.
	Url string `json:"url"`

	// AllowCookies indicates whether the eventing function should allow cookies on the session.
	AllowCookies *bool `json:"allowCookies,omitempty"`

	// ValidateTLSCertificate indicates whether the remote host's TLS certificate is validated.
	ValidateTLSCertificate *bool `json:"validateTLSCertificate,omitempty"`

	// Authentication is the authentication scheme used when calling the URL.
	Authentication *URLBindingAuthentication `json:"authentication,omitempty"`
}

// ConstantBinding exposes a fixed value to the eventing function code as a global variable.
type ConstantBinding struct {
	// Alias is the global variable name exposed to the eventing function code.
	Alias string `json:"alias"`

	// Value is the literal value the alias resolves to.
	Value string `json:"value"`
}

// Bindings groups the bucket, URL and constant bindings of an eventing function.
type Bindings struct {
	// Buckets are the bucket bindings.
	Buckets []BucketBinding `json:"buckets,omitempty"`

	// Urls are the URL bindings.
	Urls []UrlBinding `json:"urls,omitempty"`

	// Constants are the constant bindings.
	Constants []ConstantBinding `json:"constants,omitempty"`
}

// CreateEventingFunctionRequest is the payload for creating an eventing function.
type CreateEventingFunctionRequest struct {
	// Name is the name of the eventing function.
	Name string `json:"name"`

	// Description is the eventing function description.
	Description *string `json:"description,omitempty"`

	// Code is the JavaScript code executed in response to document mutations.
	Code *string `json:"code,omitempty"`

	// EventSource is the keyspace on which the function listens for document mutations.
	EventSource Keyspace `json:"eventSource"`

	// EventMetadataStorage is the keyspace used to store function metadata. Must differ from EventSource.
	EventMetadataStorage Keyspace `json:"eventMetadataStorage"`

	// Settings holds the runtime settings.
	Settings *Settings `json:"settings,omitempty"`

	// Bindings holds the bucket, URL and constant bindings.
	Bindings *Bindings `json:"bindings,omitempty"`
}

// UpdateEventingFunctionRequest is the payload for updating an eventing function. The API applies a
// partial update; the provider sends the full desired definition derived from the Terraform plan.
type UpdateEventingFunctionRequest struct {
	Description          *string   `json:"description,omitempty"`
	Code                 *string   `json:"code,omitempty"`
	EventSource          *Keyspace `json:"eventSource,omitempty"`
	EventMetadataStorage *Keyspace `json:"eventMetadataStorage,omitempty"`
	Settings             *Settings `json:"settings,omitempty"`
	Bindings             *Bindings `json:"bindings,omitempty"`
}

// GetEventingFunctionResponse is the definition of a retrieved eventing function.
type GetEventingFunctionResponse struct {
	// Name is the name of the eventing function.
	Name string `json:"name"`

	// Description is the eventing function description.
	Description *string `json:"description"`

	// Status is the read-only current runtime status of the function.
	// Enum: deployed, undeployed, paused, deploying, undeploying, pausing.
	Status string `json:"status"`

	// Code is the JavaScript code of the eventing function.
	Code string `json:"code"`

	// EventSource is the keyspace on which the function listens for document mutations.
	EventSource Keyspace `json:"eventSource"`

	// EventMetadataStorage is the keyspace used to store function metadata.
	EventMetadataStorage Keyspace `json:"eventMetadataStorage"`

	// Settings holds the runtime settings.
	Settings Settings `json:"settings"`

	// Bindings holds the bucket, URL and constant bindings.
	Bindings Bindings `json:"bindings"`
}

// SetFunctionStateRequest is the payload for changing an eventing function's activation state.
type SetFunctionStateRequest struct {
	// State is the action to take on the function. Enum: deploy, undeploy, pause, resume.
	State string `json:"state"`
}
