package api

// IndexDDLRequest is a request to run an index DDL statement.
type IndexDDLRequest struct {
	// Definition The index DDL statement.  This can be a CREATE/DROP/ALTER/BUILD statement.
	// Multiple delimited queries are not allowed.
	Definition string `json:"definition"`
}

// QueryError is the error message returned by query service.
type QueryError struct {
	// Msg The error message
	Msg string `json:"msg"`
}

// IndexDDLResponse has an array of errors returned by query service.
type IndexDDLResponse struct {
	Errors *[]QueryError `json:"errors,omitempty"`
}

// IndexDefinitionResponse represents a single index definition.
type IndexDefinitionResponse struct {
	Definition string `json:"definition"`
	IndexName  string `json:"indexName"`
	NumReplica int    `json:"numReplica"`
}

// ListIndexDefinitionsResponse represents a list of index definitions.
type ListIndexDefinitionsResponse struct {
	Definitions []IndexDefinitionResponse `json:"definitions"`
}

// IndexBuildStatusResponse is the build status for an index.
type IndexBuildStatusResponse struct {
	Status string `json:"status"`
}
