package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/goerr"
)

type entry[V any] struct {
	Value V      `json:"value"`
	Error string `json:"error"`
}

func (e entry[V]) getError() error {
	if len(e.Error) > 0 {
		return errors.New(e.Error)
	}

	return nil
}

// compliance test
var _ Cache[int, string] = NewRedis[int, string](nil, "")

// Redis кэш в redis
// Поддерживаемые опции: WithDefaultTTL
type Redis[K, V any] struct {
	client               *redis.Client
	keyFormat            string
	defaultTTL, errorTTL time.Duration
	localMX              *sync.RWMutex
}

func NewRedis[K, T any](client *redis.Client, keyFormat string, options ...Option) *Redis[K, T] {
	opt := applyAllOptions(options...)

	return &Redis[K, T]{
		client:     client,
		keyFormat:  keyFormat,
		defaultTTL: opt.DefaultTTL,
		errorTTL:   opt.ErrorTTL,
		localMX:    &sync.RWMutex{},
	}
}

func (r *Redis[K, V]) Get(ctx goctx.Context, key K) (V, error) {
	r.localMX.RLock()
	defer r.localMX.RUnlock()

	var e entry[V]

	bytes, err := r.client.Get(ctx, fmt.Sprintf(r.keyFormat, key)).Bytes()
	switch {
	case errors.Is(err, redis.Nil):
		return e.Value, goerr.ErrNotFound

	case err != nil:
		return e.Value, err
	}

	err = json.Unmarshal(bytes, &e)
	if err != nil {
		return e.Value, err
	}

	return e.Value, e.getError()
}

func (r *Redis[K, V]) GetOrAdd(ctx goctx.Context, key K, newValue NewValueFunc[V]) (V, error) {
	value, err := r.Get(ctx, key)
	switch {
	case err == nil:
		return value, err
	case !goerr.IsNotFound(err):
		return value, err
	}

	var options []SetOption
	value, options, err = newValue()
	if err != nil {
		return value, err
	}

	return value, r.Set(ctx, key, value, options...)
}

func (r *Redis[K, V]) Set(ctx goctx.Context, key K, value V, options ...SetOption) error {
	r.localMX.Lock()
	defer r.localMX.Unlock()

	var (
		opt = applyAllSetOptions(options...)
		ttl = r.defaultTTL
		e   = entry[V]{
			Value: value,
		}
	)

	if opt.Error != nil {
		e.Error = opt.Error.Error()
		ttl = r.errorTTL
	}

	if opt.TTL > -1 {
		ttl = opt.TTL
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, fmt.Sprintf(r.keyFormat, key), bytes, ttl).Err()
}

func (r *Redis[K, V]) Delete(ctx goctx.Context, key K) error {
	r.localMX.Lock()
	defer r.localMX.Unlock()

	return r.client.Del(ctx, fmt.Sprintf(r.keyFormat, key)).Err()
}
