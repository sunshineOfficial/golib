package goerr

import (
	"errors"
)

// compliance test
var _ error = &ServiceError{}

type ServiceError struct {
	internal      error
	code, message string
}

func NewServiceError(internal error, code, message string) ServiceError {
	return ServiceError{
		internal: internal,
		code:     code,
		message:  message,
	}
}

func (s ServiceError) Internal() error {
	return s.internal
}

func (s ServiceError) Error() string {
	return s.internal.Error()
}

func (s ServiceError) Unwrap() error {
	return s.internal
}

func (s ServiceError) Info() (code, message string) {
	return s.code, s.message
}

func AsServiceError(err error) *ServiceError {
	if err == nil {
		return nil
	}

	var serviceError ServiceError
	if errors.As(err, &serviceError) {
		return &serviceError
	}

	return nil
}
