package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedisOptions(t *testing.T) {
	host := "redis.go"
	pass := "pass-word-sunshine"

	options := NewRedisOptions(host, pass)
	assert.Equal(t, host, options.Addr)
	assert.Equal(t, pass, options.Password)
}

func TestNewRedis(t *testing.T) {
	client, err := NewRedis(context.Background(), NewRedisOptions("192.168.0.26:6379", "BeP1pA11lWUG5ANw7E/J/ykmqYdUHrzXP4ijFFvpOEMnetnw64/HW/b/nXtHUQvhswhRrErpeox1ATIQ"))
	require.NoError(t, err)
	require.NoError(t, client.Close())
}
