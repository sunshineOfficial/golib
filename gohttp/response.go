package gohttp

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

func WriteResponseJson(w http.ResponseWriter, status int, a any) error {
	if a == nil {
		return nil
	}

	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(a)
}

func ReadResponseJson(r *http.Response, a any) error {
	if a == nil {
		return nil
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	if err := json.NewDecoder(r.Body).Decode(a); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func WriteResponseXml(w http.ResponseWriter, status int, a any) error {
	if a == nil {
		return nil
	}

	w.Header().Set(ContentTypeHeader, ContentTypeXML)
	w.WriteHeader(status)
	return xml.NewEncoder(w).Encode(a)
}

func ReadResponseXml(r *http.Response, a any) error {
	if r.Body != nil {
		defer r.Body.Close()
	}

	if a == nil {
		return nil
	}
	return xml.NewDecoder(r.Body).Decode(a)
}
