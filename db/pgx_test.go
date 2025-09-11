package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewPgx(t *testing.T) {
	t.Skip("use correct credentials")

	db, err := NewPgx(context.Background(), "")
	require.NoError(t, err)
	require.NotNil(t, db)
	assert.NoError(t, db.Ping())
	assert.NoError(t, db.Close())
}
