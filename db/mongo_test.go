package db

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/sunshineOfficial/golib/goerr"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_NewMongo(t *testing.T) {
	t.Skip("use correct credentials")

	client, err := NewMongo(context.Background(), "")
	require.NoError(t, err)
	require.NotNil(t, client)

	list, err := client.Database("Feature1Log").ListCollectionNames(context.Background(), bson.D{})
	require.NoError(t, err)
	assert.True(t, len(list) > 0)

	log.Println(strings.Join(list, ", "))
}

func TestWrapMongoError(t *testing.T) {
	type testCase struct {
		source, expected error
	}

	warn := goerr.NewWarning("warning", nil)

	testCases := map[string]testCase{
		"nil": {},
		"ErrNoDocuments": {
			source:   mongo.ErrNoDocuments,
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

			err := WrapMongoError(data.source)
			assert.ErrorIs(t, err, data.expected)
		})
	}
}
