package gohttp

import (
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

func AddIntQuery(url string, name string, values ...int) (string, error) {
	return AddQuery(url, name, values, strconv.Itoa)
}

func AddStringQuery(url string, name string, values ...string) (string, error) {
	return AddQuery(url, name, values, func(v string) string {
		return v
	})
}

func AddUUIDQuery(url, name string, values ...uuid.UUID) (string, error) {
	return AddQuery(url, name, values, func(v uuid.UUID) string {
		return v.String()
	})
}

func AddQuery[T any](baseURL, name string, values []T, toString func(v T) string) (string, error) {
	parsedUrl, err := url.Parse(baseURL)
	if err != nil {
		return baseURL, err
	}

	query := parsedUrl.Query()
	for _, value := range values {
		query.Add(name, toString(value))
	}

	parsedUrl.RawQuery = query.Encode()
	return parsedUrl.String(), nil
}
