package reflect

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrStructOrMapRequired = errors.New("поддерживается чтение данных только в структуры и map[string]string")
	ErrNoAnonymous         = errors.New("чтение данных в анонимные типы не поддерживается")
)

type UnsupportedTypeError struct {
	message string
}

func NewUnsupportedTypeError(kind reflect.Kind) UnsupportedTypeError {
	return UnsupportedTypeError{
		message: fmt.Sprintf(`тип "%s" не поддерживается`, kind),
	}
}

func (e UnsupportedTypeError) Error() string {
	return e.message
}

type RequiredFieldError struct {
	message string
}

func NewRequiredFieldError(tagType TagType, name string) RequiredFieldError {
	return RequiredFieldError{
		message: fmt.Sprintf(`обязательный параметр %s:%s не найден`, tagType, name),
	}
}

func (e RequiredFieldError) Error() string {
	return e.message
}
