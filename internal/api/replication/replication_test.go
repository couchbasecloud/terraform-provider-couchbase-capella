package replication

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConflictLoggingJSON(t *testing.T) {
	bucket, collection := "bG9ncw==", "logs.default"
	disabled := false

	tests := []struct {
		name            string
		conflictLogging *ConflictLogging
		want            string
	}{
		{
			name:            "an empty object is sent as {}, which clears the configuration on update",
			conflictLogging: &ConflictLogging{},
			want:            `{}`,
		},
		{
			name: "disabled false is sent, and no rules are not sent",
			conflictLogging: &ConflictLogging{
				Disabled:   &disabled,
				Bucket:     &bucket,
				Collection: &collection,
			},
			want: `{"disabled": false, "bucket": "bG9ncw==", "collection": "logs.default"}`,
		},
		{
			name: "a nil rule is null, an empty rule is {}, and a full rule has both fields",
			conflictLogging: &ConflictLogging{
				Disabled:   &disabled,
				Bucket:     &bucket,
				Collection: &collection,
				LoggingRules: map[string]*ConflictLoggingRule{
					"agents.partners": nil,
					"sales":           {},
					"inventory":       {Bucket: &bucket, Collection: &collection},
				},
			},
			want: `{"disabled": false, "bucket": "bG9ncw==", "collection": "logs.default", "loggingRules": {
				"agents.partners": null,
				"sales": {},
				"inventory": {"bucket": "bG9ncw==", "collection": "logs.default"}
			}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.conflictLogging)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(body))
		})
	}
}
