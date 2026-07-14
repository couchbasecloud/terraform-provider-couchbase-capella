package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIndexKeysFromDDL(t *testing.T) {
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
			name:       "multiple keys with DESC and ASC",
			definition: "CREATE INDEX `idx1` ON `b`.`s`.`c`(`c1` INCLUDE MISSING DESC, `c2`, `c3` ASC)",
			wantKeys:   []string{"`c1` INCLUDE MISSING DESC", "`c2`", "`c3` ASC"},
		},
		{
			name:       "vector key passes through without modifier knowledge",
			definition: "CREATE INDEX `vec_idx` ON `b`.`s`.`c`(`embedding` VECTOR, `category`)",
			wantKeys:   []string{"`embedding` VECTOR", "`category`"},
		},
		{
			name:       "array expression with nested parens and commas",
			definition: "CREATE INDEX `arr_idx` ON `b`.`s`.`c`((ALL ARRAY FLATTEN_KEYS(`s`.`day` DESC, `s`.`flight`) FOR s IN `schedule` END), `type`)",
			wantKeys: []string{
				"(ALL ARRAY FLATTEN_KEYS(`s`.`day` DESC, `s`.`flight`) FOR s IN `schedule` END)",
				"`type`",
			},
		},
		{
			name:       "backticked identifier containing comma and paren",
			definition: "CREATE INDEX `weird` ON `b`.`s`.`c`(`fie,ld(1`, `f2`)",
			wantKeys:   []string{"`fie,ld(1`", "`f2`"},
		},
		{
			name:       "string literal containing comma and paren",
			definition: `CREATE INDEX ` + "`case_idx`" + ` ON ` + "`b`.`s`.`c`" + `(CASE WHEN ` + "`t`" + ` = "a,(b" THEN 1 ELSE 2 END, ` + "`f2`" + `)`,
			wantKeys:   []string{`CASE WHEN ` + "`t`" + ` = "a,(b" THEN 1 ELSE 2 END`, "`f2`"},
		},
		{
			name:       "trailing partition by, where and with clauses are ignored",
			definition: "CREATE INDEX `idx2` ON `b`.`s`.`c`(`airline`, `destinationairport`) PARTITION BY HASH(`airline`) WHERE (`id` IN [1000, 2000]) WITH {\"num_replica\":1}",
			wantKeys:   []string{"`airline`", "`destinationairport`"},
		},
		{
			name:       "primary index has no key list",
			definition: "CREATE PRIMARY INDEX `#primary` ON `b`.`s`.`c`",
			wantErr:    true,
		},
		{
			name:       "unbalanced parentheses",
			definition: "CREATE INDEX `broken` ON `b`.`s`.`c`(`f1`",
			wantErr:    true,
		},
		{
			name:       "empty key list",
			definition: "CREATE INDEX `empty` ON `b`.`s`.`c`()",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := parseIndexKeysFromDDL(tt.definition)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantKeys, keys)
		})
	}
}
