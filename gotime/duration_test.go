package gotime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration_MarshalJSON(t *testing.T) {
	type testCase struct {
		source   string
		expected time.Duration
	}

	testCases := map[string]testCase{
		"1ms": {
			source:   "1ms",
			expected: time.Millisecond,
		},
		"5ms": {
			source:   "5ms",
			expected: 5 * time.Millisecond,
		},
		"1s": {
			source:   "1s",
			expected: time.Second,
		},
		"5s": {
			source:   "5s",
			expected: 5 * time.Second,
		},
		"1m": {
			source:   "1m",
			expected: time.Minute,
		},
		"2m45s": {
			source:   "2m45s",
			expected: 2*time.Minute + 45*time.Second,
		},
		"5m": {
			source:   "5m",
			expected: 5 * time.Minute,
		},
	}

	type wrapper struct {
		Duration Duration
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var item wrapper

			err := json.Unmarshal([]byte(`{ "duration": "`+tc.source+`" }`), &item)
			require.NoError(t, err)
			require.NotNil(t, item)
			require.NotNil(t, item.Duration)

			assert.Equal(t, tc.expected, time.Duration(item.Duration))
		})
	}
}
