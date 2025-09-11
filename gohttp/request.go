package gohttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/sunshineOfficial/golib/goctx"
)

func NewRequest(ctx goctx.Context, method string, url string, body io.Reader) (*http.Request, error) {
	rq, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	return RequestWithContext(rq, ctx), nil
}

func WriteRequestJson(r *http.Request, a any) (err error) {
	if a == nil {
		return nil
	}

	r.Header.Set(ContentTypeHeader, ContentTypeJSON)

	r.Body, err = newReadCloser(func(b io.Writer, a any) error {
		return json.NewEncoder(b).Encode(a)
	}, a)
	return err
}

func ReadRequestJson(r *http.Request, a any) error {
	if a == nil {
		return nil
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(a)
}

func WriteRequestXml(r *http.Request, a any) (err error) {
	if a == nil {
		return nil
	}

	r.Header.Set(ContentTypeHeader, ContentTypeXML)

	r.Body, err = newReadCloser(func(b io.Writer, a any) error {
		return xml.NewEncoder(b).Encode(a)
	}, a)
	return err
}

func ReadRequestXml(r *http.Request, a any) error {
	if a == nil {
		return nil
	}

	defer r.Body.Close()
	return xml.NewDecoder(r.Body).Decode(a)
}

func newReadCloser(e func(b io.Writer, a any) error, a any) (io.ReadCloser, error) {
	body := bytes.NewBuffer(nil)
	if err := e(body, a); err != nil {
		return nil, err
	}

	return io.NopCloser(body), nil
}
