package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_normalizeDeletionProtectionImportID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "id= key is replaced with cluster_id=",
			input:    "id=533bc6d6-1628-49f5-9384-2dd1c5d025fe,organization_id=5ae82d34-5da8-4313-af5f-90ea4142a4c6,project_id=b02d46ce-c3b1-4a5e-8f3f-ad49112b1dc6",
			expected: "cluster_id=533bc6d6-1628-49f5-9384-2dd1c5d025fe,organization_id=5ae82d34-5da8-4313-af5f-90ea4142a4c6,project_id=b02d46ce-c3b1-4a5e-8f3f-ad49112b1dc6",
		},
		{
			name:     "cluster_id= is left unchanged",
			input:    "cluster_id=533bc6d6-1628-49f5-9384-2dd1c5d025fe,organization_id=5ae82d34-5da8-4313-af5f-90ea4142a4c6,project_id=b02d46ce-c3b1-4a5e-8f3f-ad49112b1dc6",
			expected: "cluster_id=533bc6d6-1628-49f5-9384-2dd1c5d025fe,organization_id=5ae82d34-5da8-4313-af5f-90ea4142a4c6,project_id=b02d46ce-c3b1-4a5e-8f3f-ad49112b1dc6",
		},
		{
			name:     "id= key in non-leading position is replaced",
			input:    "organization_id=5ae82d34-5da8-4313-af5f-90ea4142a4c6,id=533bc6d6-1628-49f5-9384-2dd1c5d025fe,project_id=b02d46ce-c3b1-4a5e-8f3f-ad49112b1dc6",
			expected: "organization_id=5ae82d34-5da8-4313-af5f-90ea4142a4c6,cluster_id=533bc6d6-1628-49f5-9384-2dd1c5d025fe,project_id=b02d46ce-c3b1-4a5e-8f3f-ad49112b1dc6",
		},
		{
			name:     "organization_id and project_id substrings containing 'id' are untouched",
			input:    "cluster_id=aaa,organization_id=bbb,project_id=ccc",
			expected: "cluster_id=aaa,organization_id=bbb,project_id=ccc",
		},
		{
			name:     "empty string is returned unchanged",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeDeletionProtectionImportID(tt.input))
		})
	}
}
