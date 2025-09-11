package goerr

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestWarning_Error(t *testing.T) {
	testCases := map[string]struct {
		err  Warning
		text string
	}{
		"empty": {
			err:  NewWarning("", nil),
			text: "",
		},
		"string": {
			err:  NewWarning("string", nil),
			text: "string",
		},
		"error": {
			err:  NewWarning("", errors.New("error")),
			text: "error",
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.text, tc.err.Error())
		})
	}
}

func TestWarning_Unwrap(t *testing.T) {
	testCases := []error{
		sql.ErrNoRows, mongo.ErrNoDocuments, context.DeadlineExceeded, errors.New("custom"), nil,
	}

	t.Parallel()
	for _, tc := range testCases {
		if tc == nil {
			assert.Nil(t, NewWarning("", tc).Unwrap())
			continue
		}

		t.Run(tc.Error(), func(t *testing.T) {
			assert.True(t, errors.Is(NewWarning("", tc), tc))
		})
	}
}
