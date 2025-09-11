package golog

import (
	"io"
	"os"
)

type optionHolder struct {
	out        io.Writer
	err        io.Writer
	tags       []Tag
	skip       int
	stacktrace bool
	global     bool
}

type Option interface {
	apply(o optionHolder) optionHolder
}

type optionFunc func(o optionHolder) optionHolder

func (f optionFunc) apply(o optionHolder) optionHolder {
	return f(o)
}

func WithStdOut() Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.out = os.Stdout
		return o
	})
}

func WithStdErr() Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.err = os.Stderr
		return o
	})
}

// WithOut устанавливает указанный io.Writer как вывод для обычных сообщений
func WithOut(out io.Writer) Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.out = out
		return o
	})
}

// WithErr устанавливает указанный io.Writer как вывод для ошибок
func WithErr(err io.Writer) Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.err = err
		return o
	})
}

// WithWriter устанавливает указанный io.Writer как вывод для всех сообщений
func WithWriter(w io.Writer) Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.out = w
		o.err = w
		return o
	})
}

// WithTags устанавливает теги по умолчанию. Они будут добавлены во все записи логгера автоматически
func WithTags(tags ...Tag) Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.tags = make([]Tag, len(tags))
		copy(o.tags, tags)

		return o
	})
}

// WithSkip позволяет увеличить количество пропускаемых вызовов для определения места, где выполнена запись лога.
// Может быть полезно, если Вы заворачиваете логгер в свою структуру
func WithSkip(skip int) Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.skip = skip
		return o
	})
}

// WithStacktrace добавляет стектрейс в лог с ошибками
func WithStacktrace() Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.stacktrace = true
		return o
	})
}

func MakeGlobal() Option {
	return optionFunc(func(o optionHolder) optionHolder {
		o.global = true
		return o
	})
}
