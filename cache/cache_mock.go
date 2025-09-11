package cache

import (
	"errors"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/goerr"
)

// compliance test
var _ Cache[int, string] = NewMock[int, string](nil, nil, nil)

// Mock мок для тестирования кода, использующего Cache
type Mock[K comparable, V any] struct {
	get    func(goctx.Context, K) (V, error)
	save   func(goctx.Context, K, V, ...SetOption) error
	delete func(goctx.Context, K) error
}

func NewMock[K comparable, V any](
	get func(goctx.Context, K) (V, error),
	save func(goctx.Context, K, V, ...SetOption) error,
	delete func(goctx.Context, K) error,
) *Mock[K, V] {
	return &Mock[K, V]{
		get:    get,
		save:   save,
		delete: delete,
	}
}

func (e *Mock[K, V]) Get(ctx goctx.Context, key K) (V, error) {
	if e.get != nil {
		return e.get(ctx, key)
	}

	var v V
	return v, nil
}

func (e *Mock[K, V]) GetOrAdd(ctx goctx.Context, key K, newValue NewValueFunc[V]) (V, error) {
	value, err := e.Get(ctx, key)
	switch {
	case err == nil:
		return value, nil
	case !goerr.IsNotFound(err):
		return value, err
	}

	if newValue == nil {
		return value, errors.New("newValue is nil")
	}

	var options []SetOption
	value, options, err = newValue()
	if err != nil {
		return value, err
	}

	return value, e.Set(ctx, key, value, options...)
}

func (e *Mock[K, V]) Set(ctx goctx.Context, key K, value V, options ...SetOption) error {
	if e.save != nil {
		return e.save(ctx, key, value, options...)
	}

	return nil
}

func (e *Mock[K, V]) Delete(ctx goctx.Context, key K) error {
	if e.delete != nil {
		return e.delete(ctx, key)
	}

	return nil
}
