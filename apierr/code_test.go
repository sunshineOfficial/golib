package apierr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	type testCase struct {
		scope    Scope
		entity   Entity
		reason   Reason
		flags    []string
		expected string
	}

	testCases := map[string]testCase{
		"no-flags": {
			scope:  "test",
			entity: "data",
			reason: "ok",

			expected: "test.data.ok",
		},
		"2-flags": {
			scope:  "test",
			entity: "data",
			reason: "ok",
			flags: []string{
				"tag1", "tag2",
			},

			expected: "test.data.ok.tag1.tag2",
		},
		"no-scope": {
			entity: "data",
			reason: "ok",

			expected: "internal.data.ok",
		},
		"no-entity": {
			scope:  "test",
			reason: "ok",

			expected: "test.runtime.ok",
		},
		"no-reason": {
			scope:  "test",
			entity: "data",

			expected: "test.data.unknown",
		},
	}

	for name, tc := range testCases {
		data := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := Build(data.scope, data.entity, data.reason, data.flags...)
			assert.Equal(t, data.expected, actual)
		})
	}
}
