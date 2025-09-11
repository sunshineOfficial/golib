package gosync

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sunshineOfficial/golib/goerr"
)

func TestWaitContext(t *testing.T) {
	type testCase struct {
		/* sources */
		ctxFn  func() (context.Context, context.CancelFunc)
		waitFn func(cancel context.CancelFunc) error

		/* results */
		error       error
		waitSeconds int
	}

	testCases := map[string]testCase{
		"ctx:background/waitFn:nil/err:nil/wait:0": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			waitFn:      nil,
			error:       nil,
			waitSeconds: 0,
		},
		"ctx:timeout20/waitFn:sleep21/err:DeadlineExceeded/wait:20": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Second*20)
			},
			waitFn: func(_ context.CancelFunc) error {
				time.Sleep(time.Second * 21)
				return nil
			},
			error:       context.DeadlineExceeded,
			waitSeconds: 20,
		},
		"ctx:timeout15/waitFn:sleep5/err:nil/wait:5": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Second*15)
			},
			waitFn: func(_ context.CancelFunc) error {
				time.Sleep(time.Second * 5)
				return nil
			},
			error:       nil,
			waitSeconds: 5,
		},
		"ctx:cancel/waitFn:cancelNow/err:Canceled/wait:0": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			waitFn: func(cancel context.CancelFunc) error {
				cancel()
				return nil
			},
			error:       context.Canceled,
			waitSeconds: 0,
		},
		"ctx:background/waitFn:passErr/err:passed/wait:0": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			waitFn: func(cancel context.CancelFunc) error {
				return goerr.ErrNotFound
			},
			error:       goerr.ErrNotFound,
			waitSeconds: 0,
		},
	}

	t.Parallel()

	for name, testCaseInfo := range testCases {
		tc := testCaseInfo

		t.Run(name, func(t *testing.T) {
			start := time.Now()

			t.Parallel()

			ctx, cancel := tc.ctxFn()
			defer cancel()

			err := WaitContext(ctx, func() error {
				if tc.waitFn == nil {
					return nil
				}

				return tc.waitFn(cancel)
			})

			elapsedTime := time.Since(start)

			require.ErrorIs(t, err, tc.error)
			assert.Equal(t, tc.waitSeconds, int(elapsedTime.Seconds()))
		})
	}
}

func TestNopWaitFn(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		fn := NopWaitFn(nil)
		assert.Nil(t, fn)
	})
	t.Run("fn", func(t *testing.T) {
		called := &atomic.Bool{}

		fn := NopWaitFn(func() {
			called.Store(true)
		})

		require.NotNil(t, fn)
		assert.NoError(t, fn())
		assert.True(t, called.Load())
	})
}
