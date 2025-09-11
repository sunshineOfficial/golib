package gohttp

import "net/http"

func IsOk(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

func RequireStatus(actual int, statuses ...int) error {
	for _, status := range statuses {
		if actual == status {
			return nil
		}
	}

	return NewRequestFailedError(actual, statuses...)
}

func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func Panic(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
