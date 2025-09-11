package cache

import (
	"time"

	"github.com/sunshineOfficial/golib/goctx"
)

type Cache[K comparable, V any] interface {
	// Get возвращает значение и nil, если запись с указанным ключом существует в кэше. zero-value и error в остальных случаях
	Get(goctx.Context, K) (V, error)

	// GetOrAdd возвращает значение и nil, если запись с указанным ключом существует в кэше.
	// Если записи по ключу не найдена, то метод вызывает функцию для получения значения извне и сохраняет ее в кэш по этому ключу.
	GetOrAdd(goctx.Context, K, NewValueFunc[V]) (V, error)

	// Set сохраняет значение по ключу в кэш. Возвращает ошибку, если сохранить значение не удалось, nil - в остальных случаях
	Set(goctx.Context, K, V, ...SetOption) error

	// Delete удаляет указанный ключ и его значение из кэша. Возвращает ошибку, если удаление выполнить не удалось, nil - в остальных случаях.
	Delete(goctx.Context, K) error
}

type NewValueFunc[V any] func() (V, []SetOption, error)

type Options struct {
	DefaultTTL time.Duration
	ErrorTTL   time.Duration
}

type Option interface {
	Apply(options Options) Options
}

type OptionFunc func(options Options) Options

func (c OptionFunc) Apply(options Options) Options {
	return c(options)
}

// WithDefaultTTL устанавливает TTL на все записи кэша без ошибок
func WithDefaultTTL(ttl time.Duration) Option {
	return OptionFunc(func(options Options) Options {
		options.DefaultTTL = ttl
		return options
	})
}

// WithErrorTTL устанавливает TTL на все записи кэша с ошибками
func WithErrorTTL(ttl time.Duration) Option {
	return OptionFunc(func(options Options) Options {
		options.ErrorTTL = ttl
		return options
	})
}

type SetOptions struct {
	TTL   time.Duration
	Error error
}

type SetOption interface {
	Apply(options SetOptions) SetOptions
}

type SetOptionFunc func(options SetOptions) SetOptions

func (c SetOptionFunc) Apply(options SetOptions) SetOptions {
	return c(options)
}

// WithTTL
// Переопределяет TTL для выбранной записи
func WithTTL(ttl time.Duration) SetOption {
	return SetOptionFunc(func(options SetOptions) SetOptions {
		options.TTL = ttl
		return options
	})
}

func WithError(err error) SetOption {
	return SetOptionFunc(func(options SetOptions) SetOptions {
		options.Error = err
		return options
	})
}

func applyAllOptions(options ...Option) Options {
	var result Options
	for _, o := range options {
		result = o.Apply(result)
	}

	if result.DefaultTTL > -1 && result.ErrorTTL < 0 {
		result.ErrorTTL = result.DefaultTTL
	}

	return result
}

func applyAllSetOptions(options ...SetOption) SetOptions {
	result := SetOptions{
		TTL: -1,
	}

	for _, o := range options {
		result = o.Apply(result)
	}

	return result
}
