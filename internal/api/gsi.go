package api

// IndexDDLRequest is a request to run an index DDL statement.
type IndexDDLRequest struct {
	// Definition The index DDL statement.  This can be a CREATE/DROP/ALTER/BUILD statement.
	// Multiple delimited queries are not allowed.
	Definition string `json:"definition"`
}

// QueryError is the error message returned by query service.
type QueryError struct {
	// Msg The error message.
	Msg string `json:"msg"`
}

// IndexDDLResponse has an array of errors returned by query service.
type IndexDDLResponse struct {
	Errors []QueryError `json:"errors,omitempty"`
}

// IndexDefinitionResponse represents a single index definition.
type IndexDefinitionResponse struct {
	Bucket       string   `json:"bucket"`
	Scope        string   `json:"scope"`
	Collection   string   `json:"collection"`
	IsPrimary    bool     `json:"is_primary"`
	IndexName    string   `json:"indexName"`
	SecExprs     []string `json:"secExprs"`
	PartitionBy  string   `json:"partition_by"`
	Where        string   `json:"where"`
	NumReplica   int      `json:"numReplica"`
	NumPartition int      `json:"numPartition"`
}

type IndexDefinition struct {
	IndexName  string `json:"indexName"`
	Definition string `json:"definition"`
}

// ListIndexDefinitionsResponse represents a list of index definitions.
type ListIndexDefinitionsResponse struct {
	Definitions []IndexDefinition `json:"definitions"`
}

// IndexBuildStatusResponse is the build status for an index.
type IndexBuildStatusResponse struct {
	Status string `json:"status"`
}
