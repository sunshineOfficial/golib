package goerr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/sunshineOfficial/golib/gorand"

	"github.com/stretchr/testify/assert"
)

func TestIsContextError(t *testing.T) {
	testCases := map[string]struct {
		source   error
		expected bool
	}{
		"DeadlineExceeded": {
			source:   context.DeadlineExceeded,
			expected: true,
		},
		"Canceled": {
			source:   context.Canceled,
			expected: true,
		},
		"ErrNoRows": {
			source:   sql.ErrNoRows,
			expected: false,
		},
		"wrap/DeadlineExceeded": {
			source:   fmt.Errorf("wrapped: %w", context.DeadlineExceeded),
			expected: true,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsContextError(tc.source))
		})
	}
}

func TestIsCriticalError(t *testing.T) {
	testCases := map[string]struct {
		source   error
		expected bool
	}{
		"context/DeadlineExceeded": {
			source:   context.DeadlineExceeded,
			expected: false,
		},
		"context/Canceled": {
			source:   context.Canceled,
			expected: false,
		},
		"sql/ErrNoRows": {
			source:   sql.ErrNoRows,
			expected: true,
		},
		"goerr/ErrNoRows": {
			source:   WrapNotFoundError("item not found: %w", sql.ErrNoRows),
			expected: false,
		},
		"goerr/ErrNotFound": {
			source:   ErrNotFound,
			expected: false,
		},
		"goerr/NotFoundError": {
			source:   NewNotFoundError("not found", nil),
			expected: false,
		},
		"goerr/Warning": {
			source:   NewWarning("Huston, something went wrong", nil),
			expected: false,
		},
		"random": {
			source:   errors.New(gorand.RandomString(gorand.DefaultAlphabet, 3, 6)),
			expected: true,
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsCriticalError(tc.source))
		})
	}
}
