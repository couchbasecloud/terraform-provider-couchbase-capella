package resources

import (
	"testing"
)

func TestParseIndexKeysFromDefinition(t *testing.T) {
	tests := []struct {
		name       string
		definition string
		wantKeys   []string
		wantErr    bool
	}{
		{
			name:       "single key with INCLUDE MISSING",
			definition: "CREATE INDEX `repro_include_missing` ON `travel-sample`.`_default`.`_default`(`name` INCLUDE MISSING)",
			wantKeys:   []string{"`name` INCLUDE MISSING"},
		},
		{
			name:       "single key without modifier",
			definition: "CREATE INDEX `idx1` ON `bucket`.`scope`.`collection`(`field1`)",
			wantKeys:   []string{"`field1`"},
		},
		{
			name:       "multiple keys with mixed modifiers",
			definition: "CREATE INDEX `idx` ON `bucket`.`scope`.`collection`(`key1` INCLUDE MISSING, `key2`, `key3` INCLUDE MISSING)",
			wantKeys:   []string{"`key1` INCLUDE MISSING", "`key2`", "`key3` INCLUDE MISSING"},
		},
		{
			name:       "key with WHERE clause",
			definition: "CREATE INDEX `user_notification_oktaId` ON `synergy_sv2_dev`(`oktaId` INCLUDE MISSING) WHERE (`type` = \"user_notification_status\")",
			wantKeys:   []string{"`oktaId` INCLUDE MISSING"},
		},
		{
			name:       "key with WITH clause",
			definition: "CREATE INDEX `idx` ON `bucket`.`_default`.`_default`(`name` INCLUDE MISSING) WITH {\"num_replica\":1}",
			wantKeys:   []string{"`name` INCLUDE MISSING"},
		},
		{
			name:       "key with function expression",
			definition: "CREATE INDEX `idx` ON `bucket`.`scope`.`collection`(LOWER(`name`) INCLUDE MISSING, `age`)",
			wantKeys:   []string{"LOWER(`name`) INCLUDE MISSING", "`age`"},
		},
		{
			name:       "key with nested function calls",
			definition: "CREATE INDEX `idx` ON `bucket`.`scope`.`collection`(LOWER(TRIM(`name`)) INCLUDE MISSING)",
			wantKeys:   []string{"LOWER(TRIM(`name`)) INCLUDE MISSING"},
		},
		{
			name:       "key with DESC modifier",
			definition: "CREATE INDEX `idx` ON `bucket`.`scope`.`collection`(`timestamp` DESC, `name` INCLUDE MISSING)",
			wantKeys:   []string{"`timestamp` DESC", "`name` INCLUDE MISSING"},
		},
		{
			name:       "key with PARTITION BY clause after keys",
			definition: "CREATE INDEX `idx` ON `bucket`.`scope`.`collection`(`key1` INCLUDE MISSING, `key2`) PARTITION BY HASH(`key1`)",
			wantKeys:   []string{"`key1` INCLUDE MISSING", "`key2`"},
		},
		{
			name:       "single bucket keyspace",
			definition: "CREATE INDEX `idx` ON `mybucket`(`field` INCLUDE MISSING)",
			wantKeys:   []string{"`field` INCLUDE MISSING"},
		},
		{
			name:       "array index expression",
			definition: "CREATE INDEX `idx` ON `bucket`.`scope`.`collection`(ALL ARRAY `v`.`name` FOR `v` IN `items` END INCLUDE MISSING)",
			wantKeys:   []string{"ALL ARRAY `v`.`name` FOR `v` IN `items` END INCLUDE MISSING"},
		},
		{
			name:       "missing ON clause",
			definition: "DROP INDEX `idx` ON `bucket`",
			wantErr:    true,
		},
		{
			name:       "primary index no keys",
			definition: "CREATE PRIMARY INDEX `primary_idx` ON `bucket`.`scope`.`collection`",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIndexKeysFromDefinition(tt.definition)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIndexKeysFromDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.wantKeys) {
				t.Errorf("parseIndexKeysFromDefinition() got %d keys, want %d keys\ngot:  %v\nwant: %v", len(got), len(tt.wantKeys), got, tt.wantKeys)
				return
			}
			for i := range got {
				if got[i] != tt.wantKeys[i] {
					t.Errorf("parseIndexKeysFromDefinition() key[%d] = %q, want %q", i, got[i], tt.wantKeys[i])
				}
			}
		})
	}
}
