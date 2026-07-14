package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSystemIndexesQuery(t *testing.T) {
	tests := []struct {
		name           string
		bucketName     string
		scopeName      string
		collectionName string
		indexName      string
		wantQuery      string
	}{
		{
			name:           "named scope and collection",
			bucketName:     "source",
			scopeName:      "s1",
			collectionName: "c1",
			indexName:      "user_notification_oktaId",
			wantQuery: "SELECT RAW i.index_key FROM system:indexes AS i" +
				" WHERE i.name = \"user_notification_oktaId\" AND i.`using` = \"gsi\"" +
				` AND (i.bucket_id = "source" AND i.scope_id = "s1" AND i.keyspace_id = "c1")`,
		},
		{
			name:           "default scope and collection accepts legacy keyspace form",
			bucketName:     "source",
			scopeName:      "_default",
			collectionName: "_default",
			indexName:      "idx1",
			wantQuery: "SELECT RAW i.index_key FROM system:indexes AS i" +
				" WHERE i.name = \"idx1\" AND i.`using` = \"gsi\"" +
				` AND ((i.bucket_id = "source" AND i.scope_id = "_default" AND i.keyspace_id = "_default")` +
				` OR (i.bucket_id IS MISSING AND i.keyspace_id = "source"))`,
		},
		{
			name:           "names with quotes and backslashes are escaped",
			bucketName:     `bu"cket`,
			scopeName:      "s1",
			collectionName: `col\1`,
			indexName:      `idx"1`,
			wantQuery: "SELECT RAW i.index_key FROM system:indexes AS i" +
				" WHERE i.name = \"idx\\\"1\" AND i.`using` = \"gsi\"" +
				` AND (i.bucket_id = "bu\"cket" AND i.scope_id = "s1" AND i.keyspace_id = "col\\1")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := buildSystemIndexesQuery(tt.bucketName, tt.scopeName, tt.collectionName, tt.indexName)
			assert.Equal(t, tt.wantQuery, query)
		})
	}
}

func TestExtractIndexKeys(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantKeys []string
		wantErr  bool
	}{
		{
			name:     "single key with INCLUDE MISSING",
			body:     "{\"results\": [[\"`oktaId` INCLUDE MISSING\"]], \"status\": \"success\"}",
			wantKeys: []string{"`oktaId` INCLUDE MISSING"},
		},
		{
			name: "multiple keys with modifiers",
			body: "{\"results\": [[\"`c1` INCLUDE MISSING DESC\", \"`c2`\", \"`v` VECTOR\"]]}",
			wantKeys: []string{
				"`c1` INCLUDE MISSING DESC",
				"`c2`",
				"`v` VECTOR",
			},
		},
		{
			name:    "query service returns errors",
			body:    `{"errors": [{"msg": "syntax error - invalid statement"}]}`,
			wantErr: true,
		},
		{
			name:    "no results",
			body:    `{"results": [], "status": "success"}`,
			wantErr: true,
		},
		{
			name:    "empty index key",
			body:    `{"results": [[]], "status": "success"}`,
			wantErr: true,
		},
		{
			name:    "malformed response",
			body:    `not json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := extractIndexKeys([]byte(tt.body))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantKeys, keys)
		})
	}
}
