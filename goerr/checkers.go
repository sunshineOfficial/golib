package goerr

import (
	"context"
	"errors"
)

// AsError возвращает ошибку в указаном типа или nil, если привести её к этом типу не удалось.
// Является оберткой для сокращения кода при использовании стандартного метода errors.As.
func AsError[E error](err error) (E, bool) {
	var castError E
	return castError, errors.As(err, &castError)
}

// IsContextError
// Проверяет является ли ошибка связанной с работой контекстов context.Context / goctx.Context
func IsContextError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}

// IsCriticalError
// Проверяет является ли ошибка критичной (обязательной для логирования, например).
// Под критичные ошибки не попадают: предупреждения Warning, ошибки NotFoundError и ошибки контекста, проверяемые через IsContextError
func IsCriticalError(err error) bool {
	_, isWarn := AsError[Warning](err)
	if isWarn {
		return false
	}
	_, isNotFound := AsError[NotFoundError](err)
	if isNotFound {
		return false
	}

	return !IsContextError(err)
}
