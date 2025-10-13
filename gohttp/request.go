package gohttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"

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

func WriteRequestMultipart(r *http.Request, data *MultipartData) (err error) {
	if data == nil {
		return nil
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, field := range data.Fields {
		if err = writer.WriteField(field.Name, field.Value); err != nil {
			err = fmt.Errorf("write field: %w", err)
			return closeWriter(writer, err)
		}
	}

	for _, file := range data.Files {
		headers := make(textproto.MIMEHeader)
		headers.Set("Content-Disposition", multipart.FileContentDisposition(file.FieldName, file.FileName))
		headers.Set("Content-Type", ContentTypeByExtension(file.FileName))

		part, createErr := writer.CreatePart(headers)
		if createErr != nil {
			err = fmt.Errorf("create form file: %w", createErr)
			return closeWriter(writer, err)
		}

		_, err = io.Copy(part, file.Reader)
		if err != nil {
			err = fmt.Errorf("copy form file: %w", err)
			return closeWriter(writer, err)
		}
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("close writer: %w", err)
	}

	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.Body = io.NopCloser(body)
	r.ContentLength = int64(body.Len())

	return err
}

func ContentTypeByExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "application/octet-stream"
	}

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}

func closeWriter(writer *multipart.Writer, err error) error {
	closeErr := writer.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("close writer: %w", closeErr)
		return errors.Join(err, closeErr)
	}

	return err
}

func newReadCloser(e func(b io.Writer, a any) error, a any) (io.ReadCloser, error) {
	body := bytes.NewBuffer(nil)
	if err := e(body, a); err != nil {
		return nil, err
	}

	return io.NopCloser(body), nil
}
