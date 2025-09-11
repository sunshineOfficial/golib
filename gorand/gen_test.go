package gorand

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			for i := 0; i < 25; i++ {
				num := RandomInt(-1_000_000, 1_000_000)
				assert.Greater(t, num, -1_000_000)
				assert.Less(t, num, 1_000_000)
			}
		})
	})
	t.Run("collision", func(t *testing.T) {
		index := make(map[int]bool)
		for i := 0; i < 1000; i++ {
			num := RandomInt(-1_000_000, 1_000_000)
			require.False(t, index[num], "i: %d, value: %d", i, num)

			index[num] = true
		}
	})
}

func TestRandomString(t *testing.T) {
	t.Run("fixed length", func(t *testing.T) {
		text := RandomString(DefaultAlphabet, 16, 16)
		require.Len(t, text, 16)
	})
	t.Run("random length", func(t *testing.T) {
		text := RandomString(DefaultAlphabet, 16, 32)
		require.Greater(t, len(text), 16)
		require.Less(t, len(text), 32)
	})
	t.Run("nil alphabet", func(t *testing.T) {
		text := RandomString(nil, 16, 16)
		require.Len(t, text, 16)
	})
}
