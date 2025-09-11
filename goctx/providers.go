package goctx

import (
	"context"
	"time"
)

// ProvideWithCancel возвращает контекст с функцией отмены
type ProvideWithCancel func() (context.Context, context.CancelFunc)

// ProvideWithTimeout возвращает контекст с таймаутом
type ProvideWithTimeout func(timeout time.Duration) (context.Context, context.CancelFunc)

// ProvideWithDeadline возвращает контекст с датой истечения (дедлайн)
type ProvideWithDeadline func(deadline time.Time) (context.Context, context.CancelFunc)
