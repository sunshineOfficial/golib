package gohttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRequestFailedError(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		require.Error(t, NewRequestFailedError(http.StatusOK))
	})
	t.Run("one", func(t *testing.T) {
		require.Error(t, NewRequestFailedError(http.StatusOK, http.StatusOK))
	})
	t.Run("two", func(t *testing.T) {
		require.Error(t, NewRequestFailedError(http.StatusOK, http.StatusBadRequest, http.StatusConflict))
	})
}

func TestRequestFailedError_Error(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		require.Equal(t, "status not 200: received 200", NewRequestFailedError(http.StatusOK).Error())
	})
	t.Run("one", func(t *testing.T) {
		require.Equal(t, "status not 200: received 200", NewRequestFailedError(http.StatusOK, http.StatusOK).Error())
	})
	t.Run("two", func(t *testing.T) {
		require.Equal(t, "status not 400, 409: received 200", NewRequestFailedError(http.StatusOK, http.StatusBadRequest, http.StatusConflict).Error())
	})
}
