package gosync

import "context"

type WaitFn func() error

// WaitContext ожидает выполнения указанной функции fn пока не завершится контекст ctx
func WaitContext(ctx context.Context, fn WaitFn) error {
	c := make(chan error, 1)

	go func(c chan error, fn func() error) {
		c <- fn()
		close(c)
	}(c, fn)

	select {
	case _ = <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func NopWaitFn(fn func()) WaitFn {
	if fn == nil {
		return nil
	}

	return func() error {
		fn()
		return nil
	}
}
