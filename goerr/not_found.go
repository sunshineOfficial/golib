package goerr

import (
	"database/sql"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrNotFound ошибка NotFoundError без содержания
	ErrNotFound = NewNotFoundError("", nil)
)

type NotFoundError struct {
	origin  error
	message string
}

func NewNotFoundError(message string, err error) NotFoundError {
	return NotFoundError{
		message: message,
		origin:  err,
	}
}

// WrapNotFoundError
// Заворачивает ошибку в NotFoundError, если ее можно отнести к ошибкам не найденных данных: sql.ErrNoRows, mongo.ErrNoDocuments.
// Может быть полезным при использовании в репозиториях для отвязывания ошибки от природы репозитория:
// ошибки sql.ErrNoRows и mongo.ErrNoDocuments будут выглядеть для бизнес-логики как экземпляры NotFoundError
func WrapNotFoundError(format string, err error) error {
	switch {
	case err == nil:
		return nil

	case errors.Is(err, sql.ErrNoRows),
		errors.Is(err, mongo.ErrNoDocuments):
		return NewNotFoundError(fmt.Sprintf(format, err), err)

	default:
		return err
	}
}

func (n NotFoundError) Error() string {
	if len(n.message) == 0 {
		if n.origin == nil {
			return "not found"
		}

		return n.origin.Error()
	}

	if n.origin != nil {
		return fmt.Sprintf("%s: %v", n.message, n.origin.Error())
	}

	return n.message
}

func (n NotFoundError) Unwrap() error {
	return n.origin
}

// IsNotFound
// Возвращает true, если указанная ошибка является NotFoundError или ErrNotFound
func IsNotFound(err error) bool {
	if ok := errors.Is(err, ErrNotFound); ok {
		return ok
	}

	_, ok := AsError[NotFoundError](err)
	return ok
}
