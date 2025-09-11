package goerr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestWrapNotFoundError(t *testing.T) {
	testCases := map[string]struct {
		source          error
		shouldBeWrapped bool
	}{
		"context/DeadlineExceeded": {
			source:          context.DeadlineExceeded,
			shouldBeWrapped: false,
		},
		"sql/ErrNoRows": {
			source:          sql.ErrNoRows,
			shouldBeWrapped: true,
		},
		"mongo/ErrNoDocuments": {
			source:          mongo.ErrNoDocuments,
			shouldBeWrapped: true,
		},
		"goerr/NotFoundError": {
			source:          NewNotFoundError("", nil),
			shouldBeWrapped: false,
		},
		"nil": {
			source: nil,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			newErr := WrapNotFoundError("oops, not found: %w", tc.source)
			if tc.source == nil {
				require.Nil(t, newErr)
			} else {
				require.NotNil(t, newErr)
			}

			_, was := AsError[NotFoundError](tc.source)
			_, now := AsError[NotFoundError](newErr)

			if tc.shouldBeWrapped {
				assert.True(t, now)
				assert.Regexp(t, "oops, not found: .+", newErr.Error())
			} else {
				assert.False(t, !was && now)
				assert.Equal(t, tc.source, newErr)
			}
		})
	}
}

func TestNotFoundError_Error(t *testing.T) {
	testCases := map[string]struct {
		err  NotFoundError
		text string
	}{
		"ErrNotFound": {
			err:  ErrNotFound,
			text: "not found",
		},
		"custom/string": {
			err:  NewNotFoundError("custom", nil),
			text: "custom",
		},
		"custom/error": {
			err:  NewNotFoundError("", errors.New("custom")),
			text: "custom",
		},
		"custom/string+error": {
			err:  NewNotFoundError("customS", errors.New("customE")),
			text: "customS: customE",
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.text, tc.err.Error())
		})
	}
}

func TestIsNotFound(t *testing.T) {
	testCases := map[string]struct {
		err      error
		expected bool
	}{
		"ErrNotFound": {
			err:      ErrNotFound,
			expected: true,
		},
		"custom/string": {
			err:      NewNotFoundError("custom", nil),
			expected: true,
		},
		"custom/new": {
			err:      NewNotFoundError("", errors.New("custom")),
			expected: true,
		},
		"custom/typed": {
			err:      NewNotFoundError("", context.DeadlineExceeded),
			expected: true,
		},
		"wrapped/const": {
			err:      fmt.Errorf("wrapped: %w", ErrNotFound),
			expected: true,
		},
		"wrapped/type": {
			err:      fmt.Errorf("wrapped: %w", NotFoundError{}),
			expected: true,
		},
		"context.DeadlineExceeded/wrapped": {
			err:      fmt.Errorf("wrapped: %w", context.DeadlineExceeded),
			expected: false,
		},
		"context.DeadlineExceeded/raw": {
			err:      context.DeadlineExceeded,
			expected: false,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsNotFound(tc.err))
		})
	}
}
