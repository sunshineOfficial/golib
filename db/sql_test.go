package db

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunshineOfficial/golib/goerr"
)

func TestWrapSqlError(t *testing.T) {
	type testCase struct {
		source, expected error
	}

	warn := goerr.NewWarning("warning", nil)

	testCases := map[string]testCase{
		"nil": {},
		"ErrNoRows": {
			source:   sql.ErrNoRows,
			expected: goerr.ErrNotFound,
		},
		"Warning": {
			source:   warn,
			expected: warn,
		},
	}

	for name, tc := range testCases {
		data := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := WrapSqlError(data.source)
			assert.ErrorIs(t, err, data.expected)
		})
	}
}
