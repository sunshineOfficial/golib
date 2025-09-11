package gosync

import "sync"

type Pool[T any] struct {
	pool   *sync.Pool
	canPut func(T) bool
}

func NewPool[T any](new func() T, canPut func(T) bool) *Pool[T] {
	return &Pool[T]{
		pool: &sync.Pool{
			New: func() any {
				return new()
			},
		},
		canPut: canPut,
	}
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(t T) {
	if p.canPut != nil && !p.canPut(t) {
		return
	}

	p.pool.Put(t)
}
