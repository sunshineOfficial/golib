package gohttp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrContentTypeRequired = fmt.Errorf("не найден заголовок %s", ContentTypeHeader)

	ErrHeadersRequired = errors.New("заголовки отсутствуют")
)

type RequestFailedError struct {
	expected string
	actual   int
}

func NewRequestFailedError(actual int, expectedCodes ...int) RequestFailedError {
	expected := "200"

	if len(expectedCodes) > 0 {
		temp := make([]string, 0, len(expected))
		for _, code := range expectedCodes {
			temp = append(temp, strconv.Itoa(code))
		}
		expected = strings.Join(temp, ", ")
	}

	return RequestFailedError{
		expected: expected,
		actual:   actual,
	}
}

func (r RequestFailedError) Error() string {
	return fmt.Sprintf("status not %s: received %d", r.expected, r.actual)
}
