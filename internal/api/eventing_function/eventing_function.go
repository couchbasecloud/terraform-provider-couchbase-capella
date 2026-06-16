package eventing_function

// EventingFunction is the response returned by the get eventing function endpoint.
// When the request is made with the export query parameter set to true the read-only
type EventingFunction struct {
	Settings             *Settings `json:"settings,omitempty"`
	Bindings             *Bindings `json:"bindings,omitempty"`
	Description          *string   `json:"description,omitempty"`
	Status               *string   `json:"status,omitempty"`
	Code                 *string   `json:"code,omitempty"`
	Name                 string    `json:"name"`
	EventSource          Keyspace  `json:"eventSource"`
	EventMetadataStorage Keyspace  `json:"eventMetadataStorage"`
}

// Keyspace identifies a bucket, scope and collection. Scope and collection default to _default
// on the server when omitted, so they are optional in the response.
type Keyspace struct {
	Scope      *string `json:"scope,omitempty"`
	Collection *string `json:"collection,omitempty"`
	Bucket     string  `json:"bucket"`
}

// Settings are the runtime settings that control how the eventing function is executed.
type Settings struct {
	WorkerCount           *int64  `json:"workerCount,omitempty"`
	ScriptTimeout         *int64  `json:"scriptTimeout,omitempty"`
	SqlConsistency        *string `json:"sqlConsistency,omitempty"`
	LanguageCompatibility *string `json:"languageCompatibility,omitempty"`
	FeedBoundary          *string `json:"feedBoundary,omitempty"`
	MaxTimerContextSize   *int64  `json:"maxTimerContextSize,omitempty"`
	AllowSyncDocuments    *bool   `json:"allowSyncDocuments,omitempty"`
	CursorAware           *bool   `json:"cursorAware,omitempty"`
}

// Bindings expose buckets, external URLs and constant values to the eventing function code.
type Bindings struct {
	Buckets   []BucketBinding   `json:"buckets,omitempty"`
	Urls      []UrlBinding      `json:"urls,omitempty"`
	Constants []ConstantBinding `json:"constants,omitempty"`
}

// BucketBinding gives the eventing function direct access to a collection. Bucket, scope and collection
// may be a wildcard "*".
type BucketBinding struct {
	Scope      *string `json:"scope,omitempty"`
	Collection *string `json:"collection,omitempty"`
	Permission *string `json:"permission,omitempty"`
	Alias      string  `json:"alias"`
	Bucket     string  `json:"bucket"`
}

// UrlBinding lets the eventing function call an external endpoint.
type UrlBinding struct {
	AllowCookies           *bool           `json:"allowCookies,omitempty"`
	ValidateTLSCertificate *bool           `json:"validateTLSCertificate,omitempty"`
	Authentication         *Authentication `json:"authentication,omitempty"`
	Alias                  string          `json:"alias"`
	Url                    string          `json:"url"`
}

// Authentication is the discriminated union describing how a URL binding authenticates. The Type
// field selects the scheme (none, basic, bearer, digest) and which credential fields are populated.
// The eventing service redacts sensitive values, returning five asterisks for password and bearerToken.
type Authentication struct {
	Username    *string `json:"username,omitempty"`
	Password    *string `json:"password,omitempty"`
	BearerToken *string `json:"bearerToken,omitempty"`
	Type        string  `json:"type"`
}

// ConstantBinding exposes a fixed value to the eventing function code as a global variable.
type ConstantBinding struct {
	Alias string `json:"alias"`
	Value string `json:"value"`
}
