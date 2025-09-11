package gorouter

import (
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunshineOfficial/golib/golog"
)

func TestRouter_Add(t *testing.T) {
	var (
		log = golog.NewLogger("gorouter")
	)

	t.Run("order/good", func(t *testing.T) {
		r := NewRouter(log)

		var (
			calledCount = &atomic.Int32{}
			firstMwId   = &atomic.Int32{}
			secondMwId  = &atomic.Int32{}
			thirdMwId   = &atomic.Int32{}
		)

		r.Add(func(c Context, h Handler) Handler {
			calledCount.Add(1)
			firstMwId.Store(calledCount.Load())

			return h
		}, func(c Context, h Handler) Handler {
			calledCount.Add(1)
			secondMwId.Store(calledCount.Load())

			return h
		}, func(c Context, h Handler) Handler {
			calledCount.Add(1)
			thirdMwId.Store(calledCount.Load())

			return h
		})

		r.Handle("/", func(c Context) error {
			return nil
		})

		var (
			responseWriter = httptest.NewRecorder()
			request        = httptest.NewRequest("GET", "/", nil)
		)

		r.ServeHTTP(responseWriter, request)

		assert.EqualValues(t, 3, calledCount.Load())
		assert.EqualValues(t, 1, thirdMwId.Load())
		assert.EqualValues(t, 2, secondMwId.Load())
		assert.EqualValues(t, 3, firstMwId.Load())
	})

}
