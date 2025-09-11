package gorouter

import "net/http"

type Handler func(c Context) error

func WrapStdLib(h http.Handler) Handler {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func WrapStdLibFunc(h http.HandlerFunc) Handler {
	return WrapStdLib(h)
}
